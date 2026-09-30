package openai

import (
	"encoding/json"
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/openai/openai-go/v3/packages/ssestream"
)

// BetaAgentTurnResult is the completed root turn and its final assistant messages.
// This beta helper is experimental. It does not contain the session transcript.
type BetaAgentTurnResult struct {
	Turn     Turn
	Messages []AgentSessionMessage
}

func (r BetaAgentTurnResult) SessionID() string { return r.Turn.SessionID }
func (r BetaAgentTurnResult) TurnID() string    { return r.Turn.ID }

// OutputText joins the final messages' text without adding separators or making requests.
func (r BetaAgentTurnResult) OutputText() string {
	var text strings.Builder
	for _, message := range r.Messages {
		text.WriteString(message.OutputText())
	}
	return text.String()
}

// BetaAgentTurnResultError distinguishes an interrupted observation from a known
// hosted outcome. Turn may be nil when no root turn was observed. Messages contains
// completed final messages received before collection stopped.
// This beta helper is experimental.
type BetaAgentTurnResultError struct {
	// Reason is turn_failed, turn_cancelled, session_failed, requires_action,
	// observation_failed, observation_incomplete, or unsupported_stream.
	Reason          string
	SessionID       string
	Turn            *Turn
	Messages        []AgentSessionMessage
	RequiredActions []AgentSessionRequiredActionUnion
	Cause           error
}

func (e *BetaAgentTurnResultError) Error() string {
	return "cannot collect beta agent turn result: " + e.Reason
}
func (e *BetaAgentTurnResultError) Unwrap() error { return e.Cause }

// BetaAgentSessionFinalResult consumes a session creation stream and returns its
// initial root turn's final answer. Creation must include initial Input; an idle
// self-hosted creation without Input has no initial result to collect.
// Events consumed with Next are included. Repeated calls return the cached result
// or error. This beta helper is experimental. The free function preserves
// NewStreaming's concrete ssestream.Stream return type.
func BetaAgentSessionFinalResult(stream *ssestream.Stream[AgentSessionEventUnion]) (*BetaAgentTurnResult, error) {
	if stream == nil {
		return nil, &BetaAgentTurnResultError{Reason: "unsupported_stream"}
	}
	c, ok := stream.BetaAccumulator().(*betaAgentTurnCollector)
	if !ok {
		return nil, &BetaAgentTurnResultError{Reason: "unsupported_stream"}
	}
	if c.finalized {
		return c.result, c.err
	}
	for !c.stopped(nil) && stream.Next() {
	}
	_ = stream.Close()
	return c.finalResult(stream.Err())
}

// FinalResult consumes the follow-up stream, including registered tool handlers,
// and returns the selected root turn's final messages. It can be called after
// iteration and is cached. Unhandled actions return a BetaAgentTurnResultError
// rather than waiting indefinitely. Closing observation does not cancel the turn.
func (s *AgentSessionStream) FinalResult() (*BetaAgentTurnResult, error) {
	c := &s.collector
	if c.finalized {
		return c.result, c.err
	}
	for !c.stopped(s.handlers) && s.Next() {
	}
	_ = s.Close()
	return c.finalResult(s.Err())
}

type betaAgentTurnCollector struct {
	sessionID     string
	turn          *Turn
	terminal      bool
	boundary      bool
	sessionFailed bool
	messages      map[int64]AgentSessionMessage
	required      []AgentSessionRequiredActionUnion
	decodeErr     error
	finalized     bool
	result        *BetaAgentTurnResult
	err           error
}

