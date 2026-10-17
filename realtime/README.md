# Realtime and Translation WebSockets

`client.Realtime.Connect` and `client.Realtime.Translations.Connect` use the SDK's
shared WebSocket transport. Pass request options for the opening handshake and
`RealtimeConnectionOptions` or `TranslationConnectionOptions` for queue limits
and transport cleanup. Use `option.WithQuery("model", model)` to select a model.

Each connection offers typed `Send` and `Recv`, plus `SendRaw` and `RecvRaw` for
JSON bytes. Typed and raw reads consume the same ordered queue. Unknown event
types remain available through raw reads and the typed event's `RawJSON` method.
Close or abort every connection, even after an API operation finishes.

The opening context bounds connection setup. A receive context cancels only that
wait. A successful send confirms a local write, not server processing. Inspect
`websocket.DeliveryError` before deciding how to recover from a failed send.
Explicit `Reconnect` obtains fresh credentials and opens a new connection; it
does not replay commands or restore API session state.

## Finishing a Translation session

Use `client.Realtime.ConnectTranslation` when you want the SDK to finish and drain
a Translation session. It requires a model and rejects unsupported providers
before opening the connection. The initial `session.created` event remains in
the queue for the caller to read.

```go
session, err := client.Realtime.ConnectTranslation(
    ctx, model, realtime.TranslationConnectionOptions{},
)
if err != nil {
    return err
}
defer session.Abort()

// Send input with session.Send or session.SendRaw.
// Give Finish a deadline that allows the server to process the remaining input.
finishCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()
return session.Finish(finishCtx, handleTranslationEvent)
```

Implement `handleTranslationEvent` to process the transcript and audio events.

`Finish` stops new input, waits for admitted writes, sends `session.close` once,
and delivers events through `session.closed`. It then closes the transport.
`Close` alone closes the transport without finishing or draining the protocol.
The finish callback runs synchronously; it must return promptly. Return an error
to stop draining. Cancellation, early socket closure, invalid events, and callback
errors abort the connection. Repeating `Finish` returns its recorded result.

Only one `Recv`, `RecvRaw`, or `Finish` call may own event delivery at a time.
Sends may run concurrently up to `MaxPendingSends`. The helper never reconnects
or retries a command automatically.
