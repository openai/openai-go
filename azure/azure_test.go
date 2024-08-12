package azure

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/shared"
)

func TestJSONRoute(t *testing.T) {
	chatCompletionParams := openai.ChatCompletionNewParams{
		Model: openai.F(openai.ChatModel("arbitraryDeployment")),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.ChatCompletionAssistantMessageParam{
				Role:    openai.F(openai.ChatCompletionAssistantMessageParamRoleAssistant),
				Content: openai.F[openai.ChatCompletionAssistantMessageParamContentUnion](shared.UnionString("You are a helpful assistant")),
			},
			openai.ChatCompletionUserMessageParam{
				Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
				Content: openai.F[openai.ChatCompletionUserMessageParamContentUnion](shared.UnionString("Can you tell me another word for the universe?")),
			},
		}),
	}

	serializedBytes, err := apijson.MarshalRoot(chatCompletionParams)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/openai/chat/completions", bytes.NewReader(serializedBytes))

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

	req, err := http.NewRequest("POST", "/openai/audio/transcriptions", bytes.NewReader(buff.Bytes()))

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
		Model: openai.F(openai.ChatModel("arbitraryDeployment")),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.ChatCompletionAssistantMessageParam{
				Role:    openai.F(openai.ChatCompletionAssistantMessageParamRoleAssistant),
				Content: openai.F[openai.ChatCompletionAssistantMessageParamContentUnion](shared.UnionString("You are a helpful assistant")),
			},
			openai.ChatCompletionUserMessageParam{
				Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
				Content: openai.F[openai.ChatCompletionUserMessageParamContentUnion](shared.UnionString("Can you tell me another word for the universe?")),
			},
		}),
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
