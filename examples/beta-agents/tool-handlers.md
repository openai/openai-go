# Local tools on the first turn

Pass `ToolHandlers` to `NewStreaming` to run application code as the hosted agent
requests functions. The handlers stay local; define the tools in `Agent.Tools`.

```go
stream := client.Beta.Agents.Sessions.NewStreaming(ctx, openai.BetaAgentSessionNewParams{
    Agent: openai.BetaAgentSessionNewParamsAgent{
        Model: openai.String("gpt-6-astra"),
        Tools: []openai.AgentToolParamUnion{
            openai.AgentToolParamOfParamFunction("Look up an item", "lookup", map[string]any{
                "type": "object",
                "properties": map[string]any{"id": map[string]any{"type": "string"}},
                "required": []string{"id"}, "additionalProperties": false,
            }),
        },
    },
    Environment: openai.EnvironmentParamUnion{OfParamNone: &openai.EnvironmentParamNone{}},
    Input: openai.BetaAgentSessionNewParamsInputUnion{OfString: openai.String("Look up item A123.")},
    ToolHandlers: map[string]openai.AgentToolHandler{
        "lookup": func(ctx context.Context, args map[string]any) (any, error) {
            return inventory.Lookup(ctx, args["id"].(string))
        },
    },
})
defer stream.Close()
result, err := openai.BetaAgentSessionFinalResult(stream)
if err != nil {
    return err
}
fmt.Println(result.OutputText())
```

Handlers run sequentially as the stream advances. For progress events, call
`BetaAgentSessionWithResultCollection(stream)` before iterating; retrieve the
result afterward. A configured [typed output adapter](typed-output.md) also
accepts this stream through `output.FinalResult(stream)`.

Use `Sessions.Stream` with `AgentSessionStreamParams.ToolHandlers` for a new turn
on an existing idle session. Neither helper restarts work after a disconnect.
