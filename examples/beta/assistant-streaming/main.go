package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
)

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	question := "Write me a haiku"

	fmt.Println("> " + question)
	fmt.Println()

	stream := client.Beta.Threads.NewAndRunStreaming(ctx, openai.BetaThreadNewAndRunParams{
		AssistantID: "asst_123",
		Thread: openai.BetaThreadNewAndRunParamsThread{
			Messages: []openai.BetaThreadNewAndRunParamsThreadMessage{
				{
					Role: "user",
					Content: openai.BetaThreadNewAndRunParamsThreadMessageContentUnion{
						OfString: openai.String(question),
					},
				},
			},
		},
	})

	for stream.Next() {
		event := stream.Current()

		switch data := event.AsAny().(type) {
		case openai.AssistantStreamEventThreadMessageDelta:
			for _, content := range data.Data.Delta.Content {
				switch block := content.AsAny().(type) {
				case openai.TextDeltaBlock:
					fmt.Print(block.Text.Value)
				case openai.RefusalDeltaBlock:
					fmt.Print(block.Refusal)
				}
			}
		}
	}

	fmt.Println()

	if err := stream.Err(); err != nil {
		panic(err.Error())
	}
}
