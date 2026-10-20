package responses

import (
	"encoding/json"
	"maps"
	"slices"
	"strings"

	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/tidwall/gjson"
)

// ResponseAccumulatorDetails is mutable observed WebSocket progress, not a
// validated server Response. Response contains the last received lifecycle
// response metadata, excluding output, and is nil before a lifecycle event.
// A valid terminal is still required for TerminalEvent to be populated.
type ResponseAccumulatorDetails struct {
	StreamID      string
	ResponseID    string
	TerminalEvent string
	Response      map[string]json.RawMessage
	Output        []ResponseAccumulatedDetails
}

// ResponseAccumulatedDetails preserves actual output fields and sparse content
// and annotation indices. Item excludes content only for message items.
// Content fields exclude annotations when those are an indexed list; an
// omitted or null annotations field remains omitted or null. Unrecognized
// fields and incomplete streamed logprobs are provisional raw JSON, never
// fabricated final-output models. Every map and RawMessage is owned by the caller.
type ResponseAccumulatedDetails struct {
	OutputIndex int64
	Item        map[string]json.RawMessage
	Content     map[int64]map[string]json.RawMessage
	Annotations map[int64]map[int64]json.RawMessage
}

type responseAccumulatedPart struct {
	raw           string
	logprobs      []string
	logprobsValue string
	annotations   map[int64]string
}

func (item *responseAccumulatedOutput) part(index int64) *responseAccumulatedPart {
	if item.parts == nil {
		item.parts = make(map[int64]*responseAccumulatedPart)
	}
	part := item.parts[index]
	if part == nil {
		part = &responseAccumulatedPart{}
		item.parts[index] = part
	}
	return part
}

func (item *responseAccumulatedOutput) replacePart(index int64, raw string) {
	part := item.part(index)
	*part = responseAccumulatedPart{raw: raw}
	part.addLogprobs(gjson.Get(raw, "logprobs").Raw, true)
}

// Append only the received delta, never rebuild a growing prefix on ingestion.
func (part *responseAccumulatedPart) addLogprobs(raw string, replace bool) {
	if raw == respjson.Omitted {
		return
	}
	value := gjson.Parse(raw)
	if replace || !value.IsArray() || part.logprobsValue != "[]" {
		part.logprobs = nil
	}
	if value.IsArray() {
		part.logprobsValue = "[]"
		value.ForEach(func(_, token gjson.Result) bool {
			part.logprobs = append(part.logprobs, strings.Clone(token.Raw))
			return true
		})
	} else {
		part.logprobsValue = strings.Clone(raw)
	}
}

// DetailedSnapshot returns fresh mutable maps of the last observed metadata,
// provisional content and indexed annotations. Like Snapshot, it joins/copies
// all retained data. Read a full snapshot on demand rather than on each delta;
// original events are available directly to the caller for incremental progress.
// Annotation-only events never bind a lane or replace an unrelated item.
// The original Snapshot and FinalResponse contracts remain unchanged.
func (a *ResponseAccumulator) DetailedSnapshot() ResponseAccumulatorDetails {
	result := ResponseAccumulatorDetails{
		StreamID: a.streamID, ResponseID: a.responseID, TerminalEvent: a.terminalEvent,
	}
	if a.responseRaw != "" {
		result.Response = accumulatedFields(a.responseRaw, "output")
	}
	for _, index := range slices.Sorted(maps.Keys(a.output)) {
		item := a.output[index]
		exclude := ""
		if item.itemType == "message" {
			exclude = "content"
		}
		projected := ResponseAccumulatedDetails{
			OutputIndex: index, Item: accumulatedFields(item.raw, exclude),
			Content:     make(map[int64]map[string]json.RawMessage),
			Annotations: make(map[int64]map[int64]json.RawMessage),
		}
		if item.id != "" {
			projected.Item["id"] = accumulatedJSONString(item.id)
		}
		if item.argumentsPresent {
			projected.Item["arguments"] = accumulatedJSONString(item.arguments.String())
		}
		if item.inputPresent {
			projected.Item["input"] = accumulatedJSONString(item.input.String())
		}
		for pos, part := range item.parts {
			originalAnnotations := gjson.Get(part.raw, "annotations")
			exclude := ""
			if originalAnnotations.IsArray() {
				exclude = "annotations"
			}
			fields := accumulatedFields(part.raw, exclude)
			if part.logprobsValue == "[]" {
				fields["logprobs"] = json.RawMessage("[" + strings.Join(part.logprobs, ",") + "]")
			} else if part.logprobsValue != "" {
				fields["logprobs"] = json.RawMessage(part.logprobsValue)
			}
			projected.Content[pos] = fields
			if originalAnnotations.IsArray() || len(part.annotations) > 0 {
				indexed := make(map[int64]json.RawMessage)
				if originalAnnotations.IsArray() {
					var annoIndex int64
					originalAnnotations.ForEach(func(_, value gjson.Result) bool {
						indexed[annoIndex] = json.RawMessage(value.Raw)
						annoIndex++
						return true
					})
				}
				for i, annotation := range part.annotations {
					indexed[i] = json.RawMessage(annotation)
				}
				projected.Annotations[pos] = indexed
			}
		}
		for pos, text := range item.text {
			fields := projected.Content[pos]
			if fields == nil {
				fields = make(map[string]json.RawMessage)
				projected.Content[pos] = fields
			}
			// A newer refusal/unknown part can supersede earlier text for this
			// detailed view without changing the legacy selected-field snapshot.
			if kind := gjson.ParseBytes(fields["type"]); !kind.Exists() || kind.Str == "output_text" {
				fields["text"] = accumulatedJSONString(text.String())
			}
		}
		result.Output = append(result.Output, projected)
	}
	return result
}

// Original JSON avoids precision loss and preserves fields the current model
// cannot represent. Each []byte conversion allocates caller-owned RawMessage.
func accumulatedFields(raw, exclude string) map[string]json.RawMessage {
	fields := make(map[string]json.RawMessage)
	gjson.Parse(raw).ForEach(func(key, value gjson.Result) bool {
		if exclude == "" || key.Str != exclude {
			fields[key.Str] = json.RawMessage(value.Raw)
		}
		return true
	})
	return fields
}

func accumulatedJSONString(value string) json.RawMessage {
	data, err := json.Marshal(value)
	if err != nil {
		// encoding/json supports every Go string, even invalid UTF-8.
		panic(err)
	}
	return data
}
