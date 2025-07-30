package main

import (
	"context"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/openai/openai-go"
)

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	res, err := client.Audio.Speech.New(ctx, openai.AudioSpeechNewParams{
		Model:          openai.SpeechModelGPT4oMiniTTS,
		Input:          `Why did the chicken cross the road? To get to the other side.`,
		ResponseFormat: openai.AudioSpeechNewParamsResponseFormatPCM,
		Voice:          openai.AudioSpeechNewParamsVoiceAlloy,
	})

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	op := &oto.NewContextOptions{}
	op.SampleRate = 24000
	op.ChannelCount = 1
	op.Format = oto.FormatSignedInt16LE

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}

	<-readyChan

	player := otoCtx.NewPlayer(res.Body)
	player.Play()
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
	err = player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
}
