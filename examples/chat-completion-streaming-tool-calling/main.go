package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/openai/openai-go/v3"
)

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	question := "What is the weather in New York City? Write a short paragraph about it."

	fmt.Print("> ")
	fmt.Println(question)

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},
		Tools: []openai.ChatCompletionToolUnionParam{
			openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name:        "get_weather",
				Description: openai.String("Get weather at the given location"),
				Strict:      openai.Bool(true),
				Parameters: openai.FunctionParameters{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]any{
						"location": map[string]string{
							"type": "string",
						},
					},
					"required": []string{"location"},
				},
			}),
		},
		ParallelToolCalls: openai.Bool(false), // JustFinishedToolCall is only safe with parallel calls disabled
		Seed:              openai.Int(0),
		Model:             openai.ChatModelGPT4o,
	}

	stream := client.Chat.Completions.NewStreaming(ctx, params)
	defer func() { _ = stream.Close() }()
	acc := openai.ChatCompletionAccumulator{}

	fmt.Println("\nStreaming first response...")
	for stream.Next() {
		chunk := stream.Current()
		if !acc.AddChunk(chunk) {
			panic("failed to accumulate streaming chunk")
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}

		if tool, ok := acc.JustFinishedToolCall(); ok {
			fmt.Printf("\nTool call detected: %s with arguments %s\n", tool.Name, tool.Arguments)
		}
	}
	fmt.Println()

	if err := stream.Err(); err != nil {
		panic(err)
	}

	if len(acc.Choices) == 0 || len(acc.Choices[0].Message.ToolCalls) == 0 {
		fmt.Printf("No function call\n")
		return
	}

	message := acc.Choices[0].Message
	toolCalls := message.ToolCalls
	params.Messages = append(params.Messages, message.ToParam())
	for _, toolCall := range toolCalls {
		if toolCall.Function.Name != "get_weather" {
			panic(fmt.Sprintf("unsupported tool call: %q", toolCall.Function.Name))
		}

		var args struct {
			Location string `json:"location"`
		}
		decoder := json.NewDecoder(strings.NewReader(toolCall.Function.Arguments))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&args); err != nil {
			panic(fmt.Sprintf("failed to decode tool arguments: %v", err))
		}
		if err := decoder.Decode(&struct{}{}); err != io.EOF {
			panic("failed to decode tool arguments: expected a single JSON object")
		}
		if args.Location == "" {
			fmt.Printf("Missing or invalid 'location' argument\n")
			return
		}
		weatherData := getWeather(args.Location)
		fmt.Printf("Weather in %s: %s\n", args.Location, weatherData)

		params.Messages = append(params.Messages, openai.ToolMessage(weatherData, toolCall.ID))
	}

	// Disable tools for the second round so the model returns a final answer
	// instead of making additional tool calls we don't handle. Use nil (not an
	// empty slice) so the field is omitted from the request.
	params.Tools = nil

	responseStream := client.Chat.Completions.NewStreaming(ctx, params)
	defer func() { _ = responseStream.Close() }()

	// The second response has a new completion ID, so it needs its own
	// accumulator; reusing `acc` would make AddChunk return false for every chunk.
	responseAcc := openai.ChatCompletionAccumulator{}

	fmt.Println("\nStreaming second response...")
	for responseStream.Next() {
		evt := responseStream.Current()
		if !responseAcc.AddChunk(evt) {
			panic("failed to accumulate second response chunk")
		}
		if len(evt.Choices) > 0 {
			fmt.Print(evt.Choices[0].Delta.Content)
		}
	}
	fmt.Println()

	if err := responseStream.Err(); err != nil {
		panic(err)
	}
}

func getWeather(location string) string {
	// In a real implementation, this function would call a weather API.
	return "Sunny, 25°C"
}
