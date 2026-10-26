package live_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/packages/param"
)

type transcriptCollector struct {
	mu      sync.Mutex
	updates []live.TranscriptUpdate
}

func (c *transcriptCollector) push(u live.TranscriptUpdate) {
	c.mu.Lock()
	c.updates = append(c.updates, u)
	c.mu.Unlock()
}

func (c *transcriptCollector) snapshot() []live.TranscriptUpdate {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]live.TranscriptUpdate(nil), c.updates...)
}

func transcriptEvent(t testing.TB, id, speaker, text string, start, end int64) live.ServerEventUnion {
	t.Helper()
	kind := "session.input_transcript.delta"
	if speaker == "assistant" {
		kind = "session.output_transcript.delta"
	}
	raw, err := json.Marshal(map[string]any{"type": kind, "event_id": id, "delta": text, "start_ms": start, "end_ms": end})
	if err != nil {
		t.Fatal(err)
	}
	var e live.ServerEventUnion
	if err := json.Unmarshal(raw, &e); err != nil {
		t.Fatal(err)
	}
	return e
}

func pushTranscript(t *testing.T, g *live.TranscriptGrouper, id, speaker, text string, start, end int64) {
	t.Helper()
	if err := g.Push(transcriptEvent(t, id, speaker, text, start, end)); err != nil {
		t.Fatal(err)
	}
}

func checkTranscriptFinish(t *testing.T, c *transcriptCollector, want [][2]string) {
	t.Helper()
	var got [][2]string
	var previous string
	seen := map[string]bool{}
	latest := map[string]live.TranscriptSegment{}
	for _, u := range c.snapshot() {
		if u.Reason == "" {
			old, exists := latest[u.Segment.ID]
			if exists && !strings.HasPrefix(u.Segment.Text, old.Text) {
				t.Fatalf("revised published text: before=%q after=%q", old.Text, u.Segment.Text)
			}
			latest[u.Segment.ID] = u.Segment
			continue
		}
		if seen[u.Segment.ID] || u.Segment.PreviousID != previous || latest[u.Segment.ID] != u.Segment {
			t.Fatalf("unstable ID/previous/final snapshot: %#v latest=%#v previous=%s", u, latest, previous)
		}
		got = append(got, [2]string{u.Segment.Speaker, u.Segment.Text})
		seen[u.Segment.ID] = true
		previous = u.Segment.ID
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("segments = %#v want %#v; updates=%#v", got, want, c.snapshot())
	}
}

func TestLiveTranscriptTypedInputSnapshotsAndDedup(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		e := transcriptEvent(t, "raw1", "user", "Can you ", 0, 200)
		original := e.RawJSON()
		if err := g.Push(e); err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "raw2", "user", "run ls?", 200, 400)
		before := c.snapshot()[0]
		pushTranscript(t, g, "raw1", "user", "duplicate ignored", 400, 500)
		pushTranscript(t, g, "raw3", "assistant", "Sure.", 1100, 1200)
		g.Close()
		if e.RawJSON() != original || before.Segment.Text != "Can you " {
			t.Fatalf("raw event or earlier snapshot changed: %#v raw=%s", before, e.RawJSON())
		}
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Can you run ls?"}, {"assistant", "Sure."}})
		if err := g.Push(e); err == nil {
			t.Fatal("push after Close accepted")
		}
		n := len(c.snapshot())
		g.Close()
		time.Sleep(5 * time.Second)
		synctest.Wait()
		if len(c.snapshot()) != n {
			t.Fatal("timer/callback emitted after Close")
		}
	})
}

func TestLiveTranscriptConnectionEventTypes(t *testing.T) {
	for _, name := range []string{"live", "sideband", "fork"} {
		t.Run(name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				defer g.Close()
				push := func(raw string) error {
					switch name {
					case "sideband":
						var event live.ConnectServerEventUnion
						if err := json.Unmarshal([]byte(raw), &event); err != nil {
							t.Fatal(err)
						}
						return g.PushSideband(event)
					case "fork":
						var event live.ForkServerEventUnion
						if err := json.Unmarshal([]byte(raw), &event); err != nil {
							t.Fatal(err)
						}
						return g.PushFork(event)
					default:
						var event live.ServerEventUnion
						if err := json.Unmarshal([]byte(raw), &event); err != nil {
							t.Fatal(err)
						}
						return g.Push(event)
					}
				}
				if err := push(`{"type":"session.input_transcript.delta","event_id":"missing-fields"}`); err == nil {
					t.Fatal("accepted a malformed transcript event")
				}
				for _, raw := range []string{
					`{"type":"future.event"}`,
					`{"type":"session.input_transcript.delta","event_id":"user","delta":"Hello","start_ms":0,"end_ms":500}`,
					`{"type":"session.input_transcript.delta","event_id":"user","delta":"Hello","start_ms":0,"end_ms":500}`,
					`{"type":"session.output_transcript.delta","event_id":"assistant","delta":"The answer","start_ms":1000,"end_ms":1200}`,
					`{"type":"session.closed","event_id":"closed","reason":"remote_hangup","session":{},"usage":{}}`,
				} {
					if err := push(raw); err != nil {
						t.Fatal(err)
					}
				}
				checkTranscriptFinish(t, &c, [][2]string{{"user", "Hello"}, {"assistant", "The answer"}})
				for _, update := range c.snapshot() {
					if update.Reason != "" && update.Reason != "session_closed" && update.Reason != "speaker_change" {
						t.Fatalf("unexpected finalization: %#v", update)
					}
				}
				if err := push(`{"type":"future.event"}`); err == nil {
					t.Fatal("accepted an event after session.closed")
				}
			})
		})
	}
}

