package openai

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"maps"
	"net/http"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

// AgentToolHandler handles a function call sequentially during stream iteration.
// Arguments are a private copy of the call's JSON object. Return a string,
// map[string]any, []InputContentParamUnion, or nil. Errors and invalid outputs
// produce a generic failed tool result without exposing application error text.
type AgentToolHandler func(context.Context, map[string]any) (any, error)

// AgentSessionStreamParams configures one turn on an existing idle session.
type AgentSessionStreamParams struct {
	// Input must be a nonempty string or []AgentSessionInputMessageParam.
	Input        any
	ToolHandlers map[string]AgentToolHandler
	// IdempotencyKey applies only to input. A request header overrides this value.
	// When absent, the helper generates a unique key. Each tool result gets its own key.
	IdempotencyKey param.Opt[string]
}

// AgentSessionStream streams one turn's original typed events. Use Next, Current,
// and Err as with ssestream.Stream, and defer Close if iteration may stop early.
// It is not safe to iterate concurrently. Close may be called concurrently to
// unblock iteration; it closes local resources without cancelling the backend turn.
type AgentSessionStream struct {
	ctx          context.Context
	cancel       context.CancelFunc
	sessions     *BetaAgentSessionService
	sessionID    string
	options      []option.RequestOption
	handlers     map[string]AgentToolHandler
	stream       *ssestream.Stream[AgentSessionEventUnion]
	current      AgentSessionEventUnion
	err          error
	closed       atomic.Bool
	closeOnce    sync.Once
	closeErr     error
	turnID       string
	turnEnded    bool
	recent       [1024]string
	recentCount  int
	recentNext   int
	eventIDs     map[string]struct{}
	handledCalls map[agentCallKey]struct{}
	pending      *agentPendingCall
}

type agentCallKey struct{ turnID, callID string }

// Capture dispatch state before exposing the mutable event to the caller.
type agentPendingCall struct {
	turnID, callID string
	handler        AgentToolHandler
	arguments      map[string]any
	argumentErr    error
}

// Stream subscribes to events before submitting input, then streams through the
// selected coordinator turn's completion/failure/cancellation and subsequent
// session idle event, or a session failure. The session must be idle, with only
// one input writer while the helper runs; concurrent inputs cannot be correlated
// to turns. Initial idle and subagent completion events do not end the stream.
// Protocol error frames follow the underlying stream's Err path and close the
// connection before later events are read. Typed turn/session terminal events are
// yielded when reached. Unknown functions remain available for manual handling through Events.New.
// Registered handlers run on the next Next call after their event is returned.
// WithResponseBodyInto is incompatible with the helper's required typed decodes
// and is rejected through Err; other request options retain their normal behavior.
func (r *BetaAgentSessionService) Stream(ctx context.Context, sessionID string, params AgentSessionStreamParams, opts ...option.RequestOption) *AgentSessionStream {
	ctx, cancel := context.WithCancel(ctx)
	s := &AgentSessionStream{ctx: ctx, cancel: cancel, sessions: r, sessionID: sessionID,
		options: slices.Clone(opts), handlers: maps.Clone(params.ToolHandlers),
		eventIDs: make(map[string]struct{}), handledCalls: make(map[agentCallKey]struct{})}
	input, err := agentStreamInput(params.Input)
	if err != nil {
		s.finish(err)
		return s
	}
	probe := *r
	capture, require := agentStreamResponseGuard[AgentSession]()
	probe.Options = append([]option.RequestOption{capture}, r.Options...)
	session, err := probe.Get(ctx, sessionID, append(slices.Clone(opts), require)...)
	if err != nil {
		s.finish(err)
		return s
	}
	if session == nil {
		s.finish(errors.New("sessions.Stream received an empty session response"))
		return s
	}
	if session.Status != "idle" {
		s.finish(errors.New("sessions.Stream requires an idle session; use sessions.Events.StreamStreaming to follow an active session"))
		return s
	}
	events := r.Events
	capture, require = agentStreamResponseGuard[http.Response]()
	events.Options = append([]option.RequestOption{capture}, r.Events.Options...)
	s.stream = events.StreamStreaming(ctx, sessionID, append(slices.Clone(opts), require)...)
	if err := s.stream.Err(); err != nil {
		s.finish(err)
		return s
	}
	key := agentStreamKey()
	inputOpts := slices.Clone(opts)
	if !param.IsOmitted(params.IdempotencyKey) {
		inputOpts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", params.IdempotencyKey.Value)}, inputOpts...)
	}
	// Resolve header overrides once per request, case-insensitively, using the
	// transport's normal options. A default key survives all transport retries.
	inputOpts = append(inputOpts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
		if cfg.Request.Header.Get("Idempotency-Key") == "" {
			cfg.SetHeader("Idempotency-Key", key)
		}
		return nil
	}))
	if err := r.Events.New(ctx, sessionID, BetaAgentSessionEventNewParams{Events: []AgentSessionInputParamUnion{input}}, inputOpts...); err != nil {
		s.finish(err)
	}
	return s
}

func agentStreamInput(input any) (AgentSessionInputParamUnion, error) {
	var messages []AgentSessionInputMessageParam
	switch value := input.(type) {
	case string:
		if value != "" {
			messages = []AgentSessionInputMessageParam{{Content: []InputContentParamUnion{InputContentParamOfParamInputText(value)}}}
		}
	case []AgentSessionInputMessageParam:
		messages = value
	default:
		return AgentSessionInputParamUnion{}, errors.New("input must be a string or []AgentSessionInputMessageParam")
	}
	if len(messages) == 0 {
		return AgentSessionInputParamUnion{}, errors.New("input must not be empty")
	}
	return AgentSessionInputParamOfParamAgentSessionInputMessage(messages), nil
}

