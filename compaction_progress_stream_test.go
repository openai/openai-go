package openai_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func compactionProgressClient(t *testing.T) openai.Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/responses" {
			t.Errorf("NewStreaming request = %s %s, want POST /responses", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "event: response.compaction.compacting\ndata: {\"type\":\"response.compaction.compacting\",\"item_id\":\"cmp_synthetic\",\"output_index\":2,\"sequence_number\":7,\"agent\":{\"agent_name\":\"/root\"}}\n\n"+
			"event: response.completed\ndata: {\"type\":\"response.completed\",\"sequence_number\":8,\"response\":{\"id\":\"resp_synthetic\",\"status\":\"completed\",\"output\":[]}}\n\n"+
			"data: [DONE]\n\n")
	}))
	t.Cleanup(server.Close)
	return openai.NewClient(option.WithBaseURL(server.URL+"/"), option.WithAPIKey("synthetic-key"), option.WithMaxRetries(0))
}

func TestResponseCompactionProgressStream(t *testing.T) {
	client := compactionProgressClient(t)
	stream := client.Responses.NewStreaming(t.Context(), responses.ResponseNewParams{})
	t.Cleanup(func() { _ = stream.Close() })
	if !stream.Next() {
		t.Fatalf("NewStreaming first Next() = false, error %v, want progress event", stream.Err())
	}
	progress, ok := stream.Current().AsAny().(responses.ResponseCompactionCompactingEvent)
	if !ok {
		t.Fatalf("progress AsAny() = %T, want ResponseCompactionCompactingEvent", stream.Current().AsAny())
	}
	if progress.ItemID != "cmp_synthetic" || progress.OutputIndex != 2 || progress.SequenceNumber != 7 {
		t.Errorf("progress fields = (%q, %d, %d), want (cmp_synthetic, 2, 7)", progress.ItemID, progress.OutputIndex, progress.SequenceNumber)
	}
	if !stream.Next() {
		t.Fatalf("NewStreaming second Next() = false, error %v, want completion", stream.Err())
	}
	completed, ok := stream.Current().AsAny().(responses.ResponseCompletedEvent)
	if !ok || completed.SequenceNumber != 8 || completed.Response.ID != "resp_synthetic" || completed.Response.Status != "completed" {
		t.Errorf("completion = %#v, variant %T, want completed resp_synthetic at sequence 8", completed, stream.Current().AsAny())
	}
	if stream.Next() || stream.Err() != nil {
		t.Errorf("NewStreaming after completion error = %v, want clean end", stream.Err())
	}
}

func TestBetaResponseCompactionProgressStream(t *testing.T) {
	client := compactionProgressClient(t)
	stream := client.Beta.Responses.NewStreaming(t.Context(), openai.BetaResponseNewParams{})
	t.Cleanup(func() { _ = stream.Close() })
	if !stream.Next() {
		t.Fatalf("NewStreaming first Next() = false, error %v, want progress event", stream.Err())
	}
	progress, ok := stream.Current().AsAny().(openai.BetaResponseCompactionCompactingEvent)
	if !ok {
		t.Fatalf("progress AsAny() = %T, want BetaResponseCompactionCompactingEvent", stream.Current().AsAny())
	}
	if progress.ItemID != "cmp_synthetic" || progress.OutputIndex != 2 || progress.SequenceNumber != 7 || progress.Agent.AgentName != "/root" {
		t.Errorf("progress fields = (%q, %d, %d, %q), want (cmp_synthetic, 2, 7, /root)", progress.ItemID, progress.OutputIndex, progress.SequenceNumber, progress.Agent.AgentName)
	}
	if !stream.Next() {
		t.Fatalf("NewStreaming second Next() = false, error %v, want completion", stream.Err())
	}
	completed, ok := stream.Current().AsAny().(openai.BetaResponseCompletedEvent)
	if !ok || completed.SequenceNumber != 8 || completed.Response.ID != "resp_synthetic" || completed.Response.Status != "completed" {
		t.Errorf("completion = %#v, variant %T, want completed resp_synthetic at sequence 8", completed, stream.Current().AsAny())
	}
	if stream.Next() || stream.Err() != nil {
		t.Errorf("NewStreaming after completion error = %v, want clean end", stream.Err())
	}
}
