package auth_test

import (
	"net/http"
	"testing"

	"github.com/openai/openai-go/v3/internal/testutil"
)

func nativeX509HTTPClient(t testing.TB, roundTrip func(*http.Request) (*http.Response, error)) *http.Client {
	return testutil.NewNativeX509HTTPClient(t, roundTrip)
}
