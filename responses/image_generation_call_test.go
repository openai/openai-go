package responses_test

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3/responses"
)

func TestResponseOutputItemImageGenerationCallFields(t *testing.T) {
	const payload = `{
		"id": "ig_123",
		"type": "image_generation_call",
		"status": "completed",
		"result": "iVBORw0KGgo=",
		"action": "generate",
		"background": "opaque",
		"output_format": "png",
		"quality": "high",
		"revised_prompt": "Draw a small red fox.",
		"size": "1024x1024"
	}`

	var item responses.ResponseOutputItemUnion
	if err := json.Unmarshal([]byte(payload), &item); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	image := item.AsImageGenerationCall()
	if image.Action != "generate" {
		t.Errorf("Action = %q, want %q", image.Action, "generate")
	}
	if image.Background != "opaque" {
		t.Errorf("Background = %q, want %q", image.Background, "opaque")
	}
	if image.OutputFormat != "png" {
		t.Errorf("OutputFormat = %q, want %q", image.OutputFormat, "png")
	}
	if image.Quality != "high" {
		t.Errorf("Quality = %q, want %q", image.Quality, "high")
	}
	if image.RevisedPrompt != "Draw a small red fox." {
		t.Errorf("RevisedPrompt = %q, want %q", image.RevisedPrompt, "Draw a small red fox.")
	}
	if image.Size != "1024x1024" {
		t.Errorf("Size = %q, want %q", image.Size, "1024x1024")
	}
}