func TestLiveTranscriptOverlapAndBargeIn(t *testing.T) {
	for _, inputFirst := range []bool{true, false} {
		t.Run(fmt.Sprint(inputFirst), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				input := transcriptEvent(t, "u1", "user", "Question", 0, 200)
				output := transcriptEvent(t, "a1", "assistant", "Answer", 0, 200)
				if inputFirst {
					if err = g.Push(input); err != nil {
						t.Fatal(err)
					}
					if err = g.Push(output); err != nil {
						t.Fatal(err)
					}
				} else {
					if err = g.Push(output); err != nil {
						t.Fatal(err)
					}
					if err = g.Push(input); err != nil {
						t.Fatal(err)
					}
				}
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"user", "Question"}, {"assistant", "Answer"}})
			})
		})
		t.Run(fmt.Sprintf("barge-in-%t", inputFirst), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				pushTranscript(t, g, "a0", "assistant", "The answer is", 0, 200)
				time.Sleep(50 * time.Millisecond)
				synctest.Wait()
				input := transcriptEvent(t, "u1", "user", "Wait", 200, 400)
				output := transcriptEvent(t, "a1", "assistant", " forty-two.", 200, 400)
				if inputFirst {
					if err = g.Push(input); err != nil {
						t.Fatal(err)
					}
					if err = g.Push(output); err != nil {
						t.Fatal(err)
					}
				} else {
					if err = g.Push(output); err != nil {
						t.Fatal(err)
					}
					if err = g.Push(input); err != nil {
						t.Fatal(err)
					}
				}
				pushTranscript(t, g, "u2", "user", ", stop.", 400, 600)
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"assistant", "The answer is forty-two."}, {"user", "Wait, stop."}})
			})
		})
	}
}

func TestLiveTranscriptAcknowledgmentPolicy(t *testing.T) {
	tests := []struct {
		phrase     string
		suppressed bool
	}{
		{"mhm", true}, {"Mm-hmm.", true}, {"\ufeffYEAH\ufeff", true}, {"\u00a0uh\u2003huh!\u00a0", true},
		{"\u0085yeah\u0085", false}, {"\u001cyeah\u001c", false}, {"okay-ish", false}, {"thirteen", false},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%q", tc.phrase), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				pushTranscript(t, g, "u1", "user", "Tell me", 0, 200)
				pushTranscript(t, g, "a1", "assistant", tc.phrase, 100, 400)
				pushTranscript(t, g, "u2", "user", " more", 800, 1000)
				time.Sleep(5 * time.Second)
				synctest.Wait()
				g.Close()
				want := [][2]string{{"user", "Tell me"}, {"assistant", tc.phrase}, {"user", " more"}}
				if tc.suppressed {
					want = [][2]string{{"user", "Tell me more"}}
				}
				checkTranscriptFinish(t, &c, want)
			})
		})
	}
}

func TestLiveTranscriptExpiredAcknowledgmentKeepsVerbatimWordBoundary(t *testing.T) {
	for _, advanceTimer := range []bool{false, true} {
		t.Run(fmt.Sprintf("timer-%t", advanceTimer), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				pushTranscript(t, g, "u1", "user", "Tell me", 0, 200)
				pushTranscript(t, g, "a1", "assistant", "mhm", 100, 400)
				if advanceTimer {
					time.Sleep(3 * time.Second)
					synctest.Wait()
				}
				pushTranscript(t, g, "u2", "user", " mo", 3000, 3200)
				pushTranscript(t, g, "u3", "user", "re", 3200, 3300)
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"user", "Tell me more"}})
			})
		})
	}
}

func TestLiveTranscriptIncompleteAcknowledgmentIsNotDropped(t *testing.T) {
	t.Run("interleaved-word", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			var c transcriptCollector
			g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
			if err != nil {
				t.Fatal(err)
			}
			pushTranscript(t, g, "u1", "user", "Tell me", 0, 200)
			pushTranscript(t, g, "a1", "assistant", "y", 200, 300)
			pushTranscript(t, g, "u2", "user", " more", 330, 500)
			pushTranscript(t, g, "a2", "assistant", "ellow", 550, 650)
			g.Close()
			checkTranscriptFinish(t, &c, [][2]string{{"user", "Tell me more"}, {"assistant", "yellow"}})
		})
	})
	t.Run("prefix-at-isolation-deadline", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			var c transcriptCollector
			g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
			if err != nil {
				t.Fatal(err)
			}
			pushTranscript(t, g, "u1", "user", "Tell me", 0, 200)
			pushTranscript(t, g, "a1", "assistant", "y", 200, 300)
			time.Sleep(3 * time.Second)
			synctest.Wait()
			g.Close()
			checkTranscriptFinish(t, &c, [][2]string{{"user", "Tell me"}, {"assistant", "y"}})
		})
	})
}

func TestLiveTranscriptCompletePhrasePrefixKeepsInterleavedSuffix(t *testing.T) {
	for _, suffix := range []string{"-ish", " dokey"} {
		t.Run(suffix, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				pushTranscript(t, g, "u1", "user", "Tell", 0, 200)
				pushTranscript(t, g, "a1", "assistant", "okay", 100, 300)
				pushTranscript(t, g, "u2", "user", " me", 300, 500)
				pushTranscript(t, g, "a2", "assistant", suffix, 550, 650)
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"user", "Tell me"}, {"assistant", "okay" + suffix}})
			})
		})
	}
}

func TestLiveTranscriptResumedWordUsesOnlySuppliedBoundary(t *testing.T) {
	for _, phrase := range []string{"mhm", "actually"} {
		t.Run(phrase, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				pushTranscript(t, g, "u1", "user", "Tell", 0, 200)
				pushTranscript(t, g, "a1", "assistant", phrase, 100, 300)
				pushTranscript(t, g, "u2", "user", " mo", 310, 500)
				pushTranscript(t, g, "u3", "user", "re", 550, 650)
				g.Close()
				want := [][2]string{{"user", "Tell more"}, {"assistant", phrase}}
				if phrase == "mhm" {
					want = [][2]string{{"user", "Tell more"}}
				}
				checkTranscriptFinish(t, &c, want)
			})
		})
	}
}

