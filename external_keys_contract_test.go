package openai_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestExternalKeysProviders(t *testing.T) {
	cases := []struct {
		name   string
		params openai.AdminOrganizationExternalKeyNewParams
		want   map[string]string
	}{
		{"aws", openai.AdminOrganizationExternalKeyNewParams{OfAws: &openai.AdminOrganizationExternalKeyNewParamsBodyAws{Name: "AWS", ExternalID: "org-test", RoleArn: "arn:aws:iam::123456789012:role/test", KmsArn: "arn:aws:kms:us-east-1:123456789012:key/test"}}, map[string]string{"type": "aws", "name": "AWS", "external_id": "org-test", "role_arn": "arn:aws:iam::123456789012:role/test", "kms_arn": "arn:aws:kms:us-east-1:123456789012:key/test"}},
		{"gcp", openai.AdminOrganizationExternalKeyNewParams{OfGcp: &openai.AdminOrganizationExternalKeyNewParamsBodyGcp{Name: "GCP", Audience: "org-test", KmsKeyLocation: "global", KmsKeyName: "test", KmsKeyRingName: "ring", KmsProjectID: "project", WorkloadIdentityPoolID: "pool", WorkloadIdentityProjectNumber: "123", WorkloadIdentityProviderID: "provider", Purpose: openai.String("psp_decryption"), KmsKeyVersion: openai.String("1")}}, map[string]string{"type": "gcp", "name": "GCP", "audience": "org-test", "kms_key_location": "global", "kms_key_name": "test", "kms_key_ring_name": "ring", "kms_project_id": "project", "workload_identity_pool_id": "pool", "workload_identity_project_number": "123", "workload_identity_provider_id": "provider", "purpose": "psp_decryption", "kms_key_version": "1"}},
		{"azure", openai.AdminOrganizationExternalKeyNewParams{OfAzure: &openai.AdminOrganizationExternalKeyNewParamsBodyAzure{Name: "Azure", KeyName: "org-test--key", TenantID: "00000000-0000-0000-0000-000000000000", VaultURI: "https://example.vault.azure.net", Purpose: openai.String("psp_decryption"), KeyVersion: openai.String("00000000000000000000000000000000")}}, map[string]string{"type": "azure", "name": "Azure", "key_name": "org-test--key", "tenant_id": "00000000-0000-0000-0000-000000000000", "vault_uri": "https://example.vault.azure.net", "purpose": "psp_decryption", "key_version": "00000000000000000000000000000000"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" || r.URL.Path != "/v1/organization/external_keys" {
					t.Errorf("unexpected route: %s %s", r.Method, r.URL.Path)
				}
				if got := r.Header.Get("Authorization"); got != "Bearer admin-test" {
					t.Errorf("authorization = %q", got)
				}
				var got map[string]string
				if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("body = %#v, want %#v", got, tc.want)
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]any{"id": "key-test", "object": "organization.external_key", "type": tc.name, "name": "test", "api_project_ids": []string{}, "workspace_id": "", "created_at": 1}); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL+"/v1/"), option.WithAPIKey("project-test"), option.WithAdminAPIKey("admin-test"), option.WithMaxRetries(0))
			key, err := client.Admin.Organization.ExternalKeys.New(t.Context(), tc.params)
			if err != nil {
				t.Fatal(err)
			}
			if key.ID != "key-test" {
				t.Errorf("ID = %q", key.ID)
			}
		})
	}
}

func TestExternalKeysLifecycle(t *testing.T) {
	var calls atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer admin-test" {
			t.Errorf("authorization = %q", got)
		}
		if calls.Add(1) > 5 {
			t.Error("unexpected extra request")
			http.Error(w, "unexpected extra request", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		var body any = map[string]any{"id": "key-test", "object": "organization.external_key", "type": "aws", "name": "test", "api_project_ids": []string{}, "workspace_id": "", "created_at": 1}
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/organization/external_keys":
			if r.URL.Query().Get("purpose") != "ekm" || r.URL.Query().Get("limit") != "1" {
				t.Errorf("query = %s", r.URL.RawQuery)
			}
			id := "key-one"
			more := true
			if r.URL.Query().Get("after") == "key-one" {
				id = "key-two"
				more = false
			}
			body = map[string]any{"object": "list", "data": []any{map[string]any{"id": id}}, "first_id": id, "last_id": id, "has_more": more}
		case r.Method == http.MethodPost && r.URL.Path == "/v1/organization/external_keys/key-test/validate":
			body = map[string]string{"status": "valid"}
		case (r.Method == http.MethodGet || r.Method == http.MethodDelete) && r.URL.Path == "/v1/organization/external_keys/key-test":
		default:
			t.Errorf("unexpected route: %s %s", r.Method, r.URL.Path)
			http.Error(w, "unexpected route", http.StatusBadRequest)
			return
		}
		if err := json.NewEncoder(w).Encode(body); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL+"/v1/"), option.WithAPIKey("project-test"), option.WithAdminAPIKey("admin-test"), option.WithMaxRetries(0))
	keys := client.Admin.Organization.ExternalKeys
	ctx := t.Context()
	retrieved, err := keys.Get(ctx, "key-test")
	if err != nil {
		t.Fatal(err)
	}
	if retrieved.ID != "key-test" {
		t.Errorf("retrieved ID = %q", retrieved.ID)
	}
	deleted, err := keys.Delete(ctx, "key-test")
	if err != nil {
		t.Fatal(err)
	}
	if deleted.ID != "key-test" || deleted.Object != "organization.external_key" || deleted.Name != "test" {
		t.Errorf("delete response = %#v", deleted)
	}
	validated, err := keys.Validate(ctx, "key-test")
	if err != nil {
		t.Fatal(err)
	}
	if validated.Status != "valid" {
		t.Errorf("status = %q", validated.Status)
	}
	pager := keys.ListAutoPaging(ctx, openai.AdminOrganizationExternalKeyListParams{Limit: openai.Int(1), Purpose: openai.AdminOrganizationExternalKeyListParamsPurposeEkm})
	var ids []string
	for pager.Next() {
		ids = append(ids, pager.Current().ID)
	}
	if err := pager.Err(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ids, []string{"key-one", "key-two"}) {
		t.Errorf("IDs = %v", ids)
	}
	before := calls.Load()
	if _, err := keys.Get(ctx, ""); err == nil {
		t.Error("Get accepted empty ID")
	}
	if _, err := keys.Delete(ctx, ""); err == nil {
		t.Error("Delete accepted empty ID")
	}
	if _, err := keys.Validate(ctx, ""); err == nil {
		t.Error("Validate accepted empty ID")
	}
	if calls.Load() != before {
		t.Error("empty ID sent a request")
	}
}
