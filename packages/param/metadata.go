package param

import "encoding/json"

type MetadataProvider interface {
	IsNull() bool
	RawResponse() json.RawMessage
	IsOverridden() (any, bool)
}
