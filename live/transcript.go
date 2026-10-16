package live

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math"
	"sync"
	"time"

	"github.com/openai/openai-go/v3/packages/param"
)

// ErrTranscriptBackpressure means a slow callback has left at least 1024
// updates waiting. Push has not consumed the event; it may be retried after
// callbacks catch up. Close and session.closed still flush accepted input.
var ErrTranscriptBackpressure = errors.New("live transcript callback backlog is full")

// TranscriptSegment is a complete display snapshot, not a server turn/item.
// Replace text on an existing ID. Previously delivered values remain unchanged.
type TranscriptSegment struct {
	ID         string
	PreviousID string
	Speaker    string
	Text       string
	StartMs    int64
	EndMs      int64
}

// TranscriptUpdate has an empty Reason for a display update, or a final
// "speaker_change", "inactivity", "timestamp_reset", "session_closed", "manual".
type TranscriptUpdate struct {
	Segment TranscriptSegment
	Reason  string
}

// TranscriptGrouperOptions controls display grouping. Zero options select
// defaults. Set BackchannelMaxDuration to param.NewOpt(time.Duration(0)) to
// disable acknowledgment suppression. Callbacks run in order outside the
// grouper lock, on the Push/Close caller or an owned timer; they may reenter it.
type TranscriptGrouperOptions struct {
	OnUpdate                  func(TranscriptUpdate)
	MinTurnSeparation         param.Opt[time.Duration]
	AssistantSilence          param.Opt[time.Duration]
	BackchannelMaxDuration    param.Opt[time.Duration]
	BackchannelIsolation      param.Opt[time.Duration]
	AdditionalAcknowledgments []string
}

// TranscriptGrouper accepts public Live server events explicitly, so callers
// keep raw events and complete control of transport lifetime. Close it on
// transport loss or teardown; it flushes once and stops only its owned timer.
// Duplicate tracking retains hashes of the latest 4096 accepted event IDs;
// older IDs may be accepted again. Initial possible acknowledgments wait for
// source-timestamp context or Close, rather than a network-arrival deadline.
type TranscriptGrouper struct {
	mu              sync.Mutex
	group           transcriptGrouping
	onUpdate        func(TranscriptUpdate)
	seen            map[[sha256.Size]byte]struct{}
	recentIDs       [][sha256.Size]byte
	nextID          int
	lastStart       map[string]int64
	pending         []transcriptFragment
	pendingAck      transcriptAcknowledgment
	pendingAckCount int
	pendingAckEnd   int64
	updates         []TranscriptUpdate
	dispatching     bool
	closed          bool
	timer           *time.Timer
	generation      uint64
	anchorMs        int64
	anchored        bool
	receivedAt      time.Time
	updateDue       time.Time
}

type transcriptFragment struct {
	speaker    string
	text       string
	start, end int64
	received   time.Time
}

// NewTranscriptGrouper creates a display projection with one identity space.
// It does not accept or open a connection. Negative or null durations fail.
func NewTranscriptGrouper(opts TranscriptGrouperOptions) (*TranscriptGrouper, error) {
	for _, duration := range []param.Opt[time.Duration]{opts.MinTurnSeparation, opts.AssistantSilence, opts.BackchannelMaxDuration, opts.BackchannelIsolation} {
		if param.IsNull(duration) || duration.Or(0) < 0 {
			return nil, errors.New("live transcript timings must be nonnegative durations")
		}
	}
	var id [8]byte
	if _, err := rand.Read(id[:]); err != nil {
		return nil, err
	}
	g := &TranscriptGrouper{
		onUpdate: opts.OnUpdate, seen: make(map[[sha256.Size]byte]struct{}), lastStart: make(map[string]int64),
		group: newTranscriptGrouping(opts, "segment_"+hex.EncodeToString(id[:])),
	}
	return g, nil
}

