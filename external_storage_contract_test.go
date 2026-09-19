package openai_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestExternalStorageCreateProviders(t *testing.T) {
	tests := []struct {
		name     string
		provider openai.AdminOrganizationExternalStorageNewParamsProviderUnion
		want     map[string]any
		metadata map[string]any
	}{
		{
			name: "aws",
			provider: openai.AdminOrganizationExternalStorageNewParamsProviderUnion{
				OfAws: &openai.AdminOrganizationExternalStorageNewParamsProviderAws{
					Bucket: "test-bucket", RoleArn: "arn:aws:iam::000000000000:role/test",
				},
			},
			want: map[string]any{
				"type": "aws", "bucket": "test-bucket", "role_arn": "arn:aws:iam::000000000000:role/test",
			},
			metadata: map[string]any{
				"type": "aws", "bucket": "test-bucket", "role_arn": "arn:aws:iam::000000000000:role/test",
				"account_id": "000000000000", "region": "us-east-1", "external_id": "test-external-id",
			},
		},
		{
			name: "azure",
			provider: openai.AdminOrganizationExternalStorageNewParamsProviderUnion{
				OfAzure: &openai.AdminOrganizationExternalStorageNewParamsProviderAzure{
					AccountName: "test-account", Container: "test-container", ResourceGroup: "test-group",
					SubscriptionID: "test-subscription", TenantID: "test-tenant",
				},
			},
			want: map[string]any{
				"type": "azure", "account_name": "test-account", "container": "test-container",
				"resource_group": "test-group", "subscription_id": "test-subscription", "tenant_id": "test-tenant",
			},
			metadata: map[string]any{
				"type": "azure", "account_name": "test-account", "container": "test-container",
				"resource_group": "test-group", "subscription_id": "test-subscription", "tenant_id": "test-tenant",
				"region": "eastus",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/v1/organization/external_storage" {
					t.Errorf("ExternalStorage.New(%s) route = %s %s, want POST /v1/organization/external_storage", tt.name, r.Method, r.URL.Path)
				}
				if got := r.Header.Get("Authorization"); got != "Bearer fake-admin-key" {
					t.Errorf("ExternalStorage.New(%s) authorization = %q, want fake admin bearer", tt.name, got)
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("ExternalStorage.New(%s) decode body: %v", tt.name, err)
					return
				}
				want := map[string]any{"project_id": "proj_test", "provider": tt.want}
				if !reflect.DeepEqual(body, want) {
					t.Errorf("ExternalStorage.New(%s) body = %#v, want %#v", tt.name, body, want)
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(externalStorageFixture("extstorage_test", tt.metadata)); err != nil {
					t.Errorf("ExternalStorage.New(%s) encode response: %v", tt.name, err)
				}
			}))
			t.Cleanup(server.Close)
			client := openai.NewClient(
				option.WithBaseURL(server.URL+"/v1"), option.WithHTTPClient(server.Client()),
				option.WithAPIKey("fake-project-key"), option.WithAdminAPIKey("fake-admin-key"),
			)
			result, err := client.Admin.Organization.ExternalStorage.New(t.Context(), openai.AdminOrganizationExternalStorageNewParams{
				ProjectID: "proj_test", Provider: tt.provider,
			})
			if err != nil {
				t.Fatalf("ExternalStorage.New(%s) error = %v, want nil", tt.name, err)
			}
			if result.Provider.Type != tt.name || result.Provider.AsAny() == nil {
				t.Errorf("ExternalStorage.New(%s) provider = %#v, want recognized %s variant", tt.name, result.Provider, tt.name)
			}
			if got := result.JSON.ExtraFields["future_field"].Raw(); got != "true" {
				t.Errorf("ExternalStorage.New(%s) future field = %q, want true", tt.name, got)
			}
			var provider map[string]any
			if err := json.Unmarshal([]byte(result.Provider.RawJSON()), &provider); err != nil {
				t.Fatalf("ExternalStorage.New(%s) decode provider: %v", tt.name, err)
			}
			if !reflect.DeepEqual(provider, tt.metadata) {
				t.Errorf("ExternalStorage.New(%s) provider = %#v, want %#v", tt.name, provider, tt.metadata)
			}
		})
	}
}

