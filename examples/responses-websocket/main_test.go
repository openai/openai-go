package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v4"
	"github.com/openai/openai-go/v4/option"
)

func TestRunRequiresCompletedResponses(t *testing.T) {
	for _, tc := range []struct {
		name, first, second, wantError string
		wantCreates                    int
	}{
		{"failed first", "failed", "", "first response did not complete", 1},
		{"incomplete first", "incomplete", "", "first response did not complete", 1},
		{"failed continuation", "completed", "failed", "continuation did not complete", 2},
		{"incomplete continuation", "completed", "incomplete", "continuation did not complete", 2},
		{"completed", "completed", "completed", "", 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requests := make(chan []map[string]any, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					t.Errorf("upgrade: %v", err)
					return
				}
				defer func() { _ = socket.CloseNow() }()
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				var received []map[string]any
				defer func() { requests <- received }()
				for i, status := range []string{tc.first, tc.second} {
					_, data, readErr := socket.Read(ctx)
					if readErr != nil {
						return
					}
					var request map[string]any
					if json.Unmarshal(data, &request) != nil {
						t.Error("invalid request JSON")
						return
					}
					received = append(received, request)
					if status == "" {
						return
					}
					frame := fmt.Sprintf(`{"type":"response.%s","response":{"id":"resp_%d","status":"%s"}}`, status, i+1, status)
					if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
						t.Errorf("write response: %v", err)
						return
					}
				}
				_, _, _ = socket.Read(ctx)
			}))
			defer server.Close()
			client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"))
			err := runWithClient(client)
			if tc.wantError == "" && err != nil || tc.wantError != "" && (err == nil || err.Error() != tc.wantError) {
				t.Fatalf("runWithClient() error = %v, want %q", err, tc.wantError)
			}
			got := <-requests
			if len(got) != tc.wantCreates {
				t.Fatalf("sent %d requests, want %d", len(got), tc.wantCreates)
			}
			if len(got) == 2 && got[1]["previous_response_id"] != "resp_1" {
				t.Errorf("continuation previous_response_id = %v, want resp_1", got[1]["previous_response_id"])
			}
		})
	}
}
