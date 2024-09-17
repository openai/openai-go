package main

import (
	"context"
	"io"
	"os"

	"github.com/openai/openai-go"
)

func main() {

	ctx := context.Background()
	client := openai.NewClient()

	vectorStore, err := client.Beta.VectorStores.New(
		ctx,
		openai.BetaVectorStoreNewParams{
			ExpiresAfter: openai.F(openai.BetaVectorStoreNewParamsExpiresAfter{
				Anchor: openai.F(openai.BetaVectorStoreNewParamsExpiresAfterAnchorLastActiveAt),
				Days:   openai.Int(1),
			}),
			Name: openai.String("Test vector store"),
		},
	)

	if err != nil {
		panic(err.Error())
	}

	chatRdr, err := os.Open("./chat.go")
	defer chatRdr.Close()

	if err != nil {
		panic(err.Error())
	}

	audioRdr, err := os.Open("./audio.go")
	defer audioRdr.Close()

	if err != nil {
		panic(err.Error())
	}

	batch, err := client.Beta.VectorStores.FileBatches.UploadAndPoll(ctx, vectorStore.ID, []openai.FileNewParams{
		{
			Purpose: openai.F(openai.FilePurposeAssistants),
			File:    openai.F[io.Reader](chatRdr),
		},
		{
			Purpose: openai.F(openai.FilePurposeAssistants),
			File:    openai.F[io.Reader](audioRdr),
		},
	}, []string{}, 0)

	if err != nil {
		panic(err.Error())
	}

	filesCursor, err := client.Beta.VectorStores.FileBatches.ListFiles(ctx, vectorStore.ID, batch.ID, openai.BetaVectorStoreFileBatchListFilesParams{})

	if err != nil {
		panic(err.Error())
	}

	for filesCursor != nil {
		for _, f := range filesCursor.Data {
			println("Created file with ID:", f.ID)
		}
		filesCursor, err = filesCursor.GetNextPage()
		if err != nil {
			panic(err.Error())
		}
	}
}
