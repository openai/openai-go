package openai_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func betaLocalFile(t *testing.T, name, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestBetaAgentFilesPrepareAndStage(t *testing.T) {
	var uploads, stages atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-File-Test") != "kept" {
			t.Error("request options lost")
		}
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/files" {
			if err := r.ParseMultipartForm(1 << 20); err != nil {
				t.Error(err)
				return
			}
			defer func() {
				if r.MultipartForm != nil {
					_ = r.MultipartForm.RemoveAll()
				}
			}()
			f, header, err := r.FormFile("file")
			if err != nil {
				t.Error(err)
				return
			}
			data, err := io.ReadAll(f)
			_ = f.Close()
			if err != nil || string(data) != "source" || header.Filename != "input.txt" || r.FormValue("purpose") != "user_data" {
				t.Error("wrong upload", header.Filename, err)
			}
			id := uploads.Add(1)
			_, _ = fmt.Fprintf(w, `{"id":"file-%d","object":"file"}`, id)
		} else {
			stages.Add(1)
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			if r.URL.Path != "/agents/environments/env/files" || body["type"] != "file_id" || body["path"] != "/workspace/input.txt" || body["file_id"] != "file-2" {
				t.Errorf("wrong stage %s %#v", r.URL.Path, body)
			}
			_, _ = fmt.Fprint(w, `{"path":"/workspace/input.txt","type":"file_id","file_id":"file-2"}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithHeader("X-File-Test", "kept"))
	source := betaLocalFile(t, "input.txt", "source")
	prepared, err := client.Beta.Agents.Environments.Files.Prepare(context.Background(), map[string]string{"/workspace/input.txt": source})
	if err != nil || len(prepared.Uploads) != 1 || len(prepared.Files) != 1 {
		t.Fatalf("prepared=%v err=%v", prepared, err)
	}
	wire, err := json.Marshal(prepared.Files)
	if err != nil || !strings.Contains(string(wire), `"file_id":"file-1"`) {
		t.Fatal("not reusable creation inputs", err)
	}
	staged, err := client.Beta.Agents.Environments.Files.Upload(context.Background(), "env", source, "/workspace/input.txt")
	if err != nil || staged.Upload.ID != "file-2" || stages.Load() != 1 {
		t.Fatalf("staged=%v err=%v", staged, err)
	}
}

func TestBetaAgentFilesPreflightBeforeNetwork(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		t.Error("invalid selection made a request")
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	source := betaLocalFile(t, "input.txt", "source")
	for _, destination := range []string{"/tmp/source", "/workspace/a/../b", "/workspace/a//b", "/workspace/.codex/source", "/workspace/.managed-agents-x/source", "/workspace/outputs", "/workspace/"} {
		if _, err := client.Beta.Agents.Environments.Files.Prepare(context.Background(), map[string]string{destination: source}); err == nil {
			t.Errorf("accepted %q", destination)
		}
	}
	for _, files := range []map[string]string{
		{"/workspace/a": source, "/workspace/a/b": source},
		{"/workspace/a": source, "/workspace/z": "missing-source"},
	} {
		if _, err := client.Beta.Agents.Environments.Files.Prepare(context.Background(), files); err == nil {
			t.Error("accepted invalid batch")
		}
	}
	symlink := filepath.Join(t.TempDir(), "linked")
	if err := os.Symlink(source, symlink); err == nil {
		if _, err := client.Beta.Agents.Environments.Files.Prepare(context.Background(), map[string]string{"/workspace/x": symlink}); err == nil {
			t.Error("followed symlink")
		}
	}
	if _, err := client.Beta.Agents.Environments.Files.PrepareDirectory(context.Background(), filepath.Dir(source), "/workspace", nil); err == nil {
		t.Error("accepted implicit directory selection")
	}
	if requests.Load() != 0 {
		t.Fatal("preflight did not precede requests")
	}
}

func TestBetaAgentFilesPartialUploadsAndIdempotency(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		if n == 1 {
			_, _ = fmt.Fprint(w, `{"id":"owned-upload"}`)
		} else {
			w.WriteHeader(500)
			_, _ = fmt.Fprint(w, `{"error":{"message":"synthetic"}}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
	source := betaLocalFile(t, "input.txt", "source")
	files := map[string]string{"/workspace/a": source, "/workspace/b": source}
	prepared, err := client.Beta.Agents.Environments.Files.Prepare(context.Background(), files)
	var failure *openai.BetaAgentFilePreparationError
	if !errors.As(err, &failure) || len(failure.Prepared.Uploads) != 1 || failure.Prepared.Uploads[0].ID != "owned-upload" || prepared != failure.Prepared {
		t.Fatalf("lost owned uploads: %v", err)
	}
	before := requests.Load()
	for _, opts := range [][]option.RequestOption{{option.WithHeader("Idempotency-Key", "call-key")}, nil} {
		keyed := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithHeader("Idempotency-Key", "default-key"))
		if _, err := keyed.Beta.Agents.Environments.Files.Prepare(context.Background(), files, opts...); err == nil {
			t.Error("reused batch idempotency key")
		}
	}
	if requests.Load() != before {
		t.Error("idempotency rejection made requests")
	}
}

