package responses

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// Buffers belong to the existing item projection and are discarded with a full
// item replacement. They are materialized only when a caller requests details.
type responseAccumulatedProgress struct {
	code            *strings.Builder
	status          string
	summary         map[int64]*responseAccumulatedPart
	summaryChanged  bool
	commands        map[int64]*strings.Builder
	commandsChanged bool
	shellOutput     map[int64]*responseAccumulatedShellOutput
}

type responseAccumulatedShellOutput struct {
	doneOutput                   string
	stdout, stderr               strings.Builder
	stdoutPresent, stderrPresent bool
}

func accumulatorRichEvent(kind string) bool {
	switch kind {
	case "response.reasoning_summary_part.added", "response.reasoning_summary_part.done",
		"response.reasoning_summary_text.delta", "response.reasoning_summary_text.done",
		"response.reasoning_text.delta", "response.reasoning_text.done",
		"response.shell_call_command.added", "response.shell_call_command.delta", "response.shell_call_command.done",
		"response.shell_call_output_content.delta", "response.shell_call_output_content.done",
		"response.code_interpreter_call_code.delta", "response.code_interpreter_call_code.done",
		"response.code_interpreter_call.in_progress", "response.code_interpreter_call.interpreting", "response.code_interpreter_call.completed",
		"response.file_search_call.in_progress", "response.file_search_call.searching", "response.file_search_call.completed",
		"response.web_search_call.in_progress", "response.web_search_call.searching", "response.web_search_call.completed",
		"response.image_generation_call.in_progress", "response.image_generation_call.generating", "response.image_generation_call.completed",
		"response.mcp_call.in_progress", "response.mcp_call.completed", "response.mcp_call.failed":
		return true
	}
	return false
}

func accumulatedWireIndex(value gjson.Result) (int64, bool) {
	index, err := strconv.ParseInt(value.Raw, 10, 64)
	return index, value.Exists() && err == nil && index >= 0
}

func accumulatedPart(raw string) *responseAccumulatedPart {
	return &responseAccumulatedPart{raw: strings.Clone(raw)}
}

