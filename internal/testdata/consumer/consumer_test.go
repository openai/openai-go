package consumer_test

import (
	"testing"

	"github.com/openai/openai-go/v3"
)

func TestCorePackageCompiles(t *testing.T) {
	if openai.ChatModelGPT4o == "" {
		t.Error("ChatModelGPT4o = empty, want a model identifier")
	}
}
