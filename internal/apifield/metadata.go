package apifield

import "encoding/json"

type ExplicitNull struct{}
type NeverOmitted struct{}
type CustomValue struct{ Override any }
type ResponseData json.RawMessage
type ExtraFields map[string]any

func (ExplicitNull) IsNull() bool   { return true }
func (NeverOmitted) IsNull() bool   { return false }
func (v CustomValue) IsNull() bool  { return v.Override == nil }
func (r ResponseData) IsNull() bool { return string(r) == `null` }
func (ExtraFields) IsNull() bool    { return false }

func (ExplicitNull) RawResponse() json.RawMessage   { return nil }
func (NeverOmitted) RawResponse() json.RawMessage   { return nil }
func (r ResponseData) RawResponse() json.RawMessage { return json.RawMessage(r) }
func (CustomValue) RawResponse() json.RawMessage    { return nil }
func (ExtraFields) RawResponse() json.RawMessage    { return nil }

func (ExplicitNull) IsOverridden() (any, bool)  { return nil, false }
func (NeverOmitted) IsOverridden() (any, bool)  { return nil, false }
func (v CustomValue) IsOverridden() (any, bool) { return v.Override, true }
func (ResponseData) IsOverridden() (any, bool)  { return nil, false }
func (ExtraFields) IsOverridden() (any, bool)   { return nil, false }
