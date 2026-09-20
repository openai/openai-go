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

func TestBetaAgentEnvironmentTemplateNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.Environments.Templates.New(context.TODO(), openai.BetaAgentEnvironmentTemplateNewParams{
		CapabilityDirectories: []string{"string"},
		Env: map[string]string{
			"foo": "string",
		},
		Files: []openai.HostedEnvironmentFileParamUnion{{
			OfParamFileID: &openai.HostedEnvironmentFileParamFileID{
				FileID: "x",
				Path:   "x",
			},
		}},
		Name: openai.String("x"),
		Network: openai.BetaAgentEnvironmentTemplateNewParamsNetwork{
			Access:         "enabled",
			AllowedDomains: []string{"string"},
		},
		Packages: openai.BetaAgentEnvironmentTemplateNewParamsPackages{
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
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentEnvironmentTemplateGet(t *testing.T) {
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
	_, err := client.Beta.Agents.Environments.Templates.Get(context.TODO(), "environment_template_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentEnvironmentTemplateUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.Environments.Templates.Update(
		context.TODO(),
		"environment_template_id",
		openai.BetaAgentEnvironmentTemplateUpdateParams{
			CapabilityDirectories: []string{"string"},
			Env: map[string]string{
				"foo": "string",
			},
			Files: []openai.HostedEnvironmentFileParamUnion{{
				OfParamFileID: &openai.HostedEnvironmentFileParamFileID{
					FileID: "x",
					Path:   "x",
				},
			}},
			Name: openai.String("x"),
			Network: openai.BetaAgentEnvironmentTemplateUpdateParamsNetwork{
				Access:         "enabled",
				AllowedDomains: []string{"string"},
			},
			Packages: openai.BetaAgentEnvironmentTemplateUpdateParamsPackages{
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
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentEnvironmentTemplateListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.Environments.Templates.List(context.TODO(), openai.BetaAgentEnvironmentTemplateListParams{
		After: openai.String("after"),
		Limit: openai.Int(1),
		Order: openai.BetaAgentEnvironmentTemplateListParamsOrderAsc,
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentEnvironmentTemplateDelete(t *testing.T) {
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
	_, err := client.Beta.Agents.Environments.Templates.Delete(context.TODO(), "environment_template_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
