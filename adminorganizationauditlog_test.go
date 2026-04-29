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

func TestAdminOrganizationAuditLogListWithOptionalParams(t *testing.T) {
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
	_, err := client.Admin.Organization.AuditLogs.List(context.TODO(), openai.AdminOrganizationAuditLogListParams{
		ActorEmails: []string{"string"},
		ActorIDs:    []string{"string"},
		After:       openai.String("after"),
		Before:      openai.String("before"),
		EffectiveAt: openai.AdminOrganizationAuditLogListParamsEffectiveAt{
			Gt:  openai.Int(0),
			Gte: openai.Int(0),
			Lt:  openai.Int(0),
			Lte: openai.Int(0),
		},
		EventTypes:  []string{"api_key.created"},
		Limit:       openai.Int(0),
		ProjectIDs:  []string{"string"},
		ResourceIDs: []string{"string"},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
