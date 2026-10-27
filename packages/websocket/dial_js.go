package websocket

import (
	"context"
	"errors"
	"net/http"

	wire "github.com/coder/websocket"
)

func dialSocket(context.Context, *http.Request, *http.Client) (*wire.Conn, *http.Response, error) {
	return nil, nil, errors.New("websocket: SDK WebSocket connections require native HTTP transport; js/wasm is unsupported")
}
