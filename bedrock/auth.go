package bedrock

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

const (
	bedrockService = "bedrock-mantle"

	missingRegionMessage        = "Bedrock requires an AWS region. Pass `AWSRegion` in `bedrock.Config`, or set `AWS_REGION` or `AWS_DEFAULT_REGION`."
	missingCredentialsMessage   = "Could not find credentials for Bedrock. Pass a bearer credential or AWS credentials in `bedrock.Config`, set `AWS_BEARER_TOKEN_BEDROCK`, or configure the default AWS credential chain."
	credentialResolutionMessage = "Failed to resolve AWS credentials for Bedrock. Verify your AWS profile, environment variables, or runtime identity configuration and try again."
	nonReplayableBodyMessage    = "Bedrock SigV4 authentication requires a replayable request body. Buffer the body before sending or use bearer authentication."
)

var canonicalBedrockHost = regexp.MustCompile(`(?i)^bedrock-mantle\.([a-z0-9-]+)\.api\.aws$`)

type authMode int

const (
	authModeSkip authMode = iota
	authModeBearer
	authModeSigV4
)

type resolvedConfig struct {
	mode       authMode
	region     string
	baseURL    *url.URL
	middleware option.Middleware
}

type safeError struct {
	message string
	cause   error
}

func (e *safeError) Error() string { return e.message }
func (e *safeError) Unwrap() error { return e.cause }

func newClientOptions(
	ctx context.Context,
	cfg Config,
	now func() time.Time,
	userOpts ...option.RequestOption,
) ([]option.RequestOption, error) {
	if ctx == nil {
		return nil, errors.New("bedrock: nil context")
	}

	resolved, err := resolveConfig(ctx, cfg, now)
	if err != nil {
		return nil, err
	}

	opts := []option.RequestOption{
		requestconfig.WithEnvironmentDefaultsDisabled(),
		option.WithBaseURL(resolved.baseURL.String()),
	}
	opts = append(opts, userOpts...)
	opts = append(opts, requestconfig.WithRequestFinalizer(func(rc *requestconfig.RequestConfig) error {
		if rc.APIKey != "" || rc.AdminAPIKey != "" {
			return errors.New("Bedrock provider authentication cannot be combined with an OpenAI API key. Configure authentication in `bedrock.Config`.")
		}
		if !sameBaseURL(rc.BaseURL, resolved.baseURL) {
			return errors.New("Bedrock provider routing cannot be overridden with `option.WithBaseURL`. Configure `BaseURL` in `bedrock.Config`.")
		}

		if resolved.mode == authModeSigV4 && rc.CustomHTTPDoer == nil && rc.HTTPClient != nil {
			client := *rc.HTTPClient
			client.CheckRedirect = func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			}
			rc.HTTPClient = &client
		}

		if resolved.middleware != nil {
			return option.WithMiddleware(resolved.middleware).Apply(rc)
		}
		return nil
	}))
	return opts, nil
}

func resolveConfig(ctx context.Context, cfg Config, now func() time.Time) (resolvedConfig, error) {
	mode, tokenProvider, explicitAWS, err := resolveAuthMode(cfg)
	if err != nil {
		return resolvedConfig{}, err
	}

	region, err := resolveOptionalConfigValue("AWSRegion", cfg.AWSRegion, "AWS_REGION", "AWS_DEFAULT_REGION")
	if err != nil {
		return resolvedConfig{}, err
	}
	baseURLValue, err := resolveOptionalConfigValue("BaseURL", cfg.BaseURL, "AWS_BEDROCK_BASE_URL")
	if err != nil {
		return resolvedConfig{}, err
	}

	baseURL, err := parseBaseURL(baseURLValue)
	if err != nil {
		return resolvedConfig{}, err
	}
	if baseURL != nil {
		region, err = reconcileEndpointRegion(baseURL, region)
		if err != nil {
			return resolvedConfig{}, err
		}
	}

	resolved := resolvedConfig{mode: mode, region: region, baseURL: baseURL}
	var awsCfg aws.Config
	switch mode {
	case authModeSkip:
	case authModeBearer:
	case authModeSigV4:
		awsCfg, err = loadAWSConfig(ctx, cfg, region)
		if err != nil {
			return resolvedConfig{}, err
		}
		if region == "" {
			region = strings.TrimSpace(awsCfg.Region)
		}
		if region == "" {
			return resolvedConfig{}, errors.New(missingRegionMessage)
		}
		if baseURL != nil {
			if _, err := reconcileEndpointRegion(baseURL, region); err != nil {
				return resolvedConfig{}, err
			}
		}
		awsCfg.Region = region
		if err := verifyAWSCredentials(ctx, awsCfg, explicitAWS); err != nil {
			return resolvedConfig{}, err
		}
		resolved.region = region
	}

	if resolved.baseURL == nil {
		if region == "" {
			return resolvedConfig{}, errors.New(missingRegionMessage)
		}
		resolved.baseURL, _ = url.Parse(fmt.Sprintf("https://bedrock-mantle.%s.api.aws/openai/v1/", region))
	}
	resolved.baseURL = normalizeBaseURL(resolved.baseURL)

	if mode == authModeBearer {
		resolved.middleware = bearerMiddleware(resolved.baseURL, tokenProvider)
	}
	if mode == authModeSigV4 {
		resolved.middleware = sigV4Middleware(resolved.baseURL, awsCfg, v4.NewSigner(), now)
	}

	return resolved, nil
}

