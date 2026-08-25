package openai_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/conversations"
	"github.com/openai/openai-go/v3/option"
)

func TestPaginationRejectsNonReplayableStreamingBodies(t *testing.T) {
	for _, test := range []struct {
		name   string
		family string
		auto   bool
	}{
		{name: "cursor manual pagination", family: "cursor"},
		{name: "cursor auto-pagination", family: "cursor", auto: true},
		{name: "conversation manual pagination", family: "conversation"},
		{name: "conversation auto-pagination", family: "conversation", auto: true},
		{name: "next-cursor manual pagination", family: "next"},
		{name: "next-cursor auto-pagination", family: "next", auto: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests.Add(1)
				if body, err := io.ReadAll(request.Body); err != nil || string(body) != "synthetic-private-stream" {
					t.Errorf("initial streamed request body = %q, error = %v", body, err)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w,
					`{"data":[{"id":"synthetic-first","type":"message","role":"user","content":[]}],`+
						`"has_more":true,"last_id":"synthetic-next","next":"synthetic-next"}`,
				)
			}))
			t.Cleanup(server.Close)
			client := openai.NewClient(
				option.WithBaseURL(server.URL+"/v1/"),
				option.WithAPIKey("synthetic-ordinary-key"),
				option.WithAdminAPIKey("synthetic-admin-key"),
				option.WithMaxRetries(0),
			)
			body := option.WithRequestBody("application/octet-stream", strings.NewReader("synthetic-private-stream"))
			var paginationErr error
			switch test.family {
			case "cursor":
				if test.auto {
					pager := client.Batches.ListAutoPaging(t.Context(), openai.BatchListParams{}, body)
					if !pager.Next() || pager.Next() {
						t.Error("cursor auto-pager did not stop safely after its first streamed item")
					}
					paginationErr = pager.Err()
				} else {
					page, err := client.Batches.List(t.Context(), openai.BatchListParams{}, body)
					if err != nil {
						t.Fatalf("load initial cursor page: %v", err)
					}
					_, paginationErr = page.GetNextPage()
				}
			case "conversation":
				if test.auto {
					pager := client.Conversations.Items.ListAutoPaging(
						t.Context(), "synthetic-conversation", conversations.ItemListParams{}, body,
					)
					if !pager.Next() || pager.Next() {
						t.Error("conversation auto-pager did not stop safely after its first streamed item")
					}
					paginationErr = pager.Err()
				} else {
					page, err := client.Conversations.Items.List(
						t.Context(), "synthetic-conversation", conversations.ItemListParams{}, body,
					)
					if err != nil {
						t.Fatalf("load initial conversation page: %v", err)
					}
					_, paginationErr = page.GetNextPage()
				}
			case "next":
				if test.auto {
					pager := client.Admin.Organization.Groups.ListAutoPaging(
						t.Context(), openai.AdminOrganizationGroupListParams{}, body,
					)
					if !pager.Next() || pager.Next() {
						t.Error("next-cursor auto-pager did not stop safely after its first streamed item")
					}
					paginationErr = pager.Err()
				} else {
					page, err := client.Admin.Organization.Groups.List(
						t.Context(), openai.AdminOrganizationGroupListParams{}, body,
					)
					if err != nil {
						t.Fatalf("load initial next-cursor page: %v", err)
					}
					_, paginationErr = page.GetNextPage()
				}
			}
			if paginationErr == nil || !strings.Contains(paginationErr.Error(), "non-replayable request body") {
				t.Errorf("streaming pagination error = %v, want a meaningful non-replayable-body error", paginationErr)
			}
			if paginationErr != nil && strings.Contains(paginationErr.Error(), "synthetic-private") {
				t.Errorf("streaming pagination error disclosed private request data: %q", paginationErr.Error())
			}
			if got := requests.Load(); got != 1 {
				t.Errorf("streaming pagination dispatched %d requests, want only the replayable initial request", got)
			}
		})
	}
}

func TestX509PaginationRejectsNonReplayableStreamingBodies(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	var requests atomic.Int32
	api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		requests.Add(1)
		if body, err := io.ReadAll(request.Body); err != nil || string(body) != "synthetic-private-stream" {
			t.Errorf("initial mutually authenticated streaming body = %q, error = %v", body, err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"data":[{"id":"synthetic-first"}],"has_more":true}`)
	})
	client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetries(0))
	pager := client.Batches.ListAutoPaging(t.Context(), openai.BatchListParams{},
		option.WithRequestBody("application/octet-stream", strings.NewReader("synthetic-private-stream")))
	if !pager.Next() || pager.Next() {
		t.Error("mutually authenticated auto-pager did not stop after its first streamed item")
	}
	if err := pager.Err(); err == nil || !strings.Contains(err.Error(), "non-replayable request body") {
		t.Errorf("mutually authenticated streaming pagination error = %v", err)
	}
	if requests.Load() != 1 || len(issuer.requests()) != 1 {
		t.Errorf("streaming X.509 pagination API/issuer calls = %d/%d, want 1/1",
			requests.Load(), len(issuer.requests()))
	}
}
