package openai_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
	"testing/synctest"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type vectorBatchDoer func(*http.Request) (*http.Response, error)

func (f vectorBatchDoer) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

type vectorBatchContextKey struct{}

type vectorBatchResult struct {
	batch *openai.VectorStoreFileBatch
	err   error
}

func vectorBatchService(doer vectorBatchDoer) openai.VectorStoreFileBatchService {
	return openai.NewVectorStoreFileBatchService(
		option.WithBaseURL("https://vector-batches.example/v1"),
		option.WithAPIKey("synthetic-key"),
		option.WithMaxRetries(0),
		option.WithHTTPClient(doer),
		option.WithHeader("X-Service", "kept"),
		option.WithHeader("X-Precedence", "service"),
	)
}

func vectorBatchCallOptions() []option.RequestOption {
	return []option.RequestOption{
		option.WithHeader("X-Precedence", "method"),
		option.WithHeader("X-Call", "kept"),
	}
}

func vectorBatchUpload(contents string) openai.FileNewParams {
	return openai.FileNewParams{
		File:    strings.NewReader(contents),
		Purpose: openai.FilePurposeAssistants,
	}
}

func vectorBatchStage(req *http.Request) (string, error) {
	switch req.Method + " " + req.URL.Path {
	case "POST /v1/files":
		return "upload", nil
	case "POST /v1/vector_stores/vs_test/file_batches":
		return "create", nil
	case "GET /v1/vector_stores/vs_test/file_batches/batch_created":
		return "poll", nil
	default:
		return "", fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
	}
}

func vectorBatchResponse(req *http.Request, body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
}

func vectorBatchResponseBody(stage string) string {
	switch stage {
	case "upload":
		return `{"id":"file_uploaded"}`
	case "create":
		return `{"id":"batch_created","status":"in_progress"}`
	default:
		return `{"id":"batch_created","status":"completed"}`
	}
}

func checkVectorBatchRequest(t *testing.T, req *http.Request, stage string) {
	t.Helper()
	if req.Context().Value(vectorBatchContextKey{}) != "kept" {
		t.Error("request lost the caller context")
	}
	for key, want := range map[string]string{
		"Authorization": "Bearer synthetic-key",
		"X-Service":     "kept",
		"X-Precedence":  "method",
		"X-Call":        "kept",
	} {
		if got := req.Header.Get(key); got != want {
			t.Errorf("%s on %s = %q, want %q", key, stage, got, want)
		}
	}
	if stage != "upload" {
		if got := req.Header.Get("OpenAI-Beta"); got != "assistants=v2" {
			t.Errorf("beta header on %s = %q", stage, got)
		}
	}
	if stage == "poll" {
		if got := req.Header.Get("X-Stainless-Poll-Helper"); got != "true" {
			t.Errorf("poll helper header = %q", got)
		}
		if got := req.Header.Get("X-Stainless-Poll-Interval"); got != "7" {
			t.Errorf("poll interval header = %q", got)
		}
	}
}

func readVectorBatchUpload(req *http.Request) (string, error) {
	mediaType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}
	if mediaType != "multipart/form-data" {
		return "", fmt.Errorf("upload content type = %q", mediaType)
	}
	reader := multipart.NewReader(req.Body, params["boundary"])
	fields := map[string]string{}
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		data, readErr := io.ReadAll(part)
		closeErr := part.Close()
		if err := errors.Join(readErr, closeErr); err != nil {
			return "", err
		}
		fields[part.FormName()] = string(data)
	}
	if len(fields) != 2 || fields["purpose"] != "assistants" {
		return "", fmt.Errorf("unexpected upload fields: %#v", fields)
	}
	return fields["file"], nil
}

func TestVectorStoreFileBatchHelperNewAndPoll(t *testing.T) {
	ctx := context.WithValue(context.Background(), vectorBatchContextKey{}, "kept")
	var stages []string
	service := vectorBatchService(func(req *http.Request) (*http.Response, error) {
		stage, err := vectorBatchStage(req)
		if err != nil {
			return nil, err
		}
		stages = append(stages, stage)
		checkVectorBatchRequest(t, req, stage)
		if stage == "create" {
			var got, want any
			if err := json.NewDecoder(req.Body).Decode(&got); err != nil {
				return nil, err
			}
			if err := json.Unmarshal([]byte(`{"file_ids":["file_existing"],"attributes":{"tag":"kept"}}`), &want); err != nil {
				return nil, err
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("batch body = %#v, want %#v", got, want)
			}
		}
		return vectorBatchResponse(req, vectorBatchResponseBody(stage)), nil
	})
	result, err := service.NewAndPoll(ctx, "vs_test", openai.VectorStoreFileBatchNewParams{
		FileIDs: []string{"file_existing"},
		Attributes: map[string]openai.VectorStoreFileBatchNewParamsAttributeUnion{
			"tag": {OfString: openai.String("kept")},
		},
	}, 7, vectorBatchCallOptions()...)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || result.ID != "batch_created" || result.Status != openai.VectorStoreFileBatchStatusCompleted {
		t.Fatalf("result = %#v, want batch_created/completed", result)
	}
	if !slices.Equal(stages, []string{"create", "poll"}) {
		t.Fatalf("request stages = %v, want create then poll", stages)
	}
}

