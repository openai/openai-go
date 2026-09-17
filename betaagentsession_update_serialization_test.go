package openai_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
)

func TestBetaAgentSessionUpdateSerialization(t *testing.T) {
	resetReasoning := openai.BetaAgentSessionUpdateParamsAgentReasoning{}
	resetReasoning.SetExtraFields(map[string]any{"effort": nil})
	resetAgent := openai.BetaAgentSessionUpdateParamsAgent{Reasoning: resetReasoning}
	resetAgent.SetExtraFields(map[string]any{"service_tier": nil})
	clearMetadata := openai.BetaAgentSessionUpdateParams{}
	clearMetadata.SetExtraFields(map[string]any{"metadata": nil})

	tests := []struct {
		name   string
		params openai.BetaAgentSessionUpdateParams
		want   map[string]any
	}{
		{name: "omitted settings", want: map[string]any{}},
		{
			name: "empty model",
			params: openai.BetaAgentSessionUpdateParams{
				Agent: openai.BetaAgentSessionUpdateParamsAgent{Model: openai.String("")},
			},
			want: map[string]any{"agent": map[string]any{"model": ""}},
		},
		{
			name: "populated settings",
			params: openai.BetaAgentSessionUpdateParams{
				Agent: openai.BetaAgentSessionUpdateParamsAgent{
					Model: openai.String("gpt-5"),
					Reasoning: openai.BetaAgentSessionUpdateParamsAgentReasoning{
						Effort: "low",
					},
					ServiceTier: "priority",
				},
			},
			want: map[string]any{"agent": map[string]any{
				"model": "gpt-5", "reasoning": map[string]any{"effort": "low"}, "service_tier": "priority",
			}},
		},
		{
			name:   "explicit null resets",
			params: openai.BetaAgentSessionUpdateParams{Agent: resetAgent},
			want:   map[string]any{"agent": map[string]any{"reasoning": map[string]any{"effort": nil}, "service_tier": nil}},
		},
		{
			name:   "metadata update",
			params: openai.BetaAgentSessionUpdateParams{Metadata: map[string]string{"purpose": "test"}},
			want:   map[string]any{"metadata": map[string]any{"purpose": "test"}},
		},
		{name: "null metadata", params: clearMetadata, want: map[string]any{"metadata": nil}},
		{
			name:   "empty metadata",
			params: openai.BetaAgentSessionUpdateParams{Metadata: map[string]string{}},
			want:   map[string]any{"metadata": map[string]any{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("json.Marshal(%+v) error = %v, want nil", tt.params, err)
			}
			var got map[string]any
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal(%q) error = %v, want nil", data, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("json.Marshal(%+v) = %s, want %+v", tt.params, data, tt.want)
			}
		})
	}
}
