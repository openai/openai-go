package openai_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

type audioTextTransport func(*http.Request) (*http.Response, error)

func (f audioTextTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type audioTextBody struct {
	io.Reader
	closes int
}

func (b *audioTextBody) Close() error { b.closes++; return nil }

func audioTextClient(t *testing.T, media, payload string, status int, opts ...option.RequestOption) (openai.Client, *audioTextBody) {
	t.Helper()
	body := &audioTextBody{Reader: strings.NewReader(payload)}
	defaults := []option.RequestOption{
		option.WithBaseURL("https://audio.invalid/"), option.WithAPIKey("synthetic-key"), option.WithMaxRetries(0),
		option.WithHTTPClient(&http.Client{Transport: audioTextTransport(func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodPost || (r.URL.Path != "/audio/transcriptions" && r.URL.Path != "/audio/translations") {
				t.Errorf("audio request = %s %s, want POST to audio endpoint", r.Method, r.URL.Path)
			}
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data;") {
				t.Errorf("audio request Content-Type = %q, want multipart/form-data", r.Header.Get("Content-Type"))
			}
			return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": {media}}, Body: body, Request: r}, nil
		})}),
	}
	client := openai.NewClient(append(defaults, opts...)...)
	return client, body
}

func callAudioText(client *openai.Client, endpoint, format string, opts ...option.RequestOption) (text, raw string, present bool, err error) {
	if endpoint == "transcriptions" {
		res, callErr := client.Audio.Transcriptions.New(context.Background(), openai.AudioTranscriptionNewParams{
			File: bytes.NewBufferString("synthetic audio"), Model: "whisper-1", ResponseFormat: openai.AudioResponseFormat(format),
		}, opts...)
		if res == nil {
			return "", "", false, callErr
		}
		return res.Text, res.RawJSON(), res.JSON.Text.Valid(), callErr
	}
	res, callErr := client.Audio.Translations.New(context.Background(), openai.AudioTranslationNewParams{
		File: bytes.NewBufferString("synthetic audio"), Model: "whisper-1", ResponseFormat: openai.AudioTranslationNewParamsResponseFormat(format),
	}, opts...)
	if res == nil {
		return "", "", false, callErr
	}
	return res.Text, res.RawJSON(), res.JSON.Text.Valid(), callErr
}

func TestAudioTextResponse(t *testing.T) {
	for _, endpoint := range []string{"transcriptions", "translations"} {
		for _, format := range []string{"text", "srt", "vtt"} {
			for name, payload := range map[string]string{
				"unicode": "café 日本語\n\nlast\n", "subtitle": "1\n00:00:00,000 --> 00:00:01,000\ncaption\n",
				"json-looking": `{"text":"this is the whole transcript"}`, "empty": "",
			} {
				t.Run(endpoint+"/"+format+"/"+name, func(t *testing.T) {
					client, body := audioTextClient(t, "text/plain; charset=utf-8", payload, 200)
					var response *http.Response
					text, raw, present, err := callAudioText(&client, endpoint, format, option.WithResponseInto(&response))
					if err != nil || text != payload || raw != "" || present {
						t.Errorf("%s.New(%s) = text %q, raw %q, JSON.Text.Valid %v, err %v; want %q, empty raw, false, nil", endpoint, format, text, raw, present, err, payload)
					}
					if body.closes != 1 || response == nil || response.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
						t.Errorf("%s.New response ownership = closes %d, response %v; want one close and original response metadata", endpoint, body.closes, response)
					}
				})
			}
		}
	}
}

func TestAudioTextLargeResponse(t *testing.T) {
	// Sequential cases exercise complete bodies beyond small decoder buffers.
	payload := strings.Repeat("日本語 caption\n", 65536)
	for _, endpoint := range []string{"transcriptions", "translations"} {
		client, _ := audioTextClient(t, "text/plain", payload, 200)
		text, _, _, err := callAudioText(&client, endpoint, "text")
		if err != nil || text != payload {
			t.Errorf("%s.New large text length = %d, err %v; want %d exact bytes and nil", endpoint, len(text), err, len(payload))
		}
	}
}

