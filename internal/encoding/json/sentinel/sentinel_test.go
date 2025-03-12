package sentinel_test

import (
	"github.com/openai/openai-go/internal/encoding/json/sentinel"
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
	var nullSlice []int = sentinel.NullSlice[int]()

	cases := map[string]Pair{
		"nilSlice":            {sentinel.IsNullSlice(nilSlice), false},
		"nullSlice":           {sentinel.IsNullSlice(nullSlice), true},
		"newNullSlice":        {sentinel.IsNullSlice(sentinel.NullSlice[int]()), true},
		"lenNullSlice":        {len(nullSlice) == 0, true},
		"nilSliceValue":       {sentinel.IsValueNullSlice(reflect.ValueOf(nilSlice)), false},
		"nullSliceValue":      {sentinel.IsValueNullSlice(reflect.ValueOf(nullSlice)), true},
		"compareSlices":       {slices.Compare(nilSlice, nullSlice) == 0, true},
		"compareNonNilSlices": {slices.Compare(nonNilSlice, nullSlice) == 0, false},
	}

	nilSlice = append(nullSlice, 12)
	cases["append_result"] = Pair{sentinel.IsNullSlice(nilSlice), false}
	cases["mutated_result"] = Pair{sentinel.IsNullSlice(nullSlice), true}
	cases["append_result_len"] = Pair{len(nilSlice) == 1, true}
	cases["append_null_slice_len"] = Pair{len(nullSlice) == 0, true}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got, want := c.got, c.want
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestNullPtr(t *testing.T) {
	var s *string = nil
	var i *int = nil
	var slice *[]int = nil

	var nullptrStr *string = sentinel.NullPtr[string]()
	var nullptrInt *int = sentinel.NullPtr[int]()
	var nullptrSlice *[]int = sentinel.NullPtr[[]int]()

	if *nullptrStr != "" {
		t.Errorf("Failed to safely deref")
	}
	if *nullptrInt != 0 {
		t.Errorf("Failed to safely deref")
	}
	if len(*nullptrSlice) != 0 {
		t.Errorf("Failed to safely deref")
	}

	cases := map[string]Pair{
		"nilStr":  {sentinel.IsNullPtr(s), false},
		"nullStr": {sentinel.IsNullPtr(nullptrStr), true},

		"nilInt":  {sentinel.IsNullPtr(i), false},
		"nullInt": {sentinel.IsNullPtr(nullptrInt), true},

		"nilSlice":  {sentinel.IsNullPtr(slice), false},
		"nullSlice": {sentinel.IsNullPtr(nullptrSlice), true},

		"nilValuePtr":  {sentinel.IsValueNullPtr(reflect.ValueOf(i)), false},
		"nullValuePtr": {sentinel.IsValueNullPtr(reflect.ValueOf(nullptrInt)), true},
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
