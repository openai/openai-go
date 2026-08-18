package apijson

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3/packages/param"
)

type arrayDiscriminatedUnion struct {
	OfText    *arrayDiscriminatedText    `json:",omitzero,inline"`
	OfRefusal *arrayDiscriminatedRefusal `json:",omitzero,inline"`
	param.APIUnion
}

type arrayDiscriminatedText struct {
	Type        string                          `json:"type"`
	Text        string                          `json:"text"`
	Nested      arrayNestedDiscriminatedUnion   `json:"nested"`
	NestedParts []arrayNestedDiscriminatedUnion `json:"nested_parts"`
}

type arrayDiscriminatedRefusal struct {
	Type    string `json:"type"`
	Refusal string `json:"refusal"`
}

type arrayDiscriminatedHolder struct {
	Parts []arrayDiscriminatedUnion `json:"parts"`
}

type arrayNestedDiscriminatedUnion struct {
	OfKnown *arrayNestedDiscriminatedKnown `json:",omitzero,inline"`
	param.APIUnion
}

type arrayNestedDiscriminatedKnown struct {
	Type string `json:"type"`
}

func init() {
	RegisterUnion[arrayDiscriminatedUnion](
		"type",
		Discriminator[arrayDiscriminatedText]("text"),
		Discriminator[arrayDiscriminatedRefusal]("refusal"),
	)
	RegisterUnion[arrayNestedDiscriminatedUnion](
		"type",
		Discriminator[arrayNestedDiscriminatedKnown]("known"),
	)
}

func TestArrayDecoderPreservesUnknownDiscriminatedUnionVariant(t *testing.T) {
	var got arrayDiscriminatedHolder
	err := Unmarshal([]byte(`{"parts":[{"type":"extension","text":"unknown"},{"type":"text","text":"kept"}]}`), &got)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if len(got.Parts) != 2 {
		t.Fatalf("decoded parts = %d, want 2", len(got.Parts))
	}
	unknown, ok := got.Parts[0].Overrides()
	if !ok || !reflect.DeepEqual(unknown, json.RawMessage(`{"type":"extension","text":"unknown"}`)) {
		t.Fatalf("unknown union variant was not preserved: %#v", unknown)
	}
	if got.Parts[1].OfText == nil || got.Parts[1].OfText.Text != "kept" {
		t.Fatalf("supported union variant was not preserved: %#v", got.Parts[1])
	}
}

func TestArrayDecoderPropagatesKnownUnionElementExactness(t *testing.T) {
	cases := map[string]exactness{
		`{"parts":[{"type":"text","text":"kept","extra":true}]}`: extras,
		`{"parts":[{"type":"text","text":42}]}`:                  loose,
	}

	for raw, expected := range cases {
		t.Run(raw, func(t *testing.T) {
			var got arrayDiscriminatedHolder
			decoder := decoderBuilder{root: true}
			exactness, err := decoder.unmarshalWithExactness([]byte(raw), &got)
			if err != nil {
				t.Fatalf("Unmarshal returned error: %v", err)
			}
			if exactness != expected {
				t.Fatalf("exactness = %d, want %d", exactness, expected)
			}
		})
	}
}

func TestArrayDecoderDoesNotPreserveMalformedUnionElement(t *testing.T) {
	var got arrayDiscriminatedHolder
	err := Unmarshal([]byte(`{"parts":[42,{"type":"text","text":"kept"}]}`), &got)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if len(got.Parts) != 0 {
		t.Fatalf("malformed element should keep the array invalid; decoded parts = %d", len(got.Parts))
	}
}

func TestUnknownDiscriminatedUnionStillErrorsOutsideArray(t *testing.T) {
	var got arrayDiscriminatedUnion
	if err := Unmarshal([]byte(`{"type":"extension","text":"unknown"}`), &got); err == nil {
		t.Fatal("unknown discriminated union unexpectedly decoded without an error")
	}
}

func TestArrayDecoderDoesNotPreserveNestedUnknownDiscriminatedUnion(t *testing.T) {
	var got arrayDiscriminatedHolder
	err := Unmarshal([]byte(`{"parts":[{"type":"text","text":"kept","nested":{"type":"extension"}}]}`), &got)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if len(got.Parts) != 1 || got.Parts[0].OfText == nil {
		t.Fatalf("known array element was not decoded: %#v", got.Parts)
	}
	if _, ok := got.Parts[0].OfText.Nested.Overrides(); ok {
		t.Fatal("nested unknown union was unexpectedly preserved")
	}
}

func TestArrayDecoderPreservesUnknownNestedArrayElement(t *testing.T) {
	var got arrayDiscriminatedHolder
	err := Unmarshal([]byte(`{"parts":[{"type":"text","text":"kept","nested_parts":[{"type":"extension"}]}]}`), &got)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if len(got.Parts) != 1 || got.Parts[0].OfText == nil || len(got.Parts[0].OfText.NestedParts) != 1 {
		t.Fatalf("nested array element was not decoded: %#v", got.Parts)
	}
	if _, ok := got.Parts[0].OfText.NestedParts[0].Overrides(); !ok {
		t.Fatal("unknown nested array element was not preserved")
	}
}

func TestMissingDiscriminatorStillErrors(t *testing.T) {
	var got arrayDiscriminatedUnion
	if err := Unmarshal([]byte(`{"text":"missing type"}`), &got); err == nil {
		t.Fatal("missing discriminator unexpectedly decoded without an error")
	}
}

func TestArrayDecoderDoesNotPreserveMissingDiscriminator(t *testing.T) {
	var got arrayDiscriminatedHolder
	if err := Unmarshal([]byte(`{"parts":[{"text":"missing type"},{"type":"text","text":"kept"}]}`), &got); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if len(got.Parts) != 0 {
		t.Fatalf("missing discriminator should keep the array invalid; decoded parts = %d", len(got.Parts))
	}
}

func TestArrayDecoderDoesNotPreserveWrongDiscriminatorType(t *testing.T) {
	var got arrayDiscriminatedHolder
	if err := Unmarshal([]byte(`{"parts":[{"type":42,"text":"wrong type"},{"type":"text","text":"kept"}]}`), &got); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if len(got.Parts) != 0 {
		t.Fatalf("wrong discriminator type should keep the array invalid; decoded parts = %d", len(got.Parts))
	}
}
