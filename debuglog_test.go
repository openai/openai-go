package openai_test

import (
	"bytes"
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func clearOpenAIEnvironment(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"OPENAI_API_KEY",
		"OPENAI_ADMIN_KEY",
		"OPENAI_BASE_URL",
		"OPENAI_CUSTOM_HEADERS",
		"OPENAI_ORG_ID",
		"OPENAI_PROJECT_ID",
		"OPENAI_WEBHOOK_SECRET",
	} {
		t.Setenv(key, "")
	}
}

func TestWithDebugLogEmitsOnlyRedactedMetadata(t *testing.T) {
	// Prove this public-path regression does not inherit supported ambient headers.
	t.Setenv("OPENAI_CUSTOM_HEADERS", "Proxy-Authorization: ambient-proxy-secret")
	clearOpenAIEnvironment(t)

	const requestBody = "request-body-secret"
	const responseBody = "response-body-secret"

	type receivedRequest struct {
		body               string
		escapedPath        string
		host               string
		proxyAuthorization []string
		rawQuery           string
		err                error
	}
	received := make(chan receivedRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, readErr := io.ReadAll(r.Body)
		received <- receivedRequest{
			body:               string(body),
			escapedPath:        r.URL.EscapedPath(),
			host:               r.Host,
			proxyAuthorization: r.Header.Values("Proxy-Authorization"),
			rawQuery:           r.URL.RawQuery,
			err:                readErr,
		}
		w.Header().Add("Set-Cookie", "first=response-secret")
		w.Header().Add("Set-Cookie", "second=response-secret")
		w.Header().Set("X-Request-ID", "request-id-value")
		_, _ = io.WriteString(w, responseBody)
	}))
	t.Cleanup(server.Close)

	baseURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse server URL: %v", err)
	}
	const credentialBearingHostname = "customer-credential-host-secret.example.test"
	baseURL.Host = net.JoinHostPort(credentialBearingHostname, baseURL.Port())
	baseURL.User = url.UserPassword("url-user", "url-password")

	dialer := &net.Dialer{}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network string, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, server.Listener.Addr().String())
		},
	}
	t.Cleanup(transport.CloseIdleConnections)

	var output bytes.Buffer
	var response []byte
	var rawResponse *http.Response
	client := openai.NewClient(
		option.WithBaseURL(baseURL.String()),
		option.WithAPIKey("api-key-secret"),
		option.WithHTTPClient(&http.Client{Transport: transport}),
		option.WithHeaderAdd("Proxy-Authorization", "first-proxy-secret"),
		option.WithHeaderAdd("Proxy-Authorization", "second-proxy-secret"),
		option.WithHeader("X-Custom-Secret", "custom-header-secret"),
		option.WithDebugLog(log.New(&output, "", 0)),
	)
	err = client.Post(
		context.Background(),
		"debug/path-token-secret?api_key=query-secret&access_token=access-token-secret&sig=signature-secret",
		[]byte(requestBody),
		&response,
		option.WithResponseInto(&rawResponse),
	)
	if err != nil {
		t.Fatal("Post() returned an error")
	}

	gotRequest := <-received
	if gotRequest.err != nil {
		t.Fatalf("read request body: %v", gotRequest.err)
	}
	if gotRequest.body != requestBody {
		t.Fatal("request body did not reach the server unchanged")
	}
	if got, want := gotRequest.escapedPath, "/debug/path-token-secret"; got != want {
		t.Fatal("request path did not reach the server unchanged")
	}
	if got, want := gotRequest.host, baseURL.Host; got != want {
		t.Fatal("request host did not reach the server unchanged")
	}
	const wantRawQuery = "api_key=query-secret&access_token=access-token-secret&sig=signature-secret"
	if gotRequest.rawQuery != wantRawQuery {
		t.Fatal("request query did not reach the server unchanged")
	}
	wantProxyAuthorization := []string{"first-proxy-secret", "second-proxy-secret"}
	if !reflect.DeepEqual(gotRequest.proxyAuthorization, wantProxyAuthorization) {
		t.Fatal("Proxy-Authorization headers did not match the isolated fixture")
	}
	if got := string(response); got != responseBody {
		t.Fatal("response body did not reach the caller unchanged")
	}
	if got := rawResponse.Header.Values("Set-Cookie"); !reflect.DeepEqual(got, []string{"first=response-secret", "second=response-secret"}) {
		t.Fatal("Set-Cookie headers did not reach the caller unchanged")
	}

	logOutput := output.String()
	for index, secret := range []string{
		credentialBearingHostname,
		"url-user",
		"url-password",
		"path-token-secret",
		"api_key",
		"query-secret",
		"access_token",
		"access-token-secret",
		"signature-secret",
		"api-key-secret",
		"first-proxy-secret",
		"second-proxy-secret",
		"custom-header-secret",
		requestBody,
		responseBody,
		"first=response-secret",
		"second=response-secret",
		"request-id-value",
	} {
		if strings.Contains(logOutput, secret) {
			t.Errorf("debug log contains synthetic sensitive marker %d", index)
		}
	}
	for _, metadata := range []string{
		`"method":"POST"`,
		`"Authorization":["***"]`,
		`"Proxy-Authorization":["***","***"]`,
		`"status_code":200`,
		`"Set-Cookie":["***","***"]`,
	} {
		if !strings.Contains(logOutput, metadata) {
			t.Errorf("debug log missing %q", metadata)
		}
	}
	if strings.Contains(logOutput, `"url":`) {
		t.Error("debug log contains URL metadata")
	}
}

