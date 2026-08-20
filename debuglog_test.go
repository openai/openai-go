package openai_test

import (
	"bytes"
	"context"
	"io"
	"log"
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
		proxyAuthorization []string
		rawQuery           string
		err                error
	}
	received := make(chan receivedRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, readErr := io.ReadAll(r.Body)
		received <- receivedRequest{
			body:               string(body),
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
	baseURL.User = url.UserPassword("url-user", "url-password")

	var output bytes.Buffer
	var response []byte
	var rawResponse *http.Response
	client := openai.NewClient(
		option.WithBaseURL(baseURL.String()),
		option.WithAPIKey("api-key-secret"),
		option.WithHeaderAdd("Proxy-Authorization", "first-proxy-secret"),
		option.WithHeaderAdd("Proxy-Authorization", "second-proxy-secret"),
		option.WithHeader("X-Custom-Secret", "custom-header-secret"),
		option.WithDebugLog(log.New(&output, "", 0)),
	)
	err = client.Post(
		context.Background(),
		"debug?api_key=query-secret&access_token=access-token-secret&sig=signature-secret",
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
		"url-user",
		"url-password",
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
	} {
		if strings.Contains(logOutput, secret) {
			t.Errorf("debug log contains sensitive value %q: %s", secret, logOutput)
		}
	}
	for _, metadata := range []string{
		`"method":"POST"`,
		`"url":"` + server.URL + `/debug"`,
		`"Authorization":["***"]`,
		`"Proxy-Authorization":["***","***"]`,
		`"status_code":200`,
		`"Set-Cookie":["***","***"]`,
		`"X-Request-Id":["request-id-value"]`,
	} {
		if !strings.Contains(logOutput, metadata) {
			t.Errorf("debug log missing %q: %s", metadata, logOutput)
		}
	}
}
