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

func TestBetaAgentNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.New(context.TODO(), openai.BetaAgentNewParams{
		Model:        "model",
		Instructions: openai.String("instructions"),
		Metadata: map[string]string{
			"foo": "string",
		},
		MultiAgent: openai.MultiAgentConfigParam{
			Enabled:                true,
			MaxConcurrentSubagents: openai.Int(1),
		},
		Name: openai.String("name"),
		Reasoning: openai.AgentReasoningParam{
			Effort:  openai.AgentReasoningParamEffortNone,
			Summary: openai.AgentReasoningParamSummaryConcise,
		},
		ServiceTier: openai.BetaAgentNewParamsServiceTierAuto,
		Text: openai.AgentTextParam{
			Format: openai.TextFormatParamUnion{
				OfParamText: &openai.TextFormatParamText{},
			},
			Verbosity: openai.AgentTextParamVerbosityLow,
		},
		Tools: []openai.PersistedAgentToolParamUnion{{
			OfParamFunction: &openai.PersistedAgentToolParamFunction{
				Description: "description",
				Name:        "name",
				Parameters: map[string]any{
					"foo": "bar",
				},
				DeferLoading: openai.Bool(true),
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

func TestBetaAgentGet(t *testing.T) {
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
	_, err := client.Beta.Agents.Get(context.TODO(), "agent_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.Update(
		context.TODO(),
		"agent_id",
		openai.BetaAgentUpdateParams{
			Instructions: openai.String("instructions"),
			Metadata: map[string]string{
				"foo": "string",
			},
			Model: openai.String("model"),
			MultiAgent: openai.MultiAgentConfigParam{
				Enabled:                true,
				MaxConcurrentSubagents: openai.Int(1),
			},
			Name: openai.String("name"),
			Reasoning: openai.AgentReasoningParam{
				Effort:  openai.AgentReasoningParamEffortNone,
				Summary: openai.AgentReasoningParamSummaryConcise,
			},
			ServiceTier: openai.BetaAgentUpdateParamsServiceTierAuto,
			Text: openai.AgentTextParam{
				Format: openai.TextFormatParamUnion{
					OfParamText: &openai.TextFormatParamText{},
				},
				Verbosity: openai.AgentTextParamVerbosityLow,
			},
			Tools: []openai.PersistedAgentToolParamUnion{{
				OfParamFunction: &openai.PersistedAgentToolParamFunction{
					Description: "description",
					Name:        "name",
					Parameters: map[string]any{
						"foo": "bar",
					},
					DeferLoading: openai.Bool(true),
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

func TestBetaAgentListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.List(context.TODO(), openai.BetaAgentListParams{
		After: openai.String("after"),
		Limit: openai.Int(1),
		Order: openai.BetaAgentListParamsOrderAsc,
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentDelete(t *testing.T) {
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
	_, err := client.Beta.Agents.Delete(context.TODO(), "agent_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
