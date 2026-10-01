package requestconfig

import (
	"net/http"
	"reflect"
	"testing"
)

func TestCloneIsolatesProviderHeaderOwnership(t *testing.T) {
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	original := RequestConfig{Request: req}
	original.SetHeader("X-Custom-Credential", "fake-inherited")
	original.SetHeader("X-Other-Credential", "fake-other")
	clone, err := original.CloneWithError(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	clone.AddHeader("x-custom-credential", "fake-explicit")
	clone.ClearInheritedProviderHeaders("Azure")
	if got := clone.Request.Header.Values("X-Custom-Credential"); !reflect.DeepEqual(got, []string{"fake-explicit"}) {
		t.Fatalf("clone values = %v", got)
	}
	if clone.Request.Header.Get("X-Other-Credential") != "" {
		t.Fatal("clone retained inherited header")
	}
	if original.Request.Header.Get("X-Custom-Credential") != "fake-inherited" || original.Request.Header.Get("X-Other-Credential") != "fake-other" {
		t.Fatal("clone changed original headers")
	}
	// The clone must not change ownership records used by another descendant.
	sibling, err := original.CloneWithError(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	sibling.ClearInheritedProviderHeaders("Azure")
	if len(sibling.Request.Header) != 0 {
		t.Fatal("sibling retained inherited headers")
	}
}
