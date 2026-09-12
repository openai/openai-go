// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package realtime

import (
	"bytes"
	"github.com/openai/openai-go/v3/internal/apijson"
)

func (r *RealtimeToolsConfigParam) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*r = nil
		return nil
	}
	return apijson.UnmarshalRoot(data, r)
}
