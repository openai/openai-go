package responses_test

import (
	"testing"

	"github.com/openai/openai-go/v3/responses"
)

func TestNestedActionGettersPreserveReleasedFindType(t *testing.T) {
	find := responses.ResponseFunctionWebSearchActionFindParam{
		Type: "find_in_page", Pattern: "needle", URL: "https://example.com",
	}
	item := responses.ResponseInputItemUnionParam{
		OfWebSearchCall: &responses.ResponseFunctionWebSearchParam{
			Action: responses.ResponseFunctionWebSearchActionUnionParam{OfFind: &find},
		},
	}
	action := item.GetAction()
	if got := action.GetPattern(); got != &find.Pattern {
		t.Fatalf("pattern = %v, want original pattern pointer %v", got, &find.Pattern)
	}
	if got := action.GetType(); got == nil || *got != "find_in_page" {
		t.Fatalf("action type = %v, want find_in_page", got)
	}
}
