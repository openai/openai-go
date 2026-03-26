// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package apierror

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/packages/respjson"
)

// Error represents an error that originates from the API, i.e. when a request is
// made and the API returns a response with a HTTP status code. Other errors are
// not wrapped by this SDK.
type Error struct {
	Code    string `json:"code" api:"required"`
	Message string `json:"message" api:"required"`
	Param   string `json:"param" api:"required"`
	Type    string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		Param       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
	StatusCode int
	Request    *http.Request
	Response   *http.Response
}

// Returns the unmodified JSON received from the API
func (r Error) RawJSON() string { return r.JSON.raw }
func (r *Error) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (r *Error) Error() string {
	if r.Request == nil && r.Response == nil {
		return fmt.Sprintf("api error: %s", r.JSON.raw)
	}
	if r.Request == nil {
		return fmt.Sprintf("%d %s %s", r.Response.StatusCode, http.StatusText(r.Response.StatusCode), r.JSON.raw)
	}
	if r.Response == nil {
		return fmt.Sprintf("%s %q: %s", r.Request.Method, r.Request.URL, r.JSON.raw)
	}
	// Attempt to re-populate the response body
	return fmt.Sprintf("%s %q: %d %s %s", r.Request.Method, r.Request.URL, r.Response.StatusCode, http.StatusText(r.Response.StatusCode), r.JSON.raw)
}

func (r *Error) DumpRequest(body bool) []byte {
	if r.Request == nil {
		return nil
	}
	if r.Request.GetBody != nil {
		r.Request.Body, _ = r.Request.GetBody()
	}
	out, _ := httputil.DumpRequestOut(r.Request, body)
	return out
}

func (r *Error) DumpResponse(body bool) []byte {
	if r.Response == nil {
		return nil
	}
	out, _ := httputil.DumpResponse(r.Response, body)
	return out
}
