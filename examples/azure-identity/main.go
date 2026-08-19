package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client := openai.NewClient(
		azure.WithEndpoint("https://example-endpoint.openai.azure.com", "2025-03-01-preview"),
		azure.WithTokenCredential(credential),
	)

	response, err := client.Responses.New(context.Background(), responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("Write me a haiku about computers"),
		},
		Model: openai.ChatModel("model-name"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(response.OutputText())
}
