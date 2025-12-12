package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	apiKey := os.Getenv("AZURE_OPENAI_API_KEY")
	apiVersion := "2025-03-01-preview"
	endpoint := "https://example-endpoint.openai.azure.com"
	deploymentName := "model-name" // e.g. "gpt-4o"

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, apiVersion),
		azure.WithAPIKey(apiKey),
	)

	ctx := context.Background()

	question := "Write me a haiku about computers"

	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(question)},
		Model: openai.ChatModel(deploymentName),
	})

	if err != nil {
		panic(err)
	}

	println(resp.OutputText())
}
