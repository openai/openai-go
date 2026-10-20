package openai_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
)

func TestProjectAPIKeyUpdateScopeSerialization(t *testing.T) {
	tests := []struct {
		name   string
		params openai.AdminOrganizationProjectAPIKeyUpdateParams
		want   map[string]any
	}{
		{
			name:   "existing expiry-only literal",
			params: openai.AdminOrganizationProjectAPIKeyUpdateParams{ExpiresInSeconds: 60},
			want:   map[string]any{"expires_in_seconds": float64(60)},
		},
		{
			name:   "explicit all permissions",
			params: openai.AdminOrganizationProjectAPIKeyUpdateParams{Scopes: []string{}},
			want:   map[string]any{"scopes": []any{}},
		},
		{
			name:   "scope-only read access",
			params: openai.AdminOrganizationProjectAPIKeyUpdateParams{Scopes: []string{"api.all.read"}},
			want:   map[string]any{"scopes": []any{"api.all.read"}},
		},
		{
			name:   "combined update",
			params: openai.AdminOrganizationProjectAPIKeyUpdateParams{ExpiresInSeconds: 60, Scopes: []string{}},
			want:   map[string]any{"expires_in_seconds": float64(60), "scopes": []any{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("json.Marshal(%+v) error = %v, want nil", tt.params, err)
			}
			var got map[string]any
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatalf("json.Unmarshal(%s) error = %v, want nil", body, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("json.Marshal(%+v) decoded = %#v, want %#v", tt.params, got, tt.want)
			}
		})
	}
}
