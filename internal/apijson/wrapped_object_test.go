package apijson_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/realtime"
)

func TestWrappedRequestBody(t *testing.T) {
	for _, tc := range []struct {
		name   string
		params json.Unmarshaler
		body   string
	}{
		{"union", &openai.BetaAgentEnvironmentFileNewParams{}, `{"type":"file_id","file_id":"file_123","path":"/workspace/input"}`},
		{"struct", &realtime.CallAcceptParams{}, `{"type":"realtime","model":"gpt-realtime","instructions":"Reply briefly."}`},
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
