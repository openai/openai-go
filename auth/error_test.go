package auth_test

import (
	"testing"

	"github.com/openai/openai-go/v3/auth"
)

func TestOAuthErrorFormattingPreservesEmptyDescriptionCompatibility(t *testing.T) {
	withoutDescription := (&auth.OAuthError{StatusCode: 401, ErrorCode: "invalid_grant"}).Error()
	if withoutDescription != "OAuth error (status 401): invalid_grant - " {
		t.Errorf("redacted OAuth error = %q", withoutDescription)
	}

	withDescription := (&auth.OAuthError{
		StatusCode:       401,
		ErrorCode:        "invalid_grant",
		ErrorDescription: "synthetic description",
	}).Error()
	if withDescription != "OAuth error (status 401): invalid_grant - synthetic description" {
		t.Errorf("described OAuth error = %q", withDescription)
	}
}
