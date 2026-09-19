package webhooks_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/webhooks"
)

func TestWebhookService_UnwrapSafetyEvents(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("synthetic-api-key"),
		option.WithWebhookSecret(sipMediaSecurityTestSecret),
	)
	cases := []struct {
		eventType string
		wantType  reflect.Type
		caseID    func(*webhooks.UnwrapWebhookEventUnion) (string, bool)
	}{
		{
			eventType: "safety.deactivation_issued",
			wantType:  reflect.TypeOf(webhooks.SafetyDeactivationIssuedWebhookEvent{}),
			caseID: func(event *webhooks.UnwrapWebhookEventUnion) (string, bool) {
				data := event.AsSafetyDeactivationIssued().Data
				return data.ID, data.JSON.ID.Valid()
			},
		},
		{
			eventType: "safety.warning_issued",
			wantType:  reflect.TypeOf(webhooks.SafetyWarningIssuedWebhookEvent{}),
			caseID: func(event *webhooks.UnwrapWebhookEventUnion) (string, bool) {
				data := event.AsSafetyWarningIssued().Data
				return data.ID, data.JSON.ID.Valid()
			},
		},
	}
	for _, test := range cases {
		t.Run(test.eventType, func(t *testing.T) {
			body := []byte(fmt.Sprintf(
				`{"id":"evt_safety","object":"event","created_at":1750861210,"type":%q,"data":{"id":"case_synthetic"}}`,
				test.eventType,
			))
			event, err := client.Webhooks.Unwrap(body, signedSIPMediaSecurityHeaders(t, body))
			if err != nil {
				t.Fatalf("Unwrap(%s) error = %v, want nil", test.eventType, err)
			}
			if event.ID != "evt_safety" || event.Object != "event" || event.CreatedAt != 1750861210 || event.Type != test.eventType {
				t.Errorf("Unwrap(%s) envelope = (%q, %q, %d, %q), want (evt_safety, event, 1750861210, %s)",
					test.eventType, event.ID, event.Object, event.CreatedAt, event.Type, test.eventType)
			}
			if got := reflect.TypeOf(event.AsAny()); got != test.wantType {
				t.Errorf("Unwrap(%s).AsAny() type = %v, want %v", test.eventType, got, test.wantType)
			}
			if got, present := test.caseID(event); got != "case_synthetic" || !present {
				t.Errorf("Unwrap(%s) typed case ID = (%q, present=%t), want (case_synthetic, present=true)", test.eventType, got, present)
			}
			if event.Data.ID != "case_synthetic" || !event.Data.JSON.ID.Valid() {
				t.Errorf("Unwrap(%s) union case ID = (%q, present=%t), want (case_synthetic, present=true)",
					test.eventType, event.Data.ID, event.Data.JSON.ID.Valid())
			}
		})
	}
}

func TestWebhookService_UnwrapSafetyEventsRejectsInvalidSignatures(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("synthetic-api-key"),
		option.WithWebhookSecret(sipMediaSecurityTestSecret),
	)
	for _, eventType := range []string{"safety.deactivation_issued", "safety.warning_issued"} {
		payload := fmt.Sprintf(
			`{"id":"evt_safety","object":"event","created_at":1750861210,"type":%q,"data":{"id":"case_synthetic"}}`,
			eventType,
		)
		cases := []struct {
			name       string
			body       string
			signedBody string
		}{
			{name: "tampered_case", body: strings.Replace(payload, "case_synthetic", "case_other", 1), signedBody: payload},
			{name: "invalid_signature", body: payload, signedBody: "different synthetic payload"},
		}
		for _, test := range cases {
			t.Run(eventType+"/"+test.name, func(t *testing.T) {
				event, err := client.Webhooks.Unwrap(
					[]byte(test.body), signedSIPMediaSecurityHeaders(t, []byte(test.signedBody)),
				)
				if err == nil || event != nil {
					t.Errorf("Unwrap(%s, %s) = (%v, %v), want (nil, error)", eventType, test.name, event, err)
				}
			})
		}
	}
}
