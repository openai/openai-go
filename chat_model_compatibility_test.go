package openai_test

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

func TestChatCompletionModelIDs(t *testing.T) {
	configuredID := "deployment-from-config"
	for _, tc := range []struct {
		model shared.ChatModel
		want  string
	}{
		{openai.ChatModelGPT6Luna, "gpt-6-luna"},
		{shared.ChatModelGPT5_6Terra, "gpt-5.6-terra"},
		{configuredID, "deployment-from-config"},
	} {
		t.Run(tc.want, func(t *testing.T) {
			body, err := json.Marshal(openai.ChatCompletionNewParams{Model: tc.model})
			if err != nil {
				t.Fatal(err)
			}
			var request struct {
				Model string `json:"model"`
			}
			if err := json.Unmarshal(body, &request); err != nil {
				t.Fatal(err)
			}
			if request.Model != tc.want {
				t.Errorf("model = %q, want %q", request.Model, tc.want)
			}
		})
	}
}
