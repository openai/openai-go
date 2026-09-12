// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal/testutil"
	"github.com/openai/openai-go/v3/option"
)

func TestBetaAgentVaultCredentialNewWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Beta.Agents.Vaults.Credentials.New(
		context.TODO(),
		"vault_id",
		openai.BetaAgentVaultCredentialNewParams{
			Auth: openai.CredentialAuthCreateParamUnion{
				OfParamMcpOAuth: &openai.CredentialAuthCreateParamMcpOAuth{
					AccessToken:  "access_token",
					McpServerURL: "mcp_server_url",
					ExpiresAt:    openai.String("expires_at"),
					Refresh: openai.CredentialAuthCreateParamMcpOAuthRefresh{
						ClientID:      "client_id",
						RefreshToken:  "refresh_token",
						TokenEndpoint: "token_endpoint",
						TokenEndpointAuth: openai.McpOAuthTokenEndpointAuthCreateParamUnion{
							OfParamNone: &openai.McpOAuthTokenEndpointAuthCreateParamNone{},
						},
						Resource: openai.String("resource"),
						Scope:    openai.String("scope"),
					},
				},
			},
			Name: "x",
		},
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentVaultCredentialGet(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Beta.Agents.Vaults.Credentials.Get(
		context.TODO(),
		"vault_id",
		"credential_id",
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentVaultCredentialUpdateWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Beta.Agents.Vaults.Credentials.Update(
		context.TODO(),
		"vault_id",
		"credential_id",
		openai.BetaAgentVaultCredentialUpdateParams{
			Auth: openai.CredentialAuthRotateParamUnion{
				OfParamMcpOAuth: &openai.CredentialAuthRotateParamMcpOAuth{
					AccessToken: openai.String("access_token"),
					ExpiresAt:   openai.String("expires_at"),
					Refresh: openai.CredentialAuthRotateParamMcpOAuthRefresh{
						RefreshToken: openai.String("refresh_token"),
						Scope:        openai.String("scope"),
						TokenEndpointAuth: openai.McpOAuthTokenEndpointAuthRotateParamUnion{
							OfParamClientSecretBasic: &openai.McpOAuthTokenEndpointAuthRotateParamClientSecretBasic{
								ClientSecret: openai.String("client_secret"),
							},
						},
					},
				},
			},
		},
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentVaultCredentialListWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Beta.Agents.Vaults.Credentials.List(
		context.TODO(),
		"vault_id",
		openai.BetaAgentVaultCredentialListParams{
			After: openai.String("after"),
			Limit: openai.Int(0),
			Order: openai.BetaAgentVaultCredentialListParamsOrderAsc,
			Status: openai.VaultStatusFilterUnionParam{
				OfStatus: openai.Opt(openai.VaultStatusActive),
			},
		},
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentVaultCredentialDelete(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Beta.Agents.Vaults.Credentials.Delete(
		context.TODO(),
		"vault_id",
		"credential_id",
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
