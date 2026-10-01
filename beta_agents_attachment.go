package openai

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/pagination"
)

// Attachment keeps only root identity and status. The existing SSE dispatcher
// remains the sole source of calls, including pending calls replayed on attach.
type betaAgentAttachment struct {
	ctx               context.Context
	baselineID        string
	turn              *Turn
	idle              bool
	terminal          bool
	observationFailed bool
}

func betaAgentTurnTerminal(turn *Turn) bool {
	return turn != nil && (turn.Status == "completed" || turn.Status == "failed" || turn.Status == "cancelled")
}

func (s *AgentSessionStream) attachSession(ctx context.Context) *AgentSessionStream {
	s.attachment = &betaAgentAttachment{ctx: ctx}
	s.collector.sessionID = s.sessionID
	session, err := s.attachmentSession(s.ctx)
	if err != nil {
		s.finish(err)
		return s
	}
	if session.Status == "failed" {
		s.collector.sessionFailed = true
		_ = s.Close()
		return s
	}
	baseline, err := s.latestAttachmentRoot(s.ctx)
	if err != nil {
		s.finish(err)
		return s
	}
	if baseline != nil {
		s.attachment.baselineID = baseline.ID
	}
	if baseline != nil && !betaAgentTurnTerminal(baseline) {
		s.selectAttachmentTurn(baseline)
	}
	events := s.sessions.Events
	capture, require := agentStreamResponseGuard[http.Response]()
	events.Options = append([]option.RequestOption{capture}, events.Options...)
	s.stream = events.StreamStreaming(s.ctx, s.sessionID, append(slices.Clone(s.options), require)...)
	if err = s.stream.Err(); err != nil {
		s.finish(err)
		return s
	}
	// An exact selected turn can finish while a successor makes the session active.
	if s.turnID != "" {
		turn, getErr := s.attachmentTurn(s.ctx, s.turnID)
		if getErr != nil {
			s.finish(getErr)
			return s
		}
		s.selectAttachmentTurn(turn)
		if s.attachment.terminal {
			_ = s.Close()
			return s
		}
	}
	session, err = s.attachmentSession(s.ctx)
	if err != nil {
		s.finish(err)
		return s
	}
	if session.Status == "failed" {
		s.collector.sessionFailed = true
		_ = s.Close()
		return s
	}
	if s.turnID == "" {
		// Read the latest root after the session snapshot. This is the bounded
		// attachment snapshot; work starting later belongs to a future attachment.
		latest, listErr := s.latestAttachmentRoot(s.ctx)
		if listErr != nil {
			s.finish(listErr)
			return s
		}
		if latest != nil && latest.ID != s.attachment.baselineID {
			s.selectAttachmentTurn(latest)
		}
	}
	if s.attachment.terminal || (session.Status == "idle" && s.turnID == "") {
		s.attachment.idle = session.Status == "idle"
		_ = s.Close()
	}
	// A non-idle session with no new visible root may still be starting work.
	// Keep observing; an unchanged historical terminal root is not a no-work signal.
	return s
}

func (s *AgentSessionStream) selectAttachmentTurn(turn *Turn) {
	s.turnID = turn.ID
	s.attachment.turn = turn
	s.attachment.terminal = betaAgentTurnTerminal(turn)
}

func (s *AgentSessionStream) latestAttachmentRoot(ctx context.Context) (*Turn, error) {
	service := s.sessions.Turns
	capture, require := agentStreamResponseGuard[pagination.CursorPage[Turn]]()
	service.Options = append([]option.RequestOption{capture}, service.Options...)
	turns := service.ListAutoPaging(ctx, s.sessionID, BetaAgentSessionTurnListParams{Order: "desc"}, append(slices.Clone(s.options), require)...)
	for turns.Next() {
		turn := turns.Current()
		if turn.SubagentID == "" {
			return &turn, nil
		}
	}
	return nil, turns.Err()
}

func (s *AgentSessionStream) observeAttachment(event AgentSessionEventUnion) (bool, error) {
	if s.attachment == nil {
		return true, nil
	}
	if event.Type == "agent.session.idle" {
		if s.turnID != "" && !s.attachment.terminal {
			turn, err := s.attachmentTurn(s.ctx, s.turnID)
			if err != nil {
				return false, err
			}
			if !betaAgentTurnTerminal(turn) {
				return false, nil
			}
			s.selectAttachmentTurn(turn)
		}
		s.attachment.idle = true
	}
	if event.TurnID != "" && event.Item.TurnID != "" && event.TurnID != event.Item.TurnID {
		return false, errors.New("inconsistent attached item turn identity")
	}
	id := event.TurnID
	if event.Item.TurnID != "" {
		id = event.Item.TurnID
	}
	if s.turnID == "" && id != "" {
		// Pending function items may be the first frame, without turn.created.
		turn, err := s.attachmentTurn(s.ctx, id)
		if err != nil {
			return false, err
		}
		if turn.SubagentID == "" {
			if betaAgentTurnTerminal(turn) {
				if turn.ID == s.attachment.baselineID {
					return false, nil
				}
				latest, listErr := s.latestAttachmentRoot(s.ctx)
				if listErr != nil {
					return false, listErr
				}
				if latest == nil || latest.ID != turn.ID {
					return false, nil
				}
			}
			s.selectAttachmentTurn(turn)
		}
	}
	// The root session stream contains root items; successor work belongs to a
	// separate attachment. Never dispatch a call from a different root.
	if id != "" && id != s.turnID {
		return false, nil
	}
	if id == s.turnID && id != "" && (event.Type == "agent.session.turn.completed" || event.Type == "agent.session.turn.failed" || event.Type == "agent.session.turn.cancelled") {
		s.attachment.terminal = true
	}
	return true, nil
}

