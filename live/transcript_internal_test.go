package live

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

func TestTranscriptDeferredSeparatorsCoalesce(t *testing.T) {
	for _, prefix := range []string{"", "okay"} {
		t.Run(prefix, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var updates []TranscriptUpdate
				g, err := NewTranscriptGrouper(TranscriptGrouperOptions{OnUpdate: func(u TranscriptUpdate) {
					updates = append(updates, u)
				}})
				if err != nil {
					t.Fatal(err)
				}
				defer g.Close()
				push := func(id int, text string, end int) {
					t.Helper()
					var event ServerEventUnion
					raw := fmt.Sprintf(`{"type":"session.output_transcript.delta","event_id":"compact_%d","delta":%q,"start_ms":0,"end_ms":%d}`, id, text, end)
					if decodeErr := json.Unmarshal([]byte(raw), &event); decodeErr != nil {
						t.Fatal(decodeErr)
					}
					if pushErr := g.Push(event); pushErr != nil {
						t.Fatal(pushErr)
					}
				}
				if prefix != "" {
					push(-1, prefix, 100)
				}
				for i := 0; i < 12000; i++ {
					push(i, " ", 100)
				}
				if len(g.pending) > 2 {
					t.Fatalf("retained %d deferred fragments for one source start", len(g.pending))
				}
				time.Sleep(5 * time.Second)
				synctest.Wait()
				if len(updates) != 0 {
					t.Fatal("published an acknowledgment before source context")
				}
				push(12000, "more", 200)
				g.Close()
				last := updates[len(updates)-1]
				if last.Reason != "manual" || last.Segment.Text != prefix+strings.Repeat(" ", 12000)+"more" ||
					last.Segment.StartMs != 0 || last.Segment.EndMs != 200 {
					t.Fatalf("lost deferred text or timing: %#v", last)
				}
			})
		})
	}
}

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

func TestTranscriptReentrantBackpressureReleasesBurst(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var grouper *TranscriptGrouper
		var updates, accepted, finals int
		var lastReason string
		overloaded := false
		push := func(id int, speaker string) error {
			t.Helper()
			kind := "session.input_transcript.delta"
			if speaker == "assistant" {
				kind = "session.output_transcript.delta"
			}
			var event ServerEventUnion
			raw := fmt.Sprintf(`{"type":%q,"event_id":"burst_%d","delta":"word","start_ms":%d,"end_ms":%d}`, kind, id, id*5000, id*5000+100)
			if err := json.Unmarshal([]byte(raw), &event); err != nil {
				t.Fatal(err)
			}
			return grouper.Push(event)
		}
		checkBound := func() {
			grouper.mu.Lock()
			defer grouper.mu.Unlock()
			// The final accepted event and its unfinished turns need at most
			// ten completion snapshots beyond the admission threshold.
			if len(grouper.updates) > 1024+10 {
				t.Fatalf("unbounded callback backlog: %d", len(grouper.updates))
			}
		}
		var err error
		grouper, err = NewTranscriptGrouper(TranscriptGrouperOptions{OnUpdate: func(update TranscriptUpdate) {
			updates++
			lastReason = update.Reason
			if update.Reason != "" {
				finals++
			}
			if updates == 1 {
				for i := 2; i < 2048; i++ {
					speaker := "user"
					if i%2 == 0 {
						speaker = "assistant"
					}
					if pushErr := push(i, speaker); pushErr != nil {
						if !errors.Is(pushErr, ErrTranscriptBackpressure) {
							t.Fatal(pushErr)
						}
						overloaded = true
						break
					}
					accepted++
				}
				checkBound()
				grouper.Close()
				checkBound()
			}
		}})
		if err != nil {
			t.Fatal(err)
		}
		if err := push(0, "user"); err != nil {
			t.Fatal(err)
		}
		if err := push(1, "user"); err != nil {
			t.Fatal(err)
		}
		if !overloaded || finals != accepted+1 || lastReason != "manual" {
			t.Fatalf("lost accepted turns or failed to bound backlog: overloaded=%v accepted=%d finals=%d reason=%q", overloaded, accepted, finals, lastReason)
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

func TestTranscriptCallbackPanicDrainsFullBacklog(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var g *TranscriptGrouper
		calls := 0
		event := func(i int) ServerEventUnion {
			kind := "session.input_transcript.delta"
			if i%2 == 1 {
				kind = "session.output_transcript.delta"
			}
			var e ServerEventUnion
			raw := fmt.Sprintf(`{"type":%q,"event_id":"panic_%d","delta":"word","start_ms":%d,"end_ms":%d}`, kind, i, i*5000, i*5000+100)
			if err := json.Unmarshal([]byte(raw), &e); err != nil {
				t.Fatal(err)
			}
			return e
		}
		var err error
		g, err = NewTranscriptGrouper(TranscriptGrouperOptions{OnUpdate: func(TranscriptUpdate) {
			calls++
			if calls != 1 {
				return
			}
			for i := 2; i < 2048; i++ {
				if pushErr := g.Push(event(i)); errors.Is(pushErr, ErrTranscriptBackpressure) {
					panic("callback failed after filling its backlog")
				} else if pushErr != nil {
					t.Fatal(pushErr)
				}
			}
			t.Fatal("callback did not reach backpressure")
		}})
		if err != nil {
			t.Fatal(err)
		}
		if err := g.Push(event(0)); err != nil {
			t.Fatal(err)
		}
		func() {
			defer func() {
				if recover() == nil {
					t.Error("expected callback panic")
				}
			}()
			_ = g.Push(event(1))
		}()
		defer g.Close()
		if calls <= 1 {
			t.Fatal("callback backlog did not drain before panic returned")
		}
		if err := g.Push(event(3000)); err != nil {
			t.Fatalf("event was not retryable after drain: %v", err)
		}
	})
}
