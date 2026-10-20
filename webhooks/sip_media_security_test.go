package webhooks_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/webhooks"
)

const sipMediaSecurityTestSecret = "synthetic-sip-media-security-secret"

func TestWebhookService_UnwrapSIPMediaSecurity(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("synthetic-api-key"),
		option.WithWebhookSecret(sipMediaSecurityTestSecret),
	)
	events := []struct {
		eventType string
		data      string
		field     func(*webhooks.UnwrapWebhookEventUnion) (string, respjson.Field)
	}{
		{
			eventType: "live.call.incoming",
			data:      `"session_id":"sess_sip","sip_headers":[]`,
			field: func(event *webhooks.UnwrapWebhookEventUnion) (string, respjson.Field) {
				data := event.AsLiveCallIncoming().Data
				return data.SipMediaSecurity, data.JSON.SipMediaSecurity
			},
		},
		{
			eventType: "live.transport.incoming",
			data:      `"session_id":"sess_sip","sip_headers":[],"type":"sip"`,
			field: func(event *webhooks.UnwrapWebhookEventUnion) (string, respjson.Field) {
				data := event.AsLiveTransportIncoming().Data
				return data.SipMediaSecurity, data.JSON.SipMediaSecurity
			},
		},
		{
			eventType: "realtime.call.incoming",
			data:      `"call_id":"call_sip","sip_headers":[]`,
			field: func(event *webhooks.UnwrapWebhookEventUnion) (string, respjson.Field) {
				data := event.AsRealtimeCallIncoming().Data
				return data.SipMediaSecurity, data.JSON.SipMediaSecurity
			},
		},
	}
	values := []struct {
		name      string
		fieldJSON string
		want      string
		wantRaw   string
		wantValid bool
	}{
		{name: "omitted"},
		{name: "empty", fieldJSON: `,"sip_media_security":""`, wantRaw: `""`, wantValid: true},
		{name: "rtp", fieldJSON: `,"sip_media_security":"rtp"`, want: "rtp", wantRaw: `"rtp"`, wantValid: true},
		{name: "srtp", fieldJSON: `,"sip_media_security":"srtp"`, want: "srtp", wantRaw: `"srtp"`, wantValid: true},
		{name: "future", fieldJSON: `,"sip_media_security":"future-protection"`, want: "future-protection", wantRaw: `"future-protection"`, wantValid: true},
		{name: "null", fieldJSON: `,"sip_media_security":null`, wantRaw: "null"},
	}
	for _, eventCase := range events {
		for _, valueCase := range values {
			t.Run(eventCase.eventType+"/"+valueCase.name, func(t *testing.T) {
				body := []byte(fmt.Sprintf(
					`{"id":"evt_sip","object":"event","created_at":1750861210,"type":%q,"data":{%s%s}}`,
					eventCase.eventType, eventCase.data, valueCase.fieldJSON,
				))
				event, err := client.Webhooks.Unwrap(body, signedSIPMediaSecurityHeaders(t, body))
				if err != nil {
					t.Fatalf("Unwrap(%s, %s) error = %v, want nil", eventCase.eventType, valueCase.name, err)
				}
				if event.Type != eventCase.eventType {
					t.Errorf("Unwrap(%s).Type = %q, want %q", valueCase.name, event.Type, eventCase.eventType)
				}
				got, field := eventCase.field(event)
				if got != valueCase.want || field.Valid() != valueCase.wantValid || field.Raw() != valueCase.wantRaw {
					t.Errorf("Unwrap(%s, %s) variant field = (%q, valid=%t, raw=%q), want (%q, valid=%t, raw=%q)",
						eventCase.eventType, valueCase.name, got, field.Valid(), field.Raw(),
						valueCase.want, valueCase.wantValid, valueCase.wantRaw)
				}
				union := event.Data
				if union.SipMediaSecurity != valueCase.want ||
					union.JSON.SipMediaSecurity.Valid() != valueCase.wantValid ||
					union.JSON.SipMediaSecurity.Raw() != valueCase.wantRaw {
					t.Errorf("Unwrap(%s, %s) union field = (%q, valid=%t, raw=%q), want (%q, valid=%t, raw=%q)",
						eventCase.eventType, valueCase.name, union.SipMediaSecurity,
						union.JSON.SipMediaSecurity.Valid(), union.JSON.SipMediaSecurity.Raw(),
						valueCase.want, valueCase.wantValid, valueCase.wantRaw)
				}
			})
		}
	}
}

func TestWebhookService_UnwrapSIPMediaSecurityRejectsInvalidSignatures(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("synthetic-api-key"),
		option.WithWebhookSecret(sipMediaSecurityTestSecret),
	)
	const payload = `{"id":"evt_sip","object":"event","created_at":1750861210,"type":"live.call.incoming","data":{"session_id":"sess_sip","sip_headers":[],"sip_media_security":"rtp"}}`
	cases := []struct {
		name       string
		body       string
		signedBody string
	}{
		{name: "tampered_media_security", body: strings.Replace(payload, `"rtp"`, `"srtp"`, 1), signedBody: payload},
		{name: "invalid_signature", body: payload, signedBody: "different synthetic body"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			event, err := client.Webhooks.Unwrap(
				[]byte(test.body), signedSIPMediaSecurityHeaders(t, []byte(test.signedBody)),
			)
			if err == nil || event != nil {
				t.Errorf("Unwrap(%s) = (%v, %v), want (nil, error)", test.name, event, err)
			}
		})
	}
}

func signedSIPMediaSecurityHeaders(t *testing.T, body []byte) http.Header {
	t.Helper()
	const webhookID = "wh_sip_media_security"
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	mac := hmac.New(sha256.New, []byte(sipMediaSecurityTestSecret))
	if _, err := mac.Write(append([]byte(webhookID+"."+timestamp+"."), body...)); err != nil {
		t.Fatalf("HMAC.Write(synthetic SIP payload) error = %v, want nil", err)
	}
	return http.Header{
		"Webhook-Id":        []string{webhookID},
		"Webhook-Timestamp": []string{timestamp},
		"Webhook-Signature": []string{"v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))},
	}
}
