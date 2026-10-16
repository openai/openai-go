package live_test

import (
	"testing"

	"github.com/openai/openai-go/v3/live"
)

func TestNestedDelegationGetters(t *testing.T) {
	client := &live.ClientDelegationParam{Type: "client"}
	start := &live.SessionConfigDelegationResponsesParam{
		Type: "responses", Responses: live.ResponsesDelegationConfigParam{Model: "synthetic-model"},
	}
	update := &live.SessionUpdateConfigDelegationResponsesParam{Type: "responses"}
	check := func(t *testing.T, got *string, want string, responses, wantResponses any) {
		t.Helper()
		if got == nil || *got != want {
			t.Errorf("delegation type = %v, want %q", got, want)
		}
		if responses != wantResponses {
			t.Errorf("responses = %v, want original value %v", responses, wantResponses)
		}
	}
	t.Run("start client", func(t *testing.T) {
		event := live.ClientEventUnionParam{OfSessionStart: &live.SessionStartEventParam{
			Session: live.SessionConfigParam{Delegation: live.SessionConfigDelegationUnionParam{OfClient: client}},
		}}
		got := event.GetSession().GetDelegation()
		check(t, got.GetType(), "client", got.GetResponses().AsAny(), nil)
	})
	t.Run("start responses", func(t *testing.T) {
		event := live.ClientEventUnionParam{OfSessionStart: &live.SessionStartEventParam{
			Session: live.SessionConfigParam{Delegation: live.SessionConfigDelegationUnionParam{OfResponses: start}},
		}}
		got := event.GetSession().GetDelegation()
		check(t, got.GetType(), "responses", got.GetResponses().AsAny(), &start.Responses)
	})
	for _, owner := range []string{"client", "responses"} {
		t.Run("update "+owner, func(t *testing.T) {
			delegation := live.SessionUpdateConfigDelegationUnionParam{OfClient: client}
			var want any
			if owner == "responses" {
				delegation = live.SessionUpdateConfigDelegationUnionParam{OfResponses: update}
				want = &update.Responses
			}
			event := &live.SessionUpdateEventParam{Session: live.SessionUpdateConfigParam{Delegation: delegation}}
			got := (live.ClientEventUnionParam{OfSessionUpdate: event}).GetSession().GetDelegation()
			check(t, got.GetType(), owner, got.GetResponses().AsAny(), want)
			fork := (live.ForkClientEventUnionParam{OfSessionUpdate: event}).GetSession().GetDelegation()
			var forkWant *live.ResponsesDelegationUpdateConfigParam
			if owner == "responses" {
				forkWant = &update.Responses
			}
			check(t, fork.GetType(), owner, fork.GetResponses(), forkWant)
		})
	}
	t.Run("absent", func(t *testing.T) {
		got := (live.ClientEventUnionParam{}).GetSession().GetDelegation()
		if got.GetType() != nil || got.GetResponses().AsAny() != nil {
			t.Fatal("absent delegation must remain absent")
		}
	})
}
