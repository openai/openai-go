package openai_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestNamedMultiGraderPublicClient(t *testing.T) {
	leaf := func(name string) openai.MultiGraderGraderUnionParam {
		return openai.MultiGraderGraderUnionParam{OfStringCheckGrader: &openai.StringCheckGraderParam{
			Name: name, Input: "answer", Reference: "answer", Operation: openai.StringCheckGraderOperationEq,
		}}
	}
	multi := openai.MultiGraderParam{
		Name: "combined", CalculateOutput: "(first + second) / 2",
		Graders: map[string]openai.MultiGraderGraderUnionParam{"first": leaf("first"), "second": leaf("second")},
	}
	var requests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if r.Method != "POST" || r.URL.Path != "/fine_tuning/alpha/graders/validate" {
			t.Errorf("unexpected route: %s %s", r.Method, r.URL.Path)
		}
		var body struct {
			Grader struct {
				Type            string `json:"type"`
				CalculateOutput string `json:"calculate_output"`
				Graders         map[string]struct {
					Type string `json:"type"`
					Name string `json:"name"`
				} `json:"graders"`
			} `json:"grader"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if body.Grader.Type != "multi" || body.Grader.CalculateOutput != "(first + second) / 2" || len(body.Grader.Graders) != 2 {
			t.Errorf("named grader request lost: %+v", body.Grader)
		}
		for _, key := range []string{"first", "second"} {
			if g, ok := body.Grader.Graders[key]; !ok || g.Name != key || g.Type != "string_check" {
				t.Errorf("grader %q lost: %+v", key, g)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"grader":{"type":"multi","name":"combined","calculate_output":"(first + second) / 2","graders":{"first":{"type":"string_check","name":"first","input":"answer","reference":"answer","operation":"eq"},"second":{"type":"python","name":"second","source":"return 1"}}}}`))
	}))
	defer server.Close()
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
	response, err := client.FineTuning.Alpha.Graders.Validate(t.Context(), openai.FineTuningAlphaGraderValidateParams{
		Grader: openai.FineTuningAlphaGraderValidateParamsGraderUnion{OfMultiGrader: &multi},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("requests = %d, want 1", got)
	}
	got := response.Grader.AsMultiGrader()
	if len(got.Graders) != 2 || got.CalculateOutput != multi.CalculateOutput {
		t.Fatalf("named graders lost: %+v", got)
	}
	if got.Graders["first"].AsStringCheck().Name != "first" || got.Graders["second"].AsPython().Source != "return 1" {
		t.Fatalf("named grader variants lost: %+v", got.Graders)
	}
}
