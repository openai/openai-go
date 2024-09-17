package main

import (
	"context"

	"github.com/openai/openai-go"
)

func main() {
	ctx := context.Background()
	client := openai.NewClient()
	assistant, err := client.Beta.Assistants.New(ctx, openai.BetaAssistantNewParams{
		Model:        openai.F(openai.ChatModelGPT4_1106Preview),
		Name:         openai.String("Math tutor"),
		Instructions: openai.String("You are a personal math tutor. Write and run code to answer math questions."),
	})

	if err != nil {
		panic(err.Error())
	}

	assistantID := assistant.ID
	println("Created and assistant with id", assistantID)

	ths := openai.NewBetaThreadService()

	thread, err := ths.New(ctx, openai.BetaThreadNewParams{
		Messages: openai.F("I need to solve the equation `3x + 11 = 34`. Can you help me?"),
	})
}
