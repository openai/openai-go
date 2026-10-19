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

func TestVaultUpdateClearingBody(t *testing.T) {
	for _, tt := range []struct {
		name   string
		params openai.BetaAgentVaultUpdateParams
		want   map[string]any
	}{
		{"omitted", openai.BetaAgentVaultUpdateParams{}, map[string]any{}},
		{"clear name", openai.BetaAgentVaultUpdateParams{Name: param.Null[string]()}, map[string]any{"name": nil}},
		{"clear metadata", openai.BetaAgentVaultUpdateParams{Metadata: map[string]string{}}, map[string]any{"metadata": map[string]any{}}},
		{"replace values", openai.BetaAgentVaultUpdateParams{Name: openai.String("Research"), Metadata: map[string]string{"team": "research"}}, map[string]any{"name": "Research", "metadata": map[string]any{"team": "research"}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/vaults/vault_test" {
					t.Errorf("request = %s %s, want POST /vaults/vault_test", r.Method, r.URL.Path)
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("decode request: %v", err)
					http.Error(w, "invalid JSON", http.StatusBadRequest)
					return
				}
				if !reflect.DeepEqual(body, tt.want) {
					t.Errorf("body = %#v, want %#v", body, tt.want)
				}
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{"id":"vault_test","object":"vault"}`)); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			t.Cleanup(server.Close)
			client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL(server.URL), option.WithUnsafeAllowHTTP(), option.WithHTTPClient(server.Client()), option.WithMaxRetries(0))
			vault, err := client.Beta.Agents.Vaults.Update(t.Context(), "vault_test", tt.params)
			if err != nil {
				t.Fatal(err)
			}
			if vault.ID != "vault_test" {
				t.Errorf("vault ID = %q, want vault_test", vault.ID)
			}
		})
	}
}
