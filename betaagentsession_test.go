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

func TestBetaAgentSessionNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.Sessions.New(context.TODO(), openai.BetaAgentSessionNewParams{
		Environment: openai.EnvironmentParamUnion{
			OfParamNone: &openai.EnvironmentParamNone{},
		},
		Agent: openai.BetaAgentSessionNewParamsAgent{
			Instructions: openai.String("instructions"),
			Model:        openai.String("model"),
			MultiAgent: openai.MultiAgentConfigParam{
				Enabled:                true,
				MaxConcurrentSubagents: openai.Int(1),
			},
			Reasoning: openai.AgentReasoningParam{
				Effort:  openai.AgentReasoningParamEffortNone,
				Summary: openai.AgentReasoningParamSummaryConcise,
			},
			ServiceTier: "auto",
			Text: openai.AgentTextParam{
				Format: openai.TextFormatParamUnion{
					OfParamText: &openai.TextFormatParamText{},
				},
				Verbosity: openai.AgentTextParamVerbosityLow,
			},
			Tools: []openai.AgentToolParamUnion{{
				OfParamFunction: &openai.AgentToolParamFunction{
					Description: "description",
					Name:        "name",
					Parameters: map[string]any{
						"foo": "bar",
					},
					DeferLoading: openai.Bool(true),
				},
			}},
		},
		AgentID: openai.String("agent_id"),
		Input: openai.BetaAgentSessionNewParamsInputUnion{
			OfString: openai.String("string"),
		},
		Metadata: map[string]string{
			"foo": "string",
		},
		VaultIDs: []string{"string"},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentSessionGet(t *testing.T) {
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
	_, err := client.Beta.Agents.Sessions.Get(context.TODO(), "session_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentSessionUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.Sessions.Update(
		context.TODO(),
		"session_id",
		openai.BetaAgentSessionUpdateParams{
			Metadata: map[string]string{
				"foo": "string",
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

func TestBetaAgentSessionListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.Sessions.List(context.TODO(), openai.BetaAgentSessionListParams{
		After:   openai.String("after"),
		AgentID: openai.String("agent_id"),
		Limit:   openai.Int(1),
		Order:   openai.BetaAgentSessionListParamsOrderAsc,
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentSessionDelete(t *testing.T) {
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
	_, err := client.Beta.Agents.Sessions.Delete(context.TODO(), "session_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
