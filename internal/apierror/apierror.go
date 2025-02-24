// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package apierror

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/packages/resp"
)

// aliased to make param.APIUnion private when embedding
type apiunion = param.APIUnion

// aliased to make param.APIObject private when embedding
type apiobject = param.APIObject

// Error represents an error that originates from the API, i.e. when a request is
// made and the API returns a response with a HTTP status code. Other errors are
// not wrapped by this SDK.
type Error struct {
	Code    string `json:"code,omitzero,required,nullable"`
	Message string `json:"message,omitzero,required"`
	Param   string `json:"param,omitzero,required,nullable"`
	Type    string `json:"type,omitzero,required"`
	JSON    struct {
		Code    resp.Field
		Message resp.Field
		Param   resp.Field
		Type    resp.Field
		raw     string
	} `json:"-"`
	StatusCode int
	Request    *http.Request
	Response   *http.Response
}

func (r Error) RawJSON() string { return r.JSON.raw }
func (r *Error) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (r *Error) Error() string {
	// Attempt to re-populate the response body
	return fmt.Sprintf("%s %q: %d %s %s", r.Request.Method, r.Request.URL, r.Response.StatusCode, http.StatusText(r.Response.StatusCode), r.JSON.raw)
}

func (r *Error) DumpRequest(body bool) []byte {
	if r.Request.GetBody != nil {
		r.Request.Body, _ = r.Request.GetBody()
	}
	out, _ := httputil.DumpRequestOut(r.Request, body)
	return out
}

func (r *Error) DumpResponse(body bool) []byte {
	out, _ := httputil.DumpResponse(r.Response, body)
	return out
}
