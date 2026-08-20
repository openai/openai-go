package option

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestDebugLogHeaders(t *testing.T) {
	headers := http.Header{
		"Content-Type":         {"application/json"},
		"Proxy-Authorization":  {"first-secret", "second-secret"},
		"X-Amz-Security-Token": {"session-secret"},
		"X-Custom-Secret":      {"custom-secret"},
	}
	want := http.Header{
		"Content-Type":         {"application/json"},
		"Proxy-Authorization":  {"***", "***"},
		"X-Amz-Security-Token": {"***"},
	}

	if got := debugLogHeaders(headers); !reflect.DeepEqual(got, want) {
		t.Fatalf("debugLogHeaders() = %#v, want %#v", got, want)
	}
	if got := headers.Values("Proxy-Authorization"); !reflect.DeepEqual(got, []string{"first-secret", "second-secret"}) {
		t.Fatalf("original headers were modified: %#v", got)
	}
}

func TestDebugLogURL(t *testing.T) {
	u := &url.URL{
		Scheme:   "https",
		User:     url.UserPassword("url-user", "url-password"),
		Host:     "example.com",
		Path:     "/v1/responses",
		RawQuery: "api_key=query-secret&access_token=token-secret",
		Fragment: "fragment-secret",
		Opaque:   "opaque-secret",
	}

	if got, want := debugLogURL(u), "https://example.com/v1/responses"; got != want {
		t.Fatalf("debugLogURL() = %q, want %q", got, want)
	}
}
