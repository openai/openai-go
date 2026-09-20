// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package live_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal/testutil"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared/constant"
)

func TestLiveNewWithOptionalParams(t *testing.T) {
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
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Live.New(context.TODO(), live.LiveNewParams{
		Session: live.MediaSessionConfigParam{
			Model: live.MediaSessionConfigModelGPTLive1,
			Audio: live.MediaSessionConfigAudioParam{
				Output: live.MediaSessionConfigAudioOutputParam{
					Voice: live.MediaSessionConfigAudioOutputVoiceUnionParam{
						OfBuiltIn: openai.Opt(live.BuiltInVoiceAlloy),
					},
				},
			},
			Client: live.ClientConfigParam{
				DataChannel: live.DataChannelConfigParam{
					AllowedClientEvents: live.DataChannelConfigAllowedClientEventsUnionParam{
						OfAll: constant.ValueOf[constant.All](),
					},
					AllowedServerEvents: live.DataChannelConfigAllowedServerEventsUnionParam{
						OfAll: constant.ValueOf[constant.All](),
					},
				},
			},
			Delegation: live.MediaSessionConfigDelegationUnionParam{
				OfClient: &live.ClientDelegationParam{},
			},
			Input: []live.InitialItemUnionParam{{
				OfDeveloper: &live.InitialItemDeveloperParam{
					Content: []live.InitialItemDeveloperContentParam{{
						Text: "text",
						Type: "input_text",
					}},
					ID:     openai.String("id"),
					Status: "incomplete",
					Type:   "message",
				},
			}},
			Instructions: openai.String("instructions"),
			Store:        openai.Bool(true),
		},
		Transport: live.LiveNewParamsTransport{
			Sdp: "x",
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
