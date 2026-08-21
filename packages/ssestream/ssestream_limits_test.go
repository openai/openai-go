package ssestream

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDecoderEnforcesAggregateEventLimits(t *testing.T) {
	tests := map[string]struct {
		body    string
		options DecoderOptions
	}{
		"bytes across data lines": {
			body: strings.Repeat("data: x\n", 4) + "\n",
			options: DecoderOptions{
				MaxEventBytes: 3 * len("data: x\n"),
				MaxEventLines: 10,
			},
		},
		"lines without data": {
			body: strings.Repeat(": keep-alive\n", 3) + "\n",
			options: DecoderOptions{
				MaxEventBytes: 1024,
				MaxEventLines: 2,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			decoder := NewDecoderWithOptions(&http.Response{
				Header: http.Header{"Content-Type": []string{"text/event-stream"}},
				Body:   io.NopCloser(strings.NewReader(test.body)),
			}, test.options)

			if decoder.Next() {
				t.Fatalf("oversized event unexpectedly decoded: %+v", decoder.Event())
			}
			if !errors.Is(decoder.Err(), ErrEventTooLarge) {
				t.Fatalf("decoder error = %v, want %v", decoder.Err(), ErrEventTooLarge)
			}
		})
	}
}

func TestDecoderEnforcesByteLimitDuringLineScan(t *testing.T) {
	const maxEventBytes = 1024
	body := &testRepeatingReadCloser{remaining: maxEventBytes * 1024}
	decoder := NewDecoderWithOptions(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   body,
	}, DecoderOptions{MaxEventBytes: maxEventBytes})

	if decoder.Next() {
		t.Fatalf("oversized event unexpectedly decoded: %+v", decoder.Event())
	}
	if !errors.Is(decoder.Err(), ErrEventTooLarge) {
		t.Fatalf("decoder error = %v, want %v", decoder.Err(), ErrEventTooLarge)
	}
	if body.read > maxEventBytes+1 {
		t.Fatalf("decoder read %d bytes, want at most configured limit plus CRLF lookahead %d", body.read, maxEventBytes+1)
	}
}

func TestDecoderEnforcesRemainingByteLimitDuringLineScan(t *testing.T) {
	const (
		firstLine     = ": x\n"
		maxEventBytes = 8
	)
	body := &testDripReadCloser{data: firstLine + strings.Repeat("x", maxEventBytes)}
	decoder := NewDecoderWithOptions(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   body,
	}, DecoderOptions{MaxEventBytes: maxEventBytes})

	if decoder.Next() {
		t.Fatalf("oversized event unexpectedly decoded: %+v", decoder.Event())
	}
	if !errors.Is(decoder.Err(), ErrEventTooLarge) {
		t.Fatalf("decoder error = %v, want %v", decoder.Err(), ErrEventTooLarge)
	}
	if body.read > maxEventBytes+1 {
		t.Fatalf("decoder read %d bytes, want at most aggregate limit plus lookahead %d", body.read, maxEventBytes+1)
	}
}

func TestDecoderCapsBulkReadAtRemainingByteLimit(t *testing.T) {
	const maxEventBytes = 16
	firstLine := "data: " + strings.Repeat("x", maxEventBytes-len("data: \n")) + "\n"
	body := &testCountingReadCloser{Reader: io.MultiReader(
		strings.NewReader(firstLine),
		strings.NewReader(strings.Repeat("x", maxEventBytes)),
	)}
	decoder := NewDecoderWithOptions(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   body,
	}, DecoderOptions{MaxEventBytes: maxEventBytes})

	if decoder.Next() {
		t.Fatalf("oversized event unexpectedly decoded: %+v", decoder.Event())
	}
	if !errors.Is(decoder.Err(), ErrEventTooLarge) {
		t.Fatalf("decoder error = %v, want %v", decoder.Err(), ErrEventTooLarge)
	}
	if body.read > maxEventBytes+1 {
		t.Fatalf("decoder read %d bytes, want at most aggregate limit plus lookahead %d", body.read, maxEventBytes+1)
	}
}

func TestDecoderCapsStagedReadsAtRemainingByteLimit(t *testing.T) {
	const maxEventBytes = 16
	body := &testStagedReadCloser{chunks: [][]byte{
		[]byte("data: xx\n"),
		[]byte(strings.Repeat("x", 7)),
		[]byte(strings.Repeat("x", 8)),
	}}
	decoder := NewDecoderWithOptions(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   body,
	}, DecoderOptions{MaxEventBytes: maxEventBytes})

	if decoder.Next() {
		t.Fatalf("oversized event unexpectedly decoded: %+v", decoder.Event())
	}
	if !errors.Is(decoder.Err(), ErrEventTooLarge) {
		t.Fatalf("decoder error = %v, want %v", decoder.Err(), ErrEventTooLarge)
	}
	if body.read > maxEventBytes+1 {
		t.Fatalf("decoder read %d bytes, want at most aggregate limit plus lookahead %d", body.read, maxEventBytes+1)
	}
}

