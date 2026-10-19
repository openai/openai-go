package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestUnknownEventDiagnostics(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: {\"type\":\"unknown\",\"b64_json\":\"synthetic-secret-customer-content\"}\n\ndata: [DONE]\n\n")
	}))
	defer server.Close()
	originalTransport := http.DefaultTransport
	http.DefaultTransport = server.Client().Transport
	t.Cleanup(func() { http.DefaultTransport = originalTransport })
	t.Setenv("OPENAI_BASE_URL", server.URL)
	t.Setenv("OPENAI_API_KEY", "synthetic-key")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	t.Setenv("OPENAI_CUSTOM_HEADERS", "")
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = reader.Close() }()
	original := os.Stdout
	os.Stdout = writer
	defer func() { os.Stdout = original; _ = writer.Close() }()
	main()
	if closeErr := writer.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(output), "synthetic-secret") {
		t.Fatalf("unknown event exposed provider data: %s", output)
	}
	if !strings.Contains(string(output), "Received unknown event type") {
		t.Fatalf("missing diagnostic: %s", output)
	}
}
