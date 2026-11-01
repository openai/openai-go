package websockettest

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3/packages/websocket"
	"github.com/openai/openai-go/v3/responses"
)

func TestResponsesCollectionRawReconstructionCorpus(t *testing.T) {
	data, err := os.ReadFile("testdata/responses_reconstruction.json")
	if err != nil {
		t.Fatal(err)
	}
	var corpus struct {
		Cases []struct {
			Name     string            `json:"name"`
			Events   []json.RawMessage `json:"events"`
			Expected struct {
				Result   string          `json:"result"`
				Response json.RawMessage `json:"response"`
			} `json:"expected"`
		} `json:"cases"`
	}
	if err = json.Unmarshal(data, &corpus); err != nil {
		t.Fatal(err)
	}
	for _, tc := range corpus.Cases {
		t.Run(tc.Name, func(t *testing.T) {
			frames := make([]string, len(tc.Events))
			for i, raw := range tc.Events {
				frames[i] = string(raw)
			}
			_, lane, _ := reconstructionConnection(t, frames)
			collection, collectErr := lane.CollectResponse(1 << 20)
			if collectErr != nil {
				t.Fatal(collectErr)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			for _, frame := range frames {
				raw, readErr := collection.RecvRaw(ctx)
				if readErr != nil {
					t.Fatal(readErr)
				}
				if string(raw) != frame {
					t.Fatal("raw event changed")
				}
				raw[0] = 'x'
			}
			response, finalErr := collection.FinalResponse(ctx)
			if tc.Expected.Result != "response" {
				if finalErr == nil || response != nil {
					t.Fatal("error or EOF became completion")
				}
				return
			}
			if finalErr != nil {
				t.Fatal(finalErr)
			}
			assertReconstructionJSON(t, response.RawJSON(), string(tc.Expected.Response))
			response.ID = "caller mutation"
			if again, againErr := collection.FinalResponse(ctx); againErr == nil || again != nil {
				t.Fatal("cached final was returned twice")
			}
		})
	}
}

func collectionConnection(t *testing.T, peer func(context.Context, *wire.Conn)) (*responses.ResponseConnection, *responses.ResponseLane, *responses.ResponseLane) {
	t.Helper()
	start := make(chan struct{})
	conn := laneTestConnection(t, responses.ResponseConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) {
		select {
		case <-start:
		case <-ctx.Done():
			return
		}
		peer(ctx, socket)
		_, _, _ = socket.Read(ctx)
	})
	a, err := conn.Lane("a")
	if err != nil {
		t.Fatal(err)
	}
	b, err := conn.Lane("b")
	if err != nil {
		t.Fatal(err)
	}
	close(start)
	return conn, a, b
}

func collectionWrite(t *testing.T, ctx context.Context, socket *wire.Conn, frames ...string) {
	t.Helper()
	for _, frame := range frames {
		if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
			t.Error(err)
			return
		}
	}
}

const collectionCreated = `{"type":"response.created","stream_id":"a","response":{"id":"r","status":"in_progress","output":[],"metadata":{"unicode":"雪🙂"},"future_collection":{"unicode":"雪🙂","exact":9007199254740993}}}`
const collectionDone = `{"type":"response.output_item.done","stream_id":"a","output_index":0,"item":{"type":"message","id":"m","content":[{"type":"output_text","text":"雪🙂","annotations":[]}]}}`
const collectionTerminal = `{"type":"response.completed","stream_id":"a","response":{"id":"r","status":"completed","output":null}}`

