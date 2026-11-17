// Package sjson preserves the SDK's empty-path behavior across upstream updates.
package sjson

import (
	"encoding/json"
	"errors"

	upstream "github.com/tidwall/sjson"
)

// SetBytes retains the pre-1.2.7 empty-path error and value-encoding order.
func SetBytes(data []byte, path string, value any) ([]byte, error) {
	if path == "" {
		// These are exactly the upstream 1.2.5 fast-path types. Named types,
		// pointers and custom marshalers still pass through encoding/json first.
		switch value.(type) {
		case string, []byte, bool, int8, int16, int32, int64,
			uint8, uint16, uint32, uint64, float32, float64:
		default:
			if _, err := json.Marshal(value); err != nil {
				return nil, err
			}
		}
		return data, errors.New("path cannot be empty")
	}
	return upstream.SetBytes(data, path, value)
}

// SetRawBytes rejects empty paths as earlier SDK releases did.
func SetRawBytes(data []byte, path string, value []byte) ([]byte, error) {
	if path == "" {
		return data, errors.New("path cannot be empty")
	}
	return upstream.SetRawBytes(data, path, value)
}

// DeleteBytes preserves the error callers historically ignored for empty Omit paths.
func DeleteBytes(data []byte, path string) ([]byte, error) {
	if path == "" {
		return data, errors.New("path cannot be empty")
	}
	return upstream.DeleteBytes(data, path)
}
