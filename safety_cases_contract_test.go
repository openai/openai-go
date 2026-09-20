package openai_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestSafetyCaseRetrievalContract(t *testing.T) {
	tests := []struct {
		name          string
		notice        string
		reason        any
		wantReason    string
		wantReasonRaw string
	}{
		{name: "warning_null", notice: "warning", reason: nil, wantReasonRaw: "null"},
		{name: "warning_reason", notice: "warning", reason: "synthetic reason", wantReason: "synthetic reason", wantReasonRaw: `"synthetic reason"`},
		{name: "deactivation_null", notice: "deactivation", reason: nil, wantReasonRaw: "null"},
		{name: "deactivation_reason", notice: "deactivation", reason: "synthetic reason", wantReason: "synthetic reason", wantReasonRaw: `"synthetic reason"`},
		{name: "future_notice", notice: "future-notice", reason: nil, wantReasonRaw: "null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const caseID = "case/with ?#%"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				const wantPath = "/v1/safety/cases/case%2Fwith%20%3F%23%25"
				if r.Method != http.MethodGet || r.URL.EscapedPath() != wantPath {
					t.Errorf("Safety.Cases.Get(%q) route = %s %s, want GET %s", caseID, r.Method, r.URL.EscapedPath(), wantPath)
				}
				if got := r.URL.RawQuery; got != "trace=contract" {
					t.Errorf("Safety.Cases.Get(%q) query = %q, want trace=contract", caseID, got)
				}
				if got := r.Header.Get("Authorization"); got != "Bearer fake-project-key" {
					t.Errorf("Safety.Cases.Get(%q) authorization = %q, want fake project bearer", caseID, got)
				}
				if got := r.Header.Get("X-Case-Trace"); got != "caller-owned" {
					t.Errorf("Safety.Cases.Get(%q) caller header = %q, want caller-owned", caseID, got)
				}
				body, err := io.ReadAll(r.Body)
				if err != nil || len(body) != 0 {
					t.Errorf("Safety.Cases.Get(%q) body length = %d, error = %v, want empty body and nil error", caseID, len(body), err)
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]any{
					"id": caseID, "object": "safety.case", "created_at": 123,
					"entity_identifier": "synthetic-entity", "reason": tt.reason,
					"notice":            map[string]any{"type": tt.notice, "future_notice_field": true},
					"future_case_field": true,
				}); err != nil {
					t.Errorf("Safety.Cases.Get(%q) encode response error = %v, want nil", caseID, err)
				}
			}))
			t.Cleanup(server.Close)
			client := openai.NewClient(
				option.WithBaseURL(server.URL+"/v1"), option.WithHTTPClient(server.Client()),
				option.WithAPIKey("fake-project-key"), option.WithAdminAPIKey("fake-admin-key"),
			)
			result, err := client.Safety.Cases.Get(t.Context(), caseID,
				option.WithHeader("X-Case-Trace", "caller-owned"), option.WithQuery("trace", "contract"))
			if err != nil {
				t.Fatalf("Safety.Cases.Get(%q) error = %v, want nil", caseID, err)
			}
			if result.ID != caseID || result.Object != "safety.case" || result.CreatedAt != 123 || result.EntityIdentifier != "synthetic-entity" {
				t.Errorf("Safety.Cases.Get(%q) identity fields = %#v, want fixture identity", caseID, result)
			}
			if result.Notice.Type != tt.notice || result.Reason != tt.wantReason {
				t.Errorf("Safety.Cases.Get(%q) notice/reason = %q/%q, want %q/%q", caseID, result.Notice.Type, result.Reason, tt.notice, tt.wantReason)
			}
			if got := result.JSON.Reason.Raw(); got != tt.wantReasonRaw {
				t.Errorf("Safety.Cases.Get(%q) reason presence = %q, want %q", caseID, got, tt.wantReasonRaw)
			}
			if got := result.JSON.Reason.Valid(); got != (tt.reason != nil) {
				t.Errorf("Safety.Cases.Get(%q) reason valid = %t, want %t", caseID, got, tt.reason != nil)
			}
			if got := result.JSON.ExtraFields["future_case_field"].Raw(); got != "true" {
				t.Errorf("Safety.Cases.Get(%q) future case field = %q, want true", caseID, got)
			}
			if got := result.Notice.JSON.ExtraFields["future_notice_field"].Raw(); got != "true" {
				t.Errorf("Safety.Cases.Get(%q) future notice field = %q, want true", caseID, got)
			}
		})
	}
}

func TestSafetyCaseNeverUsesAdminKeyAsFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "" {
			t.Errorf("Safety.Cases.Get with no ordinary API key authorization = %q, want empty", got)
		}
		w.WriteHeader(http.StatusUnauthorized)
	}))
	t.Cleanup(server.Close)
	client := openai.NewClient(
		option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()),
		option.WithAPIKey(""), option.WithAdminAPIKey("fake-admin-key"), option.WithMaxRetries(0),
	)
	if _, err := client.Safety.Cases.Get(t.Context(), "synthetic-case"); err == nil {
		t.Fatal("Safety.Cases.Get with no ordinary API key error = nil, want authentication error")
	}
}
