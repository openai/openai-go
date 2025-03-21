// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"os"
	"testing"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/internal/testutil"
	"github.com/openai/openai-go/option"
)

func TestAutoPagination(t *testing.T) {
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
	iter := client.FineTuning.Jobs.ListAutoPaging(context.TODO(), openai.FineTuningJobListParams{
		Limit: openai.Int(20),
	})
	// Prism mock isn't going to give us real pagination
	for i := 0; i < 3 && iter.Next(); i++ {
		job := iter.Current()
		t.Logf("%+v\n", job.ID)
	}
	if err := iter.Err(); err != nil {
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
