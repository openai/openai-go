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
			name: "audio sample child JSON ignores the file and stale fields",
			request: func() openai.AudioVoiceNewParams {
				child := &openai.AudioVoiceNewParamsBodyAudioSample{AudioSample: bytes.NewBufferString("ignored"), Consent: "ignored", Name: "ignored"}
				param.SetJSON([]byte(`{"name":"replacement","type":"audio_sample"}`), child)
				return openai.AudioVoiceNewParams{OfAudioSample: child}
			},
			json: `{"name":"replacement","type":"audio_sample"}`,
		},
		{
			name: "root extras overlay typed child override",
			request: func() openai.AudioVoiceNewParams {
				child := param.Override[openai.AudioVoiceNewParamsBodyAudioSample](map[string]any{"name": "child", "type": "audio_sample", "future_field": "drop"})
				r := openai.AudioVoiceNewParams{OfAudioSample: &child}
				r.SetExtraFields(map[string]any{"name": "root", "future_field": param.Omit})
				return r
			},
			json: `{"name":"root","type":"audio_sample"}`,
		},
		{
			name: "root JSON takes priority over child JSON",
			request: func() openai.AudioVoiceNewParams {
				child := param.Override[openai.AudioVoiceNewParamsBodyAudioSample](map[string]any{"name": "child", "type": "audio_sample"})
				r := openai.AudioVoiceNewParams{OfAudioSample: &child}
				param.SetJSON([]byte(`{"name":"root","type":"audio_sample"}`), &r)
				return r
			},
			json: `{"name":"root","type":"audio_sample"}`,
		},
		{
			name: "child null uses the standard null parameter contract",
			request: func() openai.AudioVoiceNewParams {
				child := param.NullStruct[openai.AudioVoiceNewParamsBodyAudioSample]()
				return openai.AudioVoiceNewParams{OfAudioSample: &child}
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
	} {
		t.Run(test.name, func(t *testing.T) {
			req := test.request()
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
					return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(bytes.NewBufferString(`{"id":"voice_synthetic","name":"n","type":"audio_sample"}`)), Request: r}, nil
				})}))
			if _, err := client.Audio.Voices.New(context.Background(), req); err != nil {
				t.Fatal(err)
			}
		})
	}
}
