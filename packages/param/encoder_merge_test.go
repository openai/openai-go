package param_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"math"
	"slices"
	"testing"

	shimjson "github.com/openai/openai-go/v3/internal/encoding/json"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/tidwall/sjson"
)

// Keep the released algorithm as a comparison for the encoding contract:
// bytes, field position, SJSON values, omission paths and error order.
func marshalExtrasReleased(underlying any, extras map[string]any) ([]byte, error) {
	encoded, err := shimjson.Marshal(underlying)
	if err != nil {
		return nil, err
	}
	for _, key := range slices.Sorted(maps.Keys(extras)) {
		value := extras[key]
		if value == param.Omit {
			if updated, deleteErr := sjson.DeleteBytes(encoded, key); deleteErr == nil {
				encoded = updated
			}
			continue
		}
		encoded, err = sjson.SetBytes(encoded, param.EscapeSJSONKey(key), value)
		if err != nil {
			return nil, err
		}
	}
	return encoded, nil
}

func TestExtraFieldsReleasedEncoding(t *testing.T) {
	tests := []struct {
		name       string
		underlying any
		extras     map[string]any
	}{
		{
			name:       "field position, null, replaced and absent omit",
			underlying: json.RawMessage(`{"z":1,"nested":{"one":true},"b":2,"a":"old","role":"user"}`),
			extras:     map[string]any{"0": false, "a": "new", "b": nil, "middle": json.RawMessage(`[1, null, {"x":-2.00e+2}]`), "role": param.Omit, "absent": param.Omit, "zz": true},
		},
		{
			name:       "literal punctuation and number keys",
			underlying: json.RawMessage(`{"keep":1,"a.b":"before","0":0,"q?":"old","back\\\\slash":false}`),
			extras:     map[string]any{".": 1, "|": 2, "#": 3, "@": 4, "*": 5, ":": 6, "?": 7, "\\": 8, "a.b": "after", "q?": nil, "0": "zero", "-1": true, "9": json.RawMessage(`{"t":null}`), "foo.bar|#@*:?'\\": "literal"},
		},
		{
			name:       "explicit path omit removes nested path, not literal field",
			underlying: json.RawMessage(`{"a":{"b":1,"keep":2},"a.b":"literal","keep":true}`),
			extras:     map[string]any{"a": map[string]any{"b": "changed", "keep": 42}, "a.b": param.Omit, "a.c": 3},
		},
		{
			name:       "explicit empty omit follows ignored legacy error",
			underlying: json.RawMessage(`{"":"keep","z":true}`),
			extras:     map[string]any{"": param.Omit, "new": false},
		},
		{
			name:       "SJSON bytes, number and string representation",
			underlying: json.RawMessage(`{"html\\u003ctag\\u003e":"old","unicode":"原文"}`),
			extras:     map[string]any{"html<tag>": "<script>&\u2028\u2029", "bytes": []byte("not base64 <&>"), "integer": int64(-123), "float": 1.0e30, "raw": json.RawMessage(`{ "n":1.200e-02,"a":"<&>" }`), "unicode": "声"},
		},
		{
			name:       "duplicate underlying key changes first member only",
			underlying: json.RawMessage(`{"a":1,"a":2,"b":3,"b":4}`),
			extras:     map[string]any{"a": false, "b": param.Omit, "c": "tail"},
		},
		{
			name:       "raw scalar underlying",
			underlying: json.RawMessage(`42`),
			extras:     map[string]any{"name": "scalar", "role": param.Omit},
		},
		{
			name:       "raw array underlying",
			underlying: json.RawMessage(`[1,2,3]`),
			extras:     map[string]any{"0": "first", "2": param.Omit},
		},
		{
			name:       "unsupported raw value",
			underlying: json.RawMessage(`{"keep":"base"}`),
			extras:     map[string]any{"first": true, "raw": json.RawMessage(`{bad}`), "z": false},
		},
		{
			name:       "SJSON nonfinite value remains an invalid JSON error at enclosing marshal",
			underlying: json.RawMessage(`{"a":0}`),
			extras:     map[string]any{"a": math.NaN(), "b": math.Inf(1)},
		},
		{
			name:       "invalid empty set path",
			underlying: json.RawMessage(`{"a":1}`),
			extras:     map[string]any{"": "error", "b": 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want, wantErr := marshalExtrasReleased(tt.underlying, tt.extras)
			got, gotErr := param.MarshalWithExtras(Struct{}, tt.underlying, tt.extras)
			if string(got) != string(want) || fmt.Sprint(gotErr) != fmt.Sprint(wantErr) {
				t.Fatalf("released: %s / %v\ncurrent:  %s / %v", want, wantErr, got, gotErr)
			}
		})
	}
}

type recordedMarshaler struct {
	name  string
	calls *[]string
	fail  bool
}

func (v recordedMarshaler) MarshalJSON() ([]byte, error) {
	*v.calls = append(*v.calls, v.name)
	if v.fail {
		return nil, errors.New("deliberate marshal error " + v.name)
	}
	return json.Marshal(v.name)
}

func TestExtraFieldsMarshalOrderAndErrors(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(fmt.Sprintf("fail=%t", fail), func(t *testing.T) {
			var wantCalls, gotCalls []string
			makeInputs := func(calls *[]string) (any, map[string]any) {
				base := struct {
					Base recordedMarshaler `json:"base"`
					Z    string            `json:"z"`
				}{recordedMarshaler{"base", calls, false}, "underlying"}
				extras := map[string]any{
					"z": recordedMarshaler{"last", calls, false},
					"a": recordedMarshaler{"first", calls, false},
					"m": recordedMarshaler{"middle", calls, fail},
				}
				return base, extras
			}
			base, extras := makeInputs(&wantCalls)
			want, wantErr := marshalExtrasReleased(base, extras)
			base, extras = makeInputs(&gotCalls)
			got, gotErr := param.MarshalWithExtras(Struct{}, base, extras)
			if string(got) != string(want) || fmt.Sprint(gotErr) != fmt.Sprint(wantErr) || !slices.Equal(gotCalls, wantCalls) {
				t.Fatalf("released: %s / %v / %v\ncurrent:  %s / %v / %v", want, wantErr, wantCalls, got, gotErr, gotCalls)
			}
		})
	}
}

func BenchmarkExtraFields(b *testing.B) {
	for _, prefix := range []string{"field_", "future.path|#@*:?-"} {
		for _, count := range []int{512, 1024, 2048, 4096} {
			b.Run(fmt.Sprintf("%s/%d", prefix, count), func(b *testing.B) {
				value := Struct{A: "current", B: 42}
				extras := make(map[string]any, count)
				for i := 0; i < count; i++ {
					extras[fmt.Sprintf("%s%05d", prefix, i)] = json.RawMessage("1")
				}
				value.SetExtraFields(extras)
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					if _, err := json.Marshal(value); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
