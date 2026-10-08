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

func TestBetaAgentEnvironmentNewWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithUnsafeAllowHTTP(),
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Beta.Agents.Environments.New(context.TODO(), openai.BetaAgentEnvironmentNewParams{
		Environment: openai.BetaAgentEnvironmentNewParamsEnvironment{
			CapabilityDirectories: []string{"string"},
			Desktop: openai.BetaAgentEnvironmentNewParamsEnvironmentDesktop{
				Enabled: true,
			},
			Env: map[string]string{
				"foo": "string",
			},
			EnvironmentTemplateID: openai.String("environment_template_id"),
			Files: []openai.HostedEnvironmentFileParamUnion{{
				OfParamFileID: &openai.HostedEnvironmentFileParamFileID{
					FileID: "x",
					Path:   "x",
				},
			}},
			Network: openai.BetaAgentEnvironmentNewParamsEnvironmentNetwork{
				Access:         "enabled",
				AllowedDomains: []string{"string"},
				BlockedDomains: []string{"string"},
			},
			Packages: openai.BetaAgentEnvironmentNewParamsEnvironmentPackages{
				Npm:    []string{"string"},
				Python: []string{"string"},
				System: []string{"string"},
			},
			Plugins: []openai.HostedPluginParam{{
				Description: "description",
				Name:        "x",
				Source: openai.InlineCapabilitySourceParam{
					Data: "x",
				},
			}},
			SetupCommands: []openai.SetupCommandParam{{
				Command: "command",
				Cwd:     openai.String("cwd"),
			}},
			Skills: []openai.HostedSkillParamUnion{{
				OfParamSkillReference: &openai.HostedSkillParamSkillReference{
					SkillID: "x",
					Version: openai.String("version"),
				},
			}},
		},
		VaultIDs:       []string{"string"},
		IdempotencyKey: openai.String("x"),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentEnvironmentGet(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithUnsafeAllowHTTP(),
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Beta.Agents.Environments.Get(context.TODO(), "environment_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentEnvironmentListWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithUnsafeAllowHTTP(),
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Beta.Agents.Environments.List(context.TODO(), openai.BetaAgentEnvironmentListParams{
		After: openai.String("after"),
		Limit: openai.Int(1),
		Order: openai.BetaAgentEnvironmentListParamsOrderAsc,
		Type:  openai.BetaAgentEnvironmentListParamsTypeOpenAIHosted,
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
