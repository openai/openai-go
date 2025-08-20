// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package responses_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/internal/testutil"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/packages/param"
	"github.com/openai/openai-go/v2/responses"
	"github.com/openai/openai-go/v2/shared"
)

func TestResponseNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Responses.New(context.TODO(), responses.ResponseNewParams{
		Background: openai.Bool(true),
		Include:    []responses.ResponseIncludable{responses.ResponseIncludableCodeInterpreterCallOutputs},
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("string"),
		},
		Instructions:    openai.String("instructions"),
		MaxOutputTokens: openai.Int(0),
		MaxToolCalls:    openai.Int(0),
		Metadata: shared.Metadata{
			"foo": "string",
		},
		Model:              shared.ResponsesModel("gpt-4o"),
		ParallelToolCalls:  openai.Bool(true),
		PreviousResponseID: openai.String("previous_response_id"),
		Prompt: responses.ResponsePromptParam{
			ID: "id",
			Variables: map[string]responses.ResponsePromptVariableUnionParam{
				"foo": {
					OfString: openai.String("string"),
				},
			},
			Version: openai.String("version"),
		},
		PromptCacheKey: openai.String("prompt-cache-key-1234"),
		Reasoning: shared.ReasoningParam{
			Effort:          shared.ReasoningEffortMinimal,
			GenerateSummary: shared.ReasoningGenerateSummaryAuto,
			Summary:         shared.ReasoningSummaryAuto,
		},
		SafetyIdentifier: openai.String("safety-identifier-1234"),
		ServiceTier:      responses.ResponseNewParamsServiceTierAuto,
		Store:            openai.Bool(true),
		StreamOptions: responses.ResponseNewParamsStreamOptions{
			IncludeObfuscation: openai.Bool(true),
		},
		Temperature: openai.Float(1),
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfText: &shared.ResponseFormatTextParam{},
			},
			Verbosity: responses.ResponseTextConfigVerbosityLow,
		},
		ToolChoice: responses.ResponseNewParamsToolChoiceUnion{
			OfToolChoiceMode: openai.Opt(responses.ToolChoiceOptionsNone),
		},
		Tools: []responses.ToolUnionParam{{
			OfFunction: &responses.FunctionToolParam{
				Name: "name",
				Parameters: map[string]any{
					"foo": "bar",
				},
				Strict:      openai.Bool(true),
				Description: openai.String("description"),
			},
		}},
		TopLogprobs: openai.Int(0),
		TopP:        openai.Float(1),
		Truncation:  responses.ResponseNewParamsTruncationAuto,
		User:        openai.String("user-1234"),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestResponseGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Responses.Get(
		context.TODO(),
		"resp_677efb5139a88190b512bc3fef8e535d",
		responses.ResponseGetParams{
			Include:            []responses.ResponseIncludable{responses.ResponseIncludableCodeInterpreterCallOutputs},
			IncludeObfuscation: openai.Bool(true),
			StartingAfter:      openai.Int(0),
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

func TestResponseDelete(t *testing.T) {
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
	err := client.Responses.Delete(context.TODO(), "resp_677efb5139a88190b512bc3fef8e535d")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestResponseCancel(t *testing.T) {
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
	_, err := client.Responses.Cancel(context.TODO(), "resp_677efb5139a88190b512bc3fef8e535d")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestToolParamOfMcp(t *testing.T) {
	// Test with all parameters
	t.Run("with all parameters", func(t *testing.T) {
		serverLabel := "test-server"
		serverURL := "https://test-mcp-server.com"
		allowedTools := []string{"search_products", "send_email"}
		requireApproval := "never"
		
		tool := responses.ToolParamOfMcp(serverLabel, serverURL, allowedTools, requireApproval)
		
		if tool.OfMcp == nil {
			t.Fatal("Expected OfMcp to be non-nil")
		}
		
		mcp := tool.OfMcp
		if mcp.ServerLabel != serverLabel {
			t.Errorf("Expected ServerLabel to be %s, got %s", serverLabel, mcp.ServerLabel)
		}
		
		if mcp.ServerURL != serverURL {
			t.Errorf("Expected ServerURL to be %s, got %s", serverURL, mcp.ServerURL)
		}
		
		if param.IsOmitted(mcp.AllowedTools) {
			t.Error("Expected AllowedTools to be set")
		} else {
			if mcp.AllowedTools.OfMcpAllowedTools == nil {
				t.Error("Expected OfMcpAllowedTools to be set")
			} else {
				tools := mcp.AllowedTools.OfMcpAllowedTools
				if len(tools) != 2 || tools[0] != "search_products" || tools[1] != "send_email" {
					t.Errorf("Expected allowedTools to be %v, got %v", allowedTools, tools)
				}
			}
		}
		
		if param.IsOmitted(mcp.RequireApproval) {
			t.Error("Expected RequireApproval to be set")
		} else {
			if param.IsOmitted(mcp.RequireApproval.OfMcpToolApprovalSetting) {
				t.Error("Expected OfMcpToolApprovalSetting to be set")
			} else {
				if mcp.RequireApproval.OfMcpToolApprovalSetting.Value != requireApproval {
					t.Errorf("Expected requireApproval to be %s, got %s", requireApproval, mcp.RequireApproval.OfMcpToolApprovalSetting.Value)
				}
			}
		}
	})
	
	// Test with empty allowedTools
	t.Run("with empty allowed tools", func(t *testing.T) {
		serverLabel := "test-server"
		serverURL := "https://test-mcp-server.com"
		var allowedTools []string // empty slice
		requireApproval := "auto"
		
		tool := responses.ToolParamOfMcp(serverLabel, serverURL, allowedTools, requireApproval)
		
		if tool.OfMcp == nil {
			t.Fatal("Expected OfMcp to be non-nil")
		}
		
		mcp := tool.OfMcp
		if !param.IsOmitted(mcp.AllowedTools) {
			t.Error("Expected AllowedTools to be omitted when empty slice provided")
		}
		
		if param.IsOmitted(mcp.RequireApproval) {
			t.Error("Expected RequireApproval to be set")
		}
	})
	
	// Test with nil allowedTools
	t.Run("with nil allowed tools", func(t *testing.T) {
		serverLabel := "test-server"
		serverURL := "https://test-mcp-server.com"
		var allowedTools []string = nil // nil slice
		requireApproval := "always"
		
		tool := responses.ToolParamOfMcp(serverLabel, serverURL, allowedTools, requireApproval)
		
		if tool.OfMcp == nil {
			t.Fatal("Expected OfMcp to be non-nil")
		}
		
		mcp := tool.OfMcp
		if !param.IsOmitted(mcp.AllowedTools) {
			t.Error("Expected AllowedTools to be omitted when nil slice provided")
		}
	})
	
	// Test with empty requireApproval
	t.Run("with empty require approval", func(t *testing.T) {
		serverLabel := "test-server"
		serverURL := "https://test-mcp-server.com"
		allowedTools := []string{"search_products"}
		requireApproval := "" // empty string
		
		tool := responses.ToolParamOfMcp(serverLabel, serverURL, allowedTools, requireApproval)
		
		if tool.OfMcp == nil {
			t.Fatal("Expected OfMcp to be non-nil")
		}
		
		mcp := tool.OfMcp
		if !param.IsOmitted(mcp.RequireApproval) {
			t.Error("Expected RequireApproval to be omitted when empty string provided")
		}
		
		if param.IsOmitted(mcp.AllowedTools) {
			t.Error("Expected AllowedTools to be set")
		}
	})
	
	// Test minimal case (both optional parameters empty)
	t.Run("minimal parameters", func(t *testing.T) {
		serverLabel := "minimal-server"
		serverURL := "https://minimal-server.com"
		var allowedTools []string // empty
		requireApproval := "" // empty
		
		tool := responses.ToolParamOfMcp(serverLabel, serverURL, allowedTools, requireApproval)
		
		if tool.OfMcp == nil {
			t.Fatal("Expected OfMcp to be non-nil")
		}
		
		mcp := tool.OfMcp
		if mcp.ServerLabel != serverLabel {
			t.Errorf("Expected ServerLabel to be %s, got %s", serverLabel, mcp.ServerLabel)
		}
		
		if mcp.ServerURL != serverURL {
			t.Errorf("Expected ServerURL to be %s, got %s", serverURL, mcp.ServerURL)
		}
		
		if !param.IsOmitted(mcp.AllowedTools) {
			t.Error("Expected AllowedTools to be omitted")
		}
		
		if !param.IsOmitted(mcp.RequireApproval) {
			t.Error("Expected RequireApproval to be omitted")
		}
	})
}
