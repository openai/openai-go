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
			AssistantID:            "assistant_id",
			Include:                []openai.RunStepInclude{openai.RunStepIncludeStepDetailsToolCallsFileSearchResultsContent},
			AdditionalInstructions: openai.String("additional_instructions"),
			AdditionalMessages: []openai.BetaThreadRunNewParamsAdditionalMessage{{
				Content: openai.BetaThreadRunNewParamsAdditionalMessageContentUnion{
					OfString: openai.String("string"),
				},
				Role: "user",
				Attachments: []openai.BetaThreadRunNewParamsAdditionalMessageAttachment{{
					FileID: openai.String("file_id"),
					Tools: []openai.BetaThreadRunNewParamsAdditionalMessageAttachmentToolUnion{{
						OfCodeInterpreter: &openai.CodeInterpreterToolParam{},
					}},
				}},
				Metadata: shared.Metadata{
					"foo": "string",
				},
			}},
			Instructions:        openai.String("instructions"),
			MaxCompletionTokens: openai.Int(256),
			MaxPromptTokens:     openai.Int(256),
			Metadata: shared.Metadata{
				"foo": "string",
			},
			Model:             shared.ChatModelGPT5_2,
			ParallelToolCalls: openai.Bool(true),
			ReasoningEffort:   shared.ReasoningEffortNone,
			ResponseFormat: openai.AssistantResponseFormatOptionUnionParam{
				OfAuto: constant.ValueOf[constant.Auto](),
			},
			Temperature: openai.Float(1),
			ToolChoice: openai.AssistantToolChoiceOptionUnionParam{
				OfAuto: openai.String("none"),
			},
			Tools: []openai.AssistantToolUnionParam{{
				OfCodeInterpreter: &openai.CodeInterpreterToolParam{},
			}},
			TopP: openai.Float(1),
			TruncationStrategy: openai.BetaThreadRunNewParamsTruncationStrategy{
				Type:         "auto",
				LastMessages: openai.Int(1),
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
			Metadata: shared.Metadata{
				"foo": "string",
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
			After:  openai.String("after"),
			Before: openai.String("before"),
			Limit:  openai.Int(0),
			Order:  openai.BetaThreadRunListParamsOrderAsc,
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
			ToolOutputs: []openai.BetaThreadRunSubmitToolOutputsParamsToolOutput{{
				Output:     openai.String("output"),
				ToolCallID: openai.String("tool_call_id"),
			}},
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
