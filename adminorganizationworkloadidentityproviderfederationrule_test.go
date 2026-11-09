// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

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

func TestAdminOrganizationWorkloadIdentityProviderFederationRuleNewWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithUnsafeAllowHTTP(),
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Admin.Organization.WorkloadIdentityProviders.FederationRules.New(
		context.TODO(),
		"x",
		openai.AdminOrganizationWorkloadIdentityProviderFederationRuleNewParams{
			AttributeMatchers: map[string]openai.AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsAttributeMatcherUnion{
				"foo": {
					OfString: openai.String("string"),
				},
			},
			ProjectID:        "x",
			ServiceAccountID: "x",
			Description:      openai.String("description"),
			Enabled:          openai.Bool(true),
			Name:             openai.String("name"),
			Permissions: openai.AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsUnion{
				OfServiceAccountDefault: &openai.AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsServiceAccountDefault{},
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

func TestAdminOrganizationWorkloadIdentityProviderFederationRuleGet(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithUnsafeAllowHTTP(),
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Admin.Organization.WorkloadIdentityProviders.FederationRules.Get(
		context.TODO(),
		"x",
		"x",
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAdminOrganizationWorkloadIdentityProviderFederationRuleUpdateWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithUnsafeAllowHTTP(),
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Admin.Organization.WorkloadIdentityProviders.FederationRules.Update(
		context.TODO(),
		"x",
		"x",
		openai.AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParams{
			AttributeMatchers: map[string]openai.AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsAttributeMatcherUnion{
				"foo": {
					OfString: openai.String("string"),
				},
			},
			Description: openai.String("description"),
			Enabled:     openai.Bool(true),
			Name:        openai.String("name"),
			Permissions: openai.AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsUnion{
				OfServiceAccountDefault: &openai.AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsDefaultPermissionsParam{},
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

func TestAdminOrganizationWorkloadIdentityProviderFederationRuleListWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithUnsafeAllowHTTP(),
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Admin.Organization.WorkloadIdentityProviders.FederationRules.List(
		context.TODO(),
		"x",
		openai.AdminOrganizationWorkloadIdentityProviderFederationRuleListParams{
			After:            openai.String("x"),
			Limit:            openai.Int(1),
			Order:            openai.AdminOrganizationWorkloadIdentityProviderFederationRuleListParamsOrderStringAsc,
			ServiceAccountID: openai.String("x"),
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

func TestAdminOrganizationWorkloadIdentityProviderFederationRuleDelete(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithUnsafeAllowHTTP(),
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
	)
	_, err := client.Admin.Organization.WorkloadIdentityProviders.FederationRules.Delete(
		context.TODO(),
		"x",
		"x",
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
