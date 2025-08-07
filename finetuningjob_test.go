// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/internal/testutil"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/shared"
	"github.com/openai/openai-go/v2/shared/constant"
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
		Model:        openai.FineTuningJobNewParamsModelBabbage002,
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
		Metadata: shared.Metadata{
			"foo": "string",
		},
		Method: openai.FineTuningJobNewParamsMethod{
			Type: "supervised",
			Dpo: openai.DpoMethodParam{
				Hyperparameters: openai.DpoHyperparameters{
					BatchSize: openai.DpoHyperparametersBatchSizeUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					Beta: openai.DpoHyperparametersBetaUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					LearningRateMultiplier: openai.DpoHyperparametersLearningRateMultiplierUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					NEpochs: openai.DpoHyperparametersNEpochsUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
				},
			},
			Reinforcement: openai.ReinforcementMethodParam{
				Grader: openai.ReinforcementMethodGraderUnionParam{
					OfStringCheckGrader: &openai.StringCheckGraderParam{
						Input:     "input",
						Name:      "name",
						Operation: openai.StringCheckGraderOperationEq,
						Reference: "reference",
					},
				},
				Hyperparameters: openai.ReinforcementHyperparameters{
					BatchSize: openai.ReinforcementHyperparametersBatchSizeUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					ComputeMultiplier: openai.ReinforcementHyperparametersComputeMultiplierUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					EvalInterval: openai.ReinforcementHyperparametersEvalIntervalUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					EvalSamples: openai.ReinforcementHyperparametersEvalSamplesUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					LearningRateMultiplier: openai.ReinforcementHyperparametersLearningRateMultiplierUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					NEpochs: openai.ReinforcementHyperparametersNEpochsUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					ReasoningEffort: openai.ReinforcementHyperparametersReasoningEffortDefault,
				},
			},
			Supervised: openai.SupervisedMethodParam{
				Hyperparameters: openai.SupervisedHyperparameters{
					BatchSize: openai.SupervisedHyperparametersBatchSizeUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					LearningRateMultiplier: openai.SupervisedHyperparametersLearningRateMultiplierUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
					NEpochs: openai.SupervisedHyperparametersNEpochsUnion{
						OfAuto: constant.ValueOf[constant.Auto](),
					},
				},
			},
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

func TestFineTuningJobPause(t *testing.T) {
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
	_, err := client.FineTuning.Jobs.Pause(context.TODO(), "ft-AF1WoRqd3aJAHsqc9NY7iL8F")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestFineTuningJobResume(t *testing.T) {
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
	_, err := client.FineTuning.Jobs.Resume(context.TODO(), "ft-AF1WoRqd3aJAHsqc9NY7iL8F")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
