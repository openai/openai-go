# Deferred function tools

Set `DeferLoading: openai.Bool(true)` on a function definition and include a
`tool_search` tool so the model can discover it when needed. Both APIs already
support this option in Go. Explicit `openai.Bool(false)` sends `false`; leaving
the field unset omits it. The API validates the tool-search configuration.

## Responses

```go
params := responses.ResponseNewParams{
    Model: MODEL,
    Input: responses.ResponseNewParamsInputUnion{
        OfString: openai.String("Look up item A123."),
    },
    Tools: []responses.ToolUnionParam{
        {OfToolSearch: &responses.ToolSearchToolParam{}},
        {OfFunction: &responses.FunctionToolParam{
            Name: "lookup_item", Parameters: schema,
            Strict: openai.Bool(true), DeferLoading: openai.Bool(true),
        }},
    },
}
response, err := client.Responses.New(ctx, params)
```

Use the same params with `client.Responses.NewStreaming`. Function arguments
remain JSON strings; decode them into your application struct as usual:

```go
var args struct {
    ItemID string `json:"item_id"`
}
err := json.Unmarshal([]byte(call.Arguments), &args)
```

Handle the call and submit a `function_call_output` with its `CallID`, just as
for a non-deferred function. Go does not attach an automatic argument parser to
the tool definition.

## Hosted Agents

```go
tool := openai.AgentToolParamFunction{
    Name: "lookup_item", Description: "Look up a catalog item",
    Parameters: schema, DeferLoading: openai.Bool(true),
}
tools := []openai.AgentToolParamUnion{
    {OfParamToolSearch: &openai.AgentToolParamToolSearch{}},
    {OfParamFunction: &tool},
}
// Supply tools in the agent configuration when creating the session.
// Once the session is idle, submit input with the existing handler helper:
stream := client.Beta.Agents.Sessions.Stream(ctx, session.ID, openai.AgentSessionStreamParams{
    Input: "Look up item A123.",
    ToolHandlers: map[string]openai.AgentToolHandler{
        tool.Name: func(ctx context.Context, args map[string]any) (any, error) {
            itemID, ok := args["item_id"].(string)
            if !ok {
                return nil, errors.New("item_id must be a string")
            }
            return catalog.Lookup(ctx, itemID)
        },
    },
})
defer stream.Close()
result, err := stream.FinalResult()
```

The existing handler receives decoded arguments and submits the tool result.
Deferred loading changes discovery, not the application's execution path.
