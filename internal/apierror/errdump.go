// go:build !tinygo

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package apierror

import (
	"net/http/httputil"
)

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
