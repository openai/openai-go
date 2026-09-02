package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

// A struct that will be converted to a Structured Outputs response schema.
//
// NOTE: Avoid using `omitempty` in JSON tags for structured output schemas. The
// jsonschema library interprets `omitempty` as "optional", excluding the field
// from the schema's "required" array. This can cause the API to reject the
// schema (with Strict: true) or the model to silently skip populating those fields.
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

func GenerateSchema[T any]() (map[string]any, error) {
	// Structured Outputs requires object schemas to disallow additional
	// properties. Keep definitions referenced so recursive Go types terminate.
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
	}
	var v T
	schema := reflector.Reflect(v)
	if schema.Ref != "" {
		const definitionsPrefix = "#/$defs/"
		if len(schema.Ref) <= len(definitionsPrefix) || schema.Ref[:len(definitionsPrefix)] != definitionsPrefix {
			return nil, fmt.Errorf("expand root JSON schema reference %q: unsupported reference", schema.Ref)
		}
		definition, ok := schema.Definitions[schema.Ref[len(definitionsPrefix):]]
		if !ok {
			return nil, fmt.Errorf("expand root JSON schema reference %q: definition not found", schema.Ref)
		}
		expanded := *definition
		expanded.Version = schema.Version
		expanded.ID = schema.ID
		expanded.Anchor = schema.Anchor
		expanded.Definitions = schema.Definitions
		schema = &expanded
	}
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("marshal JSON schema: %w", err)
	}
	// Preserve schema values as raw JSON so integer constraints are not rounded
	// through float64 before the SDK serializes the request.
	var rawSchema map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawSchema); err != nil {
		return nil, fmt.Errorf("decode JSON schema: %w", err)
	}
	result := make(map[string]any, len(rawSchema))
	for key, value := range rawSchema {
		result[key] = value
	}
	return result, nil
}

func historicalComputerParams(question string, schema map[string]any) responses.ResponseNewParams {
	return responses.ResponseNewParams{
		Model: openai.ChatModelGPT5_2,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(question),
		},
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
					Name:        "historical_computer",
					Description: openai.String("Notable information about a computer"),
					Schema:      schema,
					Strict:      openai.Bool(true),
				},
			},
		},
	}
}

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	question := "What computer ran the first neural network?"

	print("> ")
	println(question)

	schema, err := GenerateSchema[HistoricalComputer]()
	if err != nil {
		panic(err)
	}

	response, err := client.Responses.New(ctx, historicalComputerParams(question, schema))
	if err != nil {
		panic(err)
	}

	// The model responds with a JSON string, so parse it into a struct
	var historicalComputer HistoricalComputer
	err = json.Unmarshal([]byte(response.OutputText()), &historicalComputer)
	if err != nil {
		panic(err)
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
