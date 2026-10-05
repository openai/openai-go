# Observe local tool failures

Use `OnToolError` to send local argument, handler, or output failures to your
logger or monitoring service:

```go
stream := client.Beta.Agents.Sessions.Stream(ctx, session.ID, openai.AgentSessionStreamParams{
	Input:        "Look up order A123.",
	ToolHandlers: handlers,
	OnToolError: func(ctx context.Context, failure openai.BetaAgentToolError) {
		// Choose what to log; failure.Err preserves the original error and can
		// contain sensitive application details.
		slog.ErrorContext(ctx, "Tool failed", "tool", failure.ToolName,
			"stage", failure.Stage, "call_id", failure.CallID)
	},
})
defer stream.Close()
result, err := stream.FinalResult()
```

The beta callback runs during stream iteration, once per local failure. The API
still receives only `Tool handler failed.`; request errors propagate through the
stream normally. Nothing is logged automatically. Callback panics propagate;
the helper does not recover panics or retry application code.