func TestExternalStoragePaginationPreservesFiltersAndAdminAuth(t *testing.T) {
	requests := make(chan url.Values, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		select {
		case requests <- query:
		default:
			t.Error("ExternalStorage.ListAutoPaging made more than two requests")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != "/v1/organization/external_storage" {
			t.Errorf("ExternalStorage.ListAutoPaging route = %s %s, want GET collection", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer fake-admin-key" {
			t.Errorf("ExternalStorage.ListAutoPaging authorization = %q, want fake admin bearer", got)
		}
		id, more := "extstorage_first", true
		if query.Get("after") != "" {
			id, more = "extstorage_last", false
		}
		w.Header().Set("Content-Type", "application/json")
		provider := map[string]any{
			"type": "aws", "bucket": "test-bucket", "role_arn": "arn:aws:iam::000000000000:role/test",
			"account_id": "000000000000", "region": "us-east-1", "external_id": "test-external-id",
		}
		if err := json.NewEncoder(w).Encode(map[string]any{
			"object": "list", "data": []any{externalStorageFixture(id, provider)},
			"first_id": id, "last_id": id, "has_more": more,
		}); err != nil {
			t.Errorf("ExternalStorage.ListAutoPaging encode response: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	client := openai.NewClient(
		option.WithBaseURL(server.URL+"/v1"), option.WithHTTPClient(server.Client()),
		option.WithAPIKey("fake-project-key"), option.WithAdminAPIKey("fake-admin-key"),
	)
	page := client.Admin.Organization.ExternalStorage.ListAutoPaging(t.Context(), openai.AdminOrganizationExternalStorageListParams{
		ProjectID: openai.String("proj_test"), Limit: openai.Int(1), Order: openai.AdminOrganizationExternalStorageListParamsOrderDesc,
	})
	var ids []string
	for page.Next() {
		ids = append(ids, page.Current().ID)
	}
	if err := page.Err(); err != nil {
		t.Fatalf("ExternalStorage.ListAutoPaging error = %v, want nil", err)
	}
	if want := []string{"extstorage_first", "extstorage_last"}; !reflect.DeepEqual(ids, want) {
		t.Errorf("ExternalStorage.ListAutoPaging IDs = %v, want %v", ids, want)
	}
	for _, after := range []string{"", "extstorage_first"} {
		select {
		case got := <-requests:
			want := url.Values{"project_id": {"proj_test"}, "limit": {"1"}, "order": {"desc"}}
			if after != "" {
				want.Set("after", after)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("ExternalStorage.ListAutoPaging query = %v, want %v", got, want)
			}
		default:
			t.Errorf("ExternalStorage.ListAutoPaging missing request with after=%q", after)
		}
	}
	select {
	case got := <-requests:
		t.Errorf("ExternalStorage.ListAutoPaging extra request = %v, want none", got)
	default:
	}
}

func TestExternalStorageIDRoutesAndFutureProvider(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "/v1/organization/external_storage/ext%2Fstorage%3Ftest"
		if r.Method == http.MethodPost {
			want += "/validate"
		}
		if got := r.URL.EscapedPath(); got != want || r.URL.RawQuery != "" {
			t.Errorf("ExternalStorage ID route = %s?%s, want %s without query", got, r.URL.RawQuery, want)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer fake-admin-key" {
			t.Errorf("ExternalStorage ID route authorization = %q, want fake admin bearer", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("ExternalStorage ID route read body: %v", err)
		} else if len(body) != 0 {
			t.Errorf("ExternalStorage ID route body = %q, want empty", body)
		}
		response := externalStorageFixture("ext/storage?test", map[string]any{"type": "future", "future_provider_field": true})
		if r.Method == http.MethodDelete {
			response = map[string]any{"id": "ext/storage?test", "object": "organization.external_storage.deleted", "deleted": true}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("ExternalStorage ID route encode response: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	client := openai.NewClient(
		option.WithBaseURL(server.URL+"/v1"), option.WithHTTPClient(server.Client()),
		option.WithAPIKey("fake-project-key"), option.WithAdminAPIKey("fake-admin-key"),
	)
	storage := &client.Admin.Organization.ExternalStorage
	retrieved, err := storage.Get(t.Context(), "ext/storage?test")
	if err != nil {
		t.Fatalf("ExternalStorage.Get error = %v, want nil", err)
	}
	if retrieved.Provider.Type != "future" || retrieved.Provider.AsAny() != nil {
		t.Errorf("ExternalStorage.Get future provider = %#v, want unrecognized future variant", retrieved.Provider)
	}
	validated, err := storage.Validate(t.Context(), "ext/storage?test")
	if err != nil {
		t.Fatalf("ExternalStorage.Validate error = %v, want nil", err)
	}
	if validated.ID != retrieved.ID || validated.Provider.RawJSON() != retrieved.Provider.RawJSON() {
		t.Errorf("ExternalStorage.Validate result = %#v, want matching preserved future provider", validated)
	}
	deleted, err := storage.Delete(t.Context(), "ext/storage?test")
	if err != nil {
		t.Fatalf("ExternalStorage.Delete error = %v, want nil", err)
	}
	if deleted.ID != retrieved.ID || !deleted.Deleted {
		t.Errorf("ExternalStorage.Delete result = %#v, want matching deleted configuration", deleted)
	}
}

func TestExternalStorageNeverFallsBackToProjectKey(t *testing.T) {
	t.Setenv("OPENAI_ADMIN_KEY", "")
	requests := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests <- r.Header.Get("Authorization")
		w.WriteHeader(http.StatusUnauthorized)
	}))
	t.Cleanup(server.Close)
	client := openai.NewClient(
		option.WithBaseURL(server.URL+"/v1"), option.WithHTTPClient(server.Client()),
		option.WithAPIKey("fake-project-key"), option.WithAdminAPIKey(""), option.WithMaxRetries(0),
	)
	_, err := client.Admin.Organization.ExternalStorage.Get(t.Context(), "extstorage_test")
	var apiError *openai.Error
	if !errors.As(err, &apiError) || apiError.StatusCode != http.StatusUnauthorized {
		t.Fatalf("ExternalStorage.Get without admin key error = %v, want unauthorized API error", err)
	}
	select {
	case got := <-requests:
		if got != "" {
			t.Errorf("ExternalStorage.Get without admin key authorization = %q, want no project-key fallback", got)
		}
	default:
		t.Error("ExternalStorage.Get without admin key made no request to the mock server")
	}
}

func externalStorageFixture(id string, provider map[string]any) map[string]any {
	return map[string]any{
		"id": id, "object": "organization.external_storage", "project_id": "proj_test",
		"provider": provider, "geography": "US", "status": "validated", "created_at": 123,
		"future_field": true,
	}
}
