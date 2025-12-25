// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Nordlys-Labs/openai-go/v3"
	"github.com/Nordlys-Labs/openai-go/v3/internal/testutil"
	"github.com/Nordlys-Labs/openai-go/v3/option"
	"github.com/Nordlys-Labs/openai-go/v3/shared"
	"github.com/Nordlys-Labs/openai-go/v3/shared/constant"
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
			Content: openai.BetaThreadNewParamsMessageContentUnion{
				OfString: openai.String("string"),
			},
			Role: "user",
			Attachments: []openai.BetaThreadNewParamsMessageAttachment{{
				FileID: openai.String("file_id"),
				Tools: []openai.BetaThreadNewParamsMessageAttachmentToolUnion{{
					OfCodeInterpreter: &openai.CodeInterpreterToolParam{},
				}},
			}},
			Metadata: shared.Metadata{
				"foo": "string",
			},
		}},
		Metadata: shared.Metadata{
			"foo": "string",
		},
		ToolResources: openai.BetaThreadNewParamsToolResources{
			CodeInterpreter: openai.BetaThreadNewParamsToolResourcesCodeInterpreter{
				FileIDs: []string{"string"},
			},
			FileSearch: openai.BetaThreadNewParamsToolResourcesFileSearch{
				VectorStoreIDs: []string{"string"},
				VectorStores: []openai.BetaThreadNewParamsToolResourcesFileSearchVectorStore{{
					ChunkingStrategy: openai.BetaThreadNewParamsToolResourcesFileSearchVectorStoreChunkingStrategyUnion{
						OfAuto: &openai.BetaThreadNewParamsToolResourcesFileSearchVectorStoreChunkingStrategyAuto{},
					},
					FileIDs: []string{"string"},
					Metadata: shared.Metadata{
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
			Metadata: shared.Metadata{
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
		AssistantID:         "assistant_id",
		Instructions:        openai.String("instructions"),
		MaxCompletionTokens: openai.Int(256),
		MaxPromptTokens:     openai.Int(256),
		Metadata: shared.Metadata{
			"foo": "string",
		},
		Model:             shared.ChatModelGPT5_2,
		ParallelToolCalls: openai.Bool(true),
		ResponseFormat: openai.AssistantResponseFormatOptionUnionParam{
			OfAuto: constant.ValueOf[constant.Auto](),
		},
		Temperature: openai.Float(1),
		Thread: openai.BetaThreadNewAndRunParamsThread{
			Messages: []openai.BetaThreadNewAndRunParamsThreadMessage{{
				Content: openai.BetaThreadNewAndRunParamsThreadMessageContentUnion{
					OfString: openai.String("string"),
				},
				Role: "user",
				Attachments: []openai.BetaThreadNewAndRunParamsThreadMessageAttachment{{
					FileID: openai.String("file_id"),
					Tools: []openai.BetaThreadNewAndRunParamsThreadMessageAttachmentToolUnion{{
						OfCodeInterpreter: &openai.CodeInterpreterToolParam{},
					}},
				}},
				Metadata: shared.Metadata{
					"foo": "string",
				},
			}},
			Metadata: shared.Metadata{
				"foo": "string",
			},
			ToolResources: openai.BetaThreadNewAndRunParamsThreadToolResources{
				CodeInterpreter: openai.BetaThreadNewAndRunParamsThreadToolResourcesCodeInterpreter{
					FileIDs: []string{"string"},
				},
				FileSearch: openai.BetaThreadNewAndRunParamsThreadToolResourcesFileSearch{
					VectorStoreIDs: []string{"string"},
					VectorStores: []openai.BetaThreadNewAndRunParamsThreadToolResourcesFileSearchVectorStore{{
						ChunkingStrategy: openai.BetaThreadNewAndRunParamsThreadToolResourcesFileSearchVectorStoreChunkingStrategyUnion{
							OfAuto: &openai.BetaThreadNewAndRunParamsThreadToolResourcesFileSearchVectorStoreChunkingStrategyAuto{},
						},
						FileIDs: []string{"string"},
						Metadata: shared.Metadata{
							"foo": "string",
						},
					}},
				},
			},
		},
		ToolChoice: openai.AssistantToolChoiceOptionUnionParam{
			OfAuto: openai.String("none"),
		},
		ToolResources: openai.BetaThreadNewAndRunParamsToolResources{
			CodeInterpreter: openai.BetaThreadNewAndRunParamsToolResourcesCodeInterpreter{
				FileIDs: []string{"string"},
			},
			FileSearch: openai.BetaThreadNewAndRunParamsToolResourcesFileSearch{
				VectorStoreIDs: []string{"string"},
			},
		},
		Tools: []openai.AssistantToolUnionParam{{
			OfCodeInterpreter: &openai.CodeInterpreterToolParam{},
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
