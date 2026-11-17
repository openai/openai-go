package apijson_test

import (
	"testing"

	"github.com/openai/openai-go/v3/internal/apijson"
)

func TestMarshalMapRejectsEmptyKey(t *testing.T) {
	_, err := apijson.MarshalRoot(map[string]string{"": "value"})
	if err == nil || err.Error() != "path cannot be empty" {
		t.Fatalf("expected legacy empty-key error, got %v", err)
	}
}
