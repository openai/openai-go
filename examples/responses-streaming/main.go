package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
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
		event := stream.Current()
		switch data := event.AsAny().(type) {
		case responses.ResponseTextDeltaEvent:
			print(data.Delta)
		case responses.ResponseTextDoneEvent:
			println()
			fmt.Println("Finished content")
			completeText = data.Text
		}
	}

	if stream.Err() != nil {
		panic(stream.Err())
	}

	_ = completeText
}
