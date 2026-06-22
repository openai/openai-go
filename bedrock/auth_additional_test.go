package bedrock

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/openai/openai-go/v3/option"
)

func TestNewClientOptionsRejectsNilContext(t *testing.T) {
	_, err := newClientOptions(nil, Config{}, time.Now)
	if err == nil || !strings.Contains(err.Error(), "nil context") {
		t.Fatalf("error = %v", err)
	}
}

func TestProviderFinalizerRejectsConflictingOpenAIOptions(t *testing.T) {
	tests := []struct {
		name    string
		option  option.RequestOption
		message string
	}{
		{"API key", option.WithAPIKey("openai-key"), "OpenAI API key"},
		{"admin API key", option.WithAdminAPIKey("openai-admin-key"), "OpenAI API key"},
		{"base URL", option.WithBaseURL("https://other.example/v1"), "option.WithBaseURL"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transportCalls := 0
			client, err := NewClient(context.Background(), Config{
				APIKey:  "bedrock-key",
				BaseURL: "https://bedrock.example/openai/v1",
			},
				option.WithMaxRetries(0),
				test.option,
				option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					transportCalls++
					return successfulResponse(req), nil
				})}),
			)
			if err != nil {
				t.Fatal(err)
			}
			var response *http.Response
			err = client.Get(context.Background(), "/models", nil, &response)
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("error = %v, want substring %q", err, test.message)
			}
			if transportCalls != 0 {
				t.Fatalf("transport calls = %d", transportCalls)
			}
		})
	}
}

func TestSkipAuthAllowsGatewayAuthorization(t *testing.T) {
	var request *http.Request
	client, err := NewClient(context.Background(), Config{
		SkipAuth: true,
		BaseURL:  "https://gateway.example/openai/v1",
	}, option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		request = req
		return successfulResponse(req), nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	var response *http.Response
	if err := client.Get(context.Background(), "/models", nil, &response, option.WithHeader("Authorization", "Gateway credential")); err != nil {
		t.Fatal(err)
	}
	if got := request.Header.Get("Authorization"); got != "Gateway credential" {
		t.Fatalf("Authorization = %q", got)
	}
}

func TestEnvironmentBearerRefreshesAndFailsSafelyWhenRemoved(t *testing.T) {
	clearAWSEnvironment(t)
	t.Setenv("AWS_REGION", "us-east-1")
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "first-token")

	var authorizations []string
	client, err := NewClient(context.Background(), Config{},
		option.WithMaxRetries(0),
		option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			authorizations = append(authorizations, req.Header.Get("Authorization"))
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
	if err := os.Setenv("AWS_BEARER_TOKEN_BEDROCK", "second-token"); err != nil {
		t.Fatal(err)
	}
	if err := client.Get(context.Background(), "/models", nil, &response); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(authorizations, ","); got != "Bearer first-token,Bearer second-token" {
		t.Fatalf("authorizations = %q", got)
	}

	if err := os.Unsetenv("AWS_BEARER_TOKEN_BEDROCK"); err != nil {
		t.Fatal(err)
	}
	err = client.Get(context.Background(), "/models", nil, &response)
	if err == nil || !strings.Contains(err.Error(), "Failed to resolve a bearer credential") {
		t.Fatalf("error = %v", err)
	}
}

func TestBearerProviderFailuresAreSafe(t *testing.T) {
	cause := errors.New("provider internals")
	tests := []struct {
		name     string
		provider TokenProvider
		message  string
		cause    error
	}{
		{
			name: "provider error",
			provider: func(context.Context) (string, error) {
				return "", cause
			},
			message: "Failed to resolve a bearer credential",
			cause:   cause,
		},
		{
			name: "empty credential",
			provider: func(context.Context) (string, error) {
				return " ", nil
			},
			message: "must return a non-empty string",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transportCalls := 0
			client, err := NewClient(context.Background(), Config{
				BaseURL:              "https://bedrock.example/openai/v1",
				BedrockTokenProvider: test.provider,
			},
				option.WithMaxRetries(0),
				option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					transportCalls++
					return successfulResponse(req), nil
				})}),
			)
			if err != nil {
				t.Fatal(err)
			}
			var response *http.Response
			err = client.Get(context.Background(), "/models", nil, &response)
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("error = %v, want substring %q", err, test.message)
			}
			if test.cause != nil && !errors.Is(err, test.cause) {
				t.Fatal("provider cause was not preserved")
			}
			if transportCalls != 0 {
				t.Fatalf("transport calls = %d", transportCalls)
			}
		})
	}
}

