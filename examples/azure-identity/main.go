package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go/examples/internal/errutil"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	deployment := os.Getenv("AZURE_OPENAI_DEPLOYMENT")
	if endpoint == "" || deployment == "" {
		log.Fatal("Set AZURE_OPENAI_ENDPOINT and AZURE_OPENAI_DEPLOYMENT.")
	}

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal("Could not create the Azure credential.")
	}
	client := openai.NewClient(
		azure.WithEndpoint(endpoint, "2025-03-01-preview"),
		azure.WithTokenCredential(credential),
	)
	response, err := client.Responses.New(context.Background(), responses.ResponseNewParams{
		Model: openai.ChatModel(deployment),
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("Write a haiku about computers.")},
	})
	if err != nil {
		log.Fatal(errutil.Message(err))
	}
	fmt.Println(response.OutputText())
}
