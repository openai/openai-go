package sentinel

import (
	"unsafe"
)

var isEmptyStringHackReliableOnThisTarget = false
var emptyByte byte = '\x00'

// Confirm the unsafe semantics of the EmptyString override are reliable
// on this target.
func init() {
	var zero string
	literal := ""
	empty := unsafe.String(&emptyByte, 0)

	isEmptyStringHackReliableOnThisTarget = unsafe.StringData(empty) == &emptyByte &&
		unsafe.StringData(literal) != &emptyByte &&
		unsafe.StringData(zero) != &emptyByte
}

// EmptyString returns a special empty string that can bypass the omitzero json tag,
// otherwise it can be used like a normal empty string.
//
// EmptyString will panic if this target doesn't support this feature.
// Before using this function write a quick test to confirm that it works on your target.
func EmptyString() string {
	if isEmptyStringHackReliableOnThisTarget {
		return unsafe.String(&emptyByte, 0)
	} else {
		panic("unsafe.StringData is not available on this target, use ForceEmptyString instead")
	}
}

func IsEmptyString(s string) bool {
	return isEmptyStringHackReliableOnThisTarget && len(s) == 0 && unsafe.StringData(s) == &emptyByte
}
