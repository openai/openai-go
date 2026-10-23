package openai

import (
	"context"
	"errors"
	"maps"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

func betaAgentCreationStream(ctx context.Context, cancel context.CancelFunc, sessions *BetaAgentSessionService, stream *ssestream.Stream[AgentSessionEventUnion], handlers map[string]AgentToolHandler, opts []option.RequestOption) *ssestream.Stream[AgentSessionEventUnion] {
	s := &AgentSessionStream{ctx: ctx, cancel: cancel, sessions: sessions, stream: stream,
		options: slices.Clone(opts), handlers: maps.Clone(handlers),
		allowEOF: true,
		eventIDs: make(map[string]struct{}), handledCalls: make(map[agentCallKey]struct{})}
	s.collector.handlers = s.handlers
	if err := stream.Err(); err != nil {
		s.finish(err)
	}
	return ssestream.NewStreamWithBetaIterator[AgentSessionEventUnion](s, &s.collector)
}

func betaAgentSessionNew(r *BetaAgentSessionService, ctx context.Context, body BetaAgentSessionNewParams, opts ...option.RequestOption) (res *AgentSession, err error) {
	if len(body.ToolHandlers) != 0 {
		return nil, errors.New("ToolHandlers requires NewStreaming")
	}
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	path := "agents/sessions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

func betaAgentSessionNewStreaming(r *BetaAgentSessionService, ctx context.Context, body BetaAgentSessionNewParams, opts ...option.RequestOption) (stream *ssestream.Stream[AgentSessionEventUnion]) {
	toolOpts := opts
	var cancel context.CancelFunc
	if len(body.ToolHandlers) != 0 {
		ctx, cancel = context.WithCancel(ctx)
	}
	var (
		raw *http.Response
		err error
	)
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	if cancel != nil {
		capture, require := agentStreamResponseGuard[http.Response]()
		preClientOpts = append(preClientOpts, capture)
		opts = append(slices.Clone(opts), require)
	}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	opts = append(opts, option.WithJSONSet("stream", true))
	path := "agents/sessions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &raw, opts...)
	stream = ssestream.NewStreamWithBetaAccumulator[AgentSessionEventUnion](ssestream.NewDecoder(raw), err, &betaAgentTurnCollector{})
	if cancel != nil {
		return betaAgentCreationStream(ctx, cancel, r, stream, body.ToolHandlers, toolOpts)
	}
	return stream
}