func TestResponsesCollectionExactBudgetAndIsolation(t *testing.T) {
	for _, overflow := range []bool{false, true} {
		name := "exact"
		if overflow {
			name = "one byte short"
		}
		t.Run(name, func(t *testing.T) {
			_, lane, other := collectionConnection(t, func(ctx context.Context, socket *wire.Conn) {
				collectionWrite(t, ctx, socket, collectionCreated, collectionDone, collectionTerminal, `{"type":"response.completed","stream_id":"b","response":{"id":"neighbor","status":"completed","output":[]}}`, `{"type":"future.event","stream_id":"a","value":true}`)
			})
			budget := int64(len(collectionCreated) + len(collectionDone) + len(collectionTerminal))
			if overflow {
				budget--
			}
			collection, err := lane.CollectResponse(budget)
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			for range 3 {
				if _, readErr := collection.RecvRaw(ctx); readErr != nil {
					t.Fatal(readErr)
				}
			}
			response, finalErr := collection.FinalResponse(ctx)
			if overflow {
				if !errors.Is(finalErr, responses.ErrResponseCollectionLimit) || response != nil {
					t.Fatal("cumulative overflow not isolated")
				}
			} else if finalErr != nil || response.OutputText() != "雪🙂" {
				t.Fatal("exact UTF-8 byte budget failed")
			}
			if _, readErr := collection.RecvRaw(ctx); readErr != nil {
				t.Fatal("collector failure stopped raw receiving")
			}
			neighbor, neighborErr := other.FinalResponse(ctx)
			if neighborErr != nil || neighbor.ID != "neighbor" {
				t.Fatal("collector stopped neighbor")
			}
			if overflow {
				if !errors.Is(collection.Err(), responses.ErrResponseCollectionLimit) {
					t.Fatal("raw read cleared collection failure")
				}
			}
		})
	}
}

type collectionWaitContext struct {
	context.Context
	entered chan struct{}
	once    sync.Once
}

func (c *collectionWaitContext) Done() <-chan struct{} {
	c.once.Do(func() { close(c.entered) })
	return c.Context.Done()
}

func TestResponsesCollectionCancellationAndSuccessor(t *testing.T) {
	advance := make(chan struct{})
	defer func() {
		select {
		case <-advance:
		default:
			close(advance)
		}
	}()
	_, lane, neighbor := collectionConnection(t, func(ctx context.Context, socket *wire.Conn) {
		collectionWrite(t, ctx, socket, collectionCreated, collectionDone)
		select {
		case <-advance:
		case <-ctx.Done():
			return
		}
		collectionWrite(t, ctx, socket, collectionTerminal, `{"type":"response.created","stream_id":"a","response":{"id":"next","output":[]}}`, `{"type":"response.completed","stream_id":"a","response":{"id":"next","status":"completed","output":[]}}`, `{"type":"future.event","stream_id":"b"}`)
	})
	collection, err := lane.CollectResponse(1 << 20)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for range 2 {
		if _, readErr := collection.RecvRaw(ctx); readErr != nil {
			t.Fatal(readErr)
		}
	}
	waitCtx, stop := context.WithCancel(ctx)
	observed := &collectionWaitContext{Context: waitCtx, entered: make(chan struct{})}
	settled := make(chan error, 1)
	go func() { _, finalErr := collection.FinalResponse(observed); settled <- finalErr }()
	select {
	case <-observed.entered:
	case <-ctx.Done():
		t.Fatal("idle wait did not start")
	}
	if _, overlapErr := collection.RecvRaw(ctx); overlapErr == nil {
		t.Fatal("overlapping consumer accepted")
	}
	if resetErr := collection.Reset(); resetErr == nil {
		t.Fatal("Reset mutated state during an active final receive")
	}
	stop()
	select {
	case waitErr := <-settled:
		if !errors.Is(waitErr, context.Canceled) {
			t.Fatal(waitErr)
		}
	case <-ctx.Done():
		t.Fatal("cancellation stalled")
	}
	close(advance)
	if _, readErr := collection.RecvRaw(ctx); readErr != nil {
		t.Fatal(readErr)
	}
	// The single reader routed this later marker only after queuing both
	// successor events on a. They must not reset the unconsumed final cache.
	if _, markerErr := neighbor.Recv(ctx); markerErr != nil {
		t.Fatal(markerErr)
	}
	first, firstErr := collection.FinalResponse(ctx)
	if firstErr != nil || first.ID != "r" || first.OutputText() != "雪🙂" {
		t.Fatal("cancel/resume lost first cache")
	}
	first.ID = "caller mutation"
	next, nextErr := collection.FinalResponse(ctx)
	if nextErr != nil || next.ID != "next" || len(next.Output) != 0 {
		t.Fatal("cached take or successor reset failed")
	}
}

