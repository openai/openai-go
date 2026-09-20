//go:build !js

package websocket

import (
	"context"
	"net/http"

	wire "github.com/coder/websocket"
)

func dialSocket(ctx context.Context, request *http.Request, client *http.Client) (*wire.Conn, *http.Response, error) {
	return wire.Dial(ctx, request.URL.String(), &wire.DialOptions{HTTPClient: client, HTTPHeader: request.Header, Host: request.Host, CompressionMode: wire.CompressionDisabled})
}
