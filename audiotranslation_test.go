// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/internal/testutil"
	"github.com/openai/openai-go/v2/option"
)

func TestAudioTranslationNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Audio.Translations.New(context.TODO(), openai.AudioTranslationNewParams{
		File:           io.Reader(bytes.NewBuffer([]byte("some file contents"))),
		Model:          openai.AudioModelWhisper1,
		Prompt:         openai.String("prompt"),
		ResponseFormat: openai.AudioTranslationNewParamsResponseFormatJSON,
		Temperature:    openai.Float(0),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
