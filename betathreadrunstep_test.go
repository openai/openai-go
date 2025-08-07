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

func TestBetaThreadRunStepGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Threads.Runs.Steps.Get(
		context.TODO(),
		"step_id",
		openai.BetaThreadRunStepGetParams{
			ThreadID: "thread_id",
			RunID:    "run_id",
			Include:  []openai.RunStepInclude{openai.RunStepIncludeStepDetailsToolCallsFileSearchResultsContent},
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

func TestBetaThreadRunStepListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Threads.Runs.Steps.List(
		context.TODO(),
		"run_id",
		openai.BetaThreadRunStepListParams{
			ThreadID: "thread_id",
			After:    openai.String("after"),
			Before:   openai.String("before"),
			Include:  []openai.RunStepInclude{openai.RunStepIncludeStepDetailsToolCallsFileSearchResultsContent},
			Limit:    openai.Int(0),
			Order:    openai.BetaThreadRunStepListParamsOrderAsc,
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
