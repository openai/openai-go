package main

import (
	"context"
	"os"

	"github.com/openai/openai-go"
)

func main() {

	fileParams := []openai.FileNewParams{}

	if len(os.Args) < 3 || os.Args[1] != "--" {
		panic("usage: go run ./main.go -- <file1> <file2>\n")
	}

	// get files from the command line
	for _, arg := range os.Args[2:] {
		println("File to upload:", arg)
		rdr, err := os.Open(arg)
		defer rdr.Close()
		if err != nil {
			panic("file open failed:" + err.Error())
		}

		fileParams = append(fileParams, openai.FileNewParams{
			File:    rdr,
			Purpose: openai.FilePurposeAssistants,
		})
	}

	println("Creating a new vector store and uploading files")

	ctx := context.Background()
	client := openai.NewClient()

	vectorStore, err := client.Beta.VectorStores.New(
		ctx,
		openai.BetaVectorStoreNewParams{
			ExpiresAfter: openai.BetaVectorStoreNewParamsExpiresAfter{
				Days: openai.Int(1),
			},
			Name: openai.String("Test vector store"),
		},
	)

	if err != nil {
		panic(err)
	}

	// 0 uses default polling interval
	batch, err := client.Beta.VectorStores.FileBatches.UploadAndPoll(ctx, vectorStore.ID, fileParams,
		[]string{}, 0)

	if err != nil {
		panic(err)
	}

	println("Listing the files from the vector store")

	filesCursor, err := client.Beta.VectorStores.FileBatches.ListFiles(ctx, vectorStore.ID, batch.ID,
		openai.BetaVectorStoreFileBatchListFilesParams{})

	if err != nil {
		panic(err)
	}

	for filesCursor != nil {
		for _, f := range filesCursor.Data {
			println("Created file with ID:", f.ID)
		}
		filesCursor, err = filesCursor.GetNextPage()
		if err != nil {
			panic(err)
		}
	}
}
