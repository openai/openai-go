package openai_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/respjson"
)

func TestUsageCacheWriteTransport(t *testing.T) {
	for _, endpoint := range []string{"completions", "decisions"} {
		for _, tc := range []struct {
			name  string
			raw   string
			want  int64
			valid bool
		}{
			{"value", "73", 73, true},
			{"zero", "0", 0, true},
			{"null", "null", 0, false},
			{"omitted", "", 0, false},
		} {
			t.Run(endpoint+"/"+tc.name, func(t *testing.T) {
				row := map[string]any{
					"object":                       "organization.usage." + endpoint + ".result",
					"input_tokens":                 200,
					"output_tokens":                20,
					"num_model_requests":           1,
					"input_cache_write_tokens":     17,
					"input_cache_write_12h_tokens": 19,
				}
				if tc.raw != "" {
					row["input_cache_write_30m_tokens"] = json.RawMessage(tc.raw)
				}
				payload := map[string]any{
					"object": "page", "has_more": false, "next_page": nil,
					"data": []any{map[string]any{"object": "bucket", "start_time": 0, "end_time": 60, "results": []any{row}}},
				}
				var calls atomic.Int32
				server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					if r.Method != http.MethodGet || r.URL.Path != "/organization/usage/"+endpoint || r.URL.Query().Get("start_time") != "0" {
						t.Error("unexpected usage request")
					}
					w.Header().Set("Content-Type", "application/json")
					if err := json.NewEncoder(w).Encode(payload); err != nil {
						t.Error(err)
					}
				}))
				defer server.Close()
				client := openai.NewClient(option.WithAPIKey("fake-project-key"), option.WithAdminAPIKey("fake-admin-key"), option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithMaxRetries(0))
				check := func(value int64, field respjson.Field, legacy int64, longRetention int64) {
					t.Helper()
					if value != tc.want {
						t.Fatalf("input_cache_write_30m_tokens = %d, want %d", value, tc.want)
					}
					if field.Valid() != tc.valid || field.Raw() != tc.raw {
						t.Fatalf("unexpected field presence: valid=%t raw=%q", field.Valid(), field.Raw())
					}
					if legacy != 17 || longRetention != 19 {
						t.Fatal("existing cache-write counters changed")
					}
				}
				if endpoint == "completions" {
					response, err := client.Admin.Organization.Usage.Completions(context.Background(), openai.AdminOrganizationUsageCompletionsParams{StartTime: 0})
					if err != nil {
						t.Fatal(err)
					}
					if len(response.Data) != 1 || len(response.Data[0].Results) != 1 {
						t.Fatal("missing usage result")
					}
					result := response.Data[0].Results[0]
					check(result.InputCacheWrite30mTokens, result.JSON.InputCacheWrite30mTokens, result.InputCacheWriteTokens, result.InputCacheWrite12hTokens)
					variant := result.AsOrganizationUsageCompletionsResult()
					check(variant.InputCacheWrite30mTokens, variant.JSON.InputCacheWrite30mTokens, variant.InputCacheWriteTokens, variant.InputCacheWrite12hTokens)
				} else {
					response, err := client.Admin.Organization.Usage.Decisions(context.Background(), openai.AdminOrganizationUsageDecisionsParams{StartTime: 0})
					if err != nil {
						t.Fatal(err)
					}
					if len(response.Data) != 1 || len(response.Data[0].Results) != 1 {
						t.Fatal("missing usage result")
					}
					result := response.Data[0].Results[0]
					check(result.InputCacheWrite30mTokens, result.JSON.InputCacheWrite30mTokens, result.InputCacheWriteTokens, result.InputCacheWrite12hTokens)
				}
				if calls.Load() != 1 {
					t.Fatalf("request count = %d, want 1", calls.Load())
				}
			})
		}
	}
}
