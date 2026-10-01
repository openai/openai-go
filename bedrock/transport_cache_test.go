package bedrock

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"sync"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/option"
)

func TestLoopbackTransportCacheReusesAndEvictsIdleConnections(t *testing.T) {
	closed := make(chan string, 32)
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}))
	server.Config.ConnState = func(conn net.Conn, state http.ConnState) {
		if state == http.StateClosed {
			closed <- conn.RemoteAddr().String()
		}
	}
	server.Start()
	defer server.Close()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	client, err := NewClient(ctx, Config{APIKey: "test-bearer", BaseURL: server.URL, UnsafeAllowHTTP: true}, option.WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}
	request := func(httpClient *http.Client) httptrace.GotConnInfo {
		t.Helper()
		var connection httptrace.GotConnInfo
		requestCtx := httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{GotConn: func(info httptrace.GotConnInfo) { connection = info }})
		var response map[string]any
		if err := client.Get(requestCtx, "/models", nil, &response, option.WithHTTPClient(httpClient)); err != nil {
			t.Fatal(err)
		}
		return connection
	}

	// Keep the caller's own pool warm so eviction also proves it is not closed.
	source := &http.Transport{IdleConnTimeout: time.Hour}
	defer source.CloseIdleConnections()
	originalClient := &http.Client{Transport: source}
	originalRequest := func() httptrace.GotConnInfo {
		t.Helper()
		var connection httptrace.GotConnInfo
		requestCtx := httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{GotConn: func(info httptrace.GotConnInfo) { connection = info }})
		req, err := http.NewRequestWithContext(requestCtx, http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		response, err := originalClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if closeErr := response.Body.Close(); closeErr != nil {
				t.Error(closeErr)
			}
		}()
		if _, err := io.Copy(io.Discard, response.Body); err != nil {
			t.Fatal(err)
		}
		return connection
	}
	original := originalRequest()
	first := request(originalClient)
	if second := request(originalClient); !second.Reused || second.Conn != first.Conn {
		t.Fatal("repeated calls with the same source transport did not reuse a connection")
	}

	addresses := []string{first.Conn.LocalAddr().String()}
	for range 16 {
		connection := request(&http.Client{Transport: &http.Transport{IdleConnTimeout: time.Hour}})
		if connection.Reused {
			t.Fatal("a different source transport unexpectedly reused another transport's pool")
		}
		addresses = append(addresses, connection.Conn.LocalAddr().String())
		if len(addresses) <= 8 {
			continue
		}
		// Each ninth distinct transport must close the oldest idle clone, leaving
		// at most eight retained idle pools without waiting for their timeouts.
		wantClosed := addresses[len(addresses)-9]
		select {
		case address := <-closed:
			if address != wantClosed {
				t.Fatalf("closed connection = %s, want oldest clone %s", address, wantClosed)
			}
		case <-ctx.Done():
			t.Fatalf("oldest idle clone %s was not closed after cache eviction", wantClosed)
		}
	}
	if after := originalRequest(); !after.Reused || after.Conn != original.Conn {
		t.Fatal("eviction closed the caller's original transport pool")
	}
}

func TestLoopbackTransportCacheEvictionPreservesActiveResponse(t *testing.T) {
	release := make(chan struct{}, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/active" {
			_, _ = io.WriteString(w, "before ")
			w.(http.Flusher).Flush()
			select {
			case <-release:
				_, _ = io.WriteString(w, "after")
			case <-req.Context().Done():
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}))
	defer server.Close()
	defer close(release)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	client, err := NewClient(ctx, Config{APIKey: "test-bearer", BaseURL: server.URL, UnsafeAllowHTTP: true}, option.WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}
	var active *http.Response
	if requestErr := client.Get(ctx, "/active", nil, &active, option.WithHTTPClient(&http.Client{Transport: &http.Transport{}})); requestErr != nil {
		t.Fatal(requestErr)
	}
	defer func() {
		if closeErr := active.Body.Close(); closeErr != nil {
			t.Error(closeErr)
		}
	}()
	for range 8 {
		var response map[string]any
		if requestErr := client.Get(ctx, "/models", nil, &response, option.WithHTTPClient(&http.Client{Transport: &http.Transport{}})); requestErr != nil {
			t.Fatal(requestErr)
		}
	}
	release <- struct{}{}
	body, err := io.ReadAll(active.Body)
	if err != nil || string(body) != "before after" {
		t.Fatalf("active response after eviction = %q, %v", body, err)
	}
}

func TestLoopbackTransportCacheConcurrentChurn(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	client, err := NewClient(ctx, Config{APIKey: "test-bearer", BaseURL: server.URL, UnsafeAllowHTTP: true}, option.WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}
	var workers sync.WaitGroup
	for range 16 {
		workers.Go(func() {
			httpClient := &http.Client{Transport: &http.Transport{}}
			for range 2 {
				var response map[string]any
				if err := client.Get(ctx, "/models", nil, &response, option.WithHTTPClient(httpClient)); err != nil {
					t.Error(err)
					return
				}
			}
		})
	}
	workers.Wait()
}
