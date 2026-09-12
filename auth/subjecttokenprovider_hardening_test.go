package auth_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3/auth"
)

func TestCloudSubjectTokenProvidersBoundAndRedactMetadataResponses(t *testing.T) {
	const sensitive = "synthetic-private-metadata-access-token"
	for _, provider := range []struct {
		name     string
		identity string
		new      func() auth.SubjectTokenProvider
	}{
		{name: "Azure", identity: "azure-imds", new: func() auth.SubjectTokenProvider {
			return auth.AzureManagedIdentityTokenProvider(nil)
		}},
		{name: "GCP", identity: "gcp-metadata", new: func() auth.SubjectTokenProvider {
			return auth.GCPIDTokenProvider(nil)
		}},
	} {
		t.Run(provider.name, func(t *testing.T) {
			for _, test := range []struct {
				name          string
				status        int
				body          string
				contentLength int64
				readError     error
				canceled      bool
				want          string
				wantReads     bool
			}{
				{name: "private unsuccessful response", status: http.StatusUnauthorized,
					body: `{"access_token":"` + sensitive + `"}`, contentLength: -1,
					want: "status 401", wantReads: true},
				{name: "declared oversized unsuccessful response", status: http.StatusForbidden,
					contentLength: (4 << 10) + 1, want: "size limit"},
				{name: "streamed oversized unsuccessful response", status: http.StatusInternalServerError,
					body: strings.Repeat("a", (4<<10)+1), contentLength: -1,
					want: "size limit", wantReads: true},
				{name: "declared oversized successful response", status: http.StatusOK,
					contentLength: (1 << 20) + 1, want: "size limit"},
				{name: "streamed oversized successful response", status: http.StatusOK,
					body: strings.Repeat("a", (1<<20)+1), contentLength: -1,
					want: "size limit", wantReads: true},
				{name: "sensitive successful response read failure", status: http.StatusOK,
					contentLength: -1, readError: errors.New(sensitive),
					want: "failed to read", wantReads: true},
				{name: "sensitive unsuccessful response read failure", status: http.StatusForbidden,
					contentLength: -1, readError: errors.New(sensitive),
					want: "failed to read", wantReads: true},
				{name: "canceled metadata response read", status: http.StatusOK,
					contentLength: -1, readError: errors.New(sensitive), canceled: true,
					want: "failed to read", wantReads: true},
			} {
				t.Run(test.name, func(t *testing.T) {
					ctx, cancel := context.WithCancel(t.Context())
					defer cancel()
					body := &observedMetadataBody{
						Reader: strings.NewReader(test.body),
						err:    test.readError,
					}
					if test.canceled {
						body.cancel = cancel
					}
					doer := ordinaryOpaqueHTTPDoer(func(*http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode:    test.status,
							Header:        make(http.Header),
							Body:          body,
							ContentLength: test.contentLength,
						}, nil
					})

					token, err := provider.new().GetToken(ctx, doer)
					var typed *auth.SubjectTokenProviderError
					if token != "" || !errors.As(err, &typed) || typed.Provider != provider.identity ||
						!strings.Contains(err.Error(), test.want) {
						t.Fatalf("metadata token=%q error=%v, want provider %q and %q",
							token, err, provider.identity, test.want)
					}
					if strings.Contains(err.Error(), sensitive) {
						t.Errorf("metadata provider error disclosed its private response: %q", err.Error())
					}
					if got := body.reads.Load() != 0; got != test.wantReads {
						t.Errorf("metadata response was read=%t, want %t", got, test.wantReads)
					}
					if got := body.closes.Load(); got != 1 {
						t.Errorf("metadata response body closes=%d, want exactly one", got)
					}
					if got := errors.Is(err, context.Canceled); got != test.canceled {
						t.Errorf("metadata response preserved caller cancellation=%t, want %t", got, test.canceled)
					}
				})
			}

			t.Run("missing metadata response body", func(t *testing.T) {
				doer := ordinaryOpaqueHTTPDoer(func(*http.Request) (*http.Response, error) {
					return &http.Response{StatusCode: http.StatusOK}, nil
				})
				token, err := provider.new().GetToken(t.Context(), doer)
				var typed *auth.SubjectTokenProviderError
				if token != "" || !errors.As(err, &typed) || typed.Provider != provider.identity ||
					!strings.Contains(err.Error(), "invalid response") {
					t.Errorf("missing metadata response body token=%q error=%v", token, err)
				}
			})
		})
	}
}

