package auth_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

func TestWorkloadIdentityCloudMetadataFailuresRemainPrivate(t *testing.T) {
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
		for _, failure := range []struct {
			name      string
			redirect  bool
			oversized bool
		}{
			{name: "private unsuccessful metadata response"},
			{name: "oversized unsuccessful metadata response", oversized: true},
			{name: "cross-origin metadata redirect", redirect: true},
		} {
			t.Run(provider.name+" "+failure.name, func(t *testing.T) {
				var metadataCalls, issuerCalls, apiCalls, redirectedCalls atomic.Int32
				target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					redirectedCalls.Add(1)
					_, _ = io.WriteString(w, "synthetic-attacker-subject-token")
				}))
				t.Cleanup(target.Close)
				metadata := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					metadataCalls.Add(1)
					if failure.redirect {
						w.Header().Set("Location", target.URL+"/steal?credential=synthetic-private-metadata")
						w.WriteHeader(http.StatusTemporaryRedirect)
						return
					}
					w.WriteHeader(http.StatusUnauthorized)
					body := `{"access_token":"synthetic-private-metadata-token"}`
					if failure.oversized {
						body = strings.Repeat("a", (4<<10)+1)
					}
					_, _ = io.WriteString(w, body)
				}))
				t.Cleanup(metadata.Close)
				httpClient := &http.Client{Transport: &closureTransport{fn: func(request *http.Request) (*http.Response, error) {
					switch request.URL.Host {
					case "169.254.169.254", "metadata.google.internal":
						t.Error("metadata reached the API transport")
						return nil, errors.New("unexpected metadata request")
					case "auth.openai.com":
						issuerCalls.Add(1)
						return ordinaryWorkloadResponse(http.StatusOK,
							`{"access_token":"synthetic-workload-bearer","expires_in":3600}`), nil
					default:
						apiCalls.Add(1)
						resp := ordinaryWorkloadResponse(http.StatusOK, `{"data":[]}`)
						resp.Header.Set("Content-Type", "application/json")
						return resp, nil
					}
				}}}
				retries := 0
				if failure.redirect {
					retries = 2
				}
				client := openai.NewClient(
					option.WithHTTPClient(httpClient),
					option.WithWorkloadIdentity(auth.WorkloadIdentity{
						IdentityProviderID: "synthetic-provider",
						ServiceAccountID:   "synthetic-account",
						Provider:           auth.MetadataProviderWithDialForTest(t, provider.new(), metadata.Listener.Addr().String()),
					}),
					option.WithMaxRetries(retries),
				)

				_, err := client.Models.List(t.Context())
				var typed *auth.SubjectTokenProviderError
				if !errors.As(err, &typed) || typed.Provider != provider.identity {
					t.Fatalf("public metadata failure lost its provider identity: %v", err)
				}
				if strings.Contains(err.Error(), "synthetic-private") || strings.Contains(err.Error(), target.URL) {
					t.Errorf("public metadata failure exposed credentials: %q", err.Error())
				}
				if failure.redirect && !strings.Contains(err.Error(), "does not follow redirects") {
					t.Errorf("public metadata redirect error=%v", err)
				}
				if failure.oversized && !strings.Contains(err.Error(), "size limit") {
					t.Errorf("public oversized metadata error=%v", err)
				}
				if metadataCalls.Load() != 1 || issuerCalls.Load() != 0 || apiCalls.Load() != 0 ||
					redirectedCalls.Load() != 0 {
					t.Errorf("metadata/issuer/API/redirect requests=%d/%d/%d/%d, want 1/0/0/0",
						metadataCalls.Load(), issuerCalls.Load(), apiCalls.Load(), redirectedCalls.Load())
				}
			})
		}
	}
}

func TestWorkloadIdentitySeparatesMetadataAndAPITransports(t *testing.T) {
	for _, gcp := range []bool{false, true} {
		for _, opaque := range []bool{false, true} {
			var metadataCalls, issuerCalls, apiCalls atomic.Int32
			metadata := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				metadataCalls.Add(1)
				w.Header().Set("Metadata-Flavor", "Google")
				if gcp {
					_, _ = io.WriteString(w, "synthetic-subject")
				} else {
					_, _ = io.WriteString(w, `{"access_token":"synthetic-subject"}`)
				}
			}))
			t.Cleanup(metadata.Close)
			provider := auth.AzureManagedIdentityTokenProvider(nil)
			if gcp {
				provider = auth.GCPIDTokenProvider(nil)
			}
			provider = auth.MetadataProviderWithDialForTest(t, provider, metadata.Listener.Addr().String())
			doer := ordinaryOpaqueHTTPDoer(func(r *http.Request) (*http.Response, error) {
				switch r.URL.Host {
				case "auth.openai.com":
					issuerCalls.Add(1)
					var body struct {
						SubjectToken string `json:"subject_token"`
					}
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.SubjectToken != "synthetic-subject" {
						t.Errorf("unexpected token exchange: %v", err)
					}
					return ordinaryWorkloadResponse(http.StatusOK, `{"access_token":"synthetic-bearer","expires_in":3600}`), nil
				case "api.openai.com":
					apiCalls.Add(1)
					if r.Header.Get("Authorization") != "Bearer synthetic-bearer" {
						t.Error("API request lost workload authentication")
					}
					resp := ordinaryWorkloadResponse(http.StatusOK, `{"data":[]}`)
					resp.Header.Set("Content-Type", "application/json")
					return resp, nil
				default:
					t.Error("API client received a metadata request")
					return nil, errors.New("unexpected destination")
				}
			})
			var apiClient auth.HTTPDoer = doer
			if !opaque {
				apiClient = &http.Client{Transport: &closureTransport{fn: doer}}
			}
			client := openai.NewClient(option.WithHTTPClient(apiClient), option.WithMaxRetries(0), option.WithWorkloadIdentity(auth.WorkloadIdentity{
				IdentityProviderID: "synthetic-provider", ServiceAccountID: "synthetic-account", Provider: provider,
			}))
			if _, err := client.Models.List(t.Context()); err != nil {
				t.Fatal(err)
			}
			if metadataCalls.Load() != 1 || issuerCalls.Load() != 1 || apiCalls.Load() != 1 {
				t.Fatalf("metadata/issuer/API calls=%d/%d/%d", metadataCalls.Load(), issuerCalls.Load(), apiCalls.Load())
			}
		}
	}
}
