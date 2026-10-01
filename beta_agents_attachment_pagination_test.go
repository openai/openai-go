package openai_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestBetaAgentAttachmentRejectsCyclicPagination(t *testing.T) {
	for _, endpoint := range []string{"turns", "items"} {
		for _, cycle := range []struct {
			name, more string
			ids        []string
		}{
			{"repeated", `,"has_more":true`, []string{"a", "a"}},
			{"cycle", `,"has_more":true`, []string{"a", "b", "a"}},
			{"implicit continuation", "", []string{"a", "a"}},
		} {
			t.Run(endpoint+"/"+cycle.name, func(t *testing.T) {
				var pages atomic.Int32
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					if r.Header.Get("X-Pagination-Test") != "kept" {
						t.Error("pagination lost request options")
					}
					switch {
					case strings.HasSuffix(r.URL.Path, "/"+endpoint):
						page := int(pages.Add(1)) - 1
						if page >= len(cycle.ids) {
							// Keep the vulnerable implementation's regression test bounded.
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						if page > 0 && r.URL.Query().Get("after") != cycle.ids[page-1] {
							t.Error("wrong pagination cursor")
						}
						if endpoint == "turns" {
							_, _ = fmt.Fprintf(w, `{"data":[{"id":%q,"subagent_id":"child"}]%s}`, cycle.ids[page], cycle.more)
						} else {
							_, _ = fmt.Fprintf(w, `{"data":[{"id":%q,"type":"function_call","turn_id":"root"}]%s}`, cycle.ids[page], cycle.more)
						}
					case strings.HasSuffix(r.URL.Path, "/turns"):
						_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
					case strings.HasSuffix(r.URL.Path, "/turns/root"):
						_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"completed"}`)
					case strings.HasSuffix(r.URL.Path, "/events"):
						w.Header().Set("Content-Type", "text/event-stream")
						w.WriteHeader(http.StatusOK)
						w.(http.Flusher).Flush()
						<-r.Context().Done()
					default:
						_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
					}
				}))
				defer server.Close()
				client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
				stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{}, option.WithHeader("X-Pagination-Test", "kept"))
				defer func() { _ = stream.Close() }()
				_, err := stream.FinalResult()
				var resultError *openai.BetaAgentTurnResultError
				if !errors.As(err, &resultError) || resultError.Cause == nil || resultError.Cause.Error() != "attachment pagination cursor did not advance" {
					t.Fatalf("expected cursor failure, got %v (cause=%v)", err, errors.Unwrap(err))
				}
				if pages.Load() != int32(len(cycle.ids)) {
					t.Fatalf("followed a repeated cursor: pages=%d", pages.Load())
				}
			})
		}
	}
}
