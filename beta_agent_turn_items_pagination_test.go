package openai_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestRootTurnItemsFollowResponseCursor(t *testing.T) {
	for _, finalItemType := range []string{"message", "future_item_type"} {
		t.Run(finalItemType, func(t *testing.T) {
			requests := 0
			transport := paginationHTTPDoerFunc(func(req *http.Request) (*http.Response, error) {
				index := requests
				requests++
				if index > 1 {
					return nil, fmt.Errorf("unexpected request after terminal page")
				}
				if req.URL.Path != "/v1/agents/sessions/session_test/turns/turn_test/items" {
					return nil, fmt.Errorf("unexpected turn-items request path: %s", req.URL.Path)
				}
				query := req.URL.Query()
				if query.Get("limit") != "1" || query.Get("order") != "asc" {
					return nil, fmt.Errorf("pagination lost requested limit or order: %s", req.URL.RawQuery)
				}
				var body string
				if index == 0 {
					if query.Has("after") {
						return nil, fmt.Errorf("first request unexpectedly has an after cursor")
					}
					body = fmt.Sprintf(`{"object":"list","data":[{"type":%q,"id":null,"turn_id":"turn_test","role":"user","status":"completed","content":[]}],"last_id":"response:cursor/+first","has_more":true}`, finalItemType)
				} else {
					if query.Get("after") != "response:cursor/+first" {
						return nil, fmt.Errorf("second request must use the response last_id: got %q", query.Get("after"))
					}
					body = `{"object":"list","data":[{"type":"message","id":"terminal_item","turn_id":"turn_test","role":"user","status":"completed","content":[]}],"last_id":"valid_but_terminal_cursor","has_more":false}`
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(strings.NewReader(body)),
					Request:    req,
				}, nil
			})
			client := openai.NewClient(
				option.WithBaseURL("https://sdk-test.example/v1"),
				option.WithAPIKey("synthetic"),
				option.WithMaxRetries(0),
				option.WithHTTPClient(transport),
			)
			page := client.Beta.Agents.Sessions.Turns.Items.ListAutoPaging(
				context.Background(), "session_test", "turn_test",
				openai.BetaAgentSessionTurnItemListParams{
					Limit: openai.Int(1),
					Order: openai.BetaAgentSessionTurnItemListParamsOrderAsc,
				},
			)
			var ids []string
			for page.Next() {
				ids = append(ids, page.Current().ID)
			}
			if err := page.Err(); err != nil {
				t.Fatal(err)
			}
			if want := []string{"", "terminal_item"}; !reflect.DeepEqual(ids, want) {
				t.Fatalf("IDs = %v, want %v", ids, want)
			}
			if requests != 2 {
				t.Fatalf("request count = %d, want exactly 2", requests)
			}
		})
	}
}
