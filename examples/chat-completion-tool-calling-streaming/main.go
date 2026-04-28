package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
)

type weatherArgs struct {
	Location string `json:"location"`
}

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	question := "What should I pack for a trip to New York City this weekend?"

	fmt.Print("> ")
	fmt.Println(question)
	fmt.Println()

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},
		Tools: []openai.ChatCompletionToolUnionParam{
			openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name:        "get_weather",
				Description: openai.String("Get weather at the given location"),
				Parameters: openai.FunctionParameters{
					"type": "object",
					"properties": map[string]any{
						"location": map[string]string{
							"type": "string",
						},
					},
					"required": []string{"location"},
				},
			}),
		},
		ParallelToolCalls: openai.Bool(false),
		Seed:              openai.Int(0),
		Model:             openai.ChatModelGPT4o,
	}

	stream := client.Chat.Completions.NewStreaming(ctx, params)
	acc := openai.ChatCompletionAccumulator{}

	for stream.Next() {
		chunk := stream.Current()
		if !acc.AddChunk(chunk) {
			panic("failed to accumulate stream chunk")
		}

		if tool, ok := acc.JustFinishedToolCall(); ok {
			fmt.Println()
			fmt.Printf("tool call: %s %s\n", tool.Name, tool.Arguments)
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}
	if err := stream.Err(); err != nil {
		panic(err)
	}
	fmt.Println()

	if len(acc.Choices) == 0 || len(acc.Choices[0].Message.ToolCalls) == 0 {
		fmt.Println("No tool calls were returned.")
		return
	}

	params.Messages = append(params.Messages, acc.Choices[0].Message.ToParam())
	for _, toolCall := range acc.Choices[0].Message.ToolCalls {
		if toolCall.Function.Name != "get_weather" {
			continue
		}

		var args weatherArgs
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
			panic(err)
		}
		if args.Location == "" {
			args.Location = "New York City"
		}

		weatherData := getWeather(args.Location)
		fmt.Printf("Weather in %s: %s\n", args.Location, weatherData)

		params.Messages = append(params.Messages, openai.ToolMessage(weatherData, toolCall.ID))
	}

	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(completion.Choices[0].Message.Content)
}

func getWeather(location string) string {
	return "Cloudy, 18°C"
}
