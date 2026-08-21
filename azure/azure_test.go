package azure

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

var invalidJSONRouteBodies = []struct {
	name      string
	body      []byte
	wantError string
}{
	{name: "nil body", wantError: "requires a JSON request body"},
	{name: "empty body", body: []byte{}, wantError: "could not parse JSON request body"},
	{name: "null", body: []byte("null"), wantError: "requires a non-empty model field"},
	{name: "empty object", body: []byte("{}"), wantError: "requires a non-empty model field"},
	{name: "missing model", body: []byte(`{"input":"hello"}`), wantError: "requires a non-empty model field"},
	{name: "empty model", body: []byte(`{"model":""}`), wantError: "requires a non-empty model field"},
	{name: "scalar", body: []byte(`"gpt-4"`), wantError: "could not parse JSON request body"},
	{name: "wrong model type", body: []byte(`{"model":123}`), wantError: "could not parse JSON request body"},
	{name: "malformed", body: []byte(`{"model":`), wantError: "could not parse JSON request body"},
}

func TestJSONRoute(t *testing.T) {
	chatCompletionParams := openai.ChatCompletionNewParams{
		Model: openai.ChatModel("arbitraryDeployment"),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.AssistantMessage("You are a helpful assistant"),
			openai.UserMessage("Can you tell me another word for the universe?"),
		},
	}

	serializedBytes, err := apijson.MarshalRoot(chatCompletionParams)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/chat/completions", bytes.NewReader(serializedBytes))

	if err != nil {
		t.Fatal(err)
	}

	replacementPath, err := getReplacementPathWithDeployment(req)

	if err != nil {
		t.Fatal(err)
	}

	if replacementPath != "/openai/deployments/arbitraryDeployment/chat/completions" {
		t.Fatalf("replacementpath didn't match: %s", replacementPath)
	}

	restoredBody, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(restoredBody, serializedBytes) {
		t.Fatalf("restored body = %q, want %q", restoredBody, serializedBytes)
	}
}

func TestJSONRouteRejectsInvalidBodies(t *testing.T) {
	for _, tc := range invalidJSONRouteBodies {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/chat/completions", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tc.body != nil {
				req.Body = io.NopCloser(bytes.NewReader(tc.body))
			}

			if _, err = getJSONRoute(req); err == nil {
				t.Fatal("expected an error")
			} else if !strings.Contains(err.Error(), tc.wantError) {
				t.Fatalf("error = %q, want it to contain %q", err, tc.wantError)
			}

			if tc.body == nil {
				return
			}
			restoredBody, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(restoredBody, tc.body) {
				t.Fatalf("restored body = %q, want %q", restoredBody, tc.body)
			}
		})
	}
}

func TestAzureJSONRoutesRejectInvalidBodiesBeforeTransport(t *testing.T) {
	routes := []string{
		"/completions",
		"/chat/completions",
		"/embeddings",
		"/audio/speech",
		"/images/generations",
	}

	transportCalls := 0
	routingAttempts := 0
	client := openai.NewClient(
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			routingAttempts++
			return next(req)
		}),
		WithEndpoint("https://my-resource.openai.azure.com", "2024-10-21"),
		WithAPIKey("azure-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				transportCalls++
				return nil, errors.New("unexpected transport call")
			}),
		}),
	)

	wantRoutingAttempts := 0
	for _, route := range routes {
		for _, tc := range invalidJSONRouteBodies {
			t.Run(route+"/"+tc.name, func(t *testing.T) {
				wantRoutingAttempts++
				var body any
				if tc.body != nil {
					body = tc.body
				}

				err := client.Execute(context.Background(), http.MethodPost, route, body, nil)
				if err == nil {
					t.Fatal("expected an error")
				}
				if !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("error = %q, want it to contain %q", err, tc.wantError)
				}
				if transportCalls != 0 {
					t.Fatalf("transport called %d times, want 0", transportCalls)
				}
				if routingAttempts != wantRoutingAttempts {
					t.Fatalf("routing attempted %d times, want %d", routingAttempts, wantRoutingAttempts)
				}
			})
		}
	}
}

