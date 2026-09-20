package main

import (
	"context"
	"log"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// OPENAI_API_KEY supplies Authorization independently of custom headers.
	client := openai.NewClient(option.WithHeader("X-Application", "my-application"))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{},
		option.WithHeader("X-Request-Tag", "my-websocket-session"),
	)
	if err != nil {
		return err
	}
	defer func() { _ = connection.Close() }()
	// Custom headers are sent on the opening handshake. Reconnect preserves
	// these options and reruns any configured authentication middleware.
	first := responses.ResponsesClientEventResponseCreateParam{
		Model: "gpt-4o-mini",
		Input: responses.ResponsesClientEventResponseCreateInputUnionParam{OfString: openai.String("Say hello.")},
	}
	if err = connection.Create(ctx, first); err != nil {
		return err
	}
	response, err := connection.FinalResponse(ctx)
	if err != nil {
		return err
	}
	// Every turn is response.create. Send only new input with the previous ID.
	if err = connection.Create(ctx, responses.ResponsesClientEventResponseCreateParam{
		Model:              "gpt-4o-mini",
		PreviousResponseID: openai.String(response.ID),
		Input:              responses.ResponsesClientEventResponseCreateInputUnionParam{OfString: openai.String("Now say it in French.")},
	}); err != nil {
		return err
	}
	response, err = connection.FinalResponse(ctx)
	if err != nil {
		return err
	}
	log.Printf("Continuation finished with status %s", response.Status)
	return nil
}
