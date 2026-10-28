package openai_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestBetaAgentSessionListMetadataQuery(t *testing.T) {
	var nilMap map[string]string
	for _, tt := range []struct {
		name     string
		metadata *map[string]string
		query    url.Values
	}{
		{name: "omitted"},
		{name: "nil map", metadata: &nilMap},
		{name: "empty map", metadata: &map[string]string{}},
		{
			name:     "populated",
			metadata: &map[string]string{"foo": "string"},
			query:    url.Values{"metadata[foo]": {"string"}},
		},
		{
			name: "escaped and empty values",
			metadata: &map[string]string{
				"team name": "research & infra+ops=1/2?#%",
				"region":    "日本",
				"empty":     "",
			},
			query: url.Values{
				"metadata[team name]": {"research & infra+ops=1/2?#%"},
				"metadata[region]":    {"日本"},
				"metadata[empty]":     {""},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			want := url.Values{"agent_id": {"agent_synthetic"}, "limit": {"2"}}
			for key, values := range tt.query {
				want[key] = values
			}
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				if r.Method != http.MethodGet || r.URL.Path != "/agents/sessions" {
					t.Errorf("Sessions.List request = %s %q, want GET /agents/sessions", r.Method, r.URL.Path)
				}
				got, err := url.ParseQuery(r.URL.RawQuery)
				if err != nil {
					t.Errorf("parse query: %v", err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("Sessions.List query = %#v, want %#v", got, want)
				}
				w.Header().Set("Content-Type", "application/json")
				if _, err := fmt.Fprint(w, `{"object":"list","data":[],"has_more":false,"first_id":null,"last_id":null}`); err != nil {
					t.Errorf("write synthetic response: %v", err)
				}
			}))
			t.Cleanup(server.Close)
			client := openai.NewClient(
				option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL),
				option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic-api-key"),
				option.WithMaxRetries(0),
			)
			_, err := client.Beta.Agents.Sessions.List(t.Context(), openai.BetaAgentSessionListParams{
				AgentID: openai.String("agent_synthetic"), Limit: openai.Int(2), Metadata: tt.metadata,
			})
			if err != nil {
				t.Fatalf("Sessions.List: %v", err)
			}
			if got := requests.Load(); got != 1 {
				t.Errorf("request count = %d, want 1", got)
			}
		})
	}
}

func sessionParamsComparable[T comparable](left, right T) bool { return left == right }

func TestBetaAgentSessionListParamsRemainComparable(t *testing.T) {
	zero := openai.BetaAgentSessionListParams{}
	if !sessionParamsComparable(zero, openai.BetaAgentSessionListParams{}) || zero != (openai.BetaAgentSessionListParams{}) {
		t.Fatal("zero parameters must remain comparable and equal")
	}
	params := openai.BetaAgentSessionListParams{Metadata: &map[string]string{"foo": "string"}}
	copyOfParams := params
	if !sessionParamsComparable(params, copyOfParams) || params != copyOfParams {
		t.Fatal("parameters sharing a metadata pointer must remain equal")
	}
	keys := map[openai.BetaAgentSessionListParams]string{params: "present"}
	if keys[copyOfParams] != "present" {
		t.Fatal("parameters must remain usable as map keys")
	}
}
