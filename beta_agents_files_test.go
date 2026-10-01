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
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithHeader("X-File-Test", "kept"))
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
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	source := betaLocalFile(t, "input.txt", "source")
	for _, destination := range []string{"/workspace/\xff", "/workspace/\xfe", "/tmp/source", "/workspace/a/../b", "/workspace/a//b", "/workspace/.codex/source", "/workspace/.managed-agents-x/source", "/workspace/outputs", "/workspace/"} {
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
	if _, err := client.Beta.Agents.Environments.Files.PrepareDirectory(context.Background(), "", "/workspace", []string{"*"}); err == nil {
		t.Error("accepted implicit working directory")
	}
	if requests.Load() != 0 {
		t.Fatal("preflight did not precede requests")
	}
}

func TestBetaAgentFilesDirectoryPatternErrorBeforeNetwork(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		t.Error("invalid pattern made a request")
	}))
	defer server.Close()
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	source := betaLocalFile(t, "nomatch.txt", "source")
	_, err := client.Beta.Agents.Environments.Files.PrepareDirectory(context.Background(), filepath.Dir(source), "/workspace/data", []string{"nomatch*["})
	if !errors.Is(err, filepath.ErrBadPattern) || requests.Load() != 0 {
		t.Fatalf("err=%v requests=%d", err, requests.Load())
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
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
	source := betaLocalFile(t, "input.txt", "source")
	files := map[string]string{"/workspace/a": source, "/workspace/b": source}
	prepared, err := client.Beta.Agents.Environments.Files.Prepare(context.Background(), files)
	var failure *openai.BetaAgentFilePreparationError
	if !errors.As(err, &failure) || len(failure.Prepared.Uploads) != 1 || failure.Prepared.Uploads[0].ID != "owned-upload" || prepared != failure.Prepared {
		t.Fatalf("lost owned uploads: %v", err)
	}
	before := requests.Load()
	for _, opts := range [][]option.RequestOption{{option.WithHeader("Idempotency-Key", "call-key")}, nil} {
		keyed := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithHeader("Idempotency-Key", "default-key"))
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
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	prepared, err := client.Beta.Agents.Environments.Files.PrepareDirectory(context.Background(), directory, "/workspace/data", []string{"*.txt"})
	if err != nil || len(prepared.Files) != 1 || requests.Load() != 1 {
		t.Fatalf("selection=%v err=%v", prepared, err)
	}
	wire, err := json.Marshal(prepared.Files)
	if err != nil {
		t.Fatal(err)
	}
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
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
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
			client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
			_, err := client.Beta.Agents.Sessions.Artifacts.ForResult(&openai.BetaAgentTurnResult{Turn: openai.Turn{ID: "root", SessionID: "session"}}).Download(context.Background(), "same", io.Discard)
			if err == nil {
				t.Fatal("expected lookup error")
			}
		})
	}
}

// The mock accepts these inputs to exercise forwarding, not backend support.
func TestBetaAgentFilesLimitsAreServerOwned(t *testing.T) {
	var uploads atomic.Int32
	var uploadedBytes atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		reader, err := r.MultipartReader()
		if err != nil {
			t.Error(err)
			return
		}
		for {
			part, err := reader.NextPart()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				t.Error(err)
				return
			}
			n, err := io.Copy(io.Discard, part)
			if err != nil {
				t.Error(err)
				return
			}
			if part.FormName() == "file" {
				uploadedBytes.Add(n)
			}
		}
		_, _ = fmt.Fprintf(w, `{"id":"file-%d"}`, uploads.Add(1))
	}))
	defer server.Close()
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	source := betaLocalFile(t, "empty", "")
	files := make(map[string]string)
	for i := 0; i < 51; i++ {
		files[fmt.Sprintf("/workspace/%d", i)] = source
	}
	prepared, err := client.Beta.Agents.Environments.Files.Prepare(context.Background(), files)
	if err != nil || len(prepared.Files) != 51 || uploads.Load() != 51 {
		t.Fatalf("prepared=%v err=%v", prepared, err)
	}
	directory := t.TempDir()
	for i := 0; i < 51; i++ {
		if writeErr := os.WriteFile(filepath.Join(directory, fmt.Sprint(i)), nil, 0600); writeErr != nil {
			t.Fatal(writeErr)
		}
	}
	prepared, err = client.Beta.Agents.Environments.Files.PrepareDirectory(context.Background(), directory, "/workspace", []string{"*"})
	if err != nil || len(prepared.Files) != 51 || uploads.Load() != 102 {
		t.Fatalf("directory prepared=%v err=%v", prepared, err)
	}
	for _, size := range []int64{26 * 1024 * 1024, 51 * 1024 * 1024} {
		if err := os.Truncate(source, size); err != nil {
			t.Fatal(err)
		}
		before := uploadedBytes.Load()
		prepared, err := client.Beta.Agents.Environments.Files.Prepare(context.Background(), map[string]string{"/workspace/a": source, "/workspace/b": source})
		if err != nil || len(prepared.Files) != 2 || uploadedBytes.Load()-before != size*2 {
			t.Fatalf("size=%d prepared=%v bytes=%d err=%v", size, prepared, uploadedBytes.Load()-before, err)
		}
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
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
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
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithHeader("Idempotency-Key", "default-key"))
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
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
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

