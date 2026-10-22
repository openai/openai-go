package openai_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
)

func TestDecisionPromptCacheWireContract(t *testing.T) {
	for _, mode := range []string{"omitted", "explicit", "null"} {
		t.Run(mode, func(t *testing.T) {
			text := openai.DecisionInputPartParamOfInputText("hello")
			image := openai.DecisionInputPartParamOfInputImage("https://example.test/image.png")
			p := openai.DecisionNewParams{Model: "test-model", Questions: []openai.DecisionNewParamsQuestionUnion{{OfPredicate: &openai.DecisionNewParamsQuestionPredicate{Instructions: "check"}}}}
			var expected any
			if mode == "explicit" {
				breakpoint := openai.NewDecisionPromptCacheBreakpointParam()
				text.OfInputText.PromptCacheBreakpoint = breakpoint
				image.OfInputImage.PromptCacheBreakpoint = breakpoint
				p.PromptCacheKey = openai.String("test-prefix")
				p.PromptCacheOptions = openai.DecisionNewParamsPromptCacheOptions{Mode: openai.String("explicit"), Ttl: "30m"}
				expected = map[string]any{"mode": "explicit"}
				for _, part := range []openai.DecisionInputPartUnionParam{text, image} {
					var got *openai.DecisionPromptCacheBreakpointParam = part.GetPromptCacheBreakpoint()
					if got == nil || got.Mode != "explicit" {
						t.Fatal("shared breakpoint accessor lost the value")
					}
				}
			} else if mode == "null" {
				text.OfInputText.PromptCacheBreakpoint = param.NullStruct[openai.DecisionPromptCacheBreakpointParam]()
				image.OfInputImage.PromptCacheBreakpoint = param.NullStruct[openai.DecisionPromptCacheBreakpointParam]()
				p.PromptCacheKey = param.Null[string]()
				p.PromptCacheOptions = param.NullStruct[openai.DecisionNewParamsPromptCacheOptions]()
			}
			p.Input.OfDecisionInputMessageArray = []openai.DecisionInputMessageParam{{Role: "user", Content: openai.DecisionInputMessageContentUnionParam{OfParts: []openai.DecisionInputPartUnionParam{text, image}}}}
			requests := make(chan map[string]any, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" || r.URL.Path != "/decisions" {
					t.Errorf("unexpected route: %s %s", r.Method, r.URL.Path)
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				requests <- body
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, `{}`)
			}))
			defer server.Close()
			client := openai.NewClient(option.WithAPIKey("fake-test-key"), option.WithBaseURL(server.URL), option.WithUnsafeAllowHTTP(), option.WithMaxRetries(0))
			if _, err := client.Decisions.New(context.Background(), p); err != nil {
				t.Fatal(err)
			}
			body := <-requests
			for _, key := range []string{"prompt_cache_key", "prompt_cache_options"} {
				value, present := body[key]
				if present != (mode != "omitted") {
					t.Fatalf("%s presence = %v", key, present)
				}
				if mode == "null" && value != nil {
					t.Fatalf("%s must be JSON null: %#v", key, value)
				}
			}
			if mode == "explicit" && (body["prompt_cache_key"] != "test-prefix" || !reflect.DeepEqual(body["prompt_cache_options"], map[string]any{"mode": "explicit", "ttl": "30m"})) {
				t.Fatalf("request cache settings lost: %#v", body)
			}
			parts := body["input"].([]any)[0].(map[string]any)["content"].([]any)
			for _, part := range parts {
				value, present := part.(map[string]any)["prompt_cache_breakpoint"]
				if present != (mode != "omitted") || !reflect.DeepEqual(value, expected) {
					t.Fatalf("breakpoint = %#v, present=%v", value, present)
				}
			}
		})
	}
}

func TestDecisionPromptCacheModeWireContract(t *testing.T) {
	for _, tc := range []struct {
		name    string
		mode    param.Opt[string]
		present bool
		value   any
	}{
		{name: "omitted"},
		{name: "explicit", mode: openai.String("explicit"), present: true, value: "explicit"},
		{name: "implicit", mode: openai.String("implicit"), present: true, value: "implicit"},
		{name: "null", mode: param.Null[string](), present: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requests := make(chan map[string]any, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" || r.URL.Path != "/decisions" {
					t.Errorf("unexpected route: %s %s", r.Method, r.URL.Path)
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				requests <- body
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, `{}`)
			}))
			defer server.Close()
			client := openai.NewClient(option.WithAPIKey("fake-test-key"), option.WithBaseURL(server.URL), option.WithUnsafeAllowHTTP(), option.WithMaxRetries(0))
			_, err := client.Decisions.New(context.Background(), openai.DecisionNewParams{
				Model:              "test-model",
				Questions:          []openai.DecisionNewParamsQuestionUnion{{OfPredicate: &openai.DecisionNewParamsQuestionPredicate{Instructions: "check"}}},
				PromptCacheOptions: openai.DecisionNewParamsPromptCacheOptions{Mode: tc.mode, Ttl: "30m"},
			})
			if err != nil {
				t.Fatal(err)
			}
			body := <-requests
			options, ok := body["prompt_cache_options"].(map[string]any)
			if !ok || options["ttl"] != "30m" {
				t.Fatalf("cache options lost: %#v", body)
			}
			value, present := options["mode"]
			if present != tc.present || !reflect.DeepEqual(value, tc.value) {
				t.Fatalf("mode = %#v, present=%v; want %#v, present=%v", value, present, tc.value, tc.present)
			}
		})
	}
}
