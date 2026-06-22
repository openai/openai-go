// Package bedrock configures the standard OpenAI client for Amazon Bedrock's
// OpenAI-compatible API.
//
// New clients use the AWS SDK for Go v2 credential chain by default, including
// environment credentials, shared credentials and config files, named profiles,
// workload identity, and instance or container roles. Requests are signed with
// AWS Signature Version 4 using the bedrock-mantle service name.
package bedrock

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// TokenProvider resolves a Bedrock bearer credential. It is called before each
// request attempt so expiring credentials can be refreshed for retries.
type TokenProvider func(context.Context) (string, error)

// Config configures an Amazon Bedrock client.
//
// Authentication is selected in this order:
//
//  1. An explicit APIKey or BedrockTokenProvider.
//  2. Explicit static AWS credentials, AWSProfile, or AWSCredentialsProvider.
//  3. AWS_BEARER_TOKEN_BEDROCK.
//  4. The standard AWS SDK credential chain.
//
// Configure at most one explicit authentication mode.
type Config struct {
	// APIKey is an explicit Amazon Bedrock bearer credential.
	APIKey string

	// BedrockTokenProvider resolves a bearer credential before each request
	// attempt. It is mutually exclusive with APIKey and AWS credential options.
	BedrockTokenProvider TokenProvider

	// AWSAccessKeyID and AWSSecretAccessKey configure static AWS credentials.
	// They must be provided together.
	AWSAccessKeyID     string
	AWSSecretAccessKey string

	// AWSSessionToken is the optional session token for temporary static AWS
	// credentials. It may only be used with AWSAccessKeyID and
	// AWSSecretAccessKey.
	AWSSessionToken string

	// AWSProfile selects a named profile from the standard AWS shared config and
	// credentials files.
	AWSProfile string

	// AWSRegion is used for endpoint resolution and SigV4 signing. It takes
	// precedence over AWS_REGION, AWS_DEFAULT_REGION, and the AWS SDK region
	// provider chain.
	AWSRegion string

	// AWSCredentialsProvider supplies refreshable AWS credentials. Providers are
	// wrapped in the AWS SDK's credentials cache so expiring credentials refresh
	// according to normal AWS semantics.
	AWSCredentialsProvider aws.CredentialsProvider

	// BaseURL overrides AWS_BEDROCK_BASE_URL and the regional Mantle endpoint.
	// The default is https://bedrock-mantle.{region}.api.aws/openai/v1.
	BaseURL string

	// SkipAuth disables Bedrock bearer and SigV4 authentication. Use this only
	// when a trusted gateway or custom transport authenticates requests.
	SkipAuth bool
}

// NewClient creates a standard OpenAI client configured for Amazon Bedrock.
//
// Generic transport options, such as option.WithHTTPClient,
// option.WithMaxRetries, and option.WithMiddleware, may be supplied in opts.
// Provider authentication and BaseURL must be configured through Config. The
// Bedrock authentication middleware always runs after client- and method-level
// middleware so it signs the final request on every attempt.
func NewClient(ctx context.Context, cfg Config, opts ...option.RequestOption) (openai.Client, error) {
	resolvedOpts, err := newClientOptions(ctx, cfg, time.Now, opts...)
	if err != nil {
		return openai.Client{}, err
	}
	return openai.NewClient(resolvedOpts...), nil
}
