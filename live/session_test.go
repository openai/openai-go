// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package live_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal/testutil"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared/constant"
)

func TestSessionAcceptWithOptionalParams(t *testing.T) {
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
	err := client.Live.Sessions.Accept(
		context.TODO(),
		"session_id",
		live.SessionAcceptParams{
			Session: live.SessionAcceptParamsSession{
				Model: "gpt-live-1",
				Audio: live.SessionAcceptParamsSessionAudio{
					Output: live.SessionAcceptParamsSessionAudioOutput{
						Voice: live.SessionAcceptParamsSessionAudioOutputVoiceUnion{
							OfBuiltIn: openai.Opt(live.BuiltInVoiceAlloy),
						},
					},
				},
				Delegation: live.SessionAcceptParamsSessionDelegationUnion{
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

func TestSessionDownloadRecording(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		if _, err := w.Write([]byte("abc")); err != nil {
			t.Errorf("write response body: %v", err)
		}
	}))
	defer server.Close()
	baseURL := server.URL
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	resp, err := client.Live.Sessions.DownloadRecording(context.TODO(), "live_SQ")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			t.Errorf("close response body: %v", closeErr)
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
	if !bytes.Equal(b, []byte("abc")) {
		t.Fatalf("return value not %s: %s", "abc", b)
	}
}

func TestSessionForkWithOptionalParams(t *testing.T) {
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
	_, err := client.Live.Sessions.Fork(
		context.TODO(),
		"session_id",
		live.SessionForkParams{
			Transport: live.SessionForkParamsTransport{
				Sdp: "x",
			},
			Session: live.MediaSessionForkConfigParam{
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
				Delegation: live.MediaSessionForkConfigDelegationParam{
					Responses: live.ResponsesDelegationUpdateConfigParam{
						Instructions:      openai.String("instructions"),
						MaxOutputTokens:   openai.Int(16),
						Model:             openai.String("model"),
						ParallelToolCalls: openai.Bool(true),
						Reasoning: live.ResponsesDelegationUpdateConfigReasoningParam{
							Effort:  "none",
							Summary: "concise",
						},
						ServiceTier: live.ResponsesDelegationUpdateConfigServiceTierAuto,
						Text: live.ResponsesDelegationUpdateConfigTextParam{
							Verbosity: "low",
						},
						ToolChoice: live.ResponsesDelegationUpdateConfigToolChoiceUnionParam{
							OfLiveToolChoiceEnum: openai.String("auto"),
						},
						Tools: []live.ResponsesDelegationUpdateConfigToolUnionParam{{
							OfFunction: &live.FunctionToolParam{
								Name:        "name",
								Description: openai.String("description"),
								Parameters: map[string]any{
									"foo": "bar",
								},
								Strict: openai.Bool(true),
							},
						}},
					},
				},
				Store: openai.Bool(true),
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

func TestSessionHangup(t *testing.T) {
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
	err := client.Live.Sessions.Hangup(context.TODO(), "session_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSessionRefer(t *testing.T) {
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
	err := client.Live.Sessions.Refer(
		context.TODO(),
		"session_id",
		live.SessionReferParams{
			TargetUri: "tel:+14155550123",
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

func TestSessionReject(t *testing.T) {
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
	err := client.Live.Sessions.Reject(
		context.TODO(),
		"session_id",
		live.SessionRejectParams{
			StatusCode: 486,
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
