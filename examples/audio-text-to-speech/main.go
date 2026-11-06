//go:build !linux

package main

import (
	"context"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/openai/openai-go/examples/internal/errutil"
	"github.com/openai/openai-go/v3"
)

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	// This example demonstrates the legacy Speech endpoint. Its TTS models shut
	// down on January 6, 2027; new integrations should use the Realtime API.
	// See https://developers.openai.com/api/docs/deprecations for migration details.
	res, err := client.Audio.Speech.New(ctx, openai.AudioSpeechNewParams{
		Model:          openai.SpeechModel("tts-1"),
		Input:          `Why did the chicken cross the road? To get to the other side.`,
		ResponseFormat: openai.AudioSpeechNewParamsResponseFormatPCM,
		Voice: openai.AudioSpeechNewParamsVoiceUnion{
			OfAudioSpeechNewsVoiceString2: openai.String("alloy"),
		},
	})
	if err != nil {
		panic(errutil.Message(err))
	}
	defer func() { _ = res.Body.Close() }()

	op := &oto.NewContextOptions{}
	op.SampleRate = 24000
	op.ChannelCount = 1
	op.Format = oto.FormatSignedInt16LE

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic(errutil.Message(err))
	}

	<-readyChan

	player := otoCtx.NewPlayer(res.Body)
	player.Play()
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
}