func TestResponsesCollectionConsumedSuccessorAndDefaultScope(t *testing.T) {
	conn, _, _ := collectionConnection(t, func(ctx context.Context, socket *wire.Conn) {
		collectionWrite(t, ctx, socket, `{"type":"response.completed","stream_id":"unclaimed","response":{"id":"wrong","status":"completed","output":[]}}`, `{"type":"response.completed","response":{"id":"old","status":"completed","output":[]}}`, `{"type":"response.created","response":{"id":"new","output":[]}}`, `{"type":"response.completed","response":{"id":"new","status":"completed","output":[]}}`)
	})
	collection, err := conn.CollectResponse(1 << 20)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for range 3 {
		if _, readErr := collection.RecvRaw(ctx); readErr != nil {
			t.Fatal(readErr)
		}
	}
	response, finalErr := collection.FinalResponse(ctx)
	if finalErr != nil || response.ID != "new" {
		t.Fatal("consumed successor failed to reset, or named event contaminated default")
	}
}

func TestResponsesCollectionLocalClosureInvalidatesCache(t *testing.T) {
	for _, action := range []string{"close", "abort", "lane close", "reconnect", "recover", "failed recover"} {
		t.Run(action, func(t *testing.T) {
			conn, lane, _ := collectionConnection(t, func(ctx context.Context, socket *wire.Conn) { collectionWrite(t, ctx, socket, collectionTerminal) })
			collection, err := lane.CollectResponse(1 << 20)
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if _, readErr := collection.RecvRaw(ctx); readErr != nil {
				t.Fatal(readErr)
			}
			switch action {
			case "close":
				_ = conn.Close()
			case "abort":
				conn.Abort()
			case "lane close":
				lane.Close()
			case "reconnect":
				next, recoverErr := conn.Reconnect(ctx)
				if recoverErr != nil {
					t.Fatal(recoverErr)
				}
				defer func() { _ = next.Close() }()
			case "recover":
				next, recoverErr := conn.Recover(ctx, responses.ResponseRecoveryOptions{})
				if recoverErr != nil {
					t.Fatal(recoverErr)
				}
				defer func() { _ = next.Close() }()
			case "failed recover":
				next, recoverErr := conn.Recover(ctx, responses.ResponseRecoveryOptions{Restore: func(context.Context, *responses.ResponseConnection) error {
					return errors.New("synthetic restore failure")
				}})
				if next != nil || recoverErr == nil {
					t.Fatal("failed recovery unexpectedly succeeded")
				}
			}
			if response, finalErr := collection.FinalResponse(ctx); finalErr == nil || response != nil {
				t.Fatal("local lifecycle returned stale success")
			}
			if action == "lane close" && !errors.Is(collection.Err(), websocket.ErrLaneClosed) {
				t.Fatal("lane closure not distinguishable")
			}
		})
	}
}

