package vertex_test

import (
	"cloud.google.com/go/auth/credentials"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/vertex"
)

func TestEnvironmentDefaultsAndExplicitVertexConfiguration(t *testing.T) {
	t.Setenv("GOOGLE_CLOUD_PROJECT", "env-project")
	t.Setenv("GOOGLE_CLOUD_LOCATION", "global")
	token := func(context.Context) (string, error) { return "fake", nil }
	for _, tc := range []struct {
		name   string
		config vertex.Config
		want   string
	}{
		{"env", vertex.Config{TokenProvider: token}, "https://aiplatform.googleapis.com/v1/projects/env-project/locations/global/publishers/openai/v1/responses"},
		{"explicit project", vertex.Config{ProjectID: "my-project", TokenProvider: token}, "https://aiplatform.googleapis.com/v1/projects/my-project/locations/global/publishers/openai/v1/responses"},
		{"domain-scoped project", vertex.Config{ProjectID: "example.com:my-project", TokenProvider: token}, "https://aiplatform.googleapis.com/v1/projects/example.com:my-project/locations/global/publishers/openai/v1/responses"},
		{"explicit region", vertex.Config{Location: "us-east5", TokenProvider: token}, "https://us-east5-aiplatform.googleapis.com/v1/projects/env-project/locations/us-east5/publishers/openai/v1/responses"},
		{"explicit base", vertex.Config{BaseURL: "https://gateway.invalid/publisher/v1", TokenProvider: token}, "https://gateway.invalid/publisher/v1/responses"},
		{"escaped base", vertex.Config{BaseURL: "https://gateway.invalid/publisher%2Fv1", TokenProvider: token}, "https://gateway.invalid/publisher%2Fv1/responses"},
		{"http test base", vertex.Config{BaseURL: "http://test-server.invalid/publisher/v1", TokenProvider: token}, "http://test-server.invalid/publisher/v1/responses"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got string
			client, err := vertex.NewClient(context.Background(), tc.config, option.WithHTTPClient(&http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					got = r.URL.String()
					return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}},
						Body: io.NopCloser(strings.NewReader(`{"id":"resp_test","object":"response","output":[]}`)), Request: r}, nil
				}),
			}))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := client.Responses.New(context.Background(), responses.ResponseNewParams{Model: "model"}); err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("URL %q, want %q", got, tc.want)
			}
		})
	}
	t.Setenv("GOOGLE_CLOUD_PROJECT", "../invalid")
	if _, err := vertex.NewClient(context.Background(), vertex.Config{TokenProvider: token}); err == nil {
		t.Fatal("accepted invalid environment project")
	}
	if _, err := vertex.NewClient(context.Background(), vertex.Config{ProjectID: "my-project", Location: "us-east5.example", TokenProvider: token}); err == nil {
		t.Fatal("accepted a location containing a domain separator")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func writeTestResponse(t *testing.T, w http.ResponseWriter, body string) {
	t.Helper()
	if _, err := io.WriteString(w, body); err != nil {
		t.Error(err)
	}
}

func closeTestStream(t *testing.T, stream io.Closer) {
	t.Helper()
	if err := stream.Close(); err != nil {
		t.Errorf("stream close failed class=%T", err)
	}
}

func TestVertexResidencyConflictsBeforeCredentialUse(t *testing.T) {
	calls := 0
	cfg := vertex.Config{BaseURL: "https://trusted.invalid/v1", TokenProvider: func(context.Context) (string, error) {
		calls++
		return "fake", nil
	}}
	if _, err := vertex.NewClient(context.Background(), cfg, option.WithDataResidency(option.DataResidencyEU)); err == nil {
		t.Error("accepted OpenAI residency on a Vertex client")
	}
	client, err := vertex.NewClient(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Responses.New(context.Background(), responses.ResponseNewParams{Model: "model"},
		option.WithDataResidency(option.DataResidencyEU)); err == nil {
		t.Error("accepted per-request OpenAI residency on Vertex")
	}
	if calls != 0 {
		t.Errorf("credential callbacks=%d, want zero", calls)
	}
}

func TestFunctionReplayKeepsReasoningAndAdditionalFields(t *testing.T) {
	calls := 0
	fixture := `[{"type":"reasoning","id":"rs_test","summary":[],"encrypted_content":"fake-sealed","new_field":{"value":1}},{"type":"function_call","id":"fc_test","call_id":"call_test","name":"lookup","arguments":"{}","status":"completed"}]`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if calls == 1 {
			writeTestResponse(t, w, fmt.Sprintf(`{"id":"resp_first","object":"response","output":%s}`, fixture))
			return
		}
		var body struct {
			Input []json.RawMessage `json:"input"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
			return
		}
		var expected []json.RawMessage
		if err := json.Unmarshal([]byte(fixture), &expected); err != nil {
			t.Error(err)
			return
		}
		if len(body.Input) != 3 {
			t.Errorf("input count=%d", len(body.Input))
			return
		}
		for i := range expected {
			var got, want any
			if err := json.Unmarshal(body.Input[i], &got); err != nil {
				t.Error(err)
			}
			if err := json.Unmarshal(expected[i], &want); err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("replay item %d changed", i)
			}
		}
		var result struct {
			CallID string `json:"call_id"`
			Output string `json:"output"`
		}
		if err := json.Unmarshal(body.Input[2], &result); err != nil {
			t.Error(err)
		}
		if result.CallID != "call_test" || result.Output != "ok" {
			t.Error("function output lost call identity")
		}
		writeTestResponse(t, w, `{"id":"resp_second","object":"response","output":[]}`)
	}))
	defer server.Close()
	client, err := vertex.NewClient(context.Background(), vertex.Config{
		BaseURL: server.URL, TokenProvider: func(context.Context) (string, error) { return "fake", nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	p := responses.ResponseNewParams{Model: "model"}
	r, err := client.Responses.New(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
	var input responses.ResponseInputParam
	for _, output := range r.Output {
		input = append(input, param.Override[responses.ResponseInputItemUnionParam](json.RawMessage(output.RawJSON())))
	}
	result := responses.ResponseInputItemParamOfFunctionCallOutput("ok")
	result.OfFunctionCallOutput.CallID = openai.String(r.Output[1].AsFunctionCall().CallID)
	input = append(input, result)
	p.Input = responses.ResponseNewParamsInputUnion{OfInputItemList: input}
	if _, err := client.Responses.New(context.Background(), p); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Errorf("requests=%d", calls)
	}
}

func TestResponsesTransport(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "fake-ambient-key")
	t.Setenv("OPENAI_BASE_URL", "https://wrong.invalid/v1/")
	t.Setenv("OPENAI_ORG_ID", "fake-ambient-org")
	calls := 0
	tokens := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Header.Get("Authorization") != fmt.Sprintf("Bearer fake-%d", calls) {
			t.Error("expected a fresh synthetic credential per attempt")
		}
		if r.Header.Get("OpenAI-Organization") != "" {
			t.Error("ambient organization forwarded")
		}
		expected := "/v1/projects/test/locations/us-east5/publishers/openai/v1/responses"
		if r.Method != "POST" {
			expected += "/resp_test"
		}
		if r.URL.Path != expected {
			t.Errorf("path = %s; want %s", r.URL.Path, expected)
		}
		if calls == 1 {
			w.Header().Set("retry-after", "0")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		writeTestResponse(t, w, `{"id":"resp_test","object":"response","output":[]}`)
	}))
	defer server.Close()
	client, err := vertex.NewClient(context.Background(), vertex.Config{
		BaseURL:       server.URL + "/v1/projects/test/locations/us-east5/publishers/openai/v1",
		TokenProvider: func(context.Context) (string, error) { tokens++; return fmt.Sprintf("fake-%d", tokens), nil },
	}, option.WithMaxRetries(1))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Responses.New(context.Background(), responses.ResponseNewParams{
		Model: "model-selected-by-caller",
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("hello")},
	})
	if err != nil || resp.ID != "resp_test" {
		t.Fatalf("create failed %v", err)
	}
	if _, err := client.Responses.Get(context.Background(), "resp_test", responses.ResponseGetParams{}); err != nil {
		t.Fatal(err)
	}
	if err := client.Responses.Delete(context.Background(), "resp_test"); err != nil {
		t.Fatal(err)
	}
	if calls != 4 || tokens != 4 {
		t.Errorf("calls=%d tokens=%d", calls, tokens)
	}
}

func TestLifecyclePreservesEscapedResponseIDs(t *testing.T) {
	for _, tc := range []struct{ id, escaped string }{
		{"resp_test", "resp_test"},
		{"resp/a", "resp%2Fa"},
		{"..", "%2E%2E"},
	} {
		t.Run(tc.id, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				want := "/v1/projects/test/locations/global/publishers/openai/v1/responses/" + tc.escaped
				if r.RequestURI != want {
					t.Errorf("request URI = %q, want %q", r.RequestURI, want)
				}
				w.Header().Set("Content-Type", "application/json")
				writeTestResponse(t, w, `{}`)
			}))
			defer server.Close()
			client, err := vertex.NewClient(context.Background(), vertex.Config{
				BaseURL:       server.URL + "/v1/projects/test/locations/global/publishers/openai/v1",
				TokenProvider: func(context.Context) (string, error) { return "fake-token", nil },
			})
			if err != nil {
				t.Fatal(err)
			}
			if _, err := client.Responses.Get(context.Background(), tc.id, responses.ResponseGetParams{}); err != nil {
				t.Fatal(err)
			}
			if err := client.Responses.Delete(context.Background(), tc.id); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestCredentialIsolationAndRedirects(t *testing.T) {
	tokens, requests := 0, 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		http.Redirect(w, r, "/sink", http.StatusTemporaryRedirect)
	}))
	defer server.Close()
	cfg := vertex.Config{
		BaseURL: server.URL + "/v1", TokenProvider: func(context.Context) (string, error) {
			tokens++
			return "fake-token", nil
		},
	}
	client, err := vertex.NewClient(context.Background(), cfg, option.WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}
	args := responses.ResponseNewParams{Model: "caller-model"}
	for _, opt := range []option.RequestOption{
		option.WithAPIKey("fake-conflict"), option.WithHeader("Authorization", "Bearer fake-conflict"),
		option.WithHeader("aPi-KeY", "fake-conflict"),
		option.WithBaseURL("https://wrong.invalid/v1"),
	} {
		_, requestErr := client.Responses.New(context.Background(), args, opt)
		if requestErr == nil {
			t.Fatal("expected credential or route conflict")
		}
	}
	clientWithInheritedKey, err := vertex.NewClient(context.Background(), cfg, option.WithHeader("api-key", "fake-conflict"))
	if err != nil {
		t.Fatal(err)
	}
	if _, requestErr := clientWithInheritedKey.Responses.New(context.Background(), args); requestErr == nil {
		t.Fatal("accepted inherited api-key")
	}
	if tokens != 0 || requests != 0 {
		t.Fatal("credentials resolved or request sent before conflict rejected")
	}
	_, err = client.Responses.New(context.Background(), args)
	if err == nil {
		t.Fatal("expected the redirect response to be returned without following it")
	}
	if tokens != 1 || requests != 1 {
		t.Fatal("authenticated redirect was followed")
	}
}

func TestDeleteRejectsRedirect(t *testing.T) {
	const sensitive = "fake-sensitive-value"
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if r.URL.Path == "/sink" {
			t.Error("redirect target was contacted")
			return
		}
		http.Redirect(w, r, "/sink?token="+sensitive, http.StatusTemporaryRedirect)
	}))
	defer server.Close()
	client, err := vertex.NewClient(context.Background(), vertex.Config{
		BaseURL: server.URL, TokenProvider: func(context.Context) (string, error) { return "fake-token", nil },
	}, option.WithMaxRetries(1))
	if err != nil {
		t.Fatal(err)
	}
	err = client.Responses.Delete(context.Background(), sensitive)
	if err == nil {
		t.Fatal("expected refused redirect to fail")
	}
	if strings.Contains(err.Error(), sensitive) {
		t.Fatal("redirect error contained sensitive URL data")
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("requests = %d, want 1", got)
	}
}

func TestGoogleCredentialsRefreshQuotaAndDeadline(t *testing.T) {
	var refreshes, requests atomic.Int32
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := refreshes.Add(1)
		w.Header().Set("Content-Type", "application/json")
		// The first token is already expired, so the Google cache must refresh
		// before the parallel calls and reuse the second token.
		expiry := 3600
		if n == 1 {
			expiry = -1
		}
		writeTestResponse(t, w, fmt.Sprintf(`{"access_token":"fake-%d","token_type":"Bearer","expires_in":%d}`, n, expiry))
	}))
	defer authServer.Close()
	creds, err := credentials.DetectDefault(&credentials.DetectOptions{
		Scopes:              []string{"https://www.googleapis.com/auth/cloud-platform"},
		CredentialsJSON:     []byte(`{"type":"authorized_user","client_id":"fake-id","client_secret":"fake-secret","refresh_token":"fake-refresh","quota_project_id":"fake-billing"}`),
		TokenURL:            authServer.URL,
		DisableAsyncRefresh: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if r.Header.Get("X-Goog-User-Project") != "fake-billing" {
			t.Error("missing ADC quota project")
		}
		if r.Header.Get("OpenAI-Organization") != "" || r.Header.Get("OpenAI-Project") != "" {
			t.Error("Vertex emitted OpenAI routing metadata")
		}
		w.Header().Set("Content-Type", "application/json")
		writeTestResponse(t, w, `{"id":"resp_test","object":"response","output":[]}`)
	}))
	defer server.Close()
	client, err := vertex.NewClient(context.Background(), vertex.Config{BaseURL: server.URL + "/v1", Credentials: creds},
		option.WithOrganization("fake-openai-org"), option.WithProject("fake-openai-project"))
	if err != nil {
		t.Fatal(err)
	}
	call := func() {
		_, requestErr := client.Responses.New(context.Background(), responses.ResponseNewParams{Model: "caller-model"},
			option.WithHeader("openai-project", "fake-method-project"))
		if requestErr != nil {
			t.Errorf("request: %v", requestErr)
		}
	}
	call()
	var group sync.WaitGroup
	for range 8 {
		group.Go(call)
	}
	group.Wait()
	if refreshes.Load() != 2 || requests.Load() != 9 {
		t.Errorf("refreshes=%d requests=%d", refreshes.Load(), requests.Load())
	}
	if _, configErr := vertex.NewClient(context.Background(), vertex.Config{BaseURL: server.URL, Credentials: creds,
		TokenProvider: func(context.Context) (string, error) { return "fake", nil }}); configErr == nil {
		t.Error("accepted competing auth")
	}

	unblock := make(chan struct{})
	blocker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-unblock:
		}
	}))
	defer blocker.Close()
	defer close(unblock)
	slow, err := credentials.DetectDefault(&credentials.DetectOptions{
		Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
		CredentialsJSON: []byte(`{"type":"authorized_user","client_id":"fake","client_secret":"fake","refresh_token":"fake"}`),
		TokenURL:        blocker.URL, DisableAsyncRefresh: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	client, err = vertex.NewClient(context.Background(), vertex.Config{BaseURL: server.URL + "/v1", Credentials: slow},
		option.WithRequestTimeout(50*time.Millisecond), option.WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	_, err = client.Responses.New(context.Background(), responses.ResponseNewParams{Model: "caller-model"})
	if !errors.Is(err, context.DeadlineExceeded) || time.Since(start) > time.Second {
		t.Fatalf("credential deadline failed, type=%T elapsed=%s", err, time.Since(start))
	}
	if requests.Load() != 9 {
		t.Fatal("API request sent after auth timeout")
	}
}

func TestStreamErrorAndToolsAreNotReplayed(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["stream"] != true || body["model"] != "caller-model" || len(body["tools"].([]any)) != 1 {
			t.Error("normal generated wire shape lost")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		writeTestResponse(t, w, "event: response.created\ndata: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_test\",\"object\":\"response\",\"status\":\"in_progress\",\"output\":[]}}\n\n")
		writeTestResponse(t, w, "event: error\ndata: {\"type\":\"error\",\"code\":\"unsupported_feature\",\"message\":\"fake unsupported\",\"param\":null,\"sequence_number\":1}\n\n")
	}))
	defer server.Close()
	client, err := vertex.NewClient(context.Background(), vertex.Config{BaseURL: server.URL + "/v1",
		TokenProvider: func(context.Context) (string, error) { return "fake", nil }}, option.WithMaxRetries(3))
	if err != nil {
		t.Fatal(err)
	}
	stream := client.Responses.NewStreaming(context.Background(), responses.ResponseNewParams{
		Model: "caller-model",
		Tools: []responses.ToolUnionParam{{OfFunction: &responses.FunctionToolParam{
			Name: "lookup", Parameters: map[string]any{"type": "object", "properties": map[string]any{}},
		}}},
	})
	defer closeTestStream(t, stream)
	events := 0
	sawError := false
	for stream.Next() {
		events++
		if stream.Current().Type == "error" || stream.Current().Type == "response.failed" {
			sawError = true
		}
	}
	// The standard stream surfaces typed events or an SDK stream error.
	if events < 1 {
		t.Fatal("no initial event was surfaced")
	}
	if !sawError && stream.Err() == nil {
		t.Fatal("stream silently dropped its error")
	}
	if calls != 1 {
		t.Fatalf("partially delivered stream was replayed %d times", calls)
	}
}
