# Typed beta Agents output

Use one adapter for the schema and parser. `schema` can come from the same JSON
schema library used with Responses; use an object root.

```go
type Report struct {
    Summary string `json:"summary"`
}
output, err := openai.NewBetaAgentOutput(schema, func(data []byte) (Report, error) {
    var report Report
    err := json.Unmarshal(data, &report)
    return report, err // Add application validation here if needed.
})
if err != nil { return err }

params.Agent.Text.Format = output.Format()
stream := client.Beta.Agents.Sessions.NewStreaming(ctx, params)
defer stream.Close()
result, err := output.FinalResult(stream)
if err != nil { return err }
fmt.Println(result.OutputParsed.Summary)
```

The same `output.FinalResult(stream)` works with `Sessions.Stream` for follow-ups;
it parses the answer without changing that session's schema. Enable result
collection before iterating progress events. Parsing errors expose the completed
raw result through `BetaAgentOutputParseError.Result`.