func TestAzureJSONRoutesPreserveValidBodyAndRewritePath(t *testing.T) {
	routes := []string{
		"/completions",
		"/chat/completions",
		"/embeddings",
		"/audio/speech",
		"/images/generations",
	}
	body := []byte("{\n  \"model\" : \"deployment-name\",\n  \"input\" : \"unchanged\"\n}\n")

	for _, route := range routes {
		t.Run(route, func(t *testing.T) {
			transportCalls := 0
			client := openai.NewClient(
				WithEndpoint("https://my-resource.openai.azure.com", "2024-10-21"),
				WithAPIKey("azure-api-key"),
				option.WithMaxRetries(0),
				option.WithHTTPClient(&http.Client{
					Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
						transportCalls++
						if got, want := req.URL.EscapedPath(), "/openai/deployments/deployment-name"+route; got != want {
							t.Errorf("path = %q, want %q", got, want)
						}
						gotBody, err := io.ReadAll(req.Body)
						if err != nil {
							t.Fatal(err)
						}
						if !bytes.Equal(gotBody, body) {
							t.Errorf("body = %q, want %q", gotBody, body)
						}
						return &http.Response{
							StatusCode: http.StatusOK,
							Header:     http.Header{"Content-Type": []string{"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
							Request:    req,
						}, nil
					}),
				}),
			)

			var response map[string]any
			if err := client.Execute(context.Background(), http.MethodPost, route, body, &response); err != nil {
				t.Fatal(err)
			}
			if transportCalls != 1 {
				t.Fatalf("transport called %d times, want 1", transportCalls)
			}
		})
	}
}

func TestGetAudioMultipartRoute(t *testing.T) {
	buff := &bytes.Buffer{}
	mw := multipart.NewWriter(buff)

	fw, err := mw.CreateFormFile("file", "test.mp3")

	if err != nil {
		t.Fatal(err)
	}

	if _, err = fw.Write([]byte("ignore me")); err != nil {
		t.Fatal(err)
	}

	if writeErr := mw.WriteField("model", "arbitraryDeployment"); writeErr != nil {
		t.Fatal(writeErr)
	}

	if closeErr := mw.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}

	req, err := http.NewRequest("POST", "/audio/transcriptions", bytes.NewReader(buff.Bytes()))

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", mw.FormDataContentType())

	replacementPath, err := getReplacementPathWithDeployment(req)

	if err != nil {
		t.Fatal(err)
	}

	if replacementPath != "/openai/deployments/arbitraryDeployment/audio/transcriptions" {
		t.Fatalf("replacementpath didn't match: %s", replacementPath)
	}
}

func TestAPIKeyAuthentication(t *testing.T) {
	rc := &requestconfig.RequestConfig{
		Request: &http.Request{
			Header: make(http.Header),
			URL:    &url.URL{},
		},
	}

	if err := WithAPIKey("my-api-key").Apply(rc); err != nil {
		t.Fatal(err)
	}

	if got := rc.Request.Header.Get("Api-Key"); got != "my-api-key" {
		t.Errorf("Api-Key header: got %q, expected %q", got, "my-api-key")
	}
}

func TestAPIKeyAuthenticationSuppressesAutomaticAuthorization(t *testing.T) {
	tests := []struct {
		name        string
		apiKey      string
		adminAPIKey string
	}{
		{name: "OpenAI API key", apiKey: "normal-openai-key"},
		{name: "OpenAI admin API key", adminAPIKey: "normal-admin-key"},
		{name: "both OpenAI keys", apiKey: "normal-openai-key", adminAPIKey: "normal-admin-key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("OPENAI_API_KEY", tt.apiKey)
			t.Setenv("OPENAI_ADMIN_KEY", tt.adminAPIKey)

			var captured *http.Request
			client := openai.NewClient(
				WithEndpoint("https://my-resource.openai.azure.com", "2024-10-21"),
				WithAPIKey("azure-api-key"),
				option.WithHTTPClient(&http.Client{
					Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
						captured = req
						return &http.Response{
							StatusCode: http.StatusOK,
							Header: http.Header{
								"Content-Type": []string{"application/json"},
							},
							Body:    io.NopCloser(strings.NewReader(`{"ok":true}`)),
							Request: req,
						}, nil
					}),
				}),
			)

			var res map[string]any
			if err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res); err != nil {
				t.Fatalf("request failed: %s", err)
			}
			if captured == nil {
				t.Fatal("request was not captured")
			}
			if got := captured.Header.Get("Api-Key"); got != "azure-api-key" {
				t.Fatalf("Api-Key header = %q, want %q", got, "azure-api-key")
			}
			if got := captured.Header.Get("Authorization"); got != "" {
				t.Fatalf("Authorization header = %q, want empty", got)
			}
		})
	}
}

