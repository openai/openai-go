package bedrock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func successfulResponse(req *http.Request) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{}`)),
		Request:    req,
	}
}

func clearAWSEnvironment(t *testing.T) {
	t.Helper()
	for _, name := range []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_SESSION_TOKEN",
		"AWS_PROFILE",
		"AWS_DEFAULT_PROFILE",
		"AWS_REGION",
		"AWS_DEFAULT_REGION",
		"AWS_SHARED_CREDENTIALS_FILE",
		"AWS_CONFIG_FILE",
		"AWS_WEB_IDENTITY_TOKEN_FILE",
		"AWS_ROLE_ARN",
		"AWS_ROLE_SESSION_NAME",
		"AWS_CONTAINER_CREDENTIALS_FULL_URI",
		"AWS_CONTAINER_CREDENTIALS_RELATIVE_URI",
		"AWS_BEARER_TOKEN_BEDROCK",
		"AWS_BEDROCK_BASE_URL",
	} {
		value, present := os.LookupEnv(name)
		if err := os.Unsetenv(name); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if present {
				_ = os.Setenv(name, value)
			} else {
				_ = os.Unsetenv(name)
			}
		})
	}
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
}

func writeAWSFiles(t *testing.T, credentials, config string) (string, string) {
	t.Helper()
	dir := t.TempDir()
	credentialsPath := filepath.Join(dir, "credentials")
	configPath := filepath.Join(dir, "config")
	if err := os.WriteFile(credentialsPath, []byte(credentials), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	return credentialsPath, configPath
}

func TestSharedCredentialsFile(t *testing.T) {
	clearAWSEnvironment(t)
	credentialsPath, configPath := writeAWSFiles(t, `
[engineering]
aws_access_key_id = profile-access-key
aws_secret_access_key = profile-secret-key
aws_session_token = profile-session-token
`, `
[profile engineering]
region = eu-central-1
`)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", credentialsPath)
	t.Setenv("AWS_CONFIG_FILE", configPath)
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "ambient-bearer-must-not-win")

	var request *http.Request
	client, err := NewClient(context.Background(), Config{AWSProfile: "engineering"}, option.WithHTTPClient(&http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			request = req
			return successfulResponse(req), nil
		}),
	}))
	if err != nil {
		t.Fatal(err)
	}

	var response *http.Response
	if err := client.Get(context.Background(), "/models", nil, &response); err != nil {
		t.Fatal(err)
	}
	if got, want := request.URL.String(), "https://bedrock-mantle.eu-central-1.api.aws/openai/v1/models"; got != want {
		t.Fatalf("request URL = %q, want %q", got, want)
	}
	authorization := request.Header.Get("Authorization")
	if !strings.Contains(authorization, "Credential=profile-access-key/") || !strings.Contains(authorization, "/eu-central-1/bedrock-mantle/aws4_request") {
		t.Fatalf("unexpected authorization header: %q", authorization)
	}
	if got := request.Header.Get("X-Amz-Security-Token"); got != "profile-session-token" {
		t.Fatalf("X-Amz-Security-Token = %q", got)
	}
}

func TestDefaultCredentialChainUsesDefaultSharedProfile(t *testing.T) {
	clearAWSEnvironment(t)
	credentialsPath, configPath := writeAWSFiles(t, `
