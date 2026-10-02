// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package requestconfig

import (
	"mime/multipart"
	"net/url"

	"github.com/openai/openai-go/v3/internal/apiform"
	"github.com/openai/openai-go/v3/internal/apiquery"
)

// Keep legacy marshaling available until the SDK runtime supports writer-based
// serialization. The marker preserves custom marshaler dispatch on generic calls.
type multipartWriter interface {
	apiform.Marshaler
	MarshalMultipartTo(*multipart.Writer) error
}

// MultipartBody selects writer-based serialization for generated service calls.
func MultipartBody(body multipartWriter) any { return multipartRequestBody{body} }

type multipartRequestBody struct{ multipartWriter }

func (b multipartRequestBody) URLQuery() (url.Values, error) {
	if query, ok := b.multipartWriter.(apiquery.Queryer); ok {
		return query.URLQuery()
	}
	return nil, nil
}
