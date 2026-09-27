package responses

import (
	"context"
	"errors"
	"maps"
	"slices"
	"strings"

	"github.com/tidwall/gjson"
)

// ResponseAccumulator optionally collects output text, function-call arguments
// and custom-tool input from one Responses WebSocket turn. Feed it events from
// ResponseConnection.Recv or ResponseLane.Recv; it never reads or owns the
// socket. Use a separate accumulator per lane and Reset before the next turn.
// Its zero value is ready to use. Like a strings.Builder, do not copy an
// accumulator after first use or call its methods concurrently.
type ResponseAccumulator struct {
	bound         bool
	streamID      string
	responseID    string
	terminalEvent string
	output        map[int64]*responseAccumulatedOutput
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
}

// Reset releases only this helper's state, including its lane and turn binding.
// It never closes or changes the connection. Existing snapshots remain usable.
func (a *ResponseAccumulator) Reset() { *a = ResponseAccumulator{} }

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
// RawJSON. Done events and supplied final output override earlier deltas. When
// the terminal snapshot omits output, only this projection retains prior data;
// ResponseConnection.FinalResponse and ResponseLane.FinalResponse are unchanged.
// A different lane or response requires a separate accumulator or Reset.
func (a *ResponseAccumulator) AddEvent(event ResponsesServerEventUnion) error {
	switch event.Type {
	case "response.created", "response.in_progress", "response.completed", "response.failed", "response.incomplete",
		"response.output_item.added", "response.output_item.done",
		"response.content_part.added", "response.content_part.done",
		"response.output_text.delta", "response.output_text.done",
		"response.function_call_arguments.delta", "response.function_call_arguments.done",
		"response.custom_tool_call_input.delta", "response.custom_tool_call_input.done", "error":
	default:
		return nil
	}
	// Every selected WS event carries the same optional stream_id envelope.
	// Recv preserves it verbatim even for forward-compatible fields.
	streamID := gjson.Get(event.RawJSON(), "stream_id").String()
	if a.bound && a.streamID != streamID {
		return errors.New("responses accumulator: event belongs to another lane")
	}
	if a.terminalEvent != "" {
		return errors.New("responses accumulator: turn ended; reset before adding another event")
	}
	if event.Type == "error" {
		return &ResponseProtocolError{Event: event}
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
	a.bound, a.streamID = true, streamID
	if response != nil {
		if response.ID != "" {
			a.responseID = response.ID
		}
		if len(response.Output) > 0 {
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
		a.addItem(e.OutputIndex, e.Item)
	case "response.output_item.done":
		e := &event.OfResponsesServerEventResponseOutputItemWsDone
		a.addItem(e.OutputIndex, e.Item)
	case "response.content_part.added":
		e := &event.OfResponsesServerEventResponseContentPartWsAdded
		if e.Part.Type == "output_text" {
			replaceAccumulatedText(a.text(e.OutputIndex, e.ItemID, e.ContentIndex), e.Part.Text)
		}
	case "response.content_part.done":
		e := &event.OfResponsesServerEventResponseContentPartWsDone
		if e.Part.Type == "output_text" {
			replaceAccumulatedText(a.text(e.OutputIndex, e.ItemID, e.ContentIndex), e.Part.Text)
		}
	case "response.output_text.delta":
		e := &event.OfResponsesServerEventResponseTextWsDelta
		a.text(e.OutputIndex, e.ItemID, e.ContentIndex).WriteString(e.Delta)
	case "response.output_text.done":
		e := &event.OfResponsesServerEventResponseTextWsDone
		replaceAccumulatedText(a.text(e.OutputIndex, e.ItemID, e.ContentIndex), e.Text)
	case "response.function_call_arguments.delta":
		e := &event.OfResponsesServerEventResponseFunctionCallArgumentsWsDelta
		a.item(e.OutputIndex, e.ItemID).arguments.WriteString(e.Delta)
	case "response.function_call_arguments.done":
		e := &event.OfResponsesServerEventResponseFunctionCallArgumentsWsDone
		replaceAccumulatedText(&a.item(e.OutputIndex, e.ItemID).arguments, e.Arguments)
	case "response.custom_tool_call_input.delta":
		e := &event.OfResponsesServerEventResponseCustomToolCallInputWsDelta
		a.item(e.OutputIndex, e.ItemID).input.WriteString(e.Delta)
	case "response.custom_tool_call_input.done":
		e := &event.OfResponsesServerEventResponseCustomToolCallInputWsDone
		replaceAccumulatedText(&a.item(e.OutputIndex, e.ItemID).input, e.Input)
	}
	return nil
}

func (a *ResponseAccumulator) item(index int64, id string) *responseAccumulatedOutput {
	if a.output == nil {
		a.output = make(map[int64]*responseAccumulatedOutput)
	}
	item := a.output[index]
	if item == nil {
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
	item := a.item(index, source.ID)
	item.itemType, item.callID, item.name = source.Type, source.CallID, source.Name
	switch source.Type {
	case "message":
		for contentIndex, part := range source.Content {
			if part.Type == "output_text" {
				replaceAccumulatedText(a.text(index, source.ID, int64(contentIndex)), part.Text)
			}
		}
	case "function_call":
		replaceAccumulatedText(&item.arguments, source.Arguments.OfString)
	case "custom_tool_call":
		replaceAccumulatedText(&item.input, source.Input)
	}
}

func replaceAccumulatedText(dst *strings.Builder, text string) {
	dst.Reset()
	dst.WriteString(text)
}
