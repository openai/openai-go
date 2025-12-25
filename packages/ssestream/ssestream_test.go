package ssestream_test

import (
	"errors"
	"testing"

	"github.com/Nordlys-Labs/openai-go/v3"
	"github.com/Nordlys-Labs/openai-go/v3/packages/ssestream"
)

// mockDecoder is a test helper that implements the Decoder interface
type mockDecoder struct {
	events []ssestream.Event
	index  int
	err    error
}

func (m *mockDecoder) Next() bool {
	if m.err != nil {
		return false
	}
	if m.index >= len(m.events) {
		return false
	}
	m.index++
	return true
}

func (m *mockDecoder) Event() ssestream.Event {
	if m.index == 0 {
		return ssestream.Event{}
	}
	return m.events[m.index-1]
}

func (m *mockDecoder) Close() error {
	return nil
}

func (m *mockDecoder) Err() error {
	return m.err
}

// testStruct is a simple struct for testing JSON unmarshaling
type testStruct struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

// ============================================================================
// TestPingEventHandling - Core functionality: ping events are skipped gracefully
// ============================================================================
// These tests verify the main fix: ping events (empty/whitespace data) are
// skipped without setting an error state, allowing streams to continue.

func TestPingEventHandling(t *testing.T) {
	t.Run("empty data event is skipped", func(t *testing.T) {
		events := []ssestream.Event{
			{Data: []byte{}}, // ping event
			{Data: []byte(`{"id":"1","data":"test"}`)},
		}

		decoder := &mockDecoder{events: events}
		stream := ssestream.NewStream[testStruct](decoder, nil)

		// Should skip ping event and process valid event
		if !stream.Next() {
			t.Fatal("Expected to process valid event after ping")
		}

		chunk := stream.Current()
		if chunk.ID != "1" {
			t.Errorf("Expected ID '1', got '%s'", chunk.ID)
		}

		// Stream should complete without error
		if stream.Next() {
			t.Error("Expected stream to be exhausted")
		}
		if stream.Err() != nil {
			t.Errorf("Stream should complete without error, got: %v", stream.Err())
		}
	})

	t.Run("whitespace-only data event is skipped", func(t *testing.T) {
		events := []ssestream.Event{
			{Data: []byte("   ")}, // ping event (whitespace)
			{Data: []byte(`{"id":"1","data":"test"}`)},
		}

		decoder := &mockDecoder{events: events}
		stream := ssestream.NewStream[testStruct](decoder, nil)

		if !stream.Next() {
			t.Fatal("Expected to process valid event after whitespace ping")
		}

		chunk := stream.Current()
		if chunk.ID != "1" {
			t.Errorf("Expected ID '1', got '%s'", chunk.ID)
		}

		if stream.Err() != nil {
			t.Errorf("Stream should complete without error, got: %v", stream.Err())
		}
	})

	t.Run("multiple consecutive ping events are all skipped", func(t *testing.T) {
		events := []ssestream.Event{
			{Data: []byte{}},       // ping
			{Data: []byte(" ")},    // ping
			{Data: []byte("\t\n")}, // ping
			{Data: []byte(`{"id":"1","data":"test"}`)},
		}

		decoder := &mockDecoder{events: events}
		stream := ssestream.NewStream[testStruct](decoder, nil)

		// Should skip all ping events and process valid event
		if !stream.Next() {
			t.Fatal("Expected to process valid event after multiple pings")
		}

		chunk := stream.Current()
		if chunk.ID != "1" {
			t.Errorf("Expected ID '1', got '%s'", chunk.ID)
		}

		if stream.Err() != nil {
			t.Errorf("Stream should complete without error, got: %v", stream.Err())
		}
	})

	t.Run("ping events at start, middle, and end are all skipped", func(t *testing.T) {
		events := []ssestream.Event{
			{Data: []byte{}}, // ping at start
			{Data: []byte(`{"id":"1","data":"test1"}`)},
			{Data: []byte(" ")}, // ping in middle
			{Data: []byte(`{"id":"2","data":"test2"}`)},
			{Data: []byte("\n")}, // ping at end
		}

		decoder := &mockDecoder{events: events}
		stream := ssestream.NewStream[testStruct](decoder, nil)

		expectedIDs := []string{"1", "2"}
		var receivedIDs []string

		for stream.Next() {
			chunk := stream.Current()
			receivedIDs = append(receivedIDs, chunk.ID)
		}

		if len(receivedIDs) != len(expectedIDs) {
			t.Errorf("Expected %d chunks, got %d", len(expectedIDs), len(receivedIDs))
		}
		for i, expected := range expectedIDs {
			if i < len(receivedIDs) && receivedIDs[i] != expected {
				t.Errorf("Chunk[%d]: expected ID '%s', got '%s'", i, expected, receivedIDs[i])
			}
		}

		// Stream should complete without error
		if stream.Err() != nil {
			t.Errorf("Stream should complete without error, got: %v", stream.Err())
		}
	})

	t.Run("stream with only ping events completes without error", func(t *testing.T) {
		events := []ssestream.Event{
			{Data: []byte{}},
			{Data: []byte(" ")},
			{Data: []byte("\t\n")},
		}

		decoder := &mockDecoder{events: events}
		stream := ssestream.NewStream[testStruct](decoder, nil)

		count := 0
		for stream.Next() {
			count++
		}

		if count != 0 {
			t.Errorf("Expected 0 chunks from ping-only stream, got %d", count)
		}

		// Critical: stream should complete without error, not stuck in error state
		if stream.Err() != nil {
			t.Errorf("Stream with only ping events should complete without error, got: %v", stream.Err())
		}
	})

	t.Run("ping event does not prevent processing subsequent valid events", func(t *testing.T) {
		// This is the critical test: ping events should not cause the stream
		// to get stuck in an error state, preventing further processing
		events := []ssestream.Event{
			{Data: []byte(`{"id":"1","data":"before"}`)},
			{Data: []byte{}}, // ping event
			{Data: []byte(`{"id":"2","data":"after"}`)},
		}

		decoder := &mockDecoder{events: events}
		stream := ssestream.NewStream[testStruct](decoder, nil)

		var receivedIDs []string
		for stream.Next() {
			chunk := stream.Current()
			receivedIDs = append(receivedIDs, chunk.ID)
		}

		expectedIDs := []string{"1", "2"}
		if len(receivedIDs) != len(expectedIDs) {
			t.Errorf("Expected %d chunks, got %d. Ping event prevented processing!", len(expectedIDs), len(receivedIDs))
		}
		for i, expected := range expectedIDs {
			if i < len(receivedIDs) && receivedIDs[i] != expected {
				t.Errorf("Chunk[%d]: expected ID '%s', got '%s'", i, expected, receivedIDs[i])
			}
		}

		// Stream must complete without error
		if stream.Err() != nil {
			t.Errorf("Ping event caused stream error (this is the bug we're fixing!): %v", stream.Err())
		}
	})
}

