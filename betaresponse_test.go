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

func TestBetaResponseNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Responses.New(context.TODO(), openai.BetaResponseNewParams{
		Background: openai.Bool(true),
		ContextManagement: []openai.BetaResponseNewParamsContextManagement{{
			Type:             "type",
			CompactThreshold: openai.Int(1000),
		}},
		Conversation: openai.BetaResponseNewParamsConversationUnion{
			OfString: openai.String("string"),
		},
		Include: []openai.BetaResponseIncludable{openai.BetaResponseIncludableFileSearchCallResults},
		Input: openai.BetaResponseNewParamsInputUnion{
			OfString: openai.String("string"),
		},
		Instructions:    openai.String("instructions"),
		MaxOutputTokens: openai.Int(16),
		MaxToolCalls:    openai.Int(0),
		Metadata: map[string]string{
			"foo": "string",
		},
		Model: openai.BetaResponseNewParamsModelGPT5_1,
		Moderation: openai.BetaResponseNewParamsModeration{
			Model: "model",
			Policy: openai.BetaResponseNewParamsModerationPolicy{
				Input: openai.BetaResponseNewParamsModerationPolicyInput{
					Mode: "score",
				},
				Output: openai.BetaResponseNewParamsModerationPolicyOutput{
					Mode: "score",
				},
			},
		},
		MultiAgent: openai.BetaResponseNewParamsMultiAgent{
			Enabled:                true,
			MaxConcurrentSubagents: openai.Int(1),
		},
		ParallelToolCalls:  openai.Bool(true),
		PreviousResponseID: openai.String("previous_response_id"),
		Prompt: openai.BetaResponsePromptParam{
			ID: "id",
			Variables: map[string]openai.BetaResponsePromptVariableUnionParam{
				"foo": {
					OfString: openai.String("string"),
				},
			},
			Version: openai.String("version"),
		},
		PromptCacheKey: openai.String("prompt-cache-key-1234"),
		PromptCacheOptions: openai.BetaResponseNewParamsPromptCacheOptions{
			Mode: "implicit",
			Ttl:  "30m",
		},
		PromptCacheRetention: openai.BetaResponseNewParamsPromptCacheRetentionInMemory,
		Reasoning: openai.BetaResponseNewParamsReasoning{
			Context:         "auto",
			Effort:          "none",
			GenerateSummary: "auto",
			Mode:            "standard",
			Summary:         "auto",
		},
		SafetyIdentifier: openai.String("safety-identifier-1234"),
		ServiceTier:      openai.BetaResponseNewParamsServiceTierAuto,
		Store:            openai.Bool(true),
		StreamOptions: openai.BetaResponseNewParamsStreamOptions{
			IncludeObfuscation: openai.Bool(true),
		},
		Temperature: openai.Float(1),
		Text: openai.BetaResponseTextConfigParam{
			Format: openai.BetaResponseFormatTextConfigUnionParam{
				OfText: &openai.BetaResponseFormatTextConfigTextParam{},
			},
			Verbosity: openai.BetaResponseTextConfigVerbosityLow,
		},
		ToolChoice: openai.BetaResponseNewParamsToolChoiceUnion{
			OfToolChoiceMode: openai.Opt(openai.BetaToolChoiceOptionsNone),
		},
		Tools: []openai.BetaToolUnionParam{{
			OfFunction: &openai.BetaFunctionToolParam{
				Name: "name",
				Parameters: map[string]any{
					"foo": "bar",
				},
				Strict:         openai.Bool(true),
				AllowedCallers: []string{"direct"},
				DeferLoading:   openai.Bool(true),
				Description:    openai.String("description"),
				OutputSchema: map[string]any{
					"foo": "bar",
				},
			},
		}},
		TopLogprobs: openai.Int(0),
		TopP:        openai.Float(1),
		Truncation:  openai.BetaResponseNewParamsTruncationAuto,
		User:        openai.String("user-1234"),
		Betas:       []string{"responses_multi_agent=v1"},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaResponseGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Responses.Get(
		context.TODO(),
		"resp_677efb5139a88190b512bc3fef8e535d",
		openai.BetaResponseGetParams{
			Include:            []openai.BetaResponseIncludable{openai.BetaResponseIncludableFileSearchCallResults},
			IncludeObfuscation: openai.Bool(true),
			StartingAfter:      openai.Int(0),
			Betas:              []string{"responses_multi_agent=v1"},
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

func TestBetaResponseDeleteWithOptionalParams(t *testing.T) {
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
	err := client.Beta.Responses.Delete(
		context.TODO(),
		"resp_677efb5139a88190b512bc3fef8e535d",
		openai.BetaResponseDeleteParams{
			Betas: []string{"responses_multi_agent=v1"},
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

func TestBetaResponseCancelWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Responses.Cancel(
		context.TODO(),
		"resp_677efb5139a88190b512bc3fef8e535d",
		openai.BetaResponseCancelParams{
			Betas: []string{"responses_multi_agent=v1"},
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

func TestBetaResponseCompactWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Responses.Compact(context.TODO(), openai.BetaResponseCompactParams{
		Model: openai.BetaResponseCompactParamsModelGPT5_6Sol,
		Input: openai.BetaResponseCompactParamsInputUnion{
			OfString: openai.String("string"),
		},
		Instructions:       openai.String("instructions"),
		PreviousResponseID: openai.String("resp_123"),
		PromptCacheKey:     openai.String("prompt_cache_key"),
		PromptCacheOptions: openai.BetaResponseCompactParamsPromptCacheOptions{
			Mode: "implicit",
			Ttl:  "30m",
		},
		PromptCacheRetention: openai.BetaResponseCompactParamsPromptCacheRetentionInMemory,
		ServiceTier:          openai.BetaResponseCompactParamsServiceTierAuto,
		Betas:                []string{"responses_multi_agent=v1"},
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
