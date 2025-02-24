// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/internal/testutil"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

func TestBetaThreadNewWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Beta.Threads.New(context.TODO(), openai.BetaThreadNewParams{
		Messages: []openai.BetaThreadNewParamsMessage{{
			Content: []openai.MessageContentPartParamUnion{{
				OfImageFile: &openai.ImageFileContentBlockParam{ImageFile: openai.ImageFileParam{FileID: openai.String("file_id"), Detail: "auto"}},
			}},
			Role: "user",
			Attachments: []openai.BetaThreadNewParamsMessagesAttachment{{
				FileID: openai.String("file_id"),
				Tools: []openai.BetaThreadNewParamsMessagesAttachmentsToolUnion{{
					OfCodeInterpreter: &openai.CodeInterpreterToolParam{},
				}},
			}},
			Metadata: shared.MetadataParam{
				"foo": "string",
			},
		}},
		Metadata: shared.MetadataParam{
			"foo": "string",
		},
		ToolResources: openai.BetaThreadNewParamsToolResources{
			CodeInterpreter: openai.BetaThreadNewParamsToolResourcesCodeInterpreter{
				FileIDs: []string{"string"},
			},
			FileSearch: openai.BetaThreadNewParamsToolResourcesFileSearch{
				VectorStoreIDs: []string{"string"},
				VectorStores: []openai.BetaThreadNewParamsToolResourcesFileSearchVectorStore{{
					ChunkingStrategy: openai.FileChunkingStrategyParamUnion{
						OfAuto: &openai.AutoFileChunkingStrategyParam{},
					},
					FileIDs: []string{"string"},
					Metadata: shared.MetadataParam{
						"foo": "string",
					},
				}},
			},
		},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaThreadGet(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Beta.Threads.Get(context.TODO(), "thread_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaThreadUpdateWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Beta.Threads.Update(
		context.TODO(),
		"thread_id",
		openai.BetaThreadUpdateParams{
			Metadata: shared.MetadataParam{
				"foo": "string",
			},
			ToolResources: openai.BetaThreadUpdateParamsToolResources{
				CodeInterpreter: openai.BetaThreadUpdateParamsToolResourcesCodeInterpreter{
					FileIDs: []string{"string"},
				},
				FileSearch: openai.BetaThreadUpdateParamsToolResourcesFileSearch{
					VectorStoreIDs: []string{"string"},
				},
			},
		},
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaThreadDelete(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Beta.Threads.Delete(context.TODO(), "thread_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaThreadNewAndRunWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Beta.Threads.NewAndRun(context.TODO(), openai.BetaThreadNewAndRunParams{
		AssistantID:         openai.String("assistant_id"),
		Instructions:        openai.String("instructions"),
		MaxCompletionTokens: openai.Int(256),
		MaxPromptTokens:     openai.Int(256),
		Metadata: shared.MetadataParam{
			"foo": "string",
		},
		Model:             openai.ChatModelO3Mini,
		ParallelToolCalls: openai.Bool(true),
		Temperature:       openai.Float(1),
		Thread: openai.BetaThreadNewAndRunParamsThread{
			Messages: []openai.BetaThreadNewAndRunParamsThreadMessage{{
				Content: []openai.MessageContentPartParamUnion{{
					OfImageFile: &openai.ImageFileContentBlockParam{ImageFile: openai.ImageFileParam{FileID: openai.String("file_id"), Detail: "auto"}},
				}},
				Role: "user",
				Attachments: []openai.BetaThreadNewAndRunParamsThreadMessagesAttachment{{
					FileID: openai.String("file_id"),
					Tools: []openai.BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion{{
						OfCodeInterpreter: &openai.CodeInterpreterToolParam{},
					}},
				}},
				Metadata: shared.MetadataParam{
					"foo": "string",
				},
			}},
			Metadata: shared.MetadataParam{
				"foo": "string",
			},
			ToolResources: openai.BetaThreadNewAndRunParamsThreadToolResources{
				CodeInterpreter: openai.BetaThreadNewAndRunParamsThreadToolResourcesCodeInterpreter{
					FileIDs: []string{"string"},
				},
				FileSearch: openai.BetaThreadNewAndRunParamsThreadToolResourcesFileSearch{
					VectorStoreIDs: []string{"string"},
					VectorStores: []openai.BetaThreadNewAndRunParamsThreadToolResourcesFileSearchVectorStore{{
						ChunkingStrategy: openai.FileChunkingStrategyParamUnion{
							OfAuto: &openai.AutoFileChunkingStrategyParam{},
						},
						FileIDs: []string{"string"},
						Metadata: shared.MetadataParam{
							"foo": "string",
						},
					}},
				},
			},
		},
		ToolChoice: openai.AssistantToolChoiceOptionUnionParam{
			OfAuto: "none",
		},
		ToolResources: openai.BetaThreadNewAndRunParamsToolResources{
			CodeInterpreter: openai.BetaThreadNewAndRunParamsToolResourcesCodeInterpreter{
				FileIDs: []string{"string"},
			},
			FileSearch: openai.BetaThreadNewAndRunParamsToolResourcesFileSearch{
				VectorStoreIDs: []string{"string"},
			},
		},
		Tools: []openai.BetaThreadNewAndRunParamsToolUnion{{
			OfCodeInterpreterTool: &openai.CodeInterpreterToolParam{},
		}},
		TopP: openai.Float(1),
		TruncationStrategy: openai.BetaThreadNewAndRunParamsTruncationStrategy{
			Type:         "auto",
			LastMessages: openai.Int(1),
		},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
