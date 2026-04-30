// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal/testutil"
	"github.com/openai/openai-go/v3/option"
)

func TestAdminOrganizationUsageAudioSpeechesWithOptionalParams(t *testing.T) {
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
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Admin.Organization.Usage.AudioSpeeches(context.TODO(), openai.AdminOrganizationUsageAudioSpeechesParams{
		StartTime:   0,
		APIKeyIDs:   []string{"string"},
		BucketWidth: openai.AdminOrganizationUsageAudioSpeechesParamsBucketWidth1m,
		EndTime:     openai.Int(0),
		GroupBy:     []string{"project_id"},
		Limit:       openai.Int(0),
		Models:      []string{"string"},
		Page:        openai.String("page"),
		ProjectIDs:  []string{"string"},
		UserIDs:     []string{"string"},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAdminOrganizationUsageAudioTranscriptionsWithOptionalParams(t *testing.T) {
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
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Admin.Organization.Usage.AudioTranscriptions(context.TODO(), openai.AdminOrganizationUsageAudioTranscriptionsParams{
		StartTime:   0,
		APIKeyIDs:   []string{"string"},
		BucketWidth: openai.AdminOrganizationUsageAudioTranscriptionsParamsBucketWidth1m,
		EndTime:     openai.Int(0),
		GroupBy:     []string{"project_id"},
		Limit:       openai.Int(0),
		Models:      []string{"string"},
		Page:        openai.String("page"),
		ProjectIDs:  []string{"string"},
		UserIDs:     []string{"string"},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAdminOrganizationUsageCodeInterpreterSessionsWithOptionalParams(t *testing.T) {
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
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Admin.Organization.Usage.CodeInterpreterSessions(context.TODO(), openai.AdminOrganizationUsageCodeInterpreterSessionsParams{
		StartTime:   0,
		BucketWidth: openai.AdminOrganizationUsageCodeInterpreterSessionsParamsBucketWidth1m,
		EndTime:     openai.Int(0),
		GroupBy:     []string{"project_id"},
		Limit:       openai.Int(0),
		Page:        openai.String("page"),
		ProjectIDs:  []string{"string"},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAdminOrganizationUsageCompletionsWithOptionalParams(t *testing.T) {
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
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Admin.Organization.Usage.Completions(context.TODO(), openai.AdminOrganizationUsageCompletionsParams{
		StartTime:   0,
		APIKeyIDs:   []string{"string"},
		Batch:       openai.Bool(true),
		BucketWidth: openai.AdminOrganizationUsageCompletionsParamsBucketWidth1m,
		EndTime:     openai.Int(0),
		GroupBy:     []string{"project_id"},
		Limit:       openai.Int(0),
		Models:      []string{"string"},
		Page:        openai.String("page"),
		ProjectIDs:  []string{"string"},
		UserIDs:     []string{"string"},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAdminOrganizationUsageCostsWithOptionalParams(t *testing.T) {
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
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Admin.Organization.Usage.Costs(context.TODO(), openai.AdminOrganizationUsageCostsParams{
		StartTime:   0,
		APIKeyIDs:   []string{"string"},
		BucketWidth: openai.AdminOrganizationUsageCostsParamsBucketWidth1d,
		EndTime:     openai.Int(0),
		GroupBy:     []string{"project_id"},
		Limit:       openai.Int(0),
		Page:        openai.String("page"),
		ProjectIDs:  []string{"string"},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAdminOrganizationUsageEmbeddingsWithOptionalParams(t *testing.T) {
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
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Admin.Organization.Usage.Embeddings(context.TODO(), openai.AdminOrganizationUsageEmbeddingsParams{
		StartTime:   0,
		APIKeyIDs:   []string{"string"},
		BucketWidth: openai.AdminOrganizationUsageEmbeddingsParamsBucketWidth1m,
		EndTime:     openai.Int(0),
		GroupBy:     []string{"project_id"},
		Limit:       openai.Int(0),
		Models:      []string{"string"},
		Page:        openai.String("page"),
		ProjectIDs:  []string{"string"},
		UserIDs:     []string{"string"},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAdminOrganizationUsageImagesWithOptionalParams(t *testing.T) {
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
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Admin.Organization.Usage.Images(context.TODO(), openai.AdminOrganizationUsageImagesParams{
		StartTime:   0,
		APIKeyIDs:   []string{"string"},
		BucketWidth: openai.AdminOrganizationUsageImagesParamsBucketWidth1m,
		EndTime:     openai.Int(0),
		GroupBy:     []string{"project_id"},
		Limit:       openai.Int(0),
		Models:      []string{"string"},
		Page:        openai.String("page"),
		ProjectIDs:  []string{"string"},
		Sizes:       []string{"256x256"},
		Sources:     []string{"image.generation"},
		UserIDs:     []string{"string"},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAdminOrganizationUsageModerationsWithOptionalParams(t *testing.T) {
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
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Admin.Organization.Usage.Moderations(context.TODO(), openai.AdminOrganizationUsageModerationsParams{
		StartTime:   0,
		APIKeyIDs:   []string{"string"},
		BucketWidth: openai.AdminOrganizationUsageModerationsParamsBucketWidth1m,
		EndTime:     openai.Int(0),
		GroupBy:     []string{"project_id"},
		Limit:       openai.Int(0),
		Models:      []string{"string"},
		Page:        openai.String("page"),
		ProjectIDs:  []string{"string"},
		UserIDs:     []string{"string"},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAdminOrganizationUsageVectorStoresWithOptionalParams(t *testing.T) {
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
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Admin.Organization.Usage.VectorStores(context.TODO(), openai.AdminOrganizationUsageVectorStoresParams{
		StartTime:   0,
		BucketWidth: openai.AdminOrganizationUsageVectorStoresParamsBucketWidth1m,
		EndTime:     openai.Int(0),
		GroupBy:     []string{"project_id"},
		Limit:       openai.Int(0),
		Page:        openai.String("page"),
		ProjectIDs:  []string{"string"},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
