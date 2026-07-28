// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

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

func TestBetaResponseInputTokenCountWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Responses.InputTokens.Count(context.TODO(), openai.BetaResponseInputTokenCountParams{
		Conversation: openai.BetaResponseInputTokenCountParamsConversationUnion{
			OfString: openai.String("string"),
		},
		Input: openai.BetaResponseInputTokenCountParamsInputUnion{
			OfString: openai.String("string"),
		},
		Instructions:       openai.String("instructions"),
		Model:              openai.String("model"),
		ParallelToolCalls:  openai.Bool(true),
		Personality:        openai.BetaResponseInputTokenCountParamsPersonalityFriendly,
		PreviousResponseID: openai.String("resp_123"),
		Reasoning: openai.BetaResponseInputTokenCountParamsReasoning{
			Context:         "auto",
			Effort:          "none",
			GenerateSummary: "auto",
			Mode:            "standard",
			Summary:         "auto",
		},
		Text: openai.BetaResponseInputTokenCountParamsText{
			Format: openai.BetaResponseFormatTextConfigUnionParam{
				OfText: &openai.BetaResponseFormatTextConfigTextParam{},
			},
			Verbosity: "low",
		},
		ToolChoice: openai.BetaResponseInputTokenCountParamsToolChoiceUnion{
			OfToolChoiceMode: openai.Opt(openai.BetaToolChoiceOptionsNone),
		},
		Tools: []openai.BetaToolUnionParam{{
			OfFunction: &openai.BetaFunctionToolParam{
				Name: "name",
				Parameters: map[string]any{
					"foo": "bar",
				},
				Strict:         openai.Bool(true),
				AllowedCallers: []string{"direct"},
				DeferLoading:   openai.Bool(true),
				Description:    openai.String("description"),
				OutputSchema: map[string]any{
					"foo": "bar",
				},
			},
		}},
		Truncation: openai.BetaResponseInputTokenCountParamsTruncationAuto,
		Betas:      []string{"responses_multi_agent=v1"},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
