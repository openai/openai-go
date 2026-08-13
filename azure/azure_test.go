package azure

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

func TestJSONRoute(t *testing.T) {
	chatCompletionParams := openai.ChatCompletionNewParams{
		Model: openai.ChatModel("arbitraryDeployment"),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.AssistantMessage("You are a helpful assistant"),
			openai.UserMessage("Can you tell me another word for the universe?"),
		},
	}

	serializedBytes, err := apijson.MarshalRoot(chatCompletionParams)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/chat/completions", bytes.NewReader(serializedBytes))

	if err != nil {
		t.Fatal(err)
	}

	replacementPath, err := getReplacementPathWithDeployment(req)

	if err != nil {
		t.Fatal(err)
	}

	if replacementPath != "/openai/deployments/arbitraryDeployment/chat/completions" {
		t.Fatalf("replacementpath didn't match: %s", replacementPath)
	}
}

func TestGetAudioMultipartRoute(t *testing.T) {
	buff := &bytes.Buffer{}
	mw := multipart.NewWriter(buff)
	defer mw.Close()

	fw, err := mw.CreateFormFile("file", "test.mp3")

	if err != nil {
		t.Fatal(err)
	}

	if _, err = fw.Write([]byte("ignore me")); err != nil {
		t.Fatal(err)
	}

	if err := mw.WriteField("model", "arbitraryDeployment"); err != nil {
		t.Fatal(err)
	}

	if err := mw.Close(); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/audio/transcriptions", bytes.NewReader(buff.Bytes()))

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", mw.FormDataContentType())

	replacementPath, err := getReplacementPathWithDeployment(req)

	if err != nil {
		t.Fatal(err)
	}

	if replacementPath != "/openai/deployments/arbitraryDeployment/audio/transcriptions" {
		t.Fatalf("replacementpath didn't match: %s", replacementPath)
	}
}

func TestAPIKeyAuthentication(t *testing.T) {
	rc := &requestconfig.RequestConfig{
		Request: &http.Request{
			Header: make(http.Header),
			URL:    &url.URL{},
		},
	}

	WithAPIKey("my-api-key").Apply(rc)

	if got := rc.Request.Header.Get("Api-Key"); got != "my-api-key" {
		t.Errorf("Api-Key header: got %q, expected %q", got, "my-api-key")
	}
}

func TestAPIKeyAuthenticationSuppressesAutomaticAuthorization(t *testing.T) {
	tests := []struct {
		name        string
		apiKey      string
		adminAPIKey string
	}{
		{name: "OpenAI API key", apiKey: "normal-openai-key"},
		{name: "OpenAI admin API key", adminAPIKey: "normal-admin-key"},
		{name: "both OpenAI keys", apiKey: "normal-openai-key", adminAPIKey: "normal-admin-key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("OPENAI_API_KEY", tt.apiKey)
			t.Setenv("OPENAI_ADMIN_KEY", tt.adminAPIKey)

			var captured *http.Request
			client := openai.NewClient(
				WithEndpoint("https://my-resource.openai.azure.com", "2024-10-21"),
				WithAPIKey("azure-api-key"),
				option.WithHTTPClient(&http.Client{
					Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
						captured = req
						return &http.Response{
							StatusCode: http.StatusOK,
							Header: http.Header{
								"Content-Type": []string{"application/json"},
							},
							Body:    io.NopCloser(strings.NewReader(`{"ok":true}`)),
							Request: req,
						}, nil
					}),
				}),
			)

			var res map[string]any
			if err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res); err != nil {
				t.Fatalf("request failed: %s", err)
			}
			if captured == nil {
				t.Fatal("request was not captured")
			}
			if got := captured.Header.Get("Api-Key"); got != "azure-api-key" {
				t.Fatalf("Api-Key header = %q, want %q", got, "azure-api-key")
			}
			if got := captured.Header.Get("Authorization"); got != "" {
				t.Fatalf("Authorization header = %q, want empty", got)
			}
		})
	}
}

