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

Reattachment provides at-least-once tool-call delivery with application-owned
recovery. If a worker exits, the application restarts it, recovers the saved
session ID, and reattaches handlers. An unacknowledged call may be delivered again;
applications own idempotency when handlers perform mutations. Closing observation
leaves hosted execution running.
