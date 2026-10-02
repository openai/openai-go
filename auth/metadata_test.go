package auth

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func localMetadataProvider(t *testing.T, gcp bool, handler http.HandlerFunc) (SubjectTokenProvider, *atomic.Int32) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	provider := AzureManagedIdentityTokenProvider(nil)
	if gcp {
		provider = GCPIDTokenProvider(nil)
	}
	var dials atomic.Int32
	client := metadataClientForTest(provider)
	client.Transport.(*http.Transport).DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		dials.Add(1)
		if network != "tcp" || address != "169.254.169.254:80" {
			t.Errorf("metadata dial = %s %s, want tcp %s", network, address, "169.254.169.254:80")
			return nil, errors.New("unexpected test destination")
		}
		return (&net.Dialer{}).DialContext(ctx, network, server.Listener.Addr().String())
	}
	t.Cleanup(client.CloseIdleConnections)
	return MetadataProviderWithClientForTest(provider, client), &dials
}

func TestMetadataIgnoresAPIAndGlobalTransports(t *testing.T) {
	var unwanted atomic.Int32
	trap := metadataDoerFunc(func(*http.Request) (*http.Response, error) {
		unwanted.Add(1)
		return nil, errors.New("API or global transport received metadata")
	})
	proxy := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { unwanted.Add(1) }))
	t.Cleanup(proxy.Close)
	for _, name := range []string{"HTTP_PROXY", "http_proxy", "HTTPS_PROXY", "https_proxy"} {
		t.Setenv(name, proxy.URL)
	}
	t.Setenv("NO_PROXY", "")
	t.Setenv("no_proxy", "")
	oldClient, oldTransport := http.DefaultClient, http.DefaultTransport
	http.DefaultClient = &http.Client{Transport: trap}
	http.DefaultTransport = trap
	t.Cleanup(func() { http.DefaultClient, http.DefaultTransport = oldClient, oldTransport })

	for _, gcp := range []bool{false, true} {
		provider, dials := localMetadataProvider(t, gcp, func(w http.ResponseWriter, r *http.Request) {
			if gcp {
				if r.Host != "169.254.169.254" || r.Header.Get("Metadata-Flavor") != "Google" || r.URL.Query().Get("audience") != DefaultAudience || r.URL.Path != "/computeMetadata/v1/instance/service-accounts/default/identity" {
					t.Errorf("GCP metadata request lost its host, path, header, or audience")
				}
				w.Header().Set("Metadata-Flavor", "Google")
				_, _ = io.WriteString(w, "synthetic-subject")
			} else {
				if r.Host != "169.254.169.254" || r.Header.Get("Metadata") != "true" || r.URL.Query().Get("resource") != DefaultAzureResource || r.URL.Path != "/metadata/identity/oauth2/token" {
					t.Errorf("Azure metadata request lost its host, path, header, or resource")
				}
				_, _ = io.WriteString(w, `{"access_token":"synthetic-subject"}`)
			}
		})
		for _, apiClient := range []HTTPDoer{nil, &http.Client{Transport: trap}, trap} {
			token, err := provider.GetToken(t.Context(), apiClient)
			if err != nil || token != "synthetic-subject" {
				t.Fatalf("metadata token=%q error=%v", token, err)
			}
		}
		if dials.Load() == 0 {
			t.Fatal("metadata retrieval did not use the dedicated transport")
		}
	}
	if unwanted.Load() != 0 {
		t.Fatalf("metadata used API/global/proxy transport %d times", unwanted.Load())
	}
}

func TestGCPMetadataRequiresResponseFlavor(t *testing.T) {
	for _, values := range [][]string{nil, {"Other"}, {"Google", "Other"}, {"Google", "Google"}, {"Google, Other"}, {"Google"}} {
		t.Run(strings.Join(values, ","), func(t *testing.T) {
			provider, _ := localMetadataProvider(t, true, func(w http.ResponseWriter, _ *http.Request) {
				for _, value := range values {
					w.Header().Add("Metadata-Flavor", value)
				}
				_, _ = io.WriteString(w, "synthetic-subject")
			})
			token, err := provider.GetToken(t.Context(), nil)
			valid := len(values) == 1 && values[0] == "Google"
			if valid && (err != nil || token != "synthetic-subject") {
				t.Fatalf("valid metadata rejected: %v", err)
			}
			if !valid && (err == nil || token != "") {
				t.Fatal("invalid metadata flavor accepted")
			}
		})
	}
}

func TestMetadataDeadlinesCoverHeadersAndBody(t *testing.T) {
	for _, gcp := range []bool{false, true} {
		for _, body := range []bool{false, true} {
			for _, callerDeadline := range []bool{false, true} {
				name := "Azure"
				if gcp {
					name = "GCP"
				}
				if body {
					name += "/body"
				} else {
					name += "/headers"
				}
				if callerDeadline {
					name += "/caller"
				} else {
					name += "/default"
				}
				t.Run(name, func(t *testing.T) {
					t.Parallel()
					provider, _ := localMetadataProvider(t, gcp, func(w http.ResponseWriter, r *http.Request) {
						if body {
							w.Header().Set("Metadata-Flavor", "Google")
							w.WriteHeader(http.StatusOK)
							w.(http.Flusher).Flush()
						}
						<-r.Context().Done()
					})
					ctx := t.Context()
					if callerDeadline {
						var cancel context.CancelFunc
						ctx, cancel = context.WithTimeout(ctx, 100*time.Millisecond)
						defer cancel()
					}
					started := time.Now()
					token, err := provider.GetToken(ctx, nil)
					if token != "" || !errors.Is(err, context.DeadlineExceeded) {
						t.Fatalf("deadline token=%q err=%v", token, err)
					}
					limit := metadataRequestTimeout + 2*time.Second
					if callerDeadline {
						limit = 2 * time.Second
					}
					if time.Since(started) > limit {
						t.Fatal("metadata deadline did not bound the request")
					}
				})
			}
		}
	}
}

type metadataDoerFunc func(*http.Request) (*http.Response, error)

func (f metadataDoerFunc) Do(r *http.Request) (*http.Response, error)        { return f(r) }
func (f metadataDoerFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
