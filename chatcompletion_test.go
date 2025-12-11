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
	"github.com/openai/openai-go/v3/shared"
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
				Content: openai.ChatCompletionDeveloperMessageParamContentUnion{
					OfString: openai.String("string"),
				},
				Name: openai.String("name"),
			},
		}},
		Model: shared.ChatModelGPT5_2,
		Audio: openai.ChatCompletionAudioParam{
			Format: openai.ChatCompletionAudioParamFormatWAV,
			Voice:  openai.ChatCompletionAudioParamVoiceAlloy,
		},
		FrequencyPenalty: openai.Float(-2),
		FunctionCall: openai.ChatCompletionNewParamsFunctionCallUnion{
			OfFunctionCallMode: openai.String("none"),
		},
		Functions: []openai.ChatCompletionNewParamsFunction{{
			Name:        "name",
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
		Metadata: shared.Metadata{
			"foo": "string",
		},
		Modalities:        []string{"text"},
		N:                 openai.Int(1),
		ParallelToolCalls: openai.Bool(true),
		Prediction: openai.ChatCompletionPredictionContentParam{
			Content: openai.ChatCompletionPredictionContentContentUnionParam{
				OfString: openai.String("string"),
			},
		},
		PresencePenalty:      openai.Float(-2),
		PromptCacheKey:       openai.String("prompt-cache-key-1234"),
		PromptCacheRetention: openai.ChatCompletionNewParamsPromptCacheRetentionInMemory,
		ReasoningEffort:      shared.ReasoningEffortNone,
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfText: &shared.ResponseFormatTextParam{},
		},
		SafetyIdentifier: openai.String("safety-identifier-1234"),
		Seed:             openai.Int(-9007199254740991),
		ServiceTier:      openai.ChatCompletionNewParamsServiceTierAuto,
		Stop: openai.ChatCompletionNewParamsStopUnion{
			OfString: openai.String("\n"),
		},
		Store: openai.Bool(true),
		StreamOptions: openai.ChatCompletionStreamOptionsParam{
			IncludeObfuscation: openai.Bool(true),
			IncludeUsage:       openai.Bool(true),
		},
		Temperature: openai.Float(1),
		ToolChoice: openai.ChatCompletionToolChoiceOptionUnionParam{
			OfAuto: openai.String("none"),
		},
		Tools: []openai.ChatCompletionToolUnionParam{{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: shared.FunctionDefinitionParam{
					Name:        "name",
					Description: openai.String("description"),
					Parameters: shared.FunctionParameters{
						"foo": "bar",
					},
					Strict: openai.Bool(true),
				},
			},
		}},
		TopLogprobs: openai.Int(0),
		TopP:        openai.Float(1),
		User:        openai.String("user-1234"),
		Verbosity:   openai.ChatCompletionNewParamsVerbosityLow,
		WebSearchOptions: openai.ChatCompletionNewParamsWebSearchOptions{
			SearchContextSize: "low",
			UserLocation: openai.ChatCompletionNewParamsWebSearchOptionsUserLocation{
				Approximate: openai.ChatCompletionNewParamsWebSearchOptionsUserLocationApproximate{
					City:     openai.String("city"),
					Country:  openai.String("country"),
					Region:   openai.String("region"),
					Timezone: openai.String("timezone"),
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
			Metadata: shared.Metadata{
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
		Metadata: shared.Metadata{
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
