package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

const maxPollInterval = time.Duration(1<<63 - 1)

func mkPollingOptions(pollIntervalMs int) []option.RequestOption {
	options := []option.RequestOption{option.WithHeader("X-Stainless-Poll-Helper", "true")}
	if pollIntervalMs > 0 {
		options = append(options, option.WithHeader("X-Stainless-Poll-Interval", fmt.Sprintf("%d", pollIntervalMs)))
	}
	return options
}

func getPollInterval(raw *http.Response) time.Duration {
	const defaultPollInterval = time.Second
	if raw == nil {
		return defaultPollInterval
	}

	ms, err := strconv.ParseInt(raw.Header.Get("openai-poll-after-ms"), 10, 64)
	if err != nil {
		if errors.Is(err, strconv.ErrRange) && ms > 0 {
			return requestconfig.DefaultMaxServerDelay
		}
		return defaultPollInterval
	}
	if ms <= 0 {
		return defaultPollInterval
	}

	maxServerMilliseconds := int64(requestconfig.DefaultMaxServerDelay / time.Millisecond)
	return time.Duration(min(ms, maxServerMilliseconds)) * time.Millisecond
}

func pollInterval(pollIntervalMs int, raw *http.Response) time.Duration {
	if pollIntervalMs <= 0 {
		return getPollInterval(raw)
	}

	ms := int64(pollIntervalMs)
	if ms > int64(maxPollInterval/time.Millisecond) {
		return maxPollInterval
	}
	return time.Duration(ms) * time.Millisecond
}

// PollStatus waits until a VectorStoreFile is no longer in an incomplete state and returns it.
// Pass 0 as pollIntervalMs to use the default polling interval of 1 second.
// Server-suggested intervals are limited to 8 seconds; pass a positive value to
// explicitly use a longer interval.
func (r *VectorStoreFileService) PollStatus(ctx context.Context, vectorStoreID string, fileID string, pollIntervalMs int, opts ...option.RequestOption) (*VectorStoreFile, error) {
	var raw *http.Response
	var interval time.Duration
	opts = append(opts, mkPollingOptions(pollIntervalMs)...)
	opts = append(opts, option.WithResponseInto(&raw))
	for {
		file, err := r.Get(ctx, vectorStoreID, fileID, opts...)
		if err != nil {
			return nil, fmt.Errorf("vector store file poll: received %w", err)
		}

		switch file.Status {
		case VectorStoreFileStatusInProgress:
			if interval == 0 {
				interval = pollInterval(pollIntervalMs, raw)
			}
			if err := requestconfig.WaitForDelay(ctx, interval); err != nil {
				return nil, err
			}
		case VectorStoreFileStatusCancelled,
			VectorStoreFileStatusCompleted,
			VectorStoreFileStatusFailed:
			return file, nil
		default:
			return nil, fmt.Errorf("invalid vector store file status during polling: received %s", file.Status)
		}
	}
}

// PollStatus waits until a BetaVectorStoreFileBatch is no longer in an incomplete state and returns it.
// Pass 0 as pollIntervalMs to use the default polling interval of 1 second.
// Server-suggested intervals are limited to 8 seconds; pass a positive value to
// explicitly use a longer interval.
func (r *VectorStoreFileBatchService) PollStatus(ctx context.Context, vectorStoreID string, batchID string, pollIntervalMs int, opts ...option.RequestOption) (*VectorStoreFileBatch, error) {
	var raw *http.Response
	var interval time.Duration
	opts = append(opts, option.WithResponseInto(&raw))
	opts = append(opts, mkPollingOptions(pollIntervalMs)...)
	for {
		batch, err := r.Get(ctx, vectorStoreID, batchID, opts...)
		if err != nil {
			return nil, fmt.Errorf("vector store file batch poll: received %w", err)
		}

		switch batch.Status {
		case VectorStoreFileBatchStatusInProgress:
			if interval == 0 {
				interval = pollInterval(pollIntervalMs, raw)
			}
			if err := requestconfig.WaitForDelay(ctx, interval); err != nil {
				return nil, err
			}
		case VectorStoreFileBatchStatusCancelled,
			VectorStoreFileBatchStatusCompleted,
			VectorStoreFileBatchStatusFailed:
			return batch, nil
		default:
			return nil, fmt.Errorf("invalid vector store file batch status during polling: received %s", batch.Status)
		}
	}
}

// PollStatus waits until a VectorStoreFile is no longer in an incomplete state and returns it.
// Pass 0 as pollIntervalMs to use the default polling interval of 1 second.
// Server-suggested intervals are limited to 8 seconds; pass a positive value to
// explicitly use a longer interval.
//
// Deprecated: The Sora API is scheduled to permanently shut down on September 24,
// 2026.
func (r *VideoService) PollStatus(ctx context.Context, videoID string, pollIntervalMs int, opts ...option.RequestOption) (*Video, error) {
	var raw *http.Response
	var interval time.Duration
	opts = append(opts, mkPollingOptions(pollIntervalMs)...)
	opts = append(opts, option.WithResponseInto(&raw))
	for {
		video, err := r.Get(ctx, videoID, opts...)
		if err != nil {
			return nil, fmt.Errorf("error running video poll: received %w", err)
		}

		switch video.Status {
		case VideoStatusQueued, VideoStatusInProgress:
			if interval == 0 {
				interval = pollInterval(pollIntervalMs, raw)
			}
			if err := requestconfig.WaitForDelay(ctx, interval); err != nil {
				return nil, err
			}
		case VideoStatusCompleted,
			VideoStatusFailed:
			return video, nil
		default:
			return nil, fmt.Errorf("invalid video status during polling: received %s", video.Status)
		}
	}
}
