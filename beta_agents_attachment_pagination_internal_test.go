package openai

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3/packages/pagination"
)

func TestBetaAgentAttachmentPaginationKeepsNativeTermination(t *testing.T) {
	for _, body := range []string{
		`{"data":[],"has_more":true}`,
		`{"data":[{"id":"repeated"}],"has_more":false}`,
	} {
		t.Run(body, func(t *testing.T) {
			var page pagination.CursorPage[Turn]
			if err := json.Unmarshal([]byte(body), &page); err != nil {
				t.Fatal(err)
			}
			// These terminal pages stop without using request configuration, even
			// when their cursor has already been used earlier in the traversal.
			next, err := betaAgentNextAttachmentPage(&page, "repeated", map[string]struct{}{"repeated": {}})
			if next != nil || err != nil {
				t.Fatalf("native termination changed: next=%v err=%v", next, err)
			}
		})
	}
}
