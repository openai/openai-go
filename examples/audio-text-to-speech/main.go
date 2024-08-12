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
		Model:          openai.F(openai.AudioModelWhisper1),
		Input:          openai.String(`Why did the chicken cross the road? To get to the other side.`),
		ResponseFormat: openai.F(openai.AudioSpeechNewParamsResponseFormatPCM),
		Voice:          openai.F(openai.AudioSpeechNewParamsVoiceAlloy),
	})
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}

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
