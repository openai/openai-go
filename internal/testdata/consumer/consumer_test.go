package consumer_test

import (
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
)

func TestCoreAndAzurePackagesCompose(t *testing.T) {
	client := openai.NewClient(
		azure.WithEndpoint("https://example.openai.azure.com", "2024-06-01"),
		azure.WithAPIKey("test-key"),
	)

	if len(client.Options) == 0 {
		t.Fatal("NewClient did not retain its request options")
	}
}
