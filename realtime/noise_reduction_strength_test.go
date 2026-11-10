package realtime_test

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/realtime"
	"github.com/tidwall/gjson"
)

func TestNoiseReductionStrengthGetters(t *testing.T) {
	for _, variant := range []string{"realtime", "transcription"} {
		var update realtime.SessionUpdateEventParam
		var secret realtime.ClientSecretNewParams
		var strength *param.Opt[float64]
		if variant == "realtime" {
			session := &realtime.RealtimeSessionCreateRequestParam{}
			session.Audio.Input.NoiseReduction.Type = realtime.NoiseReductionTypeNearField
			strength = &session.Audio.Input.NoiseReduction.Strength
			update.Session.OfRealtime, secret.Session.OfRealtime = session, session
		} else {
			session := &realtime.RealtimeTranscriptionSessionCreateRequestParam{}
			session.Audio.Input.NoiseReduction.Type = realtime.NoiseReductionTypeNearField
			strength = &session.Audio.Input.NoiseReduction.Strength
			update.Session.OfTranscription, secret.Session.OfTranscription = session, session
		}

		for _, tc := range []struct {
			name    string
			request json.Marshaler
			get     func() *float64
		}{
			{"session update", update, update.Session.GetAudio().GetInput().GetNoiseReduction().GetStrength},
			{"client secret", secret, secret.Session.GetAudio().GetInput().GetNoiseReduction().GetStrength},
		} {
			t.Run(variant+"/"+tc.name, func(t *testing.T) {
				*strength = param.Opt[float64]{}
				if tc.get() != nil {
					t.Fatal("omitted strength returned a pointer")
				}
				*strength = param.Null[float64]()
				if tc.get() != nil {
					t.Fatal("null strength returned a pointer")
				}
				*strength = param.NewOpt(0.5)
				got := tc.get()
				if got == nil || *got != 0.5 {
					t.Fatalf("GetStrength() = %v, want a pointer to 0.5", got)
				}
				*got = 0
				body, err := json.Marshal(tc.request)
				if err != nil {
					t.Fatal(err)
				}
				if got := gjson.GetBytes(body, "session.audio.input.noise_reduction.strength"); got.Type != gjson.Number || got.Float() != 0 {
					t.Fatalf("marshaled strength = %q after mutation, want 0", got.Raw)
				}
			})
		}
	}
}
