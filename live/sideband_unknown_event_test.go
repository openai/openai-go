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

func TestSidebandDecodesTransportEventsAndPreservesUnknown(t *testing.T) {
	// Typed and future events must retain their complete wire payloads.
	frames := []string{
		`{"type":"transport.ringing","event_id":"ring","session_id":"live_test"}`,
		`{"type":"transport.answered","event_id":"answer","session_id":"live_test"}`,
		`{"type":"transport.failed","event_id":"fail","session_id":"live_test","error":{"type":"call_error","code":"provider_invite_failed","message":"synthetic failure","param":""}}`,
		`{"type":"transport.dtmf.received","event_id":"dtmf1","event":"5"}`,
		`{"type":"transport.dtmf.send","event_id":"dtmf2","event":"#","client_event_id":"command1"}`,
		// Null is outside the server schema; retain the existing lenient decode behavior.
		`{"type":"transport.dtmf.send","event_id":"null","event":"#","client_event_id":null}`,
		`{"type":"transport.dtmf.send","event_id":"omitted","event":"#"}`,
		`{"type":"session.input_audio.append","audio":"AAE="}`,
		`{"type":"session.output_audio.delta","delta":"AAE=","start_ms":10,"end_ms":20}`,
		`{"type":"future.transport.event","value":{"nested":42}}`,
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
		var envelope struct {
			Type          string
			ClientEventID json.RawMessage `json:"client_event_id"`
		}
		if err := json.Unmarshal([]byte(frame), &envelope); err != nil {
			t.Fatal(err)
		}
		event, err := connection.Recv(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if event.Type != envelope.Type || event.RawJSON() != frame {
			t.Fatalf("event lost data: type %q, raw %s", event.Type, event.RawJSON())
		}
		switch value := event.AsAny().(type) {
		case live.ConnectServerEventTransportRinging:
			if value.SessionID != "live_test" {
				t.Fatal(value)
			}
		case live.ConnectServerEventTransportAnswered:
			if value.SessionID != "live_test" {
				t.Fatal(value)
			}
		case live.ConnectServerEventTransportFailed:
			if value.Error.Message != "synthetic failure" {
				t.Fatal(value)
			}
		case live.ConnectServerEventTransportDtmfReceived:
			if value.Event != "5" {
				t.Fatal(value)
			}
		case live.ConnectServerEventTransportDtmfSend:
			expectedID := ""
			if value.EventID == "dtmf2" {
				expectedID = "command1"
			}
			if value.Event != "#" || value.ClientEventID != expectedID || value.JSON.ClientEventID.Raw() != string(envelope.ClientEventID) {
				t.Fatal(value)
			}
		case live.OutputAudioDeltaEvent:
			if value.Delta != "AAE=" || value.StartMs != 10 || value.EndMs != 20 {
				t.Fatal(value)
			}
		case live.ConnectServerEventSessionInputAudioAppend:
			if value.Audio != "AAE=" {
				t.Fatal(value)
			}
		default:
			if event.Type != "future.transport.event" || value != nil {
				t.Fatalf("unexpected variant %T for %s", value, event.Type)
			}
		}
	}
	if err := connection.Close(); err != nil {
		t.Fatal(err)
	}
}
