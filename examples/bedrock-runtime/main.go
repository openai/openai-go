package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/bedrock"
)

func main() {
	ctx := context.Background()
	config := bedrock.Config{Endpoint: bedrock.EndpointRuntime, AWSRegion: os.Getenv("AWS_REGION")}
	switch authentication := os.Getenv("BEDROCK_AUTH"); authentication {
	case "", "sigv4":
		// The explicit SigV4 example must not silently select an ambient bearer.
		if err := os.Unsetenv("AWS_BEARER_TOKEN_BEDROCK"); err != nil {
			panic(err)
		}
	case "bearer":
		config.APIKey = os.Getenv("AWS_BEARER_TOKEN_BEDROCK")
		if config.APIKey == "" {
			panic("bearer authentication requires AWS_BEARER_TOKEN_BEDROCK")
		}
	default:
		panic("BEDROCK_AUTH must be sigv4 or bearer")
	}
	client, err := bedrock.NewClient(ctx, config)
	if err != nil {
		panic(err)
	}

	model := os.Getenv("BEDROCK_MODEL")
	if model == "" {
		model = "us.openai.gpt-5.6-sol"
	}
	params := openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(model),
		Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage("Say hello from Amazon Bedrock Runtime!")},
	}
	if os.Getenv("BEDROCK_STREAM") == "1" {
		stream := client.Chat.Completions.NewStreaming(ctx, params)
		for stream.Next() {
			if choices := stream.Current().Choices; len(choices) != 0 {
				fmt.Print(choices[0].Delta.Content)
			}
		}
		if streamErr := stream.Err(); streamErr != nil {
			panic(streamErr)
		}
		fmt.Println()
		return
	}

	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}
	fmt.Println(completion.Choices[0].Message.Content)
}
