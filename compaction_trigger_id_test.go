package openai_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/responses"
)

type compactionIDHTTPDoer func(*http.Request) (*http.Response, error)

func (f compactionIDHTTPDoer) Do(r *http.Request) (*http.Response, error) { return f(r) }

func TestCompactionTriggerIDTransport(t *testing.T) {
	for _, beta := range []bool{false, true} {
		name := "stable"
		if beta {
			name = "beta"
		}
		for _, tc := range []struct {
			name, field, value string
			valid              bool
		}{
			{"populated", `,"id":"msg_transport_compat"`, "msg_transport_compat", true},
			{"null", `,"id":null`, "", false},
			{"omitted", "", "", false},
			{"empty", `,"id":""`, "", true},
		} {
			t.Run(name+"/"+tc.name, func(t *testing.T) {
				calls := 0
				client := openai.NewClient(
					option.WithAPIKey("synthetic-key"),
					option.WithBaseURL("https://example.invalid/v1"),
					option.WithMaxRetries(0),
					option.WithHTTPClient(compactionIDHTTPDoer(func(r *http.Request) (*http.Response, error) {
						calls++
						if r.Method != http.MethodGet || r.URL.Path != "/v1/responses/resp_compat" {
							t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
						}
						body := `{"id":"resp_compat","object":"response","instructions":[{"type":"compaction_trigger"` + tc.field + `}],"output":[]}`
						return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
					})),
				)
				var id string
				var metadata respjson.Field
				if beta {
					response, err := client.Beta.Responses.Get(context.Background(), "resp_compat", openai.BetaResponseGetParams{})
					if err != nil {
						t.Fatal(err)
					}
					items := response.Instructions.AsInputItemList()
					if len(items) != 1 {
						t.Fatalf("got %d instructions; want one", len(items))
					}
					trigger := items[0].AsCompactionTrigger()
					id, metadata = trigger.ID, trigger.JSON.ID
				} else {
					response, err := client.Responses.Get(context.Background(), "resp_compat", responses.ResponseGetParams{})
					if err != nil {
						t.Fatal(err)
					}
					items := response.Instructions.AsInputItemList()
					if len(items) != 1 {
						t.Fatalf("got %d instructions; want one", len(items))
					}
					trigger := items[0].AsCompactionTrigger()
					id, metadata = trigger.ID, trigger.JSON.ID
				}
				if calls != 1 {
					t.Fatalf("transport called %d times; want one", calls)
				}
				if id != tc.value {
					t.Fatalf("decoded compaction trigger ID = %q; want %q", id, tc.value)
				}
				if metadata.Valid() != tc.valid || (tc.name == "null" && metadata.Raw() != "null") || (tc.name == "omitted" && metadata.Raw() != "") {
					t.Fatalf("unexpected ID metadata: valid=%v raw=%q", metadata.Valid(), metadata.Raw())
				}
			})
		}
	}
}

func TestCompactionTriggerIDCompatibility(t *testing.T) {
	cases := []struct {
		name     string
		id       param.Opt[string]
		expected string
	}{
		{"omitted", param.Opt[string]{}, `{"type":"compaction_trigger"}`},
		{"null", param.Null[string](), `{"type":"compaction_trigger","id":null}`},
		{"populated", openai.String("msg_123"), `{"type":"compaction_trigger","id":"msg_123"}`},
		{"empty", openai.String(""), `{"type":"compaction_trigger","id":""}`},
	}
	assertJSON := func(t *testing.T, value any, expected string) {
		t.Helper()
		wire, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		var got, want map[string]any
		if err = json.Unmarshal(wire, &got); err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal([]byte(expected), &want); err != nil {
			t.Fatal(err)
		}
		if len(got) != len(want) {
			t.Fatalf("got %s; want %s", wire, expected)
		}
		for key, value := range want {
			if actual, exists := got[key]; !exists || actual != value {
				t.Fatalf("got %s; want %s", wire, expected)
			}
		}
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// This constructor existed before ID; keep source compatibility.
			trigger := responses.NewResponseInputItemCompactionTriggerParam()
			trigger.ID = tt.id
			stableUnion := responses.ResponseInputItemUnionParam{OfCompactionTrigger: &trigger}
			assertJSON(t, stableUnion, tt.expected)
			beta := openai.BetaResponseInputItemCompactionTriggerParam{ID: tt.id}
			betaUnion := openai.BetaResponseInputItemUnionParam{OfCompactionTrigger: &beta}
			assertJSON(t, betaUnion, tt.expected)
			for _, id := range []*string{stableUnion.GetID(), betaUnion.GetID()} {
				if tt.id.Valid() {
					if id == nil || *id != tt.id.Value {
						t.Fatalf("unexpected union ID: %v", id)
					}
				} else if id != nil {
					t.Fatalf("omitted/null ID should not have a value: %v", *id)
				}
			}
			var stableResponse responses.ResponseInputItemCompactionTrigger
			var betaResponse openai.BetaResponseInputItemCompactionTrigger
			if err := json.Unmarshal([]byte(tt.expected), &stableResponse); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(tt.expected), &betaResponse); err != nil {
				t.Fatal(err)
			}
			var raw map[string]json.RawMessage
			if err := json.Unmarshal([]byte(tt.expected), &raw); err != nil {
				t.Fatal(err)
			}
			if stableResponse.ID != tt.id.Value || betaResponse.ID != tt.id.Value ||
				stableResponse.JSON.ID.Valid() != tt.id.Valid() || betaResponse.JSON.ID.Valid() != tt.id.Valid() ||
				stableResponse.JSON.ID.Raw() != string(raw["id"]) || betaResponse.JSON.ID.Raw() != string(raw["id"]) {
				t.Fatalf("response ID value/presence did not preserve %s", tt.expected)
			}
		})
	}
	t.Run("existing_extra_field_escape_hatch", func(t *testing.T) {
		trigger := responses.NewResponseInputItemCompactionTriggerParam()
		trigger.SetExtraFields(map[string]any{"id": "msg_existing", "future": "preserved"})
		assertJSON(t, trigger, `{"type":"compaction_trigger","id":"msg_existing","future":"preserved"}`)
	})
	t.Run("beta_agent_metadata", func(t *testing.T) {
		trigger := openai.BetaResponseInputItemCompactionTriggerParam{
			ID:    openai.String("msg_123"),
			Agent: openai.BetaResponseInputItemCompactionTriggerAgentParam{AgentName: "assistant"},
		}
		wire, err := json.Marshal(trigger)
		if err != nil {
			t.Fatal(err)
		}
		var got struct {
			ID    string `json:"id"`
			Agent struct {
				AgentName string `json:"agent_name"`
			} `json:"agent"`
		}
		if err = json.Unmarshal(wire, &got); err != nil {
			t.Fatal(err)
		}
		if got.ID != "msg_123" || got.Agent.AgentName != "assistant" {
			t.Fatalf("unexpected wire body: %s", wire)
		}
	})
}
