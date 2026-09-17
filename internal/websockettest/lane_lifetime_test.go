package websockettest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/websocket"
	"github.com/openai/openai-go/v3/responses"
)

func laneTestConnection(t *testing.T, options responses.ResponseConnectionOptions, handle func(context.Context, *wire.Conn)) *responses.ResponseConnection {
	t.Helper()
	var handlers sync.WaitGroup
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.Add(1)
		defer handlers.Done()
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("accept lane test connection: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		handle(ctx, socket)
	}))
	t.Cleanup(func() {
		server.Close()
		handlers.Wait()
	})
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connection, err := client.Responses.Connect(ctx, options)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(connection.Abort)
	return connection
}

func TestResponsesLaneReservedAfterTerminalAndAutomaticSuccessor(t *testing.T) {
	releaseSuccessor := make(chan struct{})
	connection := laneTestConnection(t, responses.ResponseConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) {
		if _, _, err := socket.Read(ctx); err != nil {
			t.Errorf("read parent command: %v", err)
			return
		}
		if err := socket.Write(ctx, wire.MessageText, []byte(`{"type":"response.completed","stream_id":"parent","response":{"id":"resp_parent"}}`)); err != nil {
			t.Errorf("write parent terminal: %v", err)
			return
		}
		select {
		case <-releaseSuccessor:
		case <-ctx.Done():
			return
		}
		// Steering can automatically start a successor after a terminal event.
		// No new client command is needed, and the stream ID stays the same.
		for _, frame := range []string{
			`{"type":"response.created","stream_id":"parent","response":{"id":"resp_successor"}}`,
			`{"type":"response.completed","stream_id":"parent","response":{"id":"resp_successor"}}`,
		} {
			if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
				t.Errorf("write automatic successor: %v", err)
				return
			}
		}
		for turn := range 2 {
			if _, _, err := socket.Read(ctx); err != nil {
				t.Errorf("read sequential command %d: %v", turn, err)
				return
			}
			frame := fmt.Sprintf(`{"type":"response.completed","stream_id":"kept-open","response":{"id":"resp_sequential_%d"}}`, turn)
			if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
				t.Errorf("write sequential response %d: %v", turn, err)
				return
			}
		}
		_, _, _ = socket.Read(ctx)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	parent, laneErr := connection.Lane("parent")
	if laneErr != nil {
		t.Fatal(laneErr)
	}
	sendFixture(t, ctx, connection.Create, json.RawMessage(`{"type":"response.create","stream_id":"parent","model":"gpt-4o-mini","input":"parent"}`))
	if response, err := parent.FinalResponse(ctx); err != nil || response.ID != "resp_parent" {
		t.Fatalf("FinalResponse(parent) = %v, %v, want resp_parent", response, err)
	}
	parent.Close()
	if replacement, err := connection.Lane("parent"); err == nil {
		replacement.Close()
		t.Error("Lane(parent) succeeded after terminal and Close, want reserved-ID error")
	}
	close(releaseSuccessor)
	if _, err := parent.Recv(ctx); !errors.Is(err, websocket.ErrLaneClosed) {
		t.Errorf("Recv(closed parent) = %v, want ErrLaneClosed", err)
	}
	if event, err := connection.Recv(ctx); err != nil || event.Type != "response.created" {
		t.Fatalf("Recv(automatic successor) = %v, %v, want response.created", event.Type, err)
	}
	if response, err := connection.FinalResponse(ctx); err != nil || response.ID != "resp_successor" {
		t.Fatalf("FinalResponse(automatic successor) = %v, %v, want resp_successor", response, err)
	}
	if replacement, err := connection.Lane("parent"); err == nil {
		replacement.Close()
		t.Error("Lane(parent) succeeded after successor completed, want reserved-ID error")
	}
	sequential, err := connection.Lane("kept-open")
	if err != nil {
		t.Fatal(err)
	}
	for turn := range 2 {
		sendFixture(t, ctx, connection.Create, json.RawMessage(`{"type":"response.create","stream_id":"kept-open","model":"gpt-4o-mini","input":"sequential"}`))
		want := fmt.Sprintf("resp_sequential_%d", turn)
		if response, err := sequential.FinalResponse(ctx); err != nil || response.ID != want {
			t.Fatalf("FinalResponse(sequential turn %d) = %v, %v, want %s", turn, response, err, want)
		}
	}
}

