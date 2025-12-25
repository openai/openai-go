package azure_test

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Nordlys-Labs/openai-go/v3"
	"github.com/Nordlys-Labs/openai-go/v3/azure"
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

func Example_authentication_custom_scopes() {
	// Custom scopes can also be passed, if needed, when using Azure OpenAI endpoints.
	const azureOpenAIEndpoint = "https://<your-azureopenai-instance>.openai.azure.com"
	const azureOpenAIAPIVersion = "<api version string>"

	// For a full list of credential types look at the documentation for the Azure Identity
	// package: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		fmt.Printf("Failed to create TokenCredential: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(azureOpenAIEndpoint, azureOpenAIAPIVersion),
		azure.WithTokenCredential(tokenCredential,
			// This is an example of a custom scope. See documentation for your service
			// endpoint for the proper value to pass.
			azure.WithTokenCredentialScopes([]string{"your-custom-scope"}),
		),
	)

	_ = client
}