func TestResponsesCollectionFailureResetAndUnchargedProgress(t *testing.T) {
	for _, failure := range []string{"overflow", "malformed", "protocol"} {
		t.Run(failure, func(t *testing.T) {
			bad := collectionDone
			budget := int64(len(collectionCreated) + len(collectionTerminal))
			if failure == "malformed" {
				bad = `{"type":"response.output_item.done","stream_id":"a","output_index":0,"item":{"type":"message","id":"bad","content":{}}}`
				budget = 1 << 20
			}
			if failure == "protocol" {
				bad = `{"type":"error","stream_id":"a","error":{"type":"server_error","message":"synthetic"}}`
				budget = 1 << 20
			}
			delta := `{"type":"response.output_text.delta","stream_id":"a","output_index":0,"content_index":0,"item_id":"m","delta":"` + strings.Repeat("x", 2<<20) + `"}`
			_, lane, _ := collectionConnection(t, func(ctx context.Context, socket *wire.Conn) {
				collectionWrite(t, ctx, socket, collectionCreated, bad, delta, `{"type":"response.created","stream_id":"a","response":{"id":"fresh","output":[]}}`, `{"type":"response.completed","stream_id":"a","response":{"id":"fresh","status":"completed","output":[]}}`)
			})
			collection, err := lane.CollectResponse(budget)
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			for range 3 {
				if _, readErr := collection.RecvRaw(ctx); readErr != nil {
					t.Fatal(readErr)
				}
			}
			if collection.Err() == nil {
				t.Fatal("collection failure not latched")
			}
			response, finalErr := collection.FinalResponse(ctx)
			if finalErr == nil || response != nil {
				t.Fatal("failed collector drained into success")
			}
			if failure == "protocol" {
				if resetErr := collection.Reset(); resetErr != nil {
					t.Fatal(resetErr)
				}
				if collection.Err() != nil {
					t.Fatal("explicit Reset retained a protocol error")
				}
			}
			for range 2 {
				if _, readErr := collection.RecvRaw(ctx); readErr != nil {
					t.Fatal(readErr)
				}
			}
			fresh, freshErr := collection.FinalResponse(ctx)
			if freshErr != nil || fresh.ID != "fresh" || len(fresh.Output) != 0 {
				t.Fatal("consumed created did not reset failure")
			}
		})
	}
	_, lane, _ := collectionConnection(t, func(ctx context.Context, socket *wire.Conn) {
		collectionWrite(t, ctx, socket, `{"type":"response.output_text.delta","stream_id":"a","output_index":0,"content_index":0,"item_id":"m","delta":"`+strings.Repeat("x", 2<<20)+`"}`, collectionTerminal)
	})
	collection, err := lane.CollectResponse(int64(len(collectionTerminal)))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, readErr := collection.Recv(ctx); readErr != nil {
		t.Fatal(readErr)
	}
	response, finalErr := collection.FinalResponse(ctx)
	if finalErr != nil || len(response.Output) != 0 {
		t.Fatal("irrelevant delta charged or promoted into final")
	}
}

func TestResponsesCollectionRequiresPositiveBudget(t *testing.T) {
	conn, lane, _ := collectionConnection(t, func(context.Context, *wire.Conn) {})
	for _, budget := range []int64{0, -1} {
		if _, err := conn.CollectResponse(budget); err == nil {
			t.Fatal("connection accepted missing budget")
		}
		if _, err := lane.CollectResponse(budget); err == nil {
			t.Fatal("lane accepted missing budget")
		}
	}
}

func TestResponsesCollectionLargeFinal(t *testing.T) {
	text := strings.Repeat("x", 65<<20)
	frame, err := json.Marshal(map[string]any{"type": "response.output_item.done", "stream_id": "a", "output_index": 0, "item": map[string]any{"type": "message", "id": "large", "content": []any{map[string]any{"type": "output_text", "text": text, "annotations": []any{}}}}})
	if err != nil {
		t.Fatal(err)
	}
	_, lane, _ := collectionConnection(t, func(ctx context.Context, socket *wire.Conn) {
		collectionWrite(t, ctx, socket, string(frame), collectionTerminal)
	})
	collection, collectErr := lane.CollectResponse(int64(len(frame) + len(collectionTerminal)))
	if collectErr != nil {
		t.Fatal(collectErr)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	for range 2 {
		if _, readErr := collection.RecvRaw(ctx); readErr != nil {
			t.Fatal(readErr)
		}
	}
	response, finalErr := collection.FinalResponse(ctx)
	if finalErr != nil || response.OutputText() != text {
		t.Fatalf("large collected final failed: %v", finalErr)
	}
}
