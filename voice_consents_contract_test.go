package openai_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type consentTransport struct {
	t        *testing.T
	requests []*http.Request
	bodies   [][]byte
}

func consentJSON(id string) string {
	return `{"id":"` + id + `","object":"audio.voice_consent","name":"Owned","language":"en-US","created_at":1734220800}`
}
func (s *consentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	b := []byte{}
	if r.Body != nil {
		b, _ = io.ReadAll(r.Body)
	}
	s.requests = append(s.requests, r)
	s.bodies = append(s.bodies, b)
	status := 200
	body := consentJSON("cons_1")
	switch {
	case strings.HasSuffix(r.URL.Path, "/missing"):
		status = 404
		body = `{"error":{"message":"Missing consent","type":"invalid_request_error"}}`
	case r.Method == "DELETE":
		body = `{"id":"cons_1","object":"audio.voice_consent","deleted":true}`
	case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/voice_consents"):
		switch r.URL.Query().Get("after") {
		case "":
			body = `{"object":"list","data":[` + consentJSON("cons_1") + `],"has_more":true}`
		case "cons_1":
			body = `{"object":"list","data":[` + consentJSON("cons_2") + `],"has_more":false}`
		default:
			body = `{"object":"list","data":[],"has_more":false}`
		}
	case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/cons_1"):
		body = strings.Replace(body, "Owned", "Renamed", 1)
	}
	return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
}

func TestVoiceConsentOfflineContract(t *testing.T) {
	transport := &consentTransport{t: t}
	client := openai.NewClient(option.WithAPIKey("sk-synthetic"), option.WithBaseURL("https://unit.invalid/v1/"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{Transport: transport}))
	ctx := context.Background()
	raw := []byte{'R', 'I', 'F', 'F', 0, 255}
	created, err := client.Audio.VoiceConsents.New(ctx, openai.AudioVoiceConsentNewParams{Name: "Owned", Language: "en-US", Recording: openai.File(bytes.NewReader(raw), "consent.wav", "audio/wav")})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID != "cons_1" || created.Language != "en-US" {
		t.Fatalf("incorrect consent result")
	}
	first := transport.requests[0]
	if first.Method != "POST" || first.URL.Path != "/v1/audio/voice_consents" {
		t.Fatal("wrong upload route")
	}
	_, params, err := mime.ParseMediaType(first.Header.Get("Content-Type"))
	if err != nil {
		t.Fatal(err)
	}
	reader := multipart.NewReader(bytes.NewReader(transport.bodies[0]), params["boundary"])
	values := map[string][]byte{}
	for {
		part, e := reader.NextPart()
		if e == io.EOF {
			break
		}
		if e != nil {
			t.Fatal(e)
		}
		values[part.FormName()], _ = io.ReadAll(part)
		if part.FormName() == "recording" && (part.FileName() != "consent.wav" || part.Header.Get("Content-Type") != "audio/wav") {
			t.Fatal("lost upload filename or MIME")
		}
	}
	if !bytes.Equal(values["recording"], raw) || string(values["name"]) != "Owned" || string(values["language"]) != "en-US" {
		t.Fatal("lost multipart bytes or metadata")
	}
	got, err := client.Audio.VoiceConsents.Get(ctx, "cons_1")
	if err != nil || got.ID != "cons_1" {
		t.Fatal("retrieve failed", err)
	}
	renamed, err := client.Audio.VoiceConsents.Update(ctx, "cons_1", openai.AudioVoiceConsentUpdateParams{Name: "Renamed"})
	if err != nil || renamed.Name != "Renamed" {
		t.Fatal("update failed", err)
	}
	var update map[string]string
	if json.Unmarshal(transport.bodies[len(transport.bodies)-1], &update) != nil || update["name"] != "Renamed" {
		t.Fatal("wrong update body")
	}
	deleted, err := client.Audio.VoiceConsents.Delete(ctx, "cons_1")
	if err != nil || !deleted.Deleted {
		t.Fatal("delete failed", err)
	}
	pager := client.Audio.VoiceConsents.ListAutoPaging(ctx, openai.AudioVoiceConsentListParams{Limit: openai.Int(1)})
	ids := []string{}
	for pager.Next() {
		ids = append(ids, pager.Current().ID)
	}
	if pager.Err() != nil || strings.Join(ids, ",") != "cons_1,cons_2" {
		t.Fatal("pagination failed", pager.Err())
	}
	query := transport.requests[len(transport.requests)-1].URL.Query()
	if query.Get("after") != "cons_1" || query.Get("limit") != "1" {
		t.Fatal("wrong cursor or lost limit")
	}
	_, err = client.Audio.VoiceConsents.Get(ctx, "missing")
	var apiErr *openai.Error
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 404 {
		t.Fatal("wrong typed error")
	}
}
