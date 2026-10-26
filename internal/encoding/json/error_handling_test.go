package json

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestAddErrorContextRecognizesWrappedTypeErrors(t *testing.T) {
	type example struct{}

	cause := &UnmarshalTypeError{Field: "value"}
	wrapped := fmt.Errorf("decode: %w", cause)
	state := decodeState{
		errorContext: &errorContext{
			Struct:     reflect.TypeOf(example{}),
			FieldStack: []string{"nested"},
		},
	}

	if got := state.addErrorContext(wrapped); !errors.Is(got, wrapped) {
		t.Fatalf("addErrorContext() = %v, want original wrapped error %v", got, wrapped)
	}
	if got, want := cause.Struct, "example"; got != want {
		t.Errorf("UnmarshalTypeError.Struct = %q, want %q", got, want)
	}
	if got, want := cause.Field, "nested.value"; got != want {
		t.Errorf("UnmarshalTypeError.Field = %q, want %q", got, want)
	}
}

func TestDecoderPreservesWrappedEOF(t *testing.T) {
	decoder := NewDecoder(&wrappedEOFReader{reader: strings.NewReader("false")})
	var value bool

	if err := decoder.Decode(&value); !errors.Is(err, io.EOF) {
		t.Errorf("Decode() error = %v, want wrapped io.EOF", err)
	}
}

type wrappedEOFReader struct {
	reader *strings.Reader
}

func (reader *wrappedEOFReader) Read(p []byte) (int, error) {
	n, err := reader.reader.Read(p)
	if errors.Is(err, io.EOF) {
		return n, fmt.Errorf("read: %w", err)
	}
	return n, err
}