func (a *ResponseAccumulator) addRichEvent(event ResponsesServerEventUnion) error {
	invalid := errors.New("responses accumulator: invalid selected event fields")
	raw := gjson.Parse(event.RawJSON())
	index, ok := accumulatedWireIndex(raw.Get("output_index"))
	if !ok {
		return invalid
	}
	kind := strings.TrimPrefix(event.Type, "response.")
	group, _, _ := strings.Cut(kind, ".")
	expected := group
	switch group {
	case "reasoning_summary_part", "reasoning_summary_text", "reasoning_text":
		expected = "reasoning"
	case "shell_call_command":
		expected = "shell_call"
	case "shell_call_output_content":
		expected = "shell_call_output"
	case "code_interpreter_call_code":
		expected = "code_interpreter_call"
	}
	id := raw.Get("item_id")
	guardID := group != "shell_call_command"
	if guardID && id.Type != gjson.String {
		return invalid
	}
	item := a.output[index]
	if item == nil || item.itemType != expected || (guardID && item.id != id.Str) {
		return nil
	}
	// Only selected fields below are interpreted; all other fields stay on raw events.
	var position int64
	if group == "reasoning_summary_part" || group == "reasoning_summary_text" {
		position, ok = accumulatedWireIndex(raw.Get("summary_index"))
		if !ok {
			return invalid
		}
	} else if group == "reasoning_text" {
		position, ok = accumulatedWireIndex(raw.Get("content_index"))
		if !ok {
			return invalid
		}
	} else if group == "shell_call_command" || group == "shell_call_output_content" {
		position, ok = accumulatedWireIndex(raw.Get("command_index"))
		if !ok {
			return invalid
		}
	}
	if item.progress == nil {
		item.progress = &responseAccumulatedProgress{}
	}
	progress := item.progress
	switch group {
	case "reasoning_summary_part", "reasoning_summary_text":
		if progress.summary == nil {
			source := gjson.Get(item.raw, "summary")
			if source.Exists() && !source.IsArray() {
				return invalid
			}
			summary := make(map[int64]*responseAccumulatedPart)
			for i, part := range source.Array() {
				if !part.IsObject() {
					return invalid
				}
				summary[int64(i)] = accumulatedPart(part.Raw)
			}
			progress.summary = summary
		}
		if group == "reasoning_summary_part" {
			part := raw.Get("part")
			if !part.IsObject() || part.Get("type").Str != "summary_text" || part.Get("text").Type != gjson.String {
				return invalid
			}
			partCount := int64(len(progress.summary))
			if position > partCount || (strings.HasSuffix(kind, ".added") && position != partCount) {
				return invalid
			}
			if strings.HasSuffix(kind, ".done") && progress.summary[position] == nil {
				return nil
			}
			progress.summary[position] = accumulatedPart(part.Raw)
		} else {
			part := progress.summary[position]
			if part == nil {
				return nil
			}
			field := "delta"
			if strings.HasSuffix(kind, ".done") {
				field = "text"
			}
			value := raw.Get(field)
			if value.Type != gjson.String {
				return invalid
			}
			if err := appendReasoningText(part, value.Str, field == "text"); err != nil {
				return err
			}
		}
		progress.summaryChanged = true
	case "reasoning_text":
		part := item.parts[position]
		if part == nil {
			return nil
		}
		field := "delta"
		if strings.HasSuffix(kind, ".done") {
			field = "text"
		}
		value := raw.Get(field)
		if value.Type != gjson.String {
			return invalid
		}
		return appendReasoningText(part, value.Str, field == "text")
	case "code_interpreter_call_code":
		field := "delta"
		if strings.HasSuffix(kind, ".done") {
			field = "code"
		}
		value := raw.Get(field)
		if value.Type != gjson.String {
			return invalid
		}
		if progress.code == nil {
			seed := gjson.Get(item.raw, "code")
			if seed.Exists() && seed.Type != gjson.Null && seed.Type != gjson.String {
				return invalid
			}
			progress.code = &strings.Builder{}
			progress.code.WriteString(seed.Str)
		}
		if field == "code" {
			replaceAccumulatedText(progress.code, value.Str)
		} else {
			progress.code.WriteString(value.Str)
		}
	case "shell_call_command":
		field := "command"
		if strings.HasSuffix(kind, ".delta") {
			field = "delta"
		}
		value := raw.Get(field)
		if value.Type != gjson.String {
			return invalid
		}
		if progress.commands == nil {
			source := gjson.Get(item.raw, "action.commands")
			if !source.IsArray() {
				return invalid
			}
			commands := make(map[int64]*strings.Builder)
			for i, command := range source.Array() {
				if command.Type != gjson.String {
					return invalid
				}
				buffer := &strings.Builder{}
				buffer.WriteString(command.Str)
				commands[int64(i)] = buffer
			}
			progress.commands = commands
		}
		if strings.HasSuffix(kind, ".added") {
			if position != int64(len(progress.commands)) {
				return invalid
			}
			progress.commands[position] = &strings.Builder{}
			a.setShellCommandBound(item, int64(len(progress.commands)))
		}
		command := progress.commands[position]
		if command == nil {
			return nil
		}
		progress.commandsChanged = true
		if field == "command" {
			replaceAccumulatedText(command, value.Str)
		} else {
			command.WriteString(value.Str)
		}
	case "shell_call_output_content":
		if !a.shellOutputIndexValid(item, position) {
			return invalid
		}
		if strings.HasSuffix(kind, ".done") {
			values := raw.Get("output")
			if !values.IsArray() {
				return invalid
			}
			if progress.shellOutput == nil {
				progress.shellOutput = make(map[int64]*responseAccumulatedShellOutput)
			}
			progress.shellOutput[position] = &responseAccumulatedShellOutput{doneOutput: strings.Clone(values.Raw)}
		} else {
			delta := raw.Get("delta")
			if !delta.IsObject() {
				return invalid
			}
			stdout, stderr := delta.Get("stdout"), delta.Get("stderr")
			if (stdout.Exists() && stdout.Type != gjson.String) || (stderr.Exists() && stderr.Type != gjson.String) {
				return invalid
			}
			if progress.shellOutput == nil {
				progress.shellOutput = make(map[int64]*responseAccumulatedShellOutput)
			}
			output := progress.shellOutput[position]
			if output != nil && output.doneOutput != "" {
				return nil
			}
			if output == nil {
				output = &responseAccumulatedShellOutput{}
				progress.shellOutput[position] = output
			}
			output.stdout.WriteString(stdout.Str)
			output.stderr.WriteString(stderr.Str)
			output.stdoutPresent = output.stdoutPresent || stdout.Exists()
			output.stderrPresent = output.stderrPresent || stderr.Exists()
		}
	default:
		_, status, _ := strings.Cut(kind, ".")
		progress.status = status
	}
	return nil
}

