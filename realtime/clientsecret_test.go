// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

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
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared/constant"
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
				Audio: realtime.RealtimeAudioConfigParam{
					Input: realtime.RealtimeAudioConfigInputParam{
						Format: realtime.RealtimeAudioFormatsUnionParam{
							OfAudioPCM: &realtime.RealtimeAudioFormatsAudioPCMParam{
								Rate: 24000,
								Type: "audio/pcm",
							},
						},
						NoiseReduction: realtime.RealtimeAudioConfigInputNoiseReductionParam{
							Type: realtime.NoiseReductionTypeNearField,
						},
						Transcription: realtime.AudioTranscriptionParam{
							Language: openai.String("language"),
							Model:    realtime.AudioTranscriptionModelWhisper1,
							Prompt:   openai.String("prompt"),
						},
						TurnDetection: realtime.RealtimeAudioInputTurnDetectionUnionParam{
							OfServerVad: &realtime.RealtimeAudioInputTurnDetectionServerVadParam{
								CreateResponse:    openai.Bool(true),
								IdleTimeoutMs:     openai.Int(5000),
								InterruptResponse: openai.Bool(true),
								PrefixPaddingMs:   openai.Int(0),
								SilenceDurationMs: openai.Int(0),
								Threshold:         openai.Float(0),
							},
						},
					},
					Output: realtime.RealtimeAudioConfigOutputParam{
						Format: realtime.RealtimeAudioFormatsUnionParam{
							OfAudioPCM: &realtime.RealtimeAudioFormatsAudioPCMParam{
								Rate: 24000,
								Type: "audio/pcm",
							},
						},
						Speed: openai.Float(0.25),
						Voice: realtime.RealtimeAudioConfigOutputVoiceAlloy,
					},
				},
				Include:      []string{"item.input_audio_transcription.logprobs"},
				Instructions: openai.String("instructions"),
				MaxOutputTokens: realtime.RealtimeSessionCreateRequestMaxOutputTokensUnionParam{
					OfInt: openai.Int(0),
				},
				Model:            realtime.RealtimeSessionCreateRequestModelGPTRealtime,
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
				ToolChoice: realtime.RealtimeToolChoiceConfigUnionParam{
					OfToolChoiceMode: openai.Opt(responses.ToolChoiceOptionsNone),
				},
				Tools: realtime.RealtimeToolsConfigParam{realtime.RealtimeToolsConfigUnionParam{
					OfFunction: &realtime.RealtimeFunctionToolParam{
						Description: openai.String("description"),
						Name:        openai.String("name"),
						Parameters:  map[string]any{},
						Type:        realtime.RealtimeFunctionToolTypeFunction,
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