func TestResponsesLaneValidationCapacityAndRecovery(t *testing.T) {
	var connections, commands atomic.Int32
	connection := laneTestConnection(t, responses.ResponseConnectionOptions{MaxLanes: 3}, func(ctx context.Context, socket *wire.Conn) {
		attempt := connections.Add(1)
		for {
			_, data, err := socket.Read(ctx)
			if err != nil {
				return
			}
			commands.Add(1)
			var command struct {
				Input string `json:"input"`
			}
			if err := json.Unmarshal(data, &command); err != nil {
				t.Errorf("decode command on connection %d: %v", attempt, err)
				return
			}
			wantInput := fmt.Sprintf("connection_%d", attempt)
			if command.Input != wantInput {
				t.Errorf("command on connection %d input = %q, want %q; commands must not replay", attempt, command.Input, wantInput)
			}
			frame := fmt.Sprintf(`{"type":"response.completed","stream_id":"a","response":{"id":%q}}`, command.Input)
			if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
				t.Errorf("write response on connection %d: %v", attempt, err)
				return
			}
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, streamID := range []string{"", strings.Repeat("a", 257), "space key", "line\nkey", "tab\tkey", "slash/key", "key:part", "é", "😃", "key\x00", "key\x7f"} {
		if lane, err := connection.Lane(streamID); err == nil {
			lane.Close()
			t.Errorf("Lane(%q) succeeded, want invalid stream ID error", streamID)
		}
	}
	validIDs := []string{"a", strings.Repeat("a", 256), "Az09_.-"}
	for _, streamID := range validIDs {
		lane, err := connection.Lane(streamID)
		if err != nil {
			t.Fatalf("Lane(%q) = %v, want success after invalid IDs", streamID, err)
		}
		if streamID == "a" {
			sendFixture(t, ctx, connection.Create, json.RawMessage(`{"type":"response.create","stream_id":"a","model":"gpt-4o-mini","input":"connection_1"}`))
			if response, err := lane.FinalResponse(ctx); err != nil || response.ID != "connection_1" {
				t.Fatalf("FinalResponse(before recovery) = %v, %v, want connection_1", response, err)
			}
		}
		lane.Close()
	}
	if lane, err := connection.Lane("over-capacity"); err == nil {
		lane.Close()
		t.Error("Lane(over-capacity) succeeded after three closed lanes, want limit error")
	}
	replacement, err := connection.Reconnect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(replacement.Abort)
	for _, streamID := range validIDs {
		lane, err := replacement.Lane(streamID)
		if err != nil {
			t.Fatalf("Lane(%q) on replacement = %v, want fresh capacity", streamID, err)
		}
		if streamID == "a" {
			sendFixture(t, ctx, replacement.Create, json.RawMessage(`{"type":"response.create","stream_id":"a","model":"gpt-4o-mini","input":"connection_2"}`))
			if response, err := lane.FinalResponse(ctx); err != nil || response.ID != "connection_2" {
				t.Fatalf("FinalResponse(after recovery) = %v, %v, want connection_2 without replay", response, err)
			}
		}
		lane.Close()
	}
	if err := replacement.Close(); err != nil {
		t.Fatal(err)
	}
	if gotConnections, gotCommands := connections.Load(), commands.Load(); gotConnections != 2 || gotCommands != 2 {
		t.Errorf("recovery connections/commands = %d/%d, want 2/2 without replay", gotConnections, gotCommands)
	}
}
