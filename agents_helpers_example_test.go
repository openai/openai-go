package openai_test

import (
	"context"
	"fmt"
	"time"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func ExampleBetaAgentSessionService_Stream() {
	client := openai.NewClient()
	// The existing session must be idle, and this must be its only input writer.
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session_id", openai.AgentSessionStreamParams{
		Input: "What is the weather in Paris?",
		ToolHandlers: map[string]openai.AgentToolHandler{
			"get_weather": func(ctx context.Context, args map[string]any) (any, error) {
				// Call your application's weather service with ctx and args here.
				return map[string]any{"temperature_celsius": 21}, nil
			},
		},
	}, option.WithRequestTimeout(2*time.Minute))
	defer func() { _ = stream.Close() }()
	for stream.Next() {
		event := stream.Current()
		if event.Type == "agent.session.turn.output_text.delta" {
			fmt.Print(event.AsAgentSessionTurnOutputTextDelta().Delta)
		}
	}
	if err := stream.Err(); err != nil {
		// Handle the error through your application's error-reporting policy.
		return
	}
}

func ExampleAgentSessionMessage_OutputText() {
	message := openai.AgentSessionMessage{Content: []openai.AgentSessionMessageContentUnion{
		{Type: "output_text", Text: "Hello"},
		{Type: "output_text", Text: " world!"},
	}}
	fmt.Println(message.OutputText())
	// Output: Hello world!
}

func ExampleAgentSessionMessage_OutputText_listed() {
	client := openai.NewClient()
	items := client.Beta.Agents.Sessions.Items.ListAutoPaging(context.Background(), "session_id", openai.BetaAgentSessionItemListParams{})
	for items.Next() {
		item := items.Current()
		if item.Type == "message" {
			fmt.Println(item.AsMessage().OutputText())
		}
	}
	if err := items.Err(); err != nil {
		// Handle the error through your application's error-reporting policy.
		return
	}
}
