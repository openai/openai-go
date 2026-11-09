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

func TestAdminOrganizationWorkloadIdentityProviderNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Admin.Organization.WorkloadIdentityProviders.New(context.TODO(), openai.AdminOrganizationWorkloadIdentityProviderNewParams{
		Audience: "x",
		Issuer:   "x",
		Name:     "x",
		Type:     "oidc",
		AttributeTransformations: []openai.AdminOrganizationWorkloadIdentityProviderNewParamsAttributeTransformation{{
			Attribute:  "x",
			Expression: "x",
		}},
		Description: openai.String("description"),
		Enabled:     openai.Bool(true),
		Verification: openai.AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUnion{
			OfDiscovery: &openai.AdminOrganizationWorkloadIdentityProviderNewParamsVerificationDiscovery{
				DiscoveryURL: openai.String("discovery_url"),
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

func TestAdminOrganizationWorkloadIdentityProviderGet(t *testing.T) {
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
	_, err := client.Admin.Organization.WorkloadIdentityProviders.Get(context.TODO(), "x")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAdminOrganizationWorkloadIdentityProviderUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Admin.Organization.WorkloadIdentityProviders.Update(
		context.TODO(),
		"x",
		openai.AdminOrganizationWorkloadIdentityProviderUpdateParams{
			AttributeTransformations: []openai.AdminOrganizationWorkloadIdentityProviderUpdateParamsAttributeTransformation{{
				Attribute:  "x",
				Expression: "x",
			}},
			Audience:    openai.String("x"),
			Description: openai.String("description"),
			Enabled:     openai.Bool(true),
			Issuer:      openai.String("x"),
			Name:        openai.String("x"),
			Verification: openai.AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUnion{
				OfAdminOrganizationWorkloadIdentityProviderUpdatesVerificationDiscoveryVerificationParam: &openai.AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationDiscoveryVerificationParam{
					DiscoveryURL: openai.String("discovery_url"),
				},
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

func TestAdminOrganizationWorkloadIdentityProviderListWithOptionalParams(t *testing.T) {
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
	_, err := client.Admin.Organization.WorkloadIdentityProviders.List(context.TODO(), openai.AdminOrganizationWorkloadIdentityProviderListParams{
		After: openai.String("x"),
		Limit: openai.Int(1),
		Order: openai.AdminOrganizationWorkloadIdentityProviderListParamsOrderStringAsc,
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAdminOrganizationWorkloadIdentityProviderDelete(t *testing.T) {
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
	_, err := client.Admin.Organization.WorkloadIdentityProviders.Delete(context.TODO(), "x")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
