package websockettest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/realtime"
)

type closeableConnection interface{ Close() error }

func TestBindingHandshakes(t *testing.T) {
	variants := []struct {
		name    string
		connect func(context.Context, *openai.Client, ...option.RequestOption) (closeableConnection, error)
	}{
		{"realtime", func(ctx context.Context, c *openai.Client, opts ...option.RequestOption) (closeableConnection, error) {
			return c.Realtime.Connect(ctx, realtime.RealtimeConnectionOptions{}, opts...)
		}},
		{"live", func(ctx context.Context, c *openai.Client, opts ...option.RequestOption) (closeableConnection, error) {
			return c.Live.Connect(ctx, live.LiveConnectionOptions{}, opts...)
		}},
		{"fork", func(ctx context.Context, c *openai.Client, opts ...option.RequestOption) (closeableConnection, error) {
			return c.Live.Forks.Connect(ctx, "synthetic", live.ForkConnectionOptions{}, opts...)
		}},
		{"sideband", func(ctx context.Context, c *openai.Client, opts ...option.RequestOption) (closeableConnection, error) {
			return c.Live.Sideband.Connect(ctx, "synthetic", live.SidebandConnectionOptions{}, opts...)
		}},
	}
	for _, variant := range variants {
		t.Run(variant.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			var targetHits atomic.Int32
			target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				targetHits.Add(1)
				http.Error(w, "unexpected destination", http.StatusBadRequest)
			}))
			defer target.Close()
			headers := make(chan string, 4)
			closed := make(chan error, 3)
			readHeader := func() string {
				t.Helper()
				select {
				case header := <-headers:
					return header
				case <-ctx.Done():
					t.Fatal(ctx.Err())
					return ""
				}
			}
			source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				select {
				case headers <- r.Header.Get("Authorization"):
				case <-ctx.Done():
					return
				}
				if r.Header.Get("X-Redirect") == "yes" {
					http.Redirect(w, r, target.URL, http.StatusTemporaryRedirect)
					return
				}
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					closed <- err
					return
				}
				defer func() { _ = socket.CloseNow() }()
				_, _, err = socket.Read(ctx)
				if wire.CloseStatus(err) == wire.StatusNormalClosure {
					err = nil
				}
				closed <- err
			}))
			defer source.Close()
			client := openai.NewClient(option.WithBaseURL(source.URL), option.WithAPIKey("fake-client-key"))
			if _, err := variant.connect(ctx, &client, option.WithHeader("X-Redirect", "yes")); err == nil {
				t.Fatal("cross-origin redirect connected")
			}
			if got := readHeader(); got != "Bearer fake-client-key" {
				t.Fatalf("source credential = %q", got)
			}
			if targetHits.Load() != 0 {
				t.Fatal("redirect destination received a request")
			}
			for attempt := 1; attempt <= 3; attempt++ {
				key := fmt.Sprintf("fake-attempt-%d", attempt)
				opts := []option.RequestOption{option.WithAPIKey(key)}
				want := "Bearer " + key
				if attempt == 3 {
					want = "Bearer fake-explicit"
					opts = append(opts, option.WithHeader("Authorization", want))
				}
				connection, err := variant.connect(ctx, &client, opts...)
				if err != nil {
					t.Fatal(err)
				}
				if got := readHeader(); got != want {
					t.Errorf("handshake credential = %q, want %q", got, want)
				}
				if err := connection.Close(); err != nil {
					t.Fatal(err)
				}
				select {
				case err := <-closed:
					if err != nil {
						t.Fatal(err)
					}
				case <-ctx.Done():
					t.Fatal(ctx.Err())
				}
			}
		})
	}
}
