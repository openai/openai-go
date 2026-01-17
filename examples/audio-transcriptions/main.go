package main

import (
	"bytes"
	"context"
	"os"

	"github.com/openai/openai-go/v3"
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

	// Or if you have speech bytes, you have to wrap reader:
	var speechBytes []byte // Assume this is filled with audio data

	speechReader := openai.File(
		bytes.NewReader(speechBytes), "speech.mp3", "audio/mp3",
	)

	transcription, err = client.Audio.Transcriptions.New(ctx, openai.AudioTranscriptionNewParams{
		Model: openai.AudioModelWhisper1,
		File:  speechReader,
	})
	if err != nil {
		panic(err)
	}

	println(transcription.Text)
}
