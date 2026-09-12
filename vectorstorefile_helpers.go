package openai

import (
	"context"

	"github.com/openai/openai-go/v3/option"
)

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
