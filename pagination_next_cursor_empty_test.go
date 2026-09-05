package openai_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestNextCursorPaginationFollowsEmptyPages(t *testing.T) {
	tests := []struct {
		name  string
		pages map[string]string
		want  []string
	}{
		{
			name: "empty first page",
			pages: map[string]string{
				"":         `{"data":[],"has_more":true,"next":"cursor-1"}`,
				"cursor-1": `{"data":[{"id":"group-1","created_at":1,"group_type":"group","is_scim_managed":false,"name":"one"}],"has_more":false,"next":null}`,
			},
			want: []string{"group-1"},
		},
		{
			name: "empty intermediate page",
			pages: map[string]string{
				"":         `{"data":[{"id":"group-1","created_at":1,"group_type":"group","is_scim_managed":false,"name":"one"}],"has_more":true,"next":"cursor-1"}`,
				"cursor-1": `{"data":[],"has_more":true,"next":"cursor-2"}`,
				"cursor-2": `{"data":[{"id":"group-2","created_at":2,"group_type":"group","is_scim_managed":false,"name":"two"}],"has_more":false,"next":null}`,
			},
			want: []string{"group-1", "group-2"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			client := openai.NewClient(
				option.WithBaseURL("https://example.com/v1"),
				option.WithAdminAPIKey("test-admin-key"),
				option.WithMaxRetries(0),
				option.WithHTTPClient(paginationHTTPDoerFunc(func(req *http.Request) (*http.Response, error) {
					calls++
					cursor := req.URL.Query().Get("after")
					body, ok := test.pages[cursor]
					if !ok {
						return nil, errors.New("unexpected pagination cursor")
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": {"application/json"}},
						Body:       io.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				})),
			)

			pager := client.Admin.Organization.Groups.ListAutoPaging(
				context.Background(), openai.AdminOrganizationGroupListParams{},
			)
			var got []string
			for pager.Next() {
				got = append(got, pager.Current().ID)
			}
			if err := pager.Err(); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("group IDs = %v, want %v", got, test.want)
			}
			if calls != len(test.pages) {
				t.Fatalf("HTTP calls = %d, want %d", calls, len(test.pages))
			}
		})
	}
}