func TestAzureCredentialTransportSecurity(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")

	tests := map[string]struct {
		endpoint        string
		unsafe          bool
		requestBaseURL  string
		rewriteToRemote bool
		wantRequest     bool
	}{
		"HTTPS": {
			endpoint:    "https://my-resource.openai.azure.com",
			wantRequest: true,
		},
		"remote HTTP": {
			endpoint: "http://my-resource.openai.azure.com",
		},
		"mixed-case remote HTTP": {
			endpoint: "HtTp://my-resource.openai.azure.com",
		},
		"non-HTTP scheme": {
			endpoint: "ftp://my-resource.openai.azure.com",
		},
		"request-level remote HTTP override": {
			endpoint:       "https://my-resource.openai.azure.com",
			requestBaseURL: "http://request.example",
		},
		"middleware remote HTTP rewrite": {
			endpoint:        "https://my-resource.openai.azure.com",
			rewriteToRemote: true,
		},
		"unsafe localhost HTTP": {
			endpoint:    "http://localhost:8080",
			unsafe:      true,
			wantRequest: true,
		},
		"unsafe IPv4 loopback HTTP": {
			endpoint:    "http://127.0.0.2:8080",
			unsafe:      true,
			wantRequest: true,
		},
		"unsafe IPv6 loopback HTTP": {
			endpoint:    "http://[::1]:8080",
			unsafe:      true,
			wantRequest: true,
		},
		"unsafe remote HTTP": {
			endpoint: "http://my-resource.openai.azure.com",
			unsafe:   true,
		},
		"unsafe loopback rewritten to remote HTTP": {
			endpoint:        "http://127.0.0.1:8080",
			unsafe:          true,
			rewriteToRemote: true,
		},
	}

	authOptions := map[string]struct {
		option      func() option.RequestOption
		header      string
		headerValue string
	}{
		"API key": {
			option:      func() option.RequestOption { return WithAPIKey("azure-api-key") },
			header:      "Api-Key",
			headerValue: "azure-api-key",
		},
		"token": {
			option:      func() option.RequestOption { return WithTokenCredential(&fake.TokenCredential{}) },
			header:      "Authorization",
			headerValue: "Bearer fake_token",
		},
	}

	for authName, auth := range authOptions {
		for testName, test := range tests {
			t.Run(authName+"/"+testName, func(t *testing.T) {
				var captured *http.Request
				opts := []option.RequestOption{
					WithEndpoint(test.endpoint, "2024-10-21"),
					auth.option(),
				}
				if test.unsafe {
					opts = append(opts, WithUnsafeAllowHTTP())
				}
				if test.rewriteToRemote {
					opts = append(opts, option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
						req.URL.Scheme = "http"
						req.URL.Host = "middleware.example"
						return next(req)
					}))
				}
				opts = append(opts,
					option.WithMaxRetries(0),
					option.WithHTTPClient(&http.Client{
						Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
							captured = req
							return &http.Response{
								StatusCode: http.StatusOK,
								Header:     http.Header{"Content-Type": []string{"application/json"}},
								Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
								Request:    req,
							}, nil
						}),
					}),
				)
				client := openai.NewClient(opts...)

				var requestOptions []option.RequestOption
				if test.requestBaseURL != "" {
					requestOptions = append(requestOptions, option.WithBaseURL(test.requestBaseURL))
				}
				var res map[string]any
				err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res, requestOptions...)
				if test.wantRequest {
					if err != nil {
						t.Fatalf("request failed: %v", err)
					}
					if captured == nil {
						t.Fatal("request did not reach the transport")
					}
					if got := captured.Header.Get(auth.header); got != auth.headerValue {
						t.Fatalf("%s header = %q, want %q", auth.header, got, auth.headerValue)
					}
					return
				}

				if err == nil || !strings.Contains(err.Error(), "azure: authenticated requests require HTTPS") {
					t.Fatalf("expected HTTPS requirement error, got %v", err)
				}
				if captured != nil {
					t.Fatal("insecure credential transport reached the network")
				}
			})
		}
	}
}

