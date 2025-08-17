package azure_test

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/azure"
)

func Example_authentication() {
	// There are two ways to authenticate - using a TokenCredential (via the azidentity
	// package), or using an API Key.
	const azureOpenAIEndpoint = "https://<your-azureopenai-instance>.openai.azure.com"
	const azureOpenAIAPIVersion = "<api version string>"

	// Using a TokenCredential
	{
		// For a full list of credential types look at the documentation for the Azure Identity
		// package: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
		tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)

		if err != nil {
			fmt.Printf("Failed to create TokenCredential: %s\n", err)
			return
		}

		client := openai.NewClient(
			azure.WithEndpoint(azureOpenAIEndpoint, azureOpenAIAPIVersion),
			azure.WithTokenCredential(tokenCredential),
		)

		_ = client
	}

	// Using an API Key
	{
		const azureOpenAIAPIKey = "<key from Azure portal>"

		client := openai.NewClient(
			azure.WithEndpoint(azureOpenAIEndpoint, azureOpenAIAPIVersion),
			azure.WithAPIKey(azureOpenAIAPIKey),
		)

		_ = client
	}
}
