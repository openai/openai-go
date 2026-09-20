package openai_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestAgentsIDListPaginationPreservesFilters(t *testing.T) {
	cases := []struct {
		name    string
		path    string
		filters url.Values
		list    func(openai.Client) ([]string, error)
	}{
		{
			"sessions", "/agents/sessions", url.Values{"limit": {"1"}, "order": {"asc"}},
			func(c openai.Client) ([]string, error) {
				return collectAgentIDs(c.Beta.Agents.Sessions.ListAutoPaging(context.Background(), openai.BetaAgentSessionListParams{Limit: openai.Int(1), Order: "asc"}), func(v openai.AgentSession) string { return v.ID })
			},
		},
		{
			"vaults", "/vaults", url.Values{"limit": {"1"}, "status": {"archived"}},
			func(c openai.Client) ([]string, error) {
				return collectAgentIDs(c.Beta.Agents.Vaults.ListAutoPaging(context.Background(), openai.BetaAgentVaultListParams{Limit: openai.Int(1), Status: openai.VaultStatusFilterUnionParam{OfStatus: openai.Opt(openai.VaultStatusArchived)}}), func(v openai.Vault) string { return v.ID })
			},
		},
		{
			"credentials", "/vaults/vault_test/credentials", url.Values{"limit": {"1"}, "status[]": {"active", "archived"}},
			func(c openai.Client) ([]string, error) {
				return collectAgentIDs(c.Beta.Agents.Vaults.Credentials.ListAutoPaging(context.Background(), "vault_test", openai.BetaAgentVaultCredentialListParams{Limit: openai.Int(1), Status: openai.VaultStatusFilterUnionParam{OfArrayOfStatuses: []openai.VaultStatus{openai.VaultStatusActive, openai.VaultStatusArchived}}}), func(v openai.Credential) string { return v.ID })
			},
		},
		{
			"items", "/agents/sessions/session_test/items", url.Values{"limit": {"1"}, "order": {"asc"}},
			func(c openai.Client) ([]string, error) {
				return collectAgentIDs(c.Beta.Agents.Sessions.Items.ListAutoPaging(context.Background(), "session_test", openai.BetaAgentSessionItemListParams{Limit: openai.Int(1), Order: "asc"}), func(v openai.AgentSessionItemUnion) string { return v.ID })
			},
		},
		{
			"turns", "/agents/sessions/session_test/turns", url.Values{"limit": {"1"}, "order": {"asc"}},
			func(c openai.Client) ([]string, error) {
				return collectAgentIDs(c.Beta.Agents.Sessions.Turns.ListAutoPaging(context.Background(), "session_test", openai.BetaAgentSessionTurnListParams{Limit: openai.Int(1), Order: "asc"}), func(v openai.Turn) string { return v.ID })
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				if requests > 2 {
					t.Error("pagination fetched beyond has_more=false")
					http.Error(w, "unexpected page", http.StatusBadRequest)
					return
				}
				if r.Method != http.MethodGet || r.URL.Path != tc.path {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
				if r.Header.Get("OpenAI-Beta") != "agents=v1" {
					t.Error("missing Managed Agents beta header")
				}
				want := url.Values{}
				for key, values := range tc.filters {
					want[key] = values
				}
				if requests == 2 {
					want.Set("after", "resource_1")
				}
				if got := r.URL.Query(); !reflect.DeepEqual(got, want) {
					t.Errorf("page %d query = %v; want %v", requests, got, want)
				}
				id := "resource_1"
				if requests == 2 {
					id = "resource_2"
				}
				w.Header().Set("Content-Type", "application/json")
				page := map[string]any{
					"object": "list", "data": []map[string]string{{"id": id}},
					"first_id": id, "last_id": id, "has_more": requests == 1,
				}
				if tc.name == "items" {
					page["after"] = id
				}
				if tc.name == "turns" {
					page["last_id"] = id
				}
				if err := json.NewEncoder(w).Encode(page); err != nil {
					t.Errorf("write page: %v", err)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL(server.URL), option.WithMaxRetries(0))
			ids, err := tc.list(client)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(ids, []string{"resource_1", "resource_2"}) || requests != 2 {
				t.Fatalf("IDs = %v, requests = %d; want two IDs from two pages", ids, requests)
			}
		})
	}
}

func collectAgentIDs[T any](pager interface {
	Next() bool
	Current() T
	Err() error
}, id func(T) string) ([]string, error) {
	var ids []string
	for pager.Next() {
		ids = append(ids, id(pager.Current()))
	}
	return ids, pager.Err()
}