func TestAzureSubjectTokenProviderSanitizesInvalidMetadataJSON(t *testing.T) {
	const sensitive = "synthetic-private-malformed-access-token"
	doer := ordinaryOpaqueHTTPDoer(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(fmt.Sprintf(`{"access_token":%q,`, sensitive))),
		}, nil
	})
	token, err := auth.AzureManagedIdentityTokenProvider(nil).GetToken(t.Context(), doer)
	var typed *auth.SubjectTokenProviderError
	if token != "" || !errors.As(err, &typed) || typed.Provider != "azure-imds" ||
		!strings.Contains(err.Error(), "failed to decode") {
		t.Fatalf("malformed Azure metadata token=%q error=%v", token, err)
	}
	if strings.Contains(err.Error(), sensitive) || typed.Cause != nil {
		t.Errorf("malformed Azure metadata leaked private decoder details: %v", err)
	}
}

func TestCloudSubjectTokenProvidersNeverFollowMetadataRedirects(t *testing.T) {
	for _, provider := range []struct {
		name        string
		identity    string
		metadataKey string
		metadataVal string
		new         func() auth.SubjectTokenProvider
	}{
		{name: "Azure", identity: "azure-imds", metadataKey: "Metadata", metadataVal: "true",
			new: func() auth.SubjectTokenProvider { return auth.AzureManagedIdentityTokenProvider(nil) }},
		{name: "GCP", identity: "gcp-metadata", metadataKey: "Metadata-Flavor", metadataVal: "Google",
			new: func() auth.SubjectTokenProvider { return auth.GCPIDTokenProvider(nil) }},
	} {
		t.Run(provider.name, func(t *testing.T) {
			for _, status := range []int{
				http.StatusMovedPermanently,
				http.StatusFound,
				http.StatusSeeOther,
				http.StatusTemporaryRedirect,
				http.StatusPermanentRedirect,
			} {
				t.Run(fmt.Sprintf("HTTP %d", status), func(t *testing.T) {
					var metadataRequests, redirectedRequests, callerRedirectChecks atomic.Int32
					target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						redirectedRequests.Add(1)
						_, _ = io.WriteString(w, "synthetic-attacker-controlled-subject")
					}))
					t.Cleanup(target.Close)
					metadata := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
						metadataRequests.Add(1)
						if got := request.Header.Get(provider.metadataKey); got != provider.metadataVal {
							t.Errorf("metadata proof header %q=%q, want %q", provider.metadataKey, got, provider.metadataVal)
						}
						w.Header().Set("Location", target.URL+"/steal?credential=synthetic-private-metadata-location")
						w.WriteHeader(status)
					}))
					t.Cleanup(metadata.Close)
					caller := &http.Client{
						Transport: &closureTransport{fn: func(request *http.Request) (*http.Response, error) {
							if request.URL.Host == "169.254.169.254" || request.URL.Host == "metadata.google.internal" {
								request = request.Clone(request.Context())
								request.URL.Scheme = "http"
								request.URL.Host = metadata.Listener.Addr().String()
							}
							return http.DefaultTransport.RoundTrip(request)
						}},
						CheckRedirect: func(*http.Request, []*http.Request) error {
							callerRedirectChecks.Add(1)
							return nil
						},
					}

					token, err := provider.new().GetToken(t.Context(), caller)
					var typed *auth.SubjectTokenProviderError
					if token != "" || !errors.As(err, &typed) || typed.Provider != provider.identity ||
						!strings.Contains(err.Error(), "does not follow redirects") {
						t.Fatalf("redirected metadata token=%q error=%v", token, err)
					}
					if strings.Contains(err.Error(), "synthetic-private") || strings.Contains(err.Error(), target.URL) {
						t.Errorf("metadata redirect error leaked its destination: %q", err.Error())
					}
					if metadataRequests.Load() != 1 || redirectedRequests.Load() != 0 ||
						callerRedirectChecks.Load() != 0 {
						t.Errorf("metadata/redirect/caller-policy requests=%d/%d/%d, want 1/0/0",
							metadataRequests.Load(), redirectedRequests.Load(), callerRedirectChecks.Load())
					}
					if caller.CheckRedirect == nil {
						t.Error("metadata provider mutated its caller-owned native HTTP client")
					}
				})
			}
		})
	}
}

type observedMetadataBody struct {
	io.Reader
	err    error
	cancel context.CancelFunc
	reads  atomic.Int32
	closes atomic.Int32
}

func (body *observedMetadataBody) Read(buffer []byte) (int, error) {
	body.reads.Add(1)
	if body.cancel != nil {
		body.cancel()
	}
	if body.err != nil {
		return 0, body.err
	}
	return body.Reader.Read(buffer)
}

func (body *observedMetadataBody) Close() error {
	body.closes.Add(1)
	return nil
}