func TestAzureCredentialTransportSecurityRedirects(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")

	authOptions := map[string]struct {
		option      func() option.RequestOption
		header      string
		headerValue string
	}{
		"API key": {
			option:      func() option.RequestOption { return WithAPIKey("azure-api-key") },
			header:      "Api-Key",
			headerValue: "azure-api-key",
		},
		"token": {
			option:      func() option.RequestOption { return WithTokenCredential(&fake.TokenCredential{}) },
			header:      "Authorization",
			headerValue: "Bearer fake_token",
		},
	}

	for authName, auth := range authOptions {
		t.Run(authName+"/rejects opaque custom doer", func(t *testing.T) {
			redirectedCredential := ""
			insecureTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				redirectedCredential = req.Header.Get(auth.header)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			t.Cleanup(insecureTarget.Close)

			sourceRequests := 0
			secureSource := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				sourceRequests++
				http.Redirect(w, req, insecureTarget.URL+"/final", http.StatusTemporaryRedirect)
			}))
			t.Cleanup(secureSource.Close)

			client := openai.NewClient(
				WithEndpoint(secureSource.URL, "2024-10-21"),
				auth.option(),
				option.WithMaxRetries(0),
				option.WithHTTPClient(delegatingHTTPDoer{client: secureSource.Client()}),
			)

			var res map[string]any
			err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res)
			if err == nil || !strings.Contains(err.Error(), "custom HTTP clients") {
				t.Errorf("expected custom HTTP client error, got %v", err)
			}
			if sourceRequests != 0 {
				t.Errorf("custom HTTP client reached redirect source %d times", sourceRequests)
			}
			if redirectedCredential != "" {
				t.Errorf("credential reached insecure redirect target: %q", redirectedCredential)
			}
		})

		t.Run(authName+"/rejects HTTPS downgrade", func(t *testing.T) {
			redirectedCredential := ""
			insecureTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				redirectedCredential = req.Header.Get(auth.header)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			t.Cleanup(insecureTarget.Close)

			sourceRequests := 0
			secureSource := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				sourceRequests++
				if req.URL.Path != "/hop" {
					http.Redirect(w, req, "/hop", http.StatusTemporaryRedirect)
					return
				}
				http.Redirect(w, req, insecureTarget.URL+"/final", http.StatusTemporaryRedirect)
			}))
			t.Cleanup(secureSource.Close)

			client := openai.NewClient(
				WithEndpoint(secureSource.URL, "2024-10-21"),
				auth.option(),
				option.WithMaxRetries(0),
				option.WithHTTPClient(secureSource.Client()),
			)

			var res map[string]any
			err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res)
			if err == nil || !strings.Contains(err.Error(), "azure: authenticated requests require HTTPS") {
				t.Fatalf("expected HTTPS requirement error, got %v", err)
			}
			if redirectedCredential != "" {
				t.Fatalf("credential reached insecure redirect target: %q", redirectedCredential)
			}
			if sourceRequests != 2 {
				t.Fatalf("secure source requests = %d, want 2", sourceRequests)
			}
		})

		crossOriginTargets := map[string]func(string) string{
			"different port": func(targetURL string) string { return targetURL },
			"different host": func(targetURL string) string {
				return strings.Replace(targetURL, "127.0.0.1", "localhost", 1)
			},
		}
		for targetName, redirectTarget := range crossOriginTargets {
			t.Run(authName+"/rejects HTTPS cross-origin redirect/"+targetName, func(t *testing.T) {
				redirectedCredential := ""
				target := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					redirectedCredential = req.Header.Get(auth.header)
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"ok":true}`))
				}))
				t.Cleanup(target.Close)

				sourceRequests := 0
				source := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					sourceRequests++
					http.Redirect(w, req, redirectTarget(target.URL)+"/final", http.StatusTemporaryRedirect)
				}))
				t.Cleanup(source.Close)

				client := openai.NewClient(
					WithEndpoint(source.URL, "2024-10-21"),
					auth.option(),
					option.WithMaxRetries(0),
					option.WithHTTPClient(source.Client()),
				)

				var res map[string]any
				err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res)
				if err == nil || !strings.Contains(err.Error(), "authenticated redirects must remain on the original origin") {
					t.Fatalf("expected origin restriction error, got %v", err)
				}
				if sourceRequests != 1 {
					t.Fatalf("secure source requests = %d, want 1", sourceRequests)
				}
				if redirectedCredential != "" {
					t.Fatalf("credential reached cross-origin redirect target: %q", redirectedCredential)
				}
			})
		}

		t.Run(authName+"/unsafe remote HTTPS cannot redirect to loopback HTTP", func(t *testing.T) {
			targetReached := false
			target := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				targetReached = true
			}))
			t.Cleanup(target.Close)

			source := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				http.Redirect(w, req, target.URL+"/final", http.StatusTemporaryRedirect)
			}))
			t.Cleanup(source.Close)
			sourceURL, err := url.Parse(source.URL)
			if err != nil {
				t.Fatal(err)
			}
			transport := source.Client().Transport.(*http.Transport).Clone()
			dialer := &net.Dialer{}
			transport.DialContext = func(ctx context.Context, network string, _ string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, sourceURL.Host)
			}

			client := openai.NewClient(
				WithEndpoint("https://example.com:"+sourceURL.Port(), "2024-10-21"),
				auth.option(),
				WithUnsafeAllowHTTP(),
				option.WithMaxRetries(0),
				option.WithHTTPClient(&http.Client{Transport: transport}),
			)

			var res map[string]any
			requestErr := client.Execute(context.Background(), http.MethodGet, "models", nil, &res)
			if requestErr == nil || !strings.Contains(requestErr.Error(), "azure: authenticated requests require HTTPS") {
				t.Fatalf("expected HTTPS requirement error, got %v", requestErr)
			}
			if targetReached {
				t.Fatal("credential request reached loopback redirect target")
			}
		})

		t.Run(authName+"/preserves HTTPS redirect", func(t *testing.T) {
			redirectedCredential := ""
			secureServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path != "/final" {
					http.Redirect(w, req, "/final", http.StatusTemporaryRedirect)
					return
				}
				redirectedCredential = req.Header.Get(auth.header)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			t.Cleanup(secureServer.Close)

			client := openai.NewClient(
				WithEndpoint(secureServer.URL, "2024-10-21"),
				auth.option(),
				option.WithMaxRetries(0),
				option.WithHTTPClient(secureServer.Client()),
			)

			var res map[string]any
			if err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res); err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if redirectedCredential != auth.headerValue {
				t.Fatalf("redirected credential = %q, want %q", redirectedCredential, auth.headerValue)
			}
		})

		t.Run(authName+"/preserves unsafe loopback redirect", func(t *testing.T) {
			redirectedCredential := ""
			target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				redirectedCredential = req.Header.Get(auth.header)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			t.Cleanup(target.Close)

			source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				http.Redirect(w, req, target.URL+"/final", http.StatusTemporaryRedirect)
			}))
			t.Cleanup(source.Close)

			client := openai.NewClient(
				WithEndpoint(source.URL, "2024-10-21"),
				auth.option(),
				WithUnsafeAllowHTTP(),
				option.WithMaxRetries(0),
				option.WithHTTPClient(source.Client()),
			)

			var res map[string]any
			if err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res); err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if redirectedCredential != auth.headerValue {
				t.Fatalf("redirected credential = %q, want %q", redirectedCredential, auth.headerValue)
			}
		})

		t.Run(authName+"/preserves caller redirect policy", func(t *testing.T) {
			finalReached := false
			secureServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path != "/final" {
					http.Redirect(w, req, "/final", http.StatusTemporaryRedirect)
					return
				}
				finalReached = true
			}))
			t.Cleanup(secureServer.Close)

			redirectPolicyCalled := false
			httpClient := secureServer.Client()
			httpClient.CheckRedirect = func(*http.Request, []*http.Request) error {
				redirectPolicyCalled = true
				return http.ErrUseLastResponse
			}
			client := openai.NewClient(
				WithEndpoint(secureServer.URL, "2024-10-21"),
				auth.option(),
				option.WithMaxRetries(0),
				option.WithHTTPClient(httpClient),
			)

			var res map[string]any
			if err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res); err == nil {
				t.Fatal("expected redirect response error")
			}
			if !redirectPolicyCalled {
				t.Fatal("caller redirect policy was not invoked")
			}
			if finalReached {
				t.Fatal("redirect target was reached despite caller policy")
			}
		})
	}
}

func TestAzureUnsafeHTTPMatchesProxyBypass(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	t.Setenv("NO_PROXY", "")
	t.Setenv("no_proxy", "")

	type requestObservation struct {
		host          string
		apiKey        string
		authorization string
	}
	proxyRequests := make(chan requestObservation, 1)
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		proxyRequests <- requestObservation{
			host:          req.URL.Hostname(),
			apiKey:        req.Header.Get("Api-Key"),
			authorization: req.Header.Get("Authorization"),
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(proxy.Close)
	t.Setenv("HTTP_PROXY", proxy.URL)
	t.Setenv("http_proxy", proxy.URL)

	proxyURL, err := url.Parse(proxy.URL)
	if err != nil {
		t.Fatal(err)
	}

	hosts := map[string]struct {
		host        string
		wantAllowed bool
	}{
		"localhost":           {host: "localhost", wantAllowed: true},
		"uppercase localhost": {host: "LOCALHOST"},
		"localhost dot":       {host: "localhost."},
		"IPv4 loopback":       {host: "127.0.0.1", wantAllowed: true},
		"IPv6 loopback":       {host: "[::1]", wantAllowed: true},
	}
	authOptions := map[string]struct {
		option      func() option.RequestOption
		header      string
		headerValue string
	}{
		"API key": {
			option:      func() option.RequestOption { return WithAPIKey("azure-api-key") },
			header:      "Api-Key",
			headerValue: "azure-api-key",
		},
		"token": {
			option:      func() option.RequestOption { return WithTokenCredential(&fake.TokenCredential{}) },
			header:      "Authorization",
			headerValue: "Bearer fake_token",
		},
	}

	for authName, auth := range authOptions {
		for hostName, host := range hosts {
			t.Run(authName+"/"+hostName, func(t *testing.T) {
				select {
				case <-proxyRequests:
					t.Fatal("unexpected stale proxy request")
				default:
				}

				originRequests := make(chan requestObservation, 1)
				origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					originRequests <- requestObservation{
						host:          req.URL.Hostname(),
						apiKey:        req.Header.Get("Api-Key"),
						authorization: req.Header.Get("Authorization"),
					}
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"ok":true}`))
				}))
				t.Cleanup(origin.Close)
				originURL, parseErr := url.Parse(origin.URL)
				if parseErr != nil {
					t.Fatal(parseErr)
				}

				transport := http.DefaultTransport.(*http.Transport).Clone()
				transport.Proxy = http.ProxyFromEnvironment
				dialer := &net.Dialer{}
				transport.DialContext = func(ctx context.Context, network string, address string) (net.Conn, error) {
					if address == proxyURL.Host {
						return dialer.DialContext(ctx, network, proxyURL.Host)
					}
					return dialer.DialContext(ctx, network, originURL.Host)
				}

				endpoint := "http://" + host.host + ":" + originURL.Port()
				client := openai.NewClient(
					WithEndpoint(endpoint, "2024-10-21"),
					auth.option(),
					WithUnsafeAllowHTTP(),
					option.WithMaxRetries(0),
					option.WithHTTPClient(&http.Client{Transport: transport}),
				)

				var res map[string]any
				requestErr := client.Execute(context.Background(), http.MethodGet, "models", nil, &res)
				if !host.wantAllowed {
					if requestErr == nil || !strings.Contains(requestErr.Error(), "azure: authenticated requests require HTTPS") {
						t.Errorf("expected HTTPS requirement error, got %v", requestErr)
					}
					select {
					case got := <-proxyRequests:
						t.Errorf("rejected hostname reached proxy: %q", got.host)
					default:
					}
					select {
					case got := <-originRequests:
						t.Errorf("rejected hostname reached origin: %q", got.host)
					default:
					}
					return
				}

				if requestErr != nil {
					t.Fatalf("request failed: %v", requestErr)
				}
				select {
				case got := <-proxyRequests:
					t.Fatalf("loopback request reached proxy: %q", got.host)
				default:
				}
				select {
				case got := <-originRequests:
					if credential := map[string]string{"Api-Key": got.apiKey, "Authorization": got.authorization}[auth.header]; credential != auth.headerValue {
						t.Fatalf("%s header = %q, want %q", auth.header, credential, auth.headerValue)
					}
				default:
					t.Fatal("allowed loopback request did not reach origin")
				}
			})
		}
	}
}

