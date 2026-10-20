package openai_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type vectorFileDoer func(*http.Request) (*http.Response, error)

func (f vectorFileDoer) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

type vectorFileContextKey struct{}

type vectorFileHelperCase struct {
	name       string
	stages     []string
	attachment string
	call       func(context.Context, *openai.VectorStoreFileService, ...option.RequestOption) (*openai.VectorStoreFile, error)
}

func vectorFileHelperCases() []vectorFileHelperCase {
	upload := func() openai.FileNewParams {
		return openai.FileNewParams{
			File:    strings.NewReader("synthetic file contents"),
			Purpose: openai.FilePurposeAssistants,
		}
	}
	return []vectorFileHelperCase{
		{
			name:       "new and poll",
			stages:     []string{"attach", "poll"},
			attachment: `{"file_id":"file_existing","attributes":{"tag":"kept"}}`,
			call: func(ctx context.Context, service *openai.VectorStoreFileService, opts ...option.RequestOption) (*openai.VectorStoreFile, error) {
				return service.NewAndPoll(ctx, "vs_test", openai.VectorStoreFileNewParams{
					FileID: "file_existing",
					Attributes: map[string]openai.VectorStoreFileNewParamsAttributeUnion{
						"tag": {OfString: openai.String("kept")},
					},
				}, 7, opts...)
			},
		},
		{
			name:       "upload",
			stages:     []string{"upload", "attach"},
			attachment: `{"file_id":"file_uploaded"}`,
			call: func(ctx context.Context, service *openai.VectorStoreFileService, opts ...option.RequestOption) (*openai.VectorStoreFile, error) {
				return service.Upload(ctx, "vs_test", upload(), opts...)
			},
		},
		{
			name:       "upload and poll",
			stages:     []string{"upload", "attach", "poll"},
			attachment: `{"file_id":"file_uploaded"}`,
			call: func(ctx context.Context, service *openai.VectorStoreFileService, opts ...option.RequestOption) (*openai.VectorStoreFile, error) {
				return service.UploadAndPoll(ctx, "vs_test", upload(), 7, opts...)
			},
		},
	}
}

func vectorFileStage(t *testing.T, req *http.Request) string {
	t.Helper()
	switch req.Method + " " + req.URL.Path {
	case "POST /v1/files":
		return "upload"
	case "POST /v1/vector_stores/vs_test/files":
		return "attach"
	case "GET /v1/vector_stores/vs_test/files/file_attached":
		return "poll"
	default:
		t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
		return ""
	}
}

func vectorFileResponse(req *http.Request, body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
}

func vectorFileResponseBody(stage string) string {
	switch stage {
	case "upload":
		return `{"id":"file_uploaded"}`
	case "attach":
		return `{"id":"file_attached","status":"in_progress"}`
	default:
		return `{"id":"file_attached","status":"completed"}`
	}
}

func vectorFileService(doer vectorFileDoer) openai.VectorStoreFileService {
	return openai.NewVectorStoreFileService(
		option.WithBaseURL("https://vector-files.example/v1"),
		option.WithAPIKey("synthetic-key"),
		option.WithMaxRetries(0),
		option.WithHTTPClient(doer),
		option.WithHeader("X-Service", "kept"),
		option.WithHeader("X-Precedence", "service"),
	)
}

func TestVectorStoreFileHelperRequests(t *testing.T) {
	for _, test := range vectorFileHelperCases() {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), vectorFileContextKey{}, "kept")
			var stages []string
			service := vectorFileService(func(req *http.Request) (*http.Response, error) {
				stage := vectorFileStage(t, req)
				stages = append(stages, stage)
				if req.Context().Value(vectorFileContextKey{}) != "kept" {
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
				if stage == "upload" {
					checkVectorFileUpload(t, req)
				} else if got := req.Header.Get("OpenAI-Beta"); got != "assistants=v2" {
					t.Errorf("beta header on %s = %q", stage, got)
				}
				if stage == "attach" {
					var got, want any
					if err := json.NewDecoder(req.Body).Decode(&got); err != nil {
						t.Fatal(err)
					}
					if err := json.Unmarshal([]byte(test.attachment), &want); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, want) {
						t.Errorf("attachment = %#v, want %#v", got, want)
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
				return vectorFileResponse(req, vectorFileResponseBody(stage)), nil
			})
			result, err := test.call(ctx, &service,
				option.WithHeader("X-Precedence", "method"),
				option.WithHeader("X-Call", "kept"),
			)
			if err != nil {
				t.Fatal(err)
			}
			wantStatus := openai.VectorStoreFileStatusCompleted
			if test.name == "upload" {
				wantStatus = openai.VectorStoreFileStatusInProgress
			}
			if result == nil || result.ID != "file_attached" || result.Status != wantStatus {
				t.Fatalf("result = %#v, want file_attached/%s", result, wantStatus)
			}
			if !reflect.DeepEqual(stages, test.stages) {
				t.Fatalf("request stages = %v, want %v", stages, test.stages)
			}
		})
	}
}

func checkVectorFileUpload(t *testing.T, req *http.Request) {
	t.Helper()
	mediaType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		t.Fatal(err)
	}
	if mediaType != "multipart/form-data" {
		t.Fatalf("upload content type = %q", mediaType)
	}
	reader := multipart.NewReader(req.Body, params["boundary"])
	fields := map[string]string{}
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		data, readErr := io.ReadAll(part)
		if err := part.Close(); err != nil {
			t.Fatal(err)
		}
		if readErr != nil {
			t.Fatal(readErr)
		}
		fields[part.FormName()] = string(data)
	}
	want := map[string]string{
		"file":    "synthetic file contents",
		"purpose": "assistants",
	}
	if !reflect.DeepEqual(fields, want) {
		t.Errorf("upload fields = %#v, want %#v", fields, want)
	}
}

func TestVectorStoreFileHelperErrors(t *testing.T) {
	for _, test := range vectorFileHelperCases() {
		for failIndex, failStage := range test.stages {
			t.Run(test.name+"/"+failStage, func(t *testing.T) {
				failure := errors.New("synthetic transport failure")
				var stages []string
				service := vectorFileService(func(req *http.Request) (*http.Response, error) {
					stage := vectorFileStage(t, req)
					stages = append(stages, stage)
					if stage == failStage {
						return nil, failure
					}
					return vectorFileResponse(req, vectorFileResponseBody(stage)), nil
				})
				result, err := test.call(context.Background(), &service)
				if result != nil || !errors.Is(err, failure) {
					t.Fatalf("result = %#v, error = %v; want nil and original failure", result, err)
				}
				wantError := failure.Error()
				if failStage == "poll" {
					wantError = "vector store file poll: received " + wantError
				}
				if err.Error() != wantError {
					t.Errorf("error = %q, want %q", err.Error(), wantError)
				}
				if want := test.stages[:failIndex+1]; !reflect.DeepEqual(stages, want) {
					t.Fatalf("request stages = %v, want %v", stages, want)
				}
			})
		}
	}
}