// ============================================================================
// TestStreamNormalOperation - Regression tests: normal streams work correctly
// ============================================================================
// These tests ensure the fix doesn't break existing functionality.

func TestStreamNormalOperation(t *testing.T) {
	t.Run("valid JSON events parse correctly", func(t *testing.T) {
		events := []ssestream.Event{
			{Data: []byte(`{"id":"1","data":"test1"}`)},
			{Data: []byte(`{"id":"2","data":"test2"}`)},
			{Data: []byte(`{"id":"3","data":"test3"}`)},
		}

		decoder := &mockDecoder{events: events}
		stream := ssestream.NewStream[testStruct](decoder, nil)

		expectedIDs := []string{"1", "2", "3"}
		var receivedIDs []string

		for stream.Next() {
			chunk := stream.Current()
			receivedIDs = append(receivedIDs, chunk.ID)
		}

		if len(receivedIDs) != len(expectedIDs) {
			t.Errorf("Expected %d chunks, got %d", len(expectedIDs), len(receivedIDs))
		}
		for i, expected := range expectedIDs {
			if i < len(receivedIDs) && receivedIDs[i] != expected {
				t.Errorf("Chunk[%d]: expected ID '%s', got '%s'", i, expected, receivedIDs[i])
			}
		}

		if stream.Err() != nil {
			t.Errorf("Unexpected error: %v", stream.Err())
		}
	})

	t.Run("[DONE] event terminates stream correctly", func(t *testing.T) {
		events := []ssestream.Event{
			{Data: []byte(`{"id":"1","data":"test1"}`)},
			{Data: []byte("[DONE]")},
			{Data: []byte(`{"id":"2","data":"test2"}`)}, // Should be skipped after [DONE]
		}

		decoder := &mockDecoder{events: events}
		stream := ssestream.NewStream[testStruct](decoder, nil)

		count := 0
		for stream.Next() {
			count++
		}

		if count != 1 {
			t.Errorf("Expected 1 chunk before [DONE], got %d", count)
		}

		if stream.Err() != nil {
			t.Errorf("Unexpected error: %v", stream.Err())
		}
	})

	t.Run("empty stream completes without error", func(t *testing.T) {
		decoder := &mockDecoder{events: []ssestream.Event{}}
		stream := ssestream.NewStream[testStruct](decoder, nil)

		if stream.Next() {
			t.Error("Expected empty stream to return false")
		}

		if stream.Err() != nil {
			t.Errorf("Unexpected error: %v", stream.Err())
		}
	})
}

