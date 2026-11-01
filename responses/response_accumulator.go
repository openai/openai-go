package responses

import (
	"context"
	"encoding/json"
	"errors"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/tidwall/gjson"
)

// ResponseAccumulator optionally collects output text, function and MCP arguments
// and custom-tool input from one Responses WebSocket turn. Feed it events from
// ResponseConnection.Recv or ResponseLane.Recv; it never reads or owns the
// socket. Use a separate accumulator per lane and Reset before the next turn.
// MCP argument progress requires an observed matching mcp_call item.
// DetailedSnapshot also collects incremental refusal content; Snapshot's text
// projection does not treat a refusal as output text.
// Its zero value is ready to use. Like a strings.Builder, do not copy an
// accumulator after first use or call its methods concurrently.
type ResponseAccumulator struct {
	bound         bool
	streamID      string
	responseID    string
	terminalEvent string
	responseRaw   string
	output        map[int64]*responseAccumulatedOutput
	// Final records stay separate from the mutable provisional projection.
	finalizedOutput   map[int64]string
	reconstructionErr error
}

// ResponseAccumulatorSnapshot is a partial projection, not a server Response.
// TerminalEvent is empty until a valid response.completed, response.failed or
// response.incomplete arrives. Output may contain partial tool data; it is not
// a request to execute tools. RawJSON remains on the original received events.
type ResponseAccumulatorSnapshot struct {
	StreamID      string
	ResponseID    string
	TerminalEvent string
	Output        []ResponseAccumulatedOutput
}

// ResponseAccumulatedOutput is one output item's selected fields, ordered by
// OutputIndex in a snapshot. Text maps content indices to their current text.
// Fields not covered by this projection remain available on the received event.
type ResponseAccumulatedOutput struct {
	OutputIndex int64
	ItemID      string
	Type        string
	CallID      string
	Name        string
	Arguments   string
	Input       string
	Text        map[int64]string
}

type responseAccumulatedOutput struct {
	id, itemType, callID, name string
	arguments, input           strings.Builder
	text                       map[int64]*strings.Builder
	raw                        string
	parts                      map[int64]*responseAccumulatedPart
	argumentsPresent           bool
	inputPresent               bool
}

// Reset releases only this helper's state, including its lane and turn binding.
// It never closes or changes the connection. Existing snapshots remain usable.
func (a *ResponseAccumulator) Reset() { *a = ResponseAccumulator{} }

