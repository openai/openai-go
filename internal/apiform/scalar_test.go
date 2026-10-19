package apiform

import (
	"bytes"
	"io"
	"mime/multipart"
	"strings"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/packages/param"
)

func TestMultipartScalarNames(t *testing.T) {
	const key = "field\r\nInjected-Header: yes\r\n\r\ninjected body"
	const escaped = "field%0D%0AInjected-Header: yes%0D%0A%0D%0Ainjected body"
	tests := []struct {
		name, arrayFormat, wantName, wantBody string
		value                                 any
	}{
		{"string map", "brackets", escaped, "contents", map[string]string{key: "contents"}},
		{"true", "brackets", escaped, "true", map[string]bool{key: true}},
		{"false", "brackets", escaped, "false", map[string]bool{key: false}},
		{"signed", "brackets", escaped, "-42", map[string]int64{key: -42}},
		{"unsigned", "brackets", escaped, "42", map[string]uint64{key: 42}},
		{"float32", "brackets", escaped, "1.25", map[string]float32{key: 1.25}},
		{"float64", "brackets", escaped, "2.5", map[string]float64{key: 2.5}},
		{"time", "brackets", escaped, "2026-01-02T03:04:05Z", map[string]time.Time{key: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)}},
		{"null", "brackets", escaped, "null", map[string]param.Opt[string]{key: param.Null[string]()}},
		{"optional", "brackets", escaped, "value", map[string]param.Opt[string]{key: param.NewOpt("value")}},
		{"comma array", "comma", escaped, "one,two", map[string][]string{key: {"one", "two"}}},
		{"bracket array", "brackets", escaped + "[]", "one", map[string][]string{key: {"one"}}},
		{"nested map", "brackets", "parent[" + escaped + "]", "value", map[string]any{"parent": map[string]string{key: "value"}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var body bytes.Buffer
			writer := multipart.NewWriter(&body)
			if err := MarshalWithSettings(test.value, writer, test.arrayFormat); err != nil {
				t.Fatal(err)
			}
			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}
			if strings.Contains(body.String(), "\r\nInjected-Header:") {
				t.Fatal("field name injected a MIME header line")
			}
			reader := multipart.NewReader(&body, writer.Boundary())
			part, err := reader.NextPart()
			if err != nil {
				t.Fatal(err)
			}
			if part.FormName() != test.wantName || len(part.Header) != 1 {
				t.Fatalf("unexpected part name or headers: %q, %v", part.FormName(), part.Header)
			}
			data, err := io.ReadAll(part)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != test.wantBody {
				t.Fatalf("body = %q, want %q", data, test.wantBody)
			}
			if _, err := reader.NextPart(); err != io.EOF {
				t.Fatalf("unexpected additional part: %v", err)
			}
		})
	}
}
