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

	stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		}),
		Seed:  openai.Int(0),
		Model: openai.F(openai.ChatModelGPT4o),
	})

	for stream.Next() {
		evt := stream.Current()
		print(evt.Choices[0].Delta.Content)
	}
	println()

	if err := stream.Err(); err != nil {
		panic(err)
	}
}
