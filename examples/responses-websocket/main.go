package main

import (
	"context"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	conn, err := client.Responses.ConnectWebSocket(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			panic(closeErr)
		}
	}()

	stream, err := conn.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("Tell me a short joke")},
		Model: openai.ChatModelGPT4oMini,
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		if closeErr := stream.Close(); closeErr != nil {
			panic(closeErr)
		}
	}()

	for stream.Next() {
		event := stream.Current()
		if event.Delta != "" {
			print(event.Delta)
		}
	}
	if err := stream.Err(); err != nil {
		panic(err)
	}
}
