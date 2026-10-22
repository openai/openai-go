package auth

import (
	"context"
	"net"
	"net/http"
	"testing"
)

// MetadataProviderWithClientForTest preserves response-parser tests without
// exposing a production override for the metadata transport.
func MetadataProviderWithClientForTest(provider SubjectTokenProvider, client HTTPDoer) SubjectTokenProvider {
	switch p := provider.(type) {
	case *azureManagedIdentityTokenProvider:
		p.metadataClient = client
	case *gcpIDTokenProvider:
		p.metadataClient = client
	default:
		panic("not a metadata provider")
	}
	return provider
}

// MetadataProviderWithDialForTest routes the standard transport to a local
// fixture while checking the destination constructed by the provider.
func MetadataProviderWithDialForTest(t *testing.T, provider SubjectTokenProvider, address string) SubjectTokenProvider {
	t.Helper()
	client := metadataClientForTest(provider)
	client.Transport.(*http.Transport).DialContext = func(ctx context.Context, network, destination string) (net.Conn, error) {
		if network != "tcp" || destination != "169.254.169.254:80" {
			panic("unexpected metadata destination")
		}
		return (&net.Dialer{}).DialContext(ctx, network, address)
	}
	t.Cleanup(client.CloseIdleConnections)
	return MetadataProviderWithClientForTest(provider, client)
}

func metadataClientForTest(provider SubjectTokenProvider) *http.Client {
	switch p := provider.(type) {
	case *azureManagedIdentityTokenProvider:
		return p.metadataClient.(*http.Client)
	case *gcpIDTokenProvider:
		return p.metadataClient.(*http.Client)
	default:
		panic("not a metadata provider")
	}
}
