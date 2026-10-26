package apierror

import (
	"fmt"
	"log/slog"
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

// Error returns an HTTP status summary suitable for routine error logging.
// Request URLs and provider fields may contain secrets. Inspect the fields,
// RawJSON, DumpRequest, or DumpResponse explicitly for unsanitized diagnostics.
func (r *Error) Error() string {
	if r == nil {
		return "OpenAI API error"
	}
	if text := http.StatusText(r.StatusCode); text != "" {
		return fmt.Sprintf("OpenAI API error: %d %s", r.StatusCode, text)
	}
	return fmt.Sprintf("OpenAI API error: %d", r.StatusCode)
}

// String keeps formatting of copied Error values from exposing raw diagnostics.
func (r Error) String() string { return r.Error() }

// GoString keeps Go-syntax formatting of both values and pointers safe.
func (r Error) GoString() string { return r.Error() }

// Format uses the safe summary whenever fmt invokes the Formatter interface.
// String verbs retain their quoting, encoding, width, and precision behavior.
func (r Error) Format(state fmt.State, verb rune) {
	switch verb {
	case 's', 'q', 'x', 'X':
	default:
		verb = 's'
	}
	_, _ = fmt.Fprintf(state, fmt.FormatString(state, verb), r.Error())
}

// LogValue keeps structured logging of values and pointers on the safe summary.
func (r Error) LogValue() slog.Value { return slog.StringValue(r.Error()) }

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
