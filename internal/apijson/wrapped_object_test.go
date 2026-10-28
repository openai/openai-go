package apijson_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"runtime"
	"testing"

	"github.com/openai/openai-go/v4"
	"github.com/openai/openai-go/v4/packages/param"
	"github.com/openai/openai-go/v4/realtime"
	"github.com/openai/openai-go/v4/responses"
	"github.com/tidwall/gjson"
)

func TestWrappedRequestBody(t *testing.T) {
	for _, tc := range []struct {
		name   string
		params json.Unmarshaler
		body   string
	}{
		{"union", &openai.BetaAgentEnvironmentFileNewParams{}, `{"type":"file_id","file_id":"file_123","path":"/workspace/input"}`},
		{"struct", &realtime.CallAcceptParams{}, `{"type":"realtime","model":"gpt-realtime","instructions":"Reply briefly."}`},
		{"explicit null", &realtime.CallAcceptParams{}, `{"type":"realtime","instructions":null}`},
		{"translation", &realtime.TranslationClientSecretNewParams{}, `{"session":{"model":"gpt-realtime-translate"},"expires_after":{"anchor":"created_at","seconds":30}}`},
		{"null union field", &realtime.CallAcceptParams{}, `{"type":"realtime","tracing":null}`},
		{"null interface field", &realtime.CallAcceptParams{}, `{"type":"realtime","tracing":{"metadata":null}}`},
		{"omitted union field", &realtime.CallAcceptParams{}, `{"type":"realtime"}`},
		{"null object field", &realtime.CallAcceptParams{}, `{"type":"realtime","audio":null}`},
		{"null collections", &realtime.CallAcceptParams{}, `{"type":"realtime","include":null,"prompt":{"id":"pmpt_123","variables":null}}`},
		{"null within union", &realtime.CallAcceptParams{}, `{"type":"realtime","truncation":{"type":"retention_ratio","retention_ratio":0.8,"token_limits":null}}`},
		{"nested translation null", &realtime.TranslationClientSecretNewParams{}, `{"session":{"model":"translation-model","audio":{"input":{"noise_reduction":null}}}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(tc.body), tc.params); err != nil {
				t.Fatal(err)
			}
			before, err := json.Marshal(tc.params)
			if err != nil {
				t.Fatal(err)
			}
			var want, got map[string]any
			if err := json.Unmarshal([]byte(tc.body), &want); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(before, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("round trip = %s, want %s", before, tc.body)
			}

			t.Run("null", func(t *testing.T) {
				if err := json.Unmarshal([]byte("null"), tc.params); err != nil {
					t.Fatal(err)
				}
				after, err := json.Marshal(tc.params)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(after, before) {
					t.Fatalf("null changed the body from %s to %s", before, after)
				}
			})
			t.Run("malformed JSON", func(t *testing.T) {
				err := json.Unmarshal([]byte(`{"type":`), tc.params)
				var syntaxErr *json.SyntaxError
				if !errors.As(err, &syntaxErr) {
					t.Fatalf("error = %T %v, want a syntax error", err, err)
				}
			})
			for _, root := range []string{`[]`, `true`, `42`, `"value"`} {
				t.Run(root, func(t *testing.T) {
					err := json.Unmarshal([]byte(root), tc.params)
					var typeErr *json.UnmarshalTypeError
					if !errors.As(err, &typeErr) || typeErr.Type != reflect.TypeOf(tc.params).Elem() {
						t.Fatalf("error = %T %v, want a type error for %T", err, err, tc.params)
					}
				})
			}
		})
	}

	t.Run("nil receiver", func(t *testing.T) {
		var params *realtime.CallAcceptParams
		if err := params.UnmarshalJSON([]byte(`{}`)); err == nil {
			t.Fatal("expected an error for a nil receiver")
		}
	})
}

func TestWrappedRequestBodyNullMapIsolation(t *testing.T) {
	const body = `{"type":"realtime","prompt":{"id":"pmpt_123","variables":null}}`
	var first, second realtime.CallAcceptParams
	for _, params := range []*realtime.CallAcceptParams{&first, &second} {
		if err := json.Unmarshal([]byte(body), params); err != nil {
			t.Fatal(err)
		}
		variables := params.RealtimeSessionCreateRequest.Prompt.Variables
		if !param.IsNull(variables) || param.IsOmitted(variables) {
			t.Fatal("decoded variables should be explicitly null")
		}
	}

	runtime.GC()
	if !param.IsNull(any(second.RealtimeSessionCreateRequest.Prompt.Variables)) {
		t.Fatal("decoded map should remain null after garbage collection")
	}
	firstVariables := first.RealtimeSessionCreateRequest.Prompt.Variables
	// Keep the test isolated if a regression restores the shared sentinel.
	t.Cleanup(func() { delete(firstVariables, "tenant") })
	firstVariables["tenant"] = responses.ResponsePromptVariableUnionParam{OfString: param.NewOpt("first")}
	if _, ok := second.RealtimeSessionCreateRequest.Prompt.Variables["tenant"]; ok {
		t.Fatal("independently decoded null maps share values")
	}
	if _, ok := param.NullMap[map[string]responses.ResponsePromptVariableUnionParam]()["tenant"]; ok {
		t.Fatal("decoded null map shares values with param.NullMap")
	}
	if param.IsNull(firstVariables) {
		t.Fatal("populated map should no longer be null")
	}

	for _, tc := range []struct {
		params realtime.CallAcceptParams
		want   string
	}{
		{first, `{"tenant":"first"}`},
		{second, "null"},
	} {
		got, err := json.Marshal(tc.params)
		if err != nil {
			t.Fatal(err)
		}
		if variables := gjson.GetBytes(got, "prompt.variables").Raw; variables != tc.want {
			t.Errorf("marshaled variables = %s, want %s", variables, tc.want)
		}
	}

	clear(firstVariables)
	got, err := json.Marshal(first)
	if err != nil {
		t.Fatal(err)
	}
	if variables := gjson.GetBytes(got, "prompt.variables").Raw; variables != "null" {
		t.Fatalf("cleared variables = %s, want null", variables)
	}
}

func TestWrappedRequestBodyNullSliceIsolation(t *testing.T) {
	const body = `{"type":"realtime","include":null}`
	var first, second realtime.CallAcceptParams
	for _, params := range []*realtime.CallAcceptParams{&first, &second} {
		if err := json.Unmarshal([]byte(body), params); err != nil {
			t.Fatal(err)
		}
		if !param.IsNull(params.RealtimeSessionCreateRequest.Include) {
			t.Fatal("decoded include should be explicitly null")
		}
	}
	first.RealtimeSessionCreateRequest.Include = append(first.RealtimeSessionCreateRequest.Include, "first")
	t.Cleanup(func() { first.RealtimeSessionCreateRequest.Include[0] = "" })
	second.RealtimeSessionCreateRequest.Include = append(second.RealtimeSessionCreateRequest.Include, "second")
	if got := first.RealtimeSessionCreateRequest.Include[0]; got != "first" {
		t.Fatalf("first include = %q, want first", got)
	}
}

func TestWrappedRequestBodyTracingNull(t *testing.T) {
	var params realtime.CallAcceptParams
	for _, body := range []string{
		`{"type":"realtime","tracing":"auto"}`,
		`{"type":"realtime","tracing":null}`,
		`{"type":"realtime","tracing":"auto"}`,
	} {
		if err := json.Unmarshal([]byte(body), &params); err != nil {
			t.Fatal(err)
		}
		wantNull := bytes.Contains([]byte(body), []byte("null"))
		if param.IsNull(params.RealtimeSessionCreateRequest.Tracing) != wantNull {
			t.Fatalf("tracing null state does not match %s", body)
		}
		got, err := json.Marshal(params)
		if err != nil {
			t.Fatal(err)
		}
		var wantFields, gotFields map[string]any
		if err := json.Unmarshal([]byte(body), &wantFields); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(got, &gotFields); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(gotFields, wantFields) {
			t.Fatalf("round trip = %s, want %s", got, body)
		}
	}
}
