package websocket

import (
	"context"
	"errors"
	"net/http"

	wire "github.com/coder/websocket"
)

func dialSocket(context.Context, *http.Request, *http.Client) (*wire.Conn, *http.Response, error) {
	return nil, nil, errors.New("websocket: Responses connections require a transport that supports custom HTTP headers; js/wasm is unsupported")
}
