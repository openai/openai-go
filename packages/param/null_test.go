package param_test

import (
	"encoding/json"
	"github.com/openai/openai-go/v3/packages/param"
	"testing"
)

type Nullables struct {
	Slice []int          `json:"slice,omitzero"`
	Map   map[string]int `json:"map,omitzero"`
	param.APIObject
}

func (n Nullables) MarshalJSON() ([]byte, error) {
	type shadow Nullables
	return param.MarshalObject(n, (*shadow)(&n))
}

func TestNullMarshal(t *testing.T) {
	bytes, err := json.Marshal(Nullables{})
	if err != nil {
		t.Fatalf("json error %v", err.Error())
	}
	if string(bytes) != `{}` {
		t.Fatalf("expected empty object, got %s", string(bytes))
	}

	obj := Nullables{
		Slice: param.NullSlice[[]int](),
		Map:   param.NullMap[map[string]int](),
	}
	bytes, err = json.Marshal(obj)

	if !param.IsNull(obj.Slice) {
		t.Fatal("failed null check")
	}
	if !param.IsNull(obj.Map) {
		t.Fatal("failed null check")
	}

	if err != nil {
		t.Fatalf("json error %v", err.Error())
	}
	exp := `{"slice":null,"map":null}`
	if string(bytes) != exp {
		t.Fatalf("expected %s, got %s", exp, string(bytes))
	}

	bytes, err = json.Marshal(Nullables{Slice: []int{}, Map: map[string]int{}})
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != `{"slice":[],"map":{}}` {
		t.Fatalf("empty collections: got %s", bytes)
	}
}

func TestIsNullConcreteAndInterface(t *testing.T) {
	type namedMap map[string]int
	type namedSlice []int

	checkNullState(t, "null map", param.NullMap[map[string]int](), true, false)
	checkNullState(t, "null slice", param.NullSlice[[]int](), true, false)
	checkNullState(t, "null named map", param.NullMap[namedMap](), true, false)
	checkNullState(t, "null named slice", param.NullSlice[namedSlice](), true, false)

	checkNullState(t, "nil map", map[string]int(nil), false, true)
	checkNullState(t, "nil slice", []int(nil), false, true)
	checkNullState(t, "nil named map", namedMap(nil), false, true)
	checkNullState(t, "nil named slice", namedSlice(nil), false, true)
	checkNullState(t, "empty map", map[string]int{}, false, false)
	checkNullState(t, "empty slice", []int{}, false, false)
	checkNullState(t, "empty named map", namedMap{}, false, false)
	checkNullState(t, "empty named slice", namedSlice{}, false, false)
	checkNullState(t, "populated map", map[string]int{"a": 1}, false, false)
	checkNullState(t, "populated slice", []int{1}, false, false)

	checkNullState(t, "null opt", param.Null[string](), true, false)
	checkNullState(t, "omitted opt", param.Opt[string]{}, false, true)
	checkNullState(t, "set opt", param.NewOpt("value"), false, false)
	checkNullState(t, "null struct", param.NullStruct[Nullables](), true, false)
	checkNullState(t, "omitted struct", Nullables{}, false, true)
	checkNullState(t, "scalar", 1, false, false)
}

// Keep the concrete T when comparing it with the same value stored in any.
func checkNullState[T any](t *testing.T, name string, v T, wantNull, wantOmitted bool) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		if got := param.IsNull(v); got != wantNull {
			t.Errorf("IsNull(concrete) = %t, want %t", got, wantNull)
		}
		if got := param.IsNull(any(v)); got != wantNull {
			t.Errorf("IsNull(any) = %t, want %t", got, wantNull)
		}
		if got := param.IsOmitted(v); got != wantOmitted {
			t.Errorf("IsOmitted = %t, want %t", got, wantOmitted)
		}
	})
}

func TestIsNullNilInterface(t *testing.T) {
	if param.IsNull[any](nil) {
		t.Fatal("nil interface is not explicitly null")
	}
}
