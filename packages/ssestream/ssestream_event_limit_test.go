package ssestream

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3/internal/requestconfig"
)

func TestDecoderEnforcesExplicitAggregateEventLimit(t *testing.T) {
	const body = "data: first\ndata: second\n\n"
	for name, test := range map[string]struct {
		maxBytes int
		wantErr  bool
	}{
		"exact limit": {
			maxBytes: len("data: first\ndata: second\n"),
		},
		"overflow": {
			maxBytes: len("data: first\ndata: second\n") - 1,
			wantErr:  true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			decoder := newLimitedDecoder(t, io.NopCloser(strings.NewReader(body)), test.maxBytes)

			if test.wantErr {
				if decoder.Next() {
					t.Fatalf("oversized event unexpectedly decoded: %+v", decoder.Event())
				}
				if !errors.Is(decoder.Err(), ErrEventTooLarge) {
					t.Fatalf("decoder error = %v, want ErrEventTooLarge", decoder.Err())
				}
				return
			}

			if !decoder.Next() {
				t.Fatalf("decoder stopped before boundary-sized event: %v", decoder.Err())
			}
			if got, want := string(decoder.Event().Data), "first\nsecond\n"; got != want {
				t.Fatalf("event data = %q, want %q", got, want)
			}
		})
	}
}

func TestDecoderExplicitLimitCapsReadAhead(t *testing.T) {
	const (
		firstLine = ": x\n"
		maxBytes  = 8
	)
	body := &countingDripReadCloser{data: firstLine + strings.Repeat("x", maxBytes)}
	decoder := newLimitedDecoder(t, body, maxBytes)

	if decoder.Next() {
		t.Fatalf("oversized event unexpectedly decoded: %+v", decoder.Event())
	}
	if !errors.Is(decoder.Err(), ErrEventTooLarge) {
		t.Fatalf("decoder error = %v, want ErrEventTooLarge", decoder.Err())
	}
	if body.read > maxBytes+1 {
		t.Fatalf("decoder read %d bytes, want at most aggregate limit plus lookahead %d", body.read, maxBytes+1)
	}
}

func TestDecoderExplicitLimitCapsBulkReadAtRemainingBudget(t *testing.T) {
	const maxBytes = 16
	firstLine := "data: " + strings.Repeat("x", maxBytes-len("data: \n")) + "\n"
	body := &countingReadCloser{Reader: io.MultiReader(
		strings.NewReader(firstLine),
		strings.NewReader(strings.Repeat("x", maxBytes)),
	)}
	decoder := newLimitedDecoder(t, body, maxBytes)

	if decoder.Next() {
		t.Fatalf("oversized event unexpectedly decoded: %+v", decoder.Event())
	}
	if !errors.Is(decoder.Err(), ErrEventTooLarge) {
		t.Fatalf("decoder error = %v, want ErrEventTooLarge", decoder.Err())
	}
	if body.read > maxBytes+1 {
		t.Fatalf("decoder read %d bytes, want at most aggregate limit plus lookahead %d", body.read, maxBytes+1)
	}
}

func TestDecoderExplicitLimitPreservesIncompleteEventEOF(t *testing.T) {
	const maxBytes = 8
	for name, test := range map[string]struct {
		bodyLen int
		wantErr bool
	}{
		"exact limit is clean EOF": {bodyLen: maxBytes},
		"overflow is an error":     {bodyLen: maxBytes + 1, wantErr: true},
	} {
		t.Run(name, func(t *testing.T) {
			decoder := newLimitedDecoder(t,
				io.NopCloser(strings.NewReader(strings.Repeat("x", test.bodyLen))),
				maxBytes,
			)

			if decoder.Next() {
				t.Fatalf("incomplete event unexpectedly decoded: %+v", decoder.Event())
			}
			if test.wantErr {
				if !errors.Is(decoder.Err(), ErrEventTooLarge) {
					t.Fatalf("decoder error = %v, want ErrEventTooLarge", decoder.Err())
				}
			} else if decoder.Err() != nil {
				t.Fatalf("decoder error = %v, want clean EOF", decoder.Err())
			}
		})
	}
}

func TestDecoderExplicitLimitPreservesFramingAtBoundary(t *testing.T) {
	const eventBytes = len("data: x\n")
	for name, body := range map[string]string{
		"LF":          "data: x\n\n",
		"CRLF":        "data: x\r\n\r\n",
		"terminal CR": "data: x\n\r",
	} {
		t.Run(name, func(t *testing.T) {
			decoder := newLimitedDecoder(t,
				&countingDripReadCloser{data: body},
				eventBytes,
			)

			if !decoder.Next() {
				t.Fatalf("decoder stopped before boundary-sized event: %v", decoder.Err())
			}
			if got, want := string(decoder.Event().Data), "x\n"; got != want {
				t.Fatalf("event data = %q, want %q", got, want)
			}
		})
	}
}

