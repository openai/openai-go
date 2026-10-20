package bedrock

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

// TestRuntimeLive makes AWS requests only when BEDROCK_LIVE_TEST=1 is set.
// Select credentials with BEDROCK_LIVE_AUTH and optionally enable Responses
// and streaming with BEDROCK_LIVE_RESPONSES=1 and BEDROCK_LIVE_STREAM=1.
func TestRuntimeLive(t *testing.T) {
	if os.Getenv("BEDROCK_LIVE_TEST") != "1" {
		t.Skip("set BEDROCK_LIVE_TEST=1 to make live Amazon Bedrock requests")
	}
	mode := strings.TrimSpace(os.Getenv("BEDROCK_LIVE_AUTH"))
	if mode == "" {
		t.Fatal("set BEDROCK_LIVE_AUTH before making live Amazon Bedrock requests")
	}
	config := liveRuntimeConfig(t, mode)
	client, err := NewClient(context.Background(), config, option.WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}

	models := []string{"us.openai.gpt-5.6-sol", "us.openai.gpt-5.6-terra", "us.openai.gpt-5.6-luna"}
	if configured := os.Getenv("BEDROCK_LIVE_MODELS"); configured != "" {
		models = strings.Split(configured, ",")
	} else if configured := os.Getenv("BEDROCK_MODEL"); configured != "" {
		models = []string{configured}
	}
	for _, configured := range models {
		model := strings.TrimSpace(configured)
		if model == "" {
			t.Fatal("BEDROCK_LIVE_MODELS contains an empty inference profile")
		}
		t.Run(fmt.Sprintf("%s/%s", mode, model), func(t *testing.T) {
			chatParams := openai.ChatCompletionNewParams{
				Model:    openai.ChatModel(model),
				Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage("Reply with exactly: bedrock live test passed")},
			}
			completion, err := client.Chat.Completions.New(context.Background(), chatParams)
			if err != nil || len(completion.Choices) == 0 || completion.Choices[0].Message.Content == "" {
				t.Fatalf("chat completion = %#v, error = %v", completion, err)
			}

			if os.Getenv("BEDROCK_LIVE_STREAM") == "1" {
				stream := client.Chat.Completions.NewStreaming(context.Background(), chatParams)
				var content string
				for stream.Next() {
					if choices := stream.Current().Choices; len(choices) != 0 {
						content += choices[0].Delta.Content
					}
				}
				if streamErr := stream.Err(); streamErr != nil || content == "" {
					t.Fatalf("chat stream content = %q, error = %v", content, streamErr)
				}
			}

			if os.Getenv("BEDROCK_LIVE_RESPONSES") != "1" {
				return
			}
			responseParams := responses.ResponseNewParams{
				Model: openai.ChatModel(model),
				Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("Reply with exactly: bedrock live test passed")},
			}
			response, err := client.Responses.New(context.Background(), responseParams)
			if err != nil || response.OutputText() == "" {
				t.Fatalf("response = %#v, error = %v", response, err)
			}
			if os.Getenv("BEDROCK_LIVE_STREAM") == "1" {
				stream := client.Responses.NewStreaming(context.Background(), responseParams)
				var content string
				for stream.Next() {
					content += stream.Current().Delta
				}
				if streamErr := stream.Err(); streamErr != nil || content == "" {
					t.Fatalf("response stream content = %q, error = %v", content, streamErr)
				}
			}
		})
	}
}

func liveRuntimeConfig(t *testing.T, mode string) Config {
	t.Helper()
	config := Config{Endpoint: EndpointRuntime}
	switch mode {
	case "bearer":
		config.APIKey = requiredLiveRuntimeEnv(t, "AWS_BEARER_TOKEN_BEDROCK")
	case "environment-bearer":
		requiredLiveRuntimeEnv(t, "AWS_BEARER_TOKEN_BEDROCK")
	case "token-provider":
		requiredLiveRuntimeEnv(t, "AWS_BEARER_TOKEN_BEDROCK")
		config.BedrockTokenProvider = func(context.Context) (string, error) {
			token := strings.TrimSpace(os.Getenv("AWS_BEARER_TOKEN_BEDROCK"))
			if token == "" {
				return "", fmt.Errorf("AWS_BEARER_TOKEN_BEDROCK is no longer configured")
			}
			return token, nil
		}
	case "default-chain":
		if os.Getenv("AWS_BEARER_TOKEN_BEDROCK") != "" {
			t.Fatal("unset AWS_BEARER_TOKEN_BEDROCK to select default-chain SigV4 authentication")
		}
	case "profile":
		config.AWSProfile = requiredLiveRuntimeEnv(t, "AWS_PROFILE")
	case "static":
		config.AWSAccessKeyID = requiredLiveRuntimeEnv(t, "AWS_ACCESS_KEY_ID")
		config.AWSSecretAccessKey = requiredLiveRuntimeEnv(t, "AWS_SECRET_ACCESS_KEY")
		config.AWSSessionToken = strings.TrimSpace(os.Getenv("AWS_SESSION_TOKEN"))
	case "custom-provider":
		requiredLiveRuntimeEnv(t, "AWS_ACCESS_KEY_ID")
		requiredLiveRuntimeEnv(t, "AWS_SECRET_ACCESS_KEY")
		config.AWSCredentialsProvider = aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
			accessKeyID := strings.TrimSpace(os.Getenv("AWS_ACCESS_KEY_ID"))
			secretAccessKey := strings.TrimSpace(os.Getenv("AWS_SECRET_ACCESS_KEY"))
			if accessKeyID == "" || secretAccessKey == "" {
				return aws.Credentials{}, fmt.Errorf("AWS static credentials are no longer configured")
			}
			return aws.Credentials{
				AccessKeyID:     accessKeyID,
				SecretAccessKey: secretAccessKey,
				SessionToken:    strings.TrimSpace(os.Getenv("AWS_SESSION_TOKEN")),
			}, nil
		})
	default:
		t.Fatalf("unsupported BEDROCK_LIVE_AUTH %q; use bearer, environment-bearer, token-provider, default-chain, profile, static, or custom-provider", mode)
	}
	return config
}

func requiredLiveRuntimeEnv(t *testing.T, name string) string {
	t.Helper()
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		t.Fatalf("set %s before making live Amazon Bedrock requests", name)
	}
	return value
}