// ============================================================================
// TestStreamErrorHandling - Error cases still work correctly
// ============================================================================
// These tests ensure that real errors (not ping events) still set error state.

func TestStreamErrorHandling(t *testing.T) {
	t.Run("error field in JSON sets error state", func(t *testing.T) {
		events := []ssestream.Event{
			{Data: []byte(`{"id":"1","data":"test1"}`)},
			{Data: []byte(`{"error":"something went wrong"}`)},
		}

		decoder := &mockDecoder{events: events}
		stream := ssestream.NewStream[testStruct](decoder, nil)

		// First event should succeed
		if !stream.Next() {
			t.Fatal("Expected first event to succeed")
		}

		// Second event with error field should stop stream
		if stream.Next() {
			t.Error("Expected stream to stop on error event")
		}

		if stream.Err() == nil {
			t.Error("Expected error but got none")
		}
	})

	t.Run("malformed JSON (non-empty) sets error state", func(t *testing.T) {
		// Important: non-empty invalid JSON should error, unlike ping events
		events := []ssestream.Event{
			{Data: []byte(`{"id":"1","data":"test1"}`)},
			{Data: []byte(`{invalid json}`)}, // Non-empty but invalid JSON
		}

		decoder := &mockDecoder{events: events}
		stream := ssestream.NewStream[testStruct](decoder, nil)

		// First event should succeed
		if !stream.Next() {
			t.Fatal("Expected first event to succeed")
		}

		// Second event should fail with JSON parsing error
		if stream.Next() {
			t.Error("Expected stream to stop on malformed JSON")
		}

		if stream.Err() == nil {
			t.Error("Expected JSON parsing error but got none")
		}
	})

	t.Run("decoder errors propagate correctly", func(t *testing.T) {
		decoder := &mockDecoder{
			events: []ssestream.Event{
				{Data: []byte(`{"id":"1","data":"test1"}`)},
				{Data: []byte(`{"id":"2","data":"test2"}`)},
			},
		}
		stream := ssestream.NewStream[testStruct](decoder, nil)

		// First event should succeed
		if !stream.Next() {
			t.Fatal("Expected first event to succeed")
		}

		// Set decoder error after first event
		decoder.err = errors.New("decoder error")

		// Next call should fail with decoder error
		if stream.Next() {
			t.Error("Expected stream to stop on decoder error")
		}

		if stream.Err() == nil {
			t.Error("Expected decoder error but got none")
		}
	})

	t.Run("initial decoder error prevents stream processing", func(t *testing.T) {
		decoder := &mockDecoder{
			events: []ssestream.Event{},
			err:    errors.New("initial error"),
		}
		stream := ssestream.NewStream[testStruct](decoder, errors.New("initial error"))

		// Stream should immediately fail
		if stream.Next() {
			t.Error("Expected stream to fail immediately")
		}

		if stream.Err() == nil {
			t.Error("Expected error but got none")
		}
	})
}

// ============================================================================
// TestStreamRealWorldScenarios - Real-world use cases with actual SDK types
// ============================================================================
// These tests use actual OpenAI SDK types to verify the fix works in practice.