// Next advances the stream. Unexpected EOF is reported through Err.
func (s *AgentSessionStream) Next() (ok bool) {
	if s.closed.Load() {
		return false
	}
	// Also release the connection when application handler code panics.
	defer func() {
		if !ok {
			_ = s.Close()
		}
	}()
	if err := s.ctx.Err(); err != nil {
		return s.finish(err)
	}
	if s.pending != nil {
		call := *s.pending
		s.pending = nil
		if err := s.handle(call); err != nil {
			return s.finish(err)
		}
	}
	for s.stream.Next() {
		if s.closed.Load() {
			return false
		}
		event := s.stream.Current()
		if !s.accept(event) {
			continue
		}
		s.current = event
		if event.Type == "agent.session.failed" || (event.Type == "agent.session.idle" && s.turnEnded) {
			_ = s.Close()
		} else if event.Type == "agent.session.turn.item.added" && event.Item.Type == "function_call" {
			key := agentCallKey{event.Item.TurnID, event.Item.CallID}
			if _, exists := s.handledCalls[key]; !exists {
				s.handledCalls[key] = struct{}{}
				if handler := s.handlers[event.Item.Name]; handler != nil {
					arguments, argumentErr := agentToolArguments(event.Item.Arguments)
					s.pending = &agentPendingCall{turnID: event.Item.TurnID, callID: event.Item.CallID,
						handler: handler, arguments: arguments, argumentErr: argumentErr}
				}
			}
		}
		return true
	}
	if s.closed.Load() {
		return false
	}
	if err := s.stream.Err(); err != nil {
		return s.finish(err)
	}
	return s.finish(io.ErrUnexpectedEOF)
}

func (s *AgentSessionStream) accept(event AgentSessionEventUnion) bool {
	if _, exists := s.eventIDs[event.EventID]; exists {
		return false
	}
	if s.recentCount == len(s.recent) {
		delete(s.eventIDs, s.recent[s.recentNext])
	} else {
		s.recentCount++
	}
	s.recent[s.recentNext] = event.EventID
	s.recentNext = (s.recentNext + 1) % len(s.recent)
	s.eventIDs[event.EventID] = struct{}{}
	if event.Type == "agent.session.turn.created" && event.Turn.SubagentID == "" && s.turnID == "" {
		s.turnID = event.TurnID
	}
	switch event.Type {
	case "agent.session.turn.completed", "agent.session.turn.failed", "agent.session.turn.cancelled":
		if s.turnID != "" && event.TurnID == s.turnID {
			s.turnEnded = true
		}
	}
	return true
}

// Current returns the original typed event at the current position.
func (s *AgentSessionStream) Current() AgentSessionEventUnion { return s.current }

// Err returns the first request, stream, or tool-result submission error.
func (s *AgentSessionStream) Err() error { return s.err }

func (s *AgentSessionStream) finish(err error) bool { s.err = err; _ = s.Close(); return false }

// Close releases local resources without sending backend cancellation.
func (s *AgentSessionStream) Close() error {
	s.closeOnce.Do(func() {
		s.closed.Store(true)
		s.cancel()
		if s.stream != nil {
			s.closeErr = s.stream.Close()
		}
	})
	return s.closeErr
}

func agentStreamKey() string {
	var key [16]byte
	_, _ = rand.Read(key[:]) // crypto/rand.Read always fills the buffer or terminates the process.
	return hex.EncodeToString(key[:])
}

func (s *AgentSessionStream) handle(call agentPendingCall) error {
	result := agentToolResult(s.ctx, call)
	if err := s.ctx.Err(); err != nil {
		return err
	}
	opts := append(slices.Clone(s.options), option.WithHeader("Idempotency-Key", agentStreamKey()))
	delays := [...]time.Duration{100 * time.Millisecond, 300 * time.Millisecond, 600 * time.Millisecond}
	for attempt := 0; ; attempt++ {
		err := s.sessions.Events.New(s.ctx, s.sessionID, BetaAgentSessionEventNewParams{Events: []AgentSessionInputParamUnion{{OfParamAgentSessionInputToolResult: &result}}}, opts...)
		if err == nil {
			return nil
		}
		var apiErr *Error
		if attempt == len(delays) || !errors.As(err, &apiErr) || apiErr.StatusCode != 400 || apiErr.Code != "invalid_request_error" || apiErr.Message != "Unknown pending tool call: "+call.callID {
			return err
		}
		timer := time.NewTimer(delays[attempt])
		select {
		case <-s.ctx.Done():
			timer.Stop()
			return s.ctx.Err()
		case <-timer.C:
		}
	}
}

// Observe the generated method's destination before inherited/caller options,
// then reject replacements before transport. Keep the original options intact:
// wrapping them would lose the transport's pre-request option semantics.
func agentStreamResponseGuard[T any]() (option.RequestOption, option.RequestOption) {
	var destination **T
	capture := requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
		destination, _ = cfg.ResponseBodyInto.(**T)
		return nil
	})
	require := requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
		actual, ok := cfg.ResponseBodyInto.(**T)
		if !ok || actual != destination {
			return errors.New("sessions.Stream does not support WithResponseBodyInto; typed session and event decoding is required")
		}
		return nil
	})
	return capture, require
}
