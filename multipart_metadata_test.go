package openai_test

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
)

type multipartNamedReader struct {
	io.Reader
	name string
}

func (r multipartNamedReader) Name() string { return r.name }

type multipartFilenameReader struct {
	io.Reader
	filename string
}

func (r multipartFilenameReader) Filename() string { return r.filename }

type multipartFilenameAndNameReader struct {
	io.Reader
	filename string
	name     string
}

func (r multipartFilenameAndNameReader) Filename() string { return r.filename }
func (r multipartFilenameAndNameReader) Name() string     { return r.name }

type multipartContentTypeReader struct {
	io.Reader
	contentType string
}

func (r multipartContentTypeReader) ContentType() string { return r.contentType }

type encodedMultipartFile struct {
	header      textproto.MIMEHeader
	filename    string
	contentType string
	body        string
}

func encodeMultipartFile(t *testing.T, file io.Reader) (encodedMultipartFile, string) {
	t.Helper()

	data, contentType, err := (openai.FileNewParams{
		File:    file,
		Purpose: openai.FilePurposeAssistants,
	}).MarshalMultipart()
	if err != nil {
		t.Fatalf("MarshalMultipart() error = %v", err)
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		t.Fatalf("parse multipart content type: %v", err)
	}
	if mediaType != "multipart/form-data" {
		t.Fatalf("media type = %q, want multipart/form-data", mediaType)
	}

	reader := multipart.NewReader(bytes.NewReader(data), params["boundary"])
	var encoded *encodedMultipartFile
	for {
		part, nextErr := reader.NextPart()
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			t.Fatalf("read multipart part: %v", nextErr)
		}

		body, readErr := io.ReadAll(part)
		if readErr != nil {
			t.Fatalf("read multipart body: %v", readErr)
		}
		if part.FormName() == "file" {
			if encoded != nil {
				t.Fatal("multipart body contains more than one file part")
			}
			encoded = &encodedMultipartFile{
				header:      part.Header,
				filename:    part.FileName(),
				contentType: part.Header.Get("Content-Type"),
				body:        string(body),
			}
		}
		if closeErr := part.Close(); closeErr != nil {
			t.Fatalf("close multipart part: %v", closeErr)
		}
	}
	if encoded == nil {
		t.Fatal("multipart body does not contain a file part")
	}
	return *encoded, string(data)
}

