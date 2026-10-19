package openai

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"sync"

	"github.com/openai/openai-go/v3/option"
)

// Create a vector store file batch and polls the API until the task is complete.
// Pass 0 for pollIntervalMs to enable default polling interval.
func (r *VectorStoreFileBatchService) NewAndPoll(ctx context.Context, vectorStoreId string, body VectorStoreFileBatchNewParams, pollIntervalMs int, opts ...option.RequestOption) (res *VectorStoreFileBatch, err error) {
	return newVectorStoreFileBatchAndPoll(r, ctx, vectorStoreId, body, pollIntervalMs, opts...)
}

// Uploads the given files concurrently and then creates a vector store file batch.
//
// If you've already uploaded certain files that you want to include in this batch
// then you can pass their IDs through the file_ids argument.
//
// Pass 0 for pollIntervalMs to enable default polling interval.
//
// By default, if any file upload fails then an exception will be eagerly raised.
func (r *VectorStoreFileBatchService) UploadAndPoll(ctx context.Context, vectorStoreID string, files []FileNewParams, fileIDs []string, pollIntervalMs int, opts ...option.RequestOption) (*VectorStoreFileBatch, error) {
	return uploadVectorStoreFileBatchAndPoll(r, ctx, vectorStoreID, files, fileIDs, pollIntervalMs, opts...)
}

const defaultVectorStoreFileBatchUploadConcurrency = 16

func newVectorStoreFileBatchAndPoll(r *VectorStoreFileBatchService, ctx context.Context, vectorStoreId string, body VectorStoreFileBatchNewParams, pollIntervalMs int, opts ...option.RequestOption) (res *VectorStoreFileBatch, err error) {
	batch, err := r.New(ctx, vectorStoreId, body, opts...)
	if err != nil {
		return nil, err
	}
	return r.PollStatus(ctx, vectorStoreId, batch.ID, pollIntervalMs, opts...)
}

func uploadVectorStoreFileBatchAndPoll(r *VectorStoreFileBatchService, ctx context.Context, vectorStoreID string, files []FileNewParams, fileIDs []string, pollIntervalMs int, opts ...option.RequestOption) (*VectorStoreFileBatch, error) {
	maxConcurrency := r.MaxUploadConcurrency
	if maxConcurrency < 0 {
		return nil, errors.New("vector store file batch: MaxUploadConcurrency must not be negative")
	}
	if maxConcurrency == 0 {
		maxConcurrency = defaultVectorStoreFileBatchUploadConcurrency
	}
	if len(files) <= 0 {
		return nil, errors.New("No `files` provided to process. If you've already uploaded files you should use `.NewAndPoll()` instead")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}
	uploadCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	filesService := NewFileService(r.Options...)

	uploadedFileIDs := make(chan string, len(files))
	var firstError error
	var failureOnce sync.Once
	wg := sync.WaitGroup{}

	// Bound active uploads, including multipart bodies buffered for retries.
	var next int
	var mu sync.Mutex
	for range min(maxConcurrency, len(files)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				mu.Lock()
				if next == len(files) || uploadCtx.Err() != nil {
					mu.Unlock()
					return
				}
				file := files[next]
				next++
				mu.Unlock()
				if file.File != nil {
					file.File = vectorBatchReader{uploadCtx, file.File}
				}
				fileObj, err := filesService.New(uploadCtx, file, opts...)
				if err != nil {
					failureOnce.Do(func() {
						firstError = err
						cancel()
					})
					return
				}
				uploadedFileIDs <- fileObj.ID
			}
		}()
	}

	wg.Wait()
	close(uploadedFileIDs)

	if firstError != nil {
		return nil, firstError
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	for id := range uploadedFileIDs {
		fileIDs = append(fileIDs, id)
	}

	return r.NewAndPoll(ctx, vectorStoreID, VectorStoreFileBatchNewParams{
		FileIDs: fileIDs,
	}, pollIntervalMs, opts...)
}

// Check cancellation between reads without taking ownership of caller readers.
// An already-blocked Read still needs to be released by its owner.
type vectorBatchReader struct {
	ctx context.Context
	io.Reader
}

func (r vectorBatchReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.Reader.Read(p)
}

func (r vectorBatchReader) Filename() string {
	if named, ok := r.Reader.(interface{ Filename() string }); ok {
		return named.Filename()
	}
	if named, ok := r.Reader.(interface{ Name() string }); ok {
		return filepath.Base(named.Name())
	}
	return "anonymous_file"
}

func (r vectorBatchReader) ContentType() string {
	if typed, ok := r.Reader.(interface{ ContentType() string }); ok {
		return typed.ContentType()
	}
	return "application/octet-stream"
}
