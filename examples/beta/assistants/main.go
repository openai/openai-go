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

	println("Created and assistant with id", assistant.ID)

	prompt := "I need to solve the equation 3x + 11 = 14. Can you help me?"

	thread, err := client.Beta.Threads.New(ctx, openai.BetaThreadNewParams{
		Messages: openai.F([]openai.BetaThreadNewParamsMessage{
			{
				Content: openai.F([]openai.MessageContentPartParamUnion{
					openai.TextContentBlockParam{
						Text: openai.String(prompt),
						Type: openai.F(openai.TextContentBlockParamTypeText),
					},
				}),
				Role: openai.F(openai.BetaThreadNewParamsMessagesRoleUser),
			},
		}),
	})

	if err != nil {
		panic(err.Error())
	}

	println("Created thread with id", thread.ID)

	// pollIntervalMs of 0 uses default polling interval.
	run, err := client.Beta.Threads.Runs.NewAndPoll(ctx, thread.ID, openai.BetaThreadRunNewParams{
		AssistantID:            openai.F(assistant.ID),
		AdditionalInstructions: openai.String("Please address the user as Jane Doe. The user has a premium account."),
	}, 0)

	if err != nil {
		panic(err.Error())
	}

	if run.Status == openai.RunStatusCompleted {
		messages, err := client.Beta.Threads.Messages.List(ctx, thread.ID, openai.BetaThreadMessageListParams{})

		if err != nil {
			panic(err.Error())
		}

		for _, data := range messages.Data {
			for _, content := range data.Content {
				println(content.Text.Value)
			}
		}
	}
}
