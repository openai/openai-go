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

func TestFineTuningJobNewWithOptionalParams(t *testing.T) {
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
	_, err := client.FineTuning.Jobs.New(context.TODO(), openai.FineTuningJobNewParams{
		Model:        "babbage-002",
		TrainingFile: "file-abc123",
		Hyperparameters: openai.FineTuningJobNewParamsHyperparameters{
			BatchSize: openai.FineTuningJobNewParamsHyperparametersBatchSizeUnion{
				OfAuto: constant.ValueOf[constant.Auto](),
			},
			LearningRateMultiplier: openai.FineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion{
				OfAuto: constant.ValueOf[constant.Auto](),
			},
			NEpochs: openai.FineTuningJobNewParamsHyperparametersNEpochsUnion{
				OfAuto: constant.ValueOf[constant.Auto](),
			},
		},
		Integrations: []openai.FineTuningJobNewParamsIntegration{{
			Wandb: openai.FineTuningJobNewParamsIntegrationWandb{
				Project: "my-wandb-project",
				Entity:  openai.String("entity"),
				Name:    openai.String("name"),
				Tags:    []string{"custom-tag"},
			},
		}},
		Metadata: shared.MetadataParam{
			"foo": "string",
		},
		Method: openai.FineTuningJobNewParamsMethod{
			Dpo: openai.FineTuningJobNewParamsMethodDpo{
				Hyperparameters: openai.FineTuningJobNewParamsMethodDpoHyperparameters{
					BatchSize: openai.FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					Beta: openai.FineTuningJobNewParamsMethodDpoHyperparametersBetaUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					LearningRateMultiplier: openai.FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					NEpochs: openai.FineTuningJobNewParamsMethodDpoHyperparametersNEpochsUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
				},
			},
			Supervised: openai.FineTuningJobNewParamsMethodSupervised{
				Hyperparameters: openai.FineTuningJobNewParamsMethodSupervisedHyperparameters{
					BatchSize: openai.FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					LearningRateMultiplier: openai.FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					NEpochs: openai.FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
				},
			},
			Type: "supervised",
		},
		Seed:           openai.Int(42),
		Suffix:         openai.String("x"),
		ValidationFile: openai.String("file-abc123"),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestFineTuningJobGet(t *testing.T) {
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
	_, err := client.FineTuning.Jobs.Get(context.TODO(), "ft-AF1WoRqd3aJAHsqc9NY7iL8F")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestFineTuningJobListWithOptionalParams(t *testing.T) {
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
	_, err := client.FineTuning.Jobs.List(context.TODO(), openai.FineTuningJobListParams{
		After: openai.String("after"),
		Limit: openai.Int(0),
		Metadata: map[string]string{
			"foo": "string",
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

func TestFineTuningJobCancel(t *testing.T) {
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
	_, err := client.FineTuning.Jobs.Cancel(context.TODO(), "ft-AF1WoRqd3aJAHsqc9NY7iL8F")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestFineTuningJobListEventsWithOptionalParams(t *testing.T) {
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
	_, err := client.FineTuning.Jobs.ListEvents(
		context.TODO(),
		"ft-AF1WoRqd3aJAHsqc9NY7iL8F",
		openai.FineTuningJobListEventsParams{
			After: openai.String("after"),
			Limit: openai.Int(0),
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
