package bedrock_test

import (
	"context"

	"github.com/openai/openai-go/v3/bedrock"
)

func ExampleNewClient() {
	// The default AWS credential chain supports environment credentials,
	// ~/.aws/credentials, ~/.aws/config, SSO, web identity, and workload roles.
	client, err := bedrock.NewClient(context.Background(), bedrock.Config{
		AWSRegion: "us-west-2",
	})
	if err != nil {
		panic(err)
	}
	_ = client
}

func ExampleNewClient_profile() {
	client, err := bedrock.NewClient(context.Background(), bedrock.Config{
		AWSProfile: "production",
	})
	if err != nil {
		panic(err)
	}
	_ = client
}

func ExampleNewClient_bearer() {
	client, err := bedrock.NewClient(context.Background(), bedrock.Config{
		AWSRegion: "us-west-2",
		APIKey:    "bedrock-bearer-token",
	})
	if err != nil {
		panic(err)
	}
	_ = client
}
