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
		Model:            openai.F(openai.CompletionNewParamsModelGPT3_5TurboInstruct),
		Prompt:           openai.F[openai.CompletionNewParamsPromptUnion](shared.UnionString("This is a test.")),
		BestOf:           openai.F(int64(0)),
		Echo:             openai.F(true),
		FrequencyPenalty: openai.F(-2.000000),
		LogitBias: openai.F(map[string]int64{
			"foo": int64(0),
		}),
		Logprobs:        openai.F(int64(0)),
		MaxTokens:       openai.F(int64(16)),
		N:               openai.F(int64(1)),
		PresencePenalty: openai.F(-2.000000),
		Seed:            openai.F(int64(0)),
		Stop:            openai.F[openai.CompletionNewParamsStopUnion](shared.UnionString("\n")),
		StreamOptions: openai.F(openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: openai.F(true),
		}),
		Suffix:      openai.F("test."),
		Temperature: openai.F(1.000000),
		TopP:        openai.F(1.000000),
		User:        openai.F("user-1234"),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