func (s *AgentSessionStream) seedAttachmentCollector() {
	a := s.attachment
	c := &s.collector
	if a == nil || !c.enabled {
		return
	}
	if c.turn == nil && a.turn != nil {
		c.setTurn(*a.turn)
	}
	if betaAgentTurnTerminal(c.turn) {
		c.terminal = true
	}
}

func (s *AgentSessionStream) attachmentActions(event AgentSessionEventUnion) {
	if s.attachment == nil || !s.collector.enabled {
		return
	}
	if event.Type == "agent.session.requires_action" {
		required := s.collector.required[:0]
		for _, action := range s.collector.required {
			if action.TurnID != "" && action.TurnID != s.turnID {
				continue
			}
			if action.Type == "environment_connection" {
				latest, err := s.latestAttachmentRoot(s.ctx)
				if err != nil {
					s.collector.collectionErr = err
					return
				}
				if latest == nil || latest.ID != s.turnID || latest.Status != "waiting" {
					continue
				}
			}
			required = append(required, action)
		}
		s.collector.required = required
	}
	if event.Type == "agent.session.turn.item.added" && event.Item.Type == "function_call" && s.handlers[event.Item.Name] == nil {
		var action AgentSessionRequiredActionUnion
		if err := json.Unmarshal([]byte(event.Item.JSON.raw), &action); err != nil {
			s.collector.collectionErr = err
			return
		}
		s.collector.required = append(s.collector.required, action)
	}
}

// Only manual actions use the snapshot as a diagnostic. Function execution and
// unhandled-function errors come exclusively from delivered SSE calls.
func (s *AgentSessionStream) attachmentManualActions() error {
	a := s.attachment
	if a == nil || a.turn == nil || a.turn.Status != "waiting" || a.terminal || a.idle {
		return nil
	}
	session, err := s.attachmentSession(a.ctx)
	if err != nil {
		return err
	}
	if session.Status != "requires_action" {
		return nil
	}
	for _, action := range session.RequiredActions {
		manual := action.Type == "computer_use_approval_request" && action.TurnID == s.turnID
		if action.Type == "environment_connection" && (action.TurnID == "" || action.TurnID == s.turnID) {
			latest, listErr := s.latestAttachmentRoot(a.ctx)
			if listErr != nil {
				return listErr
			}
			manual = latest != nil && latest.ID == s.turnID && latest.Status == "waiting"
		}
		if manual {
			s.collector.required = append(s.collector.required, action)
		}
	}
	return nil
}

func (s *AgentSessionStream) reconcileAttachment() error {
	a := s.attachment
	if a == nil || s.turnID == "" {
		return nil
	}
	turn, err := s.attachmentTurn(a.ctx, s.turnID)
	if err != nil {
		return err
	}
	if turn.SubagentID != "" || turn.ID != s.turnID {
		return errors.New("attachment returned a different root turn")
	}
	c := &s.collector
	if !betaAgentTurnTerminal(c.turn) || betaAgentTurnTerminal(turn) {
		c.setTurn(*turn)
	}
	c.terminal = betaAgentTurnTerminal(c.turn)
	if c.terminal {
		c.boundary = true
		c.required = nil
	}
	if !c.terminal {
		return nil
	}
	messages := make(map[int64]AgentSessionMessage)
	service := s.sessions.Items
	capture, require := agentStreamResponseGuard[pagination.CursorPage[AgentSessionItemUnion]]()
	service.Options = append([]option.RequestOption{capture}, service.Options...)
	items := service.ListAutoPaging(a.ctx, s.sessionID, BetaAgentSessionItemListParams{Order: "asc"}, append(slices.Clone(s.options), require)...)
	for items.Next() {
		item := items.Current()
		if item.TurnID != s.turnID || item.Type != "message" || item.Role != "assistant" || item.Status != "completed" || item.Phase == "commentary" {
			continue
		}
		var message AgentSessionMessage
		if err := json.Unmarshal([]byte(item.RawJSON()), &message); err != nil {
			return err
		}
		messages[int64(len(messages))] = message
	}
	if err := items.Err(); err != nil {
		return err
	}
	c.messages = messages
	return nil
}

func (s *AgentSessionStream) attachmentSession(ctx context.Context) (*AgentSession, error) {
	service := *s.sessions
	capture, require := agentStreamResponseGuard[AgentSession]()
	service.Options = append([]option.RequestOption{capture}, service.Options...)
	session, err := service.Get(ctx, s.sessionID, append(slices.Clone(s.options), require)...)
	if err == nil && session == nil {
		return nil, errors.New("attachment received an empty session")
	}
	return session, err
}
func (s *AgentSessionStream) attachmentTurn(ctx context.Context, id string) (*Turn, error) {
	service := s.sessions.Turns
	capture, require := agentStreamResponseGuard[Turn]()
	service.Options = append([]option.RequestOption{capture}, service.Options...)
	turn, err := service.Get(ctx, s.sessionID, id, append(slices.Clone(s.options), require)...)
	if err == nil && (turn == nil || turn.ID != id || turn.SessionID != s.sessionID) {
		return nil, errors.New("attachment received an inconsistent turn")
	}
	return turn, err
}
