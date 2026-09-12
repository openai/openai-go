package webhooks_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/webhooks"
)

const webhookHelperSecret = "synthetic-webhook-helper-secret"

func signedWebhookHelperHeaders(t *testing.T, body []byte, timestamp int64) http.Header {
	t.Helper()
	timestampString := strconv.FormatInt(timestamp, 10)
	h := hmac.New(sha256.New, []byte(webhookHelperSecret))
	if _, err := h.Write([]byte("wh_helper." + timestampString + "." + string(body))); err != nil {
		t.Fatal(err)
	}
	return http.Header{
		"Webhook-Id":        {"wh_helper"},
		"Webhook-Timestamp": {timestampString},
		"Webhook-Signature": {"v1," + base64.StdEncoding.EncodeToString(h.Sum(nil))},
	}
}

func TestWebhookUnwrapHelpersVerifyBeforeDecode(t *testing.T) {
	service := webhooks.NewWebhookService(option.WithWebhookSecret(webhookHelperSecret))
	now := time.Now()
	tests := []struct {
		name string
		call func([]byte, http.Header) (*webhooks.UnwrapWebhookEventUnion, error)
	}{
		{"Unwrap", func(body []byte, headers http.Header) (*webhooks.UnwrapWebhookEventUnion, error) {
			return service.Unwrap(body, headers)
		}},
		{"UnwrapWithTolerance", func(body []byte, headers http.Header) (*webhooks.UnwrapWebhookEventUnion, error) {
			return service.UnwrapWithTolerance(body, headers, 5*time.Minute)
		}},
		{"UnwrapWithToleranceAndTime", func(body []byte, headers http.Header) (*webhooks.UnwrapWebhookEventUnion, error) {
			return service.UnwrapWithToleranceAndTime(body, headers, 5*time.Minute, now)
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run("valid raw body", func(t *testing.T) {
				body := []byte(`{"id": "evt_helper", "object": "event", "created_at": 1700000000, "type": "response.completed", "data": {"id": "resp_helper"}}`)
				headers := signedWebhookHelperHeaders(t, body, now.Unix())
				originalBody, originalHeaders := bytes.Clone(body), headers.Clone()
				result, err := test.call(body, headers)
				if err != nil {
					t.Fatal(err)
				}
				if result == nil || result.ID != "evt_helper" || result.Type != "response.completed" || result.Data.ID != "resp_helper" || result.RawJSON() != string(body) {
					t.Fatalf("unexpected decoded event: %#v", result)
				}
				if !bytes.Equal(body, originalBody) || !reflect.DeepEqual(headers, originalHeaders) {
					t.Fatal("unwrap mutated the caller's body or headers")
				}
			})
			t.Run("invalid signature and payload shape", func(t *testing.T) {
				body := []byte(`[]`)
				headers := signedWebhookHelperHeaders(t, []byte(`{}`), now.Unix())
				result, err := test.call(body, headers)
				if result != nil || err == nil || err.Error() != "webhook signature verification failed" {
					t.Fatalf("result = %#v, error = %v; want verification failure before parsing", result, err)
				}
			})
			t.Run("valid signature and invalid payload shape", func(t *testing.T) {
				body := []byte(`[]`)
				headers := signedWebhookHelperHeaders(t, body, now.Unix())
				if err := service.VerifySignatureWithToleranceAndTime(body, headers, 5*time.Minute, now); err != nil {
					t.Fatalf("raw array payload should still have a valid signature: %v", err)
				}
				result, err := test.call(body, headers)
				var typeError *json.UnmarshalTypeError
				if result == nil || !errors.As(err, &typeError) {
					t.Fatalf("result = %#v, error = %v; want allocated result and JSON type error", result, err)
				}
			})
		})
	}
}

func TestWebhookSignatureHelperValidationOrder(t *testing.T) {
	body := []byte(`{}`)
	now := time.Unix(1700000000, 0)
	for _, test := range []struct {
		name    string
		secret  string
		headers http.Header
		wantErr string
	}{
		{"missing secret", "", nil, "webhook secret must be provided either in the method call or configured on the client"},
		{"invalid secret before headers", "whsec_", nil, "invalid webhook secret: decoded secret must be at least 24 bytes"},
		{"nil headers", webhookHelperSecret, nil, "headers are required for webhook verification"},
		{"missing signature", webhookHelperSecret, http.Header{}, "missing required webhook-signature header"},
		{"missing timestamp", webhookHelperSecret, http.Header{"Webhook-Signature": {"invalid"}}, "missing required webhook-timestamp header"},
		{"missing id before timestamp parsing", webhookHelperSecret, http.Header{"Webhook-Signature": {"invalid"}, "Webhook-Timestamp": {"invalid"}}, "missing required webhook-id header"},
		{"invalid timestamp before signature", webhookHelperSecret, http.Header{"Webhook-Signature": {"invalid"}, "Webhook-Timestamp": {"invalid"}, "Webhook-Id": {"wh_helper"}}, "invalid webhook timestamp format"},
	} {
		t.Run(test.name, func(t *testing.T) {
			service := webhooks.NewWebhookService(option.WithWebhookSecret(test.secret))
			err := service.VerifySignatureWithToleranceAndTime(body, test.headers, time.Minute, now)
			if err == nil || err.Error() != test.wantErr {
				t.Fatalf("error = %v, want %q", err, test.wantErr)
			}
		})
	}
}

func TestWebhookSignatureHelperToleranceBoundaries(t *testing.T) {
	service := webhooks.NewWebhookService(option.WithWebhookSecret(webhookHelperSecret))
	body := []byte(`{}`)
	now := time.Unix(1700000000, 750000000)
	for _, test := range []struct {
		name      string
		offset    int64
		tolerance time.Duration
		wantErr   string
	}{
		{"zero tolerance", 0, 0, ""},
		{"old boundary", -5, 5 * time.Second, ""},
		{"future boundary", 5, 5 * time.Second, ""},
		{"too old", -6, 5 * time.Second, "webhook timestamp is too old"},
		{"too new", 6, 5 * time.Second, "webhook timestamp is too new"},
		{"fractional tolerance boundary", -1, 1500 * time.Millisecond, ""},
		{"fractional tolerance truncation", -2, 1500 * time.Millisecond, "webhook timestamp is too old"},
	} {
		t.Run(test.name, func(t *testing.T) {
			headers := signedWebhookHelperHeaders(t, body, now.Unix()+test.offset)
			err := service.VerifySignatureWithToleranceAndTime(body, headers, test.tolerance, now)
			if test.wantErr == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || err.Error() != test.wantErr {
				t.Fatalf("error = %v, want %q", err, test.wantErr)
			}
		})
	}
}
