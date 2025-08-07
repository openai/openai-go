package main

import (
	"context"
	"os"

	"github.com/openai/openai-go/v2"
)

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	file, err := os.Open("speech.mp3")
	if err != nil {
		panic(err)
	}

	transcription, err := client.Audio.Transcriptions.New(ctx, openai.AudioTranscriptionNewParams{
		Model: openai.AudioModelWhisper1,
		File:  file,
	})
	if err != nil {
		panic(err)
	}

	println(transcription.Text)
}
