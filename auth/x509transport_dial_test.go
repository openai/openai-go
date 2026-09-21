package auth_test

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/auth"
)

var _ = map[auth.X509Transport]struct{}{}

func TestX509TransportPreservesDeprecatedDialBehindTrackedAdapter(t *testing.T) {
	fixture := newX509TransportFixture(t)
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	contextDial := template.DialContext
	template.DialContext = nil
	var deprecatedCalls atomic.Int32
	//nolint:staticcheck // This guards source-compatible support for the original X.509 transport contract.
	template.Dial = func(network, address string) (net.Conn, error) {
		deprecatedCalls.Add(1)
		return contextDial(context.Background(), network, address)
	}
	capability := newX509Capability(t, template)
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	response, err := capability.Do(request)
	if err != nil {
		t.Fatalf("request with preserved deprecated dialer: %v", err)
	}
	if closeErr := response.Body.Close(); closeErr != nil {
		t.Fatalf("close response using preserved deprecated dialer: %v", closeErr)
	}
	if got := deprecatedCalls.Load(); got != 1 {
		t.Errorf("preserved deprecated dialer calls = %d, want 1", got)
	}
}

func TestX509TransportIgnoresDeprecatedDialShadowedByDialContext(t *testing.T) {
	fixture := newX509TransportFixture(t)
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	var deprecatedCalls atomic.Int32
	//nolint:staticcheck // The regression must prove that net/http's shadowed deprecated hook is never copied or called.
	template.Dial = func(string, string) (net.Conn, error) {
		deprecatedCalls.Add(1)
		return nil, errors.New("shadowed deprecated dialer was called")
	}
	capability := newX509Capability(t, template)
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	response, err := capability.Do(request)
	if err != nil {
		t.Fatalf("request with shadowed deprecated dialer: %v", err)
	}
	if closeErr := response.Body.Close(); closeErr != nil {
		t.Fatalf("close response using context-aware dialer: %v", closeErr)
	}
	if got := deprecatedCalls.Load(); got != 0 {
		t.Errorf("shadowed deprecated dialer calls = %d, want 0", got)
	}
}

func TestX509TransportClosesConnectionReturnedWithDialError(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		name := "context-aware"
		if legacy {
			name = "legacy"
		}
		t.Run(name, func(t *testing.T) {
			fixture := newX509TransportFixture(t)
			template := fixture.transport(t, nil)
			connection, peer := net.Pipe()
			t.Cleanup(func() { _ = peer.Close() })
			setX509TestDialer(template, legacy, func() (net.Conn, error) {
				return connection, errors.New("synthetic dial failure with connection")
			})
			capability := newX509Capability(t, template)
			request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
			response, err := capability.Do(request)
			if response != nil {
				_ = response.Body.Close()
			}
			if err == nil {
				t.Fatal("malformed dialer returned a connection and error without failing the request")
			}
			readResult := make(chan error, 1)
			go func() {
				_, err := peer.Read(make([]byte, 1))
				readResult <- err
			}()
			select {
			case err := <-readResult:
				if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrClosedPipe) {
					t.Fatalf("connection returned with dial error remained open: %v", err)
				}
			case <-time.After(time.Second):
				t.Fatal("connection returned with dial error was not closed")
			}
		})
	}
}

func TestX509TransportRejectsNilLikeDialConnection(t *testing.T) {
	var typedNil *net.TCPConn
	for _, test := range []struct {
		name       string
		legacy     bool
		connection net.Conn
	}{
		{name: "context-aware nil"},
		{name: "context-aware typed nil", connection: typedNil},
		{name: "legacy nil", legacy: true},
		{name: "legacy typed nil", legacy: true, connection: typedNil},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newX509TransportFixture(t)
			template := fixture.transport(t, nil)
			setX509TestDialer(template, test.legacy, func() (net.Conn, error) {
				return test.connection, nil
			})
			capability := newX509Capability(t, template)
			request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
			response, err := capability.Do(request)
			if response != nil {
				_ = response.Body.Close()
			}
			if err == nil {
				t.Fatal("malformed dialer returned a nil-like connection without an error")
			}
		})
	}
}

