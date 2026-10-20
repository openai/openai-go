package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestFineTuningDiagnostics(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/files":
			_, _ = io.WriteString(w, `{"id":"file-test"}`)
		case r.URL.Path == "/files/file-test":
			_, _ = io.WriteString(w, `{"id":"file-test","status":"processed"}`)
		case r.URL.Path == "/fine_tuning/jobs":
			_, _ = io.WriteString(w, `{"id":"job-test","status":"running"}`)
		case r.URL.Path == "/fine_tuning/jobs/job-test":
			_, _ = io.WriteString(w, `{"id":"job-test","status":"synthetic-secret-status"}`)
		case r.URL.Path == "/fine_tuning/jobs/job-test/events":
			_, _ = io.WriteString(w, `{"data":[{"id":"event-test","created_at":1,"level":"error","message":"synthetic-secret-error-detail"}],"has_more":false}`)
		default:
			t.Errorf("unexpected request: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
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
		t.Fatalf("fine-tuning diagnostics exposed provider data: %s", output)
	}
	if !strings.Contains(string(output), "File status: processed") || !strings.Contains(string(output), "event received") {
		t.Fatalf("missing progress diagnostics: %s", output)
	}
}
