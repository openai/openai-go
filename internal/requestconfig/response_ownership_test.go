package requestconfig_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/realtime"
)

type responseTransport func(*http.Request) (*http.Response, error)

func (f responseTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func (f responseTransport) Do(r *http.Request) (*http.Response, error) { return f(r) }

type ownedResponseBody struct {
	io.Reader
	closed bool
}

func (b *ownedResponseBody) Close() error {
	b.closed = true
	return nil
}

type unreadResponseBody struct{ t *testing.T }

func (b unreadResponseBody) Read([]byte) (int, error) {
	b.t.Error("no-result operation read an unowned body")
	return 0, errors.New("body must not be read")
}

func TestNoResultOperationsCloseUnownedResponses(t *testing.T) {
	operations := map[string]func(openai.Client, context.Context) error{
		"accept": func(c openai.Client, ctx context.Context) error {
			return c.Realtime.Calls.Accept(ctx, "call", realtime.CallAcceptParams{})
		},
		"hangup": func(c openai.Client, ctx context.Context) error { return c.Realtime.Calls.Hangup(ctx, "call") },
		"refer": func(c openai.Client, ctx context.Context) error {
			return c.Realtime.Calls.Refer(ctx, "call", realtime.CallReferParams{})
		},
		"reject": func(c openai.Client, ctx context.Context) error {
			return c.Realtime.Calls.Reject(ctx, "call", realtime.CallRejectParams{})
		},
		"beta response delete": func(c openai.Client, ctx context.Context) error {
			return c.Beta.Responses.Delete(ctx, "response", openai.BetaResponseDeleteParams{})
		},
		"container delete": func(c openai.Client, ctx context.Context) error { return c.Containers.Delete(ctx, "container") },
		"container file delete": func(c openai.Client, ctx context.Context) error {
			return c.Containers.Files.Delete(ctx, "container", "file")
		},
		"response delete": func(c openai.Client, ctx context.Context) error { return c.Responses.Delete(ctx, "response") },
	}
	for name, operation := range operations {
		t.Run(name, func(t *testing.T) {
			for _, timeout := range []time.Duration{0, time.Hour} {
				t.Run(timeout.String(), func(t *testing.T) {
					body := &ownedResponseBody{Reader: unreadResponseBody{t}}
					var requestContext context.Context
					client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL("https://example.invalid"),
						option.WithMaxRetries(0), option.WithRequestTimeout(timeout),
						option.WithHTTPClient(&http.Client{Transport: responseTransport(func(r *http.Request) (*http.Response, error) {
							requestContext = r.Context()
							return &http.Response{StatusCode: http.StatusOK, Body: body, Header: make(http.Header)}, nil
						})}))
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					if err := operation(client, ctx); err != nil {
						t.Fatal(err)
					}
					if !body.closed {
						t.Error("successful operation left its response body open")
					}
					if timeout != 0 && !errors.Is(requestContext.Err(), context.Canceled) {
						t.Error("successful operation retained its timeout context")
					}
					if ctx.Err() != nil {
						t.Error("operation canceled the caller's context")
					}
				})
			}
		})
	}
}

func TestExplicitResponseOwnership(t *testing.T) {
	for _, mode := range []string{"response option", "raw destination"} {
		t.Run(mode, func(t *testing.T) {
			body := &ownedResponseBody{Reader: strings.NewReader("response payload")}
			var requestContext context.Context
			client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL("https://example.invalid"),
				option.WithRequestTimeout(time.Hour),
				option.WithHTTPClient(&http.Client{Transport: responseTransport(func(r *http.Request) (*http.Response, error) {
					requestContext = r.Context()
					return &http.Response{StatusCode: http.StatusOK, Body: body, Header: make(http.Header)}, nil
				})}))
			var response *http.Response
			var err error
			if mode == "response option" {
				err = client.Containers.Delete(context.Background(), "container", option.WithResponseInto(&response))
			} else {
				err = client.Get(context.Background(), "raw", nil, &response)
			}
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = response.Body.Close() }()
			if body.closed || requestContext.Err() != nil {
				t.Fatal("SDK closed or canceled a caller-owned response")
			}
			contents, err := io.ReadAll(response.Body)
			if err != nil || string(contents) != "response payload" {
				t.Fatalf("caller read %q, error %v", contents, err)
			}
			if err := response.Body.Close(); err != nil {
				t.Fatal(err)
			}
			if !body.closed || !errors.Is(requestContext.Err(), context.Canceled) {
				t.Error("caller Close did not release body and timeout")
			}
		})
	}
}

func TestNoResultOperationWithBodylessCustomResponse(t *testing.T) {
	var requestContext context.Context
	client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL("https://example.invalid"),
		option.WithRequestTimeout(time.Hour),
		option.WithHTTPClient(responseTransport(func(r *http.Request) (*http.Response, error) {
			requestContext = r.Context()
			return &http.Response{StatusCode: http.StatusNoContent}, nil
		})))
	if err := client.Containers.Delete(context.Background(), "container"); err != nil {
		t.Fatal(err)
	}
	if !errors.Is(requestContext.Err(), context.Canceled) {
		t.Error("successful bodyless operation retained its timeout context")
	}
}
