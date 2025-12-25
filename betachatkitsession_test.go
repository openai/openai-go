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

func TestBetaChatKitSessionNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.ChatKit.Sessions.New(context.TODO(), openai.BetaChatKitSessionNewParams{
		User: "x",
		Workflow: openai.ChatSessionWorkflowParam{
			ID: "id",
			StateVariables: map[string]openai.ChatSessionWorkflowParamStateVariableUnion{
				"foo": {
					OfString: openai.String("string"),
				},
			},
			Tracing: openai.ChatSessionWorkflowParamTracing{
				Enabled: openai.Bool(true),
			},
			Version: openai.String("version"),
		},
		ChatKitConfiguration: openai.ChatSessionChatKitConfigurationParam{
			AutomaticThreadTitling: openai.ChatSessionChatKitConfigurationParamAutomaticThreadTitling{
				Enabled: openai.Bool(true),
			},
			FileUpload: openai.ChatSessionChatKitConfigurationParamFileUpload{
				Enabled:     openai.Bool(true),
				MaxFileSize: openai.Int(1),
				MaxFiles:    openai.Int(1),
			},
			History: openai.ChatSessionChatKitConfigurationParamHistory{
				Enabled:       openai.Bool(true),
				RecentThreads: openai.Int(1),
			},
		},
		ExpiresAfter: openai.ChatSessionExpiresAfterParam{
			Seconds: 1,
		},
		RateLimits: openai.ChatSessionRateLimitsParam{
			MaxRequestsPer1Minute: openai.Int(1),
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

func TestBetaChatKitSessionCancel(t *testing.T) {
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
	_, err := client.Beta.ChatKit.Sessions.Cancel(context.TODO(), "cksess_123")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
