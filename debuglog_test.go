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

func TestWithDebugLogEmitsOnlyRedactedMetadata(t *testing.T) {
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
		t.Fatalf("Post() error = %v", err)
	}

	gotRequest := <-received
	if gotRequest.err != nil {
		t.Fatalf("read request body: %v", gotRequest.err)
	}
	if gotRequest.body != requestBody {
		t.Fatalf("request body = %q, want %q", gotRequest.body, requestBody)
	}
	if got, want := gotRequest.escapedPath, "/debug/path-token-secret"; got != want {
		t.Fatalf("request path = %q, want %q", got, want)
	}
	if got, want := gotRequest.host, baseURL.Host; got != want {
		t.Fatalf("request host = %q, want %q", got, want)
	}
	const wantRawQuery = "api_key=query-secret&access_token=access-token-secret&sig=signature-secret"
	if gotRequest.rawQuery != wantRawQuery {
		t.Fatalf("request query = %q, want %q", gotRequest.rawQuery, wantRawQuery)
	}
	wantProxyAuthorization := []string{"first-proxy-secret", "second-proxy-secret"}
	if !reflect.DeepEqual(gotRequest.proxyAuthorization, wantProxyAuthorization) {
		t.Fatalf("Proxy-Authorization = %#v, want %#v", gotRequest.proxyAuthorization, wantProxyAuthorization)
	}
	if got := string(response); got != responseBody {
		t.Fatalf("response body = %q, want %q", got, responseBody)
	}
	if got := rawResponse.Header.Values("Set-Cookie"); !reflect.DeepEqual(got, []string{"first=response-secret", "second=response-secret"}) {
		t.Fatalf("Set-Cookie = %#v", got)
	}

	logOutput := output.String()
	for _, secret := range []string{
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
			t.Errorf("debug log contains sensitive value %q: %s", secret, logOutput)
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
			t.Errorf("debug log missing %q: %s", metadata, logOutput)
		}
	}
	if strings.Contains(logOutput, `"url":`) {
		t.Errorf("debug log contains URL metadata: %s", logOutput)
	}
}

func TestWithDebugLogOmitsResponseStatusText(t *testing.T) {
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
		t.Fatalf("Get() error = %v", err)
	}
	if rawResponse.Status != responseStatus {
		t.Fatalf("response status = %q, want %q", rawResponse.Status, responseStatus)
	}

	logOutput := output.String()
	if strings.Contains(logOutput, responseStatus) {
		t.Fatalf("debug log contains response status text: %s", logOutput)
	}
	if !strings.Contains(logOutput, `"status_code":200`) {
		t.Fatalf("debug log missing numeric status code: %s", logOutput)
	}
}