func TestLiveTranscriptDiscriminatorChangesDoNotConsumeIDOrSpeaker(t *testing.T) {
	for _, original := range []string{"user", "assistant"} {
		t.Run(original, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				event := transcriptEvent(t, "usable", original, "Hello", 0, 200)
				changed := event
				changed.Type = "session.input_transcript.delta"
				if original == "user" {
					changed.Type = "session.output_transcript.delta"
				}
				if err = g.Push(changed); err == nil {
					t.Fatal("rewritten discriminator changed the speaker")
				}
				if err = g.Push(event); err != nil {
					t.Fatal(err)
				}
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{original, "Hello"}})
			})
		})
	}
}

func TestLiveTranscriptInitialOverlapAcrossUnequalAndSplitRanges(t *testing.T) {
	for _, fragmented := range []bool{false, true} {
		for _, userFirst := range []bool{false, true} {
			t.Run(fmt.Sprintf("split=%t/userFirst=%t", fragmented, userFirst), func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					var c transcriptCollector
					g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
					if err != nil {
						t.Fatal(err)
					}
					user := func() { pushTranscript(t, g, "u", "user", "Hello", 0, 500) }
					assistant := func() {
						if fragmented {
							pushTranscript(t, g, "a1", "assistant", "o", 0, 50)
							pushTranscript(t, g, "a2", "assistant", "kay", 50, 100)
						} else {
							pushTranscript(t, g, "a", "assistant", "okay", 0, 100)
						}
					}
					if userFirst {
						user()
						assistant()
					} else {
						assistant()
						user()
					}
					g.Close()
					checkTranscriptFinish(t, &c, [][2]string{{"user", "Hello"}})
				})
			})
		}
	}
}

func TestLiveTranscriptVerbatimFragmentsAcrossInterjection(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "u1", "user", "hel", 0, 100)
		pushTranscript(t, g, "a", "assistant", "mhm", 30, 100)
		pushTranscript(t, g, "u2", "user", "lo", 100, 200)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "hello"}})
	})
}

func TestLiveTranscriptDisjointAcknowledgmentAlwaysPreserved(t *testing.T) {
	for _, userFirst := range []bool{false, true} {
		t.Run(fmt.Sprint(userFirst), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				user := func() { pushTranscript(t, g, "u1", "user", "Tell", 100, 550) }
				assistant := func() { pushTranscript(t, g, "a", "assistant", "okay", 0, 100) }
				if userFirst {
					user()
					assistant()
				} else {
					assistant()
					user()
				}
				pushTranscript(t, g, "u2", "user", " me", 550, 700)
				g.Close()
				var assistantText string
				for _, update := range c.snapshot() {
					if update.Reason != "" && update.Segment.Speaker == "assistant" {
						assistantText += update.Segment.Text
					}
				}
				if assistantText != "okay" {
					t.Fatalf("assistant text = %q; updates = %#v", assistantText, c.snapshot())
				}
			})
		})
	}
}

func TestLiveTranscriptEqualEndContinuationAtTerminal(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "u1", "user", "Tell", 0, 100)
		pushTranscript(t, g, "a", "assistant", "okay", 0, 200)
		pushTranscript(t, g, "u2", "user", " me", 100, 200)
		var terminal live.ServerEventUnion
		if err = json.Unmarshal([]byte(`{"type":"session.closed","event_id":"closed","reason":"remote_hangup","session":{},"usage":{}}`), &terminal); err != nil {
			t.Fatal(err)
		}
		if err = g.Push(terminal); err != nil {
			t.Fatal(err)
		}
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Tell me"}})
	})
}

func TestLiveTranscriptGappedInitialAcknowledgment(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "a1", "assistant", "o", 0, 40)
		pushTranscript(t, g, "a2", "assistant", "kay", 50, 100)
		pushTranscript(t, g, "u", "user", "Hello", 0, 500)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Hello"}})
	})
}

func TestLiveTranscriptInitialAssistantPrefixBeforeLaterOverlap(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "a1", "assistant", "Earlier", 0, 100)
		pushTranscript(t, g, "a2", "assistant", "Later", 2000, 2100)
		pushTranscript(t, g, "u1", "user", "Question", 1900, 2200)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"assistant", "Earlier"}, {"user", "Question"}, {"assistant", "Later"}})
	})
}

func TestLiveTranscriptInitialBatchRespectsConfiguredSilence(t *testing.T) {
	for _, silence := range []time.Duration{2 * time.Second, 300 * time.Millisecond} {
		t.Run(silence.String(), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{
					OnUpdate: c.push, AssistantSilence: param.NewOpt(silence),
				})
				if err != nil {
					t.Fatal(err)
				}
				pushTranscript(t, g, "a1", "assistant", "First", 0, 100)
				start := 100 + silence.Milliseconds()
				pushTranscript(t, g, "a2", "assistant", "Second", start, start+100)
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"assistant", "First"}, {"assistant", "Second"}})
			})
		})
	}
}

func TestLiveTranscriptNestedAudioDoesNotMoveAnchor(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "a1", "assistant", "First", 0, 1000)
		time.Sleep(time.Second)
		synctest.Wait()
		pushTranscript(t, g, "a2", "assistant", " nested", 500, 600)
		time.Sleep(1100 * time.Millisecond)
		synctest.Wait()
		checkTranscriptFinish(t, &c, [][2]string{{"assistant", "First nested"}})
		if updates := c.snapshot(); updates[len(updates)-1].Reason != "inactivity" {
			t.Fatalf("last update %#v", updates[len(updates)-1])
		}
		g.Close()
	})
}

func TestLiveTranscriptEarlierUserAudioDoesNotSuppressLaterAck(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "u1", "user", "Tell", 0, 250)
		pushTranscript(t, g, "a1", "assistant", "okay", 200, 400)
		pushTranscript(t, g, "u2", "user", " me", 100, 150)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Tell me"}, {"assistant", "okay"}})
	})
}

