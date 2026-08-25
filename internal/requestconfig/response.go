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
	return cfg.withResponseBodyTimeout(lifecycle, func(body *responseBodyLifecycle) error {
		contents, overflow, readErr := readBodyUpTo(body, cfg.MaxErrorResponseBodyBytes)
		var compressedLimitErr *compressedResponseBodyLimitError
		if readErr != nil && !errors.As(readErr, &compressedLimitErr) {
			return readErr
		}

		// Preserve the complete diagnostic body by default. An explicit limit
		// retains only its bounded prefix, including for Error.DumpResponse.
		// ContentLength describes that retained body after truncation.
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
	return cfg.withResponseBodyTimeout(lifecycle, func(body *responseBodyLifecycle) error {
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
	body          io.ReadCloser
	stop          func()
	stopReadTimer func()
	cancelOnce    sync.Once
	reachedEOF    bool
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
	n, err := l.body.Read(p)
	if err == io.EOF {
		l.reachedEOF = true
		l.finishReading()
	}
	return n, err
}

func (l *responseBodyLifecycle) finishReading() {
	if l.stopReadTimer != nil {
		l.stopReadTimer()
	}
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
	read func(*responseBodyLifecycle) error,
) error {
	var timedOut atomic.Bool
	if cfg.ResponseBodyTimeout > 0 {
		interrupted := make(chan struct{})
		timer := time.AfterFunc(cfg.ResponseBodyTimeout, func() {
			timedOut.Store(true)
			lifecycle.interrupt(interrupted)
		})
		var stopOnce sync.Once
		lifecycle.stopReadTimer = func() {
			stopOnce.Do(func() {
				if !timer.Stop() {
					<-interrupted
				}
			})
		}
	}

	err := read(lifecycle)
	lifecycle.finishReading()
	if timedOut.Load() {
		return fmt.Errorf("response body read timed out after %s: %w", cfg.ResponseBodyTimeout, context.DeadlineExceeded)
	}
	if err != nil && !lifecycle.reachedEOF {
		// Any failed read can leave unread bytes, so HTTP/2 Close can block
		// while resetting its stream. Preserve the established failure without
		// making foreground completion wait for transport cleanup.
		lifecycle.abort()
		return err
	}
	_ = lifecycle.Close()
	return err
}

func readBodyUpTo(body *responseBodyLifecycle, limit int64) (contents []byte, overflow bool, err error) {
	defer body.finishReading()
	if limit == 0 {
		contents, err = io.ReadAll(body)
		return contents, false, err
	}
	contents, err = io.ReadAll(io.LimitReader(body, limit+1))
	if int64(len(contents)) > limit {
		return contents[:limit], true, nil
	}
	return contents, false, err
}

func decodeJSONUpTo(body *responseBodyLifecycle, limit int64, dst any) (overflow bool, err error) {
	contents, overflow, readErr := readBodyUpTo(body, limit)
	if overflow {
		return true, nil
	}
	if readErr != nil {
		return false, fmt.Errorf("error reading response body: %w", readErr)
	}
	if decodeErr := json.NewDecoder(bytes.NewReader(contents)).Decode(dst); decodeErr != nil {
		return false, fmt.Errorf("error parsing response json: %w", decodeErr)
	}
	return false, nil
}