func TestJSONRoutePathConstruction(t *testing.T) {
	cases := []struct {
		path     string
		expected string
	}{
		{"/chat/completions", "/openai/deployments/gpt-4/chat/completions"},
		{"/completions", "/openai/deployments/gpt-4/completions"},
		{"/embeddings", "/openai/deployments/gpt-4/embeddings"},
		{"/audio/speech", "/openai/deployments/gpt-4/audio/speech"},
		{"/images/generations", "/openai/deployments/gpt-4/images/generations"},
		{"/models", "/openai/models"}, // endpoint without a deployment
		{"/files", "/openai/files"},   // endpoint without a deployment
	}
	for _, tc := range cases {
		req, _ := http.NewRequest("POST", tc.path, bytes.NewReader([]byte(`{"model":"gpt-4"}`)))
		got, _ := getReplacementPathWithDeployment(req)
		if got != tc.expected {
			t.Errorf("%s: got %q, expected %q", tc.path, got, tc.expected)
		}
	}
}

func TestModelWithSpecialCharsIsEscaped(t *testing.T) {
	tests := map[string]string{
		"slash":               "my-model/v1",
		"query and fragment":  "my-model?api-version=old#frag",
		"bare dot":            ".",
		"bare dot dot":        "..",
		"dot dot slash":       "../my-model",
		"encoded dot segment": "%2e%2e/my-model",
	}
	wantDeployment := map[string]string{
		"slash":               "my-model%2Fv1",
		"query and fragment":  "my-model%3Fapi-version=old%23frag",
		"bare dot":            "%2E",
		"bare dot dot":        "%2E%2E",
		"dot dot slash":       "..%2Fmy-model",
		"encoded dot segment": "%252e%252e%2Fmy-model",
	}

	for name, model := range tests {
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/chat/completions", bytes.NewReader([]byte(`{"model":`+strconv.Quote(model)+`}`)))
			got, err := getReplacementPathWithDeployment(req)
			if err != nil {
				t.Fatal(err)
			}

			expected := "/openai/deployments/" + wantDeployment[name] + "/chat/completions"
			if got != expected {
				t.Errorf("got %q, expected %q", got, expected)
			}
		})
	}
}

func TestMultipartModelWithSpecialCharsIsEscaped(t *testing.T) {
	tests := map[string]string{
		"slash":               "my-model/v1",
		"query and fragment":  "my-model?api-version=old#frag",
		"bare dot":            ".",
		"bare dot dot":        "..",
		"dot dot slash":       "../my-model",
		"encoded dot segment": "%2e%2e/my-model",
	}
	wantDeployment := map[string]string{
		"slash":               "my-model%2Fv1",
		"query and fragment":  "my-model%3Fapi-version=old%23frag",
		"bare dot":            "%2E",
		"bare dot dot":        "%2E%2E",
		"dot dot slash":       "..%2Fmy-model",
		"encoded dot segment": "%252e%252e%2Fmy-model",
	}

	for name, model := range tests {
		t.Run(name, func(t *testing.T) {
			req := newMultipartRouteRequest(t, "/audio/transcriptions", model)
			got, err := getReplacementPathWithDeployment(req)
			if err != nil {
				t.Fatal(err)
			}

			expected := "/openai/deployments/" + wantDeployment[name] + "/audio/transcriptions"
			if got != expected {
				t.Errorf("got %q, expected %q", got, expected)
			}
		})
	}
}

