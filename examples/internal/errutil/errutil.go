// Package errutil provides safe error messages for runnable examples.
package errutil

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/openai/openai-go/v4"
	"github.com/openai/openai-go/v4/responses"
)

// Message omits URLs, provider fields, and arbitrary error text, including
// wrapper messages. Inspect raw errors only in a trusted debugging environment.
func Message(err error) string {
	var apiErr *openai.Error
	if errors.As(err, &apiErr) && apiErr != nil {
		return fmt.Sprintf("OpenAI API request failed (HTTP %d %s)", apiErr.StatusCode, http.StatusText(apiErr.StatusCode))
	}
	var protocolErr *responses.ResponseProtocolError
	if errors.As(err, &protocolErr) && protocolErr != nil {
		return "Responses WebSocket API error (inspect ResponseProtocolError.Event in a trusted debugging environment)"
	}
	return "Example operation failed"
}

// Status returns a known progress state without reflecting arbitrary provider text.
func Status(status string) string {
	switch status {
	case "uploaded", "processed", "error", "validating_files", "queued", "running",
		"succeeded", "failed", "cancelled", "completed", "incomplete", "in_progress",
		"cancelling", "expired":
		return status
	default:
		return "unknown"
	}
}
