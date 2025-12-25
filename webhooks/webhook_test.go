package webhooks_test

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/Nordlys-Labs/openai-go/v3"
	"github.com/Nordlys-Labs/openai-go/v3/option"
)

// Standardized test constants (matches TypeScript implementation)
const (
	testPayload   = `{"id": "evt_685c059ae3a481909bdc86819b066fb6", "object": "event", "created_at": 1750861210, "type": "response.completed", "data": {"id": "resp_123"}}`
	testSecret    = "whsec_RdvaYFYUXuIFuEbvZHwMfYFhUf7aMYjYcmM24+Aj40c="
	testTimestamp = 1750861210
	testWebhookID = "wh_685c059ae39c8190af8c71ed1022a24d"
	testSignature = "v1,gUAg4R2hWouRZqRQG4uJypNS8YK885G838+EHb4nKBY="
)

// Helper function to create test headers with standardized timestamp
func createTestHeaders() http.Header {
	return http.Header{
		"Webhook-Signature": []string{testSignature},
		"Webhook-Timestamp": []string{strconv.FormatInt(testTimestamp, 10)},
		"Webhook-Id":        []string{testWebhookID},
	}
}

func TestWebhookService_VerifySignature_ValidSignature(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithWebhookSecret(testSecret),
	)

	// Use the time-aware method with the fixed timestamp for testing
	fixedTime := time.Unix(testTimestamp, 0)
	err := client.Webhooks.VerifySignatureWithToleranceAndTime([]byte(testPayload), createTestHeaders(), 5*time.Minute, fixedTime)
	if err != nil {
		t.Errorf("VerifySignatureWithToleranceAndTime should have succeeded with valid signature: %v", err)
	}
}

func TestWebhookService_VerifySignature_InvalidSignature(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithWebhookSecret(testSecret),
	)

	// Create headers with an invalid signature
	invalidHeaders := http.Header{
		"Webhook-Signature": []string{"v1,invalid_signature_here"},
		"Webhook-Timestamp": []string{strconv.FormatInt(testTimestamp, 10)},
		"Webhook-Id":        []string{testWebhookID},
	}

	// Use time-aware method to avoid timestamp issues
	fixedTime := time.Unix(testTimestamp, 0)
	err := client.Webhooks.VerifySignatureWithToleranceAndTime([]byte(testPayload), invalidHeaders, 5*time.Minute, fixedTime)
	if err == nil {
		t.Error("VerifySignatureWithToleranceAndTime should have failed with invalid signature")
	}

	expectedError := "webhook signature verification failed"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestWebhookService_VerifySignature_MissingSecret(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
	)

	err := client.Webhooks.VerifySignature([]byte(testPayload), createTestHeaders())
	if err == nil {
		t.Error("VerifySignature should have failed with missing secret")
	}

	expectedError := "webhook secret must be provided either in the method call or configured on the client"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestWebhookService_VerifySignature_InvalidSignatureForInvalidTimestamp(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithWebhookSecret(testSecret),
	)

	invalidHeaders := http.Header{
		"Webhook-Signature": []string{"v1,signature"},
		"Webhook-Timestamp": []string{"invalid"},
		"Webhook-Id":        []string{testWebhookID},
	}

	err := client.Webhooks.VerifySignature([]byte(testPayload), invalidHeaders)
	if err == nil {
		t.Error("VerifySignature should have failed with invalid timestamp")
	}

	// Should fail on timestamp format validation first
	expectedError := "invalid webhook timestamp format"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestWebhookService_VerifySignature_SignatureWithoutV1Prefix(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithWebhookSecret(testSecret),
	)

	currentTimestamp := strconv.FormatInt(testTimestamp, 10)
	headersWithoutV1 := http.Header{
		"Webhook-Signature": []string{"9WlByKQUfBVM08XRYmo3WqR/dQXtjGJkV1edShZZ+C0="},
		"Webhook-Timestamp": []string{currentTimestamp},
		"Webhook-Id":        []string{testWebhookID},
	}

	// Use time-aware method to avoid timestamp issues
	fixedTime := time.Unix(testTimestamp, 0)
	err := client.Webhooks.VerifySignatureWithToleranceAndTime([]byte(testPayload), headersWithoutV1, 5*time.Minute, fixedTime)
	if err == nil {
		t.Error("VerifySignatureWithToleranceAndTime should have failed on signature verification")
	}

	// Should fail on signature verification
	expectedError := "webhook signature verification failed"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestWebhookService_VerifySignature_TimestampTooOld(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithWebhookSecret(testSecret),
	)

	// Use a timestamp that's older than 5 minutes
	oldTimestamp := strconv.FormatInt(testTimestamp-400, 10) // 6 minutes 40 seconds ago
	headersOld := http.Header{
		"Webhook-Signature": []string{"v1,signature"},
		"Webhook-Timestamp": []string{oldTimestamp},
		"Webhook-Id":        []string{testWebhookID},
	}

	// Use the current time for comparison - this should fail because the timestamp is too old
	currentTime := time.Unix(testTimestamp, 0)
	err := client.Webhooks.VerifySignatureWithToleranceAndTime([]byte(testPayload), headersOld, 5*time.Minute, currentTime)
	if err == nil {
		t.Error("VerifySignatureWithToleranceAndTime should have failed with old timestamp")
	}

	expectedError := "webhook timestamp is too old"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestWebhookService_VerifySignature_TimestampTooNew(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithWebhookSecret(testSecret),
	)

	// Use a timestamp that's in the future beyond tolerance
	futureTimestamp := strconv.FormatInt(testTimestamp+400, 10) // 6 minutes 40 seconds in the future
	headersFuture := http.Header{
		"Webhook-Signature": []string{"v1,signature"},
		"Webhook-Timestamp": []string{futureTimestamp},
		"Webhook-Id":        []string{testWebhookID},
	}

	// Use the current time for comparison - this should fail because the timestamp is too new
	currentTime := time.Unix(testTimestamp, 0)
	err := client.Webhooks.VerifySignatureWithToleranceAndTime([]byte(testPayload), headersFuture, 5*time.Minute, currentTime)
	if err == nil {
		t.Error("VerifySignatureWithToleranceAndTime should have failed with future timestamp")
	}

	expectedError := "webhook timestamp is too new"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestWebhookService_VerifySignature_InvalidTimestampFormat(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithWebhookSecret(testSecret),
	)

	headersInvalid := http.Header{
		"Webhook-Signature": []string{"v1,signature"},
		"Webhook-Timestamp": []string{"not_a_number"},
		"Webhook-Id":        []string{testWebhookID},
	}

	err := client.Webhooks.VerifySignature([]byte(testPayload), headersInvalid)
	if err == nil {
		t.Error("VerifySignature should have failed with invalid timestamp format")
	}

	expectedError := "invalid webhook timestamp format"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestWebhookService_VerifySignatureWithTolerance_CustomTolerance(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithWebhookSecret(testSecret),
	)

	// Use a timestamp that's older than default tolerance but within custom tolerance
	oldTimestamp := strconv.FormatInt(testTimestamp-400, 10) // 6 minutes 40 seconds ago
	headersOld := http.Header{
		"Webhook-Signature": []string{"v1,signature"},
		"Webhook-Timestamp": []string{oldTimestamp},
		"Webhook-Id":        []string{testWebhookID},
	}

	// Should fail with default tolerance using time-aware method
	currentTime := time.Unix(testTimestamp, 0)
	err := client.Webhooks.VerifySignatureWithToleranceAndTime([]byte(testPayload), headersOld, 5*time.Minute, currentTime)
	if err == nil {
		t.Error("VerifySignatureWithToleranceAndTime should have failed with default tolerance")
	}
	if err.Error() != "webhook timestamp is too old" {
		t.Errorf("Expected 'webhook timestamp is too old', got '%s'", err.Error())
	}

	// Should pass timestamp validation with custom tolerance but fail on signature verification
	err = client.Webhooks.VerifySignatureWithToleranceAndTime([]byte(testPayload), headersOld, 10*time.Minute, currentTime)
	if err == nil {
		t.Error("VerifySignatureWithToleranceAndTime should have failed on signature verification")
	}
	if err.Error() != "webhook signature verification failed" {
		t.Errorf("Expected 'webhook signature verification failed', got '%s'", err.Error())
	}
}

