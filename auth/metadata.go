package auth

import (
	"net"
	"net/http"
	"time"
)

const (
	// Both providers document this link-local IPv4 endpoint. Like Google's
	// compute/metadata client, use the IP in the URL rather than resolving a
	// hostname or overriding socket destinations in a custom transport.
	// https://learn.microsoft.com/azure/virtual-machines/instance-metadata-service
	// https://cloud.google.com/compute/docs/metadata/querying-metadata
	azureMetadataURL       = "http://169.254.169.254/metadata/identity/oauth2/token"
	gcpMetadataURL         = "http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/identity"
	metadataRequestTimeout = 5 * time.Second
)

// newMetadataHTTPClient owns its transport instead of inheriting API proxies,
// global transports, or caller code. Providers construct requests from the fixed
// URLs above; rejecting redirects keeps those requests at their original endpoint.
func newMetadataHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: nil,
			DialContext: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          2,
			MaxIdleConnsPerHost:   2,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: metadataRequestTimeout,
			DisableCompression:    true,
		},
		CheckRedirect: func(*http.Request, []*http.Request) error { return errWorkloadIdentityRedirect },
	}
}
