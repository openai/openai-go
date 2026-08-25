package option

import (
	"net/http"
	"net/url"
	"testing"
)

func TestX509WorkloadIdentityRequiresExactGlobalBaseURL(t *testing.T) {
	for _, test := range []struct {
		name   string
		target string
		valid  bool
	}{
		{name: "exact global endpoint", target: "https://mtls.api.openai.com/v1/", valid: true},
		{name: "non-mTLS global endpoint", target: "https://api.openai.com/v1/"},
		{name: "regional endpoint", target: "https://eu.api.openai.com/v1/"},
		{name: "attacker suffix", target: "https://mtls.api.openai.com.attacker.test/v1/"},
		{name: "explicit port", target: "https://mtls.api.openai.com:443/v1/"},
		{name: "plaintext", target: "http://mtls.api.openai.com/v1/"},
		{name: "credentials", target: "https://user@mtls.api.openai.com/v1/"},
		{name: "query", target: "https://mtls.api.openai.com/v1/?signature=synthetic"},
		{name: "fragment", target: "https://mtls.api.openai.com/v1/#unsafe"},
		{name: "wrong path", target: "https://mtls.api.openai.com/other/"},
	} {
		t.Run(test.name, func(t *testing.T) {
			target, err := url.Parse(test.target)
			if err != nil {
				t.Fatalf("parse synthetic endpoint: %v", err)
			}
			if got := validX509WorkloadAPIBaseURL(target); got != test.valid {
				t.Errorf("endpoint validation = %v, want %v", got, test.valid)
			}
		})
	}
	if validX509WorkloadAPIBaseURL(nil) {
		t.Error("nil endpoint unexpectedly validated")
	}
}

func TestX509WorkloadIdentityDetectsCredentialHeaderAliases(t *testing.T) {
	for _, name := range []string{
		"Authorization", "authorization", "Api-Key", "api_key", "X-Api-Key", "x_api_key",
		"Proxy-Authorization", "proxy_authorization", "Cookie", "Set-Cookie", "Host",
		"X-Amz-Security-Token", ":authority",
	} {
		t.Run(name, func(t *testing.T) {
			if !unsafeX509CredentialHeaders(http.Header{name: []string{"synthetic-credential"}}) {
				t.Errorf("credential header alias %q was accepted", name)
			}
		})
	}
	if unsafeX509CredentialHeaders(http.Header{"OpenAI-Organization": []string{"synthetic-organization"}}) {
		t.Error("safe API organization metadata was rejected")
	}
}

func TestX509WorkloadIdentityRedactsResponseRequestCredentials(t *testing.T) {
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet,
		"https://mtls.api.openai.com/v1/models", nil)
	if err != nil {
		t.Fatalf("construct synthetic response request: %v", err)
	}
	request.Header = http.Header{
		"Authorization":       {"Bearer synthetic-token"},
		"authorization":       {"Bearer synthetic-alias"},
		"proxy_authorization": {"synthetic-proxy"},
		"api_key":             {"synthetic-api-key"},
		"Cookie":              {"synthetic-cookie"},
		"Openai-Project":      {"synthetic-project"},
	}
	response := &http.Response{StatusCode: http.StatusUnauthorized, Request: request, Body: http.NoBody}
	redacted := redactX509Response(response)
	defer func() { _ = redacted.Body.Close() }()
	if redacted == response || redacted.Request == request {
		t.Fatal("response metadata was not independently cloned")
	}
	if len(redacted.Request.Header) != 1 || redacted.Request.Header.Get("OpenAI-Project") != "synthetic-project" {
		t.Errorf("redacted response headers = %v, want only safe project metadata", redacted.Request.Header)
	}
	if request.Header.Get("Authorization") != "Bearer synthetic-token" || len(request.Header) != 6 {
		t.Error("redacting returned metadata changed the original wire request")
	}
	if redacted.StatusCode != response.StatusCode || redacted.Request.URL.String() != request.URL.String() {
		t.Error("response status or request URL changed during metadata redaction")
	}
}
