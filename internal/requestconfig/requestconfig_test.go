package requestconfig

import "testing"

func TestNormalizeOS(t *testing.T) {
	tests := map[string]string{
		"android": "Android",
		"darwin":  "MacOS",
		"freebsd": "FreeBSD",
		"ios":     "iOS",
		"linux":   "Linux",
		"openbsd": "OpenBSD",
		"solaris": "Other:solaris",
		"windows": "Windows",
	}

	for goos, expected := range tests {
		if actual := normalizeOS(goos); actual != expected {
			t.Errorf("normalizeOS(%q) = %q, want %q", goos, actual, expected)
		}
	}
}
