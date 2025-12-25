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
)

func TestFineTuningAlphaGraderRunWithOptionalParams(t *testing.T) {
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
	_, err := client.FineTuning.Alpha.Graders.Run(context.TODO(), openai.FineTuningAlphaGraderRunParams{
		Grader: openai.FineTuningAlphaGraderRunParamsGraderUnion{
			OfStringCheck: &openai.StringCheckGraderParam{
				Input:     "input",
				Name:      "name",
				Operation: openai.StringCheckGraderOperationEq,
				Reference: "reference",
			},
		},
		ModelSample: "model_sample",
		Item:        map[string]any{},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestFineTuningAlphaGraderValidateWithOptionalParams(t *testing.T) {
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
	_, err := client.FineTuning.Alpha.Graders.Validate(context.TODO(), openai.FineTuningAlphaGraderValidateParams{
		Grader: openai.FineTuningAlphaGraderValidateParamsGraderUnion{
			OfStringCheckGrader: &openai.StringCheckGraderParam{
				Input:     "input",
				Name:      "name",
				Operation: openai.StringCheckGraderOperationEq,
				Reference: "reference",
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
