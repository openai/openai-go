package openai_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestExternalStorageOCIContract(t *testing.T) {
	providerJSON := map[string]any{"type": "oci", "bucket": "test-bucket", "region": "us-ashburn-1", "tenancy_ocid": "ocid1.tenancy.oc1..synthetic"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/organization/external_storage":
			if r.Method != http.MethodPost {
				t.Errorf("method = %s, want POST", r.Method)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode request: %v", err)
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			want := map[string]any{"project_id": "proj_test", "provider": providerJSON}
			if !reflect.DeepEqual(body, want) {
				t.Errorf("body = %#v, want %#v", body, want)
			}
			if err := json.NewEncoder(w).Encode(externalStorageFixture("extstorage_test", providerJSON)); err != nil {
				t.Errorf("write response: %v", err)
			}
		case "/organization/audit_logs":
			if r.Method != http.MethodGet {
				t.Errorf("method = %s, want GET", r.Method)
			}
			if err := json.NewEncoder(w).Encode(map[string]any{
				"object": "list", "has_more": false,
				"data": []any{map[string]any{
					"id": "audit_test", "type": "external_storage.registered", "effective_at": 123,
					"external_storage.registered": map[string]any{
						"id": "extstorage_test", "data": map[string]any{"geography": "US", "provider": providerJSON},
					},
				}},
			}); err != nil {
				t.Errorf("write audit response: %v", err)
			}
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	client := openai.NewClient(option.WithAPIKey("test-key"), option.WithAdminAPIKey("test-admin-key"), option.WithBaseURL(server.URL), option.WithUnsafeAllowHTTP(), option.WithHTTPClient(server.Client()), option.WithMaxRetries(0))
	provider := openai.AdminOrganizationExternalStorageNewParamsProviderUnion{
		OfOci: &openai.AdminOrganizationExternalStorageNewParamsProviderOci{Bucket: "test-bucket", Region: "us-ashburn-1", TenancyOcid: "ocid1.tenancy.oc1..synthetic"},
	}
	storage, err := client.Admin.Organization.ExternalStorage.New(t.Context(), openai.AdminOrganizationExternalStorageNewParams{ProjectID: "proj_test", Provider: provider})
	if err != nil {
		t.Fatal(err)
	}
	oci, ok := storage.Provider.AsAny().(openai.OciExternalStorageProvider)
	if !ok {
		t.Fatalf("provider variant = %T, want OciExternalStorageProvider", storage.Provider.AsAny())
	}
	assertOCIProvider(t, oci)
	assertOCIProvider(t, storage.Provider.AsOci())

	logs, err := client.Admin.Organization.AuditLogs.List(t.Context(), openai.AdminOrganizationAuditLogListParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(logs.Data) != 1 || logs.Data[0].Type != "external_storage.registered" {
		t.Fatalf("audit data = %#v, want one external_storage.registered event", logs.Data)
	}
	registered := logs.Data[0].ExternalStorageRegistered
	if registered.ID != "extstorage_test" {
		t.Errorf("registered ID = %q, want extstorage_test", registered.ID)
	}
	oci, ok = registered.Data.Provider.AsAny().(openai.OciExternalStorageProvider)
	if !ok {
		t.Fatalf("audit provider variant = %T, want OciExternalStorageProvider", registered.Data.Provider.AsAny())
	}
	assertOCIProvider(t, oci)
	assertOCIProvider(t, registered.Data.Provider.AsOci())

	// Parameter union decoding has a separate discriminator registration.
	var decoded openai.AdminOrganizationExternalStorageNewParamsProviderUnion
	if err := json.Unmarshal([]byte(`{"type":"oci","bucket":"test-bucket","region":"us-ashburn-1","tenancy_ocid":"ocid1.tenancy.oc1..synthetic"}`), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.OfOci == nil {
		t.Fatal("OCI parameter discriminator did not select OfOci")
	}
	if got := decoded.GetBucket(); got == nil || *got != "test-bucket" {
		t.Errorf("GetBucket() = %v, want test-bucket", got)
	}
	if got := decoded.GetRegion(); got == nil || *got != "us-ashburn-1" {
		t.Errorf("GetRegion() = %v, want us-ashburn-1", got)
	}
	if got := decoded.GetTenancyOcid(); got == nil || *got != "ocid1.tenancy.oc1..synthetic" {
		t.Errorf("GetTenancyOcid() = %v, want synthetic tenancy OCID", got)
	}
}

func assertOCIProvider(t *testing.T, provider openai.OciExternalStorageProvider) {
	t.Helper()
	if provider.Type != "oci" || provider.Bucket != "test-bucket" || provider.Region != "us-ashburn-1" || provider.TenancyOcid != "ocid1.tenancy.oc1..synthetic" {
		t.Errorf("OCI provider = %#v, want synthetic OCI fields", provider)
	}
}
