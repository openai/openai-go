package openai_test

import (
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

func TestResponseNewParamsTopLevelAlias(t *testing.T) {
	params := openai.ResponseNewParams{
		Model: "gpt-5.2",
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("Say this is a test"),
		},
	}
	var _ responses.ResponseNewParams = params

	if !params.Input.OfString.Valid() || params.Input.OfString.Or("") != "Say this is a test" {
		t.Fatalf("expected input string to round trip, got %#v", params.Input.OfString)
	}
}