// Push consumes one typed event without mutating it. Other events and duplicate
// IDs are ignored. Invalid transcript intervals, empty event IDs, or Push after
// Close return an error. If callbacks fall behind, ErrTranscriptBackpressure
// rejects a new delta without consuming it. A real session.closed flushes and
// closes the grouper, even while callbacks are behind.
func (g *TranscriptGrouper) Push(event ServerEventUnion) error {
	now := time.Now()
	g.mu.Lock()
	if g.closed {
		g.mu.Unlock()
		return errors.New("live transcript grouper is closed")
	}
	var f transcriptFragment
	var id string
	var valid bool
	switch event.Type {
	case "session.closed":
		e := event.AsSessionClosedEvent()
		if !e.JSON.Type.Valid() || e.Type != "session.closed" ||
			!e.JSON.EventID.Valid() || e.EventID == "" ||
			!e.JSON.Reason.Valid() || e.Reason == "" || !e.JSON.Session.Valid() || !e.JSON.Usage.Valid() {
			g.mu.Unlock()
			return errors.New("live transcript terminal requires type, event ID, reason, session and usage")
		}
		if _, exists := g.seen[sha256.Sum256([]byte(e.EventID))]; exists {
			g.mu.Unlock()
			return nil
		}
		shouldDrain := g.enqueue(g.finish(now, "session_closed"))
		g.mu.Unlock()
		if shouldDrain {
			g.drain()
		}
		return nil
	case "session.input_transcript.delta":
		e := event.AsInputTranscriptDeltaEvent()
		f = transcriptFragment{"user", e.Delta, e.StartMs, e.EndMs, now}
		id = e.EventID
		valid = e.JSON.Type.Valid() && e.Type == "session.input_transcript.delta" &&
			e.JSON.EventID.Valid() && e.JSON.Delta.Valid() && e.JSON.StartMs.Valid() && e.JSON.EndMs.Valid()
	case "session.output_transcript.delta":
		e := event.AsOutputTranscriptDeltaEvent()
		f = transcriptFragment{"assistant", e.Delta, e.StartMs, e.EndMs, now}
		id = e.EventID
		valid = e.JSON.Type.Valid() && e.Type == "session.output_transcript.delta" &&
			e.JSON.EventID.Valid() && e.JSON.Delta.Valid() && e.JSON.StartMs.Valid() && e.JSON.EndMs.Valid()
	default:
		g.mu.Unlock()
		return nil
	}
	if !valid || id == "" || f.start < 0 || f.end < f.start {
		g.mu.Unlock()
		return errors.New("live transcript delta requires its matching type, text, an event ID and a nonnegative ordered interval")
	}
	if f.text == "" {
		g.mu.Unlock()
		return nil
	}
	digest := sha256.Sum256([]byte(id))
	if _, duplicate := g.seen[digest]; duplicate {
		g.mu.Unlock()
		return nil
	}
	// Reject before changing grouping, timer or duplicate state. Callback
	// reentry cannot block waiting for the callback itself to drain the queue.
	// Already accepted input can still finish its current and buffered turns.
	if len(g.updates) >= 1024 {
		// A recovered callback panic may leave a full but idle queue.
		// Restart that dispatcher without admitting the new event.
		shouldDrain := g.enqueue(nil)
		g.mu.Unlock()
		if shouldDrain {
			g.drain()
		}
		return ErrTranscriptBackpressure
	}
	g.remember(digest)
	var emitted []TranscriptUpdate
	firstOverlap, firstAfter, initialAssistant := -1, -1, false
	if len(g.pending) > 0 {
		// Pending fragments share a speaker: any opposite speaker flushes
		// them in this Push, whether the ranges overlap or not.
		if g.pending[0].speaker != f.speaker {
			for i, p := range g.pending {
				if transcriptOverlap(p.start, p.end, f.start, f.end) {
					firstOverlap = i
					break
				}
				if firstAfter < 0 && p.start > f.start {
					firstAfter = i
				}
			}
		}
		last := g.pending[len(g.pending)-1]
		initialAssistant = g.group.current == nil && last.speaker == "assistant" && f.speaker == "assistant" &&
			f.start >= last.start && float64(f.start) < float64(last.end)+g.group.silence
	}
	previous, hasPrevious := g.lastStart[f.speaker]
	if hasPrevious && f.start < previous {
		var before, after []transcriptFragment
		for _, part := range g.pending {
			// Move only positively overlapping fragments into the new epoch.
			// A disjoint fragment remains on the old timeline when its age is ambiguous.
			newOverlap := transcriptOverlap(part.start, part.end, f.start, f.end)
			if part.start < previous && newOverlap {
				after = append(after, part)
			} else {
				before = append(before, part)
			}
		}
		g.pending = nil
		g.pendingAckCount = 0
		emitted = g.commit(before)
		emitted = append(emitted, g.group.close(g.sourceNow(now), "timestamp_reset")...)
		clear(g.lastStart)
		g.anchored = false
		emitted = append(emitted, g.commit(append(after, f))...)
	} else if firstOverlap < 0 && firstAfter >= 0 {
		emitted = g.commit(g.pending[:firstAfter])
		g.pending = g.pending[firstAfter:]
		g.pendingAckCount = 0
		emitted = append(emitted, g.commit([]transcriptFragment{f})...)
		emitted = append(emitted, g.flushPending()...)
	} else if firstOverlap >= 0 || initialAssistant {
		if firstOverlap > 0 && g.group.current == nil && g.pending[0].speaker == "assistant" {
			prefix := firstOverlap
			end := g.pending[0].end
			ack := transcriptAcknowledgment{node: g.group.ackRoot}
			for _, part := range g.pending {
				end = max(end, part.end)
				g.group.advanceAck(&ack, part.text)
			}
			// A short utterance may still complete an acknowledgment. Keep
			// its prefix until the same silence boundary used to batch it.
			if ack.seen && ack.node != nil && float64(end-g.pending[0].start) < g.group.backchannelDuration {
				for prefix > 0 && float64(g.pending[prefix].start) < float64(g.pending[prefix-1].end)+g.group.silence {
					prefix--
				}
			}
			if prefix > 0 {
				emitted = g.commit(g.pending[:prefix])
				g.pending = g.pending[prefix:]
				g.pendingAckCount = 0
			}
		}
		g.pending = append(g.pending, f)
		if firstOverlap >= 0 {
			emitted = append(emitted, g.flushPending()...)
		}
	} else {
		emitted = g.flushPending()
		if g.group.speaker() == f.speaker {
			emitted = append(emitted, g.commit([]transcriptFragment{f})...)
		} else {
			g.pending = append(g.pending, f)
		}
	}
	g.schedule(now)
	shouldDrain := g.enqueue(emitted)
	g.mu.Unlock()
	if shouldDrain {
		g.drain()
	}
	return nil
}