func TestBetaAgentFilesDirectorySelection(t *testing.T) {
	directory := t.TempDir()
	for _, name := range []string{"include.txt", "exclude.bin"} {
		if err := os.WriteFile(filepath.Join(directory, name), []byte(name), 0600); err != nil {
			t.Fatal(err)
		}
	}
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"selected"}`)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	prepared, err := client.Beta.Agents.Environments.Files.PrepareDirectory(context.Background(), directory, "/workspace/data", []string{"*.txt"})
	if err != nil || len(prepared.Files) != 1 || requests.Load() != 1 {
		t.Fatalf("selection=%v err=%v", prepared, err)
	}
	wire, _ := json.Marshal(prepared.Files)
	if !strings.Contains(string(wire), `"path":"/workspace/data/include.txt"`) {
		t.Fatal("wrong destination", string(wire))
	}
}

func TestBetaAgentResultArtifactDownload(t *testing.T) {
	var pages, contents atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Artifact") != "kept" {
			t.Error("request options lost")
		}
		if strings.HasSuffix(r.URL.Path, "/content") {
			contents.Add(1)
			if !strings.Contains(r.URL.Path, "/exact/content") {
				t.Error("wrong artifact content", r.URL.Path)
			}
			_, _ = fmt.Fprint(w, "downloaded")
			return
		}
		pages.Add(1)
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("after") == "" {
			_, _ = fmt.Fprint(w, `{"data":[{"id":"previous","turn_id":"old","path":"/workspace/outputs/report.md"}],"has_more":true,"last_id":"previous"}`)
		} else {
			_, _ = fmt.Fprint(w, `{"data":[{"id":"exact","turn_id":"root","path":"/workspace/outputs/report.md"}],"has_more":false}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	result := &openai.BetaAgentTurnResult{Turn: openai.Turn{ID: "root", SessionID: "session"}}
	artifacts := client.Beta.Agents.Sessions.Artifacts.ForResult(result)
	if pages.Load() != 0 {
		t.Fatal("artifact accessor listed eagerly")
	}
	result.Turn.ID = "mutated"
	var destination bytes.Buffer
	artifact, err := artifacts.Download(context.Background(), "/workspace/outputs/report.md", &destination, option.WithHeader("X-Artifact", "kept"))
	if err != nil || artifact.ID != "exact" || destination.String() != "downloaded" || pages.Load() != 2 || contents.Load() != 1 {
		t.Fatalf("artifact=%v err=%v", artifact, err)
	}
}

func TestBetaAgentResultArtifactMissingAndAmbiguous(t *testing.T) {
	for _, body := range []string{`{"data":[],"has_more":false}`, `{"data":[{"id":"one","turn_id":"root","path":"same"},{"id":"two","turn_id":"root","path":"same"}],"has_more":false}`} {
		t.Run(body, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasSuffix(r.URL.Path, "/content") {
					t.Error("ambiguous download opened content")
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprint(w, body)
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
			_, err := client.Beta.Agents.Sessions.Artifacts.ForResult(&openai.BetaAgentTurnResult{Turn: openai.Turn{ID: "root", SessionID: "session"}}).Download(context.Background(), "same", io.Discard)
			if err == nil {
				t.Fatal("expected lookup error")
			}
		})
	}
}

func TestBetaAgentFilesKnownLimits(t *testing.T) {
	client := openai.NewClient(option.WithBaseURL("http://127.0.0.1:1"), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
	source := betaLocalFile(t, "empty", "")
	files := make(map[string]string)
	for i := 0; i < 51; i++ {
		files[fmt.Sprintf("/workspace/%d", i)] = source
	}
	if _, err := client.Beta.Agents.Environments.Files.Prepare(context.Background(), files); err == nil || !strings.Contains(err.Error(), "50 files") {
		t.Fatal("missing file-count preflight", err)
	}
	if err := os.Truncate(source, 26*1024*1024); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Beta.Agents.Environments.Files.Prepare(context.Background(), map[string]string{"/workspace/a": source, "/workspace/b": source}); err == nil || !strings.Contains(err.Error(), "total") {
		t.Fatal("missing aggregate preflight", err)
	}
	if err := os.Truncate(source, 51*1024*1024); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Beta.Agents.Environments.Files.Upload(context.Background(), "environment", source, "/workspace/a"); err == nil || !strings.Contains(err.Error(), "50 MiB") {
		t.Fatal("missing singleton preflight", err)
	}
}

func TestBetaAgentFilesStageFailureRetainsUpload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/files" {
			_, _ = fmt.Fprint(w, `{"id":"owned"}`)
			return
		}
		w.WriteHeader(500)
		_, _ = fmt.Fprint(w, `{"error":{"message":"staging failed"}}`)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
	_, err := client.Beta.Agents.Environments.Files.Upload(context.Background(), "env", betaLocalFile(t, "input", "source"), "/workspace/source")
	var failure *openai.BetaAgentFilePreparationError
	if !errors.As(err, &failure) || len(failure.Prepared.Uploads) != 1 || failure.Prepared.Uploads[0].ID != "owned" {
		t.Fatal("staging error lost uploaded file", err)
	}
}