func TestMultipartFilenameMetadata(t *testing.T) {
	tests := []struct {
		name         string
		file         io.Reader
		wantFilename string
		wantRaw      string
	}{
		{
			name:         "openai.File carriage return",
			file:         openai.File(strings.NewReader("contents"), "report\r.csv", "text/csv"),
			wantFilename: "report%0D.csv",
			wantRaw:      `filename="report%0D.csv"`,
		},
		{
			name: "Name line feed",
			file: multipartNamedReader{
				Reader: strings.NewReader("contents"),
				name:   "directory/report\n.csv",
			},
			wantFilename: "report%0A.csv",
			wantRaw:      `filename="report%0A.csv"`,
		},
		{
			name: "Filename CRLF",
			file: multipartFilenameReader{
				Reader:   strings.NewReader("contents"),
				filename: "report\r\n.csv",
			},
			wantFilename: "report%0D%0A.csv",
			wantRaw:      `filename="report%0D%0A.csv"`,
		},
		{
			name: "blank line and header",
			file: multipartFilenameReader{
				Reader: strings.NewReader("contents"),
				filename: "report.csv\r\nInjected-Header: yes\r\n\r\n" +
					"injected body",
			},
			wantFilename: "report.csv%0D%0AInjected-Header: yes%0D%0A%0D%0Ainjected body",
			wantRaw: `filename="report.csv%0D%0AInjected-Header: yes` +
				`%0D%0A%0D%0Ainjected body"`,
		},
		{
			name: "quote",
			file: multipartFilenameReader{
				Reader:   strings.NewReader("contents"),
				filename: `quo"te.txt`,
			},
			wantFilename: `quo"te.txt`,
			wantRaw:      `filename="quo\"te.txt"`,
		},
		{
			name: "backslash",
			file: multipartFilenameReader{
				Reader:   strings.NewReader("contents"),
				filename: `back\slash.txt`,
			},
			wantFilename: `back\slash.txt`,
			wantRaw:      `filename="back\\slash.txt"`,
		},
		{
			name: "Filename takes precedence over Name",
			file: multipartFilenameAndNameReader{
				Reader:   strings.NewReader("contents"),
				filename: "from-filename.txt",
				name:     "from-name.txt",
			},
			wantFilename: "from-filename.txt",
			wantRaw:      `filename="from-filename.txt"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			part, raw := encodeMultipartFile(t, test.file)
			if part.filename != test.wantFilename {
				t.Errorf("filename = %q, want %q", part.filename, test.wantFilename)
			}
			if !strings.Contains(raw, test.wantRaw) {
				t.Errorf("multipart body does not contain %q", test.wantRaw)
			}
			if part.header.Get("Injected-Header") != "" {
				t.Errorf("injected header = %q, want empty", part.header.Get("Injected-Header"))
			}
			if strings.Contains(raw, "\r\nInjected-Header:") {
				t.Error("multipart body contains an injected header line")
			}
			if part.body != "contents" {
				t.Errorf("file body = %q, want contents", part.body)
			}
		})
	}
}

func TestMultipartFilenameRejectsOtherControlCharacters(t *testing.T) {
	tests := []string{
		"nul\x00.txt",
		"unit-separator\x1f.txt",
		"delete\x7f.txt",
		"next-line\u0085.txt",
	}

	for _, filename := range tests {
		t.Run(filename, func(t *testing.T) {
			data, contentType, err := (openai.FileNewParams{
				File:    openai.File(strings.NewReader("contents"), filename, "text/plain"),
				Purpose: openai.FilePurposeAssistants,
			}).MarshalMultipart()
			if err == nil {
				t.Fatal("MarshalMultipart() error = nil, want invalid filename error")
			}
			if !strings.Contains(err.Error(), "invalid multipart filename") {
				t.Errorf("MarshalMultipart() error = %q, want invalid filename", err)
			}
			if data != nil || contentType != "" {
				t.Errorf("MarshalMultipart() = (%d bytes, %q), want (nil, empty)", len(data), contentType)
			}
		})
	}
}

func TestMultipartContentTypeMetadata(t *testing.T) {
	tests := []struct {
		name        string
		file        io.Reader
		contentType string
	}{
		{
			name: "openai.File vendor type",
			file: openai.File(
				strings.NewReader("contents"),
				"data.json",
				`application/vnd.openai.dataset+json; version="1"`,
			),
			contentType: `application/vnd.openai.dataset+json; version="1"`,
		},
		{
			name: "custom ContentType implementation",
			file: multipartContentTypeReader{
				Reader:      strings.NewReader("contents"),
				contentType: "Text/Plain; charset=utf-8",
			},
			contentType: "Text/Plain; charset=utf-8",
		},
		{
			name: "horizontal tab whitespace",
			file: multipartContentTypeReader{
				Reader:      strings.NewReader("contents"),
				contentType: "text/plain;\tcharset=utf-8",
			},
			contentType: "text/plain;\tcharset=utf-8",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			part, _ := encodeMultipartFile(t, test.file)
			if part.contentType != test.contentType {
				t.Errorf("Content-Type = %q, want %q", part.contentType, test.contentType)
			}
		})
	}
}

func TestMultipartContentTypeRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
	}{
		{name: "carriage return", contentType: "text/plain\rInjected: yes"},
		{name: "line feed", contentType: "text/plain\nInjected: yes"},
		{name: "CRLF", contentType: "text/plain\r\nInjected: yes"},
		{name: "blank line", contentType: "text/plain\r\nInjected: yes\r\n\r\nbody"},
		{name: "nul", contentType: "text/plain\x00"},
		{name: "delete", contentType: "text/plain\x7f"},
		{name: "unicode control", contentType: "text/plain\u0085"},
		{name: "invalid tab placement", contentType: "text/\tplain"},
		{name: "missing slash", contentType: "text"},
		{name: "missing subtype", contentType: "text/"},
		{name: "invalid parameter", contentType: "text/plain; charset"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, contentType, err := (openai.FileNewParams{
				File: multipartContentTypeReader{
					Reader:      strings.NewReader("contents"),
					contentType: test.contentType,
				},
				Purpose: openai.FilePurposeAssistants,
			}).MarshalMultipart()
			if err == nil {
				t.Fatal("MarshalMultipart() error = nil, want invalid content type error")
			}
			if !strings.Contains(err.Error(), "invalid content type") {
				t.Errorf("MarshalMultipart() error = %q, want invalid content type", err)
			}
			if strings.ContainsAny(err.Error(), "\r\n") {
				t.Errorf("MarshalMultipart() error includes a newline: %q", err)
			}
			if data != nil || contentType != "" {
				t.Errorf("MarshalMultipart() = (%d bytes, %q), want (nil, empty)", len(data), contentType)
			}
		})
	}
}
