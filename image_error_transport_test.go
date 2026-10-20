package openai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/conversations"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

// These assignments deliberately compile unchanged against the released numeric APIs.
func TestImageErrorPreservesReleasedProjectionTypes(t *testing.T) {
	var c conversations.ConversationItemUnion
	var i responses.ResponseInputItemUnion
	var r responses.ResponseItemUnion
	var o responses.ResponseOutputItemUnion
	var p responses.ResponseInputItemUnionParam
	var _ int64 = c.Error.Code
	var _ int64 = i.Error.Code
	var _ int64 = r.Error.Code
	var _ int64 = o.Error.Code
	var _ *int64 = p.GetError().GetCode()
	_ = c.Error.Code == -32603
}

type imageErrorHTTPDoer func(*http.Request) (*http.Response, error)

func (f imageErrorHTTPDoer) Do(r *http.Request) (*http.Response, error) { return f(r) }

func imageErrorClient(t *testing.T, method, path, body string, inspect func(*http.Request)) openai.Client {
	t.Helper()
	calls := 0
	t.Cleanup(func() {
		if calls != 1 {
			t.Errorf("transport calls = %d, want 1", calls)
		}
	})
	return openai.NewClient(option.WithAPIKey("synthetic-key"), option.WithBaseURL("https://example.test/v1"), option.WithMaxRetries(0), option.WithHTTPClient(imageErrorHTTPDoer(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != method || r.URL.Path != path {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if inspect != nil {
			inspect(r)
		}
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})))
}

func TestImageErrorPreservesNumericMcpResponses(t *testing.T) {
	for _, code := range []int64{-32603, 0, 9007199254740993, math.MaxInt64} {
		for _, route := range []string{"response", "input_items", "conversation"} {
			t.Run(route+"/"+strconv.FormatInt(code, 10), func(t *testing.T) {
				item := fmt.Sprintf(`{"id":"item_synthetic","type":"mcp_call","error":{"type":"mcp_protocol_error","code":%d,"message":"synthetic"}}`, code)
				var actual int64
				var present bool
				switch route {
				case "response":
					client := imageErrorClient(t, "GET", "/v1/responses/resp_synthetic", `{"id":"resp_synthetic","output":[`+item+`],"instructions":[`+item+`]}`, nil)
					res, err := client.Responses.Get(context.Background(), "resp_synthetic", responses.ResponseGetParams{})
					if err != nil {
						t.Fatal(err)
					}
					actual, present = res.Output[0].Error.Code, res.Output[0].Error.JSON.Code.Valid()
					input := res.Instructions.AsInputItemList()[0]
					if input.Error.Code != code || !input.Error.JSON.Code.Valid() {
						t.Fatalf("input projection lost numeric code %d", code)
					}
				case "input_items":
					client := imageErrorClient(t, "GET", "/v1/responses/resp_synthetic/input_items", `{"data":[`+item+`],"has_more":false}`, nil)
					res, err := client.Responses.InputItems.List(context.Background(), "resp_synthetic", responses.InputItemListParams{})
					if err != nil {
						t.Fatal(err)
					}
					actual, present = res.Data[0].Error.Code, res.Data[0].Error.JSON.Code.Valid()
				case "conversation":
					client := imageErrorClient(t, "GET", "/v1/conversations/conv_synthetic/items/item_synthetic", item, nil)
					res, err := client.Conversations.Items.Get(context.Background(), "conv_synthetic", "item_synthetic", conversations.ItemGetParams{})
					if err != nil {
						t.Fatal(err)
					}
					actual, present = res.Error.Code, res.Error.JSON.Code.Valid()
				}
				if actual != code || !present {
					t.Fatalf("numeric code = %d (present=%v), want %d", actual, present, code)
				}
			})
		}
	}
}

type imageErrorView struct {
	code, param, message, kind, stage, rawCode, rawParam string
	categories                                           []string
	errorPresent, moderationPresent                      bool
}

