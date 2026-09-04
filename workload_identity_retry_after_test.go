package openai_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

func TestWorkloadIdentityUnauthorizedRetryAfter(t *testing.T) {
	for _, x509 := range []bool{false, true} {
		for _, standalone := range []bool{false, true} {
			for _, test := range []struct {
				name              string
				headers           http.Header
				minimum           time.Duration
				maximum           time.Duration
				refused, canceled bool
				cancelOnResponse  bool
				retries           int
			}{
				{name: "fractional seconds default retries", headers: http.Header{"Retry-After": {"0.025"}}, minimum: 25 * time.Millisecond, retries: -1},
				{name: "fractional milliseconds one retry", headers: http.Header{"Retry-After-Ms": {"1.001"}}, minimum: 1001 * time.Microsecond, retries: 1},
				{name: "milliseconds precedence", headers: http.Header{"Retry-After-Ms": {"1.001"}, "Retry-After": {"90"}}, minimum: 1001 * time.Microsecond, retries: 1},
				{name: "invalid milliseconds fall back", headers: http.Header{"Retry-After-Ms": {"invalid"}, "Retry-After": {"0.025"}}, minimum: 25 * time.Millisecond, retries: 1},
				{name: "custom cap", headers: http.Header{"Retry-After-Ms": {"2"}}, maximum: time.Millisecond, refused: true, retries: 1},
				{name: "overcap", headers: http.Header{"Retry-After": {"90"}}, refused: true, retries: -1},
				{name: "overflow", headers: http.Header{"Retry-After-Ms": {"1e999"}}, refused: true, retries: 1},
				{name: "invalid immediate", headers: http.Header{"Retry-After": {"NaN"}}, retries: 1},
				{name: "absent immediate", retries: 1},
				{name: "auth replay opt out compatibility", headers: http.Header{"Retry-After": {"0.025"}, "X-Should-Retry": {"false"}}, minimum: 25 * time.Millisecond, retries: 1},
				{name: "canceled wait", headers: http.Header{"Retry-After": {"1"}}, canceled: true, retries: 1},
				{name: "canceled overcap response", headers: http.Header{"Retry-After": {"90"}}, canceled: true, cancelOnResponse: true, retries: 1},
				{name: "canceled overflow response", headers: http.Header{"Retry-After-Ms": {"1e999"}}, canceled: true, cancelOnResponse: true, retries: 1},
				{name: "zero retry budget", headers: http.Header{"Retry-After": {"0.025"}}, minimum: 25 * time.Millisecond, retries: 0},
			} {
				t.Run(fmt.Sprintf("x509=%t/standalone=%t/%s", x509, standalone, test.name), func(t *testing.T) {
					if standalone && (test.retries == 0 || test.maximum != 0) {
						t.Skip("standalone middleware owns its existing retry count and default cap")
					}
					var issuerCalls atomic.Int32
					var lastIssuer atomic.Int64
					var apiServer *x509ConformanceServer
					issuerResponse := func() *http.Response {
						issuerCalls.Add(1)
						lastIssuer.Store(time.Now().UnixNano())
						return rootWorkloadResponse(http.StatusOK, fmt.Sprintf(`{"access_token":"synthetic-token-%d","expires_in":3600}`, issuerCalls.Load()))
					}
					var opts []option.RequestOption
					var direct func(*http.Request, func(*http.Request) (*http.Response, error)) (*http.Response, error)
					if x509 {
						t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
						config, issuer, api := newX509WorkloadIdentityIntegration(t)
						apiServer = api
						issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							count := issuerCalls.Add(1)
							lastIssuer.Store(time.Now().UnixNano())
							_, _ = io.WriteString(w, x509IntegrationTokenResponse(fmt.Sprintf("synthetic-token-%d", count)))
						})
						opts = append(opts, option.WithX509WorkloadIdentity(config))
						identity, err := auth.NewX509WorkloadIdentityAuth(config)
						if err != nil {
							t.Fatal(err)
						}
						direct = func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
							return auth.X509WorkloadIdentityMiddleware(identity, config.Transport, req, next)
						}
					} else {
						provider := &mockSubjectTokenProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT}
						config := testWorkloadIdentity(provider)
						httpClient := &http.Client{Transport: retryDelayRoundTripperFunc(func(*http.Request) (*http.Response, error) { return issuerResponse(), nil })}
						opts = append(opts, option.WithHTTPClient(httpClient), option.WithWorkloadIdentity(config))
						identity, err := auth.NewWorkloadIdentityAuth(config)
						if err != nil {
							t.Fatal(err)
						}
						direct = func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
							return auth.WorkloadIdentityMiddleware(identity, httpClient, req, next)
						}
					}
					ctx, cancel := context.WithCancel(t.Context())
					defer cancel()
					attempts := 0
					var observations sync.Mutex
					var firstAt, secondAt time.Time
					body := &retryDelayResponseBody{Reader: strings.NewReader(`{"error":{"message":"synthetic rejected bearer"}}`)}
					first := &http.Response{StatusCode: 401, Header: test.headers.Clone(), Body: body}
					next := func(req *http.Request) (*http.Response, error) {
						observations.Lock()
						defer observations.Unlock()
						attempts++
						if attempts == 1 {
							firstAt = time.Now()
							first.Request = req
							if test.cancelOnResponse {
								cancel()
							} else if test.canceled {
								time.AfterFunc(time.Millisecond, cancel)
							}
							return first, nil
						}
						secondAt = time.Now()
						return rootWorkloadResponse(200, `{"data":[]}`), nil
					}
					var response *http.Response
					var err error
					if standalone {
						req, requestErr := http.NewRequestWithContext(ctx, http.MethodGet, "https://mtls.api.openai.com/v1/models", nil)
						if requestErr != nil {
							t.Fatal(requestErr)
						}
						response, err = direct(req, next)
						if response != nil && response.Body != nil {
							defer func() { _ = response.Body.Close() }()
						}
					} else {
						if x509 {
							apiServer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
								result, _ := next(req)
								for name, values := range result.Header {
									w.Header()[name] = values
								}
								w.WriteHeader(result.StatusCode)
								_, _ = io.Copy(w, result.Body)
								_ = result.Body.Close()
							})
						} else {
							opts = append(opts, option.WithMiddleware(func(req *http.Request, _ option.MiddlewareNext) (*http.Response, error) { return next(req) }))
						}
						if test.retries >= 0 {
							opts = append(opts, option.WithMaxRetries(test.retries))
						}
						if test.maximum != 0 {
							opts = append(opts, option.WithMaxRetryDelay(test.maximum))
						}
						client := openai.NewClient(opts...)
						_, err = client.Models.List(ctx, option.WithResponseInto(&response))
						if response != nil && response.Body != nil {
							defer func() { _ = response.Body.Close() }()
						}
					}
					observations.Lock()
					defer observations.Unlock()
					refused := test.refused || (!standalone && test.retries == 0)
					if test.canceled {
						if standalone && body.closes != 1 {
							t.Errorf("canceled response closes=%d, want 1", body.closes)
						}
						if !errors.Is(err, context.Canceled) {
							t.Errorf("canceled wait error=%v", err)
						}
					} else if refused {
						if response == nil || response.StatusCode != 401 {
							t.Fatalf("refused response=%v error=%v", response, err)
						}
						if response.Header.Get("Retry-After") != test.headers.Get("Retry-After") || response.Header.Get("Retry-After-Ms") != test.headers.Get("Retry-After-Ms") {
							t.Error("refusal changed retry headers")
						}
						if standalone && (response.Body != first.Body || body.closes != 0) {
							t.Error("standalone refusal changed or closed original response")
						}
						payload, readErr := io.ReadAll(response.Body)
						if readErr != nil || !strings.Contains(string(payload), "synthetic rejected bearer") {
							t.Errorf("refused body=%q error=%v", payload, readErr)
						}
					} else if err != nil || response == nil || response.StatusCode != 200 {
						t.Fatalf("replay response=%v error=%v", response, err)
					}
					want := 2
					if refused || test.canceled {
						want = 1
					}
					if attempts != want || issuerCalls.Load() != int32(want) {
						t.Errorf("API/issuer calls=%d/%d want=%d", attempts, issuerCalls.Load(), want)
					}
					if want == 2 && secondAt.Sub(firstAt) < test.minimum {
						t.Errorf("API replay wait=%s minimum=%s", secondAt.Sub(firstAt), test.minimum)
					}
					if want == 2 && time.Unix(0, lastIssuer.Load()).Sub(firstAt) < test.minimum {
						t.Errorf("issuer refresh wait=%s minimum=%s", time.Unix(0, lastIssuer.Load()).Sub(firstAt), test.minimum)
					}
				})
			}
		}
	}
}

