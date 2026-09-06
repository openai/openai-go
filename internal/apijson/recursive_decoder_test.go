package apijson

import (
	"strings"
	"testing"

	"github.com/openai/openai-go/v3/packages/respjson"
)

var recursiveNodeUnmarshalCalls int

type recursiveNode struct {
	Children []recursiveNode `json:"children"`
	Value    string          `json:"value"`
	JSON     struct {
		Children respjson.Field
		Value    respjson.Field
		raw      string
	} `json:"-"`
}

func (r *recursiveNode) UnmarshalJSON(data []byte) error {
	recursiveNodeUnmarshalCalls++
	return UnmarshalRoot(data, r)
}

func TestRecursiveGeneratedTypeRetainsDecoderState(t *testing.T) {
	const depth = 1000
	recursiveNodeUnmarshalCalls = 0

	var payload strings.Builder
	for range depth {
		payload.WriteString(`{"children":[`)
	}
	payload.WriteString(`{"value":"leaf"}`)
	for range depth {
		payload.WriteString(`]}`)
	}

	var got struct {
		Root recursiveNode `json:"root"`
	}
	if err := Unmarshal([]byte(`{"root":`+payload.String()+`}`), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if recursiveNodeUnmarshalCalls != 0 {
		t.Fatalf("recursive UnmarshalJSON calls = %d, want 0", recursiveNodeUnmarshalCalls)
	}

	node := &got.Root
	for level := 0; level < depth; level++ {
		if node.JSON.raw == "" {
			t.Fatalf("level %d did not retain raw JSON metadata", level)
		}
		if len(node.Children) != 1 {
			t.Fatalf("level %d children = %d, want 1", level, len(node.Children))
		}
		node = &node.Children[0]
	}
	if node.Value != "leaf" {
		t.Fatalf("leaf value = %q, want leaf", node.Value)
	}
}
