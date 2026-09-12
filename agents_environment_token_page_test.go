package openai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestEnvironmentFilesTokenPagination(t *testing.T) {
	const token = "synthetic:token/+="
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		index := requests
		requests++
		if index >= 2 {
			t.Error("requested past terminal page")
			http.Error(w, "extra page", 400)
			return
		}
		if r.URL.Path != "/v1/agents/environments/env_test/files" {
			t.Errorf("path: %s", r.URL.Path)
		}
		want := url.Values{"path": {"/workspace/test"}, "order": {"asc"}, "limit": {"1"}}
		if index > 0 {
			want.Set("page", token)
		}
		if !reflect.DeepEqual(r.URL.Query(), want) {
			t.Errorf("query: got %v want %v", r.URL.Query(), want)
		}
		if r.Header.Get("X-Pagination-Test") != "preserved" {
			t.Error("lost request header")
		}
		if r.Header.Get("OpenAI-Beta") != "agents=v1" {
			t.Error("lost beta header")
		}
		var next any
		if index == 0 {
			next = token
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"object": "page", "has_more": index == 0, "next": next,
			"data": []map[string]any{{"object": "agent.environment.file", "environment_id": "env_test", "path": fmt.Sprintf("/workspace/test/%d.txt", index), "size_bytes": 1}},
		})
	}))
	defer server.Close()
	client := openai.NewClient(option.WithAPIKey("synthetic"), option.WithBaseURL(server.URL+"/v1"), option.WithMaxRetries(0))
	pager := client.Beta.Agents.Environments.Files.ListAutoPaging(context.Background(), "env_test", openai.BetaAgentEnvironmentFileListParams{
		Path: openai.String("/workspace/test"), Order: "asc", Limit: openai.Int(1),
	}, option.WithHeader("X-Pagination-Test", "preserved"))
	var paths []string
	for pager.Next() {
		paths = append(paths, pager.Current().Path)
	}
	if err := pager.Err(); err != nil {
		t.Fatal(err)
	}
	if want := []string{"/workspace/test/0.txt", "/workspace/test/1.txt"}; !reflect.DeepEqual(paths, want) {
		t.Errorf("paths: %v", paths)
	}
	if requests != 2 {
		t.Errorf("requests: %d", requests)
	}
}
