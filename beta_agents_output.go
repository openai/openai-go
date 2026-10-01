package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/openai/openai-go/v3/packages/ssestream"
)

// BetaAgentOutput binds an Agents output schema to the application's typed parser.
// This beta helper is experimental. It does not change existing sessions' schemas.
type BetaAgentOutput[T any] struct {
	schema []byte
	parse  func([]byte) (T, error)
}

// NewBetaAgentOutput normalizes a JSON schema for an object-root Agents output.
// schema may be a schema map or a JSON-marshalable schema from a reflection library.
// Objects become closed with all properties required. Defaults, composition other
// than nested anyOf, remote references, and unsupported keywords are rejected.
// parse must decode and validate the output according to the application's type.
func NewBetaAgentOutput[T any](schema any, parse func([]byte) (T, error)) (*BetaAgentOutput[T], error) {
	if parse == nil {
		return nil, errors.New("beta agent output requires a parser")
	}
	raw, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("beta agent output schema: %w", err)
	}
	var root map[string]any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err = decoder.Decode(&root); err != nil {
		return nil, fmt.Errorf("beta agent output schema: %w", err)
	}
	if root["type"] != "object" {
		return nil, errors.New("beta agent output schema must have root type object")
	}
	for _, key := range []string{"anyOf", "$ref", "enum", "const"} {
		if _, ok := root[key]; ok {
			return nil, fmt.Errorf("beta agent output schema does not support root %s", key)
		}
	}
	if err = betaAgentNormalizeSchema(root, root, "schema"); err != nil {
		return nil, err
	}
	raw, err = json.Marshal(root)
	if err != nil {
		return nil, err
	}
	return &BetaAgentOutput[T]{schema: raw, parse: parse}, nil
}

// Format returns an independent Agents text.format value for creation params.
// Reuse this adapter's FinalResult or ParseResult to parse the completed answer.
func (o *BetaAgentOutput[T]) Format() TextFormatParamUnion {
	var schema map[string]any
	decoder := json.NewDecoder(bytes.NewReader(o.schema))
	decoder.UseNumber()
	_ = decoder.Decode(&schema) // The constructor has already validated this JSON.
	return TextFormatParamOfParamJSONSchema(schema)
}

// BetaAgentParsedTurnResult retains the raw result alongside the parsed answer.
// This beta helper is experimental.
type BetaAgentParsedTurnResult[T any] struct {
	*BetaAgentTurnResult
	OutputParsed T
}

// BetaAgentOutputParseError is a local parsing failure, not a hosted turn failure.
// Result retains the completed raw answer. This beta helper is experimental.
type BetaAgentOutputParseError struct {
	Result *BetaAgentTurnResult
	Cause  error
}

func (e *BetaAgentOutputParseError) Error() string { return "cannot parse completed beta agent output" }
func (e *BetaAgentOutputParseError) Unwrap() error { return e.Cause }

// ParseResult parses a completed result. It makes no requests and does not change
// the session's output format. Use the same schema when configuring that session.
func (o *BetaAgentOutput[T]) ParseResult(result *BetaAgentTurnResult) (*BetaAgentParsedTurnResult[T], error) {
	if result == nil {
		return nil, &BetaAgentOutputParseError{Cause: errors.New("missing completed result")}
	}
	value, err := o.parse([]byte(result.OutputText()))
	if err != nil {
		return nil, &BetaAgentOutputParseError{Result: result, Cause: err}
	}
	return &BetaAgentParsedTurnResult[T]{BetaAgentTurnResult: result, OutputParsed: value}, nil
}

// FinalResult consumes a creation stream or AgentSessionStream using its existing
// final-result helper, then parses the completed answer. Progress iteration still
// requires the stream's explicit result-collection opt-in. For follow-up streams
// this only parses output; it never changes the existing session's schema.
func (o *BetaAgentOutput[T]) FinalResult(stream any) (*BetaAgentParsedTurnResult[T], error) {
	var result *BetaAgentTurnResult
	var err error
	switch stream := stream.(type) {
	case *ssestream.Stream[AgentSessionEventUnion]:
		result, err = BetaAgentSessionFinalResult(stream)
	case *AgentSessionStream:
		if stream == nil {
			return nil, &BetaAgentTurnResultError{Reason: "unsupported_stream"}
		}
		result, err = stream.FinalResult()
	default:
		return nil, &BetaAgentTurnResultError{Reason: "unsupported_stream"}
	}
	if err != nil {
		return nil, err
	}
	return o.ParseResult(result)
}

