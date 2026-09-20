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

func TestBetaAgentSessionEventNewWithOptionalParams(t *testing.T) {
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
	err := client.Beta.Agents.Sessions.Events.New(
		context.TODO(),
		"session_id",
		openai.BetaAgentSessionEventNewParams{
			Events: []openai.AgentSessionInputParamUnion{{
				OfParamAgentSessionInputMessage: &openai.AgentSessionInputParamAgentSessionInputMessage{
					Input: []openai.AgentSessionInputMessageParam{{
						Content: []openai.InputContentParamUnion{{
							OfParamInputText: &openai.InputContentParamInputText{
								Text: "text",
							},
						}},
						Type: openai.AgentSessionInputMessageParamTypeMessage,
					}},
				},
			}},
			IdempotencyKey: openai.String("x"),
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
