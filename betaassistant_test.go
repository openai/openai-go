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
	"github.com/openai/openai-go/shared/constant"
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
		Model:        shared.ChatModelO3Mini,
		Description:  openai.String("description"),
		Instructions: openai.String("instructions"),
		Metadata: shared.MetadataParam{
			"foo": "string",
		},
		Name:            openai.String("name"),
		ReasoningEffort: shared.ReasoningEffortLow,
		ResponseFormat: openai.AssistantResponseFormatOptionUnionParam{
			OfAuto: constant.ValueOf[constant.Auto](),
		},
		Temperature: openai.Float(1),
		ToolResources: openai.AssistantToolResourcesParam{
			CodeInterpreter: openai.AssistantToolResourcesCodeInterpreterParam{
				FileIDs: []string{"string"},
			},
			FileSearch: openai.AssistantToolResourcesFileSearchParam{
				VectorStoreIDs: []string{"string"},
				VectorStores: []openai.AssistantToolResourcesFileSearchVectorStoreParam{{
					ChunkingStrategy: openai.AssisstantVectorStoreChunkingStrategyUnionParam{
						OfAuto: &openai.AssisstantVectorStoreChunkingStrategyAutoParam{},
					},
					FileIDs: []string{"string"},
					Metadata: shared.MetadataParam{
						"foo": "string",
					},
				}},
			},
		},
		Tools: []openai.AssistantToolUnionParam{{
			OfCodeInterpreter: &openai.CodeInterpreterToolParam{},
		}},
		TopP: openai.Float(1),
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
			Description:  openai.String("description"),
			Instructions: openai.String("instructions"),
			Metadata: shared.MetadataParam{
				"foo": "string",
			},
			Model:           "o3-mini",
			Name:            openai.String("name"),
			ReasoningEffort: shared.ReasoningEffortLow,
			ResponseFormat: openai.AssistantResponseFormatOptionUnionParam{
				OfAuto: constant.ValueOf[constant.Auto](),
			},
			Temperature: openai.Float(1),
			ToolResources: openai.BetaAssistantUpdateParamsToolResources{
				CodeInterpreter: openai.BetaAssistantUpdateParamsToolResourcesCodeInterpreter{
					FileIDs: []string{"string"},
				},
				FileSearch: openai.BetaAssistantUpdateParamsToolResourcesFileSearch{
					VectorStoreIDs: []string{"string"},
				},
			},
			Tools: []openai.AssistantToolUnionParam{{
				OfCodeInterpreter: &openai.CodeInterpreterToolParam{},
			}},
			TopP: openai.Float(1),
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
		After:  openai.String("after"),
		Before: openai.String("before"),
		Limit:  openai.Int(0),
		Order:  openai.BetaAssistantListParamsOrderAsc,
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