func TestVectorStoreFileBatchHelperConcurrentUploads(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx := context.WithValue(context.Background(), vectorBatchContextKey{}, "kept")
		started := make(chan string, 2)
		release := make(chan struct{})
		done := make(chan vectorBatchResult, 1)
		var mu sync.Mutex
		var completedUploads int
		var stages, createdIDs []string
		service := vectorBatchService(func(req *http.Request) (*http.Response, error) {
			stage, err := vectorBatchStage(req)
			if err != nil {
				return nil, err
			}
			checkVectorBatchRequest(t, req, stage)
			switch stage {
			case "upload":
				contents, err := readVectorBatchUpload(req)
				if err != nil {
					return nil, err
				}
				started <- contents
				<-release
				mu.Lock()
				completedUploads++
				mu.Unlock()
				return vectorBatchResponse(req, `{"id":"file_`+contents+`"}`), nil
			case "create":
				mu.Lock()
				completed := completedUploads
				mu.Unlock()
				if completed != 2 {
					t.Errorf("created batch after %d uploads, want 2", completed)
				}
				var body struct {
					FileIDs []string `json:"file_ids"`
				}
				if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
					return nil, err
				}
				createdIDs = body.FileIDs
			}
			stages = append(stages, stage)
			return vectorBatchResponse(req, vectorBatchResponseBody(stage)), nil
		})
		existingIDs := []string{"file_existing_1", "file_existing_2"}
		go func() {
			batch, err := service.UploadAndPoll(ctx, "vs_test",
				[]openai.FileNewParams{vectorBatchUpload("alpha"), vectorBatchUpload("beta")},
				existingIDs, 7, vectorBatchCallOptions()...)
			done <- vectorBatchResult{batch, err}
		}()
		synctest.Wait()
		if got := len(started); got != 2 {
			t.Errorf("concurrent uploads started = %d, want 2", got)
		}
		close(release)
		result := <-done
		if result.err != nil {
			t.Fatal(result.err)
		}
		if result.batch == nil || result.batch.ID != "batch_created" || result.batch.Status != openai.VectorStoreFileBatchStatusCompleted {
			t.Fatalf("result = %#v, want batch_created/completed", result.batch)
		}
		close(started)
		var contents []string
		for value := range started {
			contents = append(contents, value)
		}
		slices.Sort(contents)
		if !slices.Equal(contents, []string{"alpha", "beta"}) {
			t.Errorf("uploaded contents = %v", contents)
		}
		if len(createdIDs) != 4 || !slices.Equal(createdIDs[:2], existingIDs) {
			t.Fatalf("batch file IDs = %v, want existing IDs first", createdIDs)
		}
		uploadedIDs := slices.Clone(createdIDs[2:])
		slices.Sort(uploadedIDs)
		if !slices.Equal(uploadedIDs, []string{"file_alpha", "file_beta"}) {
			t.Errorf("uploaded file IDs = %v", uploadedIDs)
		}
		if !slices.Equal(stages, []string{"create", "poll"}) {
			t.Errorf("request stages after uploads = %v", stages)
		}
	})
}

