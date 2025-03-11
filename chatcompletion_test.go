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
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{openai.ChatCompletionDeveloperMessageParam{
			Content: openai.F[openai.ChatCompletionDeveloperMessageParamContentUnion](shared.UnionString("string")),
			Role:    openai.F(openai.ChatCompletionDeveloperMessageParamRoleDeveloper),
			Name:    openai.F("name"),
		}}),
		Model: openai.F(shared.ChatModelO3Mini),
		Audio: openai.F(openai.ChatCompletionAudioParam{
			Format: openai.F(openai.ChatCompletionAudioParamFormatWAV),
			Voice:  openai.F(openai.ChatCompletionAudioParamVoiceAlloy),
		}),
		FrequencyPenalty: openai.F(-2.000000),
		FunctionCall:     openai.F[openai.ChatCompletionNewParamsFunctionCallUnion](openai.ChatCompletionNewParamsFunctionCallFunctionCallMode(openai.ChatCompletionNewParamsFunctionCallFunctionCallModeNone)),
		Functions: openai.F([]openai.ChatCompletionNewParamsFunction{{
			Name:        openai.F("name"),
			Description: openai.F("description"),
			Parameters: openai.F(shared.FunctionParameters{
				"foo": "bar",
			}),
		}}),
		LogitBias: openai.F(map[string]int64{
			"foo": int64(0),
		}),
		Logprobs:            openai.F(true),
		MaxCompletionTokens: openai.F(int64(0)),
		MaxTokens:           openai.F(int64(0)),
		Metadata: openai.F(shared.MetadataParam{
			"foo": "string",
		}),
		Modalities:        openai.F([]openai.ChatCompletionNewParamsModality{openai.ChatCompletionNewParamsModalityText}),
		N:                 openai.F(int64(1)),
		ParallelToolCalls: openai.F(true),
		Prediction: openai.F(openai.ChatCompletionPredictionContentParam{
			Content: openai.F[openai.ChatCompletionPredictionContentContentUnionParam](shared.UnionString("string")),
			Type:    openai.F(openai.ChatCompletionPredictionContentTypeContent),
		}),
		PresencePenalty: openai.F(-2.000000),
		ReasoningEffort: openai.F(shared.ReasoningEffortLow),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](shared.ResponseFormatTextParam{
			Type: openai.F(shared.ResponseFormatTextTypeText),
		}),
		Seed:        openai.F(int64(-9007199254740991)),
		ServiceTier: openai.F(openai.ChatCompletionNewParamsServiceTierAuto),
		Stop:        openai.F[openai.ChatCompletionNewParamsStopUnion](shared.UnionString("\n")),
		Store:       openai.F(true),
		StreamOptions: openai.F(openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: openai.F(true),
		}),
		Temperature: openai.F(1.000000),
		ToolChoice:  openai.F[openai.ChatCompletionToolChoiceOptionUnionParam](openai.ChatCompletionToolChoiceOptionAuto(openai.ChatCompletionToolChoiceOptionAutoNone)),
		Tools: openai.F([]openai.ChatCompletionToolParam{{
			Function: openai.F(shared.FunctionDefinitionParam{
				Name:        openai.F("name"),
				Description: openai.F("description"),
				Parameters: openai.F(shared.FunctionParameters{
					"foo": "bar",
				}),
				Strict: openai.F(true),
			}),
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
		}}),
		TopLogprobs: openai.F(int64(0)),
		TopP:        openai.F(1.000000),
		User:        openai.F("user-1234"),
		WebSearchOptions: openai.F(openai.ChatCompletionNewParamsWebSearchOptions{
			SearchContextSize: openai.F(openai.ChatCompletionNewParamsWebSearchOptionsSearchContextSizeLow),
			UserLocation: openai.F(openai.ChatCompletionNewParamsWebSearchOptionsUserLocation{
				Approximate: openai.F(openai.ChatCompletionNewParamsWebSearchOptionsUserLocationApproximate{
					City:     openai.F("city"),
					Country:  openai.F("country"),
					Region:   openai.F("region"),
					Timezone: openai.F("timezone"),
				}),
				Type: openai.F(openai.ChatCompletionNewParamsWebSearchOptionsUserLocationTypeApproximate),
			}),
		}),
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
			Metadata: openai.F(shared.MetadataParam{
				"foo": "string",
			}),
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
		After: openai.F("after"),
		Limit: openai.F(int64(0)),
		Metadata: openai.F(shared.MetadataParam{
			"foo": "string",
		}),
		Model: openai.F("model"),
		Order: openai.F(openai.ChatCompletionListParamsOrderAsc),
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
