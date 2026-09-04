package openai_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

func TestWorkloadIdentityBackgroundIssuerRetryAfter(t *testing.T) {
	for _, standalone := range []bool{false, true} {
		for _, test := range []struct {
			name                        string
			header, value               string
			minimum, maximum            time.Duration
			status                      int
			refused, canceled, deadline bool
		}{
			{name: "finite unauthorized", header: "Retry-After-Ms", value: "75", minimum: 75 * time.Millisecond, status: 401},
			{name: "finite ordinary retry", header: "Retry-After-Ms", value: "75", minimum: 75 * time.Millisecond, status: 503},
			{name: "overflow unauthorized", header: "Retry-After", value: "1e999", status: 401, refused: true},
			{name: "overflow ordinary retry", header: "Retry-After", value: "1e999", status: 503, refused: true},
			{name: "caller cap", header: "Retry-After-Ms", value: "75", maximum: time.Millisecond, status: 401, refused: true},
			{name: "canceled wait", header: "Retry-After", value: "1", status: 401, canceled: true},
			{name: "deadline wait", header: "Retry-After", value: "1", status: 401, deadline: true},
		} {
			if standalone && (test.status != 401 || test.maximum != 0) {
				continue
			}
			t.Run(fmt.Sprintf("standalone=%t/%s", standalone, test.name), func(t *testing.T) {
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				var issuerCalls atomic.Int32
				var issuerFailedAt, nextIssuerAt atomic.Int64
				backgroundDone := make(chan (<-chan struct{}), 1)
				apiCalls := 0
				next := func(req *http.Request) (*http.Response, error) {
					apiCalls++
					if apiCalls == 2 {
						select {
						case done := <-backgroundDone:
							select {
							case <-done:
							case <-time.After(5 * time.Second):
								return nil, errors.New("background exchange did not finish")
							}
						case <-time.After(5 * time.Second):
							return nil, errors.New("background exchange did not start")
						}
						if test.canceled {
							time.AfterFunc(time.Millisecond, cancel)
						}
						response := rootWorkloadResponse(test.status, `{"error":{"message":"synthetic rejected bearer"}}`)
						response.Header.Set("Retry-After-Ms", "0")
						return response, nil
					}
					if apiCalls == 3 && test.status == 503 {
						// Follow the ordinary retry with bearer rejection so completion
						// necessarily observes any new issuer exchange it started.
						response := rootWorkloadResponse(401, `{"error":{"message":"synthetic rejected bearer"}}`)
						response.Header.Set("Retry-After-Ms", "0")
						return response, nil
					}
					return rootWorkloadResponse(200, `{"data":[]}`), nil
				}
				httpClient := &http.Client{Transport: retryDelayRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
					if req.URL.Host != "auth.openai.com" {
						return next(req)
					}
					switch issuerCalls.Add(1) {
					case 1:
						return rootWorkloadResponse(200, `{"access_token":"synthetic-cached-bearer","expires_in":60}`), nil
					case 2:
						issuerFailedAt.Store(time.Now().UnixNano())
						backgroundDone <- req.Context().Done()
						response := rootWorkloadResponse(503, `{"error":"server_error","error_description":"synthetic-private-description"}`)
						response.Header.Set(test.header, test.value)
						return response, nil
					default:
						nextIssuerAt.Store(time.Now().UnixNano())
						return rootWorkloadResponse(200, `{"access_token":"synthetic-new-bearer","expires_in":3600}`), nil
					}
				})}
				provider := &mockSubjectTokenProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT}
				config := testWorkloadIdentity(provider)
				var call func(context.Context) error
				if standalone {
					identity, err := auth.NewWorkloadIdentityAuth(config)
					if err != nil {
						t.Fatal(err)
					}
					call = func(ctx context.Context) error {
						req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.openai.com/v1/models", nil)
						if err != nil {
							return err
						}
						response, err := auth.WorkloadIdentityMiddleware(identity, httpClient, req, next)
						if response != nil && response.Body != nil {
							_ = response.Body.Close()
						}
						return err
					}
				} else {
					retries := 1
					if test.status == 503 {
						retries = 2
					}
					opts := []option.RequestOption{option.WithWorkloadIdentity(config), option.WithHTTPClient(httpClient), option.WithMaxRetries(retries)}
					if test.maximum != 0 {
						opts = append(opts, option.WithMaxRetryDelay(test.maximum))
					}
					client := openai.NewClient(opts...)
					call = func(ctx context.Context) error { _, err := client.Models.List(ctx); return err }
				}
				if err := call(t.Context()); err != nil {
					t.Fatalf("prime cached bearer: %v", err)
				}
				if test.deadline {
					var deadlineCancel context.CancelFunc
					ctx, deadlineCancel = context.WithTimeout(ctx, 50*time.Millisecond)
					defer deadlineCancel()
				}
				err := call(ctx)
				if test.deadline {
					if !errors.Is(err, context.DeadlineExceeded) {
						t.Errorf("background retry error=%v, want context.DeadlineExceeded", err)
					}
				} else if test.canceled {
					if !errors.Is(err, context.Canceled) {
						t.Errorf("background retry error=%v, want context.Canceled", err)
					}
				} else if test.refused {
					if err == nil || err.Error() != "token exchange failed with status 503" {
						t.Errorf("background refusal error=%v, want original sanitized issuer error", err)
					}
				} else if err != nil {
					t.Errorf("background retry error=%v, want nil", err)
				}
				if err != nil && strings.Contains(err.Error(), "synthetic-private-description") {
					t.Error("background retry exposed issuer description")
				}
				wantIssuer := int32(3)
				if test.refused || test.canceled || test.deadline {
					wantIssuer = 2
				}
				if got := issuerCalls.Load(); got != wantIssuer {
					t.Errorf("background issuer attempts=%d, want %d", got, wantIssuer)
				}
				if at := nextIssuerAt.Load(); at != 0 && time.Duration(at-issuerFailedAt.Load()) < test.minimum {
					t.Errorf("background issuer wait=%s, want at least %s", time.Duration(at-issuerFailedAt.Load()), test.minimum)
				}
				// A fresh logical request must not inherit the participant's refusal.
				if test.refused {
					if err := call(t.Context()); err != nil {
						t.Errorf("independent request error=%v, want nil", err)
					}
					if got := issuerCalls.Load(); got != 3 {
						t.Errorf("independent request issuer attempts=%d, want 3", got)
					}
				}
			})
		}
	}
}

