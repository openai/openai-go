// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package realtime_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/internal/testutil"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/realtime"
	"github.com/openai/openai-go/v2/responses"
	"github.com/openai/openai-go/v2/shared/constant"
)

func TestClientSecretNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Realtime.ClientSecrets.New(context.TODO(), realtime.ClientSecretNewParams{
		ExpiresAfter: realtime.ClientSecretNewParamsExpiresAfter{
			Anchor:  "created_at",
			Seconds: openai.Int(10),
		},
		Session: realtime.ClientSecretNewParamsSessionUnion{
			OfRealtime: &realtime.RealtimeSessionCreateRequestParam{
				Model: realtime.RealtimeSessionCreateRequestModelGPT4oRealtime,
				Audio: realtime.RealtimeAudioConfigParam{
					Input: realtime.RealtimeAudioConfigInputParam{
						Format: "pcm16",
						NoiseReduction: realtime.RealtimeAudioConfigInputNoiseReductionParam{
							Type: "near_field",
						},
						Transcription: realtime.RealtimeAudioConfigInputTranscriptionParam{
							Language: openai.String("language"),
							Model:    "whisper-1",
							Prompt:   openai.String("prompt"),
						},
						TurnDetection: realtime.RealtimeAudioConfigInputTurnDetectionParam{
							CreateResponse:    openai.Bool(true),
							Eagerness:         "low",
							IdleTimeoutMs:     openai.Int(0),
							InterruptResponse: openai.Bool(true),
							PrefixPaddingMs:   openai.Int(0),
							SilenceDurationMs: openai.Int(0),
							Threshold:         openai.Float(0),
							Type:              "server_vad",
						},
					},
					Output: realtime.RealtimeAudioConfigOutputParam{
						Format: "pcm16",
						Speed:  openai.Float(0.25),
						Voice:  "alloy",
					},
				},
				ClientSecret: realtime.RealtimeClientSecretConfigParam{
					ExpiresAfter: realtime.RealtimeClientSecretConfigExpiresAfterParam{
						Anchor:  "created_at",
						Seconds: openai.Int(0),
					},
				},
				Include:      []string{"item.input_audio_transcription.logprobs"},
				Instructions: openai.String("instructions"),
				MaxOutputTokens: realtime.RealtimeSessionCreateRequestMaxOutputTokensUnionParam{
					OfInt: openai.Int(0),
				},
				OutputModalities: []string{"text"},
				Prompt: responses.ResponsePromptParam{
					ID: "id",
					Variables: map[string]responses.ResponsePromptVariableUnionParam{
						"foo": {
							OfString: openai.String("string"),
						},
					},
					Version: openai.String("version"),
				},
				Temperature: openai.Float(0),
				ToolChoice: realtime.RealtimeToolChoiceConfigUnionParam{
					OfToolChoiceMode: openai.Opt(responses.ToolChoiceOptionsNone),
				},
				Tools: realtime.RealtimeToolsConfigParam{realtime.RealtimeToolsConfigUnionParam{
					OfFunction: &realtime.RealtimeToolsConfigUnionFunctionParam{
						Description: openai.String("description"),
						Name:        openai.String("name"),
						Parameters:  map[string]interface{}{},
						Type:        "function",
					},
				}},
				Tracing: realtime.RealtimeTracingConfigUnionParam{
					OfAuto: constant.ValueOf[constant.Auto](),
				},
				Truncation: realtime.RealtimeTruncationUnionParam{
					OfRealtimeTruncationStrategy: openai.String("auto"),
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
