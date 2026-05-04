package requestconfig

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/openai/openai-go/v3/internal/apierror"
)

func TestExecuteParsesTopLevelErrorPayload(t *testing.T) {
	t.Parallel()

	err := executeErrorRequest(t, http.StatusNotFound, `{"code":"not_found","message":"missing resource"}`)
	var apiErr *apierror.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected apierror.Error, got %T", err)
	}
	if apiErr.Code != "not_found" {
		t.Fatalf("expected code %q, got %q", "not_found", apiErr.Code)
	}
	if apiErr.Message != "missing resource" {
		t.Fatalf("expected message %q, got %q", "missing resource", apiErr.Message)
	}
}

func TestExecuteParsesNestedErrorPayload(t *testing.T) {
	t.Parallel()

	err := executeErrorRequest(t, http.StatusBadRequest, `{"error":{"code":"bad_request","message":"invalid input"}}`)
	var apiErr *apierror.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected apierror.Error, got %T", err)
	}
	if apiErr.Code != "bad_request" {
		t.Fatalf("expected code %q, got %q", "bad_request", apiErr.Code)
	}
	if apiErr.Message != "invalid input" {
		t.Fatalf("expected message %q, got %q", "invalid input", apiErr.Message)
	}
}

func executeErrorRequest(t *testing.T, statusCode int, payload string) error {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = io.WriteString(w, payload)
	}))
	defer server.Close()

	cfg, err := NewRequestConfig(context.Background(), http.MethodGet, "/test", nil, nil)
	if err != nil {
		t.Fatalf("NewRequestConfig() error = %v", err)
	}

	baseURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}

	cfg.BaseURL = baseURL
	cfg.HTTPClient = server.Client()
	cfg.MaxRetries = 0

	err = cfg.Execute()
	if err == nil {
		t.Fatal("expected Execute() to return an error")
	}

	return err
}