func TestDecoderHandlesIncompleteLineAtByteLimit(t *testing.T) {
	const maxEventBytes = 8
	for name, test := range map[string]struct {
		bodyLen int
		wantErr bool
	}{
		"exact limit": {bodyLen: maxEventBytes},
		"over limit":  {bodyLen: maxEventBytes + 1, wantErr: true},
	} {
		t.Run(name, func(t *testing.T) {
			decoder := NewDecoderWithOptions(&http.Response{
				Header: http.Header{"Content-Type": []string{"text/event-stream"}},
				Body:   io.NopCloser(strings.NewReader(strings.Repeat("x", test.bodyLen))),
			}, DecoderOptions{MaxEventBytes: maxEventBytes})

			if decoder.Next() {
				t.Fatalf("incomplete event unexpectedly decoded: %+v", decoder.Event())
			}
			if test.wantErr && !errors.Is(decoder.Err(), ErrEventTooLarge) {
				t.Fatalf("decoder error = %v, want %v", decoder.Err(), ErrEventTooLarge)
			}
			if !test.wantErr && decoder.Err() != nil {
				t.Fatalf("decoder error = %v, want clean EOF", decoder.Err())
			}
		})
	}
}

func TestDecoderPreservesByteLimitErrorWithConcurrentReaderError(t *testing.T) {
	const maxEventBytes = 8
	decoder := NewDecoderWithOptions(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body: &testDataErrorReadCloser{
			data: strings.Repeat("x", maxEventBytes+1),
			err:  io.ErrUnexpectedEOF,
		},
	}, DecoderOptions{MaxEventBytes: maxEventBytes})

	if decoder.Next() {
		t.Fatalf("oversized event unexpectedly decoded: %+v", decoder.Event())
	}
	if !errors.Is(decoder.Err(), ErrEventTooLarge) {
		t.Fatalf("decoder error = %v, want %v", decoder.Err(), ErrEventTooLarge)
	}
	if !errors.Is(decoder.Err(), io.ErrUnexpectedEOF) {
		t.Fatalf("decoder error = %v, want underlying %v", decoder.Err(), io.ErrUnexpectedEOF)
	}
}

func TestDecoderPreservesLineLimitErrorWithConcurrentReaderError(t *testing.T) {
	decoder := NewDecoderWithOptions(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body: &testDataErrorReadCloser{
			data: ": one\n: two\n\n",
			err:  io.ErrUnexpectedEOF,
		},
	}, DecoderOptions{MaxEventBytes: 1024, MaxEventLines: 1})

	if decoder.Next() {
		t.Fatalf("oversized event unexpectedly decoded: %+v", decoder.Event())
	}
	if !errors.Is(decoder.Err(), ErrEventTooLarge) {
		t.Fatalf("decoder error = %v, want %v", decoder.Err(), ErrEventTooLarge)
	}
	if !errors.Is(decoder.Err(), io.ErrUnexpectedEOF) {
		t.Fatalf("decoder error = %v, want underlying %v", decoder.Err(), io.ErrUnexpectedEOF)
	}
}

func TestDecoderReportsDefaultByteLimit(t *testing.T) {
	body := &testRepeatingReadCloser{remaining: DefaultMaxEventBytes + 1}
	decoder := NewDecoder(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   body,
	})

	if decoder.Next() {
		t.Fatalf("oversized event unexpectedly decoded: %+v", decoder.Event())
	}
	if !errors.Is(decoder.Err(), ErrEventTooLarge) {
		t.Fatalf("decoder error = %v, want %v", decoder.Err(), ErrEventTooLarge)
	}
}

func TestDecoderAllowsLineAtByteLimit(t *testing.T) {
	for name, body := range map[string]string{
		"LF":          "data: x\n\n",
		"CRLF":        "data: x\r\n\r\n",
		"terminal CR": "data: x\n\r",
	} {
		t.Run(name, func(t *testing.T) {
			decoder := NewDecoderWithOptions(&http.Response{
				Header: http.Header{"Content-Type": []string{"text/event-stream"}},
				Body:   io.NopCloser(strings.NewReader(body)),
			}, DecoderOptions{MaxEventBytes: len("data: x\n")})

			if !decoder.Next() {
				t.Fatalf("decoder stopped before boundary-sized event: %v", decoder.Err())
			}
			if got, want := string(decoder.Event().Data), "x\n"; got != want {
				t.Fatalf("event data = %q, want %q", got, want)
			}
		})
	}
}

func TestDecoderDispatchesEventWithTerminalCRDelimiter(t *testing.T) {
	decoder := NewDecoder(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   io.NopCloser(strings.NewReader("data: {\"value\":\"ok\"}\n\r")),
	})

	if !decoder.Next() {
		t.Fatalf("decoder stopped before event with terminal CR delimiter: %v", decoder.Err())
	}
	if got, want := string(decoder.Event().Data), "{\"value\":\"ok\"}\n"; got != want {
		t.Fatalf("event data = %q, want %q", got, want)
	}
}

