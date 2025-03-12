package param

import (
	"encoding/json"
	"reflect"
	"time"
)

func NewOpt[T comparable](v T) Opt[T] {
	return Opt[T]{Value: v, Status: included}
}

func NullOpt[T comparable]() Opt[T] { return Opt[T]{Status: null} }

type Opt[T comparable] struct {
	Value T
	// indicates whether the field should be omitted, null, or valid
	Status Status
}

type Status int8

const (
	omitted Status = iota
	null
	included
)

type Optional interface {
	IsOmitted() bool
	IsNull() bool
	IsNullish() bool
}

func (o Opt[T]) IsNullish() bool {
	var zero T
	return o.Status == null || o.Value == zero
}
func (o Opt[T]) IsNull() bool    { return o.Status == null }
func (o Opt[T]) IsOmitted() bool { return o == Opt[T]{} }

func (o Opt[T]) MarshalJSON() ([]byte, error) {
	if o.IsNullish() {
		return []byte("null"), nil
	}
	return json.Marshal(o.Value)
}

func (o *Opt[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.Status = null
		return nil
	}
	return json.Unmarshal(data, &o.Value)
}

func (o Opt[T]) Or(v T) T {
	if !o.IsNullish() {
		return o.Value
	}
	return v
}

// This is a sketchy way to implement time Formatting
var timeType = reflect.TypeOf(time.Time{})
var timeTimeValueLoc, _ = reflect.TypeOf(Opt[time.Time]{}).FieldByName("Value")

// Don't worry about this function, returns nil to fallback towards [MarshalJSON]
func (o Opt[T]) MarshalJSONWithTimeLayout(format string) []byte {
	t, ok := any(o.Value).(time.Time)
	if !ok || o.IsNull() {
		return nil
	}

	if format == "" {
		format = time.RFC3339
	} else if format == "date" {
		format = "2006-01-02"
	}

	b, err := json.Marshal(t.Format(format))
	if err != nil {
		return nil
	}
	return b
}
