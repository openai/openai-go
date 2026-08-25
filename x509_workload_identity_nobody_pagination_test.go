package openai_test

import (
	"io"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestX509WorkloadIdentityAutoPaginationPreservesHTTPNoBody(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	var requests atomic.Int32
	api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.ContentLength != 0 {
			t.Errorf("empty paginated request content length = %d, want zero", request.ContentLength)
		}
		w.Header().Set("Content-Type", "application/json")
		if requests.Add(1) == 1 {
			_, _ = io.WriteString(w, `{"data":[{"id":"synthetic-first"}],"has_more":true}`)
			return
		}
		_, _ = io.WriteString(w, `{"data":[{"id":"synthetic-second"}],"has_more":false}`)
	})
	client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetries(0))
	pages := client.Batches.ListAutoPaging(
		t.Context(),
		openai.BatchListParams{},
		option.WithRequestBody("application/json", http.NoBody),
	)
	var items int
	for pages.Next() {
		items++
	}
	if err := pages.Err(); err != nil || items != 2 {
		t.Fatalf("empty-body auto-pagination returned items=%d, error=%v", items, err)
	}
	if requests.Load() != 2 || len(issuer.requests()) != 1 {
		t.Errorf("empty-body auto-pagination API/issuer calls = %d/%d, want 2/1",
			requests.Load(), len(issuer.requests()))
	}
}