func (a *ResponseAccumulator) shellOutputIndexValid(item *responseAccumulatedOutput, index int64) bool {
	if bound, known := a.shellCommandBounds[item.callID]; known {
		return index < bound
	}
	return true
}

func (a *ResponseAccumulator) setShellCommandBound(item *responseAccumulatedOutput, count int64) {
	if item.callID == "" {
		return
	}
	if a.shellCommandBounds == nil {
		a.shellCommandBounds = make(map[string]int64)
	}
	a.shellCommandBounds[item.callID] = count
}

func (a *ResponseAccumulator) forgetShellCommandBound(item *responseAccumulatedOutput) {
	if item != nil && item.itemType == "shell_call" {
		delete(a.shellCommandBounds, item.callID)
	}
}

func appendReasoningText(part *responseAccumulatedPart, value string, replace bool) error {
	if !part.textPresent {
		seed := gjson.Parse(part.raw)
		if kind := seed.Get("type").Str; kind != "summary_text" && kind != "reasoning_text" {
			return nil
		}
		original := seed.Get("text")
		if original.Type != gjson.String {
			return errors.New("responses accumulator: invalid selected event fields")
		}
		part.text.WriteString(original.Str)
		part.textPresent = true
	}
	if replace {
		replaceAccumulatedText(&part.text, value)
	} else {
		part.text.WriteString(value)
	}
	return nil
}

func accumulatedJSON(value any) json.RawMessage {
	data, _ := json.Marshal(value) //nolint:errchkjson // Callers supply only strings and decoder-validated JSON containers.
	return data
}

func (item *responseAccumulatedOutput) renderProgress(projected *ResponseAccumulatedDetails) {
	// Keep the pre-existing non-message content view aligned with the indexed view.
	// A sparse observed part stays indexed; no missing array entries are invented.
	if item.itemType == "reasoning" && len(projected.Content) > 0 {
		values := gjson.Get(item.raw, "content").Array()
		content := make([]json.RawMessage, 0, len(values))
		for i, value := range values {
			row := json.RawMessage(value.Raw)
			if part, ok := projected.Content[int64(i)]; ok {
				row = accumulatedJSON(part)
			}
			content = append(content, row)
		}
		for {
			part, ok := projected.Content[int64(len(content))]
			if !ok {
				break
			}
			content = append(content, accumulatedJSON(part))
		}
		if len(content) > 0 {
			projected.Item["content"] = accumulatedJSON(content)
		}
	}
	progress := item.progress
	if progress == nil {
		return
	}
	if progress.status != "" {
		projected.Item["status"] = accumulatedJSONString(progress.status)
	}
	if progress.code != nil {
		projected.Item["code"] = accumulatedJSONString(progress.code.String())
	}
	if progress.summaryChanged {
		summary := make([]json.RawMessage, 0, len(progress.summary))
		for index := int64(0); index < int64(len(progress.summary)); index++ {
			part := progress.summary[index]
			fields := accumulatedFields(part.raw, "")
			if part.textPresent {
				fields["text"] = accumulatedJSONString(part.text.String())
			}
			summary = append(summary, accumulatedJSON(fields))
		}
		projected.Item["summary"] = accumulatedJSON(summary)
	}
	if progress.commandsChanged {
		commands := make([]string, len(progress.commands))
		for i := range commands {
			commands[i] = progress.commands[int64(i)].String()
		}
		action := accumulatedFields(gjson.Get(item.raw, "action").Raw, "")
		action["commands"] = accumulatedJSON(commands)
		projected.Item["action"] = accumulatedJSON(action)
	}
	if progress.shellOutput != nil {
		for index, part := range progress.shellOutput {
			fields := make(map[string]json.RawMessage)
			if part.doneOutput != "" {
				fields["output"] = json.RawMessage(part.doneOutput)
			} else if part.stdoutPresent {
				fields["stdout"] = accumulatedJSONString(part.stdout.String())
			}
			if part.stderrPresent {
				fields["stderr"] = accumulatedJSONString(part.stderr.String())
			}
			projected.Content[index] = fields
		}
	}
}