// ReconstructedResponse returns an independent typed terminal response. When
// terminal output is omitted or null, it fills output using only previously
// received response.output_item.done records, sorted by output_index. Sparse
// indices become an ordered array; no finalized records produces an empty array.
// Supplied non-null output, including an empty array, remains authoritative.
//
// This optional compatibility view never promotes unfinished deltas. It returns
// an error before a terminal event, or until Reset after a protocol error
// observed before the terminal. Failed and incomplete responses retain their
// status and error details.
// It does not read the socket or change Snapshot, DetailedSnapshot, original
// events, or ResponseConnection/ResponseLane.FinalResponse.
func (a *ResponseAccumulator) ReconstructedResponse() (*Response, error) {
	if a.reconstructionErr != nil {
		return nil, a.reconstructionErr
	}
	if a.terminalEvent == "" {
		return nil, errors.New("responses accumulator: no terminal response collected")
	}
	raw := a.responseRaw
	if output := gjson.Get(raw, "output"); !output.Exists() || output.Type == gjson.Null {
		items := make([]json.RawMessage, 0, len(a.finalizedOutput))
		for _, index := range slices.Sorted(maps.Keys(a.finalizedOutput)) {
			items = append(items, json.RawMessage(a.finalizedOutput[index]))
		}
		output, err := json.Marshal(items)
		if err != nil {
			return nil, err
		}
		fields := accumulatedFields(raw, "")
		fields["output"] = output
		data, err := json.Marshal(fields)
		if err != nil {
			return nil, err
		}
		raw = string(data)
	}
	var response Response
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Snapshot returns an independent view of the data collected so far. A missing
// terminal response or a receive error never becomes a successful completion.
func (a *ResponseAccumulator) Snapshot() ResponseAccumulatorSnapshot {
	result := ResponseAccumulatorSnapshot{
		StreamID: a.streamID, ResponseID: a.responseID, TerminalEvent: a.terminalEvent,
	}
	for _, index := range slices.Sorted(maps.Keys(a.output)) {
		item := a.output[index]
		projected := ResponseAccumulatedOutput{
			OutputIndex: index, ItemID: item.id, Type: item.itemType,
			CallID: item.callID, Name: item.name, Arguments: item.arguments.String(),
			Input: item.input.String(), Text: make(map[int64]string, len(item.text)),
		}
		for contentIndex, text := range item.text {
			projected.Text[contentIndex] = text.String()
		}
		result.Output = append(result.Output, projected)
	}
	return result
}

// OutputText returns current output text in output/content index order.
func (s ResponseAccumulatorSnapshot) OutputText() string {
	var text strings.Builder
	for _, item := range s.Output {
		for _, index := range slices.Sorted(maps.Keys(item.Text)) {
			text.WriteString(item.Text[index])
		}
	}
	return text.String()
}

// AddEvent observes a typed event from Recv without modifying it. The helper
// ignores unknown/unselected events; callers retain the original event and its
// RawJSON. Done events and supplied final output (including an empty array)
// override earlier deltas. When terminal output is omitted or null, the
// projections retain prior data; ReconstructedResponse can recover finalized
// items. ResponseConnection.FinalResponse and ResponseLane.FinalResponse are
// unchanged. A different lane or response
// requires a separate accumulator or Reset. A new nonempty item ID starts a
// fresh projection at that output index; it never appends to another item.
func (a *ResponseAccumulator) AddEvent(event ResponsesServerEventUnion) error {
	// Newly collected annotations only enrich a known item. Historically these
	// events were ignored: they must not bind a lane or replace a different item,
	// even if the API decoder accepted malformed optional fields.
	if event.Type == "response.output_text.annotation.added" {
		e := &event.OfResponsesServerEventResponseOutputTextAnnotationWsAdded
		if !a.bound || a.terminalEvent != "" ||
			!accumulatorIndicesValid(e.JSON.OutputIndex, e.JSON.ContentIndex, e.JSON.AnnotationIndex) ||
			!accumulatorStringsValid(e.JSON.ItemID) || e.JSON.Annotation.Raw() == respjson.Omitted {
			return nil
		}
		envelope, err := decodeResponseWebsocketEvent([]byte(event.RawJSON()))
		item := a.output[e.OutputIndex]
		if err == nil && envelope.streamID == a.streamID && item != nil && item.id == e.ItemID {
			part := item.part(e.ContentIndex)
			if part.annotations == nil {
				part.annotations = make(map[int64]string)
			}
			part.annotations[e.AnnotationIndex] = strings.Clone(e.JSON.Annotation.Raw())
		}
		return nil
	}
	switch event.Type {
	case "response.created", "response.in_progress", "response.completed", "response.failed", "response.incomplete",
		"response.output_item.added", "response.output_item.done",
		"response.content_part.added", "response.content_part.done",
		"response.output_text.delta", "response.output_text.done",
		"response.refusal.delta", "response.refusal.done",
		"response.function_call_arguments.delta", "response.function_call_arguments.done",
		"response.mcp_call_arguments.delta", "response.mcp_call_arguments.done",
		"response.custom_tool_call_input.delta", "response.custom_tool_call_input.done", "error":
	default:
		return nil
	}
	// Bind to the same scope Recv used, including its permissive treatment of
	// missing or invalid stream IDs. Never coerce a numeric ID into a named lane.
	envelope, err := decodeResponseWebsocketEvent([]byte(event.RawJSON()))
	if err != nil {
		return err
	}
	streamID := envelope.streamID
	if a.bound && a.streamID != streamID {
		return errors.New("responses accumulator: event belongs to another lane")
	}
	if a.terminalEvent != "" {
		return errors.New("responses accumulator: turn ended; reset before adding another event")
	}
	if event.Type == "error" {
		err := &ResponseProtocolError{Event: event}
		a.reconstructionErr = err
		return err
	}
	var response *Response
	switch event.Type {
	case "response.created":
		response = &event.OfResponsesServerEventResponseWsCreated.Response
	case "response.in_progress":
		response = &event.OfResponsesServerEventResponseInWsProgress.Response
	case "response.completed", "response.failed", "response.incomplete":
		// Use the same missing/null/invalid terminal policy as FinalResponse.
		var err error
		response, err = finalWebsocketResponse(context.Background(), func(context.Context) (ResponsesServerEventUnion, error) {
			return event, nil
		})
		if err != nil {
			return err
		}
	}
	if response != nil && a.responseID != "" && response.ID != "" && response.ID != a.responseID {
		return errors.New("responses accumulator: event belongs to another response")
	}
	invalid := errors.New("responses accumulator: invalid selected event fields")
	if response != nil {
		if !accumulatorStringsValid(response.JSON.ID) {
			return invalid
		}
		output := response.JSON.Output
		if !output.Valid() && output.Raw() != respjson.Omitted && output.Raw() != respjson.Null {
			return invalid
		}
		for _, item := range response.Output {
			if !accumulatorItemValid(item) {
				return invalid
			}
		}
		a.bound, a.streamID = true, streamID
		a.responseRaw = response.RawJSON()
		if response.ID != "" {
			a.responseID = response.ID
		}
		if response.JSON.Output.Valid() || len(response.Output) > 0 {
			// A supplied response output supersedes the entire earlier projection,
			// including items and content absent from this newer snapshot.
			a.output = nil
			for index, item := range response.Output {
				a.addItem(int64(index), item)
			}
		}
		switch event.Type {
		case "response.completed", "response.failed", "response.incomplete":
			a.terminalEvent = event.Type
		}
		return nil
	}
	switch event.Type {
	case "response.output_item.added":
		e := &event.OfResponsesServerEventResponseOutputItemWsAdded
		if !accumulatorIndicesValid(e.JSON.OutputIndex) || !e.JSON.Item.Valid() || !accumulatorItemValid(e.Item) {
			return invalid
		}
		a.addItem(e.OutputIndex, e.Item)
	case "response.output_item.done":
		e := &event.OfResponsesServerEventResponseOutputItemWsDone
		if !accumulatorIndicesValid(e.JSON.OutputIndex) || !e.JSON.Item.Valid() || !accumulatorItemValid(e.Item) {
			return invalid
		}
		a.addItem(e.OutputIndex, e.Item)
		if a.finalizedOutput == nil {
			a.finalizedOutput = make(map[int64]string)
		}
		a.finalizedOutput[e.OutputIndex] = strings.Clone(e.Item.RawJSON())
	case "response.content_part.added":
		e := &event.OfResponsesServerEventResponseContentPartWsAdded
		if !accumulatorIndicesValid(e.JSON.OutputIndex, e.JSON.ContentIndex) ||
			!accumulatorStringsValid(e.JSON.ItemID) || !e.JSON.Part.Valid() {
			return invalid
		}
		if e.Part.Type == "output_text" {
			if !accumulatorStringsValid(e.Part.JSON.Text) {
				return invalid
			}
			replaceAccumulatedText(a.text(e.OutputIndex, e.ItemID, e.ContentIndex), e.Part.Text)
		}
		if item := a.output[e.OutputIndex]; item != nil && item.id == e.ItemID {
			item.replacePart(e.ContentIndex, e.Part.RawJSON())
		}
	case "response.content_part.done":
		e := &event.OfResponsesServerEventResponseContentPartWsDone
		if !accumulatorIndicesValid(e.JSON.OutputIndex, e.JSON.ContentIndex) ||
			!accumulatorStringsValid(e.JSON.ItemID) || !e.JSON.Part.Valid() {
			return invalid
		}
		if e.Part.Type == "output_text" {
			if !accumulatorStringsValid(e.Part.JSON.Text) {
				return invalid
			}
			replaceAccumulatedText(a.text(e.OutputIndex, e.ItemID, e.ContentIndex), e.Part.Text)
		}
		if item := a.output[e.OutputIndex]; item != nil && item.id == e.ItemID {
			item.replacePart(e.ContentIndex, e.Part.RawJSON())
		}
	case "response.output_text.delta":
		e := &event.OfResponsesServerEventResponseTextWsDelta
		if !accumulatorIndicesValid(e.JSON.OutputIndex, e.JSON.ContentIndex) ||
			!accumulatorStringsValid(e.JSON.ItemID, e.JSON.Delta) {
			return invalid
		}
		a.text(e.OutputIndex, e.ItemID, e.ContentIndex).WriteString(e.Delta)
		a.output[e.OutputIndex].part(e.ContentIndex).addLogprobs(e.JSON.Logprobs.Raw(), false)
	case "response.output_text.done":
		e := &event.OfResponsesServerEventResponseTextWsDone
		if !accumulatorIndicesValid(e.JSON.OutputIndex, e.JSON.ContentIndex) ||
			!accumulatorStringsValid(e.JSON.ItemID, e.JSON.Text) {
			return invalid
		}
		replaceAccumulatedText(a.text(e.OutputIndex, e.ItemID, e.ContentIndex), e.Text)
		a.output[e.OutputIndex].part(e.ContentIndex).addLogprobs(e.JSON.Logprobs.Raw(), true)
	case "response.refusal.delta":
		e := &event.OfResponsesServerEventResponseRefusalWsDelta
		if !accumulatorIndicesValid(e.JSON.OutputIndex, e.JSON.ContentIndex) ||
			!accumulatorStringsValid(e.JSON.ItemID, e.JSON.Delta) {
			return invalid
		}
		// Previously raw refusal parts remain accepted until an update needs
		// their string. Validate that seed before changing an item's identity.
		var seed string
		if item := a.output[e.OutputIndex]; item != nil && (item.id == "" || item.id == e.ItemID) {
			if part := item.parts[e.ContentIndex]; part != nil && !part.refusalPresent &&
				gjson.Get(part.raw, "type").Str == "refusal" {
				value := gjson.Get(part.raw, "refusal")
				if value.Type != gjson.String {
					return invalid
				}
				seed = value.Str
			}
		}
		part := a.item(e.OutputIndex, e.ItemID).part(e.ContentIndex)
		if !part.refusalPresent {
			part.refusal.WriteString(seed)
			part.refusalPresent = true
		}
		part.refusal.WriteString(e.Delta)
	case "response.refusal.done":
		e := &event.OfResponsesServerEventResponseRefusalWsDone
		if !accumulatorIndicesValid(e.JSON.OutputIndex, e.JSON.ContentIndex) ||
			!accumulatorStringsValid(e.JSON.ItemID, e.JSON.Refusal) {
			return invalid
		}
		part := a.item(e.OutputIndex, e.ItemID).part(e.ContentIndex)
		replaceAccumulatedText(&part.refusal, e.Refusal)
		part.refusalPresent = true
	case "response.function_call_arguments.delta":
		e := &event.OfResponsesServerEventResponseFunctionCallArgumentsWsDelta
		if !accumulatorIndicesValid(e.JSON.OutputIndex) || !accumulatorStringsValid(e.JSON.ItemID, e.JSON.Delta) {
			return invalid
		}
		item := a.item(e.OutputIndex, e.ItemID)
		item.arguments.WriteString(e.Delta)
		item.argumentsPresent = true
	case "response.function_call_arguments.done":
		e := &event.OfResponsesServerEventResponseFunctionCallArgumentsWsDone
		if !accumulatorIndicesValid(e.JSON.OutputIndex) || !accumulatorStringsValid(e.JSON.ItemID, e.JSON.Arguments) {
			return invalid
		}
		item := a.item(e.OutputIndex, e.ItemID)
		replaceAccumulatedText(&item.arguments, e.Arguments)
		item.argumentsPresent = true
	case "response.mcp_call_arguments.delta":
		e := &event.OfResponsesServerEventResponseMcpCallArgumentsWsDelta
		if !accumulatorIndicesValid(e.JSON.OutputIndex) || !accumulatorStringsValid(e.JSON.ItemID, e.JSON.Delta) {
			return invalid
		}
		item := a.output[e.OutputIndex]
		if item == nil || item.itemType != "mcp_call" || item.id != e.ItemID {
			return nil
		}
		item.arguments.WriteString(e.Delta)
		item.argumentsPresent = true
	case "response.mcp_call_arguments.done":
		e := &event.OfResponsesServerEventResponseMcpCallArgumentsWsDone
		if !accumulatorIndicesValid(e.JSON.OutputIndex) || !accumulatorStringsValid(e.JSON.ItemID, e.JSON.Arguments) {
			return invalid
		}
		item := a.output[e.OutputIndex]
		if item == nil || item.itemType != "mcp_call" || item.id != e.ItemID {
			return nil
		}
		replaceAccumulatedText(&item.arguments, e.Arguments)
		item.argumentsPresent = true
	case "response.custom_tool_call_input.delta":
		e := &event.OfResponsesServerEventResponseCustomToolCallInputWsDelta
		if !accumulatorIndicesValid(e.JSON.OutputIndex) || !accumulatorStringsValid(e.JSON.ItemID, e.JSON.Delta) {
			return invalid
		}
		item := a.item(e.OutputIndex, e.ItemID)
		item.input.WriteString(e.Delta)
		item.inputPresent = true
	case "response.custom_tool_call_input.done":
		e := &event.OfResponsesServerEventResponseCustomToolCallInputWsDone
		if !accumulatorIndicesValid(e.JSON.OutputIndex) || !accumulatorStringsValid(e.JSON.ItemID, e.JSON.Input) {
			return invalid
		}
		item := a.item(e.OutputIndex, e.ItemID)
		replaceAccumulatedText(&item.input, e.Input)
		item.inputPresent = true
	}
	a.bound, a.streamID = true, streamID
	return nil
}

// Field.Valid alone permits loose conversions (numeric strings, booleans or
// truncated fractions). Require real, nonnegative int64 wire indices; sparse
// indices remain map keys and carry no allocation limit.
func accumulatorIndicesValid(fields ...respjson.Field) bool {
	for _, field := range fields {
		index, err := strconv.ParseInt(field.Raw(), 10, 64)
		if !field.Valid() || err != nil || index < 0 {
			return false
		}
	}
	return true
}

func accumulatorStringsValid(fields ...respjson.Field) bool {
	for _, field := range fields {
		if !field.Valid() || gjson.Parse(field.Raw()).Type != gjson.String {
			return false
		}
	}
	return true
}

// Only validate the fields this projection consumes. Unknown item and content
// kinds remain available on the raw event; tool IDs are optional in the wire
// models, but an ID that was supplied must not be coerced to a string.
func accumulatorItemValid(item ResponseOutputItemUnion) bool {
	if !accumulatorStringsValid(item.JSON.Type) ||
		(item.JSON.ID.Raw() != respjson.Omitted && !accumulatorStringsValid(item.JSON.ID)) {
		return false
	}
	switch item.Type {
	case "message":
		if !item.JSON.Content.Valid() {
			return false
		}
		for _, part := range item.Content {
			if !accumulatorStringsValid(part.JSON.Type) ||
				(part.Type == "output_text" && !accumulatorStringsValid(part.JSON.Text)) {
				return false
			}
		}
	case "function_call":
		return accumulatorStringsValid(item.JSON.CallID, item.JSON.Name, item.JSON.Arguments)
	case "custom_tool_call":
		return accumulatorStringsValid(item.JSON.CallID, item.JSON.Name, item.JSON.Input)
	}
	return true
}

func (a *ResponseAccumulator) item(index int64, id string) *responseAccumulatedOutput {
	if a.output == nil {
		a.output = make(map[int64]*responseAccumulatedOutput)
	}
	item := a.output[index]
	if item == nil || (id != "" && item.id != "" && id != item.id) {
		item = &responseAccumulatedOutput{text: make(map[int64]*strings.Builder)}
		a.output[index] = item
	}
	if id != "" {
		item.id = id
	}
	return item
}

func (a *ResponseAccumulator) text(index int64, id string, contentIndex int64) *strings.Builder {
	item := a.item(index, id)
	text := item.text[contentIndex]
	if text == nil {
		text = &strings.Builder{}
		item.text[contentIndex] = text
	}
	return text
}

func (a *ResponseAccumulator) addItem(index int64, source ResponseOutputItemUnion) {
	// Full item snapshots replace the entire projection, even when the new
	// content is empty or the corrected item has a different type.
	delete(a.output, index)
	item := a.item(index, source.ID)
	item.itemType, item.callID, item.name = source.Type, source.CallID, source.Name
	item.raw = source.RawJSON()
	switch source.Type {
	case "message":
		for contentIndex, part := range source.Content {
			item.replacePart(int64(contentIndex), part.RawJSON())
			if part.Type == "output_text" {
				replaceAccumulatedText(a.text(index, source.ID, int64(contentIndex)), part.Text)
			}
		}
	case "function_call":
		replaceAccumulatedText(&item.arguments, source.Arguments.OfString)
	case "mcp_call":
		// Preserve previously accepted raw MCP items without coercing non-string arguments.
		if arguments := gjson.Parse(source.JSON.Arguments.Raw()); arguments.Type == gjson.String {
			replaceAccumulatedText(&item.arguments, arguments.Str)
			item.argumentsPresent = true
		}
	case "custom_tool_call":
		replaceAccumulatedText(&item.input, source.Input)
	}
}

func replaceAccumulatedText(dst *strings.Builder, text string) {
	dst.Reset()
	dst.WriteString(text)
}
