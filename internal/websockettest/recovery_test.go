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
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"), option.WithMaxRetries(9))
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
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"))
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

func TestResponsesRecoveryReleasesFailedRestoration(t *testing.T) {
	for _, tc := range []struct {
		name   string
		panics bool
	}{
		{"returned error", false},
		{"recovered panic", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var upgrades atomic.Int32
			replacementClosed := make(chan struct{})
			serverCtx, cancelServer := context.WithCancel(context.Background())
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempt := upgrades.Add(1)
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					t.Errorf("upgrade: %v", err)
					return
				}
				defer func() { _ = socket.CloseNow() }()
				kind, data, readErr := socket.Read(serverCtx)
				if attempt == 2 {
					close(replacementClosed)
				}
				if attempt < 3 {
					if readErr == nil {
						t.Error("unsuccessful restoration unexpectedly sent a command")
					}
					return
				}
				if readErr != nil {
					t.Errorf("read after explicit retry: %v", readErr)
					return
				}
				if err = socket.Write(serverCtx, kind, data); err != nil {
					t.Errorf("reply after explicit retry: %v", err)
					return
				}
				_, _, _ = socket.Read(serverCtx)
			}))
			defer server.Close()
			// Cancel peers during teardown even if the SDK fails to close one.
			defer cancelServer()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"))
			connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
			if err != nil {
				t.Fatal(err)
			}
			defer connection.Abort()

			panicMarker := &struct{ message string }{"restore interrupted"}
			wantErr := errors.New("restore failed")
			var recovered any
			var recoveryErr error
			var next *responses.ResponseConnection
			func() {
				defer func() { recovered = recover() }()
				next, recoveryErr = connection.Recover(ctx, responses.ResponseRecoveryOptions{
					Restore: func(context.Context, *responses.ResponseConnection) error {
						if tc.panics {
							panic(panicMarker)
						}
						return wantErr
					},
				})
			}()
			if tc.panics {
				if recovered != panicMarker {
					t.Fatalf("panic = %v, want original marker", recovered)
				}
			} else if recovered != nil || !errors.Is(recoveryErr, wantErr) {
				t.Fatalf("recovery error = %v, panic = %v, want restoration error", recoveryErr, recovered)
			}
			if next != nil {
				t.Fatal("failed restoration returned a connection")
			}
			// Only the retired original is reachable. The SDK must close the
			// unreturned replacement and release admission for another attempt.
			select {
			case <-replacementClosed:
			case <-time.After(time.Second):
				t.Error("failed restoration left its replacement open")
			}
			next, err = connection.Recover(ctx, responses.ResponseRecoveryOptions{})
			if err != nil {
				t.Fatalf("explicit retry after failed restoration: %v", err)
			}
			defer next.Abort()
			want := []byte(`{"type":"client.future","value":"one explicit command"}`)
			if err = next.SendRaw(ctx, want); err != nil {
				t.Fatal(err)
			}
			got, err := next.RecvRaw(ctx)
			if err != nil || string(got) != string(want) {
				t.Fatalf("retry received = %s, %v", got, err)
			}
			if got := upgrades.Load(); got != 3 {
				t.Fatalf("upgrade attempts = %d, want initial plus failed and successful restorations", got)
			}
		})
	}
}
