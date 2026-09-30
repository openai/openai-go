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

type balanceArguments struct {
	Asset string `json:"asset" jsonschema:"enum=USDC,enum=ETH"`
}

type balanceReceipt struct {
	Wallet  string `json:"wallet"`
	Asset   string `json:"asset"`
	Balance string `json:"balance"`
}

// The wallet (and any credentials) stays in the application, outside model arguments.
type walletService interface {
	Balance(context.Context, balanceArguments) (balanceReceipt, error)
}

// This local fixture illustrates a Coinbase-style bound action, not an AgentKit
// integration. It performs no network requests or financial transactions.
type demoWallet struct{ id string }

func (w demoWallet) Balance(_ context.Context, args balanceArguments) (balanceReceipt, error) {
	return balanceReceipt{Wallet: w.id, Asset: args.Asset, Balance: "0.00"}, nil
}

type balanceAction struct{ wallet walletService }

func (a balanceAction) tool() (openai.AgentToolParamUnion, error) {
	// Choose a schema library in your application, or supply a JSON schema directly.
	reflector := jsonschema.Reflector{DoNotReference: true, AllowAdditionalProperties: false}
	data, err := json.Marshal(reflector.Reflect(balanceArguments{}))
	if err != nil {
		return openai.AgentToolParamUnion{}, err
	}
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		return openai.AgentToolParamUnion{}, err
	}
	return openai.AgentToolParamOfParamFunction("Read the application's wallet balance.", "wallet_balance", schema), nil
}

func (a balanceAction) handle(ctx context.Context, arguments map[string]any) (any, error) {
	if len(arguments) != 1 || arguments["asset"] == nil {
		return nil, errors.New("expected an asset argument and no other fields")
	}
	data, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args balanceArguments
	if err := json.Unmarshal(data, &args); err != nil {
		return nil, err
	}
	// Decoding into a struct does not validate JSON schema constraints.
	if args.Asset != "USDC" && args.Asset != "ETH" {
		return nil, errors.New("asset must be USDC or ETH")
	}
	receipt, err := a.wallet.Balance(ctx, args)
	if err != nil {
		return nil, err
	}
	// The existing dispatcher accepts JSON text, maps, and input content, not structs.
	output, err := json.Marshal(receipt)
	return string(output), err
}

func main() {
	action := balanceAction{wallet: demoWallet{id: "demo-wallet"}}
	tool, err := action.tool()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not build the tool definition.")
		os.Exit(1)
	}
	sessionID := os.Getenv("AGENT_SESSION_ID")
	if sessionID == "" {
		// Configure agent.tools with this definition when creating your session.
		data, err := json.MarshalIndent(tool, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not encode the tool definition.")
			os.Exit(1)
		}
		fmt.Println(string(data))
		fmt.Fprintln(os.Stderr, "Set AGENT_SESSION_ID to an idle session configured with this tool.")
		return
	}

	client := openai.NewClient()
	// This must be the session's only input writer. Keep iterating to dispatch calls.
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), sessionID, openai.AgentSessionStreamParams{
		Input: "What is my USDC balance?",
		ToolHandlers: map[string]openai.AgentToolHandler{
			tool.OfParamFunction.Name: action.handle,
		},
	})
	defer func() { _ = stream.Close() }()
	for stream.Next() {
		if event := stream.Current(); event.Type == "agent.session.turn.output_text.delta" {
			fmt.Print(event.AsAgentSessionTurnOutputTextDelta().Delta)
		}
	}
	if stream.Err() != nil {
		// Apply your application's error-reporting policy before logging raw errors.
		fmt.Fprintln(os.Stderr, "The agent stream failed.")
	}
}