func TestLiveTranscriptInitialAckWithOnlySuffixOverlapping(t *testing.T) {
	for _, split := range []bool{false, true} {
		t.Run(fmt.Sprint(split), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				if split {
					pushTranscript(t, g, "a1", "assistant", "o", 0, 40)
					pushTranscript(t, g, "a2", "assistant", "kay", 50, 100)
				} else {
					pushTranscript(t, g, "a1", "assistant", "okay", 0, 100)
				}
				pushTranscript(t, g, "u1", "user", "Hello", 50, 500)
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"user", "Hello"}})
			})
		})
	}
}

func TestLiveTranscriptLongAlternatingDisjointRanges(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const count = 12000
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < count; i++ {
			start := int64(4 * i)
			pushTranscript(t, g, fmt.Sprint("u", i), "user", "x", start, start+1)
			pushTranscript(t, g, fmt.Sprint("a", i), "assistant", "!", start+2, start+3)
		}
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", strings.Repeat("x", count)}, {"assistant", strings.Repeat("!", count)}})
	})
}

func TestLiveTranscriptIgnoredEmptyDeltasDoNotConsumeIDs(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 4096; i++ {
			pushTranscript(t, g, fmt.Sprint("noop-", i), "user", "", 0, 0)
		}
		if got := c.snapshot(); len(got) != 0 {
			t.Fatalf("ignored events emitted %#v", got)
		}
		pushTranscript(t, g, "noop-0", "user", "Hello", 0, 200)
		pushTranscript(t, g, "noop-0", "user", "Hello", 0, 200)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Hello"}})
	})
}

func TestLiveTranscriptAcknowledgmentOverlapsAfterContextPoint(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "u1", "user", "Tell", 0, 50)
		pushTranscript(t, g, "u2", "user", " me", 100, 100)
		pushTranscript(t, g, "u3", "user", " more", 120, 200)
		pushTranscript(t, g, "a1", "assistant", "mhm", 60, 150)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Tell me more"}})
	})
}

func TestLiveTranscriptZeroLengthContextAtDifferentStartIsPreserved(t *testing.T) {
	for _, userFirst := range []bool{false, true} {
		t.Run(fmt.Sprint(userFirst), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				user := func() { pushTranscript(t, g, "u1", "user", "Hello", 0, 500) }
				assistant := func() { pushTranscript(t, g, "a1", "assistant", "okay", 100, 100) }
				if userFirst {
					user()
					assistant()
				} else {
					assistant()
					user()
				}
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"user", "Hello"}, {"assistant", "okay"}})
			})
		})
	}
}

func TestLiveTranscriptAckSuffixAfterSeparationBeforeIsolation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "u1", "user", "Tell", 0, 200)
		pushTranscript(t, g, "a1", "assistant", "okay", 100, 300)
		pushTranscript(t, g, "u2", "user", " me", 300, 500)
		pushTranscript(t, g, "a2", "assistant", "-ish", 1100, 1200)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Tell me"}, {"assistant", "okay-ish"}})
	})
}

func TestLiveTranscriptEarlierDisjointInputArrivesDuringSettle(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "a1", "assistant", "Answer", 100, 200)
		pushTranscript(t, g, "u1", "user", "Question", 0, 50)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Question"}, {"assistant", "Answer"}})
	})
}

func TestLiveTranscriptZeroContextMatchesConstituentStart(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "u1", "user", "Tell", 0, 100)
		pushTranscript(t, g, "u2", "user", " me", 100, 200)
		pushTranscript(t, g, "a1", "assistant", "okay", 100, 100)
		pushTranscript(t, g, "u3", "user", " more", 200, 300)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Tell me more"}})
	})
}

func TestLiveTranscriptBufferedAssistantRespectsSilence(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "u1", "user", "Question", 0, 6000)
		pushTranscript(t, g, "a1", "assistant", "One", 100, 200)
		pushTranscript(t, g, "a2", "assistant", "Two", 5000, 5100)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Question"}, {"assistant", "One"}, {"assistant", "Two"}})
	})
}

func TestLiveTranscriptInitialAckBetweenSeparationAndSilence(t *testing.T) {
	for _, split := range []bool{false, true} {
		t.Run(fmt.Sprint(split), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				if split {
					pushTranscript(t, g, "a1", "assistant", "o", 0, 40)
					pushTranscript(t, g, "a2", "assistant", "kay", 600, 650)
				} else {
					pushTranscript(t, g, "a1", "assistant", "okay", 0, 650)
				}
				pushTranscript(t, g, "u1", "user", "Hello", 600, 1000)
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"user", "Hello"}})
			})
		})
	}
}

func TestLiveTranscriptBufferedAckPrefixAcrossSeparation(t *testing.T) {
	for _, split := range []bool{false, true} {
		t.Run(fmt.Sprint(split), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				pushTranscript(t, g, "u1", "user", "Tell", 0, 200)
				if split {
					pushTranscript(t, g, "a1", "assistant", "o", 100, 300)
					pushTranscript(t, g, "a2", "assistant", "kay", 800, 900)
				} else {
					pushTranscript(t, g, "a1", "assistant", "okay", 100, 900)
				}
				pushTranscript(t, g, "u2", "user", " me", 850, 1000)
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"user", "Tell me"}})
			})
		})
	}
}

func TestLiveTranscriptRecentAssistantTailSurvivesContinuation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "a0", "assistant", "First", 0, 100)
		time.Sleep(50 * time.Millisecond)
		synctest.Wait()
		pushTranscript(t, g, "u1", "user", "Tell", 200, 400)
		time.Sleep(50 * time.Millisecond)
		synctest.Wait()
		pushTranscript(t, g, "a1", "assistant", "okay", 150, 250)
		pushTranscript(t, g, "u2", "user", " me", 400, 600)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"assistant", "First"}, {"user", "Tell me"}, {"assistant", "okay"}})
	})
}

func TestLiveTranscriptTailOverlapsPriorAssistantAfterUserStart(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "a0", "assistant", "First", 0, 300)
		time.Sleep(50 * time.Millisecond)
		synctest.Wait()
		pushTranscript(t, g, "u1", "user", "Tell", 200, 500)
		time.Sleep(50 * time.Millisecond)
		synctest.Wait()
		pushTranscript(t, g, "a1", "assistant", "okay", 250, 350)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"assistant", "First"}, {"user", "Tell"}, {"assistant", "okay"}})
	})
}

