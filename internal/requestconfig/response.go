package requestconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/openai/openai-go/v3/internal/apierror"
	"github.com/tidwall/gjson"
)

const (
	defaultMaxResponseBodyBytes      int64 = 64 << 20
	defaultMaxErrorResponseBodyBytes int64 = 64 << 10
	defaultResponseBodyTimeout             = 10 * time.Minute
)

type responseBodyLimitError struct {
	limit      int64
	statusCode int
	cause      error
}

func (e *responseBodyLimitError) Error() string {
	if e.statusCode != 0 {
		return fmt.Sprintf("error response body for status %d exceeded configured limit of %d bytes", e.statusCode, e.limit)
	}
	return fmt.Sprintf("response body exceeded configured limit of %d bytes", e.limit)
}

func (e *responseBodyLimitError) Unwrap() error { return e.cause }

func (cfg *RequestConfig) handleErrorResponse(res *http.Response, lifecycle *responseBodyLifecycle) error {
	return cfg.withResponseBodyTimeout(lifecycle, func(body io.Reader) error {
		contents, overflow, readErr := readBodyUpTo(body, cfg.MaxErrorResponseBodyBytes)
		var compressedLimitErr *compressedResponseBodyLimitError
		if readErr != nil && !errors.As(readErr, &compressedLimitErr) {
			return readErr
		}

		// Keep only the bounded diagnostic body so Error.DumpResponse cannot make
		// another unbounded copy. ContentLength describes the retained body when
		// the server response was larger than the configured limit.
		res.Body = io.NopCloser(bytes.NewReader(contents))
		if overflow || compressedLimitErr != nil {
			res.ContentLength = int64(len(contents))
		}

		aerr := apierror.Error{Request: cfg.Request, Response: res, StatusCode: res.StatusCode}
		unwrapped := gjson.GetBytes(contents, "error").Raw
		parseErr := aerr.UnmarshalJSON([]byte(unwrapped))
		if overflow {
			return &responseBodyLimitError{
				limit:      cfg.MaxErrorResponseBodyBytes,
				statusCode: res.StatusCode,
				cause:      &aerr,
			}
		}
		if compressedLimitErr != nil {
			return &compressedResponseBodyLimitError{
				limit: compressedLimitErr.limit,
				cause: &aerr,
			}
		}
		if parseErr != nil {
			return parseErr
		}
		return &aerr
	})
}

func (cfg *RequestConfig) handleSuccessResponse(res *http.Response, lifecycle *responseBodyLifecycle) error {
	return cfg.withResponseBodyTimeout(lifecycle, func(body io.Reader) error {
		contentType := res.Header.Get("content-type")
		mediaType, _, _ := mime.ParseMediaType(contentType)
		isJSON := strings.Contains(mediaType, "application/json") || strings.HasSuffix(mediaType, "+json")

		if isJSON {
			if dst, ok := cfg.ResponseBodyInto.(*[]byte); ok {
				contents, overflow, readErr := readBodyUpTo(body, cfg.MaxResponseBodyBytes)
				if readErr != nil {
					return fmt.Errorf("error reading response body: %w", readErr)
				}
				if overflow {
					return &responseBodyLimitError{limit: cfg.MaxResponseBodyBytes}
				}
				*dst = contents
				return nil
			}

			overflow, err := decodeJSONUpTo(body, cfg.MaxResponseBodyBytes, cfg.ResponseBodyInto)
			if overflow {
				return &responseBodyLimitError{limit: cfg.MaxResponseBodyBytes}
			}
			return err
		}

		contents, overflow, readErr := readBodyUpTo(body, cfg.MaxResponseBodyBytes)
		if readErr != nil {
			return fmt.Errorf("error reading response body: %w", readErr)
		}
		if overflow {
			return &responseBodyLimitError{limit: cfg.MaxResponseBodyBytes}
		}

		switch dst := cfg.ResponseBodyInto.(type) {
		case *string:
			*dst = string(contents)
		case **string:
			tmp := string(contents)
			*dst = &tmp
		case *[]byte:
			*dst = contents
		default:
			return fmt.Errorf("expected destination type of 'string' or '[]byte' for responses with content-type '%s' that is not 'application/json'", contentType)
		}
		return nil
	})
}

type responseBodyLifecycle struct {
	body       io.ReadCloser
	stop       func()
	cancelOnce sync.Once
}

func newResponseBodyLifecycle(body io.ReadCloser, stop func()) *responseBodyLifecycle {
	return &responseBodyLifecycle{body: body, stop: stop}
}

func (l *responseBodyLifecycle) stopAttempt() {
	if l.stop != nil {
		l.cancelOnce.Do(l.stop)
	}
}

func (l *responseBodyLifecycle) Read(p []byte) (int, error) {
	return l.body.Read(p)
}

func (l *responseBodyLifecycle) Close() error {
	defer l.stopAttempt()
	return l.body.Close()
}

func (l *responseBodyLifecycle) interrupt(interrupted chan<- struct{}) {
	// Timeout must cancel first so transports whose Close does not interrupt
	// Read can observe the request ending. Signal before Close because HTTP/2
	// cleanup is allowed to block independently of the foreground timeout.
	l.stopAttempt()
	close(interrupted)
	_ = l.body.Close()
}

func (l *responseBodyLifecycle) abort() {
	l.stopAttempt()
	go func() { _ = l.body.Close() }()
}

func (cfg *RequestConfig) withResponseBodyTimeout(
	lifecycle *responseBodyLifecycle,
	read func(io.Reader) error,
) error {
	var timedOut atomic.Bool
	var timer *time.Timer
	var interrupted chan struct{}
	if cfg.ResponseBodyTimeout > 0 {
		interrupted = make(chan struct{})
		timer = time.AfterFunc(cfg.ResponseBodyTimeout, func() {
			timedOut.Store(true)
			lifecycle.interrupt(interrupted)
		})
	}

	err := read(lifecycle.body)
	if timer != nil && !timer.Stop() {
		<-interrupted
	}
	if timedOut.Load() {
		return fmt.Errorf("response body read timed out after %s: %w", cfg.ResponseBodyTimeout, context.DeadlineExceeded)
	}
	var bodyLimitErr *responseBodyLimitError
	var compressedLimitErr *compressedResponseBodyLimitError
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) ||
		errors.As(err, &bodyLimitErr) || errors.As(err, &compressedLimitErr) {
		// Aborted reads can leave unread bytes, so HTTP/2 Close can block while
		// resetting its stream. Return the established failure immediately.
		lifecycle.abort()
		return err
	}
	_ = lifecycle.Close()
	return err
}

func readBodyUpTo(body io.Reader, limit int64) (contents []byte, overflow bool, err error) {
	contents, err = io.ReadAll(io.LimitReader(body, limit+1))
	if int64(len(contents)) > limit {
		return contents[:limit], true, nil
	}
	return contents, false, err
}

func decodeJSONUpTo(body io.Reader, limit int64, dst any) (overflow bool, err error) {
	limited := &io.LimitedReader{R: body, N: limit + 1}
	decodeErr := json.NewDecoder(limited).Decode(dst)
	_, readErr := io.Copy(io.Discard, limited)
	if limited.N == 0 {
		return true, nil
	}
	if readErr != nil {
		return false, fmt.Errorf("error reading response body: %w", readErr)
	}
	if decodeErr != nil {
		return false, fmt.Errorf("error parsing response json: %w", decodeErr)
	}
	return false, nil
}
