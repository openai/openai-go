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

	res, err := client.Audio.Speech.New(ctx, openai.AudioSpeechNewParams{
		Model:          openai.SpeechModelTTS1,
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