func TestJSONRoutePathConstruction(t *testing.T) {
	cases := []struct {
		path     string
		expected string
	}{
		{"/chat/completions", "/openai/deployments/gpt-4/chat/completions"},
		{"/completions", "/openai/deployments/gpt-4/completions"},
		{"/embeddings", "/openai/deployments/gpt-4/embeddings"},
		{"/audio/speech", "/openai/deployments/gpt-4/audio/speech"},
		{"/images/generations", "/openai/deployments/gpt-4/images/generations"},
		{"/models", "/openai/models"}, // endpoint without a deployment
		{"/files", "/openai/files"},   // endpoint without a deployment
	}
	for _, tc := range cases {
		req, _ := http.NewRequest("POST", tc.path, bytes.NewReader([]byte(`{"model":"gpt-4"}`)))
		got, _ := getReplacementPathWithDeployment(req)
		if got != tc.expected {
			t.Errorf("%s: got %q, expected %q", tc.path, got, tc.expected)
		}
	}
}

func TestModelWithSpecialCharsIsEscaped(t *testing.T) {
	tests := map[string]string{
		"slash":               "my-model/v1",
		"query and fragment":  "my-model?api-version=old#frag",
		"bare dot":            ".",
		"bare dot dot":        "..",
		"dot dot slash":       "../my-model",
		"encoded dot segment": "%2e%2e/my-model",
	}
	wantDeployment := map[string]string{
		"slash":               "my-model%2Fv1",
		"query and fragment":  "my-model%3Fapi-version=old%23frag",
		"bare dot":            "%2E",
		"bare dot dot":        "%2E%2E",
		"dot dot slash":       "..%2Fmy-model",
		"encoded dot segment": "%252e%252e%2Fmy-model",
	}

	for name, model := range tests {
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/chat/completions", bytes.NewReader([]byte(`{"model":`+strconv.Quote(model)+`}`)))
			got, err := getReplacementPathWithDeployment(req)
			if err != nil {
				t.Fatal(err)
			}

			expected := "/openai/deployments/" + wantDeployment[name] + "/chat/completions"
			if got != expected {
				t.Errorf("got %q, expected %q", got, expected)
			}
		})
	}
}