func TestMissingDefaultCredentialsUsesActionableError(t *testing.T) {
	clearAWSEnvironment(t)
	credentialsPath, configPath := writeAWSFiles(t, "", "")
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", credentialsPath)
	t.Setenv("AWS_CONFIG_FILE", configPath)

	_, err := NewClient(context.Background(), Config{AWSRegion: "us-east-1"})
	if err == nil || err.Error() != missingCredentialsMessage {
		t.Fatalf("error = %v", err)
	}
}

func TestMalformedAWSConfigUsesSafeError(t *testing.T) {
	clearAWSEnvironment(t)
	dir := t.TempDir()
	t.Setenv("AWS_CONFIG_FILE", dir)
	credentialsPath, _ := writeAWSFiles(t, "", "")
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", credentialsPath)

	_, err := NewClient(context.Background(), Config{AWSRegion: "us-east-1", AWSProfile: "broken"})
	if err == nil || err.Error() != credentialResolutionMessage {
		t.Fatalf("error = %v", err)
	}
}

func TestCanonicalEndpointInfersSigningRegion(t *testing.T) {
	clearAWSEnvironment(t)
	var request *http.Request
	client, err := NewClient(context.Background(), Config{
		BaseURL:            "https://bedrock-mantle.eu-west-1.api.aws/openai/v1",
		AWSAccessKeyID:     "access-key",
		AWSSecretAccessKey: "secret-key",
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
	if got := request.Header.Get("Authorization"); !strings.Contains(got, "/eu-west-1/bedrock-mantle/aws4_request") {
		t.Fatalf("Authorization = %q", got)
	}
}

type signerFunc func(context.Context, aws.Credentials, *http.Request, string, string, string, time.Time, ...func(*v4.SignerOptions)) error

func (f signerFunc) SignHTTP(
	ctx context.Context,
	credentials aws.Credentials,
	req *http.Request,
	payloadHash string,
	service string,
	region string,
	signingTime time.Time,
	options ...func(*v4.SignerOptions),
) error {
	return f(ctx, credentials, req, payloadHash, service, region, signingTime, options...)
}

func TestSigV4MiddlewareRejectsUnsafeRequests(t *testing.T) {
	baseURL, _ := url.Parse("https://bedrock-mantle.us-east-1.api.aws/openai/v1/")
	validCredentials := aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
		return aws.Credentials{AccessKeyID: "access-key", SecretAccessKey: "secret-key"}, nil
	})
	tests := []struct {
		name    string
		request func() *http.Request
		config  aws.Config
		message string
	}{
		{
			name: "cross origin",
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "https://attacker.example/models", nil)
				return req
			},
			config:  aws.Config{Region: "us-east-1", Credentials: validCredentials},
			message: "origin other than",
		},
		{
			name: "custom authorization",
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "https://bedrock-mantle.us-east-1.api.aws/openai/v1/models", nil)
				req.Header.Set("Authorization", "Bearer custom")
				return req
			},
			config:  aws.Config{Region: "us-east-1", Credentials: validCredentials},
			message: "custom `Authorization` header",
		},
		{
			name: "endpoint region mismatch",
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "https://bedrock-mantle.us-east-1.api.aws/openai/v1/models", nil)
				return req
			},
			config:  aws.Config{Region: "us-west-2", Credentials: validCredentials},
			message: "does not match",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transportCalls := 0
			_, err := sigV4Middleware(baseURL, test.config, v4.NewSigner(), time.Now)(test.request(), func(req *http.Request) (*http.Response, error) {
				transportCalls++
				return successfulResponse(req), nil
			})
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("error = %v, want substring %q", err, test.message)
			}
			if transportCalls != 0 {
				t.Fatalf("transport calls = %d", transportCalls)
			}
		})
	}
}

