package internal

import (
	"io"
	"sync"
)

// NewCloseOnceReadCloser lets callers clean up a body after a downstream
// error without double-closing it when the downstream consumer already did.
func NewCloseOnceReadCloser(body io.ReadCloser) io.ReadCloser {
	return &closeOnceReadCloser{ReadCloser: body}
}

type closeOnceReadCloser struct {
	io.ReadCloser
	once sync.Once
	err  error
}

func (b *closeOnceReadCloser) Close() error {
	b.once.Do(func() {
		b.err = b.ReadCloser.Close()
	})
	return b.err
}
