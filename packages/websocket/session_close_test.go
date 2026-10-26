package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	wire "github.com/coder/websocket"
)

type gatedCloseTransport struct {
	base    http.RoundTripper
	entered chan struct{}
	release chan struct{}
}

func (p gatedCloseTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	res, err := p.base.RoundTrip(r)
	if err == nil && res.StatusCode == http.StatusSwitchingProtocols {
		res.Body = &gatedCloseBody{ReadWriteCloser: res.Body.(io.ReadWriteCloser), entered: p.entered, release: p.release}
	}
	return res, err
}

type gatedCloseBody struct {
	io.ReadWriteCloser
	once    sync.Once
	entered chan struct{}
	release chan struct{}
}

func (b *gatedCloseBody) Close() error {
	b.once.Do(func() { close(b.entered); <-b.release })
	return b.ReadWriteCloser.Close()
}

func TestSessionFailureDoesNotWaitForTransportClose(t *testing.T) {
	for _, mode := range []string{"direct", "session-recovery-off", "session-recovery-on"} {
		t.Run(mode, func(t *testing.T) {
			ctx := context.Background()
			peers := make(chan *wire.Conn, 1)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ws, err := wire.Accept(w, r, nil)
				if err != nil {
					return
				}
				defer func() { _ = ws.CloseNow() }()
				peers <- ws
				_, _, _ = ws.Read(r.Context())
			}))
			defer srv.Close()
			tr := gatedCloseTransport{base: srv.Client().Transport, entered: make(chan struct{}), release: make(chan struct{})}
			client := &http.Client{Transport: tr}
			opener := func(ctx context.Context) (*Connection[testEvent], error) {
				req, _ := http.NewRequestWithContext(ctx, "GET", srv.URL, nil)
				return Dial(ctx, req, client, Options{}, func(data []byte) (testEvent, error) {
					var e testEvent
					err := json.Unmarshal(data, &e)
					return e, err
				}, func(testEvent) string { return "" })
			}
			var recv func(context.Context) (testEvent, error)
			var abort func()
			if mode != "direct" {
				s, err := OpenSession(ctx, opener, RecoveryOptions{Enabled: mode == "session-recovery-on"})
				if err != nil {
					t.Fatal(err)
				}
				recv, abort = s.Recv, s.Abort
			} else {
				c, err := opener(ctx)
				if err != nil {
					t.Fatal(err)
				}
				recv, abort = c.Recv, c.Abort
			}
			defer abort()
			received := make(chan error, 1)
			receiverDone := make(chan struct{})
			go func() {
				defer close(receiverDone)
				_, err := recv(ctx)
				received <- err
			}()
			defer func() {
				close(tr.release)
				select {
				case <-receiverDone:
				case <-time.After(5 * time.Second):
					t.Error("receiver did not exit after transport close was released")
				}
			}()
			peer := <-peers
			if err := peer.Write(ctx, wire.MessageText, []byte("{broken")); err != nil {
				t.Fatal(err)
			}
			select {
			case <-tr.entered:
			case <-time.After(5 * time.Second):
				t.Fatal("transport cleanup never entered")
			}
			select {
			case err := <-received:
				var syntax *json.SyntaxError
				if !errors.As(err, &syntax) {
					t.Fatalf("got %v, want the original JSON syntax error", err)
				}
				t.Logf("receive settled while close pending: %v", err)
			case <-time.After(5 * time.Second):
				t.Error("receive still pending behind transport close")
			}
		})
	}
}
