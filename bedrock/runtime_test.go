package bedrock

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func TestBedrockEndpointResolution(t *testing.T) {
	clearAWSEnvironment(t)
	tests := []struct {
		name     string
		config   Config
		endpoint Endpoint
		baseURL  string
	}{
		{"Mantle default", Config{APIKey: "token", AWSRegion: "us-east-1"}, EndpointMantle, "https://bedrock-mantle.us-east-1.api.aws/openai/v1/"},
		{"Runtime standard", Config{Endpoint: EndpointRuntime, APIKey: "token", AWSRegion: "us-east-1"}, EndpointRuntime, "https://bedrock-runtime.us-east-1.amazonaws.com/openai/v1/"},
		{"Runtime China", Config{Endpoint: EndpointRuntime, APIKey: "token", AWSRegion: "cn-north-1"}, EndpointRuntime, "https://bedrock-runtime.cn-north-1.amazonaws.com.cn/openai/v1/"},
		{"Runtime European sovereign", Config{Endpoint: EndpointRuntime, APIKey: "token", AWSRegion: "eusc-de-east-1"}, EndpointRuntime, "https://bedrock-runtime.eusc-de-east-1.amazonaws.eu/openai/v1/"},
		{"Runtime ISO", Config{Endpoint: EndpointRuntime, APIKey: "token", AWSRegion: "us-iso-east-1"}, EndpointRuntime, "https://bedrock-runtime.us-iso-east-1.c2s.ic.gov/openai/v1/"},
		{"Runtime ISOB", Config{Endpoint: EndpointRuntime, APIKey: "token", AWSRegion: "us-isob-east-1"}, EndpointRuntime, "https://bedrock-runtime.us-isob-east-1.sc2s.sgov.gov/openai/v1/"},
		{"Runtime ISOE", Config{Endpoint: EndpointRuntime, APIKey: "token", AWSRegion: "eu-isoe-west-1"}, EndpointRuntime, "https://bedrock-runtime.eu-isoe-west-1.cloud.adc-e.uk/openai/v1/"},
		{"Runtime ISOF", Config{Endpoint: EndpointRuntime, APIKey: "token", AWSRegion: "us-isof-south-1"}, EndpointRuntime, "https://bedrock-runtime.us-isof-south-1.csp.hci.ic.gov/openai/v1/"},
		{"Runtime inferred", Config{APIKey: "token", BaseURL: "https://bedrock-runtime.us-west-2.amazonaws.com/openai/v1"}, EndpointRuntime, "https://bedrock-runtime.us-west-2.amazonaws.com/openai/v1/"},
		{"Runtime FIPS inferred", Config{APIKey: "token", BaseURL: "https://bedrock-runtime-fips.us-west-2.amazonaws.com/openai/v1"}, EndpointRuntime, "https://bedrock-runtime-fips.us-west-2.amazonaws.com/openai/v1/"},
		{"Runtime dual stack inferred", Config{APIKey: "token", BaseURL: "https://bedrock-runtime.us-west-2.api.aws/openai/v1"}, EndpointRuntime, "https://bedrock-runtime.us-west-2.api.aws/openai/v1/"},
		{"Runtime China dual stack inferred", Config{APIKey: "token", BaseURL: "https://bedrock-runtime.cn-north-1.api.amazonwebservices.com.cn/openai/v1"}, EndpointRuntime, "https://bedrock-runtime.cn-north-1.api.amazonwebservices.com.cn/openai/v1/"},
		{"Runtime trailing dot inferred", Config{APIKey: "token", BaseURL: "https://bedrock-runtime.us-west-2.amazonaws.com./openai/v1"}, EndpointRuntime, "https://bedrock-runtime.us-west-2.amazonaws.com./openai/v1/"},
		{"Runtime alternate route", Config{Endpoint: EndpointRuntime, APIKey: "token", BaseURL: "https://bedrock-runtime.us-west-2.amazonaws.com/v1"}, EndpointRuntime, "https://bedrock-runtime.us-west-2.amazonaws.com/v1/"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resolved, err := resolveConfig(context.Background(), test.config, time.Now)
			if err != nil {
				t.Fatal(err)
			}
			if resolved.endpoint != test.endpoint || resolved.baseURL.String() != test.baseURL {
				t.Fatalf("endpoint = %q, URL = %q; want %q, %q", resolved.endpoint, resolved.baseURL, test.endpoint, test.baseURL)
			}
		})
	}
}

