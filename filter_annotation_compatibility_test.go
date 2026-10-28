package openai_test

import (
	"encoding/json"
	"testing"

	openai "github.com/openai/openai-go/v4"
	"github.com/openai/openai-go/v4/conversations"
	"github.com/openai/openai-go/v4/live"
	"github.com/openai/openai-go/v4/packages/param"
	"github.com/openai/openai-go/v4/realtime"
	"github.com/openai/openai-go/v4/responses"
	"github.com/openai/openai-go/v4/shared"
	"github.com/openai/openai-go/v4/webhooks"
)

func TestReleasedCompoundFiltersPreserveNestedValues(t *testing.T) {
	leaf := shared.ComparisonFilterParam{
		Key: "color", Type: shared.ComparisonFilterTypeEq,
		Value: shared.ComparisonFilterValueUnionParam{OfString: param.NewOpt("blue")},
	}
	filter := shared.CompoundFilterParam{
		Type: shared.CompoundFilterTypeAnd,
		Filters: []shared.CompoundFilterFilterUnionParam{{OfFilter: &shared.CompoundFilterParam{
			Type:    shared.CompoundFilterTypeOr,
			Filters: []shared.CompoundFilterFilterUnionParam{{OfComparison: &leaf}},
		}}},
	}
	// Released aliases must accept a shared nested filter in every exported package.
	var a openai.CompoundFilterFilterUnionParam = filter.Filters[0]
	var b conversations.CompoundFilterFilterUnionParam = a
	var c live.CompoundFilterFilterUnionParam = b
	var d realtime.CompoundFilterFilterUnionParam = c
	var e responses.CompoundFilterFilterUnionParam = d
	var f webhooks.CompoundFilterFilterUnionParam = e
	query := openai.VectorStoreSearchParamsFiltersUnion{OfCompoundFilter: &filter}
	tool := responses.FileSearchToolFiltersUnionParam{OfCompoundFilter: &filter}
	if key := f.GetFilters()[0].GetKey(); key == nil || *key != "color" {
		t.Fatalf("nested shared getter: %v", key)
	}
	for _, child := range []shared.CompoundFilterFilterUnionParam{query.GetFilters()[0], tool.GetFilters()[0]} {
		key := child.GetFilters()[0].GetKey()
		if key == nil || *key != "color" {
			t.Fatalf("nested projected getter: %v", key)
		}
	}

	want := `{"type":"and","filters":[{"type":"or","filters":[{"type":"eq","key":"color","value":"blue"}]}]}`
	for _, value := range []any{filter, query, tool} {
		data, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		var actual shared.CompoundFilter
		if err = json.Unmarshal(data, &actual); err != nil {
			t.Fatal(err)
		}
		// Released response aliases must all decode the same nested bytes.
		var ra openai.CompoundFilterFilterUnion = actual.Filters[0]
		var rb conversations.CompoundFilterFilterUnion = ra
		var rc live.CompoundFilterFilterUnion = rb
		var rd realtime.CompoundFilterFilterUnion = rc
		var re responses.CompoundFilterFilterUnion = rd
		var rf webhooks.CompoundFilterFilterUnion = re
		inner := rf.AsCompoundFilter()
		if inner.Type != shared.CompoundFilterTypeOr || inner.Filters[0].AsComparisonFilter().Value.AsString() != "blue" {
			t.Fatalf("%T lost released nested fields: %s", value, data)
		}
		requireJSONEqual(t, data, []byte(want))
	}
}