func resolveAuthMode(cfg Config) (authMode, TokenProvider, bool, error) {
	if cfg.APIKey != "" && strings.TrimSpace(cfg.APIKey) == "" {
		return 0, nil, false, errors.New("The Bedrock bearer credential must not be empty.")
	}
	if cfg.APIKey != "" && cfg.BedrockTokenProvider != nil {
		return 0, nil, false, errors.New("The `APIKey` and `BedrockTokenProvider` options are mutually exclusive. Configure only one.")
	}

	hasAccessKey := cfg.AWSAccessKeyID != ""
	hasSecretKey := cfg.AWSSecretAccessKey != ""
	hasSessionToken := cfg.AWSSessionToken != ""
	if hasAccessKey != hasSecretKey || (hasSessionToken && !hasAccessKey) {
		return 0, nil, false, errors.New("Static AWS credentials require both `AWSAccessKeyID` and `AWSSecretAccessKey`. An `AWSSessionToken` may only be used with both.")
	}
	if hasAccessKey && (strings.TrimSpace(cfg.AWSAccessKeyID) == "" || strings.TrimSpace(cfg.AWSSecretAccessKey) == "") {
		return 0, nil, false, errors.New("Static AWS credentials require non-empty `AWSAccessKeyID` and `AWSSecretAccessKey` values.")
	}
	if hasSessionToken && strings.TrimSpace(cfg.AWSSessionToken) == "" {
		return 0, nil, false, errors.New("A static AWS `AWSSessionToken` must not be empty when provided.")
	}
	profile := strings.TrimSpace(cfg.AWSProfile)
	if cfg.AWSProfile != "" && profile == "" {
		return 0, nil, false, errors.New("The Bedrock AWS `AWSProfile` must not be empty.")
	}

	awsModes := 0
	if hasAccessKey {
		awsModes++
	}
	if profile != "" {
		awsModes++
	}
	if cfg.AWSCredentialsProvider != nil {
		awsModes++
	}
	if awsModes > 1 {
		return 0, nil, false, errors.New("Bedrock authentication is ambiguous. Configure exactly one explicit AWS mode: static credentials, profile, or credentials provider.")
	}

	hasBearer := cfg.APIKey != "" || cfg.BedrockTokenProvider != nil
	if hasBearer && awsModes != 0 {
		return 0, nil, false, errors.New("Bearer and AWS credential authentication are mutually exclusive. Configure exactly one explicit mode: bearer credential, static AWS credentials, profile, or credentials provider.")
	}
	if cfg.SkipAuth && (hasBearer || awsModes != 0) {
		return 0, nil, false, errors.New("`SkipAuth` cannot be combined with explicit Bedrock authentication options.")
	}
	if cfg.SkipAuth {
		return authModeSkip, nil, false, nil
	}
	if cfg.BedrockTokenProvider != nil {
		return authModeBearer, cfg.BedrockTokenProvider, false, nil
	}
	if cfg.APIKey != "" {
		token := cfg.APIKey
		return authModeBearer, func(context.Context) (string, error) { return token, nil }, false, nil
	}
	if awsModes != 0 {
		return authModeSigV4, nil, true, nil
	}
	if strings.TrimSpace(os.Getenv("AWS_BEARER_TOKEN_BEDROCK")) != "" {
		return authModeBearer, func(context.Context) (string, error) {
			token := os.Getenv("AWS_BEARER_TOKEN_BEDROCK")
			if strings.TrimSpace(token) == "" {
				return "", errors.New(missingCredentialsMessage)
			}
			return token, nil
		}, false, nil
	}
	return authModeSigV4, nil, false, nil
}

