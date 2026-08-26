package auth_test

import (
	"context"
	"net"
	"net/http"
	"strings"
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