func (c *betaAgentTurnCollector) Accumulate(event AgentSessionEventUnion) {
	if c.finalized || c.boundary || c.decodeErr != nil {
		return
	}
	if c.sessionID == "" {
		c.sessionID = event.SessionID
	}
	switch event.Type {
	case "agent.session.created":
		c.sessionID = event.Session.ID
	case "agent.session.requires_action":
		c.required = nil
		for _, action := range event.Session.RequiredActions {
			var copy AgentSessionRequiredActionUnion
			if c.decodeErr = json.Unmarshal([]byte(action.RawJSON()), &copy); c.decodeErr != nil {
				return
			}
			c.required = append(c.required, copy)
		}
	case "agent.session.in_progress":
		c.required = nil
	case "agent.session.failed":
		c.sessionFailed = true
	case "agent.session.idle":
		if c.terminal {
			c.boundary = true
			c.required = nil
		}
	case "agent.session.turn.created":
		if c.turn == nil && event.Turn.SubagentID == "" {
			c.setTurn(event.Turn)
		}
	}
	turnID := event.TurnID
	if event.Type == "agent.session.turn.item.done" {
		turnID = event.Item.TurnID
	}
	if c.turn == nil || turnID != c.turn.ID {
		return
	}
	switch event.Type {
	case "agent.session.turn.in_progress":
		c.setTurn(event.Turn)
	case "agent.session.turn.completed", "agent.session.turn.failed", "agent.session.turn.cancelled":
		c.setTurn(event.Turn)
		c.terminal = true
		c.required = nil
	case "agent.session.turn.item.done":
		if event.Item.Type != "message" || event.Item.Role != "assistant" || event.Item.Status != "completed" || event.Item.Phase == "commentary" {
			return
		}
		// Completed items are authoritative, including legacy messages with no phase.
		// Decode once into the history-message shape to retain annotations and isolate
		// the result from callers mutating the event returned by Current.
		var message AgentSessionMessage
		if c.decodeErr = json.Unmarshal([]byte(event.Item.JSON.raw), &message); c.decodeErr != nil {
			return
		}
		if c.messages == nil {
			c.messages = make(map[int64]AgentSessionMessage)
		}
		c.messages[event.OutputIndex] = message
	}
}

func (c *betaAgentTurnCollector) setTurn(turn Turn) {
	var copy Turn
	c.decodeErr = json.Unmarshal([]byte(turn.RawJSON()), &copy)
	if c.decodeErr == nil {
		c.turn = &copy
		c.sessionID = copy.SessionID
	}
}

func (c *betaAgentTurnCollector) stopped(handlers map[string]AgentToolHandler) bool {
	if c.decodeErr != nil || c.boundary || c.sessionFailed || (c.terminal && c.turn != nil && c.turn.Status != "completed") {
		return true
	}
	for _, action := range c.required {
		if action.Type != "function_call" || handlers[action.Name] == nil {
			return true
		}
	}
	return false
}

func (c *betaAgentTurnCollector) finalResult(cause error) (*BetaAgentTurnResult, error) {
	if c.finalized {
		return c.result, c.err
	}
	c.finalized = true
	defer func() { c.messages = nil; c.required = nil; c.turn = nil; c.decodeErr = nil }()
	messages := make([]AgentSessionMessage, 0, len(c.messages))
	for _, index := range slices.Sorted(maps.Keys(c.messages)) {
		messages = append(messages, c.messages[index])
	}
	// A later stream error cannot invalidate an already observed turn boundary.
	if c.boundary {
		cause = nil
	}
	if c.decodeErr != nil {
		cause = c.decodeErr
	}
	reason := ""
	switch {
	case c.sessionFailed:
		reason = "session_failed"
	case c.turn != nil && c.turn.Status == "failed":
		reason = "turn_failed"
	case c.turn != nil && c.turn.Status == "cancelled":
		reason = "turn_cancelled"
	case len(c.required) > 0:
		reason = "requires_action"
	case cause != nil:
		reason = "observation_failed"
	case c.turn == nil || !c.terminal || c.turn.Status != "completed" || !c.boundary:
		reason = "observation_incomplete"
		cause = io.ErrUnexpectedEOF
	}
	if reason != "" {
		c.err = &BetaAgentTurnResultError{Reason: reason, SessionID: c.sessionID, Turn: c.turn, Messages: messages, RequiredActions: c.required, Cause: cause}
		return nil, c.err
	}
	c.result = &BetaAgentTurnResult{Turn: *c.turn, Messages: messages}
	return c.result, nil
}
