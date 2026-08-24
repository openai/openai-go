package ssestream

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
)

const maxScanTokenBytes = bufio.MaxScanTokenSize << 9

type eventLimit struct {
	reader            io.Reader
	maxBytes          int
	remainingBytes    int
	bufferedLineBytes int
	err               error
}

func newEventLimit(reader io.Reader, maxBytes int) *eventLimit {
	return &eventLimit{
		reader:         reader,
		maxBytes:       maxBytes,
		remainingBytes: maxBytes,
	}
}

func (l *eventLimit) maxScanTokenBytes() int {
	if l.maxBytes >= maxScanTokenBytes {
		return maxScanTokenBytes
	}
	return l.maxBytes + 1
}

// Read prevents Scanner from fetching beyond the current line's share of the
// event byte budget. The extra byte lets scanLines distinguish framing from
// actual overflow at the boundary.
func (l *eventLimit) Read(p []byte) (int, error) {
	maxRead := l.remainingBytes - l.bufferedLineBytes
	if maxRead < 0 {
		maxRead = 0
	}
	if maxRead < math.MaxInt {
		maxRead++
	}
	if len(p) > maxRead {
		p = p[:maxRead]
	}
	return l.reader.Read(p)
}

func (l *eventLimit) scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if newline := bytes.IndexByte(data, '\n'); newline >= 0 {
		l.bufferedLineBytes = 0
		line := data[:newline]
		if isEventDelimiter(line) {
			l.remainingBytes = l.maxBytes
			return bufio.ScanLines(data, atEOF)
		}
		if !logicalLineFits(line, l.remainingBytes) {
			return 0, nil, l.tooLargeError()
		}
		l.remainingBytes -= logicalLineBytes(line)
		return bufio.ScanLines(data, atEOF)
	}

	if atEOF {
		l.bufferedLineBytes = 0
		if isEventDelimiter(data) {
			l.remainingBytes = l.maxBytes
			return bufio.ScanLines(data, atEOF)
		}
		if len(data) > l.remainingBytes {
			return 0, nil, l.tooLargeError()
		}
		return len(data), nil, nil
	}

	if len(data) > l.remainingBytes && !(l.remainingBytes == 0 && isEventDelimiter(data)) {
		return 0, nil, l.tooLargeError()
	}
	l.bufferedLineBytes = len(data)
	return 0, nil, nil
}

func (l *eventLimit) tooLargeError() error {
	if l.err == nil {
		l.err = fmt.Errorf("%w: maximum event size is %d bytes", ErrEventTooLarge, l.maxBytes)
	}
	return l.err
}

func isEventDelimiter(line []byte) bool {
	return len(line) == 0 || len(line) == 1 && line[0] == '\r'
}

func logicalLineBytes(line []byte) int {
	if len(line) > 0 && line[len(line)-1] == '\r' {
		return len(line)
	}
	return len(line) + 1
}

// logicalLineFits counts either an omitted LF or a trailing CR as one logical
// line ending, matching bufio.ScanLines normalization and event accounting.
func logicalLineFits(line []byte, maxBytes int) bool {
	return logicalLineBytes(line) <= maxBytes
}
