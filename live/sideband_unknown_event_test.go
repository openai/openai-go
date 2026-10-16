package live_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/option"
)

func TestSidebandPreservesUnmodeledTransportEvents(t *testing.T) {
	// The SDK schema intentionally excludes these SIP event variants. They must
	// still be available as complete raw events to sideband consumers.
	frames := []string{
		`{"type":"transport.ringing","event_id":"ring","session_id":"live_test"}`,
		`{"type":"transport.answered","event_id":"answer","session_id":"live_test"}`,
		`{"type":"transport.failed","event_id":"fail","session_id":"live_test","error":{"type":"call_error","code":"provider_invite_failed","message":"synthetic failure","param":""}}`,
		`{"type":"transport.dtmf.received","event_id":"dtmf1","event":"5"}`,
		`{"type":"transport.dtmf.send","event_id":"dtmf2","event":"#","client_event_id":"command1"}`,
	}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		for _, frame := range frames {
			if err := socket.Write(r.Context(), wire.MessageText, []byte(frame)); err != nil {
				t.Error(err)
				return
			}
		}
		_, _, _ = socket.Read(r.Context())
	}))
	defer server.Close()
	client := openai.NewClient(option.WithAPIKey("synthetic"), option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connection, err := client.Live.Sideband.Connect(ctx, "live_test", live.SidebandConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Abort()
	for _, frame := range frames {
		var envelope struct{ Type string }
		if err := json.Unmarshal([]byte(frame), &envelope); err != nil {
			t.Fatal(err)
		}
		event, err := connection.Recv(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if event.Type != envelope.Type || event.RawJSON() != frame || event.AsAny() != nil {
			t.Fatalf("unmodeled event lost data: type %q, raw %s", event.Type, event.RawJSON())
		}
	}
	if err := connection.Close(); err != nil {
		t.Fatal(err)
	}
}