func TestSigV4MiddlewareCredentialAndSignerFailures(t *testing.T) {
	baseURL, _ := url.Parse("https://bedrock-mantle.us-east-1.api.aws/openai/v1/")
	cause := errors.New("internal failure")
	tests := []struct {
		name    string
		config  aws.Config
		signer  httpSigner
		message string
	}{
		{
			name: "credential provider failure",
			config: aws.Config{Region: "us-east-1", Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
				return aws.Credentials{}, cause
			})},
			signer:  v4.NewSigner(),
			message: credentialResolutionMessage,
		},
		{
			name: "empty credentials",
			config: aws.Config{Region: "us-east-1", Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
				return aws.Credentials{}, nil
			})},
			signer:  v4.NewSigner(),
			message: credentialResolutionMessage,
		},
		{
			name: "signer failure",
			config: aws.Config{Region: "us-east-1", Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
				return aws.Credentials{AccessKeyID: "access-key", SecretAccessKey: "secret-key"}, nil
			})},
			signer: signerFunc(func(context.Context, aws.Credentials, *http.Request, string, string, string, time.Time, ...func(*v4.SignerOptions)) error {
				return cause
			}),
			message: "Failed to sign",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "https://bedrock-mantle.us-east-1.api.aws/openai/v1/responses", strings.NewReader("body"))
			contentLength := req.ContentLength
			_, err := sigV4Middleware(baseURL, test.config, test.signer, time.Now)(req, func(req *http.Request) (*http.Response, error) {
				t.Fatal("transport must not be called")
				return nil, nil
			})
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("error = %v, want substring %q", err, test.message)
			}
			if test.name != "empty credentials" && !errors.Is(err, cause) {
				t.Fatal("underlying cause was not preserved")
			}
			if req.ContentLength != contentLength {
				t.Fatalf("ContentLength = %d, want %d", req.ContentLength, contentLength)
			}
		})
	}
}

type testReadCloser struct {
	reader   io.Reader
	readErr  error
	closeErr error
}

func (r *testReadCloser) Read(p []byte) (int, error) {
	if r.readErr != nil {
		return 0, r.readErr
	}
	return r.reader.Read(p)
}

func (r *testReadCloser) Close() error { return r.closeErr }

func TestMaterializeReplayableBodyFailuresAndReset(t *testing.T) {
	readCause := errors.New("read failure")
	closeCause := errors.New("close failure")
	tests := []struct {
		name  string
		body  io.ReadCloser
		cause error
	}{
		{"read failure", &testReadCloser{reader: strings.NewReader("body"), readErr: readCause}, readCause},
		{"close failure", &testReadCloser{reader: strings.NewReader("body"), closeErr: closeCause}, closeCause},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := &http.Request{
				Body: test.body,
				GetBody: func() (io.ReadCloser, error) {
					return io.NopCloser(strings.NewReader("body")), nil
				},
			}
			_, err := materializeReplayableBody(req)
			if err == nil || !errors.Is(err, test.cause) {
				t.Fatalf("error = %v", err)
			}
		})
	}

	req, _ := http.NewRequest(http.MethodPost, "https://bedrock.example/responses", strings.NewReader("replay me"))
	body, err := materializeReplayableBody(req)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "replay me" {
		t.Fatalf("body = %q", body)
	}
	replayed, err := req.GetBody()
	if err != nil {
		t.Fatal(err)
	}
	defer replayed.Close()
	replayedBody, err := io.ReadAll(replayed)
	if err != nil {
		t.Fatal(err)
	}
	if string(replayedBody) != "replay me" {
		t.Fatalf("replayed body = %q", replayedBody)
	}
}

func TestURLHelpers(t *testing.T) {
	if normalizeBaseURL(nil) != nil {
		t.Fatal("normalizeBaseURL(nil) must return nil")
	}
	if !sameBaseURL(nil, nil) || sameBaseURL(nil, &url.URL{}) {
		t.Fatal("unexpected nil base URL comparison")
	}
	httpDefault, _ := url.Parse("http://example.com/path")
	httpExplicit, _ := url.Parse("http://example.com:80/other")
	if !sameOrigin(httpDefault, httpExplicit) {
		t.Fatal("default and explicit HTTP ports should match")
	}
	httpsDefault, _ := url.Parse("https://example.com/path")
	httpsExplicit, _ := url.Parse("https://example.com:443/other")
	if !sameOrigin(httpsDefault, httpsExplicit) {
		t.Fatal("default and explicit HTTPS ports should match")
	}
	ftpURL, _ := url.Parse("ftp://example.com/path")
	if got := effectivePort(ftpURL); got != "" {
		t.Fatalf("FTP effective port = %q", got)
	}
}
