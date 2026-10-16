package live

import (
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type transcriptTurn struct {
	id, previous, speaker string
	text                  strings.Builder
	start, end            int64
	received              time.Time
	emitted               int
	canDrop               bool
	continued             bool
	overlap               bool
	ack                   transcriptAcknowledgment
	ranges                []transcriptRange
	starts                []int64
}

type transcriptRange struct {
	start, end int64
}

func (t *transcriptTurn) snapshot() TranscriptSegment {
	return TranscriptSegment{t.id, t.previous, t.speaker, t.text.String(), t.start, t.end}
}

type transcriptGrouping struct {
	current, buffered                                   *transcriptTurn
	prefix, lastID                                      string
	nextID                                              uint64
	lastAssistantEnd                                    int64
	hasAssistantEnd                                     bool
	separation, silence, backchannelDuration, isolation float64
	ackRoot                                             *transcriptAckNode
	maxPhrase                                           int
}

func newTranscriptGrouping(opts TranscriptGrouperOptions, prefix string) transcriptGrouping {
	ms := func(d time.Duration) float64 { return float64(d) / float64(time.Millisecond) }
	g := transcriptGrouping{
		prefix: prefix, ackRoot: &transcriptAckNode{},
		separation:          ms(opts.MinTurnSeparation.Or(500 * time.Millisecond)),
		silence:             ms(opts.AssistantSilence.Or(2000 * time.Millisecond)),
		backchannelDuration: ms(opts.BackchannelMaxDuration.Or(1000 * time.Millisecond)),
		isolation:           ms(opts.BackchannelIsolation.Or(2000 * time.Millisecond)),
	}
	for _, phrase := range []string{"aha", "alright", "gotcha", "hm", "hmm", "mhm", "mm", "mm hmm", "okay", "ok", "right", "sure", "uh huh", "yeah", "yep", "yes"} {
		g.addPhrase(phrase)
	}
	for _, phrase := range opts.AdditionalAcknowledgments {
		g.addPhrase(phrase)
	}
	return g
}

func (g *transcriptGrouping) speaker() string {
	if g.current == nil {
		return ""
	}
	return g.current.speaker
}

func (g *transcriptGrouping) process(parts []transcriptFragment) []TranscriptUpdate {
	var first *transcriptFragment
	for i := range parts {
		if first == nil || parts[i].start < first.start {
			first = &parts[i]
		}
	}
	if first == nil {
		return nil
	}
	var emitted []TranscriptUpdate
	for {
		due, ok := g.deadline()
		if !ok || due > float64(first.start) {
			break
		}
		emitted = append(emitted, g.advance(due, false)...)
	}
	same, incomingUser := false, false
	for _, p := range parts {
		same = same || p.speaker == g.speaker()
		incomingUser = incomingUser || p.speaker == "user"
	}
	if g.current != nil && !same {
		emitted = append(emitted, g.advance(float64(first.start), incomingUser)...)
	}
	preferred := g.speaker()
	if preferred == "" {
		preferred = "user"
	} else {
		// Keep the current speaker's priority for tied ranges, but do not
		// ingest their later fragment ahead of speech that already started.
		preferredAtStart := false
		for _, part := range parts {
			preferredAtStart = preferredAtStart || part.speaker == preferred && part.start == first.start
		}
		if !preferredAtStart {
			preferred = first.speaker
		}
	}
	for _, wanted := range []bool{true, false} {
		for _, part := range parts {
			if (part.speaker == preferred) == wanted {
				emitted = append(emitted, g.ingest(part)...)
			}
		}
	}
	return emitted
}

func (g *transcriptGrouping) deadline() (float64, bool) {
	if g.current == nil {
		return 0, false
	}
	if g.buffered != nil {
		separation := float64(g.current.end) + g.separation
		if g.mightBackchannel() && !g.recentAssistant() {
			return max(separation, float64(g.buffered.end)+g.isolation), true
		}
		return separation, true
	}
	return float64(g.current.end) + g.silence, g.speaker() == "assistant"
}

func (g *transcriptGrouping) advance(now float64, incomingUser bool) []TranscriptUpdate {
	if g.current != nil && g.buffered != nil && now >= float64(g.current.end)+g.separation && !g.keepBackchannel(now) {
		g.buffered = g.maybeDrop(now, nil)
		if g.buffered != nil {
			return g.promote()
		}
	}
	if g.speaker() == "assistant" && !incomingUser && now >= float64(g.current.end)+g.silence {
		return g.finishCurrent("inactivity")
	}
	return nil
}

func (g *transcriptGrouping) ingest(f transcriptFragment) []TranscriptUpdate {
	if g.current == nil {
		g.current = g.newTurn(f)
		return g.emit(g.current, false)
	}
	if f.speaker == g.speaker() {
		g.appendTurn(g.current, f)
		if f.speaker == "user" && g.buffered != nil {
			g.buffered.continued = g.buffered.continued || f.end >= g.buffered.end
			g.buffered.overlap = g.buffered.overlap || g.buffered.overlaps(f)
		}
		return g.emit(g.current, false)
	}
	if g.speaker() == "assistant" {
		g.buffer(f)
		return g.promote()
	}
	if g.buffered != nil && float64(f.start) >= float64(g.buffered.end)+g.silence {
		g.buffered = g.maybeDrop(float64(f.start), &f)
		if g.buffered == nil {
			return g.ingest(f)
		}
		emitted := g.promote()
		emitted = append(emitted, g.finishCurrent("inactivity")...)
		return append(emitted, g.ingest(f)...)
	}
	start := f.start
	if g.buffered != nil {
		start = g.buffered.start
	}
	within := float64(f.end-start) < g.backchannelDuration
	var ack transcriptAcknowledgment
	if g.buffered == nil {
		ack.node = g.ackRoot
	} else {
		ack = g.buffered.ack
	}
	if within {
		g.advanceAck(&ack, f.text)
	} else {
		ack.node = nil
	}
	if float64(f.start) < float64(g.current.end)+g.separation ||
		g.mightBackchannel() && float64(f.start) < float64(g.buffered.end)+g.isolation {
		g.buffer(f)
		g.buffered.canDrop = ack.seen && ack.node != nil && ack.node.terminal
		g.buffered.ack = ack
		return nil
	}
	if ack.seen && ack.node != nil && ack.node.terminal {
		hadBuffered := g.buffered != nil
		g.buffer(f)
		g.buffered.canDrop = hadBuffered
		g.buffered.ack = ack
		return nil
	}
	g.buffered = g.maybeDrop(-1, &f)
	emitted := g.finishCurrent("speaker_change")
	if g.buffered != nil {
		g.current, g.buffered = g.buffered, nil
		g.appendTurn(g.current, f)
	} else {
		g.current = g.newTurn(f)
	}
	return append(emitted, g.emit(g.current, false)...)
}

func (g *transcriptGrouping) newTurn(f transcriptFragment) *transcriptTurn {
	t := &transcriptTurn{
		id:      g.prefix + "_" + strconv.FormatUint(g.nextID, 10),
		speaker: f.speaker, start: f.start, end: f.end, received: f.received, canDrop: true,
	}
	g.nextID++
	g.appendTurn(t, f)
	return t
}

func (g *transcriptGrouping) appendTurn(t *transcriptTurn, f transcriptFragment) {
	t.text.WriteString(f.text)
	if f.end > t.end {
		t.end, t.received = f.end, f.received
	}
	if len(t.starts) == 0 || t.starts[len(t.starts)-1] != f.start {
		t.starts = append(t.starts, f.start)
	}
	// Points match only an exact fragment start and must not slow range searches.
	if f.start == f.end {
		return
	}
	if len(t.ranges) == 0 || f.start > t.ranges[len(t.ranges)-1].end {
		t.ranges = append(t.ranges, transcriptRange{f.start, f.end})
	} else {
		last := &t.ranges[len(t.ranges)-1]
		last.end = max(last.end, f.end)
	}
}

func (g *transcriptGrouping) buffer(f transcriptFragment) {
	if g.buffered == nil {
		g.buffered = g.newTurn(f)
	} else {
		g.appendTurn(g.buffered, f)
	}
	g.buffered.overlap = g.buffered.overlap || g.current.overlaps(f)
}

func (g *transcriptGrouping) promote() []TranscriptUpdate {
	emitted := g.finishCurrent("speaker_change")
	g.current, g.buffered = g.buffered, nil
	return append(emitted, g.emit(g.current, false)...)
}

func (g *transcriptGrouping) finishCurrent(reason string) []TranscriptUpdate {
	t := g.current
	g.current = nil
	if t == nil {
		return nil
	}
	return g.finish(t, reason)
}

func (g *transcriptGrouping) finish(t *transcriptTurn, reason string) []TranscriptUpdate {
	if t.speaker == "assistant" {
		g.lastAssistantEnd, g.hasAssistantEnd = t.end, true
	}
	if t.emitted == 0 {
		return nil
	}
	return append(g.emit(t, true), TranscriptUpdate{t.snapshot(), reason})
}

func (g *transcriptGrouping) emit(t *transcriptTurn, force bool) []TranscriptUpdate {
	size := t.text.Len()
	if size == 0 || t.emitted != 0 && (size == t.emitted || !force && size < t.emitted*2) {
		return nil
	}
	if t.emitted == 0 {
		t.previous, g.lastID = g.lastID, t.id
	}
	t.emitted = size
	return []TranscriptUpdate{{Segment: t.snapshot()}}
}

func (g *transcriptGrouping) updatePending() bool {
	return g.current != nil && g.current.emitted != 0 && g.current.text.Len() > g.current.emitted
}

func (g *transcriptGrouping) flushUpdate() []TranscriptUpdate {
	if !g.updatePending() {
		return nil
	}
	return g.emit(g.current, true)
}

func (g *transcriptGrouping) close(now float64, reason string) []TranscriptUpdate {
	pending := g.maybeDrop(now, nil)
	emitted := g.finishCurrent(reason)
	g.buffered = nil
	if pending != nil {
		emitted = append(emitted, g.emit(pending, false)...)
		emitted = append(emitted, g.finish(pending, reason)...)
	}
	if reason == "timestamp_reset" {
		g.hasAssistantEnd = false
	}
	return emitted
}

func (g *transcriptGrouping) mightBackchannel() bool {
	return g.speaker() == "user" && g.buffered != nil && g.buffered.speaker == "assistant" &&
		g.buffered.ack.seen && g.buffered.ack.node != nil && g.buffered.overlap &&
		float64(g.buffered.end-g.buffered.start) < g.backchannelDuration
}

func (t *transcriptTurn) overlaps(f transcriptFragment) bool {
	i := sort.Search(len(t.starts), func(i int) bool { return t.starts[i] >= f.start })
	if i < len(t.starts) && t.starts[i] == f.start {
		return true
	}
	if f.start == f.end {
		return false
	}
	i = sort.Search(len(t.ranges), func(i int) bool { return t.ranges[i].end > f.start })
	return i < len(t.ranges) && t.ranges[i].start < f.end
}

// Timeline intervals touching at an end do not overlap. Appended context with
// the same start may have zero duration and retains initial user priority.
func transcriptOverlap(start, end, otherStart, otherEnd int64) bool {
	return start == otherStart || start < end && otherStart < otherEnd && start < otherEnd && otherStart < end
}

func (g *transcriptGrouping) userContinued() bool {
	return g.mightBackchannel() && g.buffered.canDrop && (g.buffered.continued || g.current.end > g.buffered.end)
}
func (g *transcriptGrouping) recentAssistant() bool {
	return g.current != nil && g.buffered != nil && g.hasAssistantEnd &&
		(g.buffered.start <= g.lastAssistantEnd ||
			float64(g.buffered.start) < float64(g.lastAssistantEnd)+g.isolation && g.buffered.start <= g.current.start)
}
func (g *transcriptGrouping) keepBackchannel(now float64) bool {
	return g.mightBackchannel() && !g.recentAssistant() &&
		now < float64(g.buffered.end)+g.isolation
}
func (g *transcriptGrouping) maybeDrop(now float64, following *transcriptFragment) *transcriptTurn {
	if !g.mightBackchannel() {
		return g.buffered
	}
	if g.recentAssistant() {
		return g.buffered
	}
	if g.userContinued() {
		return nil
	}
	if following != nil && float64(following.start) < float64(g.buffered.end)+g.isolation {
		return g.buffered
	}
	if following == nil && (now < 0 || now < float64(g.buffered.end)+g.isolation) {
		return g.buffered
	}
	if g.buffered.canDrop {
		return nil
	}
	return g.buffered
}

// Only Live-policy whitespace and punctuation trim acknowledgment boundaries;
// notably Unicode NEL is not a separator. '-' behaves as a space.
func transcriptSpace(r rune) bool {
	return strings.ContainsRune("\t\n\v\f\r \u00a0\u1680\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200a\u2028\u2029\u202f\u205f\u3000\ufeff", r)
}
func transcriptSeparator(r rune) bool {
	return transcriptSpace(r) || strings.ContainsRune(".,!?;:\"'()[]{}", r)
}

type transcriptAckNode struct {
	children map[rune]*transcriptAckNode
	terminal bool
}
type transcriptAcknowledgment struct {
	node           *transcriptAckNode
	pending        []rune
	overflow, seen bool
}

func (g *transcriptGrouping) addPhrase(raw string) {
	raw = strings.ReplaceAll(strings.ToLower(raw), "-", " ")
	raw = strings.TrimFunc(raw, transcriptSeparator)
	if raw == "" {
		return
	}
	node, size, space := g.ackRoot, 0, false
	for _, r := range raw {
		if transcriptSpace(r) {
			if space {
				continue
			}
			r, space = ' ', true
		} else {
			space = false
		}
		if node.children == nil {
			node.children = make(map[rune]*transcriptAckNode)
		}
		next := node.children[r]
		if next == nil {
			next = &transcriptAckNode{}
			node.children[r] = next
		}
		node = next
		size++
	}
	node.terminal = true
	g.maxPhrase = max(g.maxPhrase, size)
}

// Consume each rune once. Deferred trailing separators preserve trim semantics,
// while their storage never exceeds the longest configured normalized phrase.
func (g *transcriptGrouping) advanceAck(ack *transcriptAcknowledgment, text string) {
	for _, raw := range text {
		if ack.node == nil {
			return
		}
		r := unicode.ToLower(raw)
		if r == '-' {
			r = ' '
		}
		if transcriptSeparator(r) {
			if !ack.seen {
				continue
			}
			if transcriptSpace(r) {
				r = ' '
				if len(ack.pending) > 0 && ack.pending[len(ack.pending)-1] == ' ' {
					continue
				}
			}
			if len(ack.pending) < g.maxPhrase {
				ack.pending = append(ack.pending, r)
			} else {
				ack.overflow = true
			}
			continue
		}
		if ack.overflow {
			ack.node = nil
			return
		}
		for _, p := range ack.pending {
			ack.node = ack.node.children[p]
			if ack.node == nil {
				return
			}
		}
		ack.pending = nil
		ack.node = ack.node.children[r]
		ack.seen = true
	}
}
