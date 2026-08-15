package auth

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

const tokenExchangeMaxRetryDelay = 60 * time.Second

func exchangeToken(
	ctx context.Context,
	httpClient HTTPDoer,
	url string,
	body []byte,
	maxRetries int,
) (*http.Response, error) {
	for retry := 0; ; retry++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("failed to create token exchange request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err == nil && (resp == nil || resp.Body == nil) {
			return nil, fmt.Errorf("token exchange returned an invalid HTTP response")
		}
		if !shouldRetryTokenExchange(resp, err) || retry >= maxRetries {
			if err != nil {
				if resp != nil && resp.Body != nil {
					_ = resp.Body.Close()
				}
				return nil, fmt.Errorf("failed to exchange token: %w", err)
			}
			return resp, nil
		}

		if resp != nil && resp.Body != nil {
			_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
			_ = resp.Body.Close()
		}
		if err := waitForTokenExchangeRetry(ctx, tokenExchangeRetryDelay(resp, retry)); err != nil {
			return nil, err
		}
	}
}

func readTokenExchangeResponse(reader io.Reader, maxBytes int64) ([]byte, error) {
	responseBody := reader
	if maxBytes > 0 {
		responseBody = io.LimitReader(reader, maxBytes+1)
	}
	body, err := io.ReadAll(responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to read token exchange response: %w", err)
	}
	if maxBytes > 0 && int64(len(body)) > maxBytes {
		return nil, fmt.Errorf("token exchange response exceeded %d bytes", maxBytes)
	}
	return body, nil
}

func shouldRetryTokenExchange(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	return resp != nil && (resp.StatusCode == http.StatusRequestTimeout ||
		resp.StatusCode == http.StatusConflict ||
		resp.StatusCode == http.StatusTooManyRequests ||
		resp.StatusCode >= http.StatusInternalServerError)
}

func tokenExchangeRetryDelay(resp *http.Response, retry int) time.Duration {
	if resp != nil {
		if value := resp.Header.Get("Retry-After"); value != "" {
			if seconds, err := strconv.ParseFloat(value, 64); err == nil && !math.IsNaN(seconds) {
				switch {
				case seconds <= 0:
					return 0
				case math.IsInf(seconds, 1) || seconds >= tokenExchangeMaxRetryDelay.Seconds():
					return tokenExchangeMaxRetryDelay
				default:
					return time.Duration(seconds * float64(time.Second))
				}
			}
			if retryAt, err := http.ParseTime(value); err == nil {
				return min(max(0, time.Until(retryAt)), tokenExchangeMaxRetryDelay)
			}
		}
	}
	return min(time.Duration(1<<retry)*250*time.Millisecond, tokenExchangeMaxRetryDelay)
}

func waitForTokenExchangeRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
