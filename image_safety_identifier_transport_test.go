package openai_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/respjson"
)

func TestImageSafetyIdentifierTransport(t *testing.T) {
	for _, route := range []string{"generate", "edit", "generate_stream", "edit_stream"} {
		for _, shape := range []string{"value", "null", "omitted"} {
			t.Run(route+"/"+shape, func(t *testing.T) {
				field := `,"safety_identifier":"synthetic-user"`
				if shape == "null" {
					field = `,"safety_identifier":null`
				} else if shape == "omitted" {
					field = ""
				}
				editing := strings.HasPrefix(route, "edit")
				streaming := strings.HasSuffix(route, "_stream")
				path, event := "/v1/images/generations", "image_generation.completed"
				if editing {
					path, event = "/v1/images/edits", "image_edit.completed"
				}
				body, contentType := `{"created":1,"data":[]`+field+`}`, "application/json"
				if streaming {
					body = "event: " + event + "\ndata: " + `{"type":"` + event + `","b64_json":"c3ludGhldGlj","created_at":1` + field + "}\n\n"
					contentType = "text/event-stream"
				}
				calls := 0
				client := openai.NewClient(option.WithAPIKey("synthetic-key"), option.WithBaseURL("https://example.test/v1"), option.WithMaxRetries(0), option.WithHTTPClient(imageErrorHTTPDoer(func(r *http.Request) (*http.Response, error) {
					calls++
					if r.Method != http.MethodPost || r.URL.Path != path {
						t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
					}
					return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {contentType}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
				})))
				var value string
				var presence respjson.Field
				generate := openai.ImageGenerateParams{Prompt: "Synthetic fixture", SafetyIdentifier: openai.String("synthetic-user")}
				edit := openai.ImageEditParams{Prompt: "Synthetic fixture", Image: openai.ImageEditParamsImageUnion{OfFile: strings.NewReader("synthetic-image")}, SafetyIdentifier: openai.String("synthetic-user")}
				switch route {
				case "generate", "edit":
					var response *openai.ImagesResponse
					var err error
					if editing {
						response, err = client.Images.Edit(context.Background(), edit)
					} else {
						response, err = client.Images.Generate(context.Background(), generate)
					}
					if err != nil {
						t.Fatal(err)
					}
					value, presence = response.SafetyIdentifier, response.JSON.SafetyIdentifier
				case "generate_stream":
					stream := client.Images.GenerateStreaming(context.Background(), generate)
					if !stream.Next() {
						t.Fatalf("missing completion event: %v", stream.Err())
					}
					completed := stream.Current().AsImageGenerationCompleted()
					value, presence = completed.SafetyIdentifier, completed.JSON.SafetyIdentifier
					if stream.Next() || stream.Err() != nil {
						t.Fatal("expected exactly one successful completion event")
					}
					if err := stream.Close(); err != nil {
						t.Fatal(err)
					}
				case "edit_stream":
					stream := client.Images.EditStreaming(context.Background(), edit)
					if !stream.Next() {
						t.Fatalf("missing completion event: %v", stream.Err())
					}
					completed := stream.Current().AsImageEditCompleted()
					value, presence = completed.SafetyIdentifier, completed.JSON.SafetyIdentifier
					if stream.Next() || stream.Err() != nil {
						t.Fatal("expected exactly one successful completion event")
					}
					if err := stream.Close(); err != nil {
						t.Fatal(err)
					}
				}
				want, raw := "", ""
				if shape == "value" {
					want, raw = "synthetic-user", `"synthetic-user"`
				} else if shape == "null" {
					raw = "null"
				}
				if calls != 1 || value != want || presence.Raw() != raw || presence.Valid() != (shape == "value") {
					t.Fatalf("calls=%d value=%q raw=%q valid=%v; want 1, %q, %q, %v", calls, value, presence.Raw(), presence.Valid(), want, raw, shape == "value")
				}
			})
		}
	}
}
