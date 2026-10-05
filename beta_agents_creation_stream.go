package openai

import (
	"context"
	"io"
	"maps"
	"slices"

	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

func betaAgentCreationStream(ctx context.Context, cancel context.CancelFunc, sessions *BetaAgentSessionService, stream *ssestream.Stream[AgentSessionEventUnion], handlers map[string]AgentToolHandler, opts []option.RequestOption) *ssestream.Stream[AgentSessionEventUnion] {
	// Keep transport options, but creation body overrides and response capture
	// must not replace the generated tool-result request or its response destination.
	toolSessions := *sessions
	var body io.Reader
	var contentType string
	var responseBody any
	capture := requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
		body, contentType, responseBody = cfg.Body, cfg.Request.Header.Get("Content-Type"), cfg.ResponseBodyInto
		return nil
	})
	restore := requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
		cfg.Body, body = body, nil
		cfg.ResponseBodyInto, responseBody = responseBody, nil
		cfg.Request.Header.Set("Content-Type", contentType)
		cfg.ResponseInto = nil
		return nil
	})
	inherited, eventOptions := requestconfig.SplitInheritedOptions(sessions.Events.Options, sessions.eventDefaults)
	_, sessionOptions := requestconfig.SplitInheritedOptions(sessions.Options, sessions.Options)
	toolSessions.Events.Options = append([]option.RequestOption{capture}, inherited...)
	s := &AgentSessionStream{ctx: ctx, cancel: cancel, sessions: &toolSessions, stream: stream,
		options: slices.Concat(sessionOptions, opts, []option.RequestOption{restore}, eventOptions), handlers: maps.Clone(handlers),
		allowEOF: true,
		eventIDs: make(map[string]struct{}), handledCalls: make(map[agentCallKey]struct{})}
	s.collector.handlers = s.handlers
	if err := stream.Err(); err != nil {
		s.finish(err)
	}
	return ssestream.NewStreamWithBetaIterator[AgentSessionEventUnion](s, &s.collector)
}
