package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
)

func main() {
	assistantID := os.Getenv("OPENAI_ASSISTANT_ID")
	if assistantID == "" {
		panic("OPENAI_ASSISTANT_ID must be set")
	}

	client := openai.NewClient()
	ctx := context.Background()

	question := "Write a haiku about streaming APIs."
	fmt.Printf("> %s\n", question)

	stream := client.Beta.Threads.NewAndRunStreaming(ctx, openai.BetaThreadNewAndRunParams{
		AssistantID: assistantID,
		Thread: openai.BetaThreadNewAndRunParamsThread{
			Messages: []openai.BetaThreadNewAndRunParamsThreadMessage{{
				Role: "user",
				Content: openai.BetaThreadNewAndRunParamsThreadMessageContentUnion{
					OfString: openai.String(question),
				},
			}},
		},
	})

	var finalText string

	for stream.Next() {
		event := stream.Current()

		switch event := event.AsAny().(type) {
		case openai.AssistantStreamEventThreadMessageDelta:
			for _, content := range event.Data.Delta.Content {
				if textDelta, ok := content.AsAny().(openai.TextDeltaBlock); ok {
					fmt.Print(textDelta.Text.Value)
				}
			}
		case openai.AssistantStreamEventThreadMessageCompleted:
			finalText = collectText(event.Data.Content)
		}
	}

	if err := stream.Err(); err != nil {
		panic(err)
	}

	if finalText != "" {
		fmt.Printf("\n\nFinal reply:\n%s\n", finalText)
	}
}

func collectText(content []openai.MessageContentUnion) string {
	var text string

	for _, part := range content {
		if textBlock, ok := part.AsAny().(openai.TextContentBlock); ok {
			text += textBlock.Text.Value
		}
	}

	return text
}
