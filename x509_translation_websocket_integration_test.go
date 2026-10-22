package openai_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/realtime"
)

func TestX509TranslationWebSocketPublicClient(t *testing.T) {
	var exchanges, upgrades atomic.Int32
	config := newX509WebSocketConfig(t, x509WebSocketIssuer(t, &exchanges), func(w http.ResponseWriter, r *http.Request) {
		upgrades.Add(1)
		if r.ProtoMajor != 1 || r.URL.Path != "/v1/realtime/translations" || r.URL.Query().Get("model") != "gpt-realtime-translate" ||
			r.Header.Get("Authorization") != "Bearer synthetic-websocket-token-1" || r.Header.Get("X-Custom-Handshake") != "synthetic-value" {
			t.Error("Translation X.509 handshake lost its origin, options or exchanged bearer")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("accept Translation X.509 WebSocket: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		_, payload, err := socket.Read(r.Context())
		var command map[string]any
		if err != nil || json.Unmarshal(payload, &command) != nil || command["type"] != "session.input_audio_buffer.append" {
			t.Error("Translation X.509 socket did not receive the typed command")
			return
		}
		_ = socket.Write(r.Context(), wire.MessageText, []byte(`{"type":"session.output_transcript.delta","event_id":"synthetic-event","delta":"Hola"}`))
		_, _, _ = socket.Read(r.Context())
	})
	client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetries(0))
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	var response *http.Response
	connection, err := client.Realtime.ConnectTranslation(ctx, "gpt-realtime-translate", realtime.TranslationConnectionOptions{}, option.WithHeader("X-Custom-Handshake", "synthetic-value"),
		option.WithResponseInto(&response))
	if err != nil {
		t.Fatalf("Translation Connect with attested X.509 transport: %v", err)
	}
	defer connection.Abort()
	if response == nil || response.StatusCode != http.StatusSwitchingProtocols || response.Request.Header.Get("Authorization") != "" {
		t.Error("Translation X.509 upgrade did not preserve status and redact its bearer")
	}
	if err = connection.Send(ctx, realtime.RealtimeTranslationClientEventUnionParam{
		OfSessionInputAudioBufferAppend: &realtime.RealtimeTranslationInputAudioBufferAppendEventParam{Audio: "AAAA"},
	}); err != nil {
		t.Fatalf("send Translation command: %v", err)
	}
	event, err := connection.Recv(ctx)
	if err != nil || event.Type != "session.output_transcript.delta" || event.AsRealtimeTranslationOutputTranscriptDeltaEvent().Delta != "Hola" {
		t.Fatalf("receive typed Translation event: type=%q error=%v", event.Type, err)
	}
	if exchanges.Load() != 1 || upgrades.Load() != 1 {
		t.Errorf("Translation exchange/upgrade counts=%d/%d, want 1/1", exchanges.Load(), upgrades.Load())
	}
}

func TestX509TranslationWebSocketRejectsUnsafeRequestsBeforeExchange(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*http.Request)
	}{
		{"different SDK endpoint", func(r *http.Request) { r.URL.Path = "/v1/responses" }},
		{"different endpoint", func(r *http.Request) { r.URL.Path = "/v1/models" }},
		{"different Realtime endpoint", func(r *http.Request) { r.URL.Path = "/v1/realtime/calls" }},
		{"escaped Translation", func(r *http.Request) { r.URL.RawPath = "/v1/realtime/%74ranslations" }},
		{"request body", func(r *http.Request) { r.Body = io.NopCloser(strings.NewReader("synthetic")) }},
		{"spoofed header", func(r *http.Request) { r.Header["connection"] = []string{"Upgrade"} }},
		{"different origin", func(r *http.Request) { r.URL.Host = "attacker.example.test" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var exchanges, upgrades atomic.Int32
			config := newX509WebSocketConfig(t, x509WebSocketIssuer(t, &exchanges),
				func(w http.ResponseWriter, _ *http.Request) { upgrades.Add(1); w.WriteHeader(http.StatusBadRequest) })
			client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetries(0),
				option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					tc.mutate(r)
					return next(r)
				}))
			connection, err := client.Realtime.ConnectTranslation(t.Context(), "gpt-realtime-translate", realtime.TranslationConnectionOptions{})
			if connection != nil {
				connection.Abort()
			}
			if err == nil || exchanges.Load() != 0 || upgrades.Load() != 0 {
				t.Errorf("unsafe Translation accepted=%t exchanges/upgrades=%d/%d", err == nil, exchanges.Load(), upgrades.Load())
			}
		})
	}
}
