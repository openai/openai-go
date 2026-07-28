package option

import (
	"net/http"
	"testing"
)

func TestRedactDebugHeadersRedactsAWSSessionToken(t *testing.T) {
	headers := http.Header{
		"X-Amz-Security-Token": {"secret-session-token"},
		"X-Amz-Date":           {"20250102T030405Z"},
	}
	redacted := redactDebugHeaders(headers)
	if got := redacted.Get("X-Amz-Security-Token"); got != "***" {
		t.Fatalf("X-Amz-Security-Token = %q", got)
	}
	if got := redacted.Get("X-Amz-Date"); got != "20250102T030405Z" {
		t.Fatalf("X-Amz-Date = %q", got)
	}
	if got := headers.Get("X-Amz-Security-Token"); got != "secret-session-token" {
		t.Fatalf("original header was modified: %q", got)
	}
}