func TestDecoderAllowsDripFedCRLFDelimiterAtByteLimit(t *testing.T) {
	const body = "data: x\r\n\r\n"
	decoder := NewDecoderWithOptions(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   &testDripReadCloser{data: body},
	}, DecoderOptions{MaxEventBytes: len("data: x\n")})

	if !decoder.Next() {
		t.Fatalf("decoder stopped before boundary-sized event: %v", decoder.Err())
	}
	if got, want := string(decoder.Event().Data), "x\n"; got != want {
		t.Fatalf("event data = %q, want %q", got, want)
	}
}

func TestDecoderRejectsLineOverByteLimit(t *testing.T) {
	for name, body := range map[string]string{
		"LF":   "data: xx\n\n",
		"CRLF": "data: xx\r\n\r\n",
	} {
		t.Run(name, func(t *testing.T) {
			decoder := NewDecoderWithOptions(&http.Response{
				Header: http.Header{"Content-Type": []string{"text/event-stream"}},
				Body:   io.NopCloser(strings.NewReader(body)),
			}, DecoderOptions{MaxEventBytes: len("data: x\n")})

			if decoder.Next() {
				t.Fatalf("oversized event unexpectedly decoded: %+v", decoder.Event())
			}
			if !errors.Is(decoder.Err(), ErrEventTooLarge) {
				t.Fatalf("decoder error = %v, want %v", decoder.Err(), ErrEventTooLarge)
			}
		})
	}
}

func TestDecoderOptionsAllowMultilineEventWithinLimits(t *testing.T) {
	const (
		dataBlock = "data: first\ndata: second\n"
		body      = ": one\n: two\n\n" + dataBlock + "\n"
	)
	decoder := NewDecoderWithOptions(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   io.NopCloser(strings.NewReader(body)),
	}, DecoderOptions{
		MaxEventBytes: len(dataBlock),
		MaxEventLines: 2,
	})

	if !decoder.Next() {
		t.Fatalf("decoder stopped before valid event: %v", decoder.Err())
	}
	if got, want := string(decoder.Event().Data), "first\nsecond\n"; got != want {
		t.Fatalf("event data = %q, want %q", got, want)
	}
	if decoder.Next() {
		t.Fatalf("unexpected additional event: %+v", decoder.Event())
	}
	if err := decoder.Err(); err != nil {
		t.Fatalf("decoder ended with error: %v", err)
	}
}

func TestStreamPropagatesEventLimitError(t *testing.T) {
	decoder := NewDecoderWithOptions(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body: io.NopCloser(strings.NewReader(
			"data: {\"value\":\ndata: \"too large\"}\n\n",
		)),
	}, DecoderOptions{
		MaxEventBytes: len("data: {\"value\":\n"),
		MaxEventLines: 10,
	})
	stream := NewStream[map[string]any](decoder, nil)

	if stream.Next() {
		t.Fatalf("oversized event unexpectedly decoded: %+v", stream.Current())
	}
	if !errors.Is(stream.Err(), ErrEventTooLarge) {
		t.Fatalf("stream error = %v, want %v", stream.Err(), ErrEventTooLarge)
	}
}

type testRepeatingReadCloser struct {
	remaining int
	read      int
}

func (r *testRepeatingReadCloser) Read(p []byte) (int, error) {
	if r.remaining == 0 {
		return 0, io.EOF
	}
	if len(p) > r.remaining {
		p = p[:r.remaining]
	}
	for i := range p {
		p[i] = 'x'
	}
	r.remaining -= len(p)
	r.read += len(p)
	return len(p), nil
}

func (r *testRepeatingReadCloser) Close() error { return nil }

type testDripReadCloser struct {
	data string
	read int
}

func (r *testDripReadCloser) Read(p []byte) (int, error) {
	if r.read == len(r.data) {
		return 0, io.EOF
	}
	p[0] = r.data[r.read]
	r.read++
	return 1, nil
}

func (r *testDripReadCloser) Close() error { return nil }

type testCountingReadCloser struct {
	io.Reader
	read int
}

func (r *testCountingReadCloser) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.read += n
	return n, err
}

func (r *testCountingReadCloser) Close() error { return nil }

type testStagedReadCloser struct {
	chunks [][]byte
	read   int
}

func (r *testStagedReadCloser) Read(p []byte) (int, error) {
	if len(r.chunks) == 0 {
		return 0, io.EOF
	}
	chunk := r.chunks[0]
	n := copy(p, chunk)
	r.read += n
	if n == len(chunk) {
		r.chunks = r.chunks[1:]
	} else {
		r.chunks[0] = chunk[n:]
	}
	return n, nil
}

func (r *testStagedReadCloser) Close() error { return nil }

type testDataErrorReadCloser struct {
	data string
	err  error
	done bool
}

func (r *testDataErrorReadCloser) Read(p []byte) (int, error) {
	if r.done {
		return 0, io.EOF
	}
	r.done = true
	return copy(p, r.data), r.err
}

func (r *testDataErrorReadCloser) Close() error { return nil }
