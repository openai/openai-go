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

func TestBetaAssistantNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Assistants.New(context.TODO(), openai.BetaAssistantNewParams{
		Model:        openai.F(openai.ChatModelGPT4o),
		Description:  openai.F("description"),
		Instructions: openai.F("instructions"),
		Metadata:     openai.F[any](map[string]interface{}{}),
		Name:         openai.F("name"),
		Temperature:  openai.F(1.000000),
		ToolResources: openai.F(openai.BetaAssistantNewParamsToolResources{
			CodeInterpreter: openai.F(openai.BetaAssistantNewParamsToolResourcesCodeInterpreter{
				FileIDs: openai.F([]string{"string", "string", "string"}),
			}),
			FileSearch: openai.F(openai.BetaAssistantNewParamsToolResourcesFileSearch{
				VectorStoreIDs: openai.F([]string{"string"}),
				VectorStores: openai.F([]openai.BetaAssistantNewParamsToolResourcesFileSearchVectorStore{{
					FileIDs: openai.F([]string{"string", "string", "string"}),
					ChunkingStrategy: openai.F[openai.BetaAssistantNewParamsToolResourcesFileSearchVectorStoresChunkingStrategyUnion](openai.BetaAssistantNewParamsToolResourcesFileSearchVectorStoresChunkingStrategyAuto{
						Type: openai.F(openai.BetaAssistantNewParamsToolResourcesFileSearchVectorStoresChunkingStrategyAutoTypeAuto),
					}),
					Metadata: openai.F[any](map[string]interface{}{}),
				}}),
			}),
		}),
		Tools: openai.F([]openai.AssistantToolUnionParam{openai.CodeInterpreterToolParam{
			Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
		}, openai.CodeInterpreterToolParam{
			Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
		}, openai.CodeInterpreterToolParam{
			Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
		}}),
		TopP: openai.F(1.000000),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAssistantGet(t *testing.T) {
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
	_, err := client.Beta.Assistants.Get(context.TODO(), "assistant_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAssistantUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Assistants.Update(
		context.TODO(),
		"assistant_id",
		openai.BetaAssistantUpdateParams{
			Description:  openai.F("description"),
			Instructions: openai.F("instructions"),
			Metadata:     openai.F[any](map[string]interface{}{}),
			Model:        openai.F("model"),
			Name:         openai.F("name"),
			Temperature:  openai.F(1.000000),
			ToolResources: openai.F(openai.BetaAssistantUpdateParamsToolResources{
				CodeInterpreter: openai.F(openai.BetaAssistantUpdateParamsToolResourcesCodeInterpreter{
					FileIDs: openai.F([]string{"string", "string", "string"}),
				}),
				FileSearch: openai.F(openai.BetaAssistantUpdateParamsToolResourcesFileSearch{
					VectorStoreIDs: openai.F([]string{"string"}),
				}),
			}),
			Tools: openai.F([]openai.AssistantToolUnionParam{openai.CodeInterpreterToolParam{
				Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
			}, openai.CodeInterpreterToolParam{
				Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
			}, openai.CodeInterpreterToolParam{
				Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
			}}),
			TopP: openai.F(1.000000),
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

func TestBetaAssistantListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Assistants.List(context.TODO(), openai.BetaAssistantListParams{
		After:  openai.F("after"),
		Before: openai.F("before"),
		Limit:  openai.F(int64(0)),
		Order:  openai.F(openai.BetaAssistantListParamsOrderAsc),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAssistantDelete(t *testing.T) {
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
	_, err := client.Beta.Assistants.Delete(context.TODO(), "assistant_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
