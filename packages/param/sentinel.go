package param

import (
	"github.com/openai/openai-go/internal/encoding/json/sentinel"
)

// NullPtr returns a pointer to the zero value of the type T.
// When passed to a custom MarshalJSON method, it will be marshaled as `null`.
//
// It is unspecified behavior to mutate the value pointed to by the returned pointer.
func NullPtr[T any]() *T {
	return sentinel.NullPtr[T]()
}

// IsNullPtr returns true if the pointer was created by [NullPtr].
func IsNullPtr[T any](ptr *T) bool {
	return sentinel.IsNullPtr(ptr)
}

// NullSlice returns a non-nil slice with a length of 0.
// When passed to a custom MarshalJSON method, it will be marshaled as `null`.
//
// It is undefined behavior to mutate the slice returned by [NullSlice].
func NullSlice[T any]() []T {
	return sentinel.NullSlice[T]()
}

// IsNullSlice returns true if the slice was created by [NullSlice].
func IsNullSlice[T any](slice []T) bool {
	return sentinel.IsNullSlice(slice)
}

// EmptyString returns a special empty string that can bypass the omitzero json tag,
// otherwise it can be used like a normal empty string.
//
// EmptyString will panic if this target doesn't support this feature.
// Before using this function write a quick test to confirm that it works on your target.
// If not, use [RemappedEmptyString] instead.
func EmptyString() string {
	return sentinel.EmptyString()
}

// IsEmptyString returns true if the string was created by [EmptyString] or [RemappedEmptyString].
func IsEmptyString(s string) bool {
	return sentinel.IsEmptyString(s)
}