func (g *TranscriptGrouper) remember(digest [sha256.Size]byte) {
	const limit = 4096
	if len(g.recentIDs) == limit {
		delete(g.seen, g.recentIDs[g.nextID])
		g.recentIDs[g.nextID] = digest
		g.nextID = (g.nextID + 1) % limit
	} else {
		g.recentIDs = append(g.recentIDs, digest)
	}
	g.seen[digest] = struct{}{}
}

// Close finalizes once, cancels the owned timer and permits an in-progress
// callback to finish. It cannot close a socket because it never owns one.
func (g *TranscriptGrouper) Close() {
	g.mu.Lock()
	var emitted []TranscriptUpdate
	if !g.closed {
		emitted = g.finish(time.Now(), "manual")
	}
	shouldDrain := g.enqueue(emitted)
	g.mu.Unlock()
	if shouldDrain {
		g.drain()
	}
}

func (g *TranscriptGrouper) finish(now time.Time, reason string) []TranscriptUpdate {
	g.closed = true
	g.stopTimer()
	emitted := g.flushPending()
	emitted = append(emitted, g.group.close(g.sourceNow(now), reason)...)
	g.seen = nil
	g.recentIDs = nil
	return emitted
}

func (g *TranscriptGrouper) commit(parts []transcriptFragment) []TranscriptUpdate {
	if len(parts) == 0 {
		return nil
	}
	var emitted []TranscriptUpdate
	for _, p := range parts {
		if prev, exists := g.lastStart[p.speaker]; exists && p.start < prev {
			emitted = g.group.close(g.sourceNow(p.received), "timestamp_reset")
			clear(g.lastStart)
			g.anchored = false
			break
		}
	}
	for _, p := range parts {
		g.lastStart[p.speaker] = p.start
		if !g.anchored || g.anchorMs < p.end {
			g.anchorMs = p.end
			g.receivedAt = p.received
		}
		g.anchored = true
	}
	return append(emitted, g.group.process(parts)...)
}

func (g *TranscriptGrouper) flushPending() []TranscriptUpdate {
	parts := g.pending
	g.pending = nil
	g.pendingAckCount = 0
	return g.commit(parts)
}

