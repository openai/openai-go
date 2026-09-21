package websockettest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type scenario struct {
	ID        string `json:"id"`
	Turns     []turn `json:"turns"`
	CloseCode int    `json:"close_code"`
}
type turn struct {
	Request json.RawMessage   `json:"request"`
	Frames  []json.RawMessage `json:"frames"`
}

func loadScenarios(t *testing.T) []scenario {
	t.Helper()
	data, err := os.ReadFile("testdata/scenarios.json")
	if err != nil {
		t.Fatal(err)
	}
	var corpus struct {
		Version   int        `json:"version"`
		Scenarios []scenario `json:"scenarios"`
	}
	if err = json.Unmarshal(data, &corpus); err != nil {
		t.Fatal(err)
	}
	if corpus.Version != 1 {
		t.Fatalf("unsupported corpus version: %d", corpus.Version)
	}
	return corpus.Scenarios
}

func sendFixture[T any](t *testing.T, ctx context.Context, send func(context.Context, T) error, request json.RawMessage) {
	t.Helper()
	var event T
	if err := json.Unmarshal(request, &event); err != nil {
		t.Fatalf("decode typed command: %v", err)
	}
	if err := send(ctx, event); err != nil {
		t.Fatalf("send typed command: %v", err)
	}
}

func TestResponsesWebSocketSharedContract(t *testing.T) {
	for _, scenario := range loadScenarios(t) {
		t.Run(scenario.ID, func(t *testing.T) {
			var upgrades atomic.Int32
			serverResult := make(chan error, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				upgrades.Add(1)
				if r.URL.Path != "/v1/responses" || r.Header.Get("Authorization") != "Bearer fixture-key" {
					serverResult <- errors.New("incorrect upgrade path or authentication")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					serverResult <- err
					return
				}
				defer func() { _ = socket.CloseNow() }()
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				for _, turn := range scenario.Turns {
					_, request, readErr := socket.Read(ctx)
					if readErr != nil {
						serverResult <- readErr
						return
					}
					var want, got any
					if err = json.Unmarshal(turn.Request, &want); err != nil {
						serverResult <- err
						return
					}
					if err = json.Unmarshal(request, &got); err != nil {
						serverResult <- err
						return
					}
					if !reflect.DeepEqual(got, want) {
						serverResult <- fmt.Errorf("wire command = %s, want %s", request, turn.Request)
						return
					}
					for _, frame := range turn.Frames {
						var raw string
						payload := []byte(frame)
						if json.Unmarshal(frame, &raw) == nil {
							payload = []byte(raw)
						}
						if err = socket.Write(ctx, wire.MessageText, payload); err != nil {
							serverResult <- err
							return
						}
					}
				}
				if scenario.CloseCode != 0 {
					_ = socket.Close(wire.StatusCode(scenario.CloseCode), "")
					serverResult <- nil
					return
				}
				_, _, err = socket.Read(ctx)
				if wire.CloseStatus(err) != wire.StatusNormalClosure && wire.CloseStatus(err) != wire.StatusGoingAway {
					serverResult <- fmt.Errorf("client did not close: %w", err)
					return
				}
				serverResult <- nil
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL+"/v1"), option.WithAPIKey("fixture-key"))
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
			if err != nil {
				t.Fatal(err)
			}
			defer connection.Abort()
			for _, turn := range scenario.Turns {
				sendFixture(t, ctx, connection.Create, turn.Request)
				for _, frame := range turn.Frames {
					var malformed string
					event, receiveErr := connection.Recv(ctx)
					if json.Unmarshal(frame, &malformed) == nil {
						var syntax *json.SyntaxError
						if !errors.As(receiveErr, &syntax) {
							t.Fatalf("malformed frame error = %v, want JSON syntax error", receiveErr)
						}
						continue
					}
					if receiveErr != nil {
						t.Fatalf("receive: %v", receiveErr)
					}
					var expected, actual map[string]any
					if err = json.Unmarshal(frame, &expected); err != nil {
						t.Fatal(err)
					}
					if err = json.Unmarshal([]byte(event.RawJSON()), &actual); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(actual, expected) {
						t.Fatalf("received event = %s, want %s", event.RawJSON(), frame)
					}
					// The generated union exposes typed variants; the raw envelope also
					// preserves future event kinds that have no current typed variant.
					switch expected["type"] {
					case "response.completed":
						typed, ok := event.AsAny().(responses.ResponsesServerEventResponseWsCompleted)
						if !ok {
							t.Fatalf("completed AsAny() = %T, want typed completed event", event.AsAny())
						}
						if string(typed.Type) != expected["type"] || typed.Response.ID == "" {
							t.Errorf("completed event lost its typed response")
						}
						if expected["stream_id"] != nil && typed.StreamID != expected["stream_id"] {
							t.Errorf("typed routing = %q, want %s", typed.StreamID, expected["stream_id"])
						}
					case "error":
						typed, ok := event.AsAny().(responses.ResponsesServerEventResponseWsError)
						if !ok {
							t.Fatalf("error AsAny() = %T, want typed protocol error", event.AsAny())
						}
						if typed.Error.Message == "" || string(typed.Type) != "error" {
							t.Error("protocol error lost its typed nested details")
						}
					case "response.future_event":
						if event.AsAny() != nil {
							t.Errorf("unknown AsAny() = %T, want nil with raw JSON retained", event.AsAny())
						}
					}
				}
			}
			if scenario.CloseCode != 0 {
				_, err = connection.Recv(ctx)
				if wire.CloseStatus(err) != wire.StatusCode(scenario.CloseCode) {
					t.Fatalf("premature close = %v, want code %d", err, scenario.CloseCode)
				}
			}
			if err = connection.Close(); err != nil {
				t.Fatal(err)
			}
			select {
			case err = <-serverResult:
				// Decode errors abort the transport; the malformed scenario may therefore
				// end without a WebSocket close handshake.
				if err != nil && scenario.ID != "malformed_json" {
					t.Fatal(err)
				}
			case <-ctx.Done():
				t.Fatal("server teardown timed out")
			}
			if upgrades.Load() != 1 {
				t.Fatalf("physical upgrades = %d, want one", upgrades.Load())
			}
		})
	}
}
