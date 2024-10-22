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
		Model:        openai.F(openai.FineTuningJobNewParamsModelBabbage002),
		TrainingFile: openai.F("file-abc123"),
		Hyperparameters: openai.F(openai.FineTuningJobNewParamsHyperparameters{
			BatchSize:              openai.F[openai.FineTuningJobNewParamsHyperparametersBatchSizeUnion](openai.FineTuningJobNewParamsHyperparametersBatchSizeBehavior(openai.FineTuningJobNewParamsHyperparametersBatchSizeBehaviorAuto)),
			LearningRateMultiplier: openai.F[openai.FineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion](openai.FineTuningJobNewParamsHyperparametersLearningRateMultiplierBehavior(openai.FineTuningJobNewParamsHyperparametersLearningRateMultiplierBehaviorAuto)),
			NEpochs:                openai.F[openai.FineTuningJobNewParamsHyperparametersNEpochsUnion](openai.FineTuningJobNewParamsHyperparametersNEpochsBehavior(openai.FineTuningJobNewParamsHyperparametersNEpochsBehaviorAuto)),
		}),
		Integrations: openai.F([]openai.FineTuningJobNewParamsIntegration{{
			Type: openai.F(openai.FineTuningJobNewParamsIntegrationsTypeWandb),
			Wandb: openai.F(openai.FineTuningJobNewParamsIntegrationsWandb{
				Project: openai.F("my-wandb-project"),
				Entity:  openai.F("entity"),
				Name:    openai.F("name"),
				Tags:    openai.F([]string{"custom-tag", "custom-tag", "custom-tag"}),
			}),
		}, {
			Type: openai.F(openai.FineTuningJobNewParamsIntegrationsTypeWandb),
			Wandb: openai.F(openai.FineTuningJobNewParamsIntegrationsWandb{
				Project: openai.F("my-wandb-project"),
				Entity:  openai.F("entity"),
				Name:    openai.F("name"),
				Tags:    openai.F([]string{"custom-tag", "custom-tag", "custom-tag"}),
			}),
		}, {
			Type: openai.F(openai.FineTuningJobNewParamsIntegrationsTypeWandb),
			Wandb: openai.F(openai.FineTuningJobNewParamsIntegrationsWandb{
				Project: openai.F("my-wandb-project"),
				Entity:  openai.F("entity"),
				Name:    openai.F("name"),
				Tags:    openai.F([]string{"custom-tag", "custom-tag", "custom-tag"}),
			}),
		}}),
		Seed:           openai.F(int64(42)),
		Suffix:         openai.F("x"),
		ValidationFile: openai.F("file-abc123"),
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
		After: openai.F("after"),
		Limit: openai.F(int64(0)),
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
			After: openai.F("after"),
			Limit: openai.F(int64(0)),
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
