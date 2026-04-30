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

func TestAdminOrganizationProjectRateLimitListRateLimitsWithOptionalParams(t *testing.T) {
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
	_, err := client.Admin.Organization.Projects.RateLimits.ListRateLimits(
		context.TODO(),
		"project_id",
		openai.AdminOrganizationProjectRateLimitListRateLimitsParams{
			After:  openai.String("after"),
			Before: openai.String("before"),
			Limit:  openai.Int(0),
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

func TestAdminOrganizationProjectRateLimitUpdateRateLimitWithOptionalParams(t *testing.T) {
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
	_, err := client.Admin.Organization.Projects.RateLimits.UpdateRateLimit(
		context.TODO(),
		"project_id",
		"rate_limit_id",
		openai.AdminOrganizationProjectRateLimitUpdateRateLimitParams{
			Batch1DayMaxInputTokens:     openai.Int(0),
			MaxAudioMegabytesPer1Minute: openai.Int(0),
			MaxImagesPer1Minute:         openai.Int(0),
			MaxRequestsPer1Day:          openai.Int(0),
			MaxRequestsPer1Minute:       openai.Int(0),
			MaxTokensPer1Minute:         openai.Int(0),
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
