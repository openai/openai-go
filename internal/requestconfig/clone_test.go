package requestconfig

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestCloneWithErrorRejectsNonReplayableRequestBodies(t *testing.T) {
	for _, test := range []struct {
		name      string
		configure func(*RequestConfig)
		want      string
	}{
		{
			name: "stream without body factory",
			configure: func(cfg *RequestConfig) {
				cfg.Request.Body = io.NopCloser(strings.NewReader("synthetic-private-stream"))
				cfg.Request.GetBody = nil
			},
			want: "non-replayable request body",
		},
		{
			name: "body factory failure is sanitized",
			configure: func(cfg *RequestConfig) {
				cfg.Request.Body = io.NopCloser(strings.NewReader("synthetic-private-stream"))
				cfg.Request.GetBody = func() (io.ReadCloser, error) {
					return nil, errors.New("synthetic-private-factory-detail")
				}
			},
			want: "could not recreate request body",
		},
		{
			name: "body factory returning nil",
			configure: func(cfg *RequestConfig) {
				cfg.Request.Body = io.NopCloser(strings.NewReader("synthetic-private-stream"))
				cfg.Request.GetBody = func() (io.ReadCloser, error) { return nil, nil }
			},
			want: "could not recreate request body",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := newBodyCloseRequestConfig(t, http.NoBody)
			test.configure(cfg)
			clone, err := cfg.CloneWithError(t.Context())
			if clone != nil || err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("clone = %v, error = %v, want nil and %q", clone, err, test.want)
			}
			if strings.Contains(err.Error(), "synthetic-private") {
				t.Errorf("clone error disclosed sensitive request or body-factory content: %q", err.Error())
			}
			if legacy := cfg.Clone(t.Context()); legacy != nil {
				t.Errorf("legacy Clone returned %v for an unrecoverable request body", legacy)
			}
		})
	}
}

func TestCloneWithErrorRejectsInvalidInputs(t *testing.T) {
	var missing *RequestConfig
	if clone, err := missing.CloneWithError(t.Context()); clone != nil || err == nil {
		t.Errorf("nil configuration clone = %v, error = %v", clone, err)
	}
	if clone := missing.Clone(t.Context()); clone != nil {
		t.Errorf("legacy nil configuration clone = %v, want nil", clone)
	}
	if clone, err := (&RequestConfig{}).CloneWithError(t.Context()); clone != nil || err == nil {
		t.Errorf("nil request clone = %v, error = %v", clone, err)
	}
	cfg := newBodyCloseRequestConfig(t, http.NoBody)
	var absent context.Context
	if clone, err := cfg.CloneWithError(absent); clone != nil || err == nil {
		t.Errorf("nil context clone = %v, error = %v", clone, err)
	}
}

func TestCloneWithErrorClosesBodyReturnedWithFactoryError(t *testing.T) {
	cfg := newBodyCloseRequestConfig(t, http.NoBody)
	cfg.Request.Body = io.NopCloser(strings.NewReader("synthetic-private-stream"))
	recreated := &closeTrackingReadCloser{
		ReadCloser: io.NopCloser(strings.NewReader("synthetic-private-recreated-stream")),
	}
	cfg.Request.GetBody = func() (io.ReadCloser, error) {
		return recreated, errors.New("synthetic-private-factory-detail")
	}

	clone, err := cfg.CloneWithError(t.Context())
	if clone != nil || err == nil {
		t.Fatalf("clone = %v, error = %v, want nil and a sanitized body-factory error", clone, err)
	}
	if strings.Contains(err.Error(), "synthetic-private") {
		t.Errorf("clone error disclosed sensitive body-factory content: %q", err.Error())
	}
	if recreated.closes != 1 {
		t.Errorf("failed body-factory result was closed %d times, want once", recreated.closes)
	}
}

func TestCloneWithErrorPreservesReplayableBodyAndContext(t *testing.T) {
	cfg := newBodyCloseRequestConfig(t, http.NoBody)
	cfg.Request.Body = io.NopCloser(strings.NewReader("synthetic-replayable-body"))
	cfg.Request.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader("synthetic-replayable-body")), nil
	}
	clone, err := cfg.CloneWithError(t.Context())
	if err != nil || clone == nil {
		t.Fatalf("clone replayable request = %v, error = %v", clone, err)
	}
	body, readErr := io.ReadAll(clone.Request.Body)
	if readErr != nil || string(body) != "synthetic-replayable-body" {
		t.Errorf("cloned request body = %q, error = %v", body, readErr)
	}
	if closeErr := clone.Request.Body.Close(); closeErr != nil {
		t.Errorf("close cloned replayable body: %v", closeErr)
	}
	if clone.Request.Context() != t.Context() || clone.Context != t.Context() {
		t.Error("cloned replayable request did not preserve its selected context")
	}
}