func TestMultipartModelWithSpecialCharsIsEscaped(t *testing.T) {
	tests := map[string]string{
		"slash":               "my-model/v1",
		"query and fragment":  "my-model?api-version=old#frag",
		"bare dot":            ".",
		"bare dot dot":        "..",
		"dot dot slash":       "../my-model",
		"encoded dot segment": "%2e%2e/my-model",
	}
	wantDeployment := map[string]string{
		"slash":               "my-model%2Fv1",
		"query and fragment":  "my-model%3Fapi-version=old%23frag",
		"bare dot":            "%2E",
		"bare dot dot":        "%2E%2E",
		"dot dot slash":       "..%2Fmy-model",
		"encoded dot segment": "%252e%252e%2Fmy-model",
	}

	for name, model := range tests {
		t.Run(name, func(t *testing.T) {
			req := newMultipartRouteRequest(t, "/audio/transcriptions", model)
			got, err := getReplacementPathWithDeployment(req)
			if err != nil {
				t.Fatal(err)
			}

			expected := "/openai/deployments/" + wantDeployment[name] + "/audio/transcriptions"
			if got != expected {
				t.Errorf("got %q, expected %q", got, expected)
			}
		})
	}
}

func TestWithEndpointPreservesEscapedPathParams(t *testing.T) {
	tests := map[string]string{
		"slash traversal":     "../videos/vid_123",
		"query and fragment":  "vs_123/files/file_456?limit=1#frag",
		"encoded dot segment": "%2e%2e/videos/vid_123",
	}
	wantPaths := map[string]string{
		"slash traversal":     "https://my-resource.openai.azure.com/openai/vector_stores/..%2Fvideos%2Fvid_123?api-version=2024-10-21",
		"query and fragment":  "https://my-resource.openai.azure.com/openai/vector_stores/vs_123%2Ffiles%2Ffile_456%3Flimit=1%23frag?api-version=2024-10-21",
		"encoded dot segment": "https://my-resource.openai.azure.com/openai/vector_stores/%252e%252e%2Fvideos%2Fvid_123?api-version=2024-10-21",
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			var captured *http.Request
			client := openai.NewClient(
				WithEndpoint("https://my-resource.openai.azure.com", "2024-10-21"),
				WithAPIKey("sk-test"),
				option.WithHTTPClient(&http.Client{
					Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
						captured = req
						return &http.Response{
							StatusCode: http.StatusOK,
							Header: http.Header{
								"Content-Type": []string{"application/json"},
							},
							Body:    io.NopCloser(strings.NewReader(azureVectorStoreResponse)),
							Request: req,
						}, nil
					}),
				}),
			)

			if _, err := client.VectorStores.Get(context.Background(), input); err != nil {
				t.Fatalf("request failed: %s", err)
			}
			if captured == nil {
				t.Fatal("request was not captured")
			}
			if got, want := captured.URL.String(), wantPaths[name]; got != want {
				t.Fatalf("url = %q, want %q", got, want)
			}
		})
	}
}

