package main

import (
	"context"

	"github.com/Nordlys-Labs/openai-go/v3"
)

func main() {
	client := openai.NewClient()

	ctx := context.Background()

	question := "Write me a haiku"

	print("> ")
	println(question)
	println()
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},
		Seed:  openai.Int(0),
		Model: openai.ChatModelGPT4o,
	}

	completion, err := client.Chat.Completions.New(ctx, params)

	if err != nil {
		panic(err)
	}

	println(completion.Choices[0].Message.Content)
}
