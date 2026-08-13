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
	Type string `json:"type"`
	Text string `json:"text"`
}

type arrayDiscriminatedRefusal struct {
	Type    string `json:"type"`
	Refusal string `json:"refusal"`
}

type arrayDiscriminatedHolder struct {
	Parts []arrayDiscriminatedUnion `json:"parts"`
}

func init() {
	RegisterUnion[arrayDiscriminatedUnion](
		"type",
		Discriminator[arrayDiscriminatedText]("text"),
		Discriminator[arrayDiscriminatedRefusal]("refusal"),
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
