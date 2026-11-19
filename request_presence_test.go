package openai_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
)

func requestPresenceBody(t *testing.T, call func(openai.Client) error) map[string]any {
	t.Helper()
	req := captureRequest(t, call, "{}")
	var body map[string]any
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	return body
}

func TestRequestPresenceBoolean(t *testing.T) {
	for _, tt := range []struct {
		name    string
		enabled param.Opt[bool]
		want    map[string]any
	}{
		{"omitted", param.Opt[bool]{}, map[string]any{"max_concurrent_subagents": 2.0}},
		{"false", openai.Bool(false), map[string]any{"max_concurrent_subagents": 2.0, "enabled": false}},
		{"true", openai.Bool(true), map[string]any{"max_concurrent_subagents": 2.0, "enabled": true}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			body := requestPresenceBody(t, func(client openai.Client) error {
				_, err := client.Beta.Responses.New(t.Context(), openai.BetaResponseNewParams{
					MultiAgent: openai.BetaResponseNewParamsMultiAgent{MaxConcurrentSubagents: openai.Int(2), Enabled: tt.enabled},
				})
				return err
			})
			if got := body["multi_agent"]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("multi_agent = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestRequestPresenceNullableParent(t *testing.T) {
	body := requestPresenceBody(t, func(client openai.Client) error {
		_, err := client.Admin.Organization.Projects.HostedToolPermissions.Update(t.Context(), "proj_test", openai.AdminOrganizationProjectHostedToolPermissionUpdateParams{
			CodeInterpreter: param.NullStruct[openai.AdminOrganizationProjectHostedToolPermissionUpdateParamsCodeInterpreter](),
		})
		return err
	})
	if value, present := body["code_interpreter"]; !present || value != nil {
		t.Errorf("expected code_interpreter to be present and null, got %#v", body)
	}
}

func TestRequestPresenceZeroObjects(t *testing.T) {
	t.Run("session workflow", func(t *testing.T) {
		body := requestPresenceBody(t, func(client openai.Client) error {
			_, err := client.Beta.ChatKit.Sessions.New(t.Context(), openai.BetaChatKitSessionNewParams{User: "test"})
			return err
		})
		if want := map[string]any{"id": ""}; !reflect.DeepEqual(body["workflow"], want) {
			t.Errorf("workflow = %#v, want %#v", body["workflow"], want)
		}
		if _, present := body["chatkit_configuration"]; present {
			t.Errorf("expected the optional chatkit_configuration object to be omitted, got %#v", body)
		}
	})

	t.Run("spend alert notification channel", func(t *testing.T) {
		body := requestPresenceBody(t, func(client openai.Client) error {
			_, err := client.Admin.Organization.Projects.SpendAlerts.New(t.Context(), "proj_test", openai.AdminOrganizationProjectSpendAlertNewParams{})
			return err
		})
		if want := map[string]any{"type": "email"}; !reflect.DeepEqual(body["notification_channel"], want) {
			t.Errorf("notification_channel = %#v, want %#v", body["notification_channel"], want)
		}
	})
}