func TestLiveTranscriptInitialNonAckRetainsChronologicalOrder(t *testing.T) {
	for _, start := range []int64{100, 400, 600} {
		for _, settled := range []bool{false, true} {
			t.Run(fmt.Sprint(start, "-", settled), func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					var c transcriptCollector
					g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
					if err != nil {
						t.Fatal(err)
					}
					pushTranscript(t, g, "a0", "assistant", "Hello", 0, 100)
					if settled {
						time.Sleep(50 * time.Millisecond)
						synctest.Wait()
					}
					pushTranscript(t, g, "a1", "assistant", " world", start, start+50)
					pushTranscript(t, g, "u1", "user", "Question", start, start+400)
					g.Close()
					checkTranscriptFinish(t, &c, [][2]string{{"assistant", "Hello world"}, {"user", "Question"}})
				})
			})
		}
	}
}

func TestLiveTranscriptInterposedDisjointUserSpeech(t *testing.T) {
	for _, inSourceOrder := range []bool{false, true} {
		t.Run(fmt.Sprint(inSourceOrder), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				pushTranscript(t, g, "a0", "assistant", "First", 0, 100)
				if inSourceOrder {
					pushTranscript(t, g, "u1", "user", "Question", 500, 600)
					pushTranscript(t, g, "a1", "assistant", "Later", 1000, 1100)
				} else {
					pushTranscript(t, g, "a1", "assistant", "Later", 1000, 1100)
					pushTranscript(t, g, "u1", "user", "Question", 500, 600)
				}
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"assistant", "First"}, {"user", "Question"}, {"assistant", "Later"}})
			})
		})
	}
}

func TestLiveTranscriptLargeInitialAssistantBurst(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const count = 12000
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < count; i++ {
			start := int64(i * 2)
			pushTranscript(t, g, fmt.Sprint("a", i), "assistant", "x", start, start+1)
		}
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"assistant", strings.Repeat("x", count)}})
	})
}

func TestLiveTranscriptFlushesPendingBeforeTimestampReset(t *testing.T) {
	for _, interval := range [][2]int64{{1100, 1200}, {950, 1150}, {900, 950}} {
		t.Run(fmt.Sprint(interval), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				pushTranscript(t, g, "u1", "user", "Old question", 1000, 1100)
				time.Sleep(50 * time.Millisecond)
				synctest.Wait()
				pushTranscript(t, g, "a1", "assistant", "Old answer", interval[0], interval[1])
				pushTranscript(t, g, "u2", "user", "New question", 0, 100)
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"user", "Old question"}, {"assistant", "Old answer"}, {"user", "New question"}})
			})
		})
	}
}

func TestLiveTranscriptInitialAcknowledgmentWaitsForSourceContext(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "a1", "assistant", "okay", 0, 200)
		time.Sleep(5 * time.Second)
		synctest.Wait()
		if len(c.snapshot()) != 0 {
			t.Fatal("published an acknowledgment before source context arrived")
		}
		pushTranscript(t, g, "u1", "user", "Hello", 0, 500)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Hello"}})
	})
}

func TestLiveTranscriptResetKeepsNewEpochPendingFragments(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "old", "user", "Old question", 1000, 1100)
		time.Sleep(50 * time.Millisecond)
		synctest.Wait()
		pushTranscript(t, g, "a1", "assistant", "okay", 0, 100)
		pushTranscript(t, g, "u1", "user", "New question", 0, 500)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Old question"}, {"user", "New question"}})
	})
}

func TestLiveTranscriptAdjacentAssistantTailIsPreserved(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "a0", "assistant", "First", 0, 300)
		time.Sleep(50 * time.Millisecond)
		synctest.Wait()
		pushTranscript(t, g, "u1", "user", "Tell", 200, 500)
		time.Sleep(50 * time.Millisecond)
		synctest.Wait()
		pushTranscript(t, g, "a1", "assistant", "okay", 300, 400)
		pushTranscript(t, g, "u2", "user", " me", 500, 700)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"assistant", "First"}, {"user", "Tell me"}, {"assistant", "okay"}})
	})
}

func TestLiveTranscriptDisplayTimerPreservesSettleWindow(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "a0", "assistant", "The answer is", 0, 100)
		time.Sleep(50 * time.Millisecond)
		synctest.Wait()
		pushTranscript(t, g, "a1", "assistant", ":", 100, 120)
		time.Sleep(40 * time.Millisecond)
		synctest.Wait()
		pushTranscript(t, g, "u1", "user", "Wait", 200, 400)
		time.Sleep(20 * time.Millisecond)
		synctest.Wait()
		pushTranscript(t, g, "a2", "assistant", " forty-two", 200, 400)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"assistant", "The answer is: forty-two"}, {"user", "Wait"}})
	})
}

func TestLiveTranscriptDuplicateTerminalDoesNotCloseOrFlush(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "duplicate", "user", "First", 0, 200)
		before := c.snapshot()
		var terminal live.ServerEventUnion
		if err = json.Unmarshal([]byte(`{"type":"session.closed","event_id":"duplicate","reason":"remote_hangup","session":{},"usage":{}}`), &terminal); err != nil {
			t.Fatal(err)
		}
		if err = g.Push(terminal); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(before, c.snapshot()) {
			t.Fatalf("duplicate terminal flushed pending text: %#v", c.snapshot())
		}
		pushTranscript(t, g, "next", "user", " next", 200, 400)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "First next"}})
	})
}

func TestLiveTranscriptSpeechInsideUserGapIsPreserved(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "u1", "user", "Tell", 0, 100)
		pushTranscript(t, g, "a1", "assistant", "okay", 200, 300)
		pushTranscript(t, g, "u2", "user", " me", 400, 500)
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Tell me"}, {"assistant", "okay"}})
	})
}

