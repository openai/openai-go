package azure

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal/apijson"
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

func TestNoRouteChangeNeeded(t *testing.T) {
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

	req, err := http.NewRequest("POST", "/openai/does/not/need/a/deployment", bytes.NewReader(serializedBytes))

	if err != nil {
		t.Fatal(err)
	}

	replacementPath, err := getReplacementPathWithDeployment(req)

	if err != nil {
		t.Fatal(err)
	}

	if replacementPath != "/openai/does/not/need/a/deployment" {
		t.Fatalf("replacementpath didn't match: %s", replacementPath)
	}
}

func TestAPIKeyAuthentication(t *testing.T) {
	// Test that the API key option is created successfully
	apiKeyOption := WithAPIKey("test-api-key")

	// Verify the option is not nil
	if apiKeyOption == nil {
		t.Fatal("Expected API key option to be created")
	}

	// This test verifies the option is created correctly.
	// The actual header setting happens in the middleware chain.
	t.Log("API key option created successfully")
}
