package openai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type multiGraderCaptureTransport struct {
	bodies [][]byte
	paths  []string
}

func (c *multiGraderCaptureTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	c.bodies = append(c.bodies, body)
	c.paths = append(c.paths, r.URL.Path)
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{}`)), Request: r}, nil
}

func namedMultiGrader() openai.MultiGraderParam {
	return openai.MultiGraderParam{
		Name:            "combined",
		CalculateOutput: "(check + similarity + script + score + label) / 5",
		Graders: map[string]openai.MultiGraderGraderUnionParam{
			"check":      {OfStringCheckGrader: &openai.StringCheckGraderParam{Name: "check", Input: "{{ sample.output_text }}", Reference: "yes", Operation: openai.StringCheckGraderOperationEq}},
			"similarity": {OfTextSimilarityGrader: &openai.TextSimilarityGraderParam{Name: "similarity", Input: "{{ sample.output_text }}", Reference: "yes", EvaluationMetric: "fuzzy_match"}},
			"script":     {OfPythonGrader: &openai.PythonGraderParam{Name: "script", Source: "def grade(sample, item):\n    return 1.0"}},
			"score":      {OfScoreModelGrader: &openai.ScoreModelGraderParam{Name: "score", Model: "gpt-4.1", Input: []openai.ScoreModelGraderInputParam{{Role: "user", Content: openai.ScoreModelGraderInputContentUnionParam{OfString: openai.String("{{ sample.output_text }}")}}}}},
			"label":      {OfLabelModelGrader: &openai.LabelModelGraderParam{Name: "label", Model: "gpt-4.1", Labels: []string{"yes", "no"}, PassingLabels: []string{"yes"}, Input: []openai.LabelModelGraderInputParam{{Role: "user", Content: openai.LabelModelGraderInputContentUnionParam{OfString: openai.String("{{ sample.output_text }}")}}}}},
		},
	}
}

func TestMultiGraderNamedValuesAndRoundTrip(t *testing.T) {
	multi := namedMultiGrader()
	body, err := json.Marshal(multi)
	if err != nil {
		t.Fatal(err)
	}
	var response openai.MultiGrader
	if err = json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"check": "openai.StringCheckGrader", "similarity": "openai.TextSimilarityGrader", "script": "openai.PythonGrader", "score": "openai.ScoreModelGrader", "label": "openai.LabelModelGrader"}
	for key, typ := range want {
		if got := fmt.Sprintf("%T", response.Graders[key].AsAny()); got != typ {
			t.Errorf("Graders[%q].AsAny() = %s, want %s", key, got, typ)
		}
	}
	var decoded openai.MultiGraderParam
	if err = json.Unmarshal(body, &decoded); err != nil {
		t.Fatal(err)
	}
	roundTrip, err := json.Marshal(decoded)
	if err != nil {
		t.Fatal(err)
	}
	var before, after any
	if err := json.Unmarshal(body, &before); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(roundTrip, &after); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, after) {
		t.Errorf("typed MultiGrader round trip = %s, want %s", roundTrip, body)
	}
}

func TestMultiGraderNamedRequestBodies(t *testing.T) {
	transport := &multiGraderCaptureTransport{}
	client := openai.NewClient(option.WithBaseURL("http://example.test/v1/"), option.WithAPIKey("test-key"), option.WithHTTPClient(&http.Client{Transport: transport}))
	multi := namedMultiGrader()
	if _, err := client.FineTuning.Alpha.Graders.Run(context.Background(), openai.FineTuningAlphaGraderRunParams{Grader: openai.FineTuningAlphaGraderRunParamsGraderUnion{OfMulti: &multi}, ModelSample: "yes"}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.FineTuning.Alpha.Graders.Validate(context.Background(), openai.FineTuningAlphaGraderValidateParams{Grader: openai.FineTuningAlphaGraderValidateParamsGraderUnion{OfMultiGrader: &multi}}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.FineTuning.Alpha.Graders.Run(context.Background(), openai.FineTuningAlphaGraderRunParams{Grader: openai.FineTuningAlphaGraderRunParamsGraderUnion{OfStringCheck: multi.Graders["check"].OfStringCheckGrader}, ModelSample: "yes"}); err != nil {
		t.Fatal(err)
	}
	wantTypes := map[string]string{"check": "string_check", "similarity": "text_similarity", "script": "python", "score": "score_model", "label": "label_model"}
	for i, body := range transport.bodies {
		var payload struct {
			Grader struct {
				Type    string `json:"type"`
				Graders map[string]struct {
					Type string `json:"type"`
				} `json:"graders"`
				CalculateOutput string `json:"calculate_output"`
				Input           string `json:"input"`
			} `json:"grader"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatal(err)
		}
		if i == 2 {
			if payload.Grader.Type != "string_check" || payload.Grader.Input != "{{ sample.output_text }}" || payload.Grader.Graders != nil {
				t.Errorf("standalone control body = %s, want unchanged string_check", body)
			}
			continue
		}
		if payload.Grader.Type != "multi" || len(payload.Grader.Graders) != 5 || payload.Grader.CalculateOutput != multi.CalculateOutput {
			t.Errorf("multi request body = %s, want five named graders and matching formula", body)
		}
		for key, typ := range wantTypes {
			if got := payload.Grader.Graders[key].Type; got != typ {
				t.Errorf("request %d graders[%q].type = %q, want %q", i, key, got, typ)
			}
		}
	}
	wantPaths := []string{"/v1/fine_tuning/alpha/graders/run", "/v1/fine_tuning/alpha/graders/validate", "/v1/fine_tuning/alpha/graders/run"}
	if !reflect.DeepEqual(transport.paths, wantPaths) {
		t.Errorf("request paths = %v, want %v", transport.paths, wantPaths)
	}
}
