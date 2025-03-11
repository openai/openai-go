// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package responses

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/apiquery"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/pagination"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/packages/resp"
	"github.com/openai/openai-go/shared/constant"
)

// InputItemService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewInputItemService] method instead.
type InputItemService struct {
	Options []option.RequestOption
}

// NewInputItemService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewInputItemService(opts ...option.RequestOption) (r InputItemService) {
	r = InputItemService{}
	r.Options = opts
	return
}

// Returns a list of input items for a given response.
func (r *InputItemService) List(ctx context.Context, responseID string, query InputItemListParams, opts ...option.RequestOption) (res *pagination.CursorPage[ResponseItemListDataUnion], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if responseID == "" {
		err = errors.New("missing required response_id parameter")
		return
	}
	path := fmt.Sprintf("responses/%s/input_items", responseID)
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// Returns a list of input items for a given response.
func (r *InputItemService) ListAutoPaging(ctx context.Context, responseID string, query InputItemListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[ResponseItemListDataUnion] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, responseID, query, opts...))
}

// A list of Response items.
type ResponseItemList struct {
	// A list of items used to generate this response.
	Data []ResponseItemListDataUnion `json:"data,required"`
	// The ID of the first item in the list.
	FirstID string `json:"first_id,required"`
	// Whether there are more items available.
	HasMore bool `json:"has_more,required"`
	// The ID of the last item in the list.
	LastID string `json:"last_id,required"`
	// The type of object returned, must be `list`.
	Object constant.List `json:"object,required"`
	// Metadata and presence of fields
	JSON struct {
		Data        resp.Field
		FirstID     resp.Field
		HasMore     resp.Field
		LastID      resp.Field
		Object      resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseItemList) RawJSON() string { return r.JSON.raw }
func (r *ResponseItemList) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ResponseItemListDataUnion contains all possible properties and values from
// [ResponseItemListDataMessage, ResponseOutputMessage, ResponseFileSearchToolCall,
// ResponseComputerToolCall, ResponseItemListDataComputerCallOutput,
// ResponseFunctionWebSearch, ResponseFunctionToolCall,
// ResponseItemListDataFunctionCallOutput].
//
// Use ResponseItemListDataUnion.AsAny() to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ResponseItemListDataUnion struct {
	ID string `json:"id"`
	// This field is a union of
	// [ResponseInputMessageContentList,[]ResponseOutputMessageContentUnion]
	Content ResponseItemListDataUnionContent `json:"content"`
	Role    string                           `json:"role"`
	Status  string                           `json:"status"`
	// Any of "message", "message", "file_search_call", "computer_call",
	// "computer_call_output", "web_search_call", "function_call",
	// "function_call_output".
	Type                string                                       `json:"type"`
	Queries             []string                                     `json:"queries"`
	Results             []ResponseFileSearchToolCallResult           `json:"results"`
	Action              ResponseComputerToolCallActionUnion          `json:"action"`
	CallID              string                                       `json:"call_id"`
	PendingSafetyChecks []ResponseComputerToolCallPendingSafetyCheck `json:"pending_safety_checks"`
	// This field is a union of [ResponseItemListDataComputerCallOutputOutput,string]
	Output                   ResponseItemListDataUnionOutput                                 `json:"output"`
	AcknowledgedSafetyChecks []ResponseItemListDataComputerCallOutputAcknowledgedSafetyCheck `json:"acknowledged_safety_checks"`
	Arguments                string                                                          `json:"arguments"`
	Name                     string                                                          `json:"name"`
	JSON                     struct {
		ID                       resp.Field
		Content                  resp.Field
		Role                     resp.Field
		Status                   resp.Field
		Type                     resp.Field
		Queries                  resp.Field
		Results                  resp.Field
		Action                   resp.Field
		CallID                   resp.Field
		PendingSafetyChecks      resp.Field
		Output                   resp.Field
		AcknowledgedSafetyChecks resp.Field
		Arguments                resp.Field
		Name                     resp.Field
		raw                      string
	} `json:"-"`
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ResponseItemListDataUnion.AsAny().(type) {
//	case ResponseItemListDataMessage:
//	case ResponseOutputMessage:
//	case ResponseFileSearchToolCall:
//	case ResponseComputerToolCall:
//	case ResponseItemListDataComputerCallOutput:
//	case ResponseFunctionWebSearch:
//	case ResponseFunctionToolCall:
//	case ResponseItemListDataFunctionCallOutput:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ResponseItemListDataUnion) AsAny() any {
	switch u.Type {
	case "message":
		return u.AsOutputMessage()
	case "file_search_call":
		return u.AsFileSearchCall()
	case "computer_call":
		return u.AsComputerCall()
	case "computer_call_output":
		return u.AsComputerCallOutput()
	case "web_search_call":
		return u.AsWebSearchCall()
	case "function_call":
		return u.AsFunctionCall()
	case "function_call_output":
		return u.AsFunctionCallOutput()
	}
	return nil
}

func (u ResponseItemListDataUnion) AsMessage() (v ResponseItemListDataMessage) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponseItemListDataUnion) AsOutputMessage() (v ResponseOutputMessage) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponseItemListDataUnion) AsFileSearchCall() (v ResponseFileSearchToolCall) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponseItemListDataUnion) AsComputerCall() (v ResponseComputerToolCall) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponseItemListDataUnion) AsComputerCallOutput() (v ResponseItemListDataComputerCallOutput) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponseItemListDataUnion) AsWebSearchCall() (v ResponseFunctionWebSearch) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponseItemListDataUnion) AsFunctionCall() (v ResponseFunctionToolCall) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponseItemListDataUnion) AsFunctionCallOutput() (v ResponseItemListDataFunctionCallOutput) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ResponseItemListDataUnion) RawJSON() string { return u.JSON.raw }

