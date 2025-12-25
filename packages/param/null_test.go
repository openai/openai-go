package param_test

import (
	"encoding/json"
	"github.com/Nordlys-Labs/openai-go/v3/packages/param"
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
}
