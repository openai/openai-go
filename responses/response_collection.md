# Opt-in WebSocket final collection

`CollectResponse(maxBytes)` creates a receive wrapper for an existing Responses
connection or named lane. Choose a positive budget for your application. The SDK
does not add a default collection limit or change existing receive methods.

```go
lane, err := connection.Lane("work")
if err != nil { return err }
collection, err := lane.CollectResponse(applicationResponseByteBudget)
if err != nil { return err }

// Send through connection.Create with StreamID: openai.String("work") so the
// response routes to this lane. Consume this receiver exclusively through
// collection, using either typed events or original JSON.
raw, err := collection.RecvRaw(ctx)
if err != nil { return err }
_ = raw // Application-owned inspection or UI.

// This uses any final already consumed above, or receives until the next final.
response, err := collection.FinalResponse(ctx)
if err != nil { return err }
_ = response
```

The wrapper observes events only when you receive them. Direct reads through the
original connection or lane bypass collection. Use one wrapper as the exclusive
consumer of its receiver. Overlapping wrapper receive calls fail promptly; do
not call `Err` concurrently with another method. `Reset` returns an error without
mutation while a receive method is active. There is no extra
reader goroutine, and the wrapper does not own or close the socket.

A connection wrapper collects only the default stream. It still returns
unclaimed named-stream events without using them for collection. Attach named
lanes before sending their requests, and give each lane its own wrapper.

`FinalResponse` takes a collected result once. Its next call receives the next
turn, rather than returning the old result. A consumed `response.created`
automatically resets collection, even if the prior final has not been taken.
A successor that is merely queued does not erase the current cache. `Reset`
clears local state and errors without consuming anything; check its returned error.
It does not discard queued events: drain the prior turn before resetting and
sending a new request. Sending through the original connection does not itself
reset collection; consuming the next `response.created` does.

The collector retains only lifecycle metadata and finalized `output_item.done`
records. Missing or null terminal output is reconstructed from those records in
output-index order. Supplied output, including `[]`, wins. Unfinished deltas
never become final items. Failed and incomplete responses retain their status.
Returned responses and original JSON are caller-owned; later turns do not change
previously returned data.

The budget counts the cumulative **raw incoming message bytes** for
`response.created`, `response.output_item.done`, `response.completed`,
`response.failed` and `response.incomplete` in one turn. UTF-8, whitespace and
unknown metadata count exactly as received. This metric differs from Ruby's
reserialized event-byte metric. Other progress and API error events are not
charged. This is a final-collection budget, not a queue, message or total heap
limit; existing connection limits remain separate.

Overflow returns `ErrResponseCollectionLimit` from `Err` and `FinalResponse`.
It releases retained collection data but does not stop raw receiving or other
lanes. Collection errors stay latched until `Reset` or a consumed
`response.created`. API errors retain their typed event in `ResponseProtocolError`.
Malformed selected data and early EOF cannot become a successful final result.
Canceling a receive wait preserves progress so you can resume with a fresh context.

Closing or aborting the original connection, closing its lane, or retiring it
through `Reconnect`/`Recover` invalidates the wrapper. Create a new wrapper on the
replacement connection. A peer EOF alone does not erase an actual terminal that
was already collected. Close the original connection or lane when finished.

Existing connection/lane `FinalResponse`, `Recv`, `RecvRaw` and caller-fed
`ResponseAccumulator` behavior is unchanged. They do not enable this collection
mode automatically.
