package webhooks_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/webhooks"
)

func TestSignedAgentSessionWebhooks(t *testing.T) {
	client := openai.NewClient(option.WithWebhookSecret(webhookHelperSecret))
	timestamp := time.Now().Unix()
	for _, tc := range []struct {
		event string
		data  string
	}{
		{"agent.session.action_required", `{"id":"as_1","required_action":{"type":"function_call"}}`},
		{"agent.session.created", `{"id":"as_1","environment_type":"sandbox","environment_id":"env_1","connect":{"remote_url":"https://synthetic.example/agent"}}`},
		{"agent.session.failed", `{"id":"as_1","environment_type":"sandbox","environment_id":"env_1"}`},
		{"agent.session.idle", `{"id":"as_1","environment_type":"sandbox","environment_id":"env_1"}`},
		{"agent.session.in_progress", `{"id":"as_1","environment_type":"sandbox","environment_id":"env_1"}`},
	} {
		t.Run(tc.event, func(t *testing.T) {
			payload := fmt.Sprintf(`{"id":"evt_1","object":"event","created_at":%d,"type":%q,"data":%s}`, timestamp, tc.event, tc.data)
			event, err := client.Webhooks.Unwrap([]byte(payload), signedWebhookHelperHeaders(t, []byte(payload), timestamp))
			if err != nil {
				t.Fatal(err)
			}
			if event.ID != "evt_1" || event.Data.ID != "as_1" || event.Type != tc.event {
				t.Fatalf("wrong event envelope: %#v", event)
			}
			switch v := event.AsAny().(type) {
			case webhooks.AgentSessionActionRequiredWebhookEvent:
				if tc.event != "agent.session.action_required" || v.Data.RequiredAction.Type != "function_call" {
					t.Fatalf("action: %#v", v)
				}
			case webhooks.AgentSessionCreatedWebhookEvent:
				if tc.event != "agent.session.created" || v.Data.EnvironmentType != "sandbox" || v.Data.Connect.RemoteURL != "https://synthetic.example/agent" {
					t.Fatalf("created: %#v", v)
				}
			case webhooks.AgentSessionFailedWebhookEvent:
				if tc.event != "agent.session.failed" || v.Data.EnvironmentID != "env_1" {
					t.Fatalf("failed: %#v", v)
				}
			case webhooks.AgentSessionIdleWebhookEvent:
				if tc.event != "agent.session.idle" || v.Data.EnvironmentID != "env_1" {
					t.Fatalf("idle: %#v", v)
				}
			case webhooks.AgentSessionInProgressWebhookEvent:
				if tc.event != "agent.session.in_progress" || v.Data.EnvironmentID != "env_1" {
					t.Fatalf("progress: %#v", v)
				}
			default:
				t.Fatalf("unrecognized typed event %T", v)
			}
			if _, err := client.Webhooks.Unwrap([]byte(payload), signedWebhookHelperHeaders(t, []byte(payload+"tamper"), timestamp)); err == nil {
				t.Fatal("accepted a signature over different data")
			}
			t.Run("invalid data shape", func(t *testing.T) {
				body := []byte(fmt.Sprintf(`{"id":"evt_1","object":"event","created_at":%d,"type":%q,"data":[]}`, timestamp, tc.event))
				invalidEvent, invalidErr := client.Webhooks.Unwrap(body, signedWebhookHelperHeaders(t, body, timestamp))
				if invalidErr != nil {
					t.Fatal(invalidErr)
				}
				// The response decoder preserves invalid field data in metadata for
				// callers to inspect rather than rejecting the complete event.
				if invalidEvent.Type != tc.event || invalidEvent.AsAny() == nil || invalidEvent.JSON.Data.Valid() || invalidEvent.JSON.Data.Raw() != "[]" {
					t.Fatalf("invalid data was not preserved for %s: %#v", tc.event, invalidEvent)
				}
			})
		})
	}
}

func TestSignedAgentEnvironmentWebhooks(t *testing.T) {
	client := openai.NewClient(option.WithWebhookSecret(webhookHelperSecret))
	timestamp := time.Now().Unix()
	for _, tc := range []struct {
		event    string
		accessor func(webhooks.UnwrapWebhookEventUnion) (any, string)
	}{
		{"agent.environment.ready", func(u webhooks.UnwrapWebhookEventUnion) (any, string) {
			v := u.AsAgentEnvironmentReady()
			return v, v.Data.ID
		}},
		{"agent.environment.failed", func(u webhooks.UnwrapWebhookEventUnion) (any, string) {
			v := u.AsAgentEnvironmentFailed()
			return v, v.Data.ID
		}},
		{"agent.environment.suspended", func(u webhooks.UnwrapWebhookEventUnion) (any, string) {
			v := u.AsAgentEnvironmentSuspended()
			return v, v.Data.ID
		}},
		{"agent.environment.expired", func(u webhooks.UnwrapWebhookEventUnion) (any, string) {
			v := u.AsAgentEnvironmentExpired()
			return v, v.Data.ID
		}},
	} {
		t.Run(tc.event, func(t *testing.T) {
			payload := []byte(fmt.Sprintf(`{"id":"evt_1","object":"event","created_at":%d,"type":%q,"data":{"id":"env_1"}}`, timestamp, tc.event))
			event, err := client.Webhooks.Unwrap(payload, signedWebhookHelperHeaders(t, payload, timestamp))
			if err != nil {
				t.Fatal(err)
			}
			if event.ID != "evt_1" || event.Object != "event" || event.CreatedAt != timestamp || event.Type != tc.event || event.Data.ID != "env_1" {
				t.Fatalf("wrong environment event envelope: %#v", event)
			}
			typed, environmentID := tc.accessor(*event)
			if environmentID != "env_1" {
				t.Fatalf("typed environment ID = %q, want env_1", environmentID)
			}
			if got := event.AsAny(); !reflect.DeepEqual(got, typed) {
				t.Fatalf("AsAny() = %#v (%T), want %#v (%T)", got, got, typed, typed)
			}
		})
	}
}
