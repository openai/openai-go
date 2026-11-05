package openai_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func TestResponseSVGReasoning(t *testing.T) {
	for _, tc := range []struct {
		name, svg, effort string
		present           bool
	}{
		{"value", `{"reasoning":{"effort":"high"}}`, "high", true},
		{"null", `{"reasoning":null}`, "", false},
		{"omitted", `{}`, "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL("https://example.test/v1"), option.WithHTTPClient(imageErrorHTTPDoer(func(r *http.Request) (*http.Response, error) {
				calls++
				if r.URL.Path != "/v1/responses/resp_fixture" {
					t.Fatalf("unexpected path %s", r.URL.Path)
				}
				body := `{"tools":[{"type":"image_generation","svg":` + tc.svg + `}]}`
				return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			})))
			response, err := client.Responses.Get(context.Background(), "resp_fixture", responses.ResponseGetParams{})
			if err != nil {
				t.Fatal(err)
			}
			if calls != 1 || len(response.Tools) != 1 {
				t.Fatalf("calls=%d tools=%d", calls, len(response.Tools))
			}
			svg := response.Tools[0].AsImageGeneration().Svg
			if svg.Reasoning.Effort != tc.effort || svg.JSON.Reasoning.Valid() != tc.present {
				t.Errorf("reasoning effort=%q present=%v", svg.Reasoning.Effort, svg.JSON.Reasoning.Valid())
			}
		})
	}
}
