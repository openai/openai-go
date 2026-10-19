package openai_test

import (
	"context"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type multipartScalarTransport func(*http.Request) (*http.Response, error)

func (f multipartScalarTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestMultipartExtraFieldNames(t *testing.T) {
	tests := []struct {
		name, key, wantName string
		wantError           bool
	}{
		{"CR", "extra\rfield", "extra%0Dfield", false},
		{"LF", "extra\nfield", "extra%0Afield", false},
		{"header and body", "extra\r\nInjected-Header: yes\r\n\r\nbody", "extra%0D%0AInjected-Header: yes%0D%0A%0D%0Abody", false},
		{"quoted Unicode", "extra[\"\\雪\t%0D]", "extra[\"\\雪\t%0D]", false},
		{"empty", "", "", false},
		{"NUL", "extra\x00field", "", true},
		{"unit separator", "extra\x1ffield", "", true},
		{"DEL", "extra\x7ffield", "", true},
		{"Unicode control", "extra\u0085field", "", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requests := 0
			transport := multipartScalarTransport(func(req *http.Request) (*http.Response, error) {
				requests++
				_, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
				if err != nil {
					t.Fatal(err)
				}
				reader := multipart.NewReader(req.Body, params["boundary"])
				parts := 0
				for {
					part, err := reader.NextPart()
					if errors.Is(err, io.EOF) {
						break
					}
					if err != nil {
						t.Fatal(err)
					}
					parts++
					data, err := io.ReadAll(part)
					if err != nil {
						t.Fatal(err)
					}
					if part.FormName() == "purpose" {
						if string(data) != "assistants" {
							t.Fatal("ordinary scalar value changed")
						}
					} else if part.FormName() != test.wantName || string(data) != "value\r\nkept" {
						t.Fatalf("unexpected part name or body: %q, %q", part.FormName(), data)
					}
					if len(part.Header) != 1 || part.FileName() != "" {
						t.Fatal("field name changed part headers")
					}
				}
				if parts != 2 {
					t.Fatalf("part count = %d, want 2", parts)
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"id":"file-test"}`)),
					Request:    req,
				}, nil
			})
			service := openai.NewFileService(
				option.WithBaseURL("https://sdk.test/"), option.WithAPIKey("fake-key"),
				option.WithHTTPClient(&http.Client{Transport: transport}), option.WithMaxRetries(0),
			)
			params := openai.FileNewParams{Purpose: openai.FilePurposeAssistants}
			params.SetExtraFields(map[string]any{test.key: "value\r\nkept"})
			_, err := service.New(context.Background(), params)
			if test.wantError {
				if err == nil || !strings.Contains(err.Error(), "invalid multipart field name") {
					t.Fatalf("error = %v, want invalid field name", err)
				}
				if requests != 0 {
					t.Fatal("invalid field name reached transport")
				}
			} else if err != nil || requests != 1 {
				t.Fatalf("error = %v, requests = %d, want one successful request", err, requests)
			}
		})
	}
}
