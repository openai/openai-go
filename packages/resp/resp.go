package resp

type Field struct {
	// This implementation has more complexity than necessary, but it
	// makes this tradeoff in favor of better performance and lower footprint
	// in the happy path.
	//
	// Since this field appears so frequently, it's worth the complexity
	// to save memory.
	f *field
}

type field struct {
	status
	raw string
}

const (
	valid = iota
	invalid
)

type status int8

var fnull = field{raw: "null"}

// Returns true if the field is explicitly `null` _or_ if it is not present at all (ie, missing).
// To check if the field's key is present in the JSON with an explicit null value,
// you must check `f.IsNull() && !f.IsMissing()`.
func (j Field) IsNull() bool    { return j.f == &fnull }
func (j Field) IsMissing() bool { return j.f == nil || j.f == &fnull }
func (j Field) IsInvalid() bool { return j.f != nil && j.f.status == invalid }
func (j Field) Raw() string {
	if j.f == nil {
		return ""
	}
	return j.f.raw
}

func NewValidField(raw string) Field {
	if raw == "null" {
		return NewNullField()
	}
	return Field{f: &field{raw: raw, status: valid}}
}

func NewNullField() Field {
	return Field{f: &fnull}
}

func NewInvalidField(raw string) Field {
	return Field{f: &field{status: invalid, raw: raw}}
}
