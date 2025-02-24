package main

import (
	"context"

	"github.com/openai/openai-go"
)

func main() {
	client := openai.NewClient()

	ctx := context.Background()

	question := "Write me a haiku"

	print("> ")
	println(question)
	println()

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.NewChatCompletionMessageParamOfUser(question),
		},
		Seed:  openai.Int(1),
		Model: openai.ChatModelGPT4o,
	})
	if err != nil {
		panic(err)
	}

	println(completion.Choices[0].Message.Content)
}
