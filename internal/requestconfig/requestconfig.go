// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package requestconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/openai/openai-go/internal"
	"github.com/openai/openai-go/internal/apierror"
	"github.com/openai/openai-go/internal/apiform"
	"github.com/openai/openai-go/internal/apiquery"
)

func getDefaultHeaders() map[string]string {
	return map[string]string{
		"User-Agent": fmt.Sprintf("OpenAI/Go %s", internal.PackageVersion),
	}
}

func getNormalizedOS() string {
	switch runtime.GOOS {
	case "ios":
		return "iOS"
	case "android":
		return "Android"
	case "darwin":
		return "MacOS"
	case "window":
		return "Windows"
	case "freebsd":
		return "FreeBSD"
	case "openbsd":
		return "OpenBSD"
	case "linux":
		return "Linux"
	default:
		return fmt.Sprintf("Other:%s", runtime.GOOS)
	}
}

func getNormalizedArchitecture() string {
	switch runtime.GOARCH {
	case "386":
		return "x32"
	case "amd64":
		return "x64"
	case "arm":
		return "arm"
	case "arm64":
		return "arm64"
	default:
		return fmt.Sprintf("other:%s", runtime.GOARCH)
	}
}

func getPlatformProperties() map[string]string {
	return map[string]string{
		"X-Stainless-Lang":            "go",
		"X-Stainless-Package-Version": internal.PackageVersion,
		"X-Stainless-OS":              getNormalizedOS(),
		"X-Stainless-Arch":            getNormalizedArchitecture(),
		"X-Stainless-Runtime":         "go",
		"X-Stainless-Runtime-Version": runtime.Version(),
	}
}