func betaAgentNormalizeSchema(schema, root map[string]any, path string) error {
	fail := func(reason string) error { return fmt.Errorf("beta agent output %s: %s", path, reason) }
	for key := range schema {
		switch key {
		case "type", "properties", "required", "additionalProperties", "items", "enum", "const", "anyOf", "$defs", "$ref", "$schema", "title", "description":
		default:
			return fail("unsupported keyword " + key)
		}
	}
	if definitions, exists := schema["$defs"]; exists {
		defs, ok := definitions.(map[string]any)
		if !ok {
			return fail("$defs must be an object")
		}
		for name, definition := range defs {
			child, ok := definition.(map[string]any)
			if !ok {
				return fail("$defs entries must be schemas")
			}
			if err := betaAgentNormalizeSchema(child, root, path+".$defs."+name); err != nil {
				return err
			}
		}
	}
	if ref, exists := schema["$ref"]; exists {
		name, ok := ref.(string)
		if !ok || !strings.HasPrefix(name, "#/$defs/") {
			return fail("only local $defs references are supported")
		}
		target := strings.TrimPrefix(name, "#/$defs/")
		target = strings.ReplaceAll(strings.ReplaceAll(target, "~1", "/"), "~0", "~")
		defs, _ := root["$defs"].(map[string]any)
		if _, ok := defs[target].(map[string]any); !ok {
			return fail("unresolved $ref")
		}
		for key := range schema {
			if key != "$ref" && key != "description" && key != "title" {
				return fail("$ref siblings other than title/description are unsupported")
			}
		}
		return nil
	}
	if variants, exists := schema["anyOf"]; exists {
		list, ok := variants.([]any)
		if !ok || len(list) == 0 {
			return fail("anyOf must contain schemas")
		}
		for i, variant := range list {
			child, ok := variant.(map[string]any)
			if !ok {
				return fail("anyOf entries must be schemas")
			}
			if err := betaAgentNormalizeSchema(child, root, fmt.Sprintf("%s.anyOf[%d]", path, i)); err != nil {
				return err
			}
		}
		for key := range schema {
			if key != "anyOf" && key != "description" && key != "title" {
				return fail("anyOf siblings other than title/description are unsupported")
			}
		}
		return nil
	}
	kind, ok := schema["type"].(string)
	if !ok {
		kinds, valid := schema["type"].([]any)
		if !valid || len(kinds) != 2 {
			return fail("type must be a supported type or nullable type pair")
		}
		if kinds[0] == "null" {
			kind, ok = kinds[1].(string)
		} else if kinds[1] == "null" {
			kind, ok = kinds[0].(string)
		}
		if !ok || kind == "null" {
			return fail("type union must contain one supported type and null")
		}
	}
	switch kind {
	case "object":
		if additional, exists := schema["additionalProperties"]; exists && additional != false {
			return fail("additionalProperties must be false")
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			return fail("object properties must be an object")
		}
		names := make([]string, 0, len(properties))
		for name, property := range properties {
			child, ok := property.(map[string]any)
			if !ok {
				return fail("property must be a schema")
			}
			if err := betaAgentNormalizeSchema(child, root, path+".properties."+name); err != nil {
				return err
			}
			names = append(names, name)
		}
		slices.Sort(names)
		schema["additionalProperties"] = false
		schema["required"] = names
	case "array":
		items, ok := schema["items"].(map[string]any)
		if !ok {
			return fail("array items must be one schema")
		}
		if err := betaAgentNormalizeSchema(items, root, path+".items"); err != nil {
			return err
		}
	case "string", "number", "integer", "boolean", "null":
	default:
		return fail("unsupported type " + kind)
	}
	for _, key := range []string{"properties", "required", "additionalProperties", "items"} {
		if _, exists := schema[key]; exists && !((kind == "object" && key != "items") || (kind == "array" && key == "items")) {
			return fail(key + " does not apply to " + kind)
		}
	}
	if values, exists := schema["enum"]; exists {
		list, ok := values.([]any)
		if !ok || len(list) == 0 {
			return fail("enum must be a nonempty array")
		}
	}
	return nil
}
