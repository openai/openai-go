package sentinel_test

import (
	"github.com/Nordlys-Labs/openai-go/v3/internal/encoding/json/sentinel"
	"github.com/Nordlys-Labs/openai-go/v3/packages/param"
	"reflect"
	"slices"
	"testing"
)

type Pair struct {
	got  bool
	want bool
}

func TestNullSlice(t *testing.T) {
	var nilSlice []int = nil
	var nonNilSlice []int = []int{1, 2, 3}
	var nullSlice []int = param.NullSlice[[]int]()

	cases := map[string]Pair{
		"nilSlice":            {sentinel.IsNull(nilSlice), false},
		"nullSlice":           {sentinel.IsNull(nullSlice), true},
		"newNullSlice":        {sentinel.IsNull(param.NullSlice[[]int]()), true},
		"lenNullSlice":        {len(nullSlice) == 0, true},
		"nilSliceValue":       {sentinel.IsValueNull(reflect.ValueOf(nilSlice)), false},
		"nullSliceValue":      {sentinel.IsValueNull(reflect.ValueOf(nullSlice)), true},
		"compareSlices":       {slices.Compare(nilSlice, nullSlice) == 0, true},
		"compareNonNilSlices": {slices.Compare(nonNilSlice, nullSlice) == 0, false},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got, want := c.got, c.want
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestNullMap(t *testing.T) {
	var nilMap map[string]int = nil
	var nonNilMap map[string]int = map[string]int{"a": 1, "b": 2}
	var nullMap map[string]int = param.NullMap[map[string]int]()

	cases := map[string]Pair{
		"nilMap":            {sentinel.IsNull(nilMap), false},
		"nullMap":           {sentinel.IsNull(nullMap), true},
		"newNullMap":        {sentinel.IsNull(param.NullMap[map[string]int]()), true},
		"lenNullMap":        {len(nullMap) == 0, true},
		"nilMapValue":       {sentinel.IsValueNull(reflect.ValueOf(nilMap)), false},
		"nullMapValue":      {sentinel.IsValueNull(reflect.ValueOf(nullMap)), true},
		"compareMaps":       {reflect.DeepEqual(nilMap, nullMap), false},
		"compareNonNilMaps": {reflect.DeepEqual(nonNilMap, nullMap), false},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			got, want := test.got, test.want
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestIsNullRepeated(t *testing.T) {
	// Test for slices
	nullSlice1 := param.NullSlice[[]int]()
	nullSlice2 := param.NullSlice[[]int]()
	if !sentinel.IsNull(nullSlice1) {
		t.Errorf("IsNull(nullSlice1) = false, want true")
	}
	if !sentinel.IsNull(nullSlice2) {
		t.Errorf("IsNull(nullSlice2) = false, want true")
	}
	if !sentinel.IsNull(nullSlice1) || !sentinel.IsNull(nullSlice2) {
		t.Errorf("IsNull should return true for all NullSlice instances")
	}

	// Test for maps
	nullMap1 := param.NullMap[map[string]int]()
	nullMap2 := param.NullMap[map[string]int]()
	if !sentinel.IsNull(nullMap1) {
		t.Errorf("IsNull(nullMap1) = false, want true")
	}
	if !sentinel.IsNull(nullMap2) {
		t.Errorf("IsNull(nullMap2) = false, want true")
	}
}