func TestWebhookService_UnwrapWithToleranceAndTime(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithWebhookSecret(testSecret),
	)

	// Use the time-aware method with the fixed timestamp for testing
	fixedTime := time.Unix(testTimestamp, 0)
	webhookEvent, err := client.Webhooks.UnwrapWithToleranceAndTime([]byte(testPayload), createTestHeaders(), 5*time.Minute, fixedTime)
	if err != nil {
		t.Errorf("UnwrapWithToleranceAndTime should have succeeded with valid signature: %v", err)
	}

	if webhookEvent == nil {
		t.Error("UnwrapWithToleranceAndTime should return parsed event")
	}

	parsed := webhookEvent.AsResponseCompleted()
	if parsed.ID != "evt_685c059ae3a481909bdc86819b066fb6" {
		t.Errorf("Expected parsed event ID 'evt_685c059ae3a481909bdc86819b066fb6', got '%s'", parsed.ID)
	}
	if parsed.Type != "response.completed" {
		t.Errorf("Expected parsed event type 'response.completed', got '%s'", parsed.Type)
	}
}

func TestWebhookService_VerifySignature_MultipleSignaturesOneValid(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithWebhookSecret(testSecret),
	)

	// Test multiple signatures: one invalid, one valid
	multipleSignatures := "v1,invalid_signature " + testSignature
	multipleHeaders := http.Header{
		"Webhook-Signature": []string{multipleSignatures},
		"Webhook-Timestamp": []string{strconv.FormatInt(testTimestamp, 10)},
		"Webhook-Id":        []string{testWebhookID},
	}

	// Use time-aware method with fixed timestamp
	fixedTime := time.Unix(testTimestamp, 0)
	err := client.Webhooks.VerifySignatureWithToleranceAndTime([]byte(testPayload), multipleHeaders, 5*time.Minute, fixedTime)
	if err != nil {
		t.Errorf("VerifySignatureWithToleranceAndTime should have succeeded when at least one signature is valid: %v", err)
	}
}

func TestWebhookService_VerifySignature_MultipleSignaturesAllInvalid(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithWebhookSecret(testSecret),
	)

	// Test multiple invalid signatures
	multipleInvalidSignatures := "v1,invalid_signature1 v1,invalid_signature2"
	multipleHeaders := http.Header{
		"Webhook-Signature": []string{multipleInvalidSignatures},
		"Webhook-Timestamp": []string{strconv.FormatInt(testTimestamp, 10)},
		"Webhook-Id":        []string{testWebhookID},
	}

	// Use time-aware method with fixed timestamp
	fixedTime := time.Unix(testTimestamp, 0)
	err := client.Webhooks.VerifySignatureWithToleranceAndTime([]byte(testPayload), multipleHeaders, 5*time.Minute, fixedTime)
	if err == nil {
		t.Error("VerifySignatureWithToleranceAndTime should have failed when all signatures are invalid")
	}

	expectedError := "webhook signature verification failed"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}
