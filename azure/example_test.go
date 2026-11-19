package azure_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
)

func Example_authentication() {
	const azureOpenAIEndpoint = "https://<your-azureopenai-instance>.openai.azure.com"
	const azureOpenAIAPIVersion = "<api version string>"
	const azureOpenAIAPIKey = "<key from Azure portal>"

	// Supply a credential from your application here, such as the one returned
	// by azidentity.NewDefaultAzureCredential(nil).
	var tokenCredential azcore.TokenCredential

	tokenClient := openai.NewClient(
		azure.WithEndpoint(azureOpenAIEndpoint, azureOpenAIAPIVersion),
		azure.WithTokenCredential(tokenCredential),
	)

	client := openai.NewClient(
		azure.WithEndpoint(azureOpenAIEndpoint, azureOpenAIAPIVersion),
		azure.WithAPIKey(azureOpenAIAPIKey),
	)

	_ = tokenClient
	_ = client
}

func Example_authentication_custom_scopes() {
	const azureOpenAIEndpoint = "https://<your-azureopenai-instance>.openai.azure.com"
	const azureOpenAIAPIVersion = "<api version string>"

	// Supply a credential from your application here, such as the one returned
	// by azidentity.NewDefaultAzureCredential(nil).
	var tokenCredential azcore.TokenCredential

	client := openai.NewClient(
		azure.WithEndpoint(azureOpenAIEndpoint, azureOpenAIAPIVersion),
		azure.WithTokenCredential(
			tokenCredential,
			azure.WithTokenCredentialScopes([]string{"your-custom-scope"}),
		),
	)

	_ = client
}
