package openai

import (
	"encoding/json"
	"errors"
	"io"
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

// BetaAgentTurnResultError distinguishes an incomplete observation from a known
// hosted outcome. Turn may be nil when no root turn was observed. Messages contains
// only completed final messages received before collection stopped.
// This beta helper is experimental.
type BetaAgentTurnResultError struct {
	// Reason is turn_failed, turn_cancelled, session_failed, requires_action,
	// observation_failed, observation_incomplete, incomplete_output,
	// unclassified_output, invalid_event, or unsupported_stream.
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
// initial root turn's final answer. Events consumed with Next are included. The
// result or error is cached, and repeated calls never advance to another turn.
// Defer stream.Close when iteration might stop early. This beta helper is experimental.
// The free function preserves NewStreaming's concrete ssestream.Stream return type.
func BetaAgentSessionFinalResult(stream *ssestream.Stream[AgentSessionEventUnion]) (*BetaAgentTurnResult, error) {
	if stream == nil {
		return nil, &BetaAgentTurnResultError{Reason: "unsupported_stream"}
	}
	c, ok := stream.Accumulator().(*betaAgentTurnCollector)
	if !ok {
		return nil, &BetaAgentTurnResultError{Reason: "unsupported_stream"}
	}
	if c.finalized {
		return c.result, c.err
	}
	for !c.blocked(nil) && !c.boundary && !c.sessionFailed && stream.Next() {
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
	for !c.blocked(s.handlers) && s.Next() {
	}
	if c.blocked(s.handlers) {
		_ = s.Close()
	}
	return c.finalResult(s.Err())
}

type betaAgentResultMessage struct {
	index    int64
	hasIndex bool
	message  *AgentSessionMessage
	done     bool
}

type betaAgentTurnCollector struct {
	sessionID     string
	turn          *Turn
	terminal      bool
	boundary      bool
	sessionFailed bool
	messages      []betaAgentResultMessage
	positions     map[string]int
	required      []AgentSessionRequiredActionUnion
	invalid       error
	finalized     bool
	result        *BetaAgentTurnResult
	err           error
}

func (c *betaAgentTurnCollector) Accumulate(event AgentSessionEventUnion) {
	if c.finalized || c.boundary {
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
			if err := json.Unmarshal([]byte(action.RawJSON()), &copy); err != nil {
				c.invalid = err
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
		if c.turn == nil && event.Turn.SubagentID == "" && event.Turn.ID != "" {
			c.setTurn(event.Turn)
		}
	}
	turnID := event.TurnID
	if turnID == "" && (event.Type == "agent.session.turn.item.added" || event.Type == "agent.session.turn.item.done") {
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
	case "agent.session.turn.item.added", "agent.session.turn.item.done":
		if event.Item.Type != "message" || event.Item.Role != "assistant" {
			return
		}
		entry := c.message(event.Item.ID)
		if entry == nil {
			return
		}
		if entry.done && event.Type == "agent.session.turn.item.added" {
			return
		}
		var message AgentSessionMessage
		if err := json.Unmarshal([]byte(event.Item.JSON.raw), &message); err != nil {
			c.invalid = err
			return
		}
		entry.message = &message
		if event.Type == "agent.session.turn.item.done" {
			entry.done = true
			entry.index = event.OutputIndex
			entry.hasIndex = event.JSON.OutputIndex.Valid()
		}
	case "agent.session.turn.output_text.delta", "agent.session.turn.output_text.done":
		c.message(event.ItemID)
	}
}

func (c *betaAgentTurnCollector) setTurn(turn Turn) {
	var copy Turn
	if err := json.Unmarshal([]byte(turn.RawJSON()), &copy); err != nil {
		c.invalid = err
		return
	}
	if copy.ID == "" || copy.SessionID == "" || (c.turn != nil && copy.ID != c.turn.ID) || (c.sessionID != "" && copy.SessionID != c.sessionID) {
		c.invalid = errors.New("turn snapshot lacks identity")
		return
	}
	c.turn = &copy
	c.sessionID = copy.SessionID
}

func (c *betaAgentTurnCollector) message(id string) *betaAgentResultMessage {
	if id == "" {
		c.invalid = errors.New("message lacks identity")
		return nil
	}
	if c.positions == nil {
		c.positions = map[string]int{}
	}
	if position, ok := c.positions[id]; ok {
		return &c.messages[position]
	}
	c.positions[id] = len(c.messages)
	c.messages = append(c.messages, betaAgentResultMessage{})
	return &c.messages[len(c.messages)-1]
}

func (c *betaAgentTurnCollector) blocked(handlers map[string]AgentToolHandler) bool {
	for _, action := range c.required {
		if action.Type == "function_call" && handlers[action.Name] != nil {
			continue
		}
		return true
	}
	return false
}

func (c *betaAgentTurnCollector) finalResult(cause error) (*BetaAgentTurnResult, error) {
	if c.finalized {
		return c.result, c.err
	}
	c.finalized = true
	entries := slices.Clone(c.messages)
	slices.SortStableFunc(entries, func(a, b betaAgentResultMessage) int {
		if a.hasIndex != b.hasIndex {
			if a.hasIndex {
				return -1
			}
			return 1
		}
		if a.hasIndex && b.hasIndex {
			if a.index < b.index {
				return -1
			}
			if a.index > b.index {
				return 1
			}
		}
		return 0
	})
	messages := make([]AgentSessionMessage, 0, len(entries))
	outputReason := ""
	for _, entry := range entries {
		if entry.message != nil && entry.message.Phase == "commentary" {
			continue
		}
		if entry.message == nil || !entry.done || entry.message.Status != "completed" {
			outputReason = "incomplete_output"
		} else if entry.message.Phase != "final_answer" {
			outputReason = "unclassified_output"
		} else {
			messages = append(messages, *entry.message)
		}
	}
	// Later stream failures do not invalidate an already observed turn boundary.
	if c.boundary {
		cause = nil
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
	case c.invalid != nil:
		reason = "invalid_event"
		cause = c.invalid
	case c.turn == nil || !c.terminal || c.turn.Status != "completed" || !c.boundary:
		reason = "observation_incomplete"
		cause = io.ErrUnexpectedEOF
	case outputReason != "":
		reason = outputReason
	}
	if reason != "" {
		c.err = &BetaAgentTurnResultError{Reason: reason, SessionID: c.sessionID, Turn: c.turn, Messages: messages, RequiredActions: c.required, Cause: cause}
		return nil, c.err
	}
	c.result = &BetaAgentTurnResult{Turn: *c.turn, Messages: messages}
	return c.result, nil
}