func TestBedrockEndpointValidation(t *testing.T) {
	clearAWSEnvironment(t)
	tests := []struct {
		name    string
		config  Config
		message string
	}{
		{"unknown endpoint", Config{Endpoint: "unknown", APIKey: "token", AWSRegion: "us-east-1"}, "EndpointMantle` or `EndpointRuntime"},
		{"Runtime HTTP", Config{APIKey: "token", BaseURL: "http://bedrock-runtime.us-east-1.amazonaws.com/openai/v1"}, "require HTTPS"},
		{"Mantle HTTP", Config{APIKey: "token", BaseURL: "http://bedrock-mantle.us-east-1.api.aws/openai/v1"}, "require HTTPS"},
		{"Runtime hostname mismatch", Config{Endpoint: EndpointMantle, APIKey: "token", BaseURL: "https://bedrock-runtime.us-east-1.amazonaws.com/openai/v1"}, "does not match"},
		{"Mantle hostname mismatch", Config{Endpoint: EndpointRuntime, APIKey: "token", BaseURL: "https://bedrock-mantle.us-east-1.api.aws/openai/v1"}, "does not match"},
		{"Runtime region mismatch", Config{APIKey: "token", AWSRegion: "us-west-2", BaseURL: "https://bedrock-runtime.us-east-1.amazonaws.com/openai/v1"}, "does not match"},
		{"Runtime malformed region", Config{APIKey: "token", BaseURL: "https://bedrock-runtime.us--east-1.amazonaws.com/openai/v1"}, "invalid AWS region"},
		{"custom SigV4 endpoint", Config{AWSRegion: "us-east-1", AWSAccessKeyID: "access", AWSSecretAccessKey: "secret", BaseURL: "https://proxy.example/openai/v1"}, "requires an explicit `Endpoint`"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewClient(context.Background(), test.config)
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("error = %v, want substring %q", err, test.message)
			}
		})
	}
}

func TestRuntimeChatAndResponsesAuthentication(t *testing.T) {
	clearAWSEnvironment(t)
	models := []string{"us.openai.gpt-5.6-sol", "us.openai.gpt-5.6-terra", "us.openai.gpt-5.6-luna"}
	for _, authentication := range []string{"bearer", "SigV4"} {
		for _, api := range []string{"chat", "responses"} {
			for _, streaming := range []bool{false, true} {
				for _, model := range models {
					t.Run(fmt.Sprintf("%s/%s/stream=%t/%s", authentication, api, streaming, model), func(t *testing.T) {
						config := Config{Endpoint: EndpointRuntime, AWSRegion: "us-east-1"}
						if authentication == "bearer" {
							config.APIKey = "runtime-token"
						} else {
							config.AWSAccessKeyID = "runtime-access-key"
							config.AWSSecretAccessKey = "runtime-secret-key"
							config.AWSSessionToken = "runtime-session-token"
						}

						client, err := NewClient(context.Background(), config,
							option.WithMaxRetries(0),
							option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
								assertRuntimeRequest(t, req, authentication, api, streaming, model)
								return runtimeResponse(req, api, streaming, model), nil
							})}),
						)
						if err != nil {
							t.Fatal(err)
						}

						if api == "chat" {
							params := openai.ChatCompletionNewParams{
								Model:    openai.ChatModel(model),
								Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage("hello")},
							}
							if streaming {
								stream := client.Chat.Completions.NewStreaming(context.Background(), params)
								var content string
								for stream.Next() {
									if choices := stream.Current().Choices; len(choices) != 0 {
										content += choices[0].Delta.Content
									}
								}
								if err := stream.Err(); err != nil || content != "hello" {
									t.Fatalf("stream content = %q, error = %v", content, err)
								}
							} else {
								completion, err := client.Chat.Completions.New(context.Background(), params)
								if err != nil || len(completion.Choices) != 1 || completion.Choices[0].Message.Content != "hello" {
									t.Fatalf("completion = %#v, error = %v", completion, err)
								}
							}
							return
						}

						params := responses.ResponseNewParams{
							Model: openai.ChatModel(model),
							Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("hello")},
						}
						if streaming {
							stream := client.Responses.NewStreaming(context.Background(), params)
							var content string
							var completed bool
							for stream.Next() {
								event := stream.Current()
								content += event.Delta
								completed = completed || event.Type == "response.completed"
							}
							if err := stream.Err(); err != nil || content != "hello" || !completed {
								t.Fatalf("stream content = %q, completed = %t, error = %v", content, completed, err)
							}
						} else {
							response, err := client.Responses.New(context.Background(), params)
							if err != nil || response.OutputText() != "hello" {
								t.Fatalf("response = %#v, error = %v", response, err)
							}
						}
					})
				}
			}
		}
	}
}

