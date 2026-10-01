package bedrock

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/option"
)

func TestUnsafeHTTPPreservesImmediateExpectContinueBody(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()
	serverResult := make(chan error, 1)
	go func() {
		serverResult <- func() error {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				return acceptErr
			}
			defer func() { _ = conn.Close() }()
			if deadlineErr := conn.SetDeadline(time.Now().Add(5 * time.Second)); deadlineErr != nil {
				return deadlineErr
			}
			req, readErr := http.ReadRequest(bufio.NewReader(conn))
			if readErr != nil {
				return readErr
			}
			defer func() { _ = req.Body.Close() }()
			if req.Header.Get("Expect") != "100-continue" {
				return fmt.Errorf("missing Expect header")
			}
			// A raw local server deliberately sends no 100 response. Zero
			// ExpectContinueTimeout must send the body immediately, without the
			// one-second wait used by the SDK's default development transport.
			if deadlineErr := conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond)); deadlineErr != nil {
				return deadlineErr
			}
			body, bodyErr := io.ReadAll(req.Body)
			if bodyErr != nil {
				return bodyErr
			}
			if string(body) != "synthetic body" {
				return fmt.Errorf("unexpected request body")
			}
			_, writeErr := io.WriteString(conn, "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: 2\r\nConnection: close\r\n\r\n{}")
			return writeErr
		}()
	}()
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	client, err := NewClient(ctx, Config{APIKey: "synthetic-bearer", BaseURL: "http://" + listener.Addr().String(), UnsafeAllowHTTP: true},
		option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{Transport: &http.Transport{ExpectContinueTimeout: 0}}))
	if err != nil {
		t.Fatal(err)
	}
	var response map[string]any
	requestErr := client.Post(ctx, "/responses", strings.NewReader("synthetic body"), &response, option.WithHeader("Expect", "100-continue"))
	if requestErr != nil {
		t.Errorf("immediate request body failed: %v", requestErr)
	}
	select {
	case serverErr := <-serverResult:
		if serverErr != nil {
			t.Error(serverErr)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}
