package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v2"
)

// A struct that will be converted to a Structured Outputs response schema
type HistoricalComputer struct {
	Origin       Origin   `json:"origin" jsonschema_description:"The origin of the computer"`
	Name         string   `json:"full_name" jsonschema_description:"The name of the device model"`
	Legacy       string   `json:"legacy" jsonschema:"enum=positive,enum=neutral,enum=negative" jsonschema_description:"Its influence on the field of computing"`
	NotableFacts []string `json:"notable_facts" jsonschema_description:"A few key facts about the computer"`
}

type Origin struct {
	YearBuilt    int64  `json:"year_of_construction" jsonschema_description:"The year it was made"`
	Organization string `json:"organization" jsonschema_description:"The organization that was in charge of its development"`
}

func GenerateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

// Generate the JSON schema at initialization time
var HistoricalComputerResponseSchema = GenerateSchema[HistoricalComputer]()

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	question := "What computer ran the first neural network?"

	print("> ")
	println(question)

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "historical_computer",
		Description: openai.String("Notable information about a computer"),
		Schema:      HistoricalComputerResponseSchema,
		Strict:      openai.Bool(true),
	}

	// Query the Chat Completions API
	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
		},
		// Only certain models can perform structured outputs
		Model: openai.ChatModelGPT4o2024_08_06,
	})

	if err != nil {
		panic(err.Error())
	}

	// The model responds with a JSON string, so parse it into a struct
	var historicalComputer HistoricalComputer
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &historicalComputer)
	if err != nil {
		panic(err.Error())
	}

	// Use the model's structured response with a native Go struct
	fmt.Printf("Name: %v\n", historicalComputer.Name)
	fmt.Printf("Year: %v\n", historicalComputer.Origin.YearBuilt)
	fmt.Printf("Org: %v\n", historicalComputer.Origin.Organization)
	fmt.Printf("Legacy: %v\n", historicalComputer.Legacy)
	fmt.Printf("Facts:\n")
	for i, fact := range historicalComputer.NotableFacts {
		fmt.Printf("%v. %v\n", i+1, fact)
	}
}
