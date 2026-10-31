# Responses WebSockets

Use `client.Responses.Connect(ctx, responses.ResponseConnectionOptions{}, opts...)`
to open a reusable Responses connection. Client and per-connection request options,
including `option.WithHeader`, apply to the opening handshake. See the
[WebSocket example](../examples/responses-websocket/main.go).
`option.WithResponseInto` exposes the opening HTTP response. A rejected workload-identity
handshake can refresh its token within the configured retry budget; application
messages are never replayed. These connections require native HTTP transport
support and are unavailable in `js/wasm`; other SDK APIs remain usable there.

The default options do not cap message size or the receive buffer. Set
`MaxMessageBytes`, `MaxBufferedBytes` or `MaxBufferedEvents` to a positive value
to opt into a receive limit. Connections default to 16 pending sends and 64
distinct lane IDs over the connection's lifetime, including closed lanes.
Reuse an open lane for sequential responses. Exceeding an explicit receive bound
fails the connection without silently dropping events. Consume incoming events promptly or configure a buffer
limit appropriate for the application.
`CloseTimeout` defaults to five seconds and tears down an unresponsive connection
after that interval. Always call `Close` or `Abort` when finished.

## Turns, lanes, and recovery

Responses WebSockets use `response.create` and the normal Responses events, not
Realtime sessions. Leave `stream` and `background` unset. Use `Create` followed by
`FinalResponse` for a complete terminal snapshot, including tool calls; failed and
incomplete responses retain their status. `Recv` exposes individual typed events
and their `RawJSON`. A terminal response ends a turn, not the socket.

For opt-in incremental text and tool-input snapshots, feed the events you read
to a `responses.ResponseAccumulator`. It consumes no events itself:

```go
var acc responses.ResponseAccumulator
for {
    event, err := connection.Recv(ctx) // Or lane.Recv(ctx) for a named lane.
    if err != nil { return err }
    // The original typed event and event.RawJSON() remain available here.
    if err := acc.AddEvent(event); err != nil { return err }
    if event.Type == "response.completed" || event.Type == "response.failed" || event.Type == "response.incomplete" {
        snapshot := acc.Snapshot()
        // Inspect TerminalEvent: completed, failed and incomplete all end a turn.
        // Function arguments and custom-tool input in snapshot.Output are data;
        // the application decides whether and when to act on them.
        fmt.Println(snapshot.OutputText())
        break
    }
}
acc.Reset() // Ready for another turn; does not close the connection.
```

Use a separate accumulator for each lane. Snapshots remain unchanged as more
events arrive or after `Reset`. Text is grouped by output and content index;
function arguments and custom input retain their item and call IDs. Finalized
fields replace earlier deltas. `OutputText()` is the entire current projection,
not a delta: read a snapshot at a terminal or on explicit demand. Use the original
typed events for a live UI; later done events can shorten or correct earlier text.
A supplied final output, even an empty array,
overrides prior items; if final output is absent or null, the opt-in projection
retains collected fields. An item with a different nonempty ID starts fresh
at its output index; it cannot inherit text or tool inputs from an earlier item.
This projection is not a server `Response`. A missing or invalid terminal
response, an API error, or EOF never becomes a completed result.
`FinalResponse` keeps its original behavior and returns the server snapshot.

Use `acc.DetailedSnapshot()` for a mutable copy of all observed response metadata,
item, part and annotation fields, beyond the fields selected by `Snapshot()`.
`Response` is the last observed lifecycle response metadata, without `output`
(nil until one arrives). Output entries are sorted by `OutputIndex`;
`Content` and `Annotations` use sparse int64 content and annotation indices.
Fields are `json.RawMessage`: they preserve exact numbers, omitted versus null,
and fields from unknown variants, rather than manufacturing a typed complete
`Response`. Message content and list-valued annotations are extracted to those
indexed maps; null annotations remain null in the part fields. Unknown item
fields, including non-message content, remain in `Item`.

```go
details := acc.DetailedSnapshot() // At a terminal, or when explicitly requested.
for _, item := range details.Output {
    for contentIndex, part := range item.Content {
        text := part["text"]                   // nil if never received
        logprobs := part["logprobs"]           // provisional streamed form
        annotations := item.Annotations[contentIndex]
        _, _, _ = text, logprobs, annotations // Use in application-owned UI/state.
    }
}
```

Streamed logprobs accumulate with deltas. Supplied text-done logprobs replace them
(including empty or null); whole part/item replacements discard prior fields
and citations at those positions. Annotations and non-text parts enrich a
matching known item; an annotation-only event never binds a lane, starts a
turn or replaces an unrelated item. Standalone unsupported events and parts
without a known item remain available on the original received event. Tools and
unknown/partial data are never executed or promoted into final validated models.
Every returned map and raw JSON byte slice is a separate caller-owned copy.
Reading either snapshot joins/copies full accumulated output; reading one
after every delta would repeatedly rebuild growing prefixes.

Continue with `PreviousResponseID: openai.String(previous.ID)` and **only new input**.
For example, return a function result using its original call ID, then add a new
user message. The application executes tools and retains their results:

