// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/internal/testutil"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

func TestChatCompletionNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfDeveloper: &openai.ChatCompletionDeveloperMessageParam{
				Content: []openai.ChatCompletionContentPartTextParam{{Text: openai.String("text")}},
				Name:    openai.String("name"),
			},
		}},
		Model: openai.ChatModelO3Mini,
		Audio: openai.ChatCompletionAudioParam{
			Format: "wav",
			Voice:  "alloy",
		},
		FrequencyPenalty: openai.Float(-2),
		FunctionCall: openai.ChatCompletionNewParamsFunctionCallUnion{
			OfAuto: "none",
		},
		Functions: []openai.ChatCompletionNewParamsFunction{{
			Name:        openai.String("name"),
			Description: openai.String("description"),
			Parameters: shared.FunctionParameters{
				"foo": "bar",
			},
		}},
		LogitBias: map[string]int64{
			"foo": 0,
		},
		Logprobs:            openai.Bool(true),
		MaxCompletionTokens: openai.Int(0),
		MaxTokens:           openai.Int(0),
		Metadata: shared.MetadataParam{
			"foo": "string",
		},
		Modalities:        []openai.ChatCompletionModality{openai.ChatCompletionModalityText},
		N:                 openai.Int(1),
		ParallelToolCalls: openai.Bool(true),
		Prediction: openai.ChatCompletionPredictionContentParam{
			Content: []openai.ChatCompletionContentPartTextParam{{Text: openai.String("text")}},
		},
		PresencePenalty: openai.Float(-2),
		ReasoningEffort: openai.ChatCompletionReasoningEffortLow,
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfResponseFormatText: &shared.ResponseFormatTextParam{},
		},
		Seed:        openai.Int(0),
		ServiceTier: openai.ChatCompletionNewParamsServiceTierAuto,
		Stop: openai.ChatCompletionNewParamsStopUnion{
			OfString: openai.String("string"),
		},
		Store: openai.Bool(true),
		StreamOptions: openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: openai.Bool(true),
		},
		Temperature: openai.Float(1),
		ToolChoice: openai.ChatCompletionToolChoiceOptionUnionParam{
			OfAuto: "none",
		},
		Tools: []openai.ChatCompletionToolParam{{
			Function: shared.FunctionDefinitionParam{
				Name:        openai.String("name"),
				Description: openai.String("description"),
				Parameters: shared.FunctionParameters{
					"foo": "bar",
				},
				Strict: openai.Bool(true),
			},
		}},
		TopLogprobs: openai.Int(0),
		TopP:        openai.Float(1),
		User:        openai.String("user-1234"),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestChatCompletionGet(t *testing.T) {
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
	_, err := client.Chat.Completions.Get(context.TODO(), "completion_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestChatCompletionUpdate(t *testing.T) {
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
	_, err := client.Chat.Completions.Update(
		context.TODO(),
		"completion_id",
		openai.ChatCompletionUpdateParams{
			Metadata: shared.MetadataParam{
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

func TestChatCompletionListWithOptionalParams(t *testing.T) {
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
	_, err := client.Chat.Completions.List(context.TODO(), openai.ChatCompletionListParams{
		After: openai.String("after"),
		Limit: openai.Int(0),
		Metadata: shared.MetadataParam{
			"foo": "string",
		},
		Model: openai.String("model"),
		Order: openai.ChatCompletionListParamsOrderAsc,
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestChatCompletionDelete(t *testing.T) {
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
	_, err := client.Chat.Completions.Delete(context.TODO(), "completion_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
