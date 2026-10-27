package openai_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

func TestPromptCacheBreakpointConstructorsRemainCompatible(t *testing.T) {
	breakpoints := []any{
		openai.NewBetaResponseInputFileContentPromptCacheBreakpointParam(),
		openai.NewBetaResponseInputImageContentPromptCacheBreakpointParam(),
		openai.NewBetaResponseInputTextContentPromptCacheBreakpointParam(),
		openai.NewBetaResponseInputFilePromptCacheBreakpointParam(),
		openai.NewBetaResponseInputImagePromptCacheBreakpointParam(),
		openai.NewBetaResponseInputTextPromptCacheBreakpointParam(),
		openai.NewChatCompletionContentPartFilePromptCacheBreakpointParam(),
		openai.NewChatCompletionContentPartImagePromptCacheBreakpointParam(),
		openai.NewChatCompletionContentPartInputAudioPromptCacheBreakpointParam(),
		openai.NewChatCompletionContentPartTextPromptCacheBreakpointParam(),
		responses.NewResponseInputFileContentPromptCacheBreakpointParam(),
		responses.NewResponseInputImageContentPromptCacheBreakpointParam(),
		responses.NewResponseInputTextContentPromptCacheBreakpointParam(),
		responses.NewResponseInputFilePromptCacheBreakpointParam(),
		responses.NewResponseInputImagePromptCacheBreakpointParam(),
		responses.NewResponseInputTextPromptCacheBreakpointParam(),
		openai.NewDecisionPromptCacheBreakpointParam(),
	}
	for _, breakpoint := range breakpoints {
		t.Run(reflect.TypeOf(breakpoint).Name(), func(t *testing.T) {
			assertPromptCacheBreakpointJSON(t, breakpoint, "")
		})
	}
}

func TestPromptCacheBreakpointTTLSerialization(t *testing.T) {
	for _, ttl := range []string{"", "30m", "12h"} {
		t.Run(ttl, func(t *testing.T) {
			beta := openai.NewBetaResponseInputTextContentPromptCacheBreakpointParam()
			beta.Ttl = ttl
			response := responses.NewResponseInputTextContentPromptCacheBreakpointParam()
			response.Ttl = ttl
			primaryBeta := openai.NewBetaResponseInputTextPromptCacheBreakpointParam()
			primaryBeta.Ttl = ttl
			primaryResponse := responses.NewResponseInputTextPromptCacheBreakpointParam()
			primaryResponse.Ttl = ttl
			for _, breakpoint := range []any{beta, response, primaryBeta, primaryResponse} {
				assertPromptCacheBreakpointJSON(t, breakpoint, ttl)
			}
		})
	}
}

func assertPromptCacheBreakpointJSON(t *testing.T, breakpoint any, ttl string) {
	t.Helper()
	body, err := json.Marshal(breakpoint)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"mode": "explicit"}
	if ttl != "" {
		want["ttl"] = ttl
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%T serialized as %s, want %#v", breakpoint, body, want)
	}
}
