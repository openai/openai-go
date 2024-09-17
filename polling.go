package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/openai/openai-go/option"
)

func pollingOptions(pollIntervalMs int, raw *http.Response) (headers []option.RequestOption) {
	headers = []option.RequestOption{
		option.WithHeader("X-Stainless-Poll-Helper", "true"),
		option.WithResponseInto(&raw),
	}
	if pollIntervalMs > 0 {
		headers = append(headers, option.WithHeader("X-Stainless-Poll-Interval", fmt.Sprintf("%d", pollIntervalMs)))
	}
	return headers
}

func sleepInterval(pollIntervalMs int, raw *http.Response) {
	if pollIntervalMs <= 0 {
		if ms, err := strconv.Atoi(raw.Header.Get("openai-poll-after-ms")); err == nil {
			pollIntervalMs = ms
		}
		pollIntervalMs = 1000
	}
	time.Sleep(time.Duration(pollIntervalMs) * time.Millisecond)
}

func (r *BetaVectorStoreFileService) PollStatus(ctx context.Context, vectorStoreID string, fileID string, pollIntervalMs int, opts ...option.RequestOption) (*VectorStoreFile, error) {
	var raw *http.Response
	opts = append(opts, pollingOptions(pollIntervalMs, raw)...)
	for true {
		file, err := r.Get(ctx, vectorStoreID, fileID, opts...)
		if err != nil {
			return nil, err
		}
		switch file.Status {
		case VectorStoreFileStatusInProgress:
			sleepInterval(pollIntervalMs, raw)
		case VectorStoreFileStatusCancelled:
		case VectorStoreFileStatusCompleted:
		case VectorStoreFileStatusFailed:
			return file, nil
		default:
			break
		}
	}
	return nil, errors.New("Invalid vector store file status during polling")
}

func (r *BetaVectorStoreFileBatchService) PollStatus(ctx context.Context, vectorStoreID string, batchID string, pollIntervalMs int, opts ...option.RequestOption) (*VectorStoreFileBatch, error) {
	var raw *http.Response
	opts = append(opts, pollingOptions(pollIntervalMs, raw)...)
	for true {
		batch, err := r.Get(ctx, vectorStoreID, batchID, opts...)
		if err != nil {
			return nil, err
		}

		switch batch.Status {
		case VectorStoreFileBatchStatusInProgress:
			sleepInterval(pollIntervalMs, raw)
		case VectorStoreFileBatchStatusCancelled:
		case VectorStoreFileBatchStatusCompleted:
		case VectorStoreFileBatchStatusFailed:
			return batch, nil
		default:
			break
		}
	}
	return nil, errors.New("Invalid vector store file batch status during polling")
}

func (r *BetaThreadRunService) PollStatus(ctx context.Context, threadID string, runID string, pollIntervalMs int, opts ...option.RequestOption) (res *Run, err error) {
	var raw *http.Response
	opts = append(opts, pollingOptions(pollIntervalMs, raw)...)
	for true {
		run, err := r.Get(ctx, threadID, runID, opts...)
		if err != nil {
			return nil, err
		}

		switch run.Status {
		case RunStatusRequiresAction:
		case RunStatusCancelled:
		case RunStatusCompleted:
		case RunStatusFailed:
		case RunStatusExpired:
		case RunStatusIncomplete:
			return run, nil
		case RunStatusInProgress:
		case RunStatusQueued:
			sleepInterval(pollIntervalMs, raw)
		default:
			break
		}
	}
	return nil, errors.New("Invalid vector store file batch status during polling")
}