```go
output := responses.ResponseInputItemParamOfFunctionCallOutput("tool result")
output.OfFunctionCallOutput.CallID = openai.String(callID)
input := responses.ResponseInputParam{
    output,
    responses.ResponseInputItemParamOfMessage("Explain the result.", responses.EasyInputMessageRoleUser),
}
err := connection.Create(ctx, responses.ResponsesClientEventResponseCreateParam{
    Model: openai.ChatModelGPT5_2,
    PreviousResponseID: openai.String(previous.ID),
    Input: responses.ResponsesClientEventResponseCreateInputUnionParam{OfResponse: &input},
})
```

For an optional warmup, add `generate: false` to a create request using
`request.SetExtraFields(map[string]any{"generate": false})`. Receive its terminal
response and use its ID as the next `PreviousResponseID`; warmup produces no model
output. Clear that extra field with `request.SetExtraFields(nil)` if reusing the
request for a generated response. The
[runnable example](../examples/responses-websocket/main.go) shows a normal two-turn
exchange on one connection.

`StreamID` controls routing and server FIFO execution; `PreviousResponseID`
controls lineage. Reusing a lane without a parent ID starts a new chain. Register
`connection.Lane("fork")` before sending its create; Go lanes receive events, so
set `StreamID` on commands sent through the connection. One physical reader routes
interleaved events to their lanes. Sequential creates on one lane execute FIFO at
the service; different lanes may run concurrently. Keep consuming the connection's
unclaimed/default stream for connection-scoped errors too.

To fork a `store=false`/ZDR response, use its parent ID on a different lane and wait
for that lane's `response.in_progress` **before advancing the source lane**:

```go
fork, err := connection.Lane("fork")
if err != nil { return err }
defer fork.Close()
if err := connection.Create(ctx, responses.ResponsesClientEventResponseCreateParam{
    Model: openai.ChatModelGPT5_2, Store: openai.Bool(false),
    StreamID: openai.String("fork"), PreviousResponseID: openai.String(previous.ID),
    Input: responses.ResponsesClientEventResponseCreateInputUnionParam{OfString: openai.String("Explore another approach.")},
}); err != nil { return err }
for {
    event, err := fork.Recv(ctx)
    if err != nil { return err }
    if event.Type == "error" { return &responses.ResponseProtocolError{Event: event} }
    if event.Type == "response.in_progress" { break }
    if event.Type == "response.completed" || event.Type == "response.failed" || event.Type == "response.incomplete" {
        return fmt.Errorf("fork ended before becoming ready")
    }
}
// The source lane may now advance; continue draining both lanes.
```

The service keeps recent parents in a connection-local cache. With `store=true`,
older IDs may be hydrated from persisted state; with `store=false` or ZDR an
uncached ID fails with `previous_response_not_found`. A same-lane continuation
returning 4xx/5xx evicts its referenced parent; a failed cross-lane fork preserves
the shared parent for the source lane. The SDK does not maintain this server cache.

## Compaction and service limits

Server compaction uses `ContextManagement` entries with `Type: "compaction"` and
`CompactThreshold: openai.Int(20000)`. Continue normally with the latest response
ID and new input. Standalone `client.Responses.Compact` instead returns a compacted
input window. Preserve **all** its output items, including unknown fields, and
start a new chain with `PreviousResponseID` omitted or null. Its compaction resource
ID is not a continuation response ID. Using `encoding/json` and `packages/param`:

```go
var input responses.ResponseInputParam
for _, item := range compacted.Output {
    input = append(input, param.Override[responses.ResponseInputItemUnionParam](json.RawMessage(item.RawJSON())))
}
err := connection.Create(ctx, responses.ResponsesClientEventResponseCreateParam{
    Model: openai.ChatModelGPT5_2,
    Input: responses.ResponsesClientEventResponseCreateInputUnionParam{OfResponse: &input},
})
```

The [WebSocket mode guide](https://developers.openai.com/api/docs/guides/websocket-mode)
defines service limits separately from SDK buffering: up to 16 active responses
(additional creates queue), 32 distinct named streams per connection (the default
lane does not count), and a 60-minute connection lifetime. Stream IDs contain
1–256 letters, digits, underscores, hyphens or periods. Omit the ID for the default
lane; an empty string is invalid. Closing a local lane does not free a distinct
stream ID at the service. SDK `MaxLanes` bounds distinct lane IDs over the
connection lifetime, including closed lanes; `MaxPendingSends` bounds local
writes, not active server responses.

`Recv` returns protocol errors as events; `FinalResponse` returns a
`*responses.ResponseProtocolError` retaining the event. Inspect the nested code:
`previous_response_not_found` requires a usable persisted parent or a fresh chain
with saved full context; `invalid_stream_id` requires correcting the ID;
`websocket_stream_limit_reached` requires reusing an existing stream or explicitly
opening a new connection; `websocket_connection_limit_reached` requires a new
connection. Named request errors route to their lane; no-ID errors route to the
connection. Errors do not by themselves imply a successful turn.

`Reconnect`/`Recover` are explicit and never replay old commands or ambiguous
writes. Every lane's server cache is lost, even if the old ID was recently used.
On the replacement, register needed lanes again and resume from a persisted ID
when available, or send saved full context as a new chain without a parent ID.
Retain tool results so recovery does not execute tools twice. A `Restore` callback
may run for each explicitly configured recovery attempt; it owns any commands it
sends and their replay/idempotency decisions. Do not automatically retry a
`websocket.DeliveryError` whose `MayHaveBeenSent` is true.

