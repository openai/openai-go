package ssestream

import (
	"bytes"
	"net/http"
	"testing"
)

type mockReadCloser struct {
	*bytes.Reader
}

func (m mockReadCloser) Close() error {
	return nil
}

// TestStream_EmptyEvents tests that the stream correctly handles empty SSE events
// (e.g., from retry: directives or comment lines) without crashing on JSON unmarshal
func TestStream_EmptyEvents(t *testing.T) {
	// Simulate SSE stream with retry directive that creates empty event
	sseData := `retry: 3000

data: {"id":"msg_01ABC","type":"content_block_delta","delta":{"type":"text","text":"Hello"}}

data: [DONE]

`

	resp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:       mockReadCloser{bytes.NewReader([]byte(sseData))},
	}

	decoder := NewDecoder(resp)
	if decoder == nil {
		t.Fatal("Expected decoder to be created, got nil")
	}

	type testMsg struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Delta struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"delta"`
	}

	stream := NewStream[testMsg](decoder, nil)

	// Should successfully iterate without crashing on empty event
	var receivedMessages int
	for stream.Next() {
		msg := stream.Current()
		receivedMessages++

		if msg.ID != "msg_01ABC" {
			t.Errorf("Expected ID 'msg_01ABC', got '%s'", msg.ID)
		}
		if msg.Delta.Text != "Hello" {
			t.Errorf("Expected text 'Hello', got '%s'", msg.Delta.Text)
		}
	}

	if err := stream.Err(); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if receivedMessages != 1 {
		t.Errorf("Expected 1 message, got %d", receivedMessages)
	}
}

// TestStream_OnlyRetryDirective tests stream with only retry directive (no data events)
func TestStream_OnlyRetryDirective(t *testing.T) {
	sseData := `retry: 3000

`

	resp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:       mockReadCloser{bytes.NewReader([]byte(sseData))},
	}

	decoder := NewDecoder(resp)
	type testMsg struct {
		ID string `json:"id"`
	}
	stream := NewStream[testMsg](decoder, nil)

	// Should handle gracefully without any messages
	var count int
	for stream.Next() {
		count++
	}

	if err := stream.Err(); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 messages, got %d", count)
	}
}

// TestStream_MultipleEmptyEvents tests handling of multiple empty events
func TestStream_MultipleEmptyEvents(t *testing.T) {
	sseData := `retry: 3000

: comment line

data: {"id":"1","text":"first"}

retry: 5000

data: {"id":"2","text":"second"}

`

	resp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:       mockReadCloser{bytes.NewReader([]byte(sseData))},
	}

	decoder := NewDecoder(resp)
	type testMsg struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	}
	stream := NewStream[testMsg](decoder, nil)

	messages := []testMsg{}
	for stream.Next() {
		messages = append(messages, stream.Current())
	}

	if err := stream.Err(); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(messages))
	}

	if messages[0].ID != "1" || messages[0].Text != "first" {
		t.Errorf("First message incorrect: %+v", messages[0])
	}

	if messages[1].ID != "2" || messages[1].Text != "second" {
		t.Errorf("Second message incorrect: %+v", messages[1])
	}
}
