package webhooks_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

const (
	readmeWebhookBodyLimit int64 = 1 << 20
	readmeWebhookSecret          = "readme-test-secret"
	readmeWebhookID              = "wh_readme_test"
)

func readmeWebhookHandler(client openai.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, readmeWebhookBodyLimit)
		defer func() { _ = r.Body.Close() }()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(err, &maxBytesError) {
				http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, "error reading request body", http.StatusBadRequest)
			return
		}

		if err := client.Webhooks.VerifySignature(body, r.Header); err != nil {
			http.Error(w, "invalid signature", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func signReadmeWebhook(body []byte, timestamp string) string {
	signedPayload := readmeWebhookID + "." + timestamp + "." + string(body)
	mac := hmac.New(sha256.New, []byte(readmeWebhookSecret))
	_, _ = mac.Write([]byte(signedPayload))
	return "v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func TestREADMEWebhookHandlerRejectsOversizedBodyBeforeVerification(t *testing.T) {
	client := openai.NewClient(option.WithWebhookSecret(readmeWebhookSecret))
	body := bytes.Repeat([]byte("x"), int(readmeWebhookBodyLimit)+1)
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	recorder := httptest.NewRecorder()

	readmeWebhookHandler(client).ServeHTTP(recorder, req)

	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d, got %d", http.StatusRequestEntityTooLarge, recorder.Code)
	}
}

func TestREADMEWebhookHandlerAcceptsSignedBodyAtLimit(t *testing.T) {
	client := openai.NewClient(option.WithWebhookSecret(readmeWebhookSecret))
	body := bytes.Repeat([]byte("x"), int(readmeWebhookBodyLimit))
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("webhook-id", readmeWebhookID)
	req.Header.Set("webhook-timestamp", timestamp)
	req.Header.Set("webhook-signature", signReadmeWebhook(body, timestamp))
	recorder := httptest.NewRecorder()

	readmeWebhookHandler(client).ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
}

func TestREADMEWebhookExamplesEnforceRequestLimits(t *testing.T) {
	contents, err := os.ReadFile("../README.md")
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	_, webhookSection, found := strings.Cut(string(contents), "## Webhook Verification")
	if !found {
		t.Fatal("README does not contain the webhook verification section")
	}
	webhookSection, _, _ = strings.Cut(webhookSection, "### Retries")

	checks := map[string]int{
		"const maxWebhookBodySize = 1 << 20": 2,
		"http.MaxBytesReader":                2,
		"http.StatusRequestEntityTooLarge":   2,
		"ReadHeaderTimeout:":                 2,
		"ReadTimeout:":                       2,
	}
	for snippet, want := range checks {
		if got := strings.Count(webhookSection, snippet); got != want {
			t.Errorf("README webhook examples contain %q %d times, want %d", snippet, got, want)
		}
	}
}
