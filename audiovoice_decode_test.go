package openai_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestAudioVoiceDecodeAndSend(t *testing.T) {
	for _, tc := range []struct {
		name   string
		body   string
		prompt bool
	}{
		{"prompt", `{"type":"prompt","name":"Demo","prompt":"Warm narrator"}`, true},
		{"audio sample", `{"type":"audio_sample","name":"Demo","consent":"consent_demo"}`, false},
		{"default audio sample", `{"name":"Demo","consent":"consent_demo"}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var params openai.AudioVoiceNewParams
			if err := json.Unmarshal([]byte(tc.body), &params); err != nil {
				t.Fatal(err)
			}
			if (params.OfPrompt != nil) != tc.prompt || (params.OfAudioSample != nil) == tc.prompt {
				t.Fatalf("wrong branches: prompt=%t audio_sample=%t", params.OfPrompt != nil, params.OfAudioSample != nil)
			}
			if !tc.prompt {
				params.OfAudioSample.AudioSample = strings.NewReader("synthetic audio")
			}
			called := false
			client := openai.NewClient(option.WithAPIKey("test-key"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{
				Transport: voiceMetadataTransport(func(r *http.Request) (*http.Response, error) {
					called = true
					if r.URL.Path != "/v1/audio/voices" {
						t.Errorf("unexpected path: %s", r.URL.Path)
					}
					if err := r.ParseMultipartForm(1 << 20); err != nil {
						t.Fatal(err)
					}
					defer func() {
						if err := r.MultipartForm.RemoveAll(); err != nil {
							t.Errorf("remove multipart temporary files: %v", err)
						}
					}()
					if r.FormValue("name") != "Demo" {
						t.Error("name did not survive decoding")
					}
					if tc.prompt {
						if r.FormValue("prompt") != "Warm narrator" || r.FormValue("type") != "prompt" {
							t.Error("prompt fields did not survive decoding")
						}
					} else {
						if r.FormValue("consent") != "consent_demo" {
							t.Error("consent did not survive decoding")
						}
						file, _, err := r.FormFile("audio_sample")
						if err != nil {
							t.Fatal(err)
						}
						defer func() {
							if closeErr := file.Close(); closeErr != nil {
								t.Errorf("close audio file: %v", closeErr)
							}
						}()
						data, err := io.ReadAll(file)
						if err != nil || string(data) != "synthetic audio" {
							t.Error("audio file did not survive request construction")
						}
					}
					return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"id":"voice_demo","name":"Demo","created_at":0,"object":"audio.voice"}`)), Request: r}, nil
				}),
			}))
			if _, err := client.Audio.Voices.New(context.Background(), params); err != nil {
				t.Fatal(err)
			}
			if !called {
				t.Fatal("request never reached transport")
			}
		})
	}
}

func TestAudioVoiceDecodeReusedReceiver(t *testing.T) {
	var params openai.AudioVoiceNewParams
	for _, body := range []string{
		`{"type":"prompt","name":"Demo","prompt":"Warm narrator"}`,
		`{"name":"Demo","consent":"consent_demo"}`,
		`{"type":"prompt","name":"Demo","prompt":"Warm narrator"}`,
	} {
		if err := json.Unmarshal([]byte(body), &params); err != nil {
			t.Fatal(err)
		}
		wantPrompt := strings.Contains(body, `"type":"prompt"`)
		if (params.OfPrompt != nil) != wantPrompt || (params.OfAudioSample != nil) == wantPrompt {
			t.Fatal("stale union branch retained")
		}
	}
	previous := params.OfPrompt
	for _, body := range []string{`{"type":"unknown"}`, `{"type":null}`, `{"type":12}`, `{`} {
		if err := json.Unmarshal([]byte(body), &params); err == nil {
			t.Fatalf("accepted invalid input %s", body)
		}
		if params.OfPrompt != previous || params.OfAudioSample != nil {
			t.Fatal("failed decode changed receiver")
		}
	}
}
