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
		Messages: openai.F([]openai.BetaThreadNewParamsMessage{{
			Content: openai.F([]openai.MessageContentPartParamUnion{openai.ImageFileContentBlockParam{ImageFile: openai.F(openai.ImageFileParam{FileID: openai.F("file_id"), Detail: openai.F(openai.ImageFileDetailAuto)}), Type: openai.F(openai.ImageFileContentBlockTypeImageFile)}}),
			Role:    openai.F(openai.BetaThreadNewParamsMessagesRoleUser),
			Attachments: openai.F([]openai.BetaThreadNewParamsMessagesAttachment{{
				FileID: openai.F("file_id"),
				Tools: openai.F([]openai.BetaThreadNewParamsMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}}),
			}, {
				FileID: openai.F("file_id"),
				Tools: openai.F([]openai.BetaThreadNewParamsMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}}),
			}, {
				FileID: openai.F("file_id"),
				Tools: openai.F([]openai.BetaThreadNewParamsMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}}),
			}}),
			Metadata: openai.F[any](map[string]interface{}{}),
		}, {
			Content: openai.F([]openai.MessageContentPartParamUnion{openai.ImageFileContentBlockParam{ImageFile: openai.F(openai.ImageFileParam{FileID: openai.F("file_id"), Detail: openai.F(openai.ImageFileDetailAuto)}), Type: openai.F(openai.ImageFileContentBlockTypeImageFile)}}),
			Role:    openai.F(openai.BetaThreadNewParamsMessagesRoleUser),
			Attachments: openai.F([]openai.BetaThreadNewParamsMessagesAttachment{{
				FileID: openai.F("file_id"),
				Tools: openai.F([]openai.BetaThreadNewParamsMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}}),
			}, {
				FileID: openai.F("file_id"),
				Tools: openai.F([]openai.BetaThreadNewParamsMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}}),
			}, {
				FileID: openai.F("file_id"),
				Tools: openai.F([]openai.BetaThreadNewParamsMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}}),
			}}),
			Metadata: openai.F[any](map[string]interface{}{}),
		}, {
			Content: openai.F([]openai.MessageContentPartParamUnion{openai.ImageFileContentBlockParam{ImageFile: openai.F(openai.ImageFileParam{FileID: openai.F("file_id"), Detail: openai.F(openai.ImageFileDetailAuto)}), Type: openai.F(openai.ImageFileContentBlockTypeImageFile)}}),
			Role:    openai.F(openai.BetaThreadNewParamsMessagesRoleUser),
			Attachments: openai.F([]openai.BetaThreadNewParamsMessagesAttachment{{
				FileID: openai.F("file_id"),
				Tools: openai.F([]openai.BetaThreadNewParamsMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}}),
			}, {
				FileID: openai.F("file_id"),
				Tools: openai.F([]openai.BetaThreadNewParamsMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}}),
			}, {
				FileID: openai.F("file_id"),
				Tools: openai.F([]openai.BetaThreadNewParamsMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}, openai.CodeInterpreterToolParam{
					Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
				}}),
			}}),
			Metadata: openai.F[any](map[string]interface{}{}),
		}}),
		Metadata: openai.F[any](map[string]interface{}{}),
		ToolResources: openai.F(openai.BetaThreadNewParamsToolResources{
			CodeInterpreter: openai.F(openai.BetaThreadNewParamsToolResourcesCodeInterpreter{
				FileIDs: openai.F([]string{"string", "string", "string"}),
			}),
			FileSearch: openai.F(openai.BetaThreadNewParamsToolResourcesFileSearch{
				VectorStoreIDs: openai.F([]string{"string"}),
				VectorStores: openai.F([]openai.BetaThreadNewParamsToolResourcesFileSearchVectorStore{{
					ChunkingStrategy: openai.F[openai.BetaThreadNewParamsToolResourcesFileSearchVectorStoresChunkingStrategyUnion](openai.BetaThreadNewParamsToolResourcesFileSearchVectorStoresChunkingStrategyAuto{
						Type: openai.F(openai.BetaThreadNewParamsToolResourcesFileSearchVectorStoresChunkingStrategyAutoTypeAuto),
					}),
					FileIDs:  openai.F([]string{"string", "string", "string"}),
					Metadata: openai.F[any](map[string]interface{}{}),
				}}),
			}),
		}),
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
			Metadata: openai.F[any](map[string]interface{}{}),
			ToolResources: openai.F(openai.BetaThreadUpdateParamsToolResources{
				CodeInterpreter: openai.F(openai.BetaThreadUpdateParamsToolResourcesCodeInterpreter{
					FileIDs: openai.F([]string{"string", "string", "string"}),
				}),
				FileSearch: openai.F(openai.BetaThreadUpdateParamsToolResourcesFileSearch{
					VectorStoreIDs: openai.F([]string{"string"}),
				}),
			}),
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
		AssistantID:         openai.F("assistant_id"),
		Instructions:        openai.F("instructions"),
		MaxCompletionTokens: openai.F(int64(256)),
		MaxPromptTokens:     openai.F(int64(256)),
		Metadata:            openai.F[any](map[string]interface{}{}),
		Model:               openai.F(openai.ChatModelGPT4o),
		ParallelToolCalls:   openai.F(true),
		Temperature:         openai.F(1.000000),
		Thread: openai.F(openai.BetaThreadNewAndRunParamsThread{
			Messages: openai.F([]openai.BetaThreadNewAndRunParamsThreadMessage{{
				Content: openai.F([]openai.MessageContentPartParamUnion{openai.ImageFileContentBlockParam{ImageFile: openai.F(openai.ImageFileParam{FileID: openai.F("file_id"), Detail: openai.F(openai.ImageFileDetailAuto)}), Type: openai.F(openai.ImageFileContentBlockTypeImageFile)}}),
				Role:    openai.F(openai.BetaThreadNewAndRunParamsThreadMessagesRoleUser),
				Attachments: openai.F([]openai.BetaThreadNewAndRunParamsThreadMessagesAttachment{{
					FileID: openai.F("file_id"),
					Tools: openai.F([]openai.BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}}),
				}, {
					FileID: openai.F("file_id"),
					Tools: openai.F([]openai.BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}}),
				}, {
					FileID: openai.F("file_id"),
					Tools: openai.F([]openai.BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}}),
				}}),
				Metadata: openai.F[any](map[string]interface{}{}),
			}, {
				Content: openai.F([]openai.MessageContentPartParamUnion{openai.ImageFileContentBlockParam{ImageFile: openai.F(openai.ImageFileParam{FileID: openai.F("file_id"), Detail: openai.F(openai.ImageFileDetailAuto)}), Type: openai.F(openai.ImageFileContentBlockTypeImageFile)}}),
				Role:    openai.F(openai.BetaThreadNewAndRunParamsThreadMessagesRoleUser),
				Attachments: openai.F([]openai.BetaThreadNewAndRunParamsThreadMessagesAttachment{{
					FileID: openai.F("file_id"),
					Tools: openai.F([]openai.BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}}),
				}, {
					FileID: openai.F("file_id"),
					Tools: openai.F([]openai.BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}}),
				}, {
					FileID: openai.F("file_id"),
					Tools: openai.F([]openai.BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}}),
				}}),
				Metadata: openai.F[any](map[string]interface{}{}),
			}, {
				Content: openai.F([]openai.MessageContentPartParamUnion{openai.ImageFileContentBlockParam{ImageFile: openai.F(openai.ImageFileParam{FileID: openai.F("file_id"), Detail: openai.F(openai.ImageFileDetailAuto)}), Type: openai.F(openai.ImageFileContentBlockTypeImageFile)}}),
				Role:    openai.F(openai.BetaThreadNewAndRunParamsThreadMessagesRoleUser),
				Attachments: openai.F([]openai.BetaThreadNewAndRunParamsThreadMessagesAttachment{{
					FileID: openai.F("file_id"),
					Tools: openai.F([]openai.BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}}),
				}, {
					FileID: openai.F("file_id"),
					Tools: openai.F([]openai.BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}}),
				}, {
					FileID: openai.F("file_id"),
					Tools: openai.F([]openai.BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}, openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}}),
				}}),
				Metadata: openai.F[any](map[string]interface{}{}),
			}}),
			Metadata: openai.F[any](map[string]interface{}{}),
			ToolResources: openai.F(openai.BetaThreadNewAndRunParamsThreadToolResources{
				CodeInterpreter: openai.F(openai.BetaThreadNewAndRunParamsThreadToolResourcesCodeInterpreter{
					FileIDs: openai.F([]string{"string", "string", "string"}),
				}),
				FileSearch: openai.F(openai.BetaThreadNewAndRunParamsThreadToolResourcesFileSearch{
					VectorStoreIDs: openai.F([]string{"string"}),
					VectorStores: openai.F([]openai.BetaThreadNewAndRunParamsThreadToolResourcesFileSearchVectorStore{{
						ChunkingStrategy: openai.F[openai.BetaThreadNewAndRunParamsThreadToolResourcesFileSearchVectorStoresChunkingStrategyUnion](openai.BetaThreadNewAndRunParamsThreadToolResourcesFileSearchVectorStoresChunkingStrategyAuto{
							Type: openai.F(openai.BetaThreadNewAndRunParamsThreadToolResourcesFileSearchVectorStoresChunkingStrategyAutoTypeAuto),
						}),
						FileIDs:  openai.F([]string{"string", "string", "string"}),
						Metadata: openai.F[any](map[string]interface{}{}),
					}}),
				}),
			}),
		}),
		ToolChoice: openai.F[openai.AssistantToolChoiceOptionUnionParam](openai.AssistantToolChoiceOptionString(openai.AssistantToolChoiceOptionStringNone)),
		ToolResources: openai.F(openai.BetaThreadNewAndRunParamsToolResources{
			CodeInterpreter: openai.F(openai.BetaThreadNewAndRunParamsToolResourcesCodeInterpreter{
				FileIDs: openai.F([]string{"string", "string", "string"}),
			}),
			FileSearch: openai.F(openai.BetaThreadNewAndRunParamsToolResourcesFileSearch{
				VectorStoreIDs: openai.F([]string{"string"}),
			}),
		}),
		Tools: openai.F([]openai.BetaThreadNewAndRunParamsToolUnion{openai.CodeInterpreterToolParam{
			Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
		}, openai.CodeInterpreterToolParam{
			Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
		}, openai.CodeInterpreterToolParam{
			Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
		}}),
		TopP: openai.F(1.000000),
		TruncationStrategy: openai.F(openai.BetaThreadNewAndRunParamsTruncationStrategy{
			Type:         openai.F(openai.BetaThreadNewAndRunParamsTruncationStrategyTypeAuto),
			LastMessages: openai.F(int64(1)),
		}),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