func TestLiveTranscriptDeadlinesIndependentOfSettleTimer(t *testing.T) {
	for _, timerFirst := range []bool{false, true} {
		t.Run(fmt.Sprintf("inactivity/timer=%t", timerFirst), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				pushTranscript(t, g, "a0", "assistant", "First", 0, 100)
				time.Sleep(50 * time.Millisecond)
				synctest.Wait()
				pushTranscript(t, g, "u1", "user", "Question", 1500, 2500)
				if timerFirst {
					time.Sleep(50 * time.Millisecond)
					synctest.Wait()
				}
				pushTranscript(t, g, "a1", "assistant", "Answer", 2200, 2300)
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"assistant", "First"}, {"user", "Question"}, {"assistant", "Answer"}})
			})
		})
		t.Run(fmt.Sprintf("promotion/timer=%t", timerFirst), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				pushTranscript(t, g, "u1", "user", "Start", 0, 200)
				pushTranscript(t, g, "a0", "assistant", "Actually", 100, 300)
				if timerFirst {
					time.Sleep(500 * time.Millisecond)
					synctest.Wait()
				}
				pushTranscript(t, g, "a1", "assistant", " more", 800, 900)
				pushTranscript(t, g, "u2", "user", "Question", 850, 1000)
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"user", "Start"}, {"assistant", "Actually more"}, {"user", "Question"}})
			})
		})
	}
}

func TestLiveTranscriptSuppressionZeroAndLargeAdditionalPhrase(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		for _, disabled := range []bool{true, false} {
			var c transcriptCollector
			phrase := strings.Repeat("arbitrary", 1024)
			phrases := []string{phrase}
			opts := live.TranscriptGrouperOptions{OnUpdate: c.push, AdditionalAcknowledgments: phrases}
			if disabled {
				opts.BackchannelMaxDuration = param.NewOpt(time.Duration(0))
			}
			g, err := live.NewTranscriptGrouper(opts)
			if err != nil {
				t.Fatal(err)
			}
			phrases[0] = "changed"
			pushTranscript(t, g, "u1", "user", "Tell me", 0, 200)
			// Split beyond any built-in acknowledgment length without truncation.
			for i := 0; i < len(phrase); i += 32 {
				pushTranscript(t, g, fmt.Sprint("ack", i), "assistant", phrase[i:min(i+32, len(phrase))], 100, 400)
			}
			pushTranscript(t, g, "u2", "user", " more", 800, 1000)
			time.Sleep(5 * time.Second)
			synctest.Wait()
			g.Close()
			want := [][2]string{{"user", "Tell me more"}}
			if disabled {
				want = [][2]string{{"user", "Tell me"}, {"assistant", phrase}, {"user", " more"}}
			}
			checkTranscriptFinish(t, &c, want)
		}
	})
}

func TestLiveTranscriptTimersResetAndTerminal(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{
			OnUpdate: c.push, AssistantSilence: param.NewOpt(100 * time.Microsecond),
			MinTurnSeparation: param.NewOpt(time.Nanosecond),
		})
		if err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "a1", "assistant", "One", 1000, 1200)
		time.Sleep(51 * time.Millisecond)
		synctest.Wait()
		pushTranscript(t, g, "u1", "user", "next", 2000, 2200)
		pushTranscript(t, g, "u2", "user", "reset", 10, 20)
		var terminal live.ServerEventUnion
		if err = json.Unmarshal([]byte(`{"type":"session.closed","event_id":"terminal","reason":"remote_hangup","session":{"id":"ls_synthetic"},"usage":{}}`), &terminal); err != nil {
			t.Fatal(err)
		}
		if err = g.Push(terminal); err != nil {
			t.Fatal(err)
		}
		checkTranscriptFinish(t, &c, [][2]string{{"assistant", "One"}, {"user", "next"}, {"user", "reset"}})
		var reasons []string
		for _, u := range c.snapshot() {
			if u.Reason != "" {
				reasons = append(reasons, u.Reason)
			}
		}
		if !reflect.DeepEqual(reasons, []string{"inactivity", "timestamp_reset", "session_closed"}) {
			t.Fatalf("reasons=%#v", reasons)
		}
		n := len(c.snapshot())
		time.Sleep(time.Hour)
		synctest.Wait()
		if len(c.snapshot()) != n {
			t.Fatal("stale timer after terminal")
		}
	})
}

func TestLiveTranscriptInvalidIgnoredAndReentrant(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		if _, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{AssistantSilence: param.NewOpt(-time.Nanosecond)}); err == nil {
			t.Fatal("negative duration accepted")
		}
		var c transcriptCollector
		var g *live.TranscriptGrouper
		var err error
		g, err = live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: func(u live.TranscriptUpdate) {
			c.push(u)
			if u.Reason == "" {
				g.Close()
			}
		}})
		if err != nil {
			t.Fatal(err)
		}
		for _, event := range []live.ServerEventUnion{
			transcriptEvent(t, "", "user", "missing id", 0, 20),
			transcriptEvent(t, "bad_time", "assistant", "invalid", -1, 20),
			transcriptEvent(t, "reverse", "user", "invalid", 100, 20),
		} {
			if err = g.Push(event); err == nil {
				t.Fatal("invalid native event accepted")
			}
		}
		var other live.ServerEventUnion
		if err = json.Unmarshal([]byte(`{"type":"session.usage.updated","event_id":"ignored"}`), &other); err != nil {
			t.Fatal(err)
		}
		if err = g.Push(other); err != nil {
			t.Fatal(err)
		}
		pushTranscript(t, g, "valid", "user", "Hello", 0, 20)
		time.Sleep(50 * time.Millisecond)
		synctest.Wait()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Hello"}})
	})
}

