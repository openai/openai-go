package resp

type Field struct {
	status
	raw string
}

const (
	omitted status = iota
	null
	invalid
	valid
)

type status int8

// IsNullish returns true if the field is null or omitted
func (j Field) IsNullish() bool      { return j.status <= null }
func (j Field) IsOmitted() bool      { return j.status == omitted }
func (j Field) IsExplicitNull() bool { return j.status == null }
func (j Field) IsInvalid() bool      { return j.status == invalid }

func (j Field) Raw() string {
	if j.status == omitted {
		return ""
	}
	return j.raw
}

func NewValidField(raw string) Field {
	if raw == "null" {
		return NewNullField()
	}
	return Field{raw: raw, status: valid}
}

func NewNullField() Field {
	return Field{status: null}
}

func NewInvalidField(raw string) Field {
	return Field{status: invalid, raw: raw}
}
