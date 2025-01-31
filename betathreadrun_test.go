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

func TestBetaThreadRunNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Threads.Runs.New(
		context.TODO(),
		"thread_id",
		openai.BetaThreadRunNewParams{
			AssistantID:            openai.F("assistant_id"),
			Include:                openai.F([]openai.RunStepInclude{openai.RunStepIncludeStepDetailsToolCallsFileSearchResultsContent}),
			AdditionalInstructions: openai.F("additional_instructions"),
			AdditionalMessages: openai.F([]openai.BetaThreadRunNewParamsAdditionalMessage{{
				Content: openai.F([]openai.MessageContentPartParamUnion{openai.ImageFileContentBlockParam{ImageFile: openai.F(openai.ImageFileParam{FileID: openai.F("file_id"), Detail: openai.F(openai.ImageFileDetailAuto)}), Type: openai.F(openai.ImageFileContentBlockTypeImageFile)}}),
				Role:    openai.F(openai.BetaThreadRunNewParamsAdditionalMessagesRoleUser),
				Attachments: openai.F([]openai.BetaThreadRunNewParamsAdditionalMessagesAttachment{{
					FileID: openai.F("file_id"),
					Tools: openai.F([]openai.BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolUnion{openai.CodeInterpreterToolParam{
						Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
					}}),
				}}),
				Metadata: openai.F(shared.MetadataParam{
					"foo": "string",
				}),
			}}),
			Instructions:        openai.F("instructions"),
			MaxCompletionTokens: openai.F(int64(256)),
			MaxPromptTokens:     openai.F(int64(256)),
			Metadata: openai.F(shared.MetadataParam{
				"foo": "string",
			}),
			Model:             openai.F(openai.ChatModelO3Mini),
			ParallelToolCalls: openai.F(true),
			Temperature:       openai.F(1.000000),
			ToolChoice:        openai.F[openai.AssistantToolChoiceOptionUnionParam](openai.AssistantToolChoiceOptionAuto(openai.AssistantToolChoiceOptionAutoNone)),
			Tools: openai.F([]openai.AssistantToolUnionParam{openai.CodeInterpreterToolParam{
				Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter),
			}}),
			TopP: openai.F(1.000000),
			TruncationStrategy: openai.F(openai.BetaThreadRunNewParamsTruncationStrategy{
				Type:         openai.F(openai.BetaThreadRunNewParamsTruncationStrategyTypeAuto),
				LastMessages: openai.F(int64(1)),
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

func TestBetaThreadRunGet(t *testing.T) {
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
	_, err := client.Beta.Threads.Runs.Get(
		context.TODO(),
		"thread_id",
		"run_id",
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaThreadRunUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Threads.Runs.Update(
		context.TODO(),
		"thread_id",
		"run_id",
		openai.BetaThreadRunUpdateParams{
			Metadata: openai.F(shared.MetadataParam{
				"foo": "string",
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

func TestBetaThreadRunListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Threads.Runs.List(
		context.TODO(),
		"thread_id",
		openai.BetaThreadRunListParams{
			After:  openai.F("after"),
			Before: openai.F("before"),
			Limit:  openai.F(int64(0)),
			Order:  openai.F(openai.BetaThreadRunListParamsOrderAsc),
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

func TestBetaThreadRunCancel(t *testing.T) {
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
	_, err := client.Beta.Threads.Runs.Cancel(
		context.TODO(),
		"thread_id",
		"run_id",
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaThreadRunSubmitToolOutputsWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Threads.Runs.SubmitToolOutputs(
		context.TODO(),
		"thread_id",
		"run_id",
		openai.BetaThreadRunSubmitToolOutputsParams{
			ToolOutputs: openai.F([]openai.BetaThreadRunSubmitToolOutputsParamsToolOutput{{
				Output:     openai.F("output"),
				ToolCallID: openai.F("tool_call_id"),
			}}),
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
