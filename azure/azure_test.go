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
	"sync/atomic"
	"testing"
	"time"

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
			endpoint:    "http://127.0.0.1:8080",
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
				var captured atomic.Pointer[http.Request]
				endpoint := test.endpoint
				transport := http.RoundTripper(roundTripFunc(func(req *http.Request) (*http.Response, error) {
					captured.Store(req)
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
						Request:    req,
					}, nil
				}))
				if test.unsafe && test.wantRequest {
					endpointURL, err := url.Parse(endpoint)
					if err != nil {
						t.Fatal(err)
					}
					network := "tcp4"
					listenAddress := "127.0.0.1:0"
					if ip := net.ParseIP(endpointURL.Hostname()); ip != nil && ip.To4() == nil {
						network = "tcp6"
						listenAddress = "[::1]:0"
					}
					origin := newLoopbackHTTPServer(t, network, listenAddress, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						captured.Store(req)
						w.Header().Set("Content-Type", "application/json")
						_, _ = w.Write([]byte(`{"ok":true}`))
					}))
					originURL, err := url.Parse(origin.URL)
					if err != nil {
						t.Fatal(err)
					}
					endpointURL.Host = net.JoinHostPort(endpointURL.Hostname(), originURL.Port())
					endpoint = endpointURL.String()
				}
				opts := []option.RequestOption{
					WithEndpoint(endpoint, "2024-10-21"),
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
					option.WithHTTPClient(&http.Client{Transport: transport}),
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
					request := captured.Load()
					if request == nil {
						t.Fatal("request did not reach the transport")
					}
					if got := request.Header.Get(auth.header); got != auth.headerValue {
						t.Fatalf("%s header = %q, want %q", auth.header, got, auth.headerValue)
					}
					return
				}

				if err == nil || !strings.Contains(err.Error(), "azure: authenticated requests require HTTPS") {
					t.Fatalf("expected HTTPS requirement error, got %v", err)
				}
				if captured.Load() != nil {
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
			var redirectedCredential atomicString
			insecureTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				redirectedCredential.Store(req.Header.Get(auth.header))
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			t.Cleanup(insecureTarget.Close)

			var sourceRequests atomic.Int32
			secureSource := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				sourceRequests.Add(1)
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
			if got := sourceRequests.Load(); got != 0 {
				t.Errorf("custom HTTP client reached redirect source %d times", got)
			}
			if got := redirectedCredential.Load(); got != "" {
				t.Errorf("credential reached insecure redirect target: %q", got)
			}
		})

		t.Run(authName+"/rejects HTTPS downgrade", func(t *testing.T) {
			var redirectedCredential atomicString
			insecureTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				redirectedCredential.Store(req.Header.Get(auth.header))
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			t.Cleanup(insecureTarget.Close)

			var sourceRequests atomic.Int32
			secureSource := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				sourceRequests.Add(1)
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
			if got := redirectedCredential.Load(); got != "" {
				t.Fatalf("credential reached insecure redirect target: %q", got)
			}
			if got := sourceRequests.Load(); got != 2 {
				t.Fatalf("secure source requests = %d, want 2", got)
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
				redirectBodyClosed := make(chan struct{}, 1)
				var redirectedCredential atomicString
				target := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					redirectedCredential.Store(req.Header.Get(auth.header))
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"ok":true}`))
				}))
				t.Cleanup(target.Close)

				var sourceRequests atomic.Int32
				source := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					sourceRequests.Add(1)
					http.Redirect(w, req, redirectTarget(target.URL)+"/final", http.StatusTemporaryRedirect)
				}))
				t.Cleanup(source.Close)

				client := openai.NewClient(
					WithEndpoint(source.URL, "2024-10-21"),
					auth.option(),
					option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
						req.GetBody = func() (io.ReadCloser, error) {
							return &closeTrackingBody{
								Reader: strings.NewReader(`{"model":"test"}`),
								closed: redirectBodyClosed,
							}, nil
						}
						return next(req)
					}),
					// Keep the normal retry count so the source-request assertion
					// proves deterministic redirect failures skip outer retries.
					option.WithMaxRetryDelay(time.Millisecond),
					option.WithHTTPClient(source.Client()),
				)

				var res map[string]any
				err := client.Execute(context.Background(), http.MethodPost, "models", []byte(`{"model":"test"}`), &res)
				if err == nil || !strings.Contains(err.Error(), "authenticated redirects must remain on the original origin") {
					t.Fatalf("expected origin restriction error, got %v", err)
				}
				if got := sourceRequests.Load(); got != 1 {
					t.Fatalf("secure source requests = %d, want 1", got)
				}
				if got := redirectedCredential.Load(); got != "" {
					t.Fatalf("credential reached cross-origin redirect target: %q", got)
				}
				// The token policy rebuilds the request body; the direct API-key path
				// exposes the close-tracking replay body to this transport wrapper.
				if authName == "API key" {
					select {
					case <-redirectBodyClosed:
					default:
						t.Fatal("rejected redirect body was not closed")
					}
				}
			})
		}

		redirectTargetMutations := map[string]func(*http.Request){
			"Host override": func(req *http.Request) {
				req.Host = "attacker.invalid"
			},
			"opaque target": func(req *http.Request) {
				req.URL.Opaque = "//attacker.invalid/final"
			},
		}
		for mutationName, mutate := range redirectTargetMutations {
			t.Run(authName+"/rejects redirect "+mutationName, func(t *testing.T) {
				redirectBodyClosed := make(chan struct{}, 1)
				var finalReached atomic.Bool
				server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					if req.URL.Path == "/final" {
						finalReached.Store(true)
						w.Header().Set("Content-Type", "application/json")
						_, _ = w.Write([]byte(`{"ok":true}`))
						return
					}
					http.Redirect(w, req, "/final", http.StatusTemporaryRedirect)
				}))
				t.Cleanup(server.Close)

				httpClient := server.Client()
				baseTransport := httpClient.Transport
				var transportCalls atomic.Int32
				httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
					transportCalls.Add(1)
					return baseTransport.RoundTrip(req)
				})
				httpClient.CheckRedirect = func(req *http.Request, _ []*http.Request) error {
					mutate(req)
					req.Body = &closeTrackingBody{
						Reader: strings.NewReader(`{"model":"test"}`),
						closed: redirectBodyClosed,
					}
					return nil
				}

				client := openai.NewClient(
					WithEndpoint(server.URL, "2024-10-21"),
					auth.option(),
					option.WithMaxRetries(0),
					option.WithHTTPClient(httpClient),
				)

				var res map[string]any
				err := client.Execute(context.Background(), http.MethodPost, "models", []byte(`{"model":"test"}`), &res)
				if err == nil || !strings.Contains(err.Error(), "request URL origin must match the configured base URL") {
					t.Fatalf("expected canonical target error, got %v", err)
				}
				if got := transportCalls.Load(); got != 1 {
					t.Fatalf("underlying transport calls = %d, want 1", got)
				}
				if finalReached.Load() {
					t.Fatal("mutated redirect target reached the server")
				}
				select {
				case <-redirectBodyClosed:
				default:
					t.Fatal("rejected redirect body was not closed")
				}
			})
		}

		t.Run(authName+"/unsafe remote HTTPS cannot redirect to loopback HTTP", func(t *testing.T) {
			var targetReached atomic.Bool
			target := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				targetReached.Store(true)
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
			if targetReached.Load() {
				t.Fatal("credential request reached loopback redirect target")
			}
		})

		t.Run(authName+"/preserves HTTPS redirect", func(t *testing.T) {
			var redirectedCredential atomicString
			secureServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path != "/final" {
					http.Redirect(w, req, "/final", http.StatusTemporaryRedirect)
					return
				}
				redirectedCredential.Store(req.Header.Get(auth.header))
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
			if got := redirectedCredential.Load(); got != auth.headerValue {
				t.Fatalf("redirected credential = %q, want %q", got, auth.headerValue)
			}
		})

		t.Run(authName+"/preserves unsafe loopback redirect", func(t *testing.T) {
			var redirectedCredential atomicString
			target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				redirectedCredential.Store(req.Header.Get(auth.header))
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
			if got := redirectedCredential.Load(); got != auth.headerValue {
				t.Fatalf("redirected credential = %q, want %q", got, auth.headerValue)
			}
		})

		for _, targetScheme := range []string{"HTTP", "HTTPS"} {
			t.Run(authName+"/rejects unsafe cross-host loopback "+targetScheme+" redirect", func(t *testing.T) {
				targetReached := make(chan struct{}, 1)
				target := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					targetReached <- struct{}{}
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"ok":true}`))
				}))
				if targetScheme == "HTTPS" {
					target.StartTLS()
				} else {
					target.Start()
				}
				t.Cleanup(target.Close)

				source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					http.Redirect(w, req, target.URL+"/final", http.StatusTemporaryRedirect)
				}))
				t.Cleanup(source.Close)
				sourceEndpoint := strings.Replace(source.URL, "127.0.0.1", "localhost", 1)

				client := openai.NewClient(
					WithEndpoint(sourceEndpoint, "2024-10-21"),
					auth.option(),
					WithUnsafeAllowHTTP(),
					option.WithMaxRetries(0),
					option.WithHTTPClient(target.Client()),
				)

				var res map[string]any
				err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res)
				if err == nil || !strings.Contains(err.Error(), "authenticated redirects must remain on the original origin") {
					t.Fatalf("expected origin restriction error, got %v", err)
				}
				select {
				case <-targetReached:
					t.Fatal("cross-host loopback redirect reached target")
				default:
				}
			})
		}

		t.Run(authName+"/preserves unsafe loopback HTTPS upgrade", func(t *testing.T) {
			var redirectedCredential atomicString
			target := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				redirectedCredential.Store(req.Header.Get(auth.header))
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			t.Cleanup(target.Close)

			var sourceRequests atomic.Int32
			source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				sourceRequests.Add(1)
				http.Redirect(w, req, target.URL+"/final", http.StatusTemporaryRedirect)
			}))
			t.Cleanup(source.Close)

			transport := target.Client().Transport.(*http.Transport).Clone()
			transport.Proxy = nil
			client := openai.NewClient(
				WithEndpoint(source.URL, "2024-10-21"),
				auth.option(),
				WithUnsafeAllowHTTP(),
				option.WithMaxRetries(0),
				option.WithHTTPClient(&http.Client{Transport: transport}),
			)

			var res map[string]any
			if err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res); err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if got := sourceRequests.Load(); got != 1 {
				t.Fatalf("source requests = %d, want 1", got)
			}
			if got := redirectedCredential.Load(); got != auth.headerValue {
				t.Fatalf("redirected credential = %q, want %q", got, auth.headerValue)
			}
		})

		t.Run(authName+"/preserves caller redirect policy", func(t *testing.T) {
			var finalReached atomic.Bool
			secureServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path != "/final" {
					http.Redirect(w, req, "/final", http.StatusTemporaryRedirect)
					return
				}
				finalReached.Store(true)
			}))
			t.Cleanup(secureServer.Close)

			var redirectPolicyCalled atomic.Bool
			httpClient := secureServer.Client()
			httpClient.CheckRedirect = func(*http.Request, []*http.Request) error {
				redirectPolicyCalled.Store(true)
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
			if !redirectPolicyCalled.Load() {
				t.Fatal("caller redirect policy was not invoked")
			}
			if finalReached.Load() {
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

	hosts := map[string]struct {
		host          string
		network       string
		listenAddress string
		wantAllowed   bool
	}{
		"localhost":                  {host: "localhost", network: "tcp4", listenAddress: "127.0.0.1:0", wantAllowed: true},
		"localhost with IPv6 origin": {host: "localhost", network: "tcp6", listenAddress: "[::1]:0", wantAllowed: true},
		"uppercase localhost":        {host: "LOCALHOST", network: "tcp4", listenAddress: "127.0.0.1:0"},
		"localhost dot":              {host: "localhost.", network: "tcp4", listenAddress: "127.0.0.1:0"},
		"IPv4 loopback":              {host: "127.0.0.1", network: "tcp4", listenAddress: "127.0.0.1:0", wantAllowed: true},
		"IPv6 loopback":              {host: "[::1]", network: "tcp6", listenAddress: "[::1]:0", wantAllowed: true},
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
				origin := newLoopbackHTTPServer(t, host.network, host.listenAddress, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					originRequests <- requestObservation{
						host:          req.URL.Hostname(),
						apiKey:        req.Header.Get("Api-Key"),
						authorization: req.Header.Get("Authorization"),
					}
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"ok":true}`))
				}))
				originURL, parseErr := url.Parse(origin.URL)
				if parseErr != nil {
					t.Fatal(parseErr)
				}

				transport := http.DefaultTransport.(*http.Transport).Clone()
				transport.Proxy = http.ProxyFromEnvironment

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

func TestAzureUnsafeHTTPForcesDirectTransport(t *testing.T) {
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
		t.Run(authName, func(t *testing.T) {
			originCredential := make(chan string, 1)
			origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				originCredential <- req.Header.Get(auth.header)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			t.Cleanup(origin.Close)
			originURL, err := url.Parse(origin.URL)
			if err != nil {
				t.Fatal(err)
			}

			proxyCredential := make(chan string, 1)
			proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				proxyCredential <- req.Header.Get(auth.header)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			t.Cleanup(proxy.Close)
			proxyURL, err := url.Parse(proxy.URL)
			if err != nil {
				t.Fatal(err)
			}

			transport := http.DefaultTransport.(*http.Transport).Clone()
			transport.Proxy = http.ProxyURL(proxyURL)
			dialer := &net.Dialer{}
			transport.DialContext = func(ctx context.Context, network string, _ string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, proxyURL.Host)
			}
			client := openai.NewClient(
				WithEndpoint("http://localhost:"+originURL.Port(), "2024-10-21"),
				auth.option(),
				WithUnsafeAllowHTTP(),
				option.WithMaxRetries(0),
				option.WithHTTPClient(&http.Client{Transport: transport}),
			)

			var res map[string]any
			if err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res); err != nil {
				t.Fatalf("request failed: %v", err)
			}
			select {
			case credential := <-proxyCredential:
				t.Fatalf("loopback credential reached explicit proxy: %q", credential)
			default:
			}
			select {
			case credential := <-originCredential:
				if credential != auth.headerValue {
					t.Fatalf("%s header = %q, want %q", auth.header, credential, auth.headerValue)
				}
			default:
				t.Fatal("direct loopback origin did not receive request")
			}
		})
	}
}

func TestAzureUnsafeHTTPBypassesOpaqueRoundTripper(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")

	var originReached atomic.Bool
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		originReached.Store(true)
		if got := req.Header.Get("Api-Key"); got != "azure-api-key" {
			t.Errorf("Api-Key header = %q, want %q", got, "azure-api-key")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(origin.Close)
	originURL, err := url.Parse(origin.URL)
	if err != nil {
		t.Fatal(err)
	}

	transportCalled := false
	client := openai.NewClient(
		WithEndpoint("http://localhost:"+originURL.Port(), "2024-10-21"),
		WithAPIKey("azure-api-key"),
		WithUnsafeAllowHTTP(),
		option.WithMaxRetries(0),
		option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			transportCalled = true
			return nil, errors.New("unexpected transport call")
		})}),
	)

	var res map[string]any
	if err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res); err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if transportCalled {
		t.Fatal("opaque transport received unsafe loopback request")
	}
	if !originReached.Load() {
		t.Fatal("direct loopback origin did not receive request")
	}
}

func TestAzureUnsafeHTTPPreservesResponseHeaderTimeout(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")

	origin := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		<-req.Context().Done()
	}))
	t.Cleanup(origin.Close)
	originURL, err := url.Parse(origin.URL)
	if err != nil {
		t.Fatal(err)
	}

	const responseHeaderTimeout = 50 * time.Millisecond
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = responseHeaderTimeout
	client := openai.NewClient(
		WithEndpoint("http://localhost:"+originURL.Port(), "2024-10-21"),
		WithAPIKey("azure-api-key"),
		WithUnsafeAllowHTTP(),
		option.WithMaxRetries(0),
		option.WithHTTPClient(&http.Client{Transport: transport}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	started := time.Now()
	var res map[string]any
	err = client.Execute(ctx, http.MethodGet, "models", nil, &res)
	if err == nil || !strings.Contains(err.Error(), "timeout awaiting response headers") {
		t.Fatalf("expected response header timeout, got %v", err)
	}
	if elapsed := time.Since(started); elapsed >= time.Second {
		t.Fatalf("response header timeout took %v, want less than 1s", elapsed)
	}
}

func TestAzureUnsafeHTTPReusesDirectTransport(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")

	var connections atomic.Int32
	origin := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	origin.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			connections.Add(1)
		}
	}
	origin.Start()
	t.Cleanup(origin.Close)
	originURL, err := url.Parse(origin.URL)
	if err != nil {
		t.Fatal(err)
	}

	client := openai.NewClient(
		WithEndpoint("http://localhost:"+originURL.Port(), "2024-10-21"),
		WithAPIKey("azure-api-key"),
		WithUnsafeAllowHTTP(),
		option.WithMaxRetries(0),
	)

	for range 2 {
		var res map[string]any
		if err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res); err != nil {
			t.Fatalf("request failed: %v", err)
		}
	}
	if got := connections.Load(); got != 1 {
		t.Fatalf("loopback connections = %d, want 1 shared connection", got)
	}
}

func TestAzureHTTPSDoesNotCacheDirectTransports(t *testing.T) {
	origin := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(origin.Close)
	originURL, err := url.Parse(origin.URL)
	if err != nil {
		t.Fatal(err)
	}

	directTransports := &azureDirectLoopbackTransportCache{}
	for range 10 {
		base := origin.Client().Transport.(*http.Transport).Clone()
		transport := newAzureCredentialTransport(base, directTransports)
		req, err := http.NewRequest(http.MethodGet, origin.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		ctx := context.WithValue(req.Context(), azureCredentialOriginContextKey{}, azureCredentialOriginFromURL(originURL))
		res, err := transport.RoundTrip(req.WithContext(ctx))
		if err != nil {
			t.Fatalf("HTTPS request failed: %v", err)
		}
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
		base.CloseIdleConnections()
	}

	cached := 0
	directTransports.transports.Range(func(_, _ any) bool {
		cached++
		return true
	})
	if cached != 0 {
		t.Fatalf("cached direct transports after HTTPS requests = %d, want 0", cached)
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

func newLoopbackHTTPServer(t *testing.T, network string, address string, handler http.Handler) *httptest.Server {
	t.Helper()
	listener, err := net.Listen(network, address)
	if err != nil {
		if network == "tcp6" {
			t.Skipf("IPv6 loopback is unavailable: %v", err)
		}
		t.Fatal(err)
	}

	server := httptest.NewUnstartedServer(handler)
	_ = server.Listener.Close()
	server.Listener = listener
	server.Start()
	t.Cleanup(server.Close)
	return server
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

type closeTrackingBody struct {
	io.Reader
	closed chan<- struct{}
}

func (b *closeTrackingBody) Close() error {
	select {
	case b.closed <- struct{}{}:
	default:
	}
	return nil
}

type atomicString struct {
	value atomic.Value
}

func (s *atomicString) Load() string {
	value, _ := s.value.Load().(string)
	return value
}

func (s *atomicString) Store(value string) {
	s.value.Store(value)
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
