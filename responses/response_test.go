// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package responses_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/internal/testutil"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

func TestResponseNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Responses.New(context.TODO(), responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("string"),
		},
		Model:           shared.ResponsesModel("gpt-4o"),
		Background:      openai.Bool(true),
		Include:         []responses.ResponseIncludable{responses.ResponseIncludableFileSearchCallResults},
		Instructions:    openai.String("instructions"),
		MaxOutputTokens: openai.Int(0),
		Metadata: shared.Metadata{
			"foo": "string",
		},
		ParallelToolCalls:  openai.Bool(true),
		PreviousResponseID: openai.String("previous_response_id"),
		Reasoning: shared.ReasoningParam{
			Effort:          shared.ReasoningEffortLow,
			GenerateSummary: shared.ReasoningGenerateSummaryAuto,
			Summary:         shared.ReasoningSummaryAuto,
		},
		ServiceTier: responses.ResponseNewParamsServiceTierAuto,
		Store:       openai.Bool(true),
		Temperature: openai.Float(1),
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfText: &shared.ResponseFormatTextParam{},
			},
		},
		ToolChoice: responses.ResponseNewParamsToolChoiceUnion{
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
		TopP:       openai.Float(1),
		Truncation: responses.ResponseNewParamsTruncationAuto,
		User:       openai.String("user-1234"),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestResponseGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Responses.Get(
		context.TODO(),
		"resp_677efb5139a88190b512bc3fef8e535d",
		responses.ResponseGetParams{
			Include:       []responses.ResponseIncludable{responses.ResponseIncludableFileSearchCallResults},
			StartingAfter: openai.Int(0),
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

func TestResponseDelete(t *testing.T) {
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
	err := client.Responses.Delete(context.TODO(), "resp_677efb5139a88190b512bc3fef8e535d")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestResponseCancel(t *testing.T) {
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
	_, err := client.Responses.Cancel(context.TODO(), "resp_677efb5139a88190b512bc3fef8e535d")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
