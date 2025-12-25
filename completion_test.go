// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Nordlys-Labs/openai-go/v3"
	"github.com/Nordlys-Labs/openai-go/v3/internal/testutil"
	"github.com/Nordlys-Labs/openai-go/v3/option"
)

func TestCompletionNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Completions.New(context.TODO(), openai.CompletionNewParams{
		Model: openai.CompletionNewParamsModelGPT3_5TurboInstruct,
		Prompt: openai.CompletionNewParamsPromptUnion{
			OfString: openai.String("This is a test."),
		},
		BestOf:           openai.Int(0),
		Echo:             openai.Bool(true),
		FrequencyPenalty: openai.Float(-2),
		LogitBias: map[string]int64{
			"foo": 0,
		},
		Logprobs:        openai.Int(0),
		MaxTokens:       openai.Int(16),
		N:               openai.Int(1),
		PresencePenalty: openai.Float(-2),
		Seed:            openai.Int(0),
		Stop: openai.CompletionNewParamsStopUnion{
			OfString: openai.String("\n"),
		},
		StreamOptions: openai.ChatCompletionStreamOptionsParam{
			IncludeObfuscation: openai.Bool(true),
			IncludeUsage:       openai.Bool(true),
		},
		Suffix:      openai.String("test."),
		Temperature: openai.Float(1),
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