func TestAudioTextJSONAndMediaCompatibility(t *testing.T) {
	for _, endpoint := range []string{"transcriptions", "translations"} {
		for _, media := range []string{"application/json", "application/problem+json", "application/json; broken"} {
			client, _ := audioTextClient(t, media, ` {"text":"ok","extra":{"value":42}} `, 200)
			text, raw, present, err := callAudioText(&client, endpoint, "text")
			if err != nil || text != "ok" || raw != `{"text":"ok","extra":{"value":42}}` || !present {
				t.Errorf("%s.New media %q = %q, %q, present %v, err %v; want intact JSON text/raw/presence", endpoint, media, text, raw, present, err)
			}
		}
		for _, media := range []string{"", "text/plain; broken", "application/octet-stream", "text/vtt"} {
			client, _ := audioTextClient(t, media, `{"text":"still not JSON media"}`, 200)
			_, _, _, err := callAudioText(&client, endpoint, "text")
			if err == nil || !strings.Contains(err.Error(), "expected destination type") {
				t.Errorf("%s.New media %q error = %v, want unsupported destination error", endpoint, media, err)
			}
		}
		client, _ := audioTextClient(t, "application/json", "{", 200)
		if _, _, _, err := callAudioText(&client, endpoint, "json"); err == nil || !strings.Contains(err.Error(), "error parsing response json") {
			t.Errorf("%s.New invalid JSON error = %v, want JSON parse error", endpoint, err)
		}
		// Response media controls decoding even when the request format differs.
		client, _ = audioTextClient(t, "text/plain", "complete transcript", 200)
		if text, _, _, err := callAudioText(&client, endpoint, "json"); err != nil || text != "complete transcript" {
			t.Errorf("%s.New(json) with plaintext = %q, %v; want full transcript and nil", endpoint, text, err)
		}
	}
}

func TestAudioTextResponseOverrides(t *testing.T) {
	for _, endpoint := range []string{"transcriptions", "translations"} {
		var clientDestination, inherited, call []byte
		client, body := audioTextClient(t, "text/plain", "complete\n", 200, option.WithResponseBodyInto(&clientDestination))
		if endpoint == "transcriptions" {
			client.Audio.Transcriptions.Options = append(client.Audio.Transcriptions.Options, option.WithResponseBodyInto(&inherited))
		} else {
			client.Audio.Translations.Options = append(client.Audio.Translations.Options, option.WithResponseBodyInto(&inherited))
		}
		text, _, _, err := callAudioText(&client, endpoint, "text", option.WithResponseBodyInto(&call))
		if err != nil || text != "" || len(clientDestination) != 0 || len(inherited) != 0 || string(call) != "complete\n" || body.closes != 1 {
			t.Errorf("%s.New override = text %q, inherited %q, call %q, closes %d, err %v; want only call destination filled and one close", endpoint, text, inherited, call, body.closes, err)
		}
		for _, nilDestination := range []bool{false, true} {
			client, body = audioTextClient(t, "text/plain", "raw\n", 200)
			var response *http.Response
			opts := []option.RequestOption{option.WithResponseBodyInto(&response)}
			if nilDestination {
				opts = []option.RequestOption{option.WithResponseBodyInto(nil), option.WithResponseInto(&response)}
			}
			text, _, _, err = callAudioText(&client, endpoint, "text", opts...)
			if err != nil || text != "" || response == nil || body.closes != 0 {
				t.Fatalf("%s.New raw override = text %q, response %v, closes %d, err %v; want unread caller-owned body", endpoint, text, response, body.closes, err)
			}
			contents, readErr := io.ReadAll(response.Body)
			closeErr := response.Body.Close()
			if readErr != nil || closeErr != nil || string(contents) != "raw\n" || body.closes != 1 {
				t.Errorf("%s.New raw body = %q, read %v, close %v, closes %d; want raw bytes and one close", endpoint, contents, readErr, closeErr, body.closes)
			}
		}
	}
}

func TestAudioTextOptionCallbacks(t *testing.T) {
	for _, endpoint := range []string{"transcriptions", "translations"} {
		var order []string
		observe := func(label string) option.RequestOption {
			return requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
				order = append(order, label)
				switch cfg.ResponseBodyInto.(type) {
				case **openai.AudioTranscriptionNewResponseUnion, **openai.Translation:
				default:
					t.Errorf("%s.New option sees %T, want original typed destination", endpoint, cfg.ResponseBodyInto)
				}
				return nil
			})
		}
		client, _ := audioTextClient(t, "text/plain", "text", 200, observe("client"))
		if endpoint == "transcriptions" {
			client.Audio.Transcriptions.Options = append(client.Audio.Transcriptions.Options, observe("service"))
		} else {
			client.Audio.Translations.Options = append(client.Audio.Translations.Options, observe("service"))
		}
		_, _, _, err := callAudioText(&client, endpoint, "text", observe("call"))
		if err != nil || strings.Join(order, ",") != "client,service,call" {
			t.Errorf("%s.New option order = %v, err %v; want client,service,call exactly once", endpoint, order, err)
		}
		wantErr := errors.New("synthetic option error")
		client, body := audioTextClient(t, "text/plain", "text", 200)
		_, _, _, err = callAudioText(&client, endpoint, "text", requestconfig.RequestOptionFunc(func(*requestconfig.RequestConfig) error { return wantErr }))
		if !errors.Is(err, wantErr) || body.closes != 0 {
			t.Errorf("%s.New option error = %v, closes %d; want original error before transport", endpoint, err, body.closes)
		}
	}
}

