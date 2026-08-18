package requestconfig

import (
	"errors"
	"net/url"
	"reflect"
	"testing"
)

func TestInheritedOptionsPreservePreRequestOptions(t *testing.T) {
	marker := errors.New("pre-request failure")
	seen := []string{}
	pre := func(value string) RequestOption {
		return PreRequestOptionFunc(func(cfg *RequestConfig) error { seen = append(seen, value); cfg.WebhookSecret = value; return nil })
	}
	ordinary := RequestOptionFunc(func(*RequestConfig) error { t.Fatal("ordinary option executed during pre-request scan"); return nil })
	inner := InheritedOptions(pre("inner"), ordinary, WithEnvironmentDefaultsDisabled())
	opts := InheritedOptions(append(inner, pre("outer"))...)
	cfg, err := PreRequestOptions(opts...)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.WebhookSecret != "outer" || !reflect.DeepEqual(seen, []string{"inner", "outer"}) {
		t.Fatalf("pre-request order lost: %#v / %q", seen, cfg.WebhookSecret)
	}
	if !EnvironmentDefaultsDisabled(opts...) {
		t.Fatal("nested provider marker lost")
	}
	_, err = PreRequestOptions(InheritedOptions(PreRequestOptionFunc(func(*RequestConfig) error { return marker }))...)
	if !errors.Is(err, marker) {
		t.Fatalf("pre-request error lost: %v", err)
	}
}

func TestInheritedOptionsCaptureSlice(t *testing.T) {
	first := RequestOptionFunc(func(cfg *RequestConfig) error { cfg.WebhookSecret = "first"; return nil })
	second := RequestOptionFunc(func(cfg *RequestConfig) error { cfg.WebhookSecret = "second"; return nil })
	original := []RequestOption{first}
	captured := InheritedOptions(original...)
	original[0] = second
	cfg := RequestConfig{}
	if err := cfg.Apply(captured...); err != nil {
		t.Fatal(err)
	}
	if cfg.WebhookSecret != "first" {
		t.Fatal("inherited slice was mutated")
	}
}

func TestInheritedOptionsRestoreOuterLayer(t *testing.T) {
	endpoint := &url.URL{Scheme: "https", Host: "example.com"}
	residency := RequestOptionFunc(func(cfg *RequestConfig) error { return cfg.SetEndpoint("data_residency", endpoint) })
	base := RequestOptionFunc(func(cfg *RequestConfig) error { return cfg.SetEndpoint("base_url", endpoint) })
	cfg := RequestConfig{}
	opts := append([]RequestOption{base}, InheritedOptions(residency)...)
	opts = append(opts, residency)
	if err := cfg.Apply(opts...); err == nil {
		t.Fatal("nested layer concealed an outer conflict")
	}
}