func resolveOptionalConfigValue(name string, explicit string, environment ...string) (string, error) {
	if explicit != "" {
		value := strings.TrimSpace(explicit)
		if value == "" {
			return "", fmt.Errorf("the Bedrock `%s` option must not be empty", name)
		}
		return value, nil
	}
	for _, envName := range environment {
		if value := strings.TrimSpace(os.Getenv(envName)); value != "" {
			return value, nil
		}
	}
	return "", nil
}

func parseBaseURL(value string) (*url.URL, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || parsed.User != nil {
		return nil, errors.New("the Bedrock `BaseURL` must be an absolute HTTP or HTTPS URL without user information")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, errors.New("the Bedrock `BaseURL` must use HTTP or HTTPS")
	}
	return normalizeBaseURL(parsed), nil
}

func normalizeBaseURL(value *url.URL) *url.URL {
	if value == nil {
		return nil
	}
	copy := *value
	if copy.Path == "" {
		copy.Path = "/"
	} else if !strings.HasSuffix(copy.Path, "/") {
		copy.Path += "/"
	}
	return &copy
}

func reconcileEndpointRegion(baseURL *url.URL, region string) (string, error) {
	match := canonicalBedrockHost.FindStringSubmatch(baseURL.Hostname())
	if len(match) != 2 {
		return region, nil
	}
	endpointRegion := strings.ToLower(match[1])
	if region != "" && !strings.EqualFold(endpointRegion, region) {
		return "", fmt.Errorf("the Bedrock endpoint region `%s` does not match the SigV4 region `%s`", endpointRegion, region)
	}
	if region == "" {
		return endpointRegion, nil
	}
	return region, nil
}

func loadAWSConfig(ctx context.Context, cfg Config, region string) (aws.Config, error) {
	loadOptions := make([]func(*awsconfig.LoadOptions) error, 0, 3)
	if region != "" {
		loadOptions = append(loadOptions, awsconfig.WithRegion(region))
	}
	if profile := strings.TrimSpace(cfg.AWSProfile); profile != "" {
		loadOptions = append(loadOptions, awsconfig.WithSharedConfigProfile(profile))
	}

	var explicitProvider aws.CredentialsProvider
	if cfg.AWSAccessKeyID != "" {
		credentials := aws.Credentials{
			AccessKeyID:     cfg.AWSAccessKeyID,
			SecretAccessKey: cfg.AWSSecretAccessKey,
			SessionToken:    cfg.AWSSessionToken,
			Source:          "bedrock.Config",
		}
		explicitProvider = aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
			return credentials, nil
		})
	} else if cfg.AWSCredentialsProvider != nil {
		explicitProvider = cfg.AWSCredentialsProvider
	}
	if explicitProvider != nil {
		loadOptions = append(loadOptions, awsconfig.WithCredentialsProvider(explicitProvider))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, loadOptions...)
	if err != nil {
		return aws.Config{}, &safeError{message: credentialResolutionMessage, cause: err}
	}
	if awsCfg.Credentials == nil {
		return aws.Config{}, errors.New(missingCredentialsMessage)
	}
	if _, ok := awsCfg.Credentials.(*aws.CredentialsCache); !ok {
		awsCfg.Credentials = aws.NewCredentialsCache(awsCfg.Credentials)
	}
	return awsCfg, nil
}

func verifyAWSCredentials(ctx context.Context, awsCfg aws.Config, explicitAWS bool) error {
	if _, err := awsCfg.Credentials.Retrieve(ctx); err != nil {
		message := credentialResolutionMessage
		if !explicitAWS {
			message = missingCredentialsMessage
		}
		return &safeError{message: message, cause: err}
	}
	return nil
}