func TestBetaAgentFilesIdempotencyHeaderOverride(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Idempotency-Key") == "default-key" {
			t.Error("removed default header survived")
		}
		id := requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"id":"file-%d"}`, id)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithHeader("Idempotency-Key", "default-key"))
	source := betaLocalFile(t, "source", "synthetic")
	prepared, err := client.Beta.Agents.Environments.Files.Prepare(context.Background(), map[string]string{"/workspace/a": source, "/workspace/b": source}, option.WithHeaderDel("Idempotency-Key"))
	if err != nil || len(prepared.Uploads) != 2 || requests.Load() != 2 {
		t.Fatalf("per-call header removal ignored: prepared=%v err=%v", prepared, err)
	}
}

func TestBetaAgentFilesTrustedParentAlias(t *testing.T) {
	root := t.TempDir()
	actual := filepath.Join(root, "actual")
	directory := filepath.Join(actual, "selected")
	if err := os.MkdirAll(directory, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "input.txt"), []byte("input"), 0600); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(root, "alias")
	if err := os.Symlink(actual, alias); err != nil {
		t.Skip("symlinks unavailable:", err)
	}
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"file","object":"file"}`)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	selected := filepath.Join(alias, "selected")
	if _, err := client.Beta.Agents.Environments.Files.Prepare(context.Background(), map[string]string{"/workspace/input.txt": filepath.Join(selected, "input.txt")}); err != nil {
		t.Fatal("trusted file parent alias rejected", err)
	}
	if _, err := client.Beta.Agents.Environments.Files.PrepareDirectory(context.Background(), selected, "/workspace/data", []string{"*.txt"}); err != nil {
		t.Fatal("trusted directory parent alias rejected", err)
	}
	if _, err := client.Beta.Agents.Environments.Files.PrepareDirectory(context.Background(), alias, "/workspace/data", []string{"*"}); err == nil {
		t.Fatal("selected symlink root accepted")
	}
	if _, err := client.Beta.Agents.Environments.Files.PrepareDirectory(context.Background(), alias+string(filepath.Separator), "/workspace/data", []string{"*"}); err == nil {
		t.Fatal("selected symlink root with trailing separator accepted")
	}
	if requests.Load() != 2 {
		t.Fatal("unexpected requests", requests.Load())
	}
}

type betaArtifactTransport func(*http.Request) (*http.Response, error)

func (f betaArtifactTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type betaArtifactBody struct {
	io.Reader
	closed bool
}

func (b *betaArtifactBody) Close() error { b.closed = true; return nil }

type betaArtifactFailWriter struct{ err error }

func (w betaArtifactFailWriter) Write([]byte) (int, error) { return 0, w.err }

func TestBetaAgentResultArtifactWriterFailureClosesBody(t *testing.T) {
	body := &betaArtifactBody{Reader: strings.NewReader("output")}
	client := openai.NewClient(option.WithBaseURL("https://example.invalid"), option.WithAPIKey("synthetic"), option.WithHTTPClient(&http.Client{Transport: betaArtifactTransport(func(r *http.Request) (*http.Response, error) {
		response := &http.Response{StatusCode: 200, Header: http.Header{}, Request: r}
		if strings.HasSuffix(r.URL.Path, "/content") {
			response.Body = body
		} else {
			response.Header.Set("Content-Type", "application/json")
			response.Body = io.NopCloser(strings.NewReader(`{"data":[{"id":"artifact","turn_id":"root","path":"/workspace/report"}],"has_more":false}`))
		}
		return response, nil
	})}))
	sentinel := errors.New("destination unavailable")
	artifact, err := client.Beta.Agents.Sessions.Artifacts.ForResult(&openai.BetaAgentTurnResult{Turn: openai.Turn{ID: "root", SessionID: "session"}}).Download(context.Background(), "/workspace/report", betaArtifactFailWriter{sentinel})
	if !errors.Is(err, sentinel) || artifact == nil || artifact.ID != "artifact" || !body.closed {
		t.Fatalf("writer failure lost metadata or left body open: %v %v", artifact, err)
	}
}