func (g *TranscriptGrouper) sourceNow(now time.Time) float64 {
	if !g.anchored {
		return 0
	}
	return float64(g.anchorMs) + max(0, float64(now.Sub(g.receivedAt))/float64(time.Millisecond))
}

func (g *TranscriptGrouper) stopTimer() {
	g.generation++
	if g.timer != nil {
		g.timer.Stop()
		g.timer = nil
	}
}

func (g *TranscriptGrouper) pendingAcknowledgment() bool {
	if g.group.current != nil || len(g.pending) == 0 || g.pending[0].speaker != "assistant" {
		return false
	}
	// Only an overlapping or adjacent tail is known to belong to the prior
	// assistant without waiting for the other speaker's source interval.
	if g.group.hasAssistantEnd && g.pending[0].start <= g.group.lastAssistantEnd {
		return false
	}
	if g.pendingAckCount == 0 {
		g.pendingAck = transcriptAcknowledgment{node: g.group.ackRoot}
		g.pendingAckEnd = g.pending[0].end
	}
	// Advance only over new fragments. Deferred whitespace and punctuation may
	// last arbitrarily long without changing the acknowledgment decision.
	for _, part := range g.pending[g.pendingAckCount:] {
		g.pendingAckEnd = max(g.pendingAckEnd, part.end)
		g.group.advanceAck(&g.pendingAck, part.text)
	}
	g.pendingAckCount = len(g.pending)
	return g.pendingAck.node != nil &&
		float64(g.pendingAckEnd-g.pending[0].start) < g.group.backchannelDuration
}

func (g *TranscriptGrouper) schedule(now time.Time) {
	g.stopTimer()
	if g.closed {
		return
	}
	wait := math.Inf(1)
	if len(g.pending) > 0 {
		if !g.pendingAcknowledgment() {
			wait = float64(g.pending[0].received.Add(50 * time.Millisecond).Sub(now))
		}
	} else if due, ok := g.group.deadline(); ok {
		wait = (due - g.sourceNow(now)) * float64(time.Millisecond)
	}
	if g.group.updatePending() {
		if g.updateDue.IsZero() {
			g.updateDue = now.Add(50 * time.Millisecond)
		}
		wait = min(wait, float64(g.updateDue.Sub(now)))
	} else {
		g.updateDue = time.Time{}
	}
	if math.IsInf(wait, 1) {
		return
	}
	delay := time.Duration(max(0, wait))
	if wait >= float64(math.MaxInt64) {
		delay = time.Duration(math.MaxInt64)
	}
	generation := g.generation
	g.timer = time.AfterFunc(delay, func() {
		g.mu.Lock()
		if g.closed || generation != g.generation {
			g.mu.Unlock()
			return
		}
		tick := time.Now()
		var emitted []TranscriptUpdate
		if len(g.pending) == 0 || !tick.Before(g.pending[0].received.Add(50*time.Millisecond)) {
			emitted = g.flushPending()
			emitted = append(emitted, g.group.advance(g.sourceNow(tick), false)...)
		}
		emitted = append(emitted, g.group.flushUpdate()...)
		g.schedule(tick)
		shouldDrain := g.enqueue(emitted)
		g.mu.Unlock()
		if shouldDrain {
			g.drain()
		}
	})
}

func (g *TranscriptGrouper) enqueue(emitted []TranscriptUpdate) bool {
	for _, update := range emitted {
		last := len(g.updates) - 1
		if last >= 0 && update.Reason == "" && g.updates[last].Reason == "" &&
			g.updates[last].Segment.ID == update.Segment.ID {
			g.updates[last] = update
		} else {
			g.updates = append(g.updates, update)
		}
	}
	if g.dispatching || len(g.updates) == 0 {
		return false
	}
	g.dispatching = true
	return true
}

func (g *TranscriptGrouper) drain() {
	complete := false
	defer func() {
		if !complete {
			g.mu.Lock()
			g.dispatching = false
			g.mu.Unlock()
		}
	}()
	for {
		g.mu.Lock()
		if len(g.updates) == 0 {
			g.updates = nil
			g.dispatching = false
			complete = true
			g.mu.Unlock()
			return
		}
		update := g.updates[0]
		g.updates[0] = TranscriptUpdate{}
		g.updates = g.updates[1:]
		g.mu.Unlock()
		if g.onUpdate != nil {
			g.onUpdate(update)
		}
	}
}
