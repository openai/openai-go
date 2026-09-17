package websockettest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func TestResponsesRecoveryExhaustion(t *testing.T) {
	var upgrades atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if upgrades.Add(1) > 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("upgrade: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _, _ = socket.Read(ctx)
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"), option.WithMaxRetries(9))
	connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Abort()
	next, err := connection.Recover(ctx, responses.ResponseRecoveryOptions{MaxAttempts: 2})
	if err == nil || next != nil {
		t.Fatalf("exhausted recovery returned connection=%v error=%v", next, err)
	}
	if got := upgrades.Load(); got != 3 {
		t.Fatalf("upgrade attempts = %d, want initial plus two explicit attempts", got)
	}
}

func TestResponsesCloseCancelsRestoration(t *testing.T) {
	var upgrades atomic.Int32
	replacementClosed := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := upgrades.Add(1)
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("upgrade: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _, _ = socket.Read(ctx)
		if attempt == 2 {
			close(replacementClosed)
		}
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"))
	connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Abort()
	restoring := make(chan struct{})
	result := make(chan error, 1)
	go func() {
		_, recoverErr := connection.Recover(ctx, responses.ResponseRecoveryOptions{MaxAttempts: 3, Restore: func(restoreCtx context.Context, _ *responses.ResponseConnection) error {
			close(restoring)
			<-restoreCtx.Done()
			return restoreCtx.Err()
		}})
		result <- recoverErr
	}()
	select {
	case <-restoring:
	case <-ctx.Done():
		t.Fatal("restoration did not start")
	}
	_ = connection.Close()
	select {
	case err = <-result:
		if err == nil || errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("close should interrupt recovery before the deadline: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("close did not cancel restoration")
	}
	select {
	case <-replacementClosed:
	case <-ctx.Done():
		t.Fatal("failed restoration leaked its replacement socket")
	}
	if got := upgrades.Load(); got != 2 {
		t.Fatalf("upgrade attempts = %d, want initial plus one canceled restoration", got)
	}
}
