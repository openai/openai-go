// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/internal/testutil"
	"github.com/openai/openai-go/option"
)

func TestAudioTranscriptionNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Audio.Transcriptions.New(context.TODO(), openai.AudioTranscriptionNewParams{
		File:                   openai.F(io.Reader(bytes.NewBuffer([]byte("some file contents")))),
		Model:                  openai.F(openai.AudioModelWhisper1),
		Language:               openai.F("language"),
		Prompt:                 openai.F("prompt"),
		ResponseFormat:         openai.F(openai.AudioTranscriptionNewParamsResponseFormatJSON),
		Temperature:            openai.F(0.000000),
		TimestampGranularities: openai.F([]openai.AudioTranscriptionNewParamsTimestampGranularity{openai.AudioTranscriptionNewParamsTimestampGranularityWord, openai.AudioTranscriptionNewParamsTimestampGranularitySegment}),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