func bearerMiddleware(baseURL *url.URL, provider TokenProvider) option.Middleware {
	return func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		if err := validateProviderRequest(req, baseURL); err != nil {
			return nil, err
		}
		if req.Header.Get("Authorization") != "" {
			return nil, errors.New("Bedrock provider authentication cannot be combined with a custom `Authorization` header.")
		}

		token, err := provider(req.Context())
		if err != nil {
			return nil, &safeError{message: "Failed to resolve a bearer credential for Bedrock.", cause: err}
		}
		if strings.TrimSpace(token) == "" {
			return nil, errors.New("The Bedrock bearer credential provider must return a non-empty string.")
		}
		req.Header.Set("Authorization", "Bearer "+token)
		return next(req)
	}
}

type httpSigner interface {
	SignHTTP(context.Context, aws.Credentials, *http.Request, string, string, string, time.Time, ...func(*v4.SignerOptions)) error
}

func sigV4Middleware(baseURL *url.URL, cfg aws.Config, signer httpSigner, now func() time.Time) option.Middleware {
	return func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		if err := validateProviderRequest(req, baseURL); err != nil {
			return nil, err
		}
		if req.Header.Get("Authorization") != "" {
			return nil, errors.New("Bedrock provider authentication cannot be combined with a custom `Authorization` header.")
		}
		if _, err := reconcileEndpointRegion(req.URL, cfg.Region); err != nil {
			return nil, err
		}

		body, err := materializeReplayableBody(req)
		if err != nil {
			return nil, err
		}
		credentials, err := cfg.Credentials.Retrieve(req.Context())
		if err != nil {
			return nil, &safeError{message: credentialResolutionMessage, cause: err}
		}
		if strings.TrimSpace(credentials.AccessKeyID) == "" || strings.TrimSpace(credentials.SecretAccessKey) == "" {
			return nil, errors.New(credentialResolutionMessage)
		}

		req.Method = strings.ToUpper(req.Method)
		req.Header.Del("X-Amz-Date")
		req.Header.Del("X-Amz-Security-Token")
		req.Header.Del("X-Amz-Content-Sha256")
		payloadHash := sha256.Sum256(body)
		encodedHash := hex.EncodeToString(payloadHash[:])
		req.Header.Set("X-Amz-Content-Sha256", encodedHash)
		// Content-Length is transmitted by net/http but does not need to be part of
		// SigV4's signed-header set. Temporarily hide it from the AWS signer so the
		// signature matches the shared cross-SDK fixture, then restore the exact
		// wire length before the request is sent.
		contentLength := req.ContentLength
		req.ContentLength = -1
		signErr := signer.SignHTTP(req.Context(), credentials, req, encodedHash, bedrockService, cfg.Region, now().UTC())
		req.ContentLength = contentLength
		if signErr != nil {
			return nil, &safeError{message: "Failed to sign the Bedrock request.", cause: signErr}
		}
		return next(req)
	}
}

func materializeReplayableBody(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return nil, nil
	}
	if req.GetBody == nil {
		return nil, errors.New(nonReplayableBodyMessage)
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, &safeError{message: nonReplayableBodyMessage, cause: err}
	}
	if err := req.Body.Close(); err != nil {
		return nil, &safeError{message: nonReplayableBodyMessage, cause: err}
	}
	body = bytes.Clone(body)
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}
	req.ContentLength = int64(len(body))
	return body, nil
}

func validateProviderRequest(req *http.Request, baseURL *url.URL) error {
	if req.URL == nil || !sameOrigin(req.URL, baseURL) {
		return errors.New("Bedrock provider authentication cannot send credentials to an origin other than the configured provider URL.")
	}
	return nil
}

func sameBaseURL(left, right *url.URL) bool {
	if left == nil || right == nil {
		return left == right
	}
	return normalizeBaseURL(left).String() == normalizeBaseURL(right).String()
}

func sameOrigin(left, right *url.URL) bool {
	if left == nil || right == nil || !strings.EqualFold(left.Scheme, right.Scheme) || !strings.EqualFold(left.Hostname(), right.Hostname()) {
		return false
	}
	return effectivePort(left) == effectivePort(right)
}

func effectivePort(value *url.URL) string {
	if port := value.Port(); port != "" {
		return port
	}
	if strings.EqualFold(value.Scheme, "https") {
		return "443"
	}
	if strings.EqualFold(value.Scheme, "http") {
		return "80"
	}
	return ""
}