[default]
aws_access_key_id = default-access-key
aws_secret_access_key = default-secret-key
`, `
[default]
region = us-west-2
`)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", credentialsPath)
	t.Setenv("AWS_CONFIG_FILE", configPath)

	var request *http.Request
	client, err := NewClient(context.Background(), Config{}, option.WithHTTPClient(&http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			request = req
			return successfulResponse(req), nil
		}),
	}))
	if err != nil {
		t.Fatal(err)
	}
	var response *http.Response
	if err := client.Get(context.Background(), "/models", nil, &response); err != nil {
		t.Fatal(err)
	}
	if got := request.Header.Get("Authorization"); !strings.Contains(got, "Credential=default-access-key/") {
		t.Fatalf("authorization did not use the default shared profile: %q", got)
	}
	if got := request.URL.Host; got != "bedrock-mantle.us-west-2.api.aws" {
		t.Fatalf("request host = %q", got)
	}
}

func TestSigV4Fixture(t *testing.T) {
	var fixture struct {
		SigningDate string `json:"signingDate"`
		Region      string `json:"region"`
		Service     string `json:"service"`
		Credentials struct {
			AccessKeyID     string `json:"accessKeyId"`
			SecretAccessKey string `json:"secretAccessKey"`
			SessionToken    string `json:"sessionToken"`
		} `json:"credentials"`
		Request struct {
			Method      string `json:"method"`
			URL         string `json:"url"`
			Body        string `json:"body"`
			ContentType string `json:"contentType"`
		} `json:"request"`
		Expected struct {
			Date          string `json:"date"`
			PayloadHash   string `json:"payloadHash"`
			Authorization string `json:"authorization"`
		} `json:"expected"`
	}
	contents, err := os.ReadFile("testdata/sigv4.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(contents, &fixture); err != nil {
		t.Fatal(err)
	}
	if fixture.Service != bedrockService {
		t.Fatalf("fixture service = %q", fixture.Service)
	}
	signingDate, err := time.Parse(time.RFC3339Nano, fixture.SigningDate)
	if err != nil {
		t.Fatal(err)
	}
	baseURL, _ := url.Parse("https://bedrock-mantle.us-east-1.api.aws/openai/v1/")
	request, err := http.NewRequest(fixture.Request.Method, fixture.Request.URL, strings.NewReader(fixture.Request.Body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", fixture.Request.ContentType)
	awsConfig := aws.Config{
		Region: fixture.Region,
		Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID: fixture.Credentials.AccessKeyID, SecretAccessKey: fixture.Credentials.SecretAccessKey,
				SessionToken: fixture.Credentials.SessionToken,
			}, nil
		}),
	}

	called := false
	_, err = sigV4Middleware(baseURL, awsConfig, v4.NewSigner(), func() time.Time { return signingDate })(request, func(req *http.Request) (*http.Response, error) {
		called = true
		return successfulResponse(req), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("transport was not called")
	}
	if got := request.Header.Get("X-Amz-Date"); got != fixture.Expected.Date {
		t.Fatalf("X-Amz-Date = %q, want %q", got, fixture.Expected.Date)
	}
	if got := request.Header.Get("X-Amz-Content-Sha256"); got != fixture.Expected.PayloadHash {
		t.Fatalf("payload hash = %q, want %q", got, fixture.Expected.PayloadHash)
	}
	if got := request.Header.Get("Authorization"); got != fixture.Expected.Authorization {
		t.Fatalf("authorization =\n%s\nwant\n%s", got, fixture.Expected.Authorization)
	}
}

func TestRetriesAreResigned(t *testing.T) {
	clearAWSEnvironment(t)
	times := []time.Time{
		time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC),
		time.Date(2025, 1, 2, 3, 5, 5, 0, time.UTC),
	}
	var clockCalls atomic.Int32
	clock := func() time.Time {
		index := int(clockCalls.Add(1)) - 1
		if index >= len(times) {
			return times[len(times)-1]
		}
		return times[index]
	}

	var dates, authorizations []string
	var attempts int
	opts, err := newClientOptions(context.Background(), Config{
		AWSRegion:          "us-east-1",
		AWSAccessKeyID:     "access-key",
		AWSSecretAccessKey: "secret-key",
	}, clock,
		option.WithMaxRetries(1),
		option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			attempts++
			dates = append(dates, req.Header.Get("X-Amz-Date"))
			authorizations = append(authorizations, req.Header.Get("Authorization"))
			if attempts == 1 {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Header:     http.Header{"Retry-After": {"0"}},
					Body:       io.NopCloser(strings.NewReader(`{}`)),
					Request:    req,
				}, nil
			}
			return successfulResponse(req), nil
		})}),
	)
	if err != nil {
		t.Fatal(err)
	}
	client := openai.NewClient(opts...)
	var response *http.Response
	if err := client.Get(context.Background(), "/models", nil, &response); err != nil {
		t.Fatal(err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
	if got, want := strings.Join(dates, ","), "20250102T030405Z,20250102T030505Z"; got != want {
		t.Fatalf("signing dates = %q, want %q", got, want)
	}
	if authorizations[0] == authorizations[1] {
		t.Fatal("retry reused the first request signature")
	}
}

func TestExpiringCredentialsRefreshAcrossRetries(t *testing.T) {
	clearAWSEnvironment(t)
	var providerCalls atomic.Int32
	provider := aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
		call := providerCalls.Add(1)
		return aws.Credentials{
			AccessKeyID:     fmt.Sprintf("refreshing-key-%d", call),
			SecretAccessKey: "secret-key",
			CanExpire:       true,
			Expires:         time.Now().Add(-time.Minute),
		}, nil
	})

	var credentials []string
	var attempts int
	client, err := NewClient(context.Background(), Config{
		AWSRegion:              "us-east-1",
		AWSCredentialsProvider: provider,
	},
		option.WithMaxRetries(1),
		option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			attempts++
			authorization := req.Header.Get("Authorization")
			start := strings.Index(authorization, "Credential=")
			end := strings.Index(authorization, "/")
			if start < 0 || end < 0 {
				t.Fatalf("unexpected authorization: %q", authorization)
			}
			credentials = append(credentials, authorization[start+len("Credential="):end])
			if attempts == 1 {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Header:     http.Header{"Retry-After": {"0"}},
					Body:       io.NopCloser(strings.NewReader(`{}`)),
					Request:    req,
				}, nil
			}
			return successfulResponse(req), nil
		})}),
	)
	if err != nil {
		t.Fatal(err)
	}
	var response *http.Response
	if err := client.Get(context.Background(), "/models", nil, &response); err != nil {
		t.Fatal(err)
	}
	if attempts != 2 || len(credentials) != 2 {
		t.Fatalf("attempts = %d, credentials = %v", attempts, credentials)
	}
	if credentials[0] == credentials[1] {
		t.Fatalf("retry reused expired credentials: %v", credentials)
	}
	if calls := providerCalls.Load(); calls < 3 {
		t.Fatalf("credential provider calls = %d, want at least 3", calls)
	}
}

func TestSigningRunsAfterClientAndMethodMiddleware(t *testing.T) {
	clearAWSEnvironment(t)
	var request *http.Request
	client, err := NewClient(context.Background(), Config{
		AWSRegion:          "us-east-1",
		AWSAccessKeyID:     "access-key",
		AWSSecretAccessKey: "secret-key",
	},
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			req.Header.Set("X-Client-Middleware", "present")
			return next(req)
		}),
		option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			request = req
			return successfulResponse(req), nil
		})}),
	)
	if err != nil {
		t.Fatal(err)
	}
	var response *http.Response
	if err := client.Get(context.Background(), "/models", nil, &response,
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			req.Header.Set("X-Method-Middleware", "present")
			return next(req)
		}),
	); err != nil {
		t.Fatal(err)
	}
	authorization := request.Header.Get("Authorization")
	if !strings.Contains(authorization, "x-client-middleware") || !strings.Contains(authorization, "x-method-middleware") {
		t.Fatalf("middleware headers were not signed: %q", authorization)
	}
}

func TestAmbientOpenAIConfigurationIsNotInherited(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "openai-key")
	t.Setenv("OPENAI_ADMIN_KEY", "openai-admin-key")
	t.Setenv("OPENAI_BASE_URL", "https://openai-environment.example/v1")
	t.Setenv("OPENAI_ORG_ID", "org-environment")
	t.Setenv("OPENAI_PROJECT_ID", "project-environment")
	t.Setenv("OPENAI_CUSTOM_HEADERS", "X-Ambient: must-not-leak")

	var request *http.Request
	client, err := NewClient(context.Background(), Config{
		APIKey:  "bedrock-key",
		BaseURL: "https://bedrock.example/openai/v1",
	}, option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		request = req
		return successfulResponse(req), nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	var response *http.Response
	if err := client.Get(context.Background(), "/models", nil, &response); err != nil {
		t.Fatal(err)
	}
	if got := request.URL.Host; got != "bedrock.example" {
		t.Fatalf("request host = %q", got)
	}
	if got := request.Header.Get("Authorization"); got != "Bearer bedrock-key" {
		t.Fatalf("Authorization = %q", got)
	}
	for _, header := range []string{"OpenAI-Organization", "OpenAI-Project", "X-Ambient"} {
		if value := request.Header.Get(header); value != "" {
			t.Fatalf("ambient header %s leaked with value %q", header, value)
		}
	}
}

func TestRejectsCrossOriginBeforeBearerResolution(t *testing.T) {
	var attempts, providerCalls, transportCalls int
	client, err := NewClient(context.Background(), Config{
		BaseURL: "https://bedrock.example/openai/v1",
		BedrockTokenProvider: func(context.Context) (string, error) {
			providerCalls++
			return "token", nil
		},
	},
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			attempts++
			return next(req)
		}),
		option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			transportCalls++
			return successfulResponse(req), nil
		})}),
	)
	if err != nil {
		t.Fatal(err)
	}
	var response *http.Response
	err = client.Get(context.Background(), "https://attacker.example/steal", nil, &response)
	if err == nil || !strings.Contains(err.Error(), "origin other than") {
		t.Fatalf("error = %v", err)
	}
	if attempts != 1 || providerCalls != 0 || transportCalls != 0 {
		t.Fatalf("attempts = %d, provider calls = %d, transport calls = %d", attempts, providerCalls, transportCalls)
	}
}

func TestRejectsCustomAuthorizationBeforeBearerResolution(t *testing.T) {
	var providerCalls, transportCalls int
	client, err := NewClient(context.Background(), Config{
		BaseURL: "https://bedrock.example/openai/v1",
		BedrockTokenProvider: func(context.Context) (string, error) {
			providerCalls++
			return "token", nil
		},
	}, option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		transportCalls++
		return successfulResponse(req), nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	var response *http.Response
	err = client.Get(context.Background(), "/models", nil, &response, option.WithHeader("Authorization", "Bearer custom"))
	if err == nil || !strings.Contains(err.Error(), "custom `Authorization` header") {
		t.Fatalf("error = %v", err)
	}
	if providerCalls != 0 || transportCalls != 0 {
		t.Fatalf("provider calls = %d, transport calls = %d", providerCalls, transportCalls)
	}
}

func TestSigV4RejectsNonReplayableBody(t *testing.T) {
	clearAWSEnvironment(t)
	var transportCalls int
	client, err := NewClient(context.Background(), Config{
		AWSRegion:          "us-east-1",
		AWSAccessKeyID:     "access-key",
		AWSSecretAccessKey: "secret-key",
	}, option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		transportCalls++
		return successfulResponse(req), nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	var response *http.Response
	err = client.Post(context.Background(), "/responses", io.LimitReader(strings.NewReader("body"), 4), &response)
	if err == nil || !strings.Contains(err.Error(), "replayable request body") {
		t.Fatalf("error = %v", err)
	}
	if transportCalls != 0 {
		t.Fatalf("transport calls = %d", transportCalls)
	}
}

func TestSigV4DisablesRedirects(t *testing.T) {
	clearAWSEnvironment(t)
	var targetCalls atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		targetCalls.Add(1)
	}))
	defer target.Close()
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, target.URL, http.StatusFound)
	}))
	defer source.Close()

	client, err := NewClient(context.Background(), Config{
		AWSRegion:          "us-east-1",
		AWSAccessKeyID:     "access-key",
		AWSSecretAccessKey: "secret-key",
		BaseURL:            source.URL,
	})
	if err != nil {
		t.Fatal(err)
	}
	var response *http.Response
	if err := client.Get(context.Background(), "/models", nil, &response); err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusFound {
		t.Fatalf("status = %d", response.StatusCode)
	}
	if calls := targetCalls.Load(); calls != 0 {
		t.Fatalf("redirect target calls = %d", calls)
	}
}

func TestBearerDisablesOriginChangingRedirects(t *testing.T) {
	var sourceAuthorization string
	var targetCalls atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		targetCalls.Add(1)
	}))
	defer target.Close()
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		sourceAuthorization = req.Header.Get("Authorization")
		http.Redirect(w, req, target.URL, http.StatusFound)
	}))
	defer source.Close()

	client, err := NewClient(context.Background(), Config{
		APIKey:  "secret-bearer",
		BaseURL: source.URL,
	})
	if err != nil {
		t.Fatal(err)
	}
	var response *http.Response
	if err := client.Get(context.Background(), "/models", nil, &response); err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusFound {
		t.Fatalf("status = %d", response.StatusCode)
	}
	if sourceAuthorization != "Bearer secret-bearer" {
		t.Fatalf("source Authorization = %q", sourceAuthorization)
	}
	if calls := targetCalls.Load(); calls != 0 {
		t.Fatalf("redirect target calls = %d", calls)
	}
}

func TestConfigValidation(t *testing.T) {
	clearAWSEnvironment(t)
	tests := []struct {
		name    string
		config  Config
		message string
	}{
		{"partial static credentials", Config{AWSAccessKeyID: "access"}, "require both"},
		{"session token alone", Config{AWSSessionToken: "session"}, "may only be used"},
		{"empty bearer credential", Config{APIKey: " "}, "must not be empty"},
		{"empty static access key", Config{AWSAccessKeyID: " ", AWSSecretAccessKey: "secret"}, "non-empty"},
		{"empty static secret key", Config{AWSAccessKeyID: "access", AWSSecretAccessKey: " "}, "non-empty"},
		{"empty static session token", Config{AWSAccessKeyID: "access", AWSSecretAccessKey: "secret", AWSSessionToken: " "}, "must not be empty"},
		{"empty profile", Config{AWSProfile: " "}, "must not be empty"},
		{"multiple AWS modes", Config{AWSAccessKeyID: "access", AWSSecretAccessKey: "secret", AWSProfile: "profile"}, "ambiguous"},
		{"bearer and AWS", Config{APIKey: "bearer", AWSProfile: "profile"}, "mutually exclusive"},
		{"bearer modes", Config{APIKey: "bearer", BedrockTokenProvider: func(context.Context) (string, error) { return "", nil }}, "mutually exclusive"},
		{"skip auth and bearer", Config{SkipAuth: true, APIKey: "bearer", BaseURL: "https://bedrock.example/openai/v1"}, "cannot be combined"},
		{"empty region", Config{APIKey: "bearer", AWSRegion: " "}, "must not be empty"},
		{"empty base URL", Config{APIKey: "bearer", AWSRegion: "us-east-1", BaseURL: " "}, "must not be empty"},
		{"relative base URL", Config{APIKey: "bearer", BaseURL: "/openai/v1"}, "absolute HTTP or HTTPS"},
		{"unsupported base URL scheme", Config{APIKey: "bearer", BaseURL: "ftp://bedrock.example/openai/v1"}, "must use HTTP or HTTPS"},
		{"invalid region", Config{APIKey: "bearer", AWSRegion: "us-east-1%"}, "invalid AWS region"},
		{"invalid endpoint region", Config{APIKey: "bearer", BaseURL: "https://bedrock-mantle.us--east-1.api.aws/openai/v1"}, "invalid AWS region"},
		{"missing bearer region", Config{APIKey: "bearer"}, "AWS region is required"},
		{"missing SigV4 region", Config{AWSAccessKeyID: "access", AWSSecretAccessKey: "secret"}, "AWS region is required"},
		{
			"endpoint region mismatch",
			Config{AWSRegion: "us-west-2", AWSAccessKeyID: "access", AWSSecretAccessKey: "secret", BaseURL: "https://bedrock-mantle.us-east-1.api.aws/openai/v1"},
			"does not match",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewClient(context.Background(), test.config)
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("error = %v, want substring %q", err, test.message)
			}
		})
	}
}

func TestCredentialErrorsKeepSafeMessageAndCause(t *testing.T) {
	clearAWSEnvironment(t)
	cause := errors.New("provider detail")
	_, err := NewClient(context.Background(), Config{
		AWSRegion: "us-east-1",
		AWSCredentialsProvider: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
			return aws.Credentials{}, cause
		}),
	})
	if err == nil {
		t.Fatal("expected an error")
	}
	if got := err.Error(); got != credentialResolutionMessage {
		t.Fatalf("error message = %q", got)
	}
	if !errors.Is(err, cause) {
		t.Fatal("credential provider cause was not preserved")
	}
}
