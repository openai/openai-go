package main

import (
	"context"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/responses"
)

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	question := "Tell me about briefly about Doug Engelbart"

	stream := client.Responses.NewStreaming(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(question)},
		Model: openai.ChatModelGPT4,
	})

	var completeText string

	for stream.Next() {
		data := stream.Current()
		print(data.Delta)
		if data.JSON.Text.Valid() {
			println()
			println("Finished Content")
			completeText = data.Text
			break
		}
	}

	if stream.Err() != nil {
		panic(stream.Err())
	}

	_ = completeText
}
