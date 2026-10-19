package openai

import (
	"context"

	"github.com/openai/openai-go/v3/option"
)

// Create Video and Poll for Completion
//
// Polls the API and blocks until the task is complete.
// Default polling interval is 1 second.
//
// Deprecated: The Sora API is scheduled to permanently shut down on September 24,
// 2026.
func (r *VideoService) NewAndPoll(ctx context.Context, body VideoNewParams, pollIntervalMs int, opts ...option.RequestOption) (res *Video, err error) {
	video, err := r.New(ctx, body, opts...)
	if err != nil {
		return nil, err
	}
	return r.PollStatus(ctx, video.ID, pollIntervalMs, opts...)
}