func TestLiveTranscriptMalformedWireDoesNotConsumeID(t *testing.T) {
	for _, kind := range []string{"session.input_transcript.delta", "session.output_transcript.delta"} {
		for _, change := range []struct {
			field string
			value any
		}{
			{"delta", nil}, {"delta", map[string]string{"unexpected": "value"}}, {"start_ms", nil}, {"start_ms", "now"}, {"end_ms", nil}, {"end_ms", "later"},
		} {
			t.Run(kind+"/"+fmt.Sprint(change.field, "-", change.value), func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				defer g.Close()
				obj := map[string]any{"type": kind, "event_id": "reusable", "delta": "Hello", "start_ms": 0, "end_ms": 20}
				obj[change.field] = change.value
				raw, err := json.Marshal(obj)
				if err != nil {
					t.Fatal(err)
				}
				var event live.ServerEventUnion
				if err = json.Unmarshal(raw, &event); err != nil {
					t.Fatal(err)
				}
				if err = g.Push(event); err == nil {
					t.Fatalf("malformed native field accepted: %s", raw)
				}
				speaker := "user"
				if kind == "session.output_transcript.delta" {
					speaker = "assistant"
				}
				pushTranscript(t, g, "reusable", speaker, "Hello", 0, 20)
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{speaker, "Hello"}})
			})
		}
	}
}

func TestLiveTranscriptRejectsFalseTerminalAndRetainsPending(t *testing.T) {
	for _, tc := range []struct {
		name, raw  string
		changeType bool
	}{
		{"type-only", `{"type":"session.closed"}`, false},
		{"no-event-id", `{"type":"session.closed","reason":"remote_hangup","session":{},"usage":{}}`, false},
		{"no-reason", `{"type":"session.closed","event_id":"bad","session":{},"usage":{}}`, false},
		{"empty-reason", `{"type":"session.closed","event_id":"bad","reason":"","session":{},"usage":{}}`, false},
		{"no-session", `{"type":"session.closed","event_id":"bad","reason":"remote_hangup","usage":{}}`, false},
		{"null-usage", `{"type":"session.closed","event_id":"bad","reason":"remote_hangup","session":{},"usage":null}`, false},
		{"rewritten-discriminator", `{"type":"session.input_transcript.delta","event_id":"forged","delta":"not a terminal","start_ms":0,"end_ms":20,"reason":"remote_hangup","session":{},"usage":{}}`, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				defer g.Close()
				pushTranscript(t, g, "first", "user", "Keep ", 0, 20)
				var bad, terminal live.ServerEventUnion
				if err = json.Unmarshal([]byte(tc.raw), &bad); err != nil {
					t.Fatal(err)
				}
				if tc.changeType {
					bad.Type = "session.closed"
				}
				if err = g.Push(bad); err == nil {
					t.Fatalf("nonterminal stopped grouping: %s", tc.raw)
				}
				if len(c.snapshot()) != 0 {
					t.Fatal("malformed terminal flushed pending text")
				}
				pushTranscript(t, g, "second", "user", "going.", 20, 40)
				if err = json.Unmarshal([]byte(`{"type":"session.closed","event_id":"closed","reason":"future_reason","session":{},"usage":{}}`), &terminal); err != nil {
					t.Fatal(err)
				}
				if err = g.Push(terminal); err != nil {
					t.Fatal(err)
				}
				checkTranscriptFinish(t, &c, [][2]string{{"user", "Keep going."}})
				updates := c.snapshot()
				if updates[len(updates)-1].Reason != "session_closed" {
					t.Fatalf("real terminal not delivered: %#v", updates)
				}
				time.Sleep(time.Hour)
				synctest.Wait()
				if len(c.snapshot()) != len(updates) {
					t.Fatal("stale timers delivered after the real terminal")
				}
			})
		})
	}
}

func TestLiveTranscriptLongDeferredAcknowledgment(t *testing.T) {
	for _, prefix := range []string{"okay", ""} {
		t.Run(prefix, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				if prefix != "" {
					pushTranscript(t, g, "prefix", "assistant", prefix, 0, 100)
				}
				for i := 0; i < 12000; i++ {
					pushTranscript(t, g, fmt.Sprint(i), "assistant", " ", 0, 100)
				}
				time.Sleep(5 * time.Second)
				synctest.Wait()
				if len(c.snapshot()) != 0 {
					t.Fatal("published before source context")
				}
				pushTranscript(t, g, "continuation", "assistant", "more", 100, 200)
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"assistant", prefix + strings.Repeat(" ", 12000) + "more"}})
			})
		})
	}
}

func BenchmarkLiveTranscriptDeferredAcknowledgment(b *testing.B) {
	for _, count := range []int{1000, 10000} {
		b.Run(fmt.Sprint(count), func(b *testing.B) {
			events := make([]live.ServerEventUnion, count)
			for i := range events {
				events[i] = transcriptEvent(b, fmt.Sprint(i), "assistant", " ", 0, 100)
			}
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{})
				if err != nil {
					b.Fatal(err)
				}
				for _, event := range events {
					if err := g.Push(event); err != nil {
						b.Fatal(err)
					}
				}
				g.Close()
			}
		})
	}
}

func TestLiveTranscriptAcknowledgmentAfterAssistantTimestampReset(t *testing.T) {
	for _, tc := range []struct {
		name      string
		userFirst bool
		settle    time.Duration
	}{
		{"user-pending", true, 0},
		{"user-committed", true, 50 * time.Millisecond},
		{"user-delayed", true, 5 * time.Second},
		{"assistant-first", false, 5 * time.Second},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				defer g.Close()
				pushTranscript(t, g, "old", "assistant", "Old answer", 1000, 1100)
				time.Sleep(3 * time.Second)
				synctest.Wait()
				previous := len(c.snapshot())
				if tc.userFirst {
					pushTranscript(t, g, "user", "user", "Hello", 0, 500)
					time.Sleep(tc.settle)
					synctest.Wait()
				}
				pushTranscript(t, g, "ack", "assistant", "okay", 0, 100)
				if !tc.userFirst {
					time.Sleep(tc.settle)
					synctest.Wait()
					if len(c.snapshot()) != previous {
						t.Fatal("published acknowledgment after reset before source context")
					}
					pushTranscript(t, g, "user", "user", "Hello", 0, 500)
				}
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"assistant", "Old answer"}, {"user", "Hello"}})
			})
		})
	}
}