func TestDecoderExplicitLimitResetsBetweenEvents(t *testing.T) {
	const (
		line     = "data: x\n"
		maxBytes = len(line)
	)
	decoder := newLimitedDecoder(t, io.NopCloser(strings.NewReader(line+"\n"+line+"\n")), maxBytes)

	for event := 1; event <= 2; event++ {
		if !decoder.Next() {
			t.Fatalf("decoder stopped before exact-budget event %d: %v", event, decoder.Err())
		}
		if got, want := string(decoder.Event().Data), "x\n"; got != want {
			t.Fatalf("event %d data = %q, want %q", event, got, want)
		}
	}
}

func TestDecoderExplicitLimitResetsAfterIgnoredBlock(t *testing.T) {
	const (
		ignored  = "event: ignored\n"
		data     = "data: x\n"
		maxBytes = len(ignored)
	)
	decoder := newLimitedDecoder(t, io.NopCloser(strings.NewReader(ignored+"\n"+data+"\n")), maxBytes)

	if !decoder.Next() {
		t.Fatalf("decoder stopped before data following ignored block: %v", decoder.Err())
	}
	if got, want := string(decoder.Event().Data), "x\n"; got != want {
		t.Fatalf("event data = %q, want %q", got, want)
	}
}

func TestDecoderExplicitLimitPreservesLegacyLineLimit(t *testing.T) {
	line := strings.Repeat("x", maxScanTokenBytes)
	defaultDecoder := newLimitedDecoder(t, io.NopCloser(strings.NewReader(line)), 0)
	limitedDecoder := newLimitedDecoder(t, io.NopCloser(strings.NewReader(line)), maxScanTokenBytes*2)

	if defaultDecoder.Next() || defaultDecoder.Err() == nil {
		t.Fatalf("default decoder error = %v, want oversized-token error", defaultDecoder.Err())
	}
	if limitedDecoder.Next() || limitedDecoder.Err() == nil {
		t.Fatalf("limited decoder error = %v, want oversized-token error", limitedDecoder.Err())
	}
	if got, want := limitedDecoder.Err().Error(), defaultDecoder.Err().Error(); got != want {
		t.Fatalf("limited decoder error = %q, want legacy error %q", got, want)
	}
}

func TestDecoderExplicitLimitPreservesReaderError(t *testing.T) {
	const maxBytes = 8
	decoder := newLimitedDecoder(t, &dataErrorReadCloser{
		data: strings.Repeat("x", maxBytes+1),
		err:  io.ErrUnexpectedEOF,
	}, maxBytes)

	if decoder.Next() {
		t.Fatalf("oversized event unexpectedly decoded: %+v", decoder.Event())
	}
	if !errors.Is(decoder.Err(), ErrEventTooLarge) {
		t.Fatalf("decoder error = %v, want ErrEventTooLarge", decoder.Err())
	}
	if !errors.Is(decoder.Err(), io.ErrUnexpectedEOF) {
		t.Fatalf("decoder error = %v, want underlying %v", decoder.Err(), io.ErrUnexpectedEOF)
	}
}

func TestRegisteredDecoderIgnoresExplicitEventLimit(t *testing.T) {
	const contentType = "application/x-openai-go-event-limit-test"
	want := &testDecoder{}
	RegisterDecoder(contentType, func(io.ReadCloser) Decoder { return want })
	t.Cleanup(func() { delete(decoderTypes, contentType) })

	response := limitedResponse(io.NopCloser(strings.NewReader("")), 1)
	t.Cleanup(func() { _ = response.Body.Close() })
	response.Header.Set("Content-Type", contentType)
	if got := NewDecoder(response); got != want {
		t.Fatalf("decoder = %T, want registered decoder", got)
	}
}

func newLimitedDecoder(t *testing.T, body io.ReadCloser, maxBytes int) Decoder {
	t.Helper()
	response := limitedResponse(body, maxBytes)
	t.Cleanup(func() { _ = response.Body.Close() })
	return NewDecoder(response)
}

func limitedResponse(body io.ReadCloser, maxBytes int) *http.Response {
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		panic(err)
	}
	req = requestconfig.WithSSEMaxEventBytes(req, maxBytes)
	return &http.Response{
		Header:  http.Header{"Content-Type": {"text/event-stream"}},
		Body:    body,
		Request: req,
	}
}

type countingDripReadCloser struct {
	data string
	read int
}

func (r *countingDripReadCloser) Read(p []byte) (int, error) {
	if r.read == len(r.data) {
		return 0, io.EOF
	}
	p[0] = r.data[r.read]
	r.read++
	return 1, nil
}

func (r *countingDripReadCloser) Close() error { return nil }

type countingReadCloser struct {
	io.Reader
	read int
}

func (r *countingReadCloser) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.read += n
	return n, err
}

func (r *countingReadCloser) Close() error { return nil }

type dataErrorReadCloser struct {
	data string
	err  error
	done bool
}

func (r *dataErrorReadCloser) Read(p []byte) (int, error) {
	if r.done {
		return 0, io.EOF
	}
	r.done = true
	return copy(p, r.data), r.err
}

func (r *dataErrorReadCloser) Close() error { return nil }
