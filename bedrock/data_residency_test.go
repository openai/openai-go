package bedrock

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

func TestDataResidencyRejectedBeforeCredentialDiscovery(t *testing.T) {
	for _, test := range []struct {
		name    string
		opts    []option.RequestOption
		message string
	}{
		{"valid region", []option.RequestOption{option.WithDataResidency(option.DataResidencyEU)}, "Bedrock provider"},
		{"empty region", []option.RequestOption{option.WithDataResidency("")}, "invalid data residency"},
		{"unknown region", []option.RequestOption{option.WithDataResidency("unknown")}, "invalid data residency"},
		{"base URL first", []option.RequestOption{option.WithBaseURL("https://custom.example"), option.WithDataResidency(option.DataResidencyEU)}, "Bedrock provider"},
		{"residency first", []option.RequestOption{option.WithDataResidency(option.DataResidencyEU), option.WithBaseURL("https://custom.example")}, "Bedrock provider"},
		{"inherited region", requestconfig.InheritedOptions(option.WithDataResidency(option.DataResidencyEU)), "Bedrock provider"},
	} {
		t.Run(test.name, func(t *testing.T) {
			clearAWSEnvironment(t)
			credentialCalls, callbackCalls := 0, 0
			config := Config{
				AWSRegion: "us-east-1",
				AWSCredentialsProvider: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					credentialCalls++
					return aws.Credentials{AccessKeyID: "dummy-access", SecretAccessKey: "dummy-secret"}, nil
				}),
			}
			callback := requestconfig.RequestOptionFunc(func(*requestconfig.RequestConfig) error { callbackCalls++; return nil })
			opts := append([]option.RequestOption{callback}, test.opts...)
			_, err := NewClient(context.Background(), config, opts...)
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("error = %v, want %q", err, test.message)
			}
			if credentialCalls != 0 || callbackCalls != 0 {
				t.Fatalf("credential calls = %d, arbitrary callbacks = %d", credentialCalls, callbackCalls)
			}
		})
	}
}

func TestEndpointPreflightDoesNotEvaluateUserCallbacks(t *testing.T) {
	clearAWSEnvironment(t)
	callbackCalls := 0
	callback := requestconfig.RequestOptionFunc(func(*requestconfig.RequestConfig) error { callbackCalls++; return nil })
	client, err := NewClient(context.Background(), Config{SkipAuth: true, BaseURL: "https://bedrock.example/openai/v1"}, callback)
	if err != nil {
		t.Fatal(err)
	}
	if callbackCalls != 0 {
		t.Fatalf("constructor evaluated user callback %d times", callbackCalls)
	}
	// Applying the stored configuration still evaluates an ordinary callback once.
	cfg := requestconfig.RequestConfig{}
	if err := cfg.Apply(client.Options...); err != nil {
		t.Fatal(err)
	}
	if callbackCalls != 1 {
		t.Fatalf("request setup evaluated callback %d times", callbackCalls)
	}
}

func TestDataResidencyRequestOverrideRejected(t *testing.T) {
	clearAWSEnvironment(t)
	transportCalls := 0
	client, err := NewClient(context.Background(), Config{SkipAuth: true, BaseURL: "https://bedrock.example/openai/v1"},
		option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			transportCalls++
			return successfulResponse(req), nil
		})}),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = client.Get(context.Background(), "models", nil, nil, option.WithDataResidency(option.DataResidencyEU))
	if err == nil || !strings.Contains(err.Error(), "Bedrock provider") {
		t.Fatalf("expected provider error, got %v", err)
	}
	if transportCalls != 0 {
		t.Fatalf("provider conflict made %d HTTP requests", transportCalls)
	}
}
