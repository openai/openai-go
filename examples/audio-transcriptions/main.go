package main

import (
	"context"
	"os"

	"github.com/Nordlys-Labs/openai-go/v3"
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
