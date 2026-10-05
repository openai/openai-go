package openai

import (
	"context"
	"errors"
	"sync"

	"github.com/openai/openai-go/v3/option"
)

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

	filesService := NewFileService(r.Options...)

	uploadedFileIDs := make(chan string, len(files))
	fileUploadErrors := make(chan error, len(files))
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
				if next == len(files) {
					mu.Unlock()
					return
				}
				file := files[next]
				next++
				mu.Unlock()
				if err := ctx.Err(); err != nil {
					fileUploadErrors <- err
					return
				}
				fileObj, err := filesService.New(ctx, file, opts...)
				if err != nil {
					fileUploadErrors <- err
					continue
				}
				uploadedFileIDs <- fileObj.ID
			}
		}()
	}

	wg.Wait()
	close(uploadedFileIDs)
	close(fileUploadErrors)

	for err := range fileUploadErrors {
		return nil, err
	}

	for id := range uploadedFileIDs {
		fileIDs = append(fileIDs, id)
	}

	return r.NewAndPoll(ctx, vectorStoreID, VectorStoreFileBatchNewParams{
		FileIDs: fileIDs,
	}, pollIntervalMs, opts...)
}
