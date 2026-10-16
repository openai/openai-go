package live

import (
	"encoding/json"
	"fmt"
	"testing"
	"testing/synctest"
)

func TestTranscriptEventIDRetentionIsBounded(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		g, err := NewTranscriptGrouper(TranscriptGrouperOptions{})
		if err != nil {
			t.Fatal(err)
		}
		defer g.Close()
		for i := 0; i < 12000; i++ {
			kind := "session.input_transcript.delta"
			if i%2 == 1 {
				kind = "session.output_transcript.delta"
			}
			var event ServerEventUnion
			raw := fmt.Sprintf(`{"type":%q,"event_id":"retained_%d","delta":"word","start_ms":%d,"end_ms":%d}`, kind, i, i*5000, i*5000+100)
			if err := json.Unmarshal([]byte(raw), &event); err != nil {
				t.Fatal(err)
			}
			if err := g.Push(event); err != nil {
				t.Fatal(err)
			}
		}
		if len(g.seen) > 4096 {
			t.Fatalf("retained %d event IDs after finalized turns", len(g.seen))
		}
	})
}

func TestTranscriptDrainedQueueReleasesBurst(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var grouper *TranscriptGrouper
		var updates int
		var lastReason string
		push := func(id int, speaker string) {
			t.Helper()
			kind := "session.input_transcript.delta"
			if speaker == "assistant" {
				kind = "session.output_transcript.delta"
			}
			var event ServerEventUnion
			raw := fmt.Sprintf(`{"type":%q,"event_id":"burst_%d","delta":"word","start_ms":%d,"end_ms":%d}`,
				kind, id, id*5000, id*5000+100)
			if err := json.Unmarshal([]byte(raw), &event); err != nil {
				t.Fatal(err)
			}
			if err := grouper.Push(event); err != nil {
				t.Fatal(err)
			}
		}
		var err error
		grouper, err = NewTranscriptGrouper(TranscriptGrouperOptions{OnUpdate: func(update TranscriptUpdate) {
			updates++
			lastReason = update.Reason
			if updates == 1 {
				// Accumulate a backlog while a callback is in progress, then
				// close before that backlog has been delivered.
				for i := 2; i < 258; i++ {
					speaker := "user"
					if i%2 == 0 {
						speaker = "assistant"
					}
					push(i, speaker)
				}
				grouper.Close()
			}
		}})
		if err != nil {
			t.Fatal(err)
		}
		push(0, "user")
		push(1, "user")
		if updates < 256 || lastReason != "manual" {
			t.Fatalf("backlog was not delivered through Close: %d updates, last reason %q", updates, lastReason)
		}
		grouper.mu.Lock()
		defer grouper.mu.Unlock()
		if grouper.updates != nil {
			t.Fatalf("drained queue retains backing storage: len=%d cap=%d", len(grouper.updates), cap(grouper.updates))
		}
	})
}

// The user continues an open turn while the assistant emits distinct
// zero-duration contexts. None shares the user's start, so no overlap can
// suppress later overlap checks.
func BenchmarkTranscriptInterleavedPoints(b *testing.B) {
	for _, size := range []int{4096, 16384, 65536} {
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			for range b.N {
				g := newTranscriptGrouping(TranscriptGrouperOptions{}, "bench")
				g.ingest(transcriptFragment{speaker: "user", text: "x", start: 0, end: 1})
				for i := 1; i <= size; i++ {
					at := int64(i)
					g.ingest(transcriptFragment{speaker: "assistant", text: " ", start: at, end: at})
					g.ingest(transcriptFragment{speaker: "user", text: "x", start: 0, end: at + 1})
				}
				if g.buffered == nil || g.buffered.overlap {
					b.Fatal("distinct context points should remain nonoverlapping")
				}
			}
		})
	}
}
