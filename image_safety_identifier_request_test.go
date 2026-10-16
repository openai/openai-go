package openai_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
)

func TestImageEditSafetyIdentifierRequest(t *testing.T) {
	for _, streaming := range []bool{false, true} {
		method := "Edit"
		if streaming {
			method = "EditStreaming"
		}
		for _, tc := range []struct {
			name    string
			value   param.Opt[string]
			present bool
			want    string
		}{
			{name: "omitted"},
			{name: "null", value: param.Null[string]()},
			{name: "empty", value: openai.String(""), present: true},
			{name: "literal_null", value: openai.String("null"), present: true, want: "null"},
			{name: "value", value: openai.String("synthetic-identifier"), present: true, want: "synthetic-identifier"},
		} {
			t.Run(method+"/"+tc.name, func(t *testing.T) {
				calls := 0
				client := openai.NewClient(option.WithAPIKey("synthetic-key"), option.WithBaseURL("https://example.test/v1"), option.WithMaxRetries(0), option.WithHTTPClient(imageErrorHTTPDoer(func(r *http.Request) (*http.Response, error) {
					calls++
					if r.Method != http.MethodPost || r.URL.Path != "/v1/images/edits" {
						t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
					}
					reader, err := r.MultipartReader()
					if err != nil {
						t.Fatal(err)
					}
					fields := map[string]string{}
					for {
						part, err := reader.NextPart()
						if errors.Is(err, io.EOF) {
							break
						}
						if err != nil {
							t.Fatal(err)
						}
						value, err := io.ReadAll(part)
						if err != nil {
							t.Fatal(err)
						}
						if _, duplicate := fields[part.FormName()]; duplicate {
							t.Fatalf("duplicate multipart field %q", part.FormName())
						}
						if part.FormName() == "image" && part.FileName() == "" {
							t.Error("image must remain a file part")
						}
						fields[part.FormName()] = string(value)
					}
					value, present := fields["safety_identifier"]
					if present != tc.present || value != tc.want {
						t.Errorf("safety_identifier = %q (present=%v); want %q (present=%v)", value, present, tc.want, tc.present)
					}
					for key, want := range map[string]string{"user": "synthetic-user", "prompt": "Synthetic fixture", "image": "synthetic-image", "n": "null"} {
						if value, present := fields[key]; !present || value != want {
							t.Errorf("%s = %q (present=%v); want %q", key, value, present, want)
						}
					}
					if value, present := fields["stream"]; present != streaming || (streaming && value != "true") {
						t.Errorf("stream = %q (present=%v); want present=%v", value, present, streaming)
					}
					return imageEditRequestTestResponse(r, streaming), nil
				})))
				body := openai.ImageEditParams{
					Prompt:           "Synthetic fixture",
					Image:            openai.ImageEditParamsImageUnion{OfFile: strings.NewReader("synthetic-image")},
					User:             openai.String("synthetic-user"),
					N:                param.Null[int64](),
					SafetyIdentifier: tc.value,
				}
				runImageEditRequestTest(t, client, streaming, body)
				if calls != 1 || body.SafetyIdentifier != tc.value || !param.IsNull(body.N) {
					t.Fatalf("calls=%d or caller parameters changed", calls)
				}
			})
		}
		t.Run(method+"/body_override", func(t *testing.T) {
			calls := 0
			client := openai.NewClient(option.WithAPIKey("synthetic-key"), option.WithBaseURL("https://example.test/v1"), option.WithMaxRetries(0), option.WithHTTPClient(imageErrorHTTPDoer(func(r *http.Request) (*http.Response, error) {
				calls++
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				if string(body) != "synthetic-override" || r.Header.Get("Content-Type") != "application/x-image-test" {
					t.Fatal("request body override was not preserved")
				}
				return imageEditRequestTestResponse(r, streaming), nil
			})))
			body := openai.ImageEditParams{Prompt: "Synthetic fixture", Image: openai.ImageEditParamsImageUnion{OfFile: strings.NewReader("synthetic-image")}, SafetyIdentifier: param.Null[string]()}
			runImageEditRequestTest(t, client, streaming, body, option.WithRequestBody("application/x-image-test", []byte("synthetic-override")))
			if calls != 1 || !param.IsNull(body.SafetyIdentifier) {
				t.Fatalf("calls=%d or caller parameters changed", calls)
			}
		})
	}
}

func imageEditRequestTestResponse(r *http.Request, streaming bool) *http.Response {
	body, contentType := `{"created":1,"data":[]}`, "application/json"
	if streaming {
		body = "event: image_edit.completed\ndata: {\"type\":\"image_edit.completed\",\"b64_json\":\"c3ludGhldGlj\",\"created_at\":1}\n\n"
		contentType = "text/event-stream"
	}
	return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {contentType}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}
}

func runImageEditRequestTest(t *testing.T, client openai.Client, streaming bool, body openai.ImageEditParams, opts ...option.RequestOption) {
	t.Helper()
	if !streaming {
		if _, err := client.Images.Edit(context.Background(), body, opts...); err != nil {
			t.Fatal(err)
		}
		return
	}
	stream := client.Images.EditStreaming(context.Background(), body, opts...)
	defer func() {
		if err := stream.Close(); err != nil {
			t.Error(err)
		}
	}()
	if !stream.Next() {
		t.Fatalf("missing completion event: %v", stream.Err())
	}
	if stream.Current().Type != "image_edit.completed" || stream.Next() || stream.Err() != nil {
		t.Fatal("expected exactly one successful image-edit completion event")
	}
}
