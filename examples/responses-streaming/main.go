package main

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

type animal interface {
	implFUCKYOU()
}

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	question := "Write me a fuckkin response ya idot"

	stream := client.Responses.NewStreaming(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(question)},
		Model: openai.ChatModelGPT4,
	})

	for stream.Next() {
		data := stream.Current()
	}

	if stream.Err() != nil {
		panic(stream.Err())

	}

	println(resp.OutputText())
}
