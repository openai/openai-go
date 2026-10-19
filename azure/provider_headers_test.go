package azure

import (
	"context"
	"net/http"
	"slices"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestAzureProviderSwitchIsolatesHeaders(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "fake-openai-key")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	t.Setenv("OPENAI_ORG_ID", "fake-org")
	t.Setenv("OPENAI_PROJECT_ID", "fake-project")
	t.Setenv("OPENAI_CUSTOM_HEADERS", "X-Api-Key: fake-secret\nCookie: fake-cookie\nX-Amz-Security-Token: fake-token\nX-Custom-Credential: fake-custom\nX-Reapplied: inherited\nX-Added: inherited\nX-Deleted: inherited")
	for _, mode := range []string{"key", "token"} {
		for _, layer := range []string{"request", "service"} {
			for _, order := range []string{"before", "after"} {
				t.Run(mode+"/"+layer+"/"+order, func(t *testing.T) {
					var received http.Header
					client := openai.NewClient(
						option.WithHeader("X-Client-Secret", "fake-client-secret"),
						option.WithHeaderAdd("Accept", "fake-inherited-accept"),
						option.WithHeader("User-Agent", "fake-inherited-agent"),
						option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
							received = req.Header.Clone()
							return successfulAzureResponse(req), nil
						})}),
					)
					explicit := []option.RequestOption{
						option.WithHeader("x-reapplied", "explicit"),
						option.WithHeaderAdd("x-added", "explicit-one"),
						option.WithHeaderAdd("X-Added", "explicit-two"),
						option.WithHeaderDel("X-Deleted"),
						option.WithHeaderAdd("Accept", "application/custom"),
					}
					auth := WithAPIKey("fake-azure-key")
					if mode == "token" {
						auth = WithTokenCredential(&fake.TokenCredential{})
					}
					opts := []option.RequestOption{WithEndpoint("https://resource.openai.azure.com", "2024-10-21"), auth}
					if order == "before" {
						opts = append(explicit, opts...)
					} else {
						opts = append(opts, explicit...)
					}
					if layer == "request" {
						_, err := client.Models.List(context.Background(), opts...)
						if err != nil {
							t.Fatal(err)
						}
					} else {
						service := openai.NewModelService(append(slices.Clone(client.Options), opts...)...)
						_, err := service.List(context.Background())
						if err != nil {
							t.Fatal(err)
						}
					}
					if received == nil {
						t.Fatal("transport was not called")
					}
					for _, key := range []string{"X-Api-Key", "Cookie", "X-Amz-Security-Token", "X-Custom-Credential", "X-Client-Secret", "OpenAI-Organization", "OpenAI-Project", "X-Deleted"} {
						if len(received.Values(key)) != 0 {
							t.Errorf("inherited %s reached Azure", key)
						}
					}
					if got := received.Values("X-Reapplied"); !slices.Equal(got, []string{"explicit"}) {
						t.Errorf("reapplied header = %v", got)
					}
					if got := received.Values("X-Added"); !slices.Equal(got, []string{"explicit-one", "explicit-two"}) {
						t.Errorf("added header = %v", got)
					}
					if !slices.Equal(received.Values("Accept"), []string{"application/json", "application/custom"}) || received.Get("User-Agent") == "" || received.Get("User-Agent") == "fake-inherited-agent" {
						t.Error("SDK request headers missing")
					}
					if mode == "key" {
						if received.Get("Api-Key") != "fake-azure-key" || received.Get("Authorization") != "" {
							t.Error("incorrect Azure key authentication")
						}
					} else if received.Get("Authorization") != "Bearer fake_token" || received.Get("Api-Key") != "" {
						t.Error("incorrect Azure token authentication")
					}

					// Switching a request or service must not mutate the original client.
					_, err := client.Models.List(context.Background())
					if err != nil {
						t.Fatal(err)
					}
					if received.Get("X-Api-Key") != "fake-secret" || received.Get("OpenAI-Organization") != "fake-org" || received.Get("OpenAI-Project") != "fake-project" || received.Get("X-Client-Secret") != "fake-client-secret" {
						t.Error("original OpenAI headers changed")
					}
				})
			}
		}
	}
}

func TestAzureCredentialOverridePreservesProviderHeaders(t *testing.T) {
	var received http.Header
	client := openai.NewClient(
		WithEndpoint("https://resource.openai.azure.com", "2024-10-21"),
		WithAPIKey("fake-original-key"),
		option.WithHeader("X-Azure-Gateway", "fake-gateway-key"),
		option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			received = req.Header.Clone()
			return successfulAzureResponse(req), nil
		})}),
	)
	_, err := client.Models.List(context.Background(), WithTokenCredential(&fake.TokenCredential{}))
	if err != nil {
		t.Fatal(err)
	}
	if received.Get("X-Azure-Gateway") != "fake-gateway-key" || received.Get("Authorization") != "Bearer fake_token" || received.Get("Api-Key") != "" {
		t.Fatal("Azure credential override changed provider headers")
	}
}

func TestAzurePreservesOperationHeaders(t *testing.T) {
	for _, layer := range []string{"client", "service", "request"} {
		t.Run(layer, func(t *testing.T) {
			var received http.Header
			networking := option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				received = req.Header.Clone()
				return successfulAzureResponse(req), nil
			})})
			azureOptions := []option.RequestOption{WithEndpoint("https://resource.openai.azure.com", "2024-10-21"), WithAPIKey("fake-azure-key")}
			opts := []option.RequestOption{networking}
			if layer == "client" {
				opts = append(opts, azureOptions...)
			} else {
				// Inherited overrides must not hide generated defaults after switching.
				opts = append(opts, option.WithHeader("Accept", "fake-inherited"), option.WithHeader("Content-Type", "fake-inherited"), option.WithHeader("OpenAI-Beta", "fake-inherited"))
			}
			client := openai.NewClient(opts...)
			speech := client.Audio.Speech
			vectorStores := client.VectorStores
			var requestOptions []option.RequestOption
			if layer == "service" {
				serviceOptions := append(slices.Clone(client.Options), azureOptions...)
				speech = openai.NewAudioSpeechService(serviceOptions...)
				vectorStores = openai.NewVectorStoreService(serviceOptions...)
			} else if layer == "request" {
				requestOptions = azureOptions
			}
			response, err := speech.New(context.Background(), openai.AudioSpeechNewParams{Input: "test", Model: "tts-1", Voice: openai.AudioSpeechNewParamsVoiceUnion{OfString: openai.String("alloy")}}, requestOptions...)
			if err != nil {
				t.Fatal(err)
			}
			if closeErr := response.Body.Close(); closeErr != nil {
				t.Fatal(closeErr)
			}
			if received.Get("Accept") != "application/octet-stream" || received.Get("Content-Type") != "application/json" {
				t.Fatal("speech operation headers were lost")
			}
			_, err = vectorStores.Get(context.Background(), "vector-store-test", requestOptions...)
			if err != nil {
				t.Fatal(err)
			}
			if received.Get("OpenAI-Beta") != "assistants=v2" {
				t.Fatal("vector store operation header was lost")
			}
		})
	}
}