func TestWithEndpointPreservesEscapedPathParams(t *testing.T) {
	tests := map[string]string{
		"slash traversal":     "../videos/vid_123",
		"query and fragment":  "vs_123/files/file_456?limit=1#frag",
		"encoded dot segment": "%2e%2e/videos/vid_123",
	}
	wantPaths := map[string]string{
		"slash traversal":     "https://my-resource.openai.azure.com/openai/vector_stores/..%2Fvideos%2Fvid_123?api-version=2024-10-21",
		"query and fragment":  "https://my-resource.openai.azure.com/openai/vector_stores/vs_123%2Ffiles%2Ffile_456%3Flimit=1%23frag?api-version=2024-10-21",
		"encoded dot segment": "https://my-resource.openai.azure.com/openai/vector_stores/%252e%252e%2Fvideos%2Fvid_123?api-version=2024-10-21",
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			var captured *http.Request
			client := openai.NewClient(
				WithEndpoint("https://my-resource.openai.azure.com", "2024-10-21"),
				WithAPIKey("sk-test"),
				option.WithHTTPClient(&http.Client{
					Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
						captured = req
						return &http.Response{
							StatusCode: http.StatusOK,
							Header: http.Header{
								"Content-Type": []string{"application/json"},
							},
							Body:    io.NopCloser(strings.NewReader(azureVectorStoreResponse)),
							Request: req,
						}, nil
					}),
				}),
			)

			if _, err := client.VectorStores.Get(context.Background(), input); err != nil {
				t.Fatalf("request failed: %s", err)
			}
			if captured == nil {
				t.Fatal("request was not captured")
			}
			if got, want := captured.URL.String(), wantPaths[name]; got != want {
				t.Fatalf("url = %q, want %q", got, want)
			}
		})
	}
}

func TestWithEndpointBaseURL(t *testing.T) {
	tests := map[string]struct {
		endpoint        string
		apiVersion      string
		expectedBaseURL string
		expectedQuery   string
		shouldFail      bool
	}{
		"Azure endpoint": {
			endpoint:        "https://my-resource.openai.azure.com",
			apiVersion:      "2024-10-21",
			expectedBaseURL: "https://my-resource.openai.azure.com/",
			expectedQuery:   "api-version=2024-10-21",
		},
		"Azure endpoint with trailing slash": {
			endpoint:        "https://my-resource.openai.azure.com/",
			apiVersion:      "2024-10-21",
			expectedBaseURL: "https://my-resource.openai.azure.com/",
			expectedQuery:   "api-version=2024-10-21",
		},
		"Azure endpoint with path": {
			endpoint:        "https://my-resource.openai.azure.com/custom/path",
			apiVersion:      "2023-05-15",
			expectedBaseURL: "https://my-resource.openai.azure.com/custom/path/",
			expectedQuery:   "api-version=2023-05-15",
		},
		"empty apiVersion": {
			endpoint:   "https://my-resource.openai.azure.com",
			apiVersion: "",
			shouldFail: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opt := WithEndpoint(tc.endpoint, tc.apiVersion)

			rc := &requestconfig.RequestConfig{
				Request: &http.Request{
					Header: make(http.Header),
					URL:    &url.URL{},
				},
			}

			err := opt.Apply(rc)

			if tc.shouldFail {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("WithEndpoint returned error: %v", err)
			}

			if rc.BaseURL == nil {
				t.Fatal("BaseURL was not set")
			}
			if rc.BaseURL.String() != tc.expectedBaseURL {
				t.Errorf("BaseURL: got %q, expected %q", rc.BaseURL.String(), tc.expectedBaseURL)
			}

			query := rc.Request.URL.RawQuery
			if query != tc.expectedQuery {
				t.Errorf("Query: got %q, expected %q", query, tc.expectedQuery)
			}
		})
	}
}

func newMultipartRouteRequest(t *testing.T, route string, model string) *http.Request {
	t.Helper()

	buff := &bytes.Buffer{}
	mw := multipart.NewWriter(buff)
	fw, err := mw.CreateFormFile("file", "test.mp3")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = fw.Write([]byte("ignore me")); err != nil {
		t.Fatal(err)
	}
	if err := mw.WriteField("model", model); err != nil {
		t.Fatal(err)
	}
	if err := mw.Close(); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", route, bytes.NewReader(buff.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return req
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

const azureVectorStoreResponse = `{
	"id": "vs_dummy",
	"object": "vector_store",
	"created_at": 0,
	"file_counts": {
		"cancelled": 0,
		"completed": 0,
		"failed": 0,
		"in_progress": 0,
		"total": 0
	},
	"last_active_at": 0,
	"metadata": {},
	"name": "dummy",
	"status": "completed",
	"usage_bytes": 0
}`
