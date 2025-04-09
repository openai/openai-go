// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"strings"
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
				Content: openai.ChatCompletionDeveloperMessageParamContentUnion{
					OfString: openai.String("string"),
				},
				Name: openai.String("name"),
			},
		}},
		Model: shared.ChatModelO3Mini,
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
		Metadata: shared.MetadataParam{
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
		PresencePenalty: openai.Float(-2),
		ReasoningEffort: shared.ReasoningEffortLow,
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfText: &shared.ResponseFormatTextParam{},
		},
		Seed:        openai.Int(-9007199254740991),
		ServiceTier: openai.ChatCompletionNewParamsServiceTierAuto,
		Stop: openai.ChatCompletionNewParamsStopUnion{
			OfString: openai.String("\n"),
		},
		Store: openai.Bool(true),
		StreamOptions: openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: openai.Bool(true),
		},
		Temperature: openai.Float(1),
		ToolChoice: openai.ChatCompletionToolChoiceOptionUnionParam{
			OfAuto: openai.String("none"),
		},
		Tools: []openai.ChatCompletionToolParam{{
			Function: shared.FunctionDefinitionParam{
				Name:        "name",
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

func TestChatCompletionCustomBaseURL(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv := &http.Server{}

	ready := make(chan struct{})
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.String(), "/openai/v1") {
				t.Errorf("expected prefix to be /openai/v1, got %s", r.URL.String())
			}

			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "completion_id"}`))
		})
		lstr, err := net.Listen("tcp", "localhost:4011")
		if err != nil {
			t.Errorf("net.Listen: %s", err.Error())
		}
		close(ready)
		if err := srv.Serve(lstr); err != http.ErrServerClosed {
			t.Errorf("srv.Serve: %s", err.Error())
		}
	}()
	// Wait until the server is listening
	<-ready

	go func() {
		<-ctx.Done()
		srv.Shutdown(ctx)
	}()

	baseURL := "http://localhost:4011/openai/v1"
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
