package openai_test

import (
	"encoding/json"
	"reflect"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/realtime"
	"github.com/openai/openai-go/v3/responses"
)

type unionListItem interface {
	Overrides() (any, bool)
}

func assertUnknownUnionListRoot[T ~[]E, E unionListItem](t *testing.T) {
	t.Helper()

	body := []byte(`[{"type":"extension","metadata":{"source":"fixture"}}]`)
	var got T
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("decoded items = %d, want 1", len(got))
	}
	override, ok := got[0].Overrides()
	if !ok {
		t.Fatal("unknown list item was not preserved as raw JSON")
	}
	raw, ok := override.(json.RawMessage)
	if !ok {
		t.Fatalf("preserved item type = %T, want json.RawMessage", override)
	}
	var gotUnknown any
	if err := json.Unmarshal(raw, &gotUnknown); err != nil {
		t.Fatalf("decode preserved item: %v", err)
	}
	var wantUnknown any
	if err := json.Unmarshal(body[1:len(body)-1], &wantUnknown); err != nil {
		t.Fatalf("decode expected item: %v", err)
	}
	if !reflect.DeepEqual(gotUnknown, wantUnknown) {
		t.Fatalf("preserved item = %#v, want %#v", gotUnknown, wantUnknown)
	}

	encoded, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var roundTrip any
	if err := json.Unmarshal(encoded, &roundTrip); err != nil {
		t.Fatalf("decode marshaled list: %v", err)
	}
	var wantRoundTrip any
	if err := json.Unmarshal(body, &wantRoundTrip); err != nil {
		t.Fatalf("decode expected list: %v", err)
	}
	if !reflect.DeepEqual(roundTrip, wantRoundTrip) {
		t.Fatalf("round-tripped list = %#v, want %#v", roundTrip, wantRoundTrip)
	}
}

func TestParamUnionListRootsPreserveUnknownElements(t *testing.T) {
	t.Run("beta computer actions", func(t *testing.T) {
		assertUnknownUnionListRoot[openai.BetaComputerActionListParam, openai.BetaComputerActionUnionParam](t)
	})
	t.Run("beta function call output items", func(t *testing.T) {
		assertUnknownUnionListRoot[openai.BetaResponseFunctionCallOutputItemListParam, openai.BetaResponseFunctionCallOutputItemUnionParam](t)
	})
	t.Run("beta response input", func(t *testing.T) {
		assertUnknownUnionListRoot[openai.BetaResponseInputParam, openai.BetaResponseInputItemUnionParam](t)
	})
	t.Run("beta response input message content", func(t *testing.T) {
		assertUnknownUnionListRoot[openai.BetaResponseInputMessageContentListParam, openai.BetaResponseInputContentUnionParam](t)
	})
	t.Run("computer actions", func(t *testing.T) {
		assertUnknownUnionListRoot[responses.ComputerActionListParam, responses.ComputerActionUnionParam](t)
	})
	t.Run("function call output items", func(t *testing.T) {
		assertUnknownUnionListRoot[responses.ResponseFunctionCallOutputItemListParam, responses.ResponseFunctionCallOutputItemUnionParam](t)
	})
	t.Run("response input", func(t *testing.T) {
		assertUnknownUnionListRoot[responses.ResponseInputParam, responses.ResponseInputItemUnionParam](t)
	})
	t.Run("response input message content", func(t *testing.T) {
		assertUnknownUnionListRoot[responses.ResponseInputMessageContentListParam, responses.ResponseInputContentUnionParam](t)
	})
	t.Run("realtime tools", func(t *testing.T) {
		assertUnknownUnionListRoot[realtime.RealtimeToolsConfigParam, realtime.RealtimeToolsConfigUnionParam](t)
	})
}

func TestParamUnionListRootPreservesNullSliceSemantics(t *testing.T) {
	got := responses.ResponseInputParam{{OfMessage: &responses.EasyInputMessageParam{}}}
	if err := json.Unmarshal([]byte(`null`), &got); err != nil {
		t.Fatalf("unmarshal null: %v", err)
	}
	if got != nil {
		t.Fatalf("decoded null list = %#v, want nil", got)
	}
}

func TestStandaloneParamUnionStillRejectsUnknownVariant(t *testing.T) {
	var got responses.ResponseInputItemUnionParam
	if err := json.Unmarshal([]byte(`{"type":"extension"}`), &got); err == nil {
		t.Fatal("unknown standalone union unexpectedly decoded without an error")
	}
}

func TestParamUnionListRootStillRejectsMalformedElement(t *testing.T) {
	var got responses.ResponseInputParam
	if err := json.Unmarshal([]byte(`[{"type":42}]`), &got); err == nil {
		t.Fatal("wrong-typed list discriminator unexpectedly decoded without an error")
	}
	if err := json.Unmarshal([]byte(`[{"metadata":{"source":"fixture"}}]`), &got); err == nil {
		t.Fatal("missing list discriminator unexpectedly decoded without an error")
	}
}
