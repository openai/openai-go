// This example adapts an application's typed action to the beta Agents API.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v3"
)

type lookupArguments struct {
	ItemID string `json:"item_id" jsonschema:"enum=item-1,enum=item-2"`
}

type itemRecord struct {
	Catalog      string `json:"catalog"`
	ItemID       string `json:"item_id"`
	Availability string `json:"availability"`
}

// The catalog (and any credentials) stays in the application, outside model arguments.
type catalogService interface {
	Lookup(context.Context, lookupArguments) (itemRecord, error)
}

// This read-only catalog fixture performs no network requests.
type demoCatalog struct{ id string }

func (c demoCatalog) Lookup(_ context.Context, args lookupArguments) (itemRecord, error) {
	return itemRecord{Catalog: c.id, ItemID: args.ItemID, Availability: "available"}, nil
}

type lookupAction struct{ catalog catalogService }

func (a lookupAction) tool() (openai.AgentToolParamUnion, error) {
	// Choose a schema library in your application, or supply a JSON schema directly.
	reflector := jsonschema.Reflector{DoNotReference: true, AllowAdditionalProperties: false}
	data, err := json.Marshal(reflector.Reflect(lookupArguments{}))
	if err != nil {
		return openai.AgentToolParamUnion{}, err
	}
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		return openai.AgentToolParamUnion{}, err
	}
	return openai.AgentToolParamOfParamFunction("Look up an item in the application's catalog.", "catalog_lookup", schema), nil
}

func (a lookupAction) handle(ctx context.Context, arguments map[string]any) (any, error) {
	if len(arguments) != 1 || arguments["item_id"] == nil {
		return nil, errors.New("expected an item_id argument and no other fields")
	}
	data, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args lookupArguments
	if err = json.Unmarshal(data, &args); err != nil {
		return nil, err
	}
	// Decoding into a struct does not validate JSON schema constraints.
	if args.ItemID != "item-1" && args.ItemID != "item-2" {
		return nil, errors.New("item_id must be item-1 or item-2")
	}
	record, err := a.catalog.Lookup(ctx, args)
	if err != nil {
		return nil, err
	}
	// The existing dispatcher accepts JSON text, maps, and input content, not structs.
	output, err := json.Marshal(record)
	return string(output), err
}

func main() {
	if err := run(context.Background(), openai.NewClient(), os.Getenv("AGENT_SESSION_ID")); err != nil {
		// Apply your application's error-reporting policy before logging raw errors.
		fmt.Fprintln(os.Stderr, "The typed-tools example failed.")
		os.Exit(1)
	}
}

func run(ctx context.Context, client openai.Client, sessionID string) error {
	action := lookupAction{catalog: demoCatalog{id: "demo-catalog"}}
	tool, err := action.tool()
	if err != nil {
		return err
	}
	if sessionID == "" {
		// Configure agent.tools with this definition when creating your session.
		data, err := json.MarshalIndent(tool, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		fmt.Fprintln(os.Stderr, "Set AGENT_SESSION_ID to an idle session configured with this tool.")
		return nil
	}

	// This must be the session's only input writer. Keep iterating to dispatch calls.
	stream := client.Beta.Agents.Sessions.Stream(ctx, sessionID, openai.AgentSessionStreamParams{
		Input: "Look up item-1 in the catalog.",
		ToolHandlers: map[string]openai.AgentToolHandler{
			tool.OfParamFunction.Name: action.handle,
		},
	})
	defer func() { _ = stream.Close() }()
	for stream.Next() {
		event := stream.Current()
		switch event.Type {
		case "agent.session.turn.output_text.delta":
			fmt.Print(event.AsAgentSessionTurnOutputTextDelta().Delta)
		case "agent.session.failed":
			return errors.New("agent session failed")
		case "agent.session.turn.failed", "agent.session.turn.cancelled":
			if event.Turn.SubagentID == "" {
				return errors.New("agent turn did not complete")
			}
		}
	}
	return stream.Err()
}