func TestBetaAgentFilesLongDestinationAPIError(t *testing.T) {
	for _, character := range []string{"a", "🙂"} {
		t.Run(character, func(t *testing.T) {
			destination := "/workspace/" + strings.Repeat(character, 4097)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.URL.Path == "/files" {
					_, _ = fmt.Fprint(w, `{"id":"owned"}`)
					return
				}
				var body struct {
					Path string `json:"path"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				if body.Path != destination {
					t.Error("destination changed before API submission")
				}
				w.WriteHeader(http.StatusBadRequest)
				_, _ = fmt.Fprint(w, `{"error":{"message":"destination exceeds server limit","type":"invalid_request_error"}}`)
			}))
			defer server.Close()
			client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
			_, err := client.Beta.Agents.Environments.Files.Upload(context.Background(), "env", betaLocalFile(t, "input", "source"), destination)
			var failure *openai.BetaAgentFilePreparationError
			var apiError *openai.Error
			if !errors.As(err, &failure) || len(failure.Prepared.Uploads) != 1 || failure.Prepared.Uploads[0].ID != "owned" || !errors.As(err, &apiError) || apiError.StatusCode != http.StatusBadRequest {
				t.Fatalf("API rejection lost error or upload ownership: %v", err)
			}
		})
	}
}

func TestBetaAgentFilesEmptyResponsePreservesUploads(t *testing.T) {
	for _, stage := range []bool{false, true} {
		t.Run(fmt.Sprint(stage), func(t *testing.T) {
			var uploads atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.URL.Path == "/files" && uploads.Add(1) == 1 {
					_, _ = fmt.Fprint(w, `{"id":"owned"}`)
				} else {
					_, _ = fmt.Fprint(w, `null`)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
			source := betaLocalFile(t, "input", "source")
			var err error
			if stage {
				_, err = client.Beta.Agents.Environments.Files.Upload(context.Background(), "env", source, "/workspace/a")
			} else {
				_, err = client.Beta.Agents.Environments.Files.Prepare(context.Background(), map[string]string{"/workspace/a": source, "/workspace/b": source})
			}
			var failure *openai.BetaAgentFilePreparationError
			if !errors.As(err, &failure) || len(failure.Prepared.Uploads) != 1 || failure.Prepared.Uploads[0].ID != "owned" {
				t.Fatal("empty response lost earlier upload", err)
			}
		})
	}
}

func TestBetaAgentFilesCancellationBetweenUploads(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	first := betaLocalFile(t, "a.txt", "first")
	second := betaLocalFile(t, "b.txt", "second")
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"owned"}`)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0), option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		response, err := next(r)
		if err == nil {
			// Cancel only after the response body can be consumed by the generated decoder.
			response.Body = &betaCancelAfterRead{ReadCloser: response.Body, cancel: func() { cancel(); _ = os.Remove(second) }}
		}
		return response, err
	}))
	prepared, err := client.Beta.Agents.Environments.Files.Prepare(ctx, map[string]string{"/workspace/a": first, "/workspace/b": second})
	if !errors.Is(err, context.Canceled) || prepared == nil || len(prepared.Uploads) != 1 || requests.Load() != 1 {
		t.Fatalf("prepared=%v requests=%d err=%v", prepared, requests.Load(), err)
	}
}

type betaCancelAfterRead struct {
	io.ReadCloser
	cancel func()
}

func (r *betaCancelAfterRead) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	if errors.Is(err, io.EOF) {
		r.cancel()
	}
	return n, err
}

func TestBetaAgentResultArtifactContentErrorClosesBody(t *testing.T) {
	body := &betaArtifactBody{Reader: strings.NewReader("partial")}
	sentinel := errors.New("content unavailable")
	client := openai.NewClient(option.WithBaseURL("https://example.invalid"), option.WithAPIKey("synthetic"), option.WithMaxRetries(0), option.WithMiddleware(func(r *http.Request, _ option.MiddlewareNext) (*http.Response, error) {
		response := &http.Response{StatusCode: 200, Header: http.Header{}, Request: r}
		if strings.HasSuffix(r.URL.Path, "/content") {
			response.Body = body
			return response, sentinel
		}
		response.Header.Set("Content-Type", "application/json")
		response.Body = io.NopCloser(strings.NewReader(`{"data":[{"id":"artifact","turn_id":"root","path":"/workspace/report"}],"has_more":false}`))
		return response, nil
	}))
	var destination bytes.Buffer
	_, err := client.Beta.Agents.Sessions.Artifacts.ForResult(&openai.BetaAgentTurnResult{Turn: openai.Turn{ID: "root", SessionID: "session"}}).Download(context.Background(), "/workspace/report", &destination)
	if !errors.Is(err, sentinel) || !body.closed {
		t.Fatalf("closed=%t err=%v", body.closed, err)
	}
}

func TestBetaAgentFilesCanceledBeforePreflight(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client := openai.NewClient(option.WithAPIKey("synthetic"))
	missing := filepath.Join(t.TempDir(), "missing")
	_, err := client.Beta.Agents.Environments.Files.Prepare(ctx, map[string]string{"/workspace/input": missing})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("prepare should honor cancellation before accessing source: %v", err)
	}
	_, err = client.Beta.Agents.Environments.Files.PrepareDirectory(ctx, missing, "/workspace", []string{"*"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("directory preparation should honor cancellation before scanning: %v", err)
	}
	_, err = client.Beta.Agents.Environments.Files.Upload(ctx, "env", missing, "/workspace/input")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("upload should honor cancellation before accessing source: %v", err)
	}
}
