package openai

import (
	"net/http"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/internal/requestconfig"
)

func TestGetPollIntervalValidatesRemoteDelay(t *testing.T) {
	tests := map[string]struct {
		raw  *http.Response
		want time.Duration
	}{
		"missing response": {
			want: time.Second,
		},
		"missing header": {
			raw:  &http.Response{Header: make(http.Header)},
			want: time.Second,
		},
		"positive value": {
			raw:  &http.Response{Header: http.Header{"Openai-Poll-After-Ms": {"250"}}},
			want: 250 * time.Millisecond,
		},
		"zero": {
			raw:  &http.Response{Header: http.Header{"Openai-Poll-After-Ms": {"0"}}},
			want: time.Second,
		},
		"negative": {
			raw:  &http.Response{Header: http.Header{"Openai-Poll-After-Ms": {"-1"}}},
			want: time.Second,
		},
		"large value": {
			raw:  &http.Response{Header: http.Header{"Openai-Poll-After-Ms": {"60000"}}},
			want: requestconfig.DefaultMaxServerDelay,
		},
		"overflow": {
			raw:  &http.Response{Header: http.Header{"Openai-Poll-After-Ms": {"999999999999999999999999"}}},
			want: requestconfig.DefaultMaxServerDelay,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := getPollInterval(test.raw); got != test.want {
				t.Fatalf("getPollInterval() = %s, want %s", got, test.want)
			}
		})
	}
}

func TestPollIntervalAllowsExplicitLongerDelay(t *testing.T) {
	if got := pollInterval(9000, nil); got != 9*time.Second {
		t.Fatalf("pollInterval() = %s, want 9s", got)
	}
}
