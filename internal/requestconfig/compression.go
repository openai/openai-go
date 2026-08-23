package requestconfig

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

type compressedResponseBodyLimitError struct {
	limit int64
	cause error
}

func (e *compressedResponseBodyLimitError) Error() string {
	return fmt.Sprintf("compressed response body exceeded configured limit of %d bytes", e.limit)
}

func (e *compressedResponseBodyLimitError) Unwrap() error { return e.cause }

type wireLimitReader struct {
	reader    io.Reader
	remaining int64
	limit     int64
}

func (r *wireLimitReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if r.remaining == 0 {
		var extra [1]byte
		n, err := r.reader.Read(extra[:])
		if n != 0 {
			return 0, &compressedResponseBodyLimitError{limit: r.limit}
		}
		return 0, err
	}
	if int64(len(p)) > r.remaining {
		p = p[:r.remaining]
	}
	n, err := r.reader.Read(p)
	r.remaining -= int64(n)
	return n, err
}

type gzipResponseBody struct {
	body       io.ReadCloser
	compressed io.Reader
	decoded    io.Reader
	initOnce   sync.Once
	initErr    error
}

func newGzipResponseBody(body io.ReadCloser, wireLimit int64) *gzipResponseBody {
	var reader io.Reader = body
	if wireLimit > 0 {
		reader = &wireLimitReader{reader: body, remaining: wireLimit, limit: wireLimit}
	}
	return &gzipResponseBody{body: body, compressed: reader}
}

func (b *gzipResponseBody) Read(p []byte) (int, error) {
	b.initOnce.Do(func() {
		b.decoded, b.initErr = gzip.NewReader(b.compressed)
	})
	if b.initErr != nil {
		return 0, b.initErr
	}
	return b.decoded.Read(p)
}

func (b *gzipResponseBody) Close() error {
	// gzip.Reader.Close does not close its underlying reader. Closing the wire
	// body directly preserves net/http's guarantee that Close can unblock Read.
	return b.body.Close()
}

func (cfg *RequestConfig) ownsSuccessResponseBody() bool {
	if cfg.ResponseBodyInto == nil {
		return false
	}
	_, raw := cfg.ResponseBodyInto.(**http.Response)
	return !raw
}

func (cfg *RequestConfig) shouldManageGzip(req *http.Request) bool {
	if req.Method == http.MethodHead ||
		req.Header.Get("Accept-Encoding") != "" ||
		req.Header.Get("Range") != "" {
		return false
	}

	// Respect explicit stdlib transport configuration, including SDK-owned
	// wrappers that preserve the selected native transport's policy.
	if cfg.CustomHTTPDoer == nil && cfg.HTTPClient != nil {
		transport := cfg.HTTPClient.Transport
		if transport == nil {
			transport = http.DefaultTransport
		}
		switch transport := transport.(type) {
		case *http.Transport:
			return !transport.DisableCompression
		case interface{ CompressionDisabled() bool }:
			return !transport.CompressionDisabled()
		}
	}
	return true
}

func (cfg *RequestConfig) withManagedGzip(next middlewareNext) middlewareNext {
	return func(req *http.Request) (*http.Response, error) {
		manageGzip := cfg.shouldManageGzip(req)
		if manageGzip {
			// Decorate a per-attempt clone so authentication or other outer
			// middleware can replay the pristine request through this layer.
			req = req.Clone(req.Context())
			req.Header.Set("Accept-Encoding", "gzip")
		}

		res, err := next(req)
		if err != nil || res == nil || res.Body == nil || !manageGzip ||
			!strings.EqualFold(res.Header.Get("Content-Encoding"), "gzip") {
			return res, err
		}

		var wireLimit int64
		switch {
		case res.StatusCode >= http.StatusBadRequest:
			wireLimit = cfg.MaxErrorResponseBodyBytes
		case cfg.ownsSuccessResponseBody():
			wireLimit = cfg.MaxResponseBodyBytes
		}
		res.Body = newGzipResponseBody(res.Body, wireLimit)
		res.Header.Del("Content-Encoding")
		res.Header.Del("Content-Length")
		res.ContentLength = -1
		res.Uncompressed = true
		return res, nil
	}
}
