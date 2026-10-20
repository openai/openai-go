package consumer_test

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3/responses"
)

func TestCodeInterpreterNetworkPolicyParameterCompatibility(t *testing.T) {
	cases := []struct {
		name   string
		policy responses.ToolCodeInterpreterContainerCodeInterpreterToolAutoNetworkPolicyUnionParam
		want   string
	}{
		{"disabled", responses.ToolCodeInterpreterContainerCodeInterpreterToolAutoNetworkPolicyUnionParam{
			OfDisabled: &responses.ContainerNetworkPolicyDisabledParam{},
		}, `{"type":"disabled"}`},
		{"allowlist", responses.ToolCodeInterpreterContainerCodeInterpreterToolAutoNetworkPolicyUnionParam{
			OfAllowlist: &responses.ContainerNetworkPolicyAllowlistParam{AllowedDomains: []string{"example.com"}},
		}, `{"allowed_domains":["example.com"],"type":"allowlist"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			current := responses.ToolCodeInterpreterContainerCodeInterpreterContainerAutoParam{NetworkPolicy: tc.policy}
			got, err := json.Marshal(current.NetworkPolicy)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tc.want {
				t.Fatalf("network policy = %s, want %s", got, tc.want)
			}
		})
	}
}
