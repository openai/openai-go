package openai_test

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/responses"
)

func TestReleasedConstantConstructors(t *testing.T) {
	breakpoints := []any{
		openai.NewChatCompletionContentPartFilePromptCacheBreakpointParam(),
		openai.NewChatCompletionContentPartImagePromptCacheBreakpointParam(),
		openai.NewChatCompletionContentPartInputAudioPromptCacheBreakpointParam(),
		openai.NewChatCompletionContentPartTextPromptCacheBreakpointParam(),
		openai.NewBetaResponseInputFilePromptCacheBreakpointParam(),
		openai.NewBetaResponseInputFileContentPromptCacheBreakpointParam(),
		openai.NewBetaResponseInputImagePromptCacheBreakpointParam(),
		openai.NewBetaResponseInputImageContentPromptCacheBreakpointParam(),
		openai.NewBetaResponseInputTextPromptCacheBreakpointParam(),
		openai.NewBetaResponseInputTextContentPromptCacheBreakpointParam(),
		responses.NewResponseInputFilePromptCacheBreakpointParam(),
		responses.NewResponseInputFileContentPromptCacheBreakpointParam(),
		responses.NewResponseInputImagePromptCacheBreakpointParam(),
		responses.NewResponseInputImageContentPromptCacheBreakpointParam(),
		responses.NewResponseInputTextPromptCacheBreakpointParam(),
		responses.NewResponseInputTextContentPromptCacheBreakpointParam(),
	}
	for _, breakpoint := range breakpoints {
		data, err := json.Marshal(breakpoint)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != `{"mode":"explicit"}` {
			t.Fatalf("%T marshaled as %s", breakpoint, data)
		}
	}
	data, err := json.Marshal(live.NewClientDelegationParam())
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"type":"client"}` {
		t.Fatalf("delegation marshaled as %s", data)
	}
}
