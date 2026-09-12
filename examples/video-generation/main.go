package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
)

func main() {
	client := openai.NewClient()

	ctx := context.Background()

	//nolint:staticcheck // SA1019: the SDK team keeps this example until the Sora API shuts down on 2026-09-24.
	video, err := client.Videos.NewAndPoll(ctx, openai.VideoNewParams{
		Model:  openai.VideoModelSora2,
		Prompt: "A video of the words 'Thank you' in sparkling letters",
	}, 1000)
	if err != nil {
		panic(err)
	}

	if video.Status == openai.VideoStatusCompleted {
		fmt.Println("Video successfully completed: ", video)
	} else {
		fmt.Printf("Video creation failed. Status: %s\n", video.Status)
	}
}
