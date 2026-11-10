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

func TestBetaAgentAutomationNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.Automations.New(context.TODO(), openai.BetaAgentAutomationNewParams{
		CreateAutomationParams: openai.AutomationScheduleParamUnion{
			OfParamPeriodic: &openai.AutomationScheduleParamPeriodic{
				Expression: "x",
				Name:       "x",
				Target: openai.AutomationTargetParamUnion{
					OfParamNewSession: &openai.AutomationTargetParamNewSession{
						Environment: openai.EnvironmentParamUnion{
							OfParamNone: &openai.EnvironmentParamNone{},
						},
						Input: []openai.AgentSessionInputMessageParam{{
							Content: []openai.InputContentParamUnion{{
								OfParamInputText: &openai.InputContentParamInputText{
									Text: "text",
								},
							}},
							Type: openai.AgentSessionInputMessageParamTypeMessage,
						}},
						Agent: openai.AutomationTargetParamNewSessionAgent{
							Compaction: openai.AutomationTargetParamNewSessionAgentCompaction{
								Async: openai.Bool(true),
							},
							Experimental: openai.AutomationTargetParamNewSessionAgentExperimental{
								TimeAwareness: openai.Bool(true),
							},
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
									Async:        openai.Bool(true),
									DeferLoading: openai.Bool(true),
								},
							}},
						},
						AgentID: openai.String("agent_id"),
						Metadata: map[string]string{
							"foo": "string",
						},
						VaultIDs: []string{"string"},
					},
				},
				Timezone: "x",
				Metadata: map[string]string{
					"foo": "string",
				},
			},
		},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentAutomationGet(t *testing.T) {
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
	_, err := client.Beta.Agents.Automations.Get(context.TODO(), "automation_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentAutomationUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.Automations.Update(
		context.TODO(),
		"automation_id",
		openai.BetaAgentAutomationUpdateParams{
			Expression: openai.String("x"),
			Metadata: map[string]string{
				"foo": "string",
			},
			Name:  openai.String("x"),
			RunAt: openai.Int(0),
			Target: openai.AutomationTargetParamUnion{
				OfParamNewSession: &openai.AutomationTargetParamNewSession{
					Environment: openai.EnvironmentParamUnion{
						OfParamNone: &openai.EnvironmentParamNone{},
					},
					Input: []openai.AgentSessionInputMessageParam{{
						Content: []openai.InputContentParamUnion{{
							OfParamInputText: &openai.InputContentParamInputText{
								Text: "text",
							},
						}},
						Type: openai.AgentSessionInputMessageParamTypeMessage,
					}},
					Agent: openai.AutomationTargetParamNewSessionAgent{
						Compaction: openai.AutomationTargetParamNewSessionAgentCompaction{
							Async: openai.Bool(true),
						},
						Experimental: openai.AutomationTargetParamNewSessionAgentExperimental{
							TimeAwareness: openai.Bool(true),
						},
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
								Async:        openai.Bool(true),
								DeferLoading: openai.Bool(true),
							},
						}},
					},
					AgentID: openai.String("agent_id"),
					Metadata: map[string]string{
						"foo": "string",
					},
					VaultIDs: []string{"string"},
				},
			},
			Timezone: openai.String("x"),
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

func TestBetaAgentAutomationListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.Automations.List(context.TODO(), openai.BetaAgentAutomationListParams{
		After:     openai.String("after"),
		Limit:     openai.Int(1),
		Order:     openai.BetaAgentAutomationListParamsOrderAsc,
		SessionID: openai.String("session_id"),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentAutomationDelete(t *testing.T) {
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
	_, err := client.Beta.Agents.Automations.Delete(context.TODO(), "automation_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentAutomationPause(t *testing.T) {
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
	_, err := client.Beta.Agents.Automations.Pause(context.TODO(), "automation_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentAutomationResume(t *testing.T) {
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
	_, err := client.Beta.Agents.Automations.Resume(context.TODO(), "automation_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