func TestWithEndpointBaseURL(t *testing.T) {
	tests := map[string]struct {
		endpoint        string
		apiVersion      string
		expectedBaseURL string
		expectedQuery   string
		shouldFail      bool
	}{
		"Azure endpoint": {
			endpoint:        "https://my-resource.openai.azure.com",
			apiVersion:      "2024-10-21",
			expectedBaseURL: "https://my-resource.openai.azure.com/",
			expectedQuery:   "api-version=2024-10-21",
		},
		"Azure endpoint with trailing slash": {
			endpoint:        "https://my-resource.openai.azure.com/",
			apiVersion:      "2024-10-21",
			expectedBaseURL: "https://my-resource.openai.azure.com/",
			expectedQuery:   "api-version=2024-10-21",
		},
		"Azure endpoint with path": {
			endpoint:        "https://my-resource.openai.azure.com/custom/path",
			apiVersion:      "2023-05-15",
			expectedBaseURL: "https://my-resource.openai.azure.com/custom/path/",
			expectedQuery:   "api-version=2023-05-15",
		},
		"empty apiVersion": {
			endpoint:   "https://my-resource.openai.azure.com",
			apiVersion: "",
			shouldFail: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opt := WithEndpoint(tc.endpoint, tc.apiVersion)

			rc := &requestconfig.RequestConfig{
				Request: &http.Request{
					Header: make(http.Header),
					URL:    &url.URL{},
				},
			}

			err := opt.Apply(rc)

			if tc.shouldFail {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("WithEndpoint returned error: %v", err)
			}

			if rc.BaseURL == nil {
				t.Fatal("BaseURL was not set")
			}
			if rc.BaseURL.String() != tc.expectedBaseURL {
				t.Errorf("BaseURL: got %q, expected %q", rc.BaseURL.String(), tc.expectedBaseURL)
			}

			query := rc.Request.URL.RawQuery
			if query != tc.expectedQuery {
				t.Errorf("Query: got %q, expected %q", query, tc.expectedQuery)
			}
		})
	}
}

func newMultipartRouteRequest(t *testing.T, route string, model string) *http.Request {
	t.Helper()

	buff := &bytes.Buffer{}
	mw := multipart.NewWriter(buff)
	fw, err := mw.CreateFormFile("file", "test.mp3")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = fw.Write([]byte("ignore me")); err != nil {
		t.Fatal(err)
	}
	if writeErr := mw.WriteField("model", model); writeErr != nil {
		t.Fatal(writeErr)
	}
	if closeErr := mw.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}

	req, err := http.NewRequest("POST", route, bytes.NewReader(buff.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return req
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type delegatingHTTPDoer struct {
	client *http.Client
}

func (d delegatingHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	return d.client.Do(req)
}

const azureVectorStoreResponse = `{
	"id": "vs_dummy",
	"object": "vector_store",
	"created_at": 0,
	"file_counts": {
		"cancelled": 0,
		"completed": 0,
		"failed": 0,
		"in_progress": 0,
		"total": 0
	},
	"last_active_at": 0,
	"metadata": {},
	"name": "dummy",
	"status": "completed",
	"usage_bytes": 0
}`
