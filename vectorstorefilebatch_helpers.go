package openai

import (
	"context"
	"errors"
	"sync"

	"github.com/openai/openai-go/v3/option"
)

func newVectorStoreFileBatchAndPoll(r *VectorStoreFileBatchService, ctx context.Context, vectorStoreId string, body VectorStoreFileBatchNewParams, pollIntervalMs int, opts ...option.RequestOption) (res *VectorStoreFileBatch, err error) {
	batch, err := r.New(ctx, vectorStoreId, body, opts...)
	if err != nil {
		return nil, err
	}
	return r.PollStatus(ctx, vectorStoreId, batch.ID, pollIntervalMs, opts...)
}

func uploadVectorStoreFileBatchAndPoll(r *VectorStoreFileBatchService, ctx context.Context, vectorStoreID string, files []FileNewParams, fileIDs []string, pollIntervalMs int, opts ...option.RequestOption) (*VectorStoreFileBatch, error) {
	if len(files) <= 0 {
		return nil, errors.New("No `files` provided to process. If you've already uploaded files you should use `.NewAndPoll()` instead")
	}

	filesService := NewFileService(r.Options...)

	uploadedFileIDs := make(chan string, len(files))
	fileUploadErrors := make(chan error, len(files))
	wg := sync.WaitGroup{}

	for _, file := range files {
		wg.Add(1)
		go func(file FileNewParams) {
			defer wg.Done()
			fileObj, err := filesService.New(ctx, file, opts...)
			if err != nil {
				fileUploadErrors <- err
				return
			}
			uploadedFileIDs <- fileObj.ID
		}(file)
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
