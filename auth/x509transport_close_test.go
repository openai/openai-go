package auth_test

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestX509TransportCloseDoesNotCancelRequestsAlreadyInProgress(t *testing.T) {
	fixture := newX509TransportFixture(t)
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	dial := template.DialContext
	dialStarted := make(chan struct{})
	continueDial := make(chan struct{})
	template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		close(dialStarted)
		select {
		case <-continueDial:
			return dial(ctx, network, address)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	capability := newX509Capability(t, template)
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	completed := make(chan error, 1)
	go func() {
		response, err := capability.Do(request)
		if err == nil {
			err = response.Body.Close()
		}
		completed <- err
	}()

	select {
	case <-dialStarted:
	case <-time.After(5 * time.Second):
		t.Fatal("request did not begin dialing before capability closure")
	}
	if err := capability.Close(); err != nil {
		t.Fatalf("close capability with request in progress: %v", err)
	}
	future := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	if err := x509Rejected(t, capability, future); !strings.Contains(err.Error(), "closed") {
		t.Fatalf("request begun after closure error = %v", err)
	}
	close(continueDial)

	select {
	case err := <-completed:
		if err != nil {
			t.Fatalf("request already in progress when Close returned: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("request already in progress did not complete after Close")
	}
}

func TestX509TransportCloseReleasesHTTP2ConnectionAfterActiveResponseFinishes(t *testing.T) {
	fixture := newX509TransportFixture(t)
	started := make(chan struct{})
	release := make(chan struct{})
	server, closed := newX509HTTP2LifecycleServer(t, fixture, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		close(started)
		<-release
		_, _ = io.WriteString(w, "synthetic HTTP/2 response")
	}))

	template := fixture.transport(t, server)
	template.ForceAttemptHTTP2 = true
	capability := newX509Capability(t, template)
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	response, err := capability.Do(request)
	if err != nil {
		t.Fatalf("begin active HTTP/2 response: %v", err)
	}
	if response.ProtoMajor != 2 {
		t.Fatalf("negotiated HTTP/%d, want HTTP/2", response.ProtoMajor)
	}
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("HTTP/2 response did not start")
	}
	if err := capability.Close(); err != nil {
		t.Fatalf("close capability with active HTTP/2 response: %v", err)
	}
	close(release)
	if _, err := io.Copy(io.Discard, response.Body); err != nil {
		t.Fatalf("read active HTTP/2 response after Close: %v", err)
	}
	if err := response.Body.Close(); err != nil {
		t.Fatalf("close active HTTP/2 response body: %v", err)
	}
	select {
	case <-closed:
	case <-time.After(5 * time.Second):
		t.Fatal("HTTP/2 connection remained pooled after the active response finished")
	}
}

func TestX509TransportCloseReleasesHTTP2ConnectionAfterCancellation(t *testing.T) {
	for _, afterHeaders := range []bool{false, true} {
		name := "before response headers"
		if afterHeaders {
			name = "after response headers"
		}
		t.Run(name, func(t *testing.T) {
			fixture := newX509TransportFixture(t)
			started := make(chan struct{})
			server, closed := newX509HTTP2LifecycleServer(t, fixture,
				http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
					if afterHeaders {
						w.WriteHeader(http.StatusOK)
						if flusher, ok := w.(http.Flusher); ok {
							flusher.Flush()
						}
					}
					close(started)
					<-request.Context().Done()
				}))
			template := fixture.transport(t, server)
			template.ForceAttemptHTTP2 = true
			capability := newX509Capability(t, template)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models").WithContext(ctx)
			type result struct {
				response *http.Response
				err      error
			}
			completed := make(chan result, 1)
			go func() {
				response, err := capability.Do(request)
				if err != nil && response != nil {
					_ = response.Body.Close()
					response = nil
				}
				completed <- result{response: response, err: err}
			}()
			select {
			case <-started:
			case <-time.After(5 * time.Second):
				t.Fatal("canceled HTTP/2 request did not reach the server")
			}
			var got result
			if afterHeaders {
				select {
				case got = <-completed:
					if got.err != nil {
						t.Fatalf("receive HTTP/2 response headers before cancellation: %v", got.err)
					}
				case <-time.After(5 * time.Second):
					t.Fatal("HTTP/2 response headers did not reach the caller")
				}
			}
			if err := capability.Close(); err != nil {
				t.Fatalf("close capability with cancelable HTTP/2 request: %v", err)
			}
			cancel()
			if afterHeaders {
				if closeErr := got.response.Body.Close(); closeErr != nil {
					t.Fatalf("close canceled HTTP/2 response body: %v", closeErr)
				}
			} else {
				select {
				case got = <-completed:
					if got.err == nil {
						t.Fatal("HTTP/2 request canceled before response headers unexpectedly succeeded")
					}
				case <-time.After(5 * time.Second):
					t.Fatal("canceled HTTP/2 request did not finish")
				}
			}
			select {
			case <-closed:
			case <-time.After(5 * time.Second):
				t.Fatal("HTTP/2 connection remained pooled after cancellation")
			}
		})
	}
}

func TestX509TransportCloseRejectsConnectionFromLateDialCompletion(t *testing.T) {
	fixture := newX509TransportFixture(t)
	template := fixture.transport(t, nil)
	dialStarted := make(chan struct{})
	releaseDial := make(chan struct{})
	clientConnection, serverConnection := net.Pipe()
	t.Cleanup(func() {
		_ = clientConnection.Close()
		_ = serverConnection.Close()
	})
	template.DialContext = func(context.Context, string, string) (net.Conn, error) {
		close(dialStarted)
		<-releaseDial
		return clientConnection, nil
	}
	capability := newX509Capability(t, template)
	ctx, cancel := context.WithCancel(t.Context())
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models").WithContext(ctx)
	completed := make(chan error, 1)
	go func() {
		response, err := capability.Do(request)
		if response != nil {
			_ = response.Body.Close()
		}
		completed <- err
	}()
	select {
	case <-dialStarted:
	case <-time.After(5 * time.Second):
		t.Fatal("late dial did not start")
	}
	cancel()
	select {
	case err := <-completed:
		if err == nil {
			t.Fatal("request canceled during its dial unexpectedly succeeded")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("request cancellation waited for its blocked dial")
	}
	if err := capability.Close(); err != nil {
		t.Fatalf("close capability before late dial completion: %v", err)
	}
	close(releaseDial)
	if err := serverConnection.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatalf("bound late-connection observation: %v", err)
	}
	if _, err := serverConnection.Read(make([]byte, 1)); !errors.Is(err, io.EOF) {
		t.Fatalf("connection returned after capability closure ended with %v, want EOF", err)
	}
}

func newX509HTTP2LifecycleServer(
	t *testing.T, fixture *x509TransportFixture, handler http.Handler,
) (*httptest.Server, <-chan struct{}) {
	t.Helper()
	server := httptest.NewUnstartedServer(handler)
	server.EnableHTTP2 = true
	server.TLS = &tls.Config{
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    fixture.trust,
		Certificates: []tls.Certificate{fixture.certificate(t, "synthetic server", []string{x509TransportAPI}, false)},
	}
	closed := make(chan struct{})
	var closeOnce sync.Once
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateClosed {
			closeOnce.Do(func() { close(closed) })
		}
	}
	server.Config.ErrorLog = log.New(io.Discard, "", 0)
	server.StartTLS()
	t.Cleanup(server.Close)
	return server, closed
}
