package requestconfig

// ValidHTTPHeaderName reports whether name contains only HTTP token characters.
// This is internal API and may change without notice.
func ValidHTTPHeaderName(name string) bool {
	if name == "" {
		return false
	}
	for index := 0; index < len(name); index++ {
		value := name[index]
		switch {
		case value >= 'A' && value <= 'Z', value >= 'a' && value <= 'z', value >= '0' && value <= '9':
		case value == '!', value == '#', value == '$', value == '%', value == '&', value == '\'',
			value == '*', value == '+', value == '-', value == '.', value == '^', value == '_',
			value == '`', value == '|', value == '~':
		default:
			return false
		}
	}
	return true
}

// ValidHTTPHeaderValue reports whether value contains no prohibited HTTP
// control characters. This is internal API and may change without notice.
func ValidHTTPHeaderValue(value string) bool {
	for index := 0; index < len(value); index++ {
		character := value[index]
		if character < ' ' && character != '\t' || character == 0x7f {
			return false
		}
	}
	return true
}
