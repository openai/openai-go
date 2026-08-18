package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/bedrock"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	ctx := context.Background()
	client, err := bedrock.NewClient(ctx, bedrock.Config{
		AWSRegion: os.Getenv("AWS_REGION"),
	})
	if err != nil {
		panic(err)
	}

	model := os.Getenv("BEDROCK_MODEL")
	if model == "" {
		model = "openai.gpt-oss-120b"
	}
	response, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Model: openai.ChatModel(model),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("Write a haiku about cloud credentials."),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(response.OutputText())
}
