package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

func unwrapWebhook(r *WebhookService, body []byte, headers http.Header, opts ...option.RequestOption) (*UnwrapWebhookEventUnion, error) {
	// Always perform signature verification
	err := r.VerifySignature(body, headers, opts...)
	if err != nil {
		return nil, err
	}

	res := &UnwrapWebhookEventUnion{}
	err = res.UnmarshalJSON(body)
	if err != nil {
		return res, err
	}
	return res, nil
}

func unwrapWebhookWithTolerance(r *WebhookService, body []byte, headers http.Header, tolerance time.Duration, opts ...option.RequestOption) (*UnwrapWebhookEventUnion, error) {
	err := r.VerifySignatureWithTolerance(body, headers, tolerance, opts...)
	if err != nil {
		return nil, err
	}

	res := &UnwrapWebhookEventUnion{}
	err = res.UnmarshalJSON(body)
	if err != nil {
		return res, err
	}
	return res, nil
}

func unwrapWebhookWithToleranceAndTime(r *WebhookService, body []byte, headers http.Header, tolerance time.Duration, now time.Time, opts ...option.RequestOption) (*UnwrapWebhookEventUnion, error) {
	err := r.VerifySignatureWithToleranceAndTime(body, headers, tolerance, now, opts...)
	if err != nil {
		return nil, err
	}

	res := &UnwrapWebhookEventUnion{}
	err = res.UnmarshalJSON(body)
	if err != nil {
		return res, err
	}
	return res, nil
}

func verifyWebhookSignatureWithToleranceAndTime(r *WebhookService, body []byte, headers http.Header, tolerance time.Duration, now time.Time, opts ...option.RequestOption) error {
	cfg, err := requestconfig.PreRequestOptions(slices.Concat(r.Options, opts)...)
	if err != nil {
		return err
	}
	webhookSecret := cfg.WebhookSecret

	if webhookSecret == "" {
		return errors.New("webhook secret must be provided either in the method call or configured on the client")
	}
	decodedSecret, err := decodeWebhookSecret(webhookSecret)
	if err != nil {
		return err
	}

	if headers == nil {
		return errors.New("headers are required for webhook verification")
	}

	// Extract required headers
	signatureHeader := headers.Get("webhook-signature")
	if signatureHeader == "" {
		return errors.New("missing required webhook-signature header")
	}

	timestampHeader := headers.Get("webhook-timestamp")
	if timestampHeader == "" {
		return errors.New("missing required webhook-timestamp header")
	}

	webhookID := headers.Get("webhook-id")
	if webhookID == "" {
		return errors.New("missing required webhook-id header")
	}

	// Validate timestamp to prevent replay attacks
	timestampSeconds, err := strconv.ParseInt(timestampHeader, 10, 64)
	if err != nil {
		return errors.New("invalid webhook timestamp format")
	}

	nowUnix := now.Unix()
	toleranceSeconds := int64(tolerance.Seconds())

	if nowUnix-timestampSeconds > toleranceSeconds {
		return errors.New("webhook timestamp is too old")
	}

	if timestampSeconds > nowUnix+toleranceSeconds {
		return errors.New("webhook timestamp is too new")
	}

	// Extract signatures from v1,<base64> format
	// The signature header can have multiple values, separated by spaces.
	// Each value is in the format v1,<base64>. We should accept if any match.
	var signatures []string
	for _, part := range strings.Fields(signatureHeader) {
		if strings.HasPrefix(part, "v1,") {
			signatures = append(signatures, part[3:])
		} else {
			signatures = append(signatures, part)
		}
	}

	// Create the signed payload: {webhook_id}.{timestamp}.{payload}
	signedPayload := fmt.Sprintf("%s.%s.%s", webhookID, timestampHeader, string(body))

	// Compute HMAC-SHA256 signature
	h := hmac.New(sha256.New, decodedSecret)
	h.Write([]byte(signedPayload))
	expectedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Accept if any signature matches using timing-safe comparison
	for _, signature := range signatures {
		if subtle.ConstantTimeCompare([]byte(expectedSignature), []byte(signature)) == 1 {
			return nil
		}
	}

	return errors.New("webhook signature verification failed")
}