func (r *ResponseItemListDataUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ResponseItemListDataUnionContent is an implicit subunion of
// [ResponseItemListDataUnion]. ResponseItemListDataUnionContent provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ResponseItemListDataUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfInputItemContentList OfResponseOutputMessageContent]
type ResponseItemListDataUnionContent struct {
	OfInputItemContentList         ResponseInputMessageContentList     `json:",inline"`
	OfResponseOutputMessageContent []ResponseOutputMessageContentUnion `json:",inline"`
	JSON                           struct {
		OfInputItemContentList         resp.Field
		OfResponseOutputMessageContent resp.Field
		raw                            string
	} `json:"-"`
}

func (r *ResponseItemListDataUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ResponseItemListDataUnionOutput is an implicit subunion of
// [ResponseItemListDataUnion]. ResponseItemListDataUnionOutput provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ResponseItemListDataUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString]
type ResponseItemListDataUnionOutput struct {
	OfString string                      `json:",inline"`
	Type     constant.ComputerScreenshot `json:"type"`
	FileID   string                      `json:"file_id"`
	ImageURL string                      `json:"image_url"`
	JSON     struct {
		OfString resp.Field
		Type     resp.Field
		FileID   resp.Field
		ImageURL resp.Field
		raw      string
	} `json:"-"`
}

func (r *ResponseItemListDataUnionOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponseItemListDataMessage struct {
	// The unique ID of the message input.
	ID string `json:"id,required"`
	// A list of one or many input items to the model, containing different content
	// types.
	Content ResponseInputMessageContentList `json:"content,required"`
	// The role of the message input. One of `user`, `system`, or `developer`.
	//
	// Any of "user", "system", "developer".
	Role string `json:"role,required"`
	// The status of item. One of `in_progress`, `completed`, or `incomplete`.
	// Populated when items are returned via API.
	//
	// Any of "in_progress", "completed", "incomplete".
	Status string `json:"status"`
	// The type of the message input. Always set to `message`.
	//
	// Any of "message".
	Type string `json:"type"`
	// Metadata and presence of fields
	JSON struct {
		ID          resp.Field
		Content     resp.Field
		Role        resp.Field
		Status      resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseItemListDataMessage) RawJSON() string { return r.JSON.raw }
func (r *ResponseItemListDataMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponseItemListDataComputerCallOutput struct {
	// The unique ID of the computer call tool output.
	ID string `json:"id,required"`
	// The ID of the computer tool call that produced the output.
	CallID string `json:"call_id,required"`
	// A computer screenshot image used with the computer use tool.
	Output ResponseItemListDataComputerCallOutputOutput `json:"output,required"`
	// The type of the computer tool call output. Always `computer_call_output`.
	Type constant.ComputerCallOutput `json:"type,required"`
	// The safety checks reported by the API that have been acknowledged by the
	// developer.
	AcknowledgedSafetyChecks []ResponseItemListDataComputerCallOutputAcknowledgedSafetyCheck `json:"acknowledged_safety_checks"`
	// The status of the message input. One of `in_progress`, `completed`, or
	// `incomplete`. Populated when input items are returned via API.
	//
	// Any of "in_progress", "completed", "incomplete".
	Status string `json:"status"`
	// Metadata and presence of fields
	JSON struct {
		ID                       resp.Field
		CallID                   resp.Field
		Output                   resp.Field
		Type                     resp.Field
		AcknowledgedSafetyChecks resp.Field
		Status                   resp.Field
		ExtraFields              map[string]resp.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseItemListDataComputerCallOutput) RawJSON() string { return r.JSON.raw }
func (r *ResponseItemListDataComputerCallOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A computer screenshot image used with the computer use tool.
type ResponseItemListDataComputerCallOutputOutput struct {
	// Specifies the event type. For a computer screenshot, this property is always set
	// to `computer_screenshot`.
	Type constant.ComputerScreenshot `json:"type,required"`
	// The identifier of an uploaded file that contains the screenshot.
	FileID string `json:"file_id"`
	// The URL of the screenshot image.
	ImageURL string `json:"image_url"`
	// Metadata and presence of fields
	JSON struct {
		Type        resp.Field
		FileID      resp.Field
		ImageURL    resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseItemListDataComputerCallOutputOutput) RawJSON() string { return r.JSON.raw }
func (r *ResponseItemListDataComputerCallOutputOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A pending safety check for the computer call.
type ResponseItemListDataComputerCallOutputAcknowledgedSafetyCheck struct {
	// The ID of the pending safety check.
	ID string `json:"id,required"`
	// The type of the pending safety check.
	Code string `json:"code,required"`
	// Details about the pending safety check.
	Message string `json:"message,required"`
	// Metadata and presence of fields
	JSON struct {
		ID          resp.Field
		Code        resp.Field
		Message     resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseItemListDataComputerCallOutputAcknowledgedSafetyCheck) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponseItemListDataComputerCallOutputAcknowledgedSafetyCheck) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponseItemListDataFunctionCallOutput struct {
	// The unique ID of the function call tool output.
	ID string `json:"id,required"`
	// The unique ID of the function tool call generated by the model.
	CallID string `json:"call_id,required"`
	// A JSON string of the output of the function tool call.
	Output string `json:"output,required"`
	// The type of the function tool call output. Always `function_call_output`.
	Type constant.FunctionCallOutput `json:"type,required"`
	// The status of the item. One of `in_progress`, `completed`, or `incomplete`.
	// Populated when items are returned via API.
	//
	// Any of "in_progress", "completed", "incomplete".
	Status string `json:"status"`
	// Metadata and presence of fields
	JSON struct {
		ID          resp.Field
		CallID      resp.Field
		Output      resp.Field
		Type        resp.Field
		Status      resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseItemListDataFunctionCallOutput) RawJSON() string { return r.JSON.raw }
func (r *ResponseItemListDataFunctionCallOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type InputItemListParams struct {
	// An item ID to list items after, used in pagination.
	After param.Opt[string] `query:"after,omitzero"`
	// An item ID to list items before, used in pagination.
	Before param.Opt[string] `query:"before,omitzero"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Opt[int64] `query:"limit,omitzero"`
	// The order to return the input items in. Default is `asc`.
	//
	// - `asc`: Return the input items in ascending order.
	// - `desc`: Return the input items in descending order.
	//
	// Any of "asc", "desc".
	Order InputItemListParamsOrder `query:"order,omitzero"`
	paramObj
}

// IsPresent returns false if the field is omitted or `null`.
func (f InputItemListParams) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }

// URLQuery serializes [InputItemListParams]'s query parameters as `url.Values`.
func (r InputItemListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The order to return the input items in. Default is `asc`.
//
// - `asc`: Return the input items in ascending order.
// - `desc`: Return the input items in descending order.
type InputItemListParamsOrder string

const (
	InputItemListParamsOrderAsc  InputItemListParamsOrder = "asc"
	InputItemListParamsOrderDesc InputItemListParamsOrder = "desc"
)
