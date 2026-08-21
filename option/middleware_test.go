package option

import (
	"net/http"
	"reflect"
	"testing"
)

func TestDebugLogMethod(t *testing.T) {
	for _, method := range []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace,
	} {
		if got := debugLogMethod(method); got != method {
			t.Errorf("debugLogMethod(%q) = %q, want unchanged", method, got)
		}
	}

	if got := debugLogMethod("CUSTOM-METHOD-SECRET"); got != debugLogRedacted {
		t.Errorf("debugLogMethod() = %q, want %q", got, debugLogRedacted)
	}
}

func TestDebugLogHeaders(t *testing.T) {
	headers := http.Header{
		"Content-Type":         {"application/json"},
		"Proxy-Authorization":  {"first-secret", "second-secret"},
		"X-Amz-Security-Token": {"session-secret"},
		"X-Custom-Secret":      {"custom-secret"},
	}
	want := http.Header{
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
