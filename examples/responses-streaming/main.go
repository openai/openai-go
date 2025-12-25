package main

import (
	"context"

	"github.com/Nordlys-Labs/openai-go/v3"
	"github.com/Nordlys-Labs/openai-go/v3/responses"
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
