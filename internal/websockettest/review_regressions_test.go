package websockettest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type websocketSubjectTokenProvider struct{ calls atomic.Int32 }

func (*websocketSubjectTokenProvider) TokenType() auth.SubjectTokenType {
	return auth.SubjectTokenTypeJWT
}

func (p *websocketSubjectTokenProvider) GetToken(context.Context, auth.HTTPDoer) (string, error) {
	p.calls.Add(1)
	return "synthetic-subject-token", nil
}

type websocketAuthTransport func(*http.Request) (*http.Response, error)

func (f websocketAuthTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestResponsesWebSocketWorkloadIdentity(t *testing.T) {
	var upgrades, exchanges atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer synthetic-access-token" || r.Header.Get("X-Application") != "synthetic-app" {
			t.Error("Responses.Connect() did not preserve workload identity and custom headers")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		upgrades.Add(1)
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("accept workload identity upgrade: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := socket.Write(ctx, wire.MessageText, []byte(`{"type":"response.completed","response":{"id":"resp_workload","status":"completed"}}`)); err != nil {
			t.Errorf("write workload identity response: %v", err)
			return
		}
		_, _, _ = socket.Read(ctx)
	}))
	defer server.Close()
	transport := server.Client().Transport
	httpClient := &http.Client{Transport: websocketAuthTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == auth.TokenExchangeURL {
			exchanges.Add(1)
			return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"access_token":"synthetic-access-token","expires_in":3600}`))}, nil
		}
		if r.URL.String() == server.URL+"/responses" {
			return transport.RoundTrip(r)
		}
		return nil, errors.New("unexpected non-test destination")
	})}
	provider := &websocketSubjectTokenProvider{}
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(httpClient), option.WithWorkloadIdentity(auth.WorkloadIdentity{
		IdentityProviderID: "synthetic-idp", ServiceAccountID: "synthetic-account", Provider: provider,
	}))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{}, option.WithHeader("X-Application", "synthetic-app"))
	if err != nil {
		t.Fatalf("Responses.Connect(workload identity) = %v, want success", err)
	}
	defer connection.Abort()
	response, err := connection.FinalResponse(ctx)
	if err != nil || response.ID != "resp_workload" {
		t.Fatalf("FinalResponse(workload identity) = %v, %v, want resp_workload", response, err)
	}
	if upgrades.Load() != 1 || exchanges.Load() != 1 || provider.calls.Load() != 1 {
		t.Errorf("workload identity upgrades/exchanges/provider calls = %d/%d/%d, want 1/1/1", upgrades.Load(), exchanges.Load(), provider.calls.Load())
	}
}

func TestResponsesWebSocketLargeFinalResponse(t *testing.T) {
	// This probes compatibility above the previous 64 MiB default; it is not
	// an API maximum and deliberately does not raise any connection limit.
	want := strings.Repeat("x", 65<<20)
	payload, err := json.Marshal(map[string]any{
		"type": "response.completed",
		"response": map[string]any{
			"id": "resp_large", "status": "completed",
			"output": []any{map[string]any{
				"type": "message", "id": "msg_large", "role": "assistant", "status": "completed",
				"content": []any{map[string]any{"type": "output_text", "text": want, "annotations": []any{}}},
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, acceptErr := wire.Accept(w, r, nil)
		if acceptErr != nil {
			t.Errorf("accept large response upgrade: %v", acceptErr)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if writeErr := socket.Write(ctx, wire.MessageText, payload); writeErr != nil {
			t.Errorf("write large response: %v", writeErr)
			return
		}
		_, _, _ = socket.Read(ctx)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic-key"))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Abort()
	response, err := connection.FinalResponse(ctx)
	if err != nil {
		t.Fatalf("FinalResponse(65 MiB output) = %v, want success", err)
	}
	if got := response.OutputText(); got != want {
		t.Errorf("FinalResponse(65 MiB output) text length = %d, want %d", len(got), len(want))
	}
}

func TestResponsesWebSocketMalformedEventFailsAllLanes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("accept upgrade: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, _, err := socket.Read(ctx); err != nil {
			t.Errorf("read command after lanes attach: %v", err)
			return
		}
		for _, frame := range []string{
			`{"type":"response.output_text.delta","stream_id":"a","delta":"accepted"}`,
			`{"type":`,
		} {
			if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
				t.Errorf("write event: %v", err)
				return
			}
		}
		_, _, _ = socket.Read(ctx)
	}))
	t.Cleanup(server.Close)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic-key"))
	connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(connection.Abort)
	laneA, err := connection.Lane("a")
	if err != nil {
		t.Fatal(err)
	}
	laneB, err := connection.Lane("b")
	if err != nil {
		t.Fatal(err)
	}
	if err = connection.Create(ctx, responses.ResponsesClientEventResponseCreateParam{}); err != nil {
		t.Fatal(err)
	}
	// A different lane must observe the reader failure before lane A is consumed.
	_, err = laneB.Recv(ctx)
	var syntax *json.SyntaxError
	if !errors.As(err, &syntax) {
		t.Fatalf("lane B Recv(malformed frame after lane A event) = %v, want JSON syntax error", err)
	}
	event, err := laneA.Recv(ctx)
	if err != nil || event.AsResponseOutputTextDelta().Delta != "accepted" {
		t.Fatalf("lane A Recv(queued event before malformed frame) type = %q, error = %v, want accepted delta", event.Type, err)
	}
	if _, err := laneA.Recv(ctx); !errors.As(err, &syntax) {
		t.Errorf("lane A Recv(after queued event) = %v, want JSON syntax error", err)
	}
}
