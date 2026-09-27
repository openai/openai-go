package openai_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// Exercise updates which must not rotate or clear the credential's auth.
func TestBetaAgentVaultCredentialMetadataOnlyUpdate(t *testing.T) {
	for _, tt := range []struct {
		name     string
		metadata map[string]string
	}{
		{"set metadata", map[string]string{"purpose": "synthetic-test"}},
		{"clear metadata", map[string]string{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				const path = "/vaults/vault_synthetic/credentials/cred_synthetic"
				if r.Method != http.MethodPost || r.URL.Path != path {
					t.Errorf("Credentials.Update request = %s %q, want POST %q", r.Method, r.URL.Path, path)
				}
				var got map[string]json.RawMessage
				if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
					t.Errorf("Credentials.Update JSON decode error = %v", err)
					http.Error(w, "invalid synthetic JSON", http.StatusBadRequest)
					return
				}
				if _, exists := got["auth"]; exists {
					t.Error("Credentials.Update sent auth during metadata-only update")
				}
				var metadata map[string]string
				if err := json.Unmarshal(got["metadata"], &metadata); err != nil {
					t.Errorf("Credentials.Update metadata decode error = %v", err)
				}
				if !reflect.DeepEqual(metadata, tt.metadata) {
					t.Errorf("Credentials.Update metadata = %#v, want %#v", metadata, tt.metadata)
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]any{
					"id": "cred_synthetic", "object": "vault.credential", "vault_id": "vault_synthetic",
					"name": "synthetic", "created_at": 1, "updated_at": 2, "metadata": tt.metadata,
					"auth": map[string]any{"type": "static_bearer", "mcp_server_url": "https://mcp.example.test"},
				}); err != nil {
					t.Errorf("synthetic credential response error = %v", err)
				}
			}))
			t.Cleanup(server.Close)
			client := openai.NewClient(
				option.WithBaseURL(server.URL),
				option.WithAPIKey("synthetic-api-key"),
				option.WithAdminAPIKey("synthetic-admin-key"),
				option.WithMaxRetries(0),
			)
			result, err := client.Beta.Agents.Vaults.Credentials.Update(t.Context(),
				"vault_synthetic", "cred_synthetic",
				openai.BetaAgentVaultCredentialUpdateParams{Metadata: tt.metadata})
			if err != nil {
				t.Fatalf("Credentials.Update metadata-only error = %v", err)
			}
			if !result.JSON.Metadata.Valid() || !reflect.DeepEqual(result.Metadata, tt.metadata) {
				t.Errorf("credential response metadata = %#v, presence = %t; want %#v with field present",
					result.Metadata, result.JSON.Metadata.Valid(), tt.metadata)
			}
		})
	}
}
