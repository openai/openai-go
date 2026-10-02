package openai_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
)

type voiceMetadataTransport func(*http.Request) (*http.Response, error)

func (f voiceMetadataTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestAudioVoicesNewRootMetadata(t *testing.T) {
	for _, test := range []struct {
		name     string
		request  func() openai.AudioVoiceNewParams
		json     string
		want     map[string][]string
		wantFile bool
	}{
		{
			name: "prompt root overrides child and retains extra fields",
			request: func() openai.AudioVoiceNewParams {
				child := &openai.AudioVoiceNewParamsBodyPrompt{Type: "prompt", Name: "old name", Prompt: "calm"}
				child.SetExtraFields(map[string]any{"child_field": "keep", "name": "child name"})
				r := openai.AudioVoiceNewParams{OfPrompt: child}
				r.SetExtraFields(map[string]any{"name": "root name", "future_count": 2, "prompt": param.Omit})
				return r
			},
			want: map[string][]string{"name": {"root name"}, "type": {"prompt"}, "child_field": {"keep"}, "future_count": {"2"}},
		},
		{
			name: "audio sample keeps file and root extras",
			request: func() openai.AudioVoiceNewParams {
				r := openai.AudioVoiceNewParams{OfAudioSample: &openai.AudioVoiceNewParamsBodyAudioSample{
					AudioSample: bytes.NewBufferString("fixture audio"), Consent: "consent_1", Name: "old name",
				}}
				r.SetExtraFields(map[string]any{"name": "root name", "future_field": "value"})
				return r
			},
			want:     map[string][]string{"name": {"root name"}, "consent": {"consent_1"}, "future_field": {"value"}},
			wantFile: true,
		},
		{
			name: "prompt child metadata overrides without root extras",
			request: func() openai.AudioVoiceNewParams {
				child := &openai.AudioVoiceNewParamsBodyPrompt{Type: "prompt", Name: "old name", Prompt: "calm"}
				child.SetExtraFields(map[string]any{"child_field": "keep", "name": "child name"})
				return openai.AudioVoiceNewParams{OfPrompt: child}
			},
			want: map[string][]string{"name": {"child name"}, "type": {"prompt"}, "prompt": {"calm"}, "child_field": {"keep"}},
		},
		{
			name: "audio sample child metadata overrides without root extras",
			request: func() openai.AudioVoiceNewParams {
				child := &openai.AudioVoiceNewParamsBodyAudioSample{
					AudioSample: bytes.NewBufferString("fixture audio"), Consent: "consent_1", Name: "old name",
				}
				child.SetExtraFields(map[string]any{"name": "child name", "consent": param.Omit, "future_field": "value"})
				return openai.AudioVoiceNewParams{OfAudioSample: child}
			},
			want:     map[string][]string{"name": {"child name"}, "future_field": {"value"}},
			wantFile: true,
		},
		{
			name: "prompt child JSON replaces reflected fields",
			request: func() openai.AudioVoiceNewParams {
				child := &openai.AudioVoiceNewParamsBodyPrompt{Type: "prompt", Name: "ignored", Prompt: "ignored"}
				param.SetJSON([]byte(`{"name":"json override","type":"prompt","prompt":"soft"}`), child)
				return openai.AudioVoiceNewParams{OfPrompt: child}
			},
			json: `{"name":"json override","type":"prompt","prompt":"soft"}`,
		},
		{
			name: "audio sample child JSON ignores the file and stale fields",
			request: func() openai.AudioVoiceNewParams {
				child := &openai.AudioVoiceNewParamsBodyAudioSample{AudioSample: bytes.NewBufferString("ignored"), Consent: "ignored", Name: "ignored"}
				param.SetJSON([]byte(`{"name":"replacement","type":"prompt","prompt":"soft"}`), child)
				return openai.AudioVoiceNewParams{OfAudioSample: child}
			},
			json: `{"name":"replacement","type":"prompt","prompt":"soft"}`,
		},
		{
			name: "root extras overlay typed child override",
			request: func() openai.AudioVoiceNewParams {
				child := param.Override[openai.AudioVoiceNewParamsBodyPrompt](map[string]any{"name": "child", "type": "prompt", "prompt": "calm", "script_hint": "drop"})
				r := openai.AudioVoiceNewParams{OfPrompt: &child}
				r.SetExtraFields(map[string]any{"name": "root", "script_hint": param.Omit})
				return r
			},
			json: `{"name":"root","prompt":"calm","type":"prompt"}`,
		},
		{
			name: "root JSON takes priority over child JSON",
			request: func() openai.AudioVoiceNewParams {
				child := param.Override[openai.AudioVoiceNewParamsBodyPrompt](map[string]any{"name": "child", "type": "prompt", "prompt": "ignored"})
				r := openai.AudioVoiceNewParams{OfPrompt: &child}
				param.SetJSON([]byte(`{"name":"root","type":"prompt","prompt":"calm"}`), &r)
				return r
			},
			json: `{"name":"root","type":"prompt","prompt":"calm"}`,
		},
		{
			name: "child null uses the standard null parameter contract",
			request: func() openai.AudioVoiceNewParams {
				child := param.NullStruct[openai.AudioVoiceNewParamsBodyPrompt]()
				return openai.AudioVoiceNewParams{OfPrompt: &child}
			},
			json: "null",
		},
		{
			name: "untyped future variant uses only root extras",
			request: func() openai.AudioVoiceNewParams {
				r := openai.AudioVoiceNewParams{}
				r.SetExtraFields(map[string]any{
					"type": "future_fixture", "name": "from extras", "description": "calm", "consent": param.Omit,
				})
				return r
			},
			json: `{"description":"calm","name":"from extras","type":"future_fixture"}`,
		},
		{
			name: "zero value keeps the empty union behavior",
			request: func() openai.AudioVoiceNewParams {
				return openai.AudioVoiceNewParams{}
			},
			json: "null",
		},
		{
			name: "set json ignores existing variants",
			request: func() openai.AudioVoiceNewParams {
				r := openai.AudioVoiceNewParams{
					OfPrompt:      &openai.AudioVoiceNewParamsBodyPrompt{Type: "prompt", Name: "ignored", Prompt: "ignored"},
					OfAudioSample: &openai.AudioVoiceNewParamsBodyAudioSample{Name: "also ignored"},
				}
				param.SetJSON([]byte(`{"name":"json override","type":"prompt","prompt":"soft"}`), &r)
				return r
			},
			json: `{"name":"json override","type":"prompt","prompt":"soft"}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := test.request()
			oldPrompt := req.OfPrompt
			client := openai.NewClient(option.WithBaseURL("https://synthetic.example/"), option.WithAPIKey("test-key"),
				option.WithHTTPClient(&http.Client{Transport: voiceMetadataTransport(func(r *http.Request) (*http.Response, error) {
					if test.json != "" {
						b, err := io.ReadAll(r.Body)
						if err != nil {
							t.Fatal(err)
						}
						if string(b) != test.json || r.Header.Get("Content-Type") != "application/json" {
							t.Fatalf("body=%s type=%s", b, r.Header.Get("Content-Type"))
						}
					} else {
						if err := r.ParseMultipartForm(1 << 20); err != nil {
							t.Fatal(err)
						}
						defer func() { _ = r.MultipartForm.RemoveAll() }()
						if len(r.MultipartForm.Value) != len(test.want) {
							t.Fatalf("fields=%v want %v", r.MultipartForm.Value, test.want)
						}
						for key, want := range test.want {
							got := r.MultipartForm.Value[key]
							if len(got) != len(want) || len(got) != 1 || got[0] != want[0] {
								t.Fatalf("%s=%v want %v", key, got, want)
							}
						}
						if test.wantFile {
							f, _, err := r.FormFile("audio_sample")
							if err != nil {
								t.Fatal(err)
							}
							defer func() { _ = f.Close() }()
							b, err := io.ReadAll(f)
							if err != nil || string(b) != "fixture audio" {
								t.Fatalf("file=%s err=%v", b, err)
							}
						} else if len(r.MultipartForm.File) != 0 {
							t.Fatalf("unexpected files %v", r.MultipartForm.File)
						}
					}
					return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(bytes.NewBufferString(`{"id":"voice_synthetic","name":"n","type":"prompt"}`)), Request: r}, nil
				})}))
			if _, err := client.Audio.Voices.New(context.Background(), req); err != nil {
				t.Fatal(err)
			}
			if oldPrompt != nil && test.json == "" && oldPrompt.ExtraFields()["name"] != "child name" {
				t.Fatalf("mutated caller metadata: %v", oldPrompt.ExtraFields())
			}
		})
	}
}