func TestImageErrorTypedTransport(t *testing.T) {
	for _, route := range []string{"response", "beta", "conversation"} {
		for _, mode := range []string{"values", "nulls", "no_moderation", "no_error"} {
			t.Run(route+"/"+mode, func(t *testing.T) {
				code, parameter := `"synthetic_code"`, `"synthetic_param"`
				if mode == "nulls" {
					code, parameter = "null", "null"
				}
				moderation := `,"moderation_details":{"moderation_stage":"output","categories":["synthetic_category"]}`
				if mode == "no_moderation" {
					moderation = ""
				}
				errorJSON := `,"error":{"code":` + code + `,"param":` + parameter + `,"message":"synthetic blocked","type":"image_generation_error"` + moderation + `}`
				if mode == "no_error" {
					errorJSON = ""
				}
				item := `{"id":"img_synthetic","type":"image_generation_call","status":"failed","result":null` + errorJSON + `}`
				body := `{"id":"resp_synthetic","output":[` + item + `]}`
				var got imageErrorView
				switch route {
				case "response":
					client := imageErrorClient(t, "GET", "/v1/responses/resp_synthetic", body, nil)
					res, err := client.Responses.Get(context.Background(), "resp_synthetic", responses.ResponseGetParams{})
					if err != nil {
						t.Fatal(err)
					}
					call := res.Output[0].AsImageGenerationCall()
					e := call.Error
					got = imageErrorView{e.Code, e.Param, e.Message, e.Type, e.ModerationDetails.ModerationStage, e.JSON.Code.Raw(), e.JSON.Param.Raw(), e.ModerationDetails.Categories, call.JSON.Error.Valid(), e.JSON.ModerationDetails.Valid()}
					if res.Output[0].Error.JSON.Code.Valid() {
						t.Fatal("image code must not masquerade as a numeric projection")
					}
				case "beta":
					client := imageErrorClient(t, "GET", "/v1/responses/resp_synthetic", body, func(r *http.Request) {
						if r.URL.Query().Get("beta") != "true" {
							t.Fatal("missing beta query")
						}
					})
					res, err := client.Beta.Responses.Get(context.Background(), "resp_synthetic", openai.BetaResponseGetParams{})
					if err != nil {
						t.Fatal(err)
					}
					call := res.Output[0].AsImageGenerationCall()
					e := call.Error
					got = imageErrorView{e.Code, e.Param, e.Message, e.Type, e.ModerationDetails.ModerationStage, e.JSON.Code.Raw(), e.JSON.Param.Raw(), e.ModerationDetails.Categories, call.JSON.Error.Valid(), e.JSON.ModerationDetails.Valid()}
				case "conversation":
					client := imageErrorClient(t, "GET", "/v1/conversations/conv_synthetic/items/img_synthetic", item, nil)
					res, err := client.Conversations.Items.Get(context.Background(), "conv_synthetic", "img_synthetic", conversations.ItemGetParams{})
					if err != nil {
						t.Fatal(err)
					}
					call := res.AsImageGenerationCall()
					e := call.Error
					got = imageErrorView{e.Code, e.Param, e.Message, e.Type, e.ModerationDetails.ModerationStage, e.JSON.Code.Raw(), e.JSON.Param.Raw(), e.ModerationDetails.Categories, call.JSON.Error.Valid(), e.JSON.ModerationDetails.Valid()}
				}
				want := imageErrorView{"synthetic_code", "synthetic_param", "synthetic blocked", "image_generation_error", "output", code, parameter, []string{"synthetic_category"}, true, true}
				if mode == "nulls" {
					want.code, want.param = "", ""
				}
				if mode == "no_moderation" {
					want.stage, want.categories, want.moderationPresent = "", nil, false
				}
				if mode == "no_error" {
					want = imageErrorView{}
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("typed image error = %#v, want %#v", got, want)
				}
			})
		}
	}
}

func TestImageErrorPreservesRequestSerialization(t *testing.T) {
	mcp := responses.ResponseInputItemUnionParam{OfMcpCall: &responses.ResponseInputItemMcpCallParam{ID: "mcp_synthetic", Arguments: "{}", Name: "synthetic", ServerLabel: "synthetic", Error: responses.McpToolCallErrorUnionParam{OfMcpProtocolError: &responses.McpToolCallErrorMcpProtocolErrorParam{Code: -32603, Message: "synthetic"}}}}
	// The released flattened GetCode helper has a separate type-switch mismatch.
	// Assert the retained concrete variant and serialized value instead of assuming it works.
	legacyError, ok := mcp.GetError().AsAny().(*responses.McpToolCallErrorMcpProtocolErrorParam)
	if !ok || legacyError.Code != -32603 {
		t.Fatal("MCP request projection lost its underlying numeric variant")
	}
	if code := mcp.OfMcpCall.Error.GetCode(); code == nil || *code != -32603 {
		t.Fatal("concrete MCP request GetCode lost its numeric value")
	}
	image := responses.ResponseInputItemUnionParam{OfImageGenerationCall: &responses.ResponseInputItemImageGenerationCallParam{ID: "img_synthetic", Result: param.Null[string](), Status: "failed", Error: responses.ResponseInputItemImageGenerationCallErrorParam{Code: openai.String("synthetic_code"), Param: openai.String("synthetic_param"), Message: "synthetic", Type: "image_generation_error"}}}
	if image.GetError().GetCode() != nil {
		t.Fatal("image request code must not be exposed as numeric")
	}
	for name, item := range map[string]responses.ResponseInputItemUnionParam{"mcp": mcp, "image": image} {
		t.Run(name, func(t *testing.T) {
			client := imageErrorClient(t, "POST", "/v1/responses", `{"id":"resp_synthetic","output":[]}`, func(r *http.Request) {
				var body struct {
					Input []struct {
						Error struct {
							Code  json.RawMessage `json:"code"`
							Param string          `json:"param"`
						} `json:"error"`
					} `json:"input"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if len(body.Input) != 1 {
					t.Fatal("expected one input item")
				}
				if name == "mcp" && string(body.Input[0].Error.Code) != "-32603" {
					t.Fatal("numeric request code changed")
				}
				if name == "image" && (string(body.Input[0].Error.Code) != `"synthetic_code"` || body.Input[0].Error.Param != "synthetic_param") {
					t.Fatal("image request fields changed")
				}
			})
			_, err := client.Responses.New(context.Background(), responses.ResponseNewParams{Model: "gpt-4o", Input: responses.ResponseNewParamsInputUnion{OfInputItemList: responses.ResponseInputParam{item}}})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
