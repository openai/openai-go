package azure

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"

	"github.com/Nordlys-Labs/openai-go/v3"
	"github.com/Nordlys-Labs/openai-go/v3/internal/apijson"
	"github.com/Nordlys-Labs/openai-go/v3/internal/requestconfig"
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
	req, _ := http.NewRequest("POST", "/chat/completions", bytes.NewReader([]byte(`{"model":"my-model/v1"}`)))
	got, _ := getReplacementPathWithDeployment(req)

	expected := "/openai/deployments/my-model%2Fv1/chat/completions"
	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
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
