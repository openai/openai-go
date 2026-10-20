// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package responses

import (
	"bytes"
	"github.com/openai/openai-go/v3/internal/apijson"
)

func (r *ComputerActionListParam) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*r = nil
		return nil
	}
	return apijson.UnmarshalRoot(data, r)
}

func (r *ResponseFunctionCallOutputItemListParam) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*r = nil
		return nil
	}
	return apijson.UnmarshalRoot(data, r)
}

func (r *ResponseInputParam) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*r = nil
		return nil
	}
	return apijson.UnmarshalRoot(data, r)
}

func (r *ResponseInputMessageContentListParam) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*r = nil
		return nil
	}
	return apijson.UnmarshalRoot(data, r)
}
