package openai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestBetaResponseConcurrentRequestOptions(t *testing.T) {
	for _, streaming := range []bool{false, true} {
		t.Run(fmt.Sprintf("streaming=%t", streaming), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var body struct {
					Instructions string `json:"instructions"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				want := []string{"caller-beta", body.Instructions}
				if got := r.Header.Values("OpenAI-Beta"); !slices.Equal(got, want) {
					t.Errorf("request %s: beta headers = %q, want %q", body.Instructions, got, want)
				}
				if got := r.Header.Get("X-Caller"); got != "preserved" {
					t.Errorf("caller header = %q, want preserved", got)
				}
				if streaming {
					w.Header().Set("Content-Type", "text/event-stream")
					if _, err := fmt.Fprint(w, "data: [DONE]\n\n"); err != nil {
						t.Error(err)
					}
				} else {
					w.Header().Set("Content-Type", "application/json")
					if _, err := fmt.Fprint(w, `{}`); err != nil {
						t.Error(err)
					}
				}
			}))
			defer server.Close()
			service := openai.NewBetaResponseService(
				option.WithBaseURL(server.URL), option.WithUnsafeAllowHTTP(),
				option.WithAPIKey("test-key"), option.WithMaxRetries(0),
			)
			// The unused slots expose writes through the caller's backing array,
			// even when scheduling happens not to produce a header mismatch.
			opts := make([]option.RequestOption, 2, 8)
			opts[0] = option.WithHeader("X-Caller", "preserved")
			opts[1] = option.WithHeaderAdd("OpenAI-Beta", "caller-beta")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			start := make(chan struct{})
			var ready, done sync.WaitGroup
			for i := range 32 {
				ready.Add(1)
				done.Add(1)
				go func() {
					defer done.Done()
					ready.Done()
					<-start
					beta := fmt.Sprintf("request-%d", i)
					params := openai.BetaResponseNewParams{Instructions: openai.String(beta), Betas: []string{beta}}
					if streaming {
						stream := service.NewStreaming(ctx, params, opts...)
						defer func() {
							if err := stream.Close(); err != nil {
								t.Error(err)
							}
						}()
						for stream.Next() {
						}
						if err := stream.Err(); err != nil {
							t.Error(err)
						}
					} else if _, err := service.New(ctx, params, opts...); err != nil {
						t.Error(err)
					}
				}()
			}
			ready.Wait()
			close(start)
			done.Wait()
			for i, opt := range opts[len(opts):cap(opts)] {
				if opt != nil {
					t.Errorf("request mutated caller option backing array at slot %d", len(opts)+i)
				}
			}
		})
	}
}
