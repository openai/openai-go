package sentinel

import (
	"reflect"
	"runtime"
	"sync"
	"weak"

	"github.com/openai/openai-go/v3/internal/encoding/json/shims"
)

type cacheEntry struct {
	x    any
	ptr  uintptr
	kind reflect.Kind
}

var nullCache sync.Map // map[reflect.Type]cacheEntry

// decodedNullMaps identifies independently allocated maps without keeping
// them alive. The weak identity prevents a reused address from matching a map
// whose cleanup has not run yet.
var decodedNullMaps sync.Map // map[uintptr]weak.Pointer[byte]

type decodedNullMap struct {
	address  uintptr
	identity weak.Pointer[byte]
}

func NewNullSentinel[T any](mk func() T) T {
	return newNullSentinel(shims.TypeFor[T](), func() any { return mk() }).(T)
}

// NullValue returns a null sentinel for a slice or map type, or an invalid
// value for other types. Decoded values cannot mutate shared storage.
func NullValue(t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.Map:
		value := reflect.MakeMap(t)
		// The pointer identifies the map allocation; it is never dereferenced.
		pointer := (*byte)(value.UnsafePointer())
		entry := decodedNullMap{address: value.Pointer(), identity: weak.Make(pointer)}
		decodedNullMaps.Store(entry.address, entry.identity)
		runtime.AddCleanup(pointer, func(entry decodedNullMap) {
			decodedNullMaps.CompareAndDelete(entry.address, entry.identity)
		}, entry)
		return value
	case reflect.Slice:
		value := reflect.ValueOf(newNullSentinel(t, func() any {
			return reflect.MakeSlice(t, 0, 1).Interface()
		}))
		// Preserve the sentinel identity without exposing its backing array.
		return value.Slice3(0, 0, 0)
	default:
		return reflect.Value{}
	}
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
		if ok && v.Pointer() == null.(cacheEntry).ptr {
			return true
		}
		if v.Kind() == reflect.Map {
			identity, ok := decodedNullMaps.Load(v.Pointer())
			return ok && identity.(weak.Pointer[byte]).Value() == (*byte)(v.UnsafePointer()) && v.Len() == 0
		}
	}
	return false
}

func IsNull[T any](v T) bool {
	return IsValueNull(reflect.ValueOf(v))
}
