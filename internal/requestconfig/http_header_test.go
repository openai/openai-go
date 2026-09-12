package requestconfig

import "testing"

func TestValidHTTPHeaderName(t *testing.T) {
	for _, test := range []struct {
		name  string
		value string
		valid bool
	}{
		{name: "letters and digits", value: "X-Request-123", valid: true},
		{name: "all token punctuation", value: "!#$%&'*+-.^_`|~", valid: true},
		{name: "empty"},
		{name: "colon", value: "X:Unsafe"},
		{name: "space", value: "X Unsafe"},
		{name: "newline", value: "X\nUnsafe"},
		{name: "non ASCII", value: "X-\xff"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := ValidHTTPHeaderName(test.value); got != test.valid {
				t.Errorf("ValidHTTPHeaderName(%q) = %v, want %v", test.value, got, test.valid)
			}
		})
	}
}

func TestValidHTTPHeaderValue(t *testing.T) {
	for _, test := range []struct {
		name  string
		value string
		valid bool
	}{
		{name: "empty", valid: true},
		{name: "printable", value: "Bearer synthetic-token", valid: true},
		{name: "horizontal tab", value: "first\tsecond", valid: true},
		{name: "extended bytes", value: "\xff", valid: true},
		{name: "null", value: "first\x00second"},
		{name: "carriage return", value: "first\rsecond"},
		{name: "newline", value: "first\nsecond"},
		{name: "delete", value: "first\x7fsecond"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := ValidHTTPHeaderValue(test.value); got != test.valid {
				t.Errorf("ValidHTTPHeaderValue(%q) = %v, want %v", test.value, got, test.valid)
			}
		})
	}
}