func TestLiveTranscriptLaterAcknowledgmentWaitsForSourceContext(t *testing.T) {
	for _, userFirst := range []bool{true, false} {
		t.Run(fmt.Sprint(userFirst), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: c.push})
				if err != nil {
					t.Fatal(err)
				}
				pushTranscript(t, g, "old", "assistant", "Old answer", 0, 100)
				time.Sleep(3 * time.Second)
				synctest.Wait()
				previous := len(c.snapshot())
				if userFirst {
					pushTranscript(t, g, "u1", "user", "Hello", 5000, 5500)
				}
				pushTranscript(t, g, "a1", "assistant", "okay", 5000, 5100)
				if !userFirst {
					time.Sleep(5 * time.Second)
					synctest.Wait()
					if len(c.snapshot()) != previous {
						t.Fatal("published a later acknowledgment before source context")
					}
					pushTranscript(t, g, "u1", "user", "Hello", 5000, 5500)
				}
				g.Close()
				checkTranscriptFinish(t, &c, [][2]string{{"assistant", "Old answer"}, {"user", "Hello"}})
			})
		})
	}
}

func TestLiveTranscriptReentrantCloseDrainsAfterCallbackPanic(t *testing.T) {
	for _, panicAgain := range []bool{false, true} {
		t.Run(fmt.Sprint(panicAgain), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				var g *live.TranscriptGrouper
				calls := 0
				var err error
				g, err = live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: func(u live.TranscriptUpdate) {
					calls++
					c.push(u)
					if calls == 1 {
						g.Close()
						panic("first callback panic")
					}
					if panicAgain {
						panic("later callback panic")
					}
				}})
				if err != nil {
					t.Fatal(err)
				}
				defer g.Close()
				pushTranscript(t, g, "user", "user", "Hello", 0, 100)
				func() {
					defer func() {
						if value := recover(); value != "first callback panic" {
							t.Errorf("panic = %v, want first callback panic", value)
						}
					}()
					pushTranscript(t, g, "assistant", "assistant", "Answer", 5000, 5100)
				}()
				checkTranscriptFinish(t, &c, [][2]string{{"user", "Hello"}, {"assistant", "Answer"}})
			})
		})
	}
}

func TestLiveTranscriptConcurrentPushDrainsAfterCallbackPanic(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c transcriptCollector
		entered, release := make(chan struct{}), make(chan struct{})
		panicked := make(chan any, 1)
		first := true
		g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: func(u live.TranscriptUpdate) {
			c.push(u)
			if first {
				first = false
				close(entered)
				<-release
				panic("first callback panic")
			}
		}})
		if err != nil {
			t.Fatal(err)
		}
		defer g.Close()
		pushTranscript(t, g, "first", "user", "Hello", 0, 100)
		go func() {
			defer func() { panicked <- recover() }()
			pushTranscript(t, g, "second", "assistant", "Answer", 5000, 5100)
		}()
		<-entered
		err = g.Push(transcriptEvent(t, "third", "user", "Next", 10000, 10100))
		close(release)
		if err != nil {
			t.Fatal(err)
		}
		if value := <-panicked; value != "first callback panic" {
			t.Fatalf("panic = %v, want first callback panic", value)
		}
		found := false
		for _, update := range c.snapshot() {
			if update.Segment.Text == "Answer" {
				found = true
			}
		}
		if !found {
			t.Fatal("accepted concurrent update was not delivered before panic returned")
		}
		g.Close()
		checkTranscriptFinish(t, &c, [][2]string{{"user", "Hello"}, {"assistant", "Answer"}, {"user", "Next"}})
	})
}

func TestLiveTranscriptSlowCallbackBackpressure(t *testing.T) {
	for _, closeWhileBlocked := range []bool{false, true} {
		t.Run(fmt.Sprintf("closeWhileBlocked=%v", closeWhileBlocked), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var c transcriptCollector
				entered, release := make(chan struct{}), make(chan struct{})
				first := true
				g, err := live.NewTranscriptGrouper(live.TranscriptGrouperOptions{OnUpdate: func(u live.TranscriptUpdate) {
					if first {
						first = false
						close(entered)
						<-release
					}
					c.push(u)
				}})
				if err != nil {
					t.Fatal(err)
				}
				pushTranscript(t, g, "initial", "user", "initial", 0, 100)
				<-entered
				want := [][2]string{{"user", "initial"}}
				var rejected live.ServerEventUnion
				for i := 1; i <= 2000; i++ {
					speaker := "assistant"
					if i%2 == 0 {
						speaker = "user"
					}
					word := fmt.Sprintf("word_%d", i)
					event := transcriptEvent(t, word, speaker, word, int64(i*5000), int64(i*5000+100))
					if err := g.Push(event); err != nil {
						if !errors.Is(err, live.ErrTranscriptBackpressure) {
							t.Fatal(err)
						}
						rejected = event
						break
					}
					want = append(want, [2]string{speaker, word})
				}
				if rejected.Type == "" {
					close(release)
					g.Close()
					t.Fatal("slow callback accepted an unbounded update backlog")
				}
				// A duplicate and unrelated event consume no queue space.
				if err := g.Push(transcriptEvent(t, "initial", "user", "initial", 0, 100)); err != nil {
					t.Fatal(err)
				}
				if err := g.Push(live.ServerEventUnion{}); err != nil {
					t.Fatal(err)
				}
				if closeWhileBlocked {
					g.Close()
				}
				close(release)
				synctest.Wait()
				if !closeWhileBlocked {
					if err := g.Push(rejected); err != nil {
						t.Fatal(err)
					}
					if rejected.Type == "session.input_transcript.delta" {
						want = append(want, [2]string{"user", rejected.AsInputTranscriptDeltaEvent().Delta})
					} else {
						want = append(want, [2]string{"assistant", rejected.AsOutputTranscriptDeltaEvent().Delta})
					}
					g.Close()
				}
				synctest.Wait()
				checkTranscriptFinish(t, &c, want)
			})
		})
	}
}
