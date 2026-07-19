package apierror

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDumpRequestRedactsSensitiveHeaders(t *testing.T) {
	body := []byte(`{"message":"hello"}`)
	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}
	req.Header.Set("Authorization", "Bearer sk-test")
	req.Header.Set("Api-Key", "azure-key")
	req.Header.Set("OpenAI-Organization", "org_123")
	req.Header.Set("OpenAI-Project", "proj_123")
	req.Header.Set("Webhook-Signature", "v1,signature")
	req.Header.Set("Webhook-Id", "wh_123")
	req.Header.Set("Webhook-Timestamp", "123")
	req.Header.Set("X-Custom", "keep-me")

	apierr := &Error{Request: req}
	dump := strings.ToLower(string(apierr.DumpRequest(true)))

	for _, want := range []string{
		"authorization: ***",
		"api-key: ***",
		"openai-organization: ***",
		"openai-project: ***",
		"webhook-signature: ***",
		"webhook-id: ***",
		"webhook-timestamp: ***",
	} {
		if !strings.Contains(dump, want) {
			t.Fatalf("request dump missing redacted header %q:\n%s", want, dump)
		}
	}
	if strings.Contains(dump, "sk-test") || strings.Contains(dump, "azure-key") || strings.Contains(dump, "org_123") {
		t.Fatalf("request dump leaked a sensitive header:\n%s", dump)
	}
	if !strings.Contains(dump, "x-custom: keep-me") {
		t.Fatalf("request dump dropped non-sensitive header:\n%s", dump)
	}
	if !strings.Contains(dump, `{"message":"hello"}`) {
		t.Fatalf("request dump missing body:\n%s", dump)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer sk-test" {
		t.Fatalf("request header mutated, got %q", got)
	}
}

func TestDumpResponseRedactsSensitiveHeaders(t *testing.T) {
	resp := &http.Response{
		Status:     "401 Unauthorized",
		StatusCode: http.StatusUnauthorized,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: http.Header{
			"Authorization":       []string{"Bearer sk-test"},
			"Set-Cookie":          []string{"session=secret"},
			"OpenAI-Organization": []string{"org_123"},
			"OpenAI-Project":      []string{"proj_123"},
			"Webhook-Signature":   []string{"v1,signature"},
			"X-Custom":            []string{"keep-me"},
		},
		Body: io.NopCloser(strings.NewReader(`{"error":"denied"}`)),
	}

	apierr := &Error{Response: resp}
	dump := strings.ToLower(string(apierr.DumpResponse(true)))

	for _, want := range []string{
		"authorization: ***",
		"set-cookie: ***",
		"openai-organization: ***",
		"openai-project: ***",
		"webhook-signature: ***",
	} {
		if !strings.Contains(dump, want) {
			t.Fatalf("response dump missing redacted header %q:\n%s", want, dump)
		}
	}
	if strings.Contains(dump, "sk-test") || strings.Contains(dump, "session=secret") || strings.Contains(dump, "org_123") {
		t.Fatalf("response dump leaked a sensitive header:\n%s", dump)
	}
	if !strings.Contains(dump, "x-custom: keep-me") {
		t.Fatalf("response dump dropped non-sensitive header:\n%s", dump)
	}
	if !strings.Contains(dump, `{"error":"denied"}`) {
		t.Fatalf("response dump missing body:\n%s", dump)
	}
	if got := resp.Header.Get("Authorization"); got != "Bearer sk-test" {
		t.Fatalf("response header mutated, got %q", got)
	}
}