func assertRuntimeRequest(t *testing.T, req *http.Request, authentication, api string, streaming bool, model string) {
	t.Helper()
	path := "/openai/v1/" + api
	if api == "chat" {
		path = "/openai/v1/chat/completions"
	}
	if req.URL.Host != "bedrock-runtime.us-east-1.amazonaws.com" || req.URL.Path != path {
		t.Fatalf("request URL = %q, want Runtime host and path %q", req.URL, path)
	}
	authorization := req.Header.Get("Authorization")
	if authentication == "bearer" {
		if authorization != "Bearer runtime-token" {
			t.Fatalf("authorization = %q", authorization)
		}
	} else if !strings.Contains(authorization, "Credential=runtime-access-key/") ||
		!strings.Contains(authorization, "/us-east-1/bedrock/aws4_request") ||
		req.Header.Get("X-Amz-Security-Token") != "runtime-session-token" {
		t.Fatalf("unexpected SigV4 credentials or session token: authorization = %q", authorization)
	}
	var body struct {
		Model  string `json:"model"`
		Stream bool   `json:"stream"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Model != model || body.Stream != streaming {
		t.Fatalf("body model = %q, stream = %t; want %q, %t", body.Model, body.Stream, model, streaming)
	}
}

func runtimeResponse(req *http.Request, api string, streaming bool, model string) *http.Response {
	var body string
	contentType := "application/json"
	if api == "chat" {
		if streaming {
			contentType = "text/event-stream"
			body = fmt.Sprintf("data: {\"id\":\"chat_runtime\",\"object\":\"chat.completion.chunk\",\"model\":%q,\"choices\":[{\"index\":0,\"delta\":{\"content\":\"hello\"}}]}\n\ndata: [DONE]\n\n", model)
		} else {
			body = fmt.Sprintf(`{"id":"chat_runtime","object":"chat.completion","model":%q,"choices":[{"index":0,"message":{"role":"assistant","content":"hello"},"finish_reason":"stop"}]}`, model)
		}
	} else {
		response := fmt.Sprintf(`{"id":"resp_runtime","object":"response","model":%q,"status":"completed","output":[{"id":"msg_runtime","type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"hello","annotations":[]}]}]}`, model)
		if streaming {
			contentType = "text/event-stream"
			body = "event: response.output_text.delta\ndata: {\"type\":\"response.output_text.delta\",\"delta\":\"hello\",\"item_id\":\"msg_runtime\",\"output_index\":0,\"content_index\":0,\"sequence_number\":1}\n\n" +
				"event: response.completed\ndata: {\"type\":\"response.completed\",\"sequence_number\":2,\"response\":" + response + "}\n\n"
		} else {
			body = response
		}
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": {contentType}, "X-Request-Id": {"req_runtime"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
}

func TestRuntimeCustomHostUsesSelectedSigner(t *testing.T) {
	clearAWSEnvironment(t)
	for _, endpoint := range []Endpoint{EndpointMantle, EndpointRuntime} {
		t.Run(string(endpoint), func(t *testing.T) {
			var authorization string
			client, err := NewClient(context.Background(), Config{
				Endpoint:           endpoint,
				AWSRegion:          "us-east-1",
				AWSAccessKeyID:     "access",
				AWSSecretAccessKey: "secret",
				BaseURL:            "https://proxy.example/openai/v1",
			}, option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				authorization = req.Header.Get("Authorization")
				return successfulResponse(req), nil
			})}))
			if err != nil {
				t.Fatal(err)
			}
			var response *http.Response
			if err := client.Get(context.Background(), "/models", nil, &response); err != nil {
				t.Fatal(err)
			}
			service := bedrockService
			if endpoint == EndpointRuntime {
				service = bedrockRuntimeService
			}
			if !strings.Contains(authorization, "/us-east-1/"+service+"/aws4_request") {
				t.Fatalf("authorization = %q, want signing service %q", authorization, service)
			}
		})
	}
}

