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

func TestAdminOrganizationProjectSpendAlertNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Admin.Organization.Projects.SpendAlerts.New(
		context.TODO(),
		"project_id",
		openai.AdminOrganizationProjectSpendAlertNewParams{
			Currency: openai.AdminOrganizationProjectSpendAlertNewParamsCurrencyUsd,
			Interval: openai.AdminOrganizationProjectSpendAlertNewParamsIntervalMonth,
			NotificationChannel: openai.AdminOrganizationProjectSpendAlertNewParamsNotificationChannel{
				Recipients:    []string{"string"},
				SubjectPrefix: openai.String("subject_prefix"),
			},
			ThresholdAmount: 0,
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

func TestAdminOrganizationProjectSpendAlertUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Admin.Organization.Projects.SpendAlerts.Update(
		context.TODO(),
		"project_id",
		"alert_id",
		openai.AdminOrganizationProjectSpendAlertUpdateParams{
			Currency: openai.AdminOrganizationProjectSpendAlertUpdateParamsCurrencyUsd,
			Interval: openai.AdminOrganizationProjectSpendAlertUpdateParamsIntervalMonth,
			NotificationChannel: openai.AdminOrganizationProjectSpendAlertUpdateParamsNotificationChannel{
				Recipients:    []string{"string"},
				SubjectPrefix: openai.String("subject_prefix"),
			},
			ThresholdAmount: 0,
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

func TestAdminOrganizationProjectSpendAlertListWithOptionalParams(t *testing.T) {
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
	_, err := client.Admin.Organization.Projects.SpendAlerts.List(
		context.TODO(),
		"project_id",
		openai.AdminOrganizationProjectSpendAlertListParams{
			After:  openai.String("after"),
			Before: openai.String("before"),
			Limit:  openai.Int(0),
			Order:  openai.AdminOrganizationProjectSpendAlertListParamsOrderAsc,
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

func TestAdminOrganizationProjectSpendAlertDelete(t *testing.T) {
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
	_, err := client.Admin.Organization.Projects.SpendAlerts.Delete(
		context.TODO(),
		"project_id",
		"alert_id",
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