func TestStreamRealWorldScenarios(t *testing.T) {
	t.Run("chat completion stream with ping event in middle", func(t *testing.T) {
		// Simulates real-world scenario from Anthropic/OpenAI API where a ping event
		// (e.g., ": ping - 2025-12-25 17:02:01.661905") appears in the middle of a stream.
		// The decoder parses this as an event with empty/whitespace data.
		events := []ssestream.Event{
			{Data: []byte(`{"id":"msg_011Rmwern9PBanyrMbeDyRD1","choices":[{"index":0,"delta":{"content":"CPU Data"},"finish_reason":null}],"created":1766682106,"model":"claude-sonnet-4-5-20250929","object":"chat.completion.chunk"}`)},
			{Data: []byte(`{"id":"msg_011Rmwern9PBanyrMbeDyRD1","choices":[{"index":0,"delta":{"content":" Fetching: Cache"},"finish_reason":null}],"created":1766682106,"model":"claude-sonnet-4-5-20250929","object":"chat.completion.chunk"}`)},
			{Data: []byte(`{"id":"msg_011Rmwern9PBanyrMbeDyRD1","choices":[{"index":0,"delta":{"content":" Lines"},"finish_reason":null}],"created":1766682106,"model":"claude-sonnet-4-5-20250929","object":"chat.completion.chunk"}`)},
			{Data: []byte{}}, // Ping event (empty data) - simulates ": ping - 2025-12-25 17:02:01.661905"
			{Data: []byte(`{"id":"msg_011Rmwern9PBanyrMbeDyRD1","choices":[{"index":0,"delta":{"content":" and Cache Line"},"finish_reason":null}],"created":1766682106,"model":"claude-sonnet-4-5-20250929","object":"chat.completion.chunk"}`)},
			{Data: []byte(`{"id":"msg_011Rmwern9PBanyrMbeDyRD1","choices":[{"index":0,"delta":{"content":" Bouncing"},"finish_reason":null}],"created":1766682106,"model":"claude-sonnet-4-5-20250929","object":"chat.completion.chunk"}`)},
			{Data: []byte(`{"id":"msg_011Rmwern9PBanyrMbeDyRD1","choices":[{"index":0,"delta":{"content":"\n\n## Word"},"finish_reason":null}],"created":1766682106,"model":"claude-sonnet-4-5-20250929","object":"chat.completion.chunk"}`)},
		}

		decoder := &mockDecoder{events: events}
		stream := ssestream.NewStream[openai.ChatCompletionChunk](decoder, nil)

		// Should process all valid events, skipping the ping event
		expectedCount := 6 // All events except the ping event
		count := 0
		var ids []string
		var contents []string

		for stream.Next() {
			count++
			chunk := stream.Current()
			ids = append(ids, chunk.ID)
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				contents = append(contents, chunk.Choices[0].Delta.Content)
			}
		}

		if count != expectedCount {
			t.Errorf("Expected %d chunks (skipping ping event), got %d", expectedCount, count)
		}

		// Verify all chunks have the same ID (from the same message)
		expectedID := "msg_011Rmwern9PBanyrMbeDyRD1"
		for i, id := range ids {
			if id != expectedID {
				t.Errorf("Chunk[%d]: expected ID '%s', got '%s'", i, expectedID, id)
			}
		}

		// Verify content chunks were processed correctly
		expectedContents := []string{
			"CPU Data",
			" Fetching: Cache",
			" Lines",
			" and Cache Line",
			" Bouncing",
			"\n\n## Word",
		}
		if len(contents) != len(expectedContents) {
			t.Errorf("Expected %d content chunks, got %d", len(expectedContents), len(contents))
		}
		for i, expected := range expectedContents {
			if i < len(contents) && contents[i] != expected {
				t.Errorf("Content[%d]: expected '%s', got '%s'", i, expected, contents[i])
			}
		}

		// Critical: Stream should complete without error (ping event should be skipped gracefully)
		if stream.Err() != nil {
			t.Errorf("Stream should complete without error, got: %v", stream.Err())
		}
	})

	t.Run("chat completion stream processes all chunks despite ping events", func(t *testing.T) {
		// Test that ping events don't interrupt the stream processing
		events := []ssestream.Event{
			{Data: []byte(`{"id":"msg_123","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}],"created":1766682106,"model":"gpt-4","object":"chat.completion.chunk"}`)},
			{Data: []byte{}}, // ping
			{Data: []byte(`{"id":"msg_123","choices":[{"index":0,"delta":{"content":" World"},"finish_reason":null}],"created":1766682106,"model":"gpt-4","object":"chat.completion.chunk"}`)},
			{Data: []byte(" ")}, // ping
			{Data: []byte(`{"id":"msg_123","choices":[{"index":0,"delta":{"content":"!"},"finish_reason":null}],"created":1766682106,"model":"gpt-4","object":"chat.completion.chunk"}`)},
		}

		decoder := &mockDecoder{events: events}
		stream := ssestream.NewStream[openai.ChatCompletionChunk](decoder, nil)

		var contents []string
		for stream.Next() {
			chunk := stream.Current()
			if len(chunk.Choices) > 0 {
				contents = append(contents, chunk.Choices[0].Delta.Content)
			}
		}

		expectedContents := []string{"Hello", " World", "!"}
		if len(contents) != len(expectedContents) {
			t.Errorf("Expected %d content chunks, got %d. Ping events interrupted processing!", len(expectedContents), len(contents))
		}
		for i, expected := range expectedContents {
			if i < len(contents) && contents[i] != expected {
				t.Errorf("Content[%d]: expected '%s', got '%s'", i, expected, contents[i])
			}
		}

		if stream.Err() != nil {
			t.Errorf("Stream should complete without error, got: %v", stream.Err())
		}
	})
}
