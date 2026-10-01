# Reattach beta Agents tool handlers

Save the session ID and omit `Input` to continue handling pending function calls
without submitting the original prompt again.

```go
stream := client.Beta.Agents.Sessions.Stream(ctx, savedSessionID,
    openai.AgentSessionStreamParams{ToolHandlers: handlers})
defer stream.Close()
result, err := stream.FinalResult()
if err != nil { return err }
fmt.Println(result.OutputText())
```

For progress, call `WithResultCollection()` before iterating `Next()`, then call
`FinalResult()`. A completed result includes durable output missed during the
disconnect. Attaching to an idle session with no selected turn drains normally;
`FinalResult()` reports `no_turn_selected`.

Reattachment uses at-least-once tool-call delivery: an unacknowledged call may be
delivered again after reconnecting. Applications are responsible for idempotency
when handlers perform mutations. Closing the stream stops local observation, not
hosted execution.
