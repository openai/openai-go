package openai_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// These handwritten regressions exercise the public HTTP entrypoints with
// synthetic credentials only; they do not depend on the API mock server.
func TestBetaAgentVaultCredentialEnvironmentVariable(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "synthetic-env-api-key")
	t.Setenv("OPENAI_ADMIN_KEY", "synthetic-env-admin-key")
	t.Setenv("OPENAI_WEBHOOK_SECRET", "synthetic-env-webhook-secret")
	t.Setenv("TEST_CREDENTIAL", "synthetic-env-value-must-not-be-read")

	const limitedResponse = `{"id":"cred_test","name":"test","created_at":1,"updated_at":2,"auth":{"type":"environment_variable","secret_name":"TEST_CREDENTIAL","networking":{"type":"limited","allowed_hosts":["api.example.test","192.0.2.1"]}}}`
	const unrestrictedResponse = `{"id":"cred_test","name":"test","created_at":1,"updated_at":2,"auth":{"type":"environment_variable","secret_name":"TEST_CREDENTIAL","networking":{"type":"unrestricted"}}}`

	tests := []struct {
		name              string
		invoke            func(context.Context, openai.Client) (*openai.Credential, error)
		path              string
		request           string
		response          string
		secret            string
		networkingType    string
		networkingVariant any
		allowedHosts      []string
	}{
		{
			name: "create limited",
			invoke: func(ctx context.Context, client openai.Client) (*openai.Credential, error) {
				return client.Beta.Agents.Vaults.Credentials.New(ctx, "vault_test", openai.BetaAgentVaultCredentialNewParams{
					Name: "test",
					Auth: openai.CredentialAuthCreateParamOfParamEnvironmentVariable(
						openai.CredentialNetworkingParamLimited{AllowedHosts: []string{"api.example.test", "192.0.2.1"}},
						"TEST_CREDENTIAL", "synthetic-create-limited",
					),
				})
			},
			path:              "/vaults/vault_test/credentials",
			request:           `{"name":"test","auth":{"type":"environment_variable","secret_name":"TEST_CREDENTIAL","secret_value":"synthetic-create-limited","networking":{"type":"limited","allowed_hosts":["api.example.test","192.0.2.1"]}}}`,
			response:          limitedResponse,
			secret:            "synthetic-create-limited",
			networkingType:    "limited",
			networkingVariant: openai.CredentialNetworkingLimited{},
			allowedHosts:      []string{"api.example.test", "192.0.2.1"},
		},
		{
			name: "create unrestricted",
			invoke: func(ctx context.Context, client openai.Client) (*openai.Credential, error) {
				return client.Beta.Agents.Vaults.Credentials.New(ctx, "vault_test", openai.BetaAgentVaultCredentialNewParams{
					Name: "test",
					Auth: openai.CredentialAuthCreateParamOfParamEnvironmentVariable(
						openai.NewCredentialNetworkingParamUnrestricted(),
						"TEST_CREDENTIAL", "synthetic-create-unrestricted",
					),
				})
			},
			path:              "/vaults/vault_test/credentials",
			request:           `{"name":"test","auth":{"type":"environment_variable","secret_name":"TEST_CREDENTIAL","secret_value":"synthetic-create-unrestricted","networking":{"type":"unrestricted"}}}`,
			response:          unrestrictedResponse,
			secret:            "synthetic-create-unrestricted",
			networkingType:    "unrestricted",
			networkingVariant: openai.CredentialNetworkingUnrestricted{},
		},
		{
			name: "rotate secret without changing configuration",
			invoke: func(ctx context.Context, client openai.Client) (*openai.Credential, error) {
				return client.Beta.Agents.Vaults.Credentials.Update(ctx, "vault_test", "cred_test", openai.BetaAgentVaultCredentialUpdateParams{
					Auth: openai.CredentialAuthRotateParamOfParamEnvironmentVariable("synthetic-rotated-secret"),
				})
			},
			path:              "/vaults/vault_test/credentials/cred_test",
			request:           `{"auth":{"type":"environment_variable","secret_value":"synthetic-rotated-secret"}}`,
			response:          limitedResponse,
			secret:            "synthetic-rotated-secret",
			networkingType:    "limited",
			networkingVariant: openai.CredentialNetworkingLimited{},
			allowedHosts:      []string{"api.example.test", "192.0.2.1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wantRequest map[string]any
			if err := json.Unmarshal([]byte(tt.request), &wantRequest); err != nil {
				t.Fatalf("json.Unmarshal(%q request fixture) error = %v, want nil", tt.name, err)
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != tt.path {
					t.Errorf("Credentials(%q) request = %s %q, want POST %q", tt.name, r.Method, r.URL.Path, tt.path)
				}
				if got, want := r.Header.Get("Authorization"), "Bearer synthetic-api-key"; got != want {
					t.Errorf("Credentials(%q) Authorization matches explicit synthetic key = %t, want true", tt.name, got == want)
				}
				var gotRequest map[string]any
				if err := json.NewDecoder(r.Body).Decode(&gotRequest); err != nil {
					t.Errorf("Credentials(%q) request JSON decode error = %v, want nil", tt.name, err)
					http.Error(w, "invalid synthetic request", http.StatusBadRequest)
					return
				}
				if !reflect.DeepEqual(gotRequest, wantRequest) {
					t.Errorf("Credentials(%q) request JSON = %#v, want %#v", tt.name, gotRequest, wantRequest)
				}
				w.Header().Set("Content-Type", "application/json")
				if _, err := io.WriteString(w, tt.response); err != nil {
					t.Errorf("Credentials(%q) synthetic response write error = %v, want nil", tt.name, err)
				}
			}))
			t.Cleanup(server.Close)

			client := openai.NewClient(
				option.WithBaseURL(server.URL),
				option.WithAPIKey("synthetic-api-key"),
				option.WithAdminAPIKey("synthetic-admin-key"),
				option.WithMaxRetries(0),
			)
			credential, err := tt.invoke(t.Context(), client)
			if err != nil {
				t.Fatalf("Credentials(%q) error = %v, want nil", tt.name, err)
			}
			if got, want := credential.ID, "cred_test"; got != want {
				t.Errorf("Credentials(%q).ID = %q, want %q", tt.name, got, want)
			}
			if got := credential.JSON.Auth.Valid(); !got {
				t.Errorf("Credentials(%q).JSON.Auth.Valid() = %t, want true", tt.name, got)
			}
			auth, ok := credential.Auth.AsAny().(openai.CredentialAuthEnvironmentVariable)
			if !ok {
				t.Fatalf("Credentials(%q).Auth.AsAny() = %T, want CredentialAuthEnvironmentVariable", tt.name, credential.Auth.AsAny())
			}
			if got, want := auth.SecretName, "TEST_CREDENTIAL"; got != want {
				t.Errorf("Credentials(%q).Auth.SecretName = %q, want %q", tt.name, got, want)
			}
			if !auth.JSON.SecretName.Valid() || !auth.JSON.Networking.Valid() {
				t.Errorf("Credentials(%q) secret-name/networking presence = %t/%t, want true/true", tt.name, auth.JSON.SecretName.Valid(), auth.JSON.Networking.Valid())
			}
			if got := auth.Networking.Type; got != tt.networkingType {
				t.Errorf("Credentials(%q).Auth.Networking.Type = %q, want %q", tt.name, got, tt.networkingType)
			}
			if got := auth.Networking.AsAny(); reflect.TypeOf(got) != reflect.TypeOf(tt.networkingVariant) {
				t.Errorf("Credentials(%q).Auth.Networking.AsAny() = %T, want %T", tt.name, got, tt.networkingVariant)
			}
			if got := auth.Networking.AllowedHosts; !slices.Equal(got, tt.allowedHosts) {
				t.Errorf("Credentials(%q).Auth.Networking.AllowedHosts = %v, want %v", tt.name, got, tt.allowedHosts)
			}
			if got, want := auth.Networking.JSON.AllowedHosts.Valid(), tt.networkingType == "limited"; got != want {
				t.Errorf("Credentials(%q).Auth.Networking.JSON.AllowedHosts.Valid() = %t, want %t", tt.name, got, want)
			}
			encoded, err := json.Marshal(credential)
			if err != nil {
				t.Fatalf("json.Marshal(Credentials(%q) metadata) error = %v, want nil", tt.name, err)
			}
			for _, metadata := range []string{string(encoded), credential.RawJSON(), auth.RawJSON()} {
				if strings.Contains(metadata, tt.secret) || strings.Contains(metadata, `"secret_value"`) {
					t.Errorf("Credentials(%q) response metadata contains secret material, want metadata only", tt.name)
				}
			}
		})
	}
}
