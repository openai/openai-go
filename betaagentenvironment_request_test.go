package openai_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
)

func TestBetaAgentEnvironmentNewRequestBody(t *testing.T) {
	for _, tt := range []struct {
		name    string
		desktop openai.BetaAgentEnvironmentNewParamsEnvironmentDesktop
		want    string
	}{
		{
			name: "minimal desktop omitted",
			want: `{"environment":{"type":"openai_hosted"}}`,
		},
		{
			name:    "desktop disabled",
			desktop: openai.BetaAgentEnvironmentNewParamsEnvironmentDesktop{Enabled: openai.Bool(false)},
			want:    `{"environment":{"type":"openai_hosted","desktop":{"enabled":false}}}`,
		},
		{
			name:    "desktop enabled",
			desktop: openai.BetaAgentEnvironmentNewParamsEnvironmentDesktop{Enabled: openai.Bool(true)},
			want:    `{"environment":{"type":"openai_hosted","desktop":{"enabled":true}}}`,
		},
		{
			name:    "desktop null",
			desktop: param.NullStruct[openai.BetaAgentEnvironmentNewParamsEnvironmentDesktop](),
			want:    `{"environment":{"type":"openai_hosted","desktop":null}}`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var want any
			if err := json.Unmarshal([]byte(tt.want), &want); err != nil {
				t.Fatal(err)
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/agents/environments" {
					t.Errorf("request = %s %s, want POST /agents/environments", r.Method, r.URL.Path)
				}
				var got any
				if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
					t.Errorf("decode request body: %v", err)
					http.Error(w, "invalid synthetic JSON", http.StatusBadRequest)
					return
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("request body = %#v, want %#v", got, want)
				}
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{"id":"environment_synthetic","object":"agent.environment","type":"openai_hosted","status":"ready","files":[]}`)); err != nil {
					t.Errorf("write synthetic response: %v", err)
				}
			}))
			t.Cleanup(server.Close)
			client := openai.NewClient(
				option.WithBaseURL(server.URL), option.WithUnsafeAllowHTTP(),
				option.WithAPIKey("synthetic-api-key"), option.WithAdminAPIKey("synthetic-admin-key"),
				option.WithMaxRetries(0),
			)
			params := openai.BetaAgentEnvironmentNewParams{}
			params.Environment.Desktop = tt.desktop
			if _, err := client.Beta.Agents.Environments.New(t.Context(), params); err != nil {
				t.Fatalf("Environments.New: %v", err)
			}
		})
	}
}
