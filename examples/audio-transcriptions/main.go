package main

import (
	"context"
	"io"
	"os"

	"github.com/openai/openai-go"
)

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	file, err := os.Open("speech.mp3")
	if err != nil {
		panic(err)
	}

	transcription, err := client.Audio.Transcriptions.New(ctx, openai.AudioTranscriptionNewParams{
		Model: openai.F(openai.AudioTranscriptionNewParamsModelWhisper1),
		File:  openai.F[io.Reader](file),
	})
	if err != nil {
		panic(err)
	}

	println(transcription.Text)
}