func TestVectorStoreFileBatchHelperErrors(t *testing.T) {
	tests := []struct {
		name   string
		stages []string
		call   func(*openai.VectorStoreFileBatchService) (*openai.VectorStoreFileBatch, error)
	}{
		{
			name:   "new and poll",
			stages: []string{"create", "poll"},
			call: func(service *openai.VectorStoreFileBatchService) (*openai.VectorStoreFileBatch, error) {
				return service.NewAndPoll(context.Background(), "vs_test", openai.VectorStoreFileBatchNewParams{
					FileIDs: []string{"file_existing"},
				}, 7)
			},
		},
		{
			name:   "upload and poll",
			stages: []string{"upload", "create", "poll"},
			call: func(service *openai.VectorStoreFileBatchService) (*openai.VectorStoreFileBatch, error) {
				return service.UploadAndPoll(context.Background(), "vs_test",
					[]openai.FileNewParams{vectorBatchUpload("synthetic")}, nil, 7)
			},
		},
	}
	for _, test := range tests {
		for failIndex, failStage := range test.stages {
			t.Run(test.name+"/"+failStage, func(t *testing.T) {
				failure := errors.New("synthetic transport failure")
				var stages []string
				service := vectorBatchService(func(req *http.Request) (*http.Response, error) {
					stage, err := vectorBatchStage(req)
					if err != nil {
						return nil, err
					}
					stages = append(stages, stage)
					if stage == failStage {
						return nil, failure
					}
					return vectorBatchResponse(req, vectorBatchResponseBody(stage)), nil
				})
				result, err := test.call(&service)
				if result != nil || !errors.Is(err, failure) {
					t.Fatalf("result = %#v, error = %v; want nil and original failure", result, err)
				}
				wantError := failure.Error()
				if failStage == "poll" {
					wantError = "vector store file batch poll: received " + wantError
				}
				if err.Error() != wantError {
					t.Errorf("error = %q, want %q", err.Error(), wantError)
				}
				if want := test.stages[:failIndex+1]; !slices.Equal(stages, want) {
					t.Fatalf("request stages = %v, want %v", stages, want)
				}
			})
		}
	}
}

func TestVectorStoreFileBatchHelperEmptyFiles(t *testing.T) {
	for _, files := range [][]openai.FileNewParams{nil, {}} {
		calls := 0
		service := vectorBatchService(func(*http.Request) (*http.Response, error) {
			calls++
			return nil, errors.New("unexpected request")
		})
		result, err := service.UploadAndPoll(context.Background(), "vs_test", files, []string{"file_existing"}, 7)
		const want = "No `files` provided to process. If you've already uploaded files you should use `.NewAndPoll()` instead"
		if result != nil || err == nil || err.Error() != want || calls != 0 {
			t.Errorf("result = %#v, error = %v, requests = %d", result, err, calls)
		}
	}
}

func TestVectorStoreFileBatchHelperWaitsForAllUploads(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		failure := errors.New("synthetic upload failure")
		releaseFailure := make(chan struct{})
		releaseOther := make(chan struct{})
		started := make(chan string, 2)
		done := make(chan vectorBatchResult, 1)
		service := vectorBatchService(func(req *http.Request) (*http.Response, error) {
			stage, err := vectorBatchStage(req)
			if err != nil {
				return nil, err
			}
			if stage != "upload" {
				return nil, fmt.Errorf("unexpected %s after upload failure", stage)
			}
			contents, err := readVectorBatchUpload(req)
			if err != nil {
				return nil, err
			}
			started <- contents
			if contents == "failure" {
				<-releaseFailure
				return nil, failure
			}
			<-releaseOther
			return vectorBatchResponse(req, `{"id":"file_other"}`), nil
		})
		go func() {
			batch, err := service.UploadAndPoll(context.Background(), "vs_test",
				[]openai.FileNewParams{vectorBatchUpload("failure"), vectorBatchUpload("other")}, nil, 7)
			done <- vectorBatchResult{batch, err}
		}()
		synctest.Wait()
		if got := len(started); got != 2 {
			t.Errorf("concurrent uploads started = %d, want 2", got)
		}
		close(releaseFailure)
		synctest.Wait()
		select {
		case result := <-done:
			t.Errorf("returned before all uploads completed: %#v", result)
			close(releaseOther)
			return
		default:
		}
		close(releaseOther)
		result := <-done
		if result.batch != nil || !errors.Is(result.err, failure) || result.err.Error() != failure.Error() {
			t.Fatalf("result = %#v, error = %v; want original upload failure", result.batch, result.err)
		}
	})
}

func TestVectorStoreFileBatchHelperCancellation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		started := make(chan struct{}, 2)
		done := make(chan vectorBatchResult, 1)
		service := vectorBatchService(func(req *http.Request) (*http.Response, error) {
			stage, err := vectorBatchStage(req)
			if err != nil {
				return nil, err
			}
			if stage != "upload" {
				return nil, fmt.Errorf("unexpected %s after cancellation", stage)
			}
			started <- struct{}{}
			<-req.Context().Done()
			return nil, req.Context().Err()
		})
		go func() {
			batch, err := service.UploadAndPoll(ctx, "vs_test",
				[]openai.FileNewParams{vectorBatchUpload("alpha"), vectorBatchUpload("beta")}, nil, 7)
			done <- vectorBatchResult{batch, err}
		}()
		synctest.Wait()
		if got := len(started); got != 2 {
			t.Errorf("concurrent uploads started = %d, want 2", got)
		}
		cancel()
		result := <-done
		if result.batch != nil || !errors.Is(result.err, context.Canceled) {
			t.Fatalf("result = %#v, error = %v; want context cancellation", result.batch, result.err)
		}
	})
}