func TestWorkloadIdentityIssuerRetryAfter(t *testing.T) {
	for _, status := range []int{429, 503, 400, 401, 403} {
		for _, test := range []struct {
			name                                       string
			headers                                    http.Header
			minimum, maximum                           time.Duration
			refused, canceled, permanent, clearHeaders bool
		}{
			{name: "seconds", headers: http.Header{"Retry-After": {"0.55"}}, minimum: 550 * time.Millisecond},
			{name: "milliseconds precedence", headers: http.Header{"Retry-After-Ms": {"550"}, "Retry-After": {"90"}}, minimum: 550 * time.Millisecond},
			{name: "default cap", headers: http.Header{"Retry-After": {"90"}}, refused: true},
			{name: "caller cap", headers: http.Header{"Retry-After-Ms": {"25"}}, maximum: time.Millisecond, refused: true},
			{name: "overflow", headers: http.Header{"Retry-After-Ms": {"1e999"}}, refused: true},
			{name: "body mutates headers", headers: http.Header{"Retry-After": {"90"}}, refused: true, clearHeaders: true},
			{name: "invalid fallback", headers: http.Header{"Retry-After": {"NaN"}}, maximum: time.Millisecond},
			{name: "absent fallback", maximum: time.Millisecond},
			{name: "cancel wait", headers: http.Header{"Retry-After": {"1"}}, canceled: true},
			{name: "permanent OAuth", headers: http.Header{"Retry-After": {"0.55"}}, permanent: true, refused: true},
		} {
			t.Run(fmt.Sprintf("status=%d/%s", status, test.name), func(t *testing.T) {
				if test.permanent && status != 400 && status != 401 && status != 403 {
					t.Skip("only OAuth statuses have permanent OAuth errors")
				}
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				issuerCalls, apiCalls := 0, 0
				var firstAt, secondAt time.Time
				code := "temporarily_unavailable"
				if test.permanent {
					code = "invalid_client"
				}
				first := &http.Response{StatusCode: status, Header: test.headers.Clone()}
				body := &retryDelayResponseBody{Reader: strings.NewReader(fmt.Sprintf(`{"error":%q,"error_description":"synthetic-private-description"}`, code))}
				first.Body = &issuerRetryAfterBody{retryDelayResponseBody: body, headers: first.Header, clear: test.clearHeaders}
				httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
					if req.URL.Host == "auth.openai.com" {
						issuerCalls++
						if issuerCalls == 1 {
							firstAt = time.Now()
							if test.canceled {
								time.AfterFunc(time.Millisecond, cancel)
							}
							return first, nil
						}
						secondAt = time.Now()
						return rootWorkloadResponse(200, `{"access_token":"synthetic-refreshed-bearer","expires_in":3600}`), nil
					}
					apiCalls++
					return rootWorkloadResponse(200, `{"data":[]}`), nil
				}}}
				provider := &mockSubjectTokenProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT}
				opts := []option.RequestOption{option.WithWorkloadIdentity(testWorkloadIdentity(provider)), option.WithHTTPClient(httpClient), option.WithMaxRetries(1)}
				if test.maximum != 0 {
					opts = append(opts, option.WithMaxRetryDelay(test.maximum))
				}
				client := openai.NewClient(opts...)
				_, err := client.Models.List(ctx)
				if test.canceled {
					if !errors.Is(err, context.Canceled) {
						t.Errorf("cancel error=%v", err)
					}
				} else if test.refused {
					if err == nil {
						t.Fatal("issuer refusal unexpectedly succeeded")
					}
					if status == 400 || status == 401 || status == 403 {
						var oauthError *auth.OAuthError
						if !errors.As(err, &oauthError) || oauthError.StatusCode != status || string(oauthError.ErrorCode) != code {
							t.Errorf("issuer OAuth identity lost: %v", err)
						}
					} else if err.Error() != fmt.Sprintf("token exchange failed with status %d", status) {
						t.Errorf("issuer error changed: %v", err)
					}
					if strings.Contains(err.Error(), "synthetic-private-description") {
						t.Error("issuer private description leaked")
					}
				} else if err != nil {
					t.Fatal(err)
				}
				wantIssuer, wantAPI := 2, 1
				if test.refused || test.canceled {
					wantIssuer, wantAPI = 1, 0
				}
				if issuerCalls != wantIssuer || apiCalls != wantAPI {
					t.Errorf("issuer/API attempts=%d/%d want=%d/%d", issuerCalls, apiCalls, wantIssuer, wantAPI)
				}
				if wantIssuer == 2 && secondAt.Sub(firstAt) < test.minimum {
					t.Errorf("issuer wait=%s minimum=%s", secondAt.Sub(firstAt), test.minimum)
				}
				if body.closes != 1 {
					t.Errorf("issuer body closes=%d want1", body.closes)
				}
			})
		}
	}
}

type issuerRetryAfterBody struct {
	*retryDelayResponseBody
	headers http.Header
	clear   bool
}

func (body *issuerRetryAfterBody) Read(p []byte) (int, error) {
	if body.clear {
		clear(body.headers)
	}
	return body.retryDelayResponseBody.Read(p)
}
