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

	stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},
		Seed:  openai.Int(0),
		Model: openai.ChatModelGPT4o,
	})

	for stream.Next() {
		evt := stream.Current()
		if len(evt.Choices) > 0 {
			print(evt.Choices[0].Delta.Content)
		}
	}
	println()

	if err := stream.Err(); err != nil {
		panic(err.Error())
	}
}