func TestRuntimeCanonicalHostInfersSigningService(t *testing.T) {
	clearAWSEnvironment(t)
	var authorization string
	client, err := NewClient(context.Background(), Config{
		AWSAccessKeyID:     "access",
		AWSSecretAccessKey: "secret",
		BaseURL:            "https://bedrock-runtime-fips.us-east-1.api.aws/openai/v1",
	}, option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		authorization = req.Header.Get("Authorization")
		return successfulResponse(req), nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	var response *http.Response
	if err := client.Get(context.Background(), "/models", nil, &response); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(authorization, "/us-east-1/bedrock/aws4_request") {
		t.Fatalf("authorization = %q, want inferred Runtime signer", authorization)
	}
}

func TestRuntimeRetriesRefreshAuthentication(t *testing.T) {
	clearAWSEnvironment(t)
	for _, authentication := range []string{"bearer", "SigV4"} {
		t.Run(authentication, func(t *testing.T) {
			config := Config{Endpoint: EndpointRuntime, AWSRegion: "us-east-1"}
			var providerCalls int
			if authentication == "bearer" {
				config.BedrockTokenProvider = func(context.Context) (string, error) {
					providerCalls++
					return fmt.Sprintf("runtime-token-%d", providerCalls), nil
				}
			} else {
				config.AWSAccessKeyID = "access"
				config.AWSSecretAccessKey = "secret"
			}

			var authorizations []string
			var clockCalls int
			clock := func() time.Time {
				clockCalls++
				return time.Date(2025, 1, 2, 3, clockCalls, 5, 0, time.UTC)
			}
			opts, err := newClientOptions(context.Background(), config, clock,
				option.WithMaxRetries(1),
				option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					authorizations = append(authorizations, req.Header.Get("Authorization"))
					if len(authorizations) == 1 {
						return &http.Response{
							StatusCode: http.StatusTooManyRequests,
							Header:     http.Header{"Content-Type": {"application/json"}, "Retry-After": {"0"}},
							Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"retry"}}`)),
							Request:    req,
						}, nil
					}
					return runtimeResponse(req, "chat", false, "us.openai.gpt-5.6-sol"), nil
				})}),
			)
			if err != nil {
				t.Fatal(err)
			}
			client := openai.NewClient(opts...)
			completion, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
				Model:    openai.ChatModel("us.openai.gpt-5.6-sol"),
				Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage("hello")},
			})
			if err != nil || completion.Choices[0].Message.Content != "hello" {
				t.Fatalf("completion = %#v, error = %v", completion, err)
			}
			if len(authorizations) != 2 || authorizations[0] == authorizations[1] {
				t.Fatalf("authorization attempts = %v, want two independently refreshed values", authorizations)
			}
			if authentication == "bearer" && providerCalls != 2 {
				t.Fatalf("bearer provider calls = %d, want 2", providerCalls)
			}
			if authentication == "SigV4" && !strings.Contains(authorizations[1], "/us-east-1/bedrock/aws4_request") {
				t.Fatalf("retry authorization = %q, want Runtime signer", authorizations[1])
			}
		})
	}
}
