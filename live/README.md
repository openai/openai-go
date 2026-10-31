# Live WebSocket sessions

Live has three connection roles: a primary session, a sideband attached to an
existing session, and a fork of a stored session. Each provides typed `Send` and
`Recv` methods, plus `SendRaw` and `RecvRaw` for JSON bytes.

These snippets assume a configured `client`, a bounded `ctx`, and an accessible
`model`. Import `github.com/openai/openai-go/v3/live`; the primary example also
uses `errors`, and the sideband example uses `github.com/openai/openai-go/v3/option`.
`handleEvent` is application code that consumes typed Live events.

## Start and finish a primary session

Opening the socket does not start the session. Send `session.start` first and
wait for `session.started` before sending application commands.

```go
conn, err := client.Live.Connect(ctx, live.LiveConnectionOptions{})
if err != nil {
	return err
}
defer conn.Abort()

start := live.ClientEventParamOfSessionStart(live.SessionConfigParam{Model: model})
if err := conn.Send(ctx, start); err != nil {
	return err
}
for {
	event, err := conn.Recv(ctx)
	if err != nil {
		return err
	}
	if event.Type == "error" || event.Type == "session.closed" {
		return errors.New("Live session did not start")
	}
	if event.Type == "session.started" {
		break
	}
}

// Send audio and consume events here before ending the session.
if err := conn.Send(ctx, live.ClientEventUnionParam{
	OfSessionClose: &live.SessionCloseEventParam{},
}); err != nil {
	return err
}
for {
	event, err := conn.Recv(ctx)
	if err != nil {
		return err
	}
	if event.Type == "error" {
		return errors.New("Live session returned an API error")
	}
	handleEvent(event)
	if event.Type == "session.closed" {
		return conn.Close()
	}
}
```

Keep receiving after `session.close` to deliver trailing output through
`session.closed`. A transport close, timeout, or canceled read is not proof of
session finalization. `Close` closes the socket; it does not send `session.close`.
The application owns audio capture/playback and handles API error events.

## Attach a sideband

Use the ID of an eligible signaling session, such as one created through the
WebRTC signaling API, and the matching API key or bound observer credential.
An arbitrary primary WebSocket session ID or membership in the same project is
not enough. Set `graceful_close=false` so closing this observer's socket leaves
the underlying session running:

```go
conn, err := client.Live.Sideband.Connect(
	ctx, sessionID, live.SidebandConnectionOptions{},
	option.WithQuery("graceful_close", "false"),
)
if err != nil {
	return err
}
defer conn.Abort()
```

Attach to the running session; do not send another `session.start` or wait for a
fresh `session.started`. The sideband uses its own typed command and event unions.
Closing the sideband socket releases
that connection; sending `session.close` acts on the session itself.

## Fork a stored session

A fork requires an eligible stored session from the same project, with a
finalized, nonempty recording. A live session ID alone is insufficient. Check
the effective storage setting and recording availability before attempting a
fork; requesting storage does not guarantee that the recording is available.

```go
conn, err := client.Live.Forks.Connect(ctx, storedSessionID, live.ForkConnectionOptions{})
if err != nil {
	return err
}
defer conn.Abort()
start := live.ForkClientEventParamOfSessionStart(live.ForkSessionConfigParam{})
if err := conn.Send(ctx, start); err != nil {
	return err
}
```

Wait for `session.started`, then use the fork's typed events. Finish and drain
through `session.closed` as with a primary session, using the fork command union.

## Connection ownership

Pass request options for the opening handshake and connection options for
transport limits. The opening context bounds setup; each receive context bounds
only that wait. Typed and raw reads share one ordered queue. Close or abort every
connection. Explicit reconnect opens a new socket with fresh credentials; it
does not replay commands or restore the Live session.
