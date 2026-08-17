package auth_test

import (
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3/auth"
)

func TestX509WorkloadIdentityAuthRejectsHTTPTransportChange(t *testing.T) {
	t.Run("explicit transport", func(t *testing.T) {
		var callsA atomic.Int32
		var callsB atomic.Int32
		transportA := &closureTransport{fn: func(*http.Request) (*http.Response, error) {
			callsA.Add(1)
			return tokenResponse("token-a", 60), nil
		}}
		transportB := &closureTransport{fn: func(*http.Request) (*http.Response, error) {
			callsB.Add(1)
			return tokenResponse("token-b", 60), nil
		}}
		httpClient := &http.Client{Transport: transportA}
		wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
		if err != nil {
			t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
		}

		token, err := wia.GetToken(t.Context(), httpClient)
		if err != nil {
			t.Fatalf("first GetToken() error = %v", err)
		}
		if token != "token-a" {
			t.Fatalf("first GetToken() = %q, want token-a", token)
		}
		token, err = wia.GetToken(t.Context(), httpClient)
		if err != nil {
			t.Fatalf("cached GetToken() error = %v", err)
		}
		if token != "token-a" {
			t.Fatalf("cached GetToken() = %q, want token-a", token)
		}

		httpClient.Transport = transportB
		ctx := t.Context()
		const callers = 16
		errs := make(chan error, callers)
		for range callers {
			go func() {
				_, getErr := wia.GetToken(ctx, httpClient)
				errs <- getErr
			}()
		}
		for range callers {
			if getErr := <-errs; getErr == nil {
				t.Error("GetToken() after HTTP transport change error = nil")
			}
		}
		if got, want := callsA.Load(), int32(1); got != want {
			t.Errorf("HTTP transport A exchange calls = %d, want %d", got, want)
		}
		if got := callsB.Load(); got != 0 {
			t.Errorf("HTTP transport B exchange calls = %d, want 0", got)
		}
	})

	t.Run("default transport", func(t *testing.T) {
		originalDefaultTransport := http.DefaultTransport
		t.Cleanup(func() {
			http.DefaultTransport = originalDefaultTransport
		})

		var callsA atomic.Int32
		var callsB atomic.Int32
		http.DefaultTransport = &closureTransport{fn: func(*http.Request) (*http.Response, error) {
			callsA.Add(1)
			return tokenResponse("token-a", 60), nil
		}}
		httpClient := &http.Client{}
		wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
		if err != nil {
			t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
		}
		if _, err := wia.GetToken(t.Context(), httpClient); err != nil {
			t.Fatalf("first GetToken() error = %v", err)
		}

		http.DefaultTransport = &closureTransport{fn: func(*http.Request) (*http.Response, error) {
			callsB.Add(1)
			return tokenResponse("token-b", 60), nil
		}}
		if _, err := wia.GetToken(t.Context(), httpClient); err == nil {
			t.Fatal("GetToken() after default HTTP transport change error = nil")
		}
		if got, want := callsA.Load(), int32(1); got != want {
			t.Errorf("default HTTP transport A exchange calls = %d, want %d", got, want)
		}
		if got := callsB.Load(); got != 0 {
			t.Errorf("default HTTP transport B exchange calls = %d, want 0", got)
		}
	})
}