func TestReleasedBetaCompoundFiltersPreserveNestedValues(t *testing.T) {
	op := "and"
	leaf := openai.BetaFileSearchToolFiltersCompoundFilterFilterParam{
		Key: "color", Type: "eq",
		Value: openai.BetaFileSearchToolFiltersCompoundFilterFilterValueUnionParam{OfString: param.NewOpt("blue")},
	}
	filter := openai.BetaFileSearchToolFiltersCompoundFilterParam{
		Type: op,
		Filters: []openai.BetaFileSearchToolFiltersCompoundFilterFilterUnionParam{{OfFilter: &openai.BetaFileSearchToolFiltersCompoundFilterParam{
			Type:    "or",
			Filters: []openai.BetaFileSearchToolFiltersCompoundFilterFilterUnionParam{{OfFileSearchToolFiltersCompound: &leaf}},
		}}},
	}
	// Released beta operation fields remain plain strings.
	_ = openai.BetaFileSearchToolFiltersCompoundFilter{Type: op}
	_ = openai.BetaFileSearchToolFiltersCompoundFilterParam{Type: op}
	tool := openai.BetaFileSearchToolFiltersUnionParam{OfCompoundFilter: &filter}
	key := tool.GetFilters()[0].GetFilters()[0].GetKey()
	if key == nil || *key != "color" {
		t.Fatalf("beta getter: %v", key)
	}

	want := `{"type":"and","filters":[{"type":"or","filters":[{"type":"eq","key":"color","value":"blue"}]}]}`
	for _, value := range []any{filter, tool} {
		data, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		var actual openai.BetaFileSearchToolFiltersCompoundFilter
		if err = json.Unmarshal(data, &actual); err != nil {
			t.Fatal(err)
		}
		inner := actual.Filters[0].AsBetaFileSearchToolFiltersCompoundFilter()
		if inner.Type != "or" || inner.Filters[0].AsFileSearchToolFiltersCompound().Value.AsString() != "blue" {
			t.Fatalf("%T lost beta nested fields: %s", value, data)
		}
		requireJSONEqual(t, data, []byte(want))
	}
}

func TestReleasedAnnotationAccess(t *testing.T) {
	for _, annotation := range []string{
		`{"type":"file_citation","file_id":"file_1","filename":"guide.md","index":2,"new_metadata":{"v":12}}`,
		`{"type":"new_citation","file_id":"file_1","future":[1,"value",{"deep":true}]}`,
	} {
		wire := []byte(`{"type":"response.output_text.annotation.added","item_id":"item_1","output_index":0,"content_index":0,"annotation_index":0,"sequence_number":1,"annotation":` + annotation + `}`)
		var stable responses.ResponseOutputTextAnnotationAddedEvent
		var stableStream responses.ResponseStreamEventUnion
		var beta openai.BetaResponseOutputTextAnnotationAddedEvent
		var betaStream openai.BetaResponseStreamEventUnion
		for _, target := range []any{&stable, &stableStream, &beta, &betaStream} {
			if err := json.Unmarshal(wire, target); err != nil {
				t.Fatal(err)
			}
		}
		for _, raw := range []string{stable.Annotation.RawJSON(), stableStream.Annotation.RawJSON(), beta.Annotation.RawJSON(), betaStream.Annotation.RawJSON()} {
			requireJSONEqual(t, []byte(raw), []byte(annotation))
		}
		if stable.Annotation.Type == "file_citation" {
			for _, id := range []string{stable.Annotation.AsFileCitation().FileID, stableStream.Annotation.AsFileCitation().FileID, beta.Annotation.AsFileCitation().FileID, betaStream.Annotation.AsFileCitation().FileID} {
				if id != "file_1" {
					t.Fatalf("lost released file citation %q", id)
				}
			}
		}

	}
}

func requireJSONEqual(t *testing.T, actual, expected []byte) {
	t.Helper()
	var a, b any
	if err := json.Unmarshal(actual, &a); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(expected, &b); err != nil {
		t.Fatal(err)
	}
	// Use canonical JSON because field order is not part of the request contract.
	ac, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	bc, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}
	if string(ac) != string(bc) {
		t.Fatalf("got %s, want %s", ac, bc)
	}
}

func TestReleasedAnnotationFromConstructedValue(t *testing.T) {
	stable := responses.ResponseOutputTextAnnotationAddedEvent{
		Annotation: responses.ResponseOutputTextAnnotationAddedEventAnnotationUnion{Type: "file_citation", FileID: "file_manual"},
	}
	stableStream := responses.ResponseStreamEventUnion{Annotation: stable.Annotation}
	beta := openai.BetaResponseOutputTextAnnotationAddedEvent{
		Annotation: openai.BetaResponseOutputTextAnnotationAddedEventAnnotationUnion{Type: "file_citation", FileID: "file_manual"},
	}
	betaStream := openai.BetaResponseStreamEventUnion{Annotation: beta.Annotation}
	for _, id := range []string{stable.Annotation.FileID, stableStream.Annotation.FileID, beta.Annotation.FileID, betaStream.Annotation.FileID} {
		if id != "file_manual" {
			t.Fatalf("constructed annotation file ID = %q", id)
		}
	}
}