func NewRequestConfig(ctx context.Context, method string, u string, body interface{}, dst interface{}, opts ...func(*RequestConfig) error) (*RequestConfig, error) {
	var reader io.Reader

	contentType := "application/json"
	hasSerializationFunc := false

	if body, ok := body.(json.Marshaler); ok {
		content, err := body.MarshalJSON()
		if err != nil {
			return nil, err
		}
		reader = bytes.NewBuffer(content)
		hasSerializationFunc = true
	}
	if body, ok := body.(apiform.Marshaler); ok {
		var (
			content []byte
			err     error
		)
		content, contentType, err = body.MarshalMultipart()
		if err != nil {
			return nil, err
		}
		reader = bytes.NewBuffer(content)
		hasSerializationFunc = true
	}
	if body, ok := body.(apiquery.Queryer); ok {
		hasSerializationFunc = true
		params := body.URLQuery().Encode()
		if params != "" {
			u = u + "?" + params
		}
	}
	if body, ok := body.([]byte); ok {
		reader = bytes.NewBuffer(body)
		hasSerializationFunc = true
	}
	if body, ok := body.(io.Reader); ok {
		reader = body
		hasSerializationFunc = true
	}

	// Fallback to json serialization if none of the serialization functions that we expect
	// to see is present.
	if body != nil && !hasSerializationFunc {
		content, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewBuffer(content)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, err
	}
	if reader != nil {
		req.Header.Set("Content-Type", contentType)
	}

	req.Header.Set("Accept", "application/json")
	for k, v := range getDefaultHeaders() {
		req.Header.Add(k, v)
	}

	for k, v := range getPlatformProperties() {
		req.Header.Add(k, v)
	}
	cfg := RequestConfig{
		MaxRetries: 2,
		Context:    ctx,
		Request:    req,
		HTTPClient: http.DefaultClient,
		Body:       reader,
	}
	cfg.ResponseBodyInto = dst
	err = cfg.Apply(opts...)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// RequestConfig represents all the state related to one request.
//
// Editing the variables inside RequestConfig directly is unstable api. Prefer
// composing func(\*RequestConfig) error instead if possible.
type RequestConfig struct {
	MaxRetries     int
	RequestTimeout time.Duration
	Context        context.Context
	Request        *http.Request
	BaseURL        *url.URL
	HTTPClient     *http.Client
	Middlewares    []middleware
	APIKey         string
	Organization   string
	Project        string
	// If ResponseBodyInto not nil, then we will attempt to deserialize into
	// ResponseBodyInto. If Destination is a []byte, then it will return the body as
	// is.
	ResponseBodyInto interface{}
	// ResponseInto copies the \*http.Response of the corresponding request into the
	// given address
	ResponseInto **http.Response
	Body         io.Reader
}

// middleware is exactly the same type as the Middleware type found in the [option] package,
// but it is redeclared here for circular dependency issues.
type middleware = func(*http.Request, middlewareNext) (*http.Response, error)

// middlewareNext is exactly the same type as the MiddlewareNext type found in the [option] package,
// but it is redeclared here for circular dependency issues.
type middlewareNext = func(*http.Request) (*http.Response, error)

func applyMiddleware(middleware middleware, next middlewareNext) middlewareNext {
	return func(req *http.Request) (res *http.Response, err error) {
		return middleware(req, next)
	}
}

func shouldRetry(req *http.Request, res *http.Response) bool {
	// If there is no way to recover the Body, then we shouldn't retry.
	if req.Body != nil && req.GetBody == nil {
		return false
	}

	// If there is no response, that indicates that there is a connection error
	// so we retry the request.
	if res == nil {
		return true
	}

	// If the header explictly wants a retry behavior, respect that over the
	// http status code.
	if res.Header.Get("x-should-retry") == "true" {
		return true
	}
	if res.Header.Get("x-should-retry") == "false" {
		return false
	}

	return res.StatusCode == http.StatusRequestTimeout ||
		res.StatusCode == http.StatusConflict ||
		res.StatusCode == http.StatusTooManyRequests ||
		res.StatusCode >= http.StatusInternalServerError
}

func parseRetryAfterHeader(resp *http.Response) (time.Duration, bool) {
	if resp == nil {
		return 0, false
	}

	type retryData struct {
		header string
		units  time.Duration

		// custom is used when the regular algorithm failed and is optional.
		// the returned duration is used verbatim (units is not applied).
		custom func(string) (time.Duration, bool)
	}

	nop := func(string) (time.Duration, bool) { return 0, false }

	// the headers are listed in order of preference
	retries := []retryData{
		{
			header: "Retry-After-Ms",
			units:  time.Millisecond,
			custom: nop,
		},
		{
			header: "Retry-After",
			units:  time.Second,

			// retry-after values are expressed in either number of
			// seconds or an HTTP-date indicating when to try again
			custom: func(ra string) (time.Duration, bool) {
				t, err := time.Parse(time.RFC1123, ra)
				if err != nil {
					return 0, false
				}
				return time.Until(t), true
			},
		},
	}

	for _, retry := range retries {
		v := resp.Header.Get(retry.header)
		if v == "" {
			continue
		}
		if retryAfter, err := strconv.ParseFloat(v, 64); err == nil {
			return time.Duration(retryAfter * float64(retry.units)), true
		}
		if d, ok := retry.custom(v); ok {
			return d, true
		}
	}

	return 0, false
}

func retryDelay(res *http.Response, retryCount int) time.Duration {
	// If the API asks us to wait a certain amount of time (and it's a reasonable amount),
	// just do what it says.

	if retryAfterDelay, ok := parseRetryAfterHeader(res); ok && 0 <= retryAfterDelay && retryAfterDelay < time.Minute {
		return retryAfterDelay
	}

	maxDelay := 8 * time.Second
	delay := time.Duration(0.5 * float64(time.Second) * math.Pow(2, float64(retryCount)))
	if delay > maxDelay {
		delay = maxDelay
	}

	jitter := rand.Int63n(int64(delay / 4))
	delay -= time.Duration(jitter)
	return delay
}

func (cfg *RequestConfig) Execute() (err error) {
	cfg.Request.URL, err = cfg.BaseURL.Parse(strings.TrimLeft(cfg.Request.URL.String(), "/"))
	if err != nil {
		return err
	}

	if cfg.Body != nil && cfg.Request.Body == nil {
		switch body := cfg.Body.(type) {
		case *bytes.Buffer:
			b := body.Bytes()
			cfg.Request.ContentLength = int64(body.Len())
			cfg.Request.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(b)), nil }
			cfg.Request.Body, _ = cfg.Request.GetBody()
		case *bytes.Reader:
			cfg.Request.ContentLength = int64(body.Len())
			cfg.Request.GetBody = func() (io.ReadCloser, error) {
				_, err := body.Seek(0, 0)
				return io.NopCloser(body), err
			}
			cfg.Request.Body, _ = cfg.Request.GetBody()
		default:
			if rc, ok := body.(io.ReadCloser); ok {
				cfg.Request.Body = rc
			} else {
				cfg.Request.Body = io.NopCloser(body)
			}
		}
	}

	handler := cfg.HTTPClient.Do
	for i := len(cfg.Middlewares) - 1; i >= 0; i -= 1 {
		handler = applyMiddleware(cfg.Middlewares[i], handler)
	}

	var res *http.Response
	for retryCount := 0; retryCount <= cfg.MaxRetries; retryCount += 1 {
		ctx := cfg.Request.Context()
		if cfg.RequestTimeout != time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, cfg.RequestTimeout)
			defer cancel()
		}

		res, err = handler(cfg.Request.Clone(ctx))
		if ctx != nil && ctx.Err() != nil {
			return ctx.Err()
		}
		if !shouldRetry(cfg.Request, res) || retryCount >= cfg.MaxRetries {
			break
		}

		// Prepare next request and wait for the retry delay
		if cfg.Request.GetBody != nil {
			cfg.Request.Body, err = cfg.Request.GetBody()
			if err != nil {
				return err
			}
		}

		// Can't actually refresh the body, so we don't attempt to retry here
		if cfg.Request.GetBody == nil && cfg.Request.Body != nil {
			break
		}

		time.Sleep(retryDelay(res, retryCount))
	}

	// Save *http.Response if it is requested to, even if there was an error making the request. This is
	// useful in cases where you might want to debug by inspecting the response. Note that if err != nil,
	// the response should be generally be empty, but there are edge cases.
	if cfg.ResponseInto != nil {
		*cfg.ResponseInto = res
	}
	if responseBodyInto, ok := cfg.ResponseBodyInto.(**http.Response); ok {
		*responseBodyInto = res
	}

	// If there was a connection error in the final request or any other transport error,
	// return that early without trying to coerce into an APIError.
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		contents, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return err
		}

		// If there is an APIError, re-populate the response body so that debugging
		// utilities can conveniently dump the response without issue.
		res.Body = io.NopCloser(bytes.NewBuffer(contents))

		// Load the contents into the error format if it is provided.
		aerr := apierror.Error{Request: cfg.Request, Response: res, StatusCode: res.StatusCode}
		err = aerr.UnmarshalJSON(contents)
		if err != nil {
			return err
		}
		return &aerr
	}

	if cfg.ResponseBodyInto == nil {
		return nil
	}
	if _, ok := cfg.ResponseBodyInto.(**http.Response); ok {
		return nil
	}

	contents, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// If we are not json, return plaintext
	contentType := res.Header.Get("content-type")
	isJSON := strings.Contains(contentType, "application/json") || strings.Contains(contentType, "application/vnd.api+json")
	if !isJSON {
		switch dst := cfg.ResponseBodyInto.(type) {
		case *string:
			*dst = string(contents)
		case **string:
			tmp := string(contents)
			*dst = &tmp
		case *[]byte:
			*dst = contents
		default:
			return fmt.Errorf("expected destination type of 'string' or '[]byte' for responses with content-type that is not 'application/json'")
		}
		return nil
	}

	// If the response happens to be a byte array, deserialize the body as-is.
	switch dst := cfg.ResponseBodyInto.(type) {
	case *[]byte:
		*dst = contents
	}

	err = json.NewDecoder(bytes.NewReader(contents)).Decode(cfg.ResponseBodyInto)
	if err != nil {
		return fmt.Errorf("error parsing response json: %w", err)
	}

	return nil
}

func ExecuteNewRequest(ctx context.Context, method string, u string, body interface{}, dst interface{}, opts ...func(*RequestConfig) error) error {
	cfg, err := NewRequestConfig(ctx, method, u, body, dst, opts...)
	if err != nil {
		return err
	}
	return cfg.Execute()
}

func (cfg *RequestConfig) Clone(ctx context.Context) *RequestConfig {
	if cfg == nil {
		return nil
	}
	req := cfg.Request.Clone(ctx)
	var err error
	if req.Body != nil {
		req.Body, err = req.GetBody()
	}
	if err != nil {
		return nil
	}
	new := &RequestConfig{
		MaxRetries:     cfg.MaxRetries,
		RequestTimeout: cfg.RequestTimeout,
		Context:        ctx,
		Request:        req,
		BaseURL:        cfg.BaseURL,
		HTTPClient:     cfg.HTTPClient,
		Middlewares:    cfg.Middlewares,
		APIKey:         cfg.APIKey,
		Organization:   cfg.Organization,
		Project:        cfg.Project,
	}

	return new
}

func (cfg *RequestConfig) Apply(opts ...func(*RequestConfig) error) error {
	for _, opt := range opts {
		err := opt(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}
