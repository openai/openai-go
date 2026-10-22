package openai

import (
	"context"

	"github.com/openai/openai-go/v3/option"
)

// Create a vector store file by attaching a
// [File](https://platform.openai.com/docs/api-reference/files) to a
// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object).
//
// Polls the API and blocks until the task is complete.
// Default polling interval is 1 second.
func (r *VectorStoreFileService) NewAndPoll(ctx context.Context, vectorStoreId string, body VectorStoreFileNewParams, pollIntervalMs int, opts ...option.RequestOption) (res *VectorStoreFile, err error) {
	return newVectorStoreFileAndPoll(r, ctx, vectorStoreId, body, pollIntervalMs, opts...)
}

// Upload a file to the `files` API and then attach it to the given vector store.
//
// Note the file will be asynchronously processed (you can use the alternative
// polling helper method to wait for processing to complete).
func (r *VectorStoreFileService) Upload(ctx context.Context, vectorStoreID string, body FileNewParams, opts ...option.RequestOption) (*VectorStoreFile, error) {
	return uploadVectorStoreFile(r, ctx, vectorStoreID, body, opts...)
}

// Add a file to a vector store and poll until processing is complete.
// Default polling interval is 1 second.
func (r *VectorStoreFileService) UploadAndPoll(ctx context.Context, vectorStoreID string, body FileNewParams, pollIntervalMs int, opts ...option.RequestOption) (*VectorStoreFile, error) {
	return uploadVectorStoreFileAndPoll(r, ctx, vectorStoreID, body, pollIntervalMs, opts...)
}

func newVectorStoreFileAndPoll(r *VectorStoreFileService, ctx context.Context, vectorStoreId string, body VectorStoreFileNewParams, pollIntervalMs int, opts ...option.RequestOption) (res *VectorStoreFile, err error) {
	file, err := r.New(ctx, vectorStoreId, body, opts...)
	if err != nil {
		return nil, err
	}
	return r.PollStatus(ctx, vectorStoreId, file.ID, pollIntervalMs, opts...)
}

func uploadVectorStoreFile(r *VectorStoreFileService, ctx context.Context, vectorStoreID string, body FileNewParams, opts ...option.RequestOption) (*VectorStoreFile, error) {
	filesService := NewFileService(r.Options...)
	fileObj, err := filesService.New(ctx, body, opts...)
	if err != nil {
		return nil, err
	}
	return r.New(ctx, vectorStoreID, VectorStoreFileNewParams{
		FileID: fileObj.ID,
	}, opts...)
}

func uploadVectorStoreFileAndPoll(r *VectorStoreFileService, ctx context.Context, vectorStoreID string, body FileNewParams, pollIntervalMs int, opts ...option.RequestOption) (*VectorStoreFile, error) {
	res, err := r.Upload(ctx, vectorStoreID, body, opts...)
	if err != nil {
		return nil, err
	}
	return r.PollStatus(ctx, vectorStoreID, res.ID, pollIntervalMs, opts...)
}
