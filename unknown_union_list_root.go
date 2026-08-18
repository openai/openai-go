package openai

import (
	"bytes"

	"github.com/openai/openai-go/v3/internal/apijson"
)

func (r *BetaComputerActionListParam) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*r = nil
		return nil
	}
	return apijson.UnmarshalRoot(data, r)
}

func (r *BetaResponseFunctionCallOutputItemListParam) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*r = nil
		return nil
	}
	return apijson.UnmarshalRoot(data, r)
}

func (r *BetaResponseInputParam) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*r = nil
		return nil
	}
	return apijson.UnmarshalRoot(data, r)
}

func (r *BetaResponseInputMessageContentListParam) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*r = nil
		return nil
	}
	return apijson.UnmarshalRoot(data, r)
}