func TestWorkloadIdentityBackgroundIssuerParticipants(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	var issuerCalls atomic.Int32
	releaseIssuer := make(chan struct{})
	backgroundDone := make(chan (<-chan struct{}), 1)
	httpClient := &http.Client{Transport: retryDelayRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		switch issuerCalls.Add(1) {
		case 1:
			return rootWorkloadResponse(200, `{"access_token":"synthetic-cached-bearer","expires_in":60}`), nil
		case 2:
			backgroundDone <- req.Context().Done()
			select {
			case <-releaseIssuer:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
			response := rootWorkloadResponse(503, `{"error":"server_error"}`)
			response.Header.Set("Retry-After", "1e999")
			return response, nil
		default:
			return rootWorkloadResponse(200, `{"access_token":"synthetic-new-bearer","expires_in":3600}`), nil
		}
	})}
	provider := &mockSubjectTokenProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT}
	identity, err := auth.NewWorkloadIdentityAuth(testWorkloadIdentity(provider))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := identity.GetToken(ctx, httpClient); err != nil {
		t.Fatalf("prime cached bearer: %v", err)
	}
	arrived := make(chan struct{})
	completed := make(chan error, 1)
	var refreshDone <-chan struct{}
	invoke := func() {
		req, requestErr := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.openai.com/v1/models", nil)
		if requestErr != nil {
			select {
			case completed <- requestErr:
			case <-ctx.Done():
			}
			return
		}
		attempt := 0
		response, err := auth.WorkloadIdentityMiddleware(identity, httpClient, req, func(*http.Request) (*http.Response, error) {
			attempt++
			if attempt > 1 {
				return rootWorkloadResponse(200, `{"data":[]}`), nil
			}
			select {
			case arrived <- struct{}{}:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
			select {
			case <-refreshDone:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
			return rootWorkloadResponse(401, `{"error":{"message":"synthetic rejected bearer"}}`), nil
		})
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
		select {
		case completed <- err:
		case <-ctx.Done():
		}
	}
	// The first caller launches the background generation. Both callers must
	// receive its cached bearer before the issuer failure is released.
	go invoke()
	select {
	case refreshDone = <-backgroundDone:
	case <-ctx.Done():
		t.Fatal("background issuer did not start")
	}
	for participant := 0; participant < 2; participant++ {
		if participant == 1 {
			go invoke()
		}
		select {
		case <-arrived:
		case <-ctx.Done():
			t.Fatal("cached bearer participant did not arrive")
		}
	}
	close(releaseIssuer)
	for participant := 0; participant < 2; participant++ {
		select {
		case err := <-completed:
			if err == nil || err.Error() != "token exchange failed with status 503" {
				t.Errorf("participant %d error=%v, want original issuer refusal", participant, err)
			}
		case <-ctx.Done():
			t.Fatal("participant did not finish")
		}
	}
	if got := issuerCalls.Load(); got != 2 {
		t.Errorf("shared background issuer calls=%d, want 2", got)
	}
	if _, err := identity.GetToken(ctx, httpClient); err != nil {
		t.Errorf("independent GetToken error=%v, want nil", err)
	}
	if got := issuerCalls.Load(); got != 3 {
		t.Errorf("independent issuer calls=%d, want 3", got)
	}
}
