package openai_test

import (
	"io"
	"mime"
	"runtime"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
)

func TestMultipartInferredFilenameOmitsDirectories(t *testing.T) {
	tests := []struct{ name, want string }{
		{"/home/user/private/report.txt", "report.txt"},
		{"report.txt", "report.txt"},
	}
	if runtime.GOOS == "windows" {
		tests = append(tests,
			struct{ name, want string }{`C:\Users\user\private\report.txt`, "report.txt"},
			struct{ name, want string }{`\\server\private\report.txt`, "report.txt"},
		)
	} else {
		tests = append(tests,
			struct{ name, want string }{`/home/user/private/customer\report.txt`, `customer\report.txt`},
			struct{ name, want string }{`C:\Users\user\private\report.txt`, `C:\Users\user\private\report.txt`},
		)
	}
	for _, test := range tests {
		name := test.name
		for _, wrapped := range []bool{false, true} {
			t.Run(name+map[bool]string{false: "/direct", true: "/wrapped"}[wrapped], func(t *testing.T) {
				var reader io.Reader = multipartNamedReader{Reader: strings.NewReader("contents"), name: name}
				if wrapped {
					reader = openai.File(reader, "", "text/plain")
				}
				part, _ := encodeMultipartFile(t, reader)
				_, params, err := mime.ParseMediaType(part.header.Get("Content-Disposition"))
				if err != nil {
					t.Fatal(err)
				}
				if params["filename"] != test.want || part.body != "contents" {
					t.Fatalf("filename = %q, body = %q", params["filename"], part.body)
				}
			})
		}
	}
}

func TestMultipartExplicitFilenamePreservesDirectories(t *testing.T) {
	const name = `chosen/path\report.txt`
	for _, reader := range []io.Reader{
		openai.File(strings.NewReader("contents"), name, "text/plain"),
		multipartFilenameReader{Reader: strings.NewReader("contents"), filename: name},
	} {
		part, _ := encodeMultipartFile(t, reader)
		_, params, err := mime.ParseMediaType(part.header.Get("Content-Disposition"))
		if err != nil {
			t.Fatal(err)
		}
		if params["filename"] != name || part.body != "contents" {
			t.Fatalf("filename = %q, body = %q", params["filename"], part.body)
		}
	}
}
