package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/openai/openai-go"
)

func main() {

	// Create an OpenAI client
	client := openai.NewClient()

	ctx := context.Background()

	// Create a Vector Store
	vector_store, err := client.Beta.VectorStores.New(ctx, openai.BetaVectorStoreNewParams{
		Name: openai.String("RAG Demo Vector Store"),
		ExpiresAfter: openai.F(openai.BetaVectorStoreNewParamsExpiresAfter{
			Anchor: openai.F(openai.BetaVectorStoreNewParamsExpiresAfterAnchorLastActiveAt),
			Days:   openai.F(int64(1)),
		}),
	})
	if err != nil {
		panic(err)
	}

	// Create and upload files
	file_contents := []string{
		"Joseph has a pet frog named Milo.",
		"Milo's birthday is on October 7th.",
	}
	for i, content := range file_contents {
		file_name := "doc_" + strconv.Itoa(i) + ".txt"
		doc_params := openai.FileNewParams{
			File:    openai.FileParam(strings.NewReader(content), file_name, "text/plain"),
			Purpose: openai.F(openai.FilePurposeAssistants),
		}
		file_id, err := client.Files.New(ctx, doc_params)
		if err != nil {
			panic(err)
		}
		_, err = client.Beta.VectorStores.Files.New(ctx, vector_store.ID, openai.BetaVectorStoreFileNewParams{
			FileID: openai.F(string(file_id.ID)),
		})
		if err != nil {
			panic(err)
		}
		println(fmt.Sprintf("File added to vector store: %v", file_name))
	}

	// Create an assistant
	const INSTRUCTIONS = `You are a helpful AI bot that answers questions for a user. Keep your response short and direct.
		You may receive a set of context and a question that will relate to the context.
		Do not give information outside the document or repeat your findings.`

	assistant, err := client.Beta.Assistants.New(ctx, openai.BetaAssistantNewParams{
		Name:         openai.String("Frog Buddy"),
		Instructions: openai.String(INSTRUCTIONS),
		Tools: openai.F([]openai.AssistantToolUnionParam{
			openai.FileSearchToolParam{Type: openai.F(openai.FileSearchToolTypeFileSearch)},
		}),
		ToolResources: openai.F(openai.BetaAssistantNewParamsToolResources{
			FileSearch: openai.F(openai.BetaAssistantNewParamsToolResourcesFileSearch{
				VectorStoreIDs: openai.F([]string{vector_store.ID}),
			}),
		}),
		Model: openai.String(openai.ChatModelGPT4o),
	})
	if err != nil {
		panic(err)
	}

	// Create a thread
	thread, err := client.Beta.Threads.New(ctx, openai.BetaThreadNewParams{})
	if err != nil {
		panic(err)
	}

	// Create a message in the thread
	request_message, err := client.Beta.Threads.Messages.New(ctx, thread.ID, openai.BetaThreadMessageNewParams{
		Role: openai.F(openai.BetaThreadMessageNewParamsRoleUser),
		Content: openai.F([]openai.MessageContentPartParamUnion{
			openai.TextContentBlockParam{
				Type: openai.F(openai.TextContentBlockParamTypeText),
				Text: openai.String("When is the birthday of Joseph's pet frog?"),
			},
		}),
	})
	if err != nil {
		panic(err)
	}

	// Create a run
	run, err := client.Beta.Threads.Runs.New(ctx, thread.ID, openai.BetaThreadRunNewParams{
		AssistantID: openai.String(assistant.ID),
		Include:     openai.F([]openai.RunStepInclude{openai.RunStepIncludeStepDetailsToolCallsFileSearchResultsContent}),
	})
	if err != nil {
		panic(err)
	}

	println("Waiting for Run to complete...")
	for run.Status != openai.RunStatusCompleted {
		run, err = client.Beta.Threads.Runs.Get(ctx, thread.ID, run.ID)
		if err != nil {
			panic(err)
		}
		println("Run status: ", run.Status)
		time.Sleep(2 * time.Second)
	}
	println("Run Completed!")

	// List messages
	messages, err := client.Beta.Threads.Messages.List(ctx, thread.ID, openai.BetaThreadMessageListParams{
		RunID: openai.F(string(run.ID)),
	})
	if err != nil {
		panic(err)
	}

	println()
	println(">", request_message.Content[0].Text.Value)
	for messages != nil {
		for _, message := range messages.Data {
			println(message.Content[0].Text.Value)
		}
		messages, err = messages.GetNextPage()
		if err != nil {
			panic(err)
		}
	}
}
