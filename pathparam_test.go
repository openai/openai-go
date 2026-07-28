package openai_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestPathParamsAreEscaped(t *testing.T) {
	t.Run("path traversal stays inside vector store id segment", func(t *testing.T) {
		req := captureRequest(t, func(client openai.Client) error {
			_, err := client.VectorStores.Get(context.Background(), "../videos/vid_123")
			return err
		}, vectorStoreResponse)

		if req.Method != http.MethodGet {
			t.Fatalf("method = %q, want GET", req.Method)
		}
		if got, want := req.URL.String(), "https://api.openai.com/v1/vector_stores/..%2Fvideos%2Fvid_123"; got != want {
			t.Fatalf("url = %q, want %q", got, want)
		}
	})

	t.Run("slash query and fragment stay inside vector store id segment", func(t *testing.T) {
		req := captureRequest(t, func(client openai.Client) error {
			_, err := client.VectorStores.Get(context.Background(), "vs_123/files/file_456?limit=1#frag")
			return err
		}, vectorStoreResponse)

		if got, want := req.URL.String(), "https://api.openai.com/v1/vector_stores/vs_123%2Ffiles%2Ffile_456%3Flimit=1%23frag"; got != want {
			t.Fatalf("url = %q, want %q", got, want)
		}
	})

	t.Run("encoded dot segments stay encoded", func(t *testing.T) {
		req := captureRequest(t, func(client openai.Client) error {
			_, err := client.VectorStores.Get(context.Background(), "%2e%2e/videos/vid_123")
			return err
		}, vectorStoreResponse)

		if got, want := req.URL.String(), "https://api.openai.com/v1/vector_stores/%252e%252e%2Fvideos%2Fvid_123"; got != want {
			t.Fatalf("url = %q, want %q", got, want)
		}
	})

	t.Run("bare dot segments stay inside vector store id segment", func(t *testing.T) {
		for input, wantURL := range map[string]string{
			".":  "https://api.openai.com/v1/vector_stores/%2E",
			"..": "https://api.openai.com/v1/vector_stores/%2E%2E",
		} {
			req := captureRequest(t, func(client openai.Client) error {
				_, err := client.VectorStores.Get(context.Background(), input)
				return err
			}, vectorStoreResponse)

			if got := req.URL.String(); got != wantURL {
				t.Fatalf("url = %q, want %q", got, wantURL)
			}
		}
	})

	t.Run("admin path traversal stays inside project id segment", func(t *testing.T) {
		req := captureRequest(t, func(client openai.Client) error {
			_, err := client.Admin.Organization.Projects.APIKeys.Delete(context.Background(), "proj_123/../../admin_api_keys/key_456?", "ignored")
			return err
		}, adminAPIKeyDeletedResponse)

		if req.Method != http.MethodDelete {
			t.Fatalf("method = %q, want DELETE", req.Method)
		}
		if got, want := req.URL.String(), "https://api.openai.com/v1/organization/projects/proj_123%2F..%2F..%2Fadmin_api_keys%2Fkey_456%3F/api_keys/ignored"; got != want {
			t.Fatalf("url = %q, want %q", got, want)
		}
	})
}

func captureRequest(t *testing.T, call func(openai.Client) error, responseBody string) *http.Request {
	t.Helper()

	var captured *http.Request
	client := openai.NewClient(
		option.WithBaseURL("https://api.openai.com/v1/"),
		option.WithAPIKey("sk-test"),
		option.WithAdminAPIKey("sk-admin-test"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				captured = req
				return &http.Response{
					StatusCode: http.StatusOK,
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body:    io.NopCloser(strings.NewReader(responseBody)),
					Request: req,
				}, nil
			}},
		}),
	)

	if err := call(client); err != nil {
		t.Fatalf("request failed: %s", err)
	}
	if captured == nil {
		t.Fatal("request was not captured")
	}
	return captured
}

const vectorStoreResponse = `{
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

const adminAPIKeyDeletedResponse = `{
	"id": "key",
	"object": "organization.project.api_key.deleted",
	"deleted": true
}`
