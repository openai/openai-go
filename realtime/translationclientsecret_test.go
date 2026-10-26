// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package realtime_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal/testutil"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/realtime"
)

func TestTranslationClientSecretNewWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithUnsafeAllowHTTP(),
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Realtime.Translations.ClientSecrets.New(context.TODO(), realtime.TranslationClientSecretNewParams{
		RealtimeTranslationClientSecretCreateRequest: realtime.RealtimeTranslationClientSecretCreateRequestParam{
			Session: realtime.RealtimeTranslationSessionCreateRequestParam{
				Model: "model",
				Audio: realtime.RealtimeTranslationSessionCreateRequestAudioParam{
					Input: realtime.RealtimeTranslationSessionCreateRequestAudioInputParam{
						NoiseReduction: realtime.RealtimeTranslationSessionCreateRequestAudioInputNoiseReductionParam{
							Type: realtime.NoiseReductionTypeNearField,
						},
						Transcription: realtime.RealtimeTranslationSessionCreateRequestAudioInputTranscriptionParam{
							Model: "model",
						},
					},
					Output: realtime.RealtimeTranslationSessionCreateRequestAudioOutputParam{
						Language: openai.String("language"),
					},
				},
			},
			ExpiresAfter: realtime.RealtimeTranslationClientSecretCreateRequestExpiresAfterParam{
				Anchor:  "created_at",
				Seconds: openai.Int(10),
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