func TestAudioTextRichJSON(t *testing.T) {
	const payload = `{"text":"ok","duration":1.5,"language":"en","segments":[{"id":0,"text":"ok","start":0,"end":1}],"words":[{"word":"ok","start":0,"end":1}]}`
	client, _ := audioTextClient(t, "application/json", payload, 200)
	res, err := client.Audio.Transcriptions.New(t.Context(), openai.AudioTranscriptionNewParams{
		File: bytes.NewBufferString("audio"), Model: "whisper-1", ResponseFormat: openai.AudioResponseFormatVerboseJSON,
	})
	if err != nil || res == nil {
		t.Fatalf("Transcriptions.New(verbose_json) = %v, %v; want response and nil", res, err)
	}
	if res.Text != "ok" || res.Duration != 1.5 || res.Language != "en" || len(res.Segments) != 1 || len(res.Words) != 1 || res.RawJSON() != payload || !res.JSON.Duration.Valid() {
		t.Errorf("Transcriptions.New(verbose_json) = %#v; want complete rich fields and raw JSON", res)
	}
	if verbose := res.AsTranscriptionVerbose(); verbose.Text != "ok" || verbose.Duration != 1.5 || len(verbose.Words) != 1 {
		t.Errorf("AsTranscriptionVerbose() = %#v; want unchanged rich JSON variant", verbose)
	}
	for _, endpoint := range []string{"transcriptions", "translations"} {
		for _, payload := range []string{`{}`, `{"text":null}`} {
			client, _ := audioTextClient(t, "application/json", payload, 200)
			text, raw, present, err := callAudioText(&client, endpoint, "json")
			if err != nil || text != "" || raw != payload || present {
				t.Errorf("%s.New(%s) = %q, %q, present %v, %v; want unchanged absent/null JSON field", endpoint, payload, text, raw, present, err)
			}
		}
	}
}

func TestAudioTextStreamingUnchanged(t *testing.T) {
	client, body := audioTextClient(t, "text/event-stream", "data: {\"type\":\"transcript.text.delta\",\"delta\":\"hello\"}\n\ndata: [DONE]\n\n", 200)
	stream := client.Audio.Transcriptions.NewStreaming(t.Context(), openai.AudioTranscriptionNewParams{
		File: bytes.NewBufferString("audio"), Model: "gpt-4o-transcribe", ResponseFormat: openai.AudioResponseFormatJSON,
	})
	var deltas []string
	for stream.Next() {
		deltas = append(deltas, stream.Current().Delta)
	}
	if err := stream.Close(); err != nil {
		t.Fatalf("Transcriptions.NewStreaming Close() = %v, want nil", err)
	}
	if err := stream.Err(); err != nil || strings.Join(deltas, "") != "hello" || body.closes != 1 {
		t.Errorf("Transcriptions.NewStreaming = %v, closes %d, err %v; want hello, one close, nil", deltas, body.closes, err)
	}
}

type audioTextErrorReader struct{}

func (audioTextErrorReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func TestAudioTextReadError(t *testing.T) {
	for _, endpoint := range []string{"transcriptions", "translations"} {
		client, body := audioTextClient(t, "text/plain", "", 200)
		body.Reader = audioTextErrorReader{}
		_, _, _, err := callAudioText(&client, endpoint, "text")
		if !errors.Is(err, io.ErrUnexpectedEOF) || body.closes != 1 {
			t.Errorf("%s.New read error = %v, closes %d; want original read error and one close", endpoint, err, body.closes)
		}
	}
}

func TestAudioTextAPIError(t *testing.T) {
	const payload = `{"error":{"message":"synthetic rejection","type":"invalid_request_error","code":"bad_audio"}}`
	for _, endpoint := range []string{"transcriptions", "translations"} {
		client, body := audioTextClient(t, "text/plain", payload, 400)
		var response *http.Response
		text, _, _, err := callAudioText(&client, endpoint, "text", option.WithResponseInto(&response))
		var apiErr *openai.Error
		if !errors.As(err, &apiErr) || apiErr.StatusCode != 400 || text != "" || body.closes != 1 || response == nil {
			t.Fatalf("%s.New API error = %T, text %q, closes %d; want API error400, empty text and one close", endpoint, err, text, body.closes)
		}
		contents, readErr := io.ReadAll(response.Body)
		closeErr := response.Body.Close()
		if readErr != nil || closeErr != nil || string(contents) != payload {
			t.Errorf("%s.New error body = %q, read %v, close %v; want original error payload", endpoint, contents, readErr, closeErr)
		}
	}
}
