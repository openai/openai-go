package openai

import (
	"context"
	"maps"
	"slices"

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
