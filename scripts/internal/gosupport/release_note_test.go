package gosupport_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3/scripts/internal/gosupport"
)

func TestCheckReleaseNoteAcceptsApprovedTextInReleaseSection(t *testing.T) {
	root, note := makeReleaseRepository(t)
	writeReleaseFile(t, root, "CHANGELOG.md", fmt.Sprintf(`# Changelog

## 3.45.0 (2026-08-01)

%s

### Features

* require Go 1.25

## 3.44.0 (2026-07-17)
`, note))

	if findings := gosupport.CheckReleaseNote(root); len(findings) != 0 {
		t.Fatalf("unexpected findings: %v", findings)
	}
}

func TestCheckReleaseNoteRejectsTextInOlderSection(t *testing.T) {
	root, note := makeReleaseRepository(t)
	writeReleaseFile(t, root, "CHANGELOG.md", fmt.Sprintf(`# Changelog

## 3.45.0 (2026-08-01)

### Features

* require Go 1.25

## 3.44.0 (2026-07-17)

%s
`, note))

	findings := gosupport.CheckReleaseNote(root)
	if !containsReleaseFinding(findings, "does not contain the approved") {
		t.Fatalf("expected missing-note finding, got: %v", findings)
	}
}

func makeReleaseRepository(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	note := "openai-go v3.45.0 requires Go 1.25 or newer. " +
		"v3.44.0 is the final compatible release."
	writeReleaseFile(t, root, ".github/go-support-policy.json", fmt.Sprintf(`{
		"minimum":"1.25",
		"grace_until":null,
		"reason":"test",
		"approved_by":"SDK team",
		"release_note":{
			"sdk_version":"v3.45.0",
			"last_compatible_sdk_version":"v3.44.0",
			"text":%q
		}
	}`, note))
	return root, note
}

func writeReleaseFile(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, relative)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func containsReleaseFinding(findings []string, needle string) bool {
	for _, finding := range findings {
		if strings.Contains(finding, needle) {
			return true
		}
	}
	return false
}