func TestX509TransportPreservesDeadlineIdentityForLiveContextTimeout(t *testing.T) {
	fixture := newX509TransportFixture(t)
	template := fixture.transport(t, nil)
	template.DialContext = func(context.Context, string, string) (net.Conn, error) {
		return nil, context.DeadlineExceeded
	}
	capability := newX509Capability(t, template)
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	response, err := capability.Do(request)
	if response != nil {
		_ = response.Body.Close()
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("live-context transport timeout error = %v, want context.DeadlineExceeded", err)
	}
	var timeout interface{ Timeout() bool }
	if !errors.As(err, &timeout) || !timeout.Timeout() {
		t.Errorf("live-context transport timeout did not retain net.Error timeout semantics: %v", err)
	}
}

func TestX509TransportWaitsForAvailableDialCapacity(t *testing.T) {
	const callers = 33
	fixture := newX509TransportFixture(t)
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	template.DisableKeepAlives = true
	dialContext := template.DialContext
	release := make(chan struct{})
	capacityReached := make(chan struct{})
	var started atomic.Int32
	template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		if started.Add(1) == callers-1 {
			close(capacityReached)
		}
		<-release
		return dialContext(ctx, network, address)
	}
	capability := newX509Capability(t, template)

	start := make(chan struct{})
	results := make(chan error, callers)
	for range callers {
		go func() {
			<-start
			request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
			response, err := capability.Do(request)
			if response != nil {
				_ = response.Body.Close()
			}
			results <- err
		}()
	}
	close(start)
	select {
	case <-capacityReached:
	case <-time.After(5 * time.Second):
		close(release)
		t.Fatal("requests did not fill the bounded dial admission capacity")
	}
	close(release)
	for range callers {
		select {
		case err := <-results:
			if err != nil {
				t.Errorf("request waiting for dial admission failed: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("request waiting for dial admission did not complete after capacity became available")
		}
	}
	if got := started.Load(); got != callers {
		t.Errorf("dial calls = %d, want %d", got, callers)
	}
}

func TestX509TransportBoundsNonCooperativeDialers(t *testing.T) {
	const maximumConcurrentDials = 32
	fixture := newX509TransportFixture(t)
	template := fixture.transport(t, nil)
	release := make(chan struct{})
	var dialers sync.WaitGroup
	var started atomic.Int32
	reachedCapacity := make(chan struct{})
	template.DialContext = func(context.Context, string, string) (net.Conn, error) {
		dialers.Add(1)
		defer dialers.Done()
		if started.Add(1) == maximumConcurrentDials {
			close(reachedCapacity)
		}
		<-release
		return nil, errors.New("synthetic non-cooperative dialer released")
	}
	capability := newX509Capability(t, template)

	const callers = maximumConcurrentDials + 8
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	start := make(chan struct{})
	results := make(chan error, callers)
	var calls sync.WaitGroup
	calls.Add(callers)
	for range callers {
		go func() {
			defer calls.Done()
			<-start
			request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models").WithContext(ctx)
			response, err := capability.Do(request)
			if response != nil {
				_ = response.Body.Close()
			}
			results <- err
		}()
	}
	close(start)
	select {
	case <-reachedCapacity:
	case <-time.After(5 * time.Second):
		close(release)
		t.Fatal("non-cooperative dialers did not reach the bounded admission limit")
	}
	cancel()
	calls.Wait()
	close(results)
	if got := started.Load(); got != maximumConcurrentDials {
		t.Errorf("non-cooperative dialers started = %d, want bounded %d", got, maximumConcurrentDials)
	}
	for err := range results {
		if err == nil {
			t.Error("request with a non-cooperative dialer unexpectedly succeeded")
		}
	}
	close(release)
	dialers.Wait()
}

func setX509TestDialer(
	template *http.Transport,
	legacy bool,
	dial func() (net.Conn, error),
) {
	if legacy {
		template.DialContext = nil
		//nolint:staticcheck // The regression must exercise the accepted legacy dialer path.
		template.Dial = func(string, string) (net.Conn, error) { return dial() }
		return
	}
	template.DialContext = func(context.Context, string, string) (net.Conn, error) { return dial() }
}
