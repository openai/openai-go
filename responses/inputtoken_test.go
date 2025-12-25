// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package responses_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Nordlys-Labs/openai-go/v3"
	"github.com/Nordlys-Labs/openai-go/v3/internal/testutil"
	"github.com/Nordlys-Labs/openai-go/v3/option"
	"github.com/Nordlys-Labs/openai-go/v3/responses"
	"github.com/Nordlys-Labs/openai-go/v3/shared"
)

func TestInputTokenCountWithOptionalParams(t *testing.T) {
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
	)
	_, err := client.Responses.InputTokens.Count(context.TODO(), responses.InputTokenCountParams{
		Conversation: responses.InputTokenCountParamsConversationUnion{
			OfString: openai.String("string"),
		},
		Input: responses.InputTokenCountParamsInputUnion{
			OfString: openai.String("string"),
		},
		Instructions:       openai.String("instructions"),
		Model:              openai.String("model"),
		ParallelToolCalls:  openai.Bool(true),
		PreviousResponseID: openai.String("resp_123"),
		Reasoning: shared.ReasoningParam{
			Effort:          shared.ReasoningEffortNone,
			GenerateSummary: shared.ReasoningGenerateSummaryAuto,
			Summary:         shared.ReasoningSummaryAuto,
		},
		Text: responses.InputTokenCountParamsText{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfText: &shared.ResponseFormatTextParam{},
			},
			Verbosity: "low",
		},
		ToolChoice: responses.InputTokenCountParamsToolChoiceUnion{
			OfToolChoiceMode: openai.Opt(responses.ToolChoiceOptionsNone),
		},
		Tools: []responses.ToolUnionParam{{
			OfFunction: &responses.FunctionToolParam{
				Name: "name",
				Parameters: map[string]any{
					"foo": "bar",
				},
				Strict:      openai.Bool(true),
				Description: openai.String("description"),
			},
		}},
		Truncation: responses.InputTokenCountParamsTruncationAuto,
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
