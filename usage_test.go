// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"os"
	"testing"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/internal/testutil"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

func TestUsage(t *testing.T) {
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
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{openai.ChatCompletionUserMessageParam{
			Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
			Content: openai.F[openai.ChatCompletionUserMessageParamContentUnion](shared.UnionString("Say this is a test")),
		}}),
		Model: openai.F(shared.ChatModelO3Mini),
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v\n", chatCompletion)
}
