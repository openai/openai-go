package sentinel

import (
	"github.com/openai/openai-go/v3/internal/encoding/json/shims"
	"reflect"
	"sync"
)

type cacheEntry struct {
	x    any
	ptr  uintptr
	kind reflect.Kind
}

var nullCache sync.Map // map[reflect.Type]cacheEntry

func NewNullSentinel[T any](mk func() T) T {
	return newNullSentinel(shims.TypeFor[T](), func() any { return mk() }).(T)
}

// NullValue returns the null sentinel for a slice or map type, or an invalid
// value for other types.
func NullValue(t reflect.Type) reflect.Value {
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Map {
		return reflect.Value{}
	}
	return reflect.ValueOf(newNullSentinel(t, func() any {
		if t.Kind() == reflect.Map {
			return reflect.MakeMapWithSize(t, 1).Interface()
		}
		return reflect.MakeSlice(t, 0, 1).Interface()
	}))
}

func newNullSentinel(t reflect.Type, mk func() any) any {
	entry, loaded := nullCache.Load(t) // avoid premature allocation
	if !loaded {
		x := mk()
		ptr := reflect.ValueOf(x).Pointer()
		entry, _ = nullCache.LoadOrStore(t, cacheEntry{x, ptr, t.Kind()})
	}
	return entry.(cacheEntry).x
}

// for internal use only
func IsValueNull(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Map, reflect.Slice:
		null, ok := nullCache.Load(v.Type())
		return ok && v.Pointer() == null.(cacheEntry).ptr
	}
	return false
}

func IsNull[T any](v T) bool {
	t := reflect.TypeOf(v)
	if t == nil {
		return false
	}
	switch t.Kind() {
	case reflect.Map, reflect.Slice:
		null, ok := nullCache.Load(t)
		return ok && reflect.ValueOf(v).Pointer() == null.(cacheEntry).ptr
	}
	return false
}
