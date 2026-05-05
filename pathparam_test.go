// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openai/openai-go/v3/option"
)

func TestPathSegment(t *testing.T) {
	tests := map[string]string{
		"id/with/slash":          "id%2Fwith%2Fslash",
		"id?query=injected":      "id%3Fquery=injected",
		"id#fragment":            "id%23fragment",
		"../videos/vid_123":      "%2E%2E%2Fvideos%2Fvid_123",
		"%2e%2e/videos/vid_123": "%252e%252e%2Fvideos%2Fvid_123",
		".":                      "%2E",
		"..":                     "%2E%2E",
	}

	for value, want := range tests {
		if got := pathSegment(value); got != want {
			t.Fatalf("pathSegment(%q) = %q, want %q", value, got, want)
		}
	}
}

func TestPathParamsAreEscapedInRequests(t *testing.T) {
	tests := map[string]string{
		"id/with/slash":          "/vector_stores/id%2Fwith%2Fslash",
		"id?query=injected":      "/vector_stores/id%3Fquery=injected",
		"id#fragment":            "/vector_stores/id%23fragment",
		"../videos/vid_123":      "/vector_stores/%2E%2E%2Fvideos%2Fvid_123",
		"%2e%2e/videos/vid_123": "/vector_stores/%252e%252e%2Fvideos%2Fvid_123",
	}

	for vectorStoreID, wantPath := range tests {
		t.Run(vectorStoreID, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.EscapedPath() != wantPath {
					t.Fatalf("request path = %q, want %q", r.URL.EscapedPath(), wantPath)
				}
				if r.URL.RawQuery != "" {
					t.Fatalf("request query = %q, want empty", r.URL.RawQuery)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"id":"vs_123","created_at":0,"file_counts":{},"last_active_at":0,"metadata":{},"name":"test","object":"vector_store","status":"completed","usage_bytes":0}`))
			}))
			defer server.Close()

			client := NewClient(
				option.WithBaseURL(server.URL),
				option.WithAPIKey("My API Key"),
			)
			_, err := client.VectorStores.Get(context.TODO(), vectorStoreID)
			if err != nil {
				t.Fatalf("err should be nil: %s", err.Error())
			}
		})
	}
}
