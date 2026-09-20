package webhooks_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"slices"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/webhooks"
)

func TestWebhookManagementRequestContract(t *testing.T) {
	const id = "whe/?#% 雪"
	cases := []struct {
		name   string
		suffix string
		body   map[string]any
		call   func(*openai.Client) error
	}{
		{name: "update omitted fields", body: map[string]any{}, call: func(c *openai.Client) error {
			_, err := c.Webhooks.Update(t.Context(), id, webhooks.WebhookUpdateParams{})
			return err
		}},
		{name: "rotate omitted flag", suffix: "/rotate_secret", body: map[string]any{}, call: func(c *openai.Client) error {
			_, err := c.Webhooks.RotateSecret(t.Context(), id, webhooks.WebhookRotateSecretParams{})
			return err
		}},
		{name: "rotate explicit false", suffix: "/rotate_secret", body: map[string]any{"keep_old_secret_active_for_24_hours": false}, call: func(c *openai.Client) error {
			_, err := c.Webhooks.RotateSecret(t.Context(), id, webhooks.WebhookRotateSecretParams{KeepOldSecretActiveFor24Hours: openai.Bool(false)})
			return err
		}},
		{name: "rotate explicit true", suffix: "/rotate_secret", body: map[string]any{"keep_old_secret_active_for_24_hours": true}, call: func(c *openai.Client) error {
			_, err := c.Webhooks.RotateSecret(t.Context(), id, webhooks.WebhookRotateSecretParams{KeepOldSecretActiveFor24Hours: openai.Bool(true)})
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantPath := "/webhook_endpoints/" + url.PathEscape(id) + tc.suffix
				if r.Method != http.MethodPost || r.URL.EscapedPath() != wantPath || r.URL.RawQuery != "" {
					t.Errorf("Webhooks request = %s %s, want POST %s without query", r.Method, r.URL.RequestURI(), wantPath)
				}
				if r.Header.Get("Authorization") != "Bearer test-api-key" {
					t.Error("Webhooks request did not use the explicit ordinary API credential")
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("Webhooks request body decode = %v, want valid JSON", err)
					return
				}
				if !reflect.DeepEqual(body, tc.body) {
					t.Errorf("Webhooks(%s) body = %#v, want %#v", tc.name, body, tc.body)
				}
				w.Header().Set("Content-Type", "application/json")
				if _, err := io.WriteString(w, `{"id":"whe_123","object":"webhook_endpoint","signing_secret":"synthetic-test-secret"}`); err != nil {
					t.Errorf("Writing synthetic webhook response = %v, want nil", err)
				}
			}))
			t.Cleanup(server.Close)
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("test-api-key"), option.WithAdminAPIKey("test-admin-key"), option.WithMaxRetries(0))
			if err := tc.call(&client); err != nil {
				t.Fatalf("Webhooks(%s, %q) = %v, want nil error", tc.name, id, err)
			}
		})
	}
}

func TestWebhookManagementCreateKeepsDestinationInBody(t *testing.T) {
	const destination = "https://webhook-destination.invalid/events"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/webhook_endpoints" {
			t.Errorf("Webhooks.New destination = %s %s, want POST /webhook_endpoints", r.Method, r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Webhooks.New body decode = %v, want nil", err)
			return
		}
		want := map[string]any{"name": "test endpoint", "url": destination, "event_types": []any{"response.completed"}}
		if !reflect.DeepEqual(body, want) {
			t.Errorf("Webhooks.New body = %#v, want %#v", body, want)
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := io.WriteString(w, `{"id":"whe_123","object":"webhook_endpoint","signing_secret":"synthetic-test-secret","signing_secret_hint":null}`); err != nil {
			t.Errorf("Writing synthetic creation response = %v, want nil", err)
		}
	}))
	t.Cleanup(server.Close)
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("test-api-key"), option.WithMaxRetries(0))
	endpoint, err := client.Webhooks.New(t.Context(), webhooks.WebhookNewParams{Name: "test endpoint", URL: destination, EventTypes: []string{"response.completed"}})
	if err != nil {
		t.Fatalf("Webhooks.New(%q) = %v, want nil error", destination, err)
	}
	if endpoint.SigningSecret != "synthetic-test-secret" || endpoint.JSON.SigningSecretHint.Raw() != "null" || endpoint.JSON.UpdatedAt.Raw() != "" {
		t.Error("Webhooks.New response did not preserve the signing secret, explicit-null hint, or omitted updated_at")
	}
}

func TestWebhookManagementPaginationAndEventTypes(t *testing.T) {
	var pages atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var body string
		switch r.URL.Path {
		case "/webhook_event_types":
			if r.URL.RawQuery != "" {
				t.Errorf("EventTypes.List query = %q, want no pagination parameters", r.URL.RawQuery)
			}
			body = `{"object":"list","data":["response.completed","future.event"]}`
		case "/webhook_endpoints":
			page := pages.Add(1)
			if r.URL.Query().Get("limit") != "1" {
				t.Errorf("Webhooks.List page %d limit = %q, want 1", page, r.URL.Query().Get("limit"))
			}
			if page == 1 {
				if r.URL.Query().Get("after") != "" {
					t.Errorf("Webhooks.List first cursor = %q, want omitted", r.URL.Query().Get("after"))
				}
				body = `{"object":"list","data":[{"id":"whe_first"}],"has_more":true,"first_id":"whe_first","last_id":"whe_first"}`
			} else {
				if r.URL.Query().Get("after") != "whe_first" {
					t.Errorf("Webhooks.List next cursor = %q, want whe_first", r.URL.Query().Get("after"))
				}
				body = `{"object":"list","data":[{"id":"whe_last"}],"has_more":false,"first_id":"whe_last","last_id":"whe_last"}`
			}
		default:
			t.Errorf("Webhook list path = %q, want an API list endpoint", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if _, err := io.WriteString(w, body); err != nil {
			t.Errorf("Writing synthetic list response = %v, want nil", err)
		}
	}))
	t.Cleanup(server.Close)
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("test-api-key"), option.WithMaxRetries(0))
	events, err := client.Webhooks.EventTypes.List(t.Context())
	if err != nil {
		t.Fatalf("EventTypes.List() = %v, want nil error", err)
	}
	if !slices.Equal(events.Data, []string{"response.completed", "future.event"}) {
		t.Errorf("EventTypes.List() data = %v, want known and future event names", events.Data)
	}
	iterator := client.Webhooks.ListAutoPaging(t.Context(), webhooks.WebhookListParams{Limit: openai.Int(1)})
	var ids []string
	for iterator.Next() {
		ids = append(ids, iterator.Current().ID)
	}
	if err := iterator.Err(); err != nil {
		t.Fatalf("Webhooks.ListAutoPaging(limit=1) = %v, want nil error", err)
	}
	if !slices.Equal(ids, []string{"whe_first", "whe_last"}) || pages.Load() != 2 {
		t.Errorf("Webhooks.ListAutoPaging(limit=1) = %v across %d pages, want [whe_first whe_last] across 2 pages", ids, pages.Load())
	}
}
