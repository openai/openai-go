package consumer_test

import (
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/conversations"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/realtime"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
	"github.com/openai/openai-go/v3/webhooks"
)

// Existing consumer code must retain access through every exported alias package.
const (
	_ openai.ChatModel        = openai.ChatModelGPT5_1Mini
	_ conversations.ChatModel = conversations.ChatModelGPT5_1Mini
	_ live.ChatModel          = live.ChatModelGPT5_1Mini
	_ realtime.ChatModel      = realtime.ChatModelGPT5_1Mini
	_ responses.ChatModel     = responses.ChatModelGPT5_1Mini
	_ shared.ChatModel        = shared.ChatModelGPT5_1Mini
	_ webhooks.ChatModel      = webhooks.ChatModelGPT5_1Mini
)

func TestCorePackageCompiles(t *testing.T) {
	if openai.ChatModelGPT4o == "" {
		t.Error("ChatModelGPT4o = empty, want a model identifier")
	}
}

// Existing callers can use strings and the exported constants without conversions.
func TestResponseSteeringReasonStringCompatibility(t *testing.T) {
	value := "future_reason"
	var reason responses.ResponseSteerPendingReason = value
	var known string = responses.ResponseSteerPendingReasonWaitingForRequiredInput
	if reason != value || known != "waiting_for_required_input" {
		t.Fatal("steering reason values changed")
	}
}