func TestWithDebugLogOmitsResponseStatusText(t *testing.T) {
	clearOpenAIEnvironment(t)

	const responseStatus = "200 response-status-secret"

	var output bytes.Buffer
	client := openai.NewClient(
		option.WithAPIKey("api-key-secret"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Status:     responseStatus,
						StatusCode: http.StatusOK,
						Header:     make(http.Header),
						Body:       io.NopCloser(strings.NewReader("response-body-secret")),
						Request:    req,
					}, nil
				},
			},
		}),
		option.WithDebugLog(log.New(&output, "", 0)),
	)

	var response []byte
	var rawResponse *http.Response
	if err := client.Get(
		context.Background(),
		"debug",
		nil,
		&response,
		option.WithResponseInto(&rawResponse),
	); err != nil {
		t.Fatal("Get() returned an error")
	}
	if rawResponse.Status != responseStatus {
		t.Fatal("response status did not reach the caller unchanged")
	}

	logOutput := output.String()
	if strings.Contains(logOutput, responseStatus) {
		t.Fatal("debug log contains response status text")
	}
	if !strings.Contains(logOutput, `"status_code":200`) {
		t.Fatal("debug log missing numeric status code")
	}
}

func TestWithDebugLogRedactsUnrecognizedMethod(t *testing.T) {
	clearOpenAIEnvironment(t)

	const requestMethod = "CUSTOM-METHOD-SECRET"

	receivedMethod := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod <- r.Method
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)

	var output bytes.Buffer
	client := openai.NewClient(
		option.WithBaseURL(server.URL),
		option.WithAPIKey("api-key-secret"),
		option.WithDebugLog(log.New(&output, "", 0)),
	)
	if err := client.Execute(context.Background(), requestMethod, "debug", nil, nil); err != nil {
		t.Fatal("Execute() returned an error")
	}
	if got := <-receivedMethod; got != requestMethod {
		t.Fatal("request method did not reach the server unchanged")
	}

	logOutput := output.String()
	if strings.Contains(logOutput, requestMethod) {
		t.Fatal("debug log contains unrecognized request method")
	}
	if !strings.Contains(logOutput, `"method":"***"`) {
		t.Fatal("debug log missing method placeholder")
	}
}
