package vertex_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
	"github.com/openai/openai-go/v3/vertex"
)

// Disabled unless explicitly opted into on an entitled machine. This test
// never prints credentials, raw errors, project ID, endpoint or response body.
func TestLiveVertex(t *testing.T) {
	if os.Getenv("OPENAI_VERTEX_LIVE") != "1" {
		t.Skip("opt-in live acceptance")
	}
	project, location, model := os.Getenv("GOOGLE_CLOUD_PROJECT"), os.Getenv("GOOGLE_CLOUD_LOCATION"), os.Getenv("OPENAI_VERTEX_MODEL")
	if project == "" || location == "" || model == "" {
		t.Fatal("configure GOOGLE_CLOUD_PROJECT, GOOGLE_CLOUD_LOCATION and OPENAI_VERTEX_MODEL")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9._:/-]+$`).MatchString(model) {
		t.Fatal("model ID not safe to include in acceptance output")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	client, err := vertex.NewClient(ctx, vertex.Config{ProjectID: project, Location: location, BaseURL: os.Getenv("OPENAI_VERTEX_BASE_URL")},
		option.WithMaxRetries(0), option.WithRequestTimeout(30*time.Second))
	if err != nil {
		t.Fatal("ADC setup failed; check local Google identity and quota permissions")
	}
	sha, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil || !regexp.MustCompile(`^[a-f0-9]{40,64}$`).MatchString(strings.TrimSpace(string(sha))) {
		t.Fatal("unable to record source revision")
	}
	env := os.Getenv("OPENAI_VERTEX_ENVIRONMENT")
	if env == "" {
		env = "unspecified"
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`).MatchString(env) {
		t.Fatal("environment must be a short non-sensitive alias")
	}
	t.Logf("utc=%s source_sha=%s sdk=%s model=%s environment=%s",
		time.Now().UTC().Format(time.RFC3339), strings.TrimSpace(string(sha)), internal.PackageVersion, model, env)
	input := responses.ResponseNewParamsInputUnion{OfString: openai.String("Reply with the word ready.")}
	params := responses.ResponseNewParams{Model: shared.ResponsesModel(model), Input: input, Store: openai.Bool(false)}
	t.Run("create", func(t *testing.T) {
		r, err := client.Responses.New(ctx, params)
		if err != nil {
			liveFailure(t, err)
		}
		if r == nil || r.ID == "" || r.Status != "completed" || len(r.Output) == 0 {
			t.Fatal("create did not return completed output")
		}
	})
	t.Run("stream", func(t *testing.T) {
		stream := client.Responses.NewStreaming(ctx, params)
		defer closeTestStream(t, stream)
		complete := false
		for stream.Next() {
			if stream.Current().Type == "response.completed" {
				complete = true
			}
		}
		if stream.Err() != nil {
			liveFailure(t, stream.Err())
		}
		if !complete {
			t.Fatal("stream did not complete")
		}
	})
	t.Run("tools", func(t *testing.T) {
		p := params
		p.Tools = []responses.ToolUnionParam{{OfFunction: &responses.FunctionToolParam{
			Name: "lookup", Parameters: map[string]any{"type": "object", "properties": map[string]any{}, "required": []string{}, "additionalProperties": false},
		}}}
		p.ToolChoice = responses.ResponseNewParamsToolChoiceUnion{OfFunctionTool: &responses.ToolChoiceFunctionParam{Name: "lookup"}}
		r, err := client.Responses.New(ctx, p)
		if err != nil {
			liveFailure(t, err)
		}
		if r == nil || r.Status != "completed" {
			t.Fatal("function tool request did not complete")
		}
		// Replay all returned context (including reasoning items) without
		// requiring storage. Never execute model-supplied arguments.
		var input responses.ResponseInputParam
		calls := 0
		for _, output := range r.Output {
			input = append(input, param.Override[responses.ResponseInputItemUnionParam](json.RawMessage(output.RawJSON())))
			if output.Type == "function_call" {
				call := output.AsFunctionCall()
				if call.Name != "lookup" || call.CallID == "" {
					t.Fatal("unexpected function call")
				}
				calls++
				result := responses.ResponseInputItemParamOfFunctionCallOutput("probe-ok")
				result.OfFunctionCallOutput.CallID = openai.String(call.CallID)
				input = append(input, result)
			}
		}
		if calls == 0 {
			t.Fatal("expected a lookup function call")
		}
		next := params
		next.Input = responses.ResponseNewParamsInputUnion{OfInputItemList: input}
		r, err = client.Responses.New(ctx, next)
		if err != nil {
			liveFailure(t, err)
		}
		if r == nil || r.Status != "completed" || len(r.Output) == 0 {
			t.Fatal("function output roundtrip did not complete")
		}
	})
	t.Run("stateful", func(t *testing.T) {
		if os.Getenv("OPENAI_VERTEX_LIVE_STATEFUL") != "1" {
			t.Skip("set OPENAI_VERTEX_LIVE_STATEFUL=1 only for a deployment with storage")
		}
		p := params
		p.Store = openai.Bool(true)
		r, err := client.Responses.New(ctx, p)
		if err != nil {
			liveFailure(t, err)
		}
		if r == nil || r.ID == "" {
			t.Fatal("stored create returned no response ID")
		}
		defer func() {
			cleanup, finish := context.WithTimeout(context.Background(), 30*time.Second)
			defer finish()
			if cleanupErr := client.Responses.Delete(cleanup, r.ID); cleanupErr != nil {
				// Do not let the original deadline prevent deletion.
				t.Errorf("cleanup delete failed class=%T", cleanupErr)
			}
		}()
		if r.Status != "completed" {
			t.Fatal("stored create did not complete")
		}
		got, err := client.Responses.Get(ctx, r.ID, responses.ResponseGetParams{})
		if err != nil {
			liveFailure(t, err)
		}
		if got == nil || got.ID != r.ID || got.Status != "completed" {
			t.Fatal("retrieve did not return the stored response")
		}
		next := params
		next.PreviousResponseID = openai.String(r.ID)
		followup, err := client.Responses.New(ctx, next)
		if err != nil {
			liveFailure(t, err)
		}
		if followup == nil || followup.Status != "completed" {
			t.Fatal("previous_response_id did not complete")
		}
	})
	t.Run("compact", func(t *testing.T) {
		if os.Getenv("OPENAI_VERTEX_LIVE_COMPACT") != "1" {
			t.Skip("set OPENAI_VERTEX_LIVE_COMPACT=1 for a deployment with compaction")
		}
		compacted, err := client.Responses.Compact(ctx, responses.ResponseCompactParams{
			Model: responses.ResponseCompactParamsModel(model),
			Input: responses.ResponseCompactParamsInputUnion{OfString: openai.String("Reply ready.")},
		})
		if err != nil {
			liveFailure(t, err)
		}
		if compacted == nil || compacted.ID == "" || compacted.Object != "response.compaction" || len(compacted.Output) == 0 {
			t.Fatal("compact did not return compaction output")
		}
	})
}

// Only operation, class, and HTTP status enter logs. Do not stringify errors:
// they may contain request/response bodies or credential details.
func liveFailure(t *testing.T, err error) {
	t.Helper()
	var status *openai.Error
	if errors.As(err, &status) {
		t.Fatalf("request failed class=%T status=%d", err, status.StatusCode)
	}
	t.Fatalf("request failed class=%T", err)
}
