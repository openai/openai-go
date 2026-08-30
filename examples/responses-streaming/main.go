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
	defer func() { _ = stream.Close() }()

	var completeText string

	for stream.Next() {
		switch event := stream.Current().AsAny().(type) {
		case responses.ResponseTextDeltaEvent:
			fmt.Print(event.Delta)
		case responses.ResponseCompletedEvent:
			completeText = event.Response.OutputText()
			fmt.Println("\nResponse completed")
		case responses.ResponseFailedEvent:
			fmt.Printf("\nResponse failed: %s\n", event.Response.ID)
		case responses.ResponseErrorEvent:
			fmt.Printf(
				"\nStream error event: %s param=%s code=%s\n",
				event.Message,
				event.Param,
				event.Code,
			)
		}
	}

	if err := stream.Err(); err != nil {
		// Transport/API-level streaming errors can surface here instead of as typed
		// ResponseErrorEvent / ResponseFailedEvent values.
		panic(err)
	}

	_ = completeText
}
