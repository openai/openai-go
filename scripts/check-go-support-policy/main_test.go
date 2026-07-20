package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCheckRepositoryAcceptsCurrentSupportWindow(t *testing.T) {
	root := makeRepository(t, "1.25", "null")
	server := releaseServer(t)
	defer server.Close()

	result, err := checkRepository(root, server.URL, time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC), server.Client())
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 0 {
		t.Fatalf("unexpected findings: %v", result.Findings)
	}
}

func TestCheckRepositoryRejectsExpiredGrace(t *testing.T) {
	root := makeRepository(t, "1.24", `"2026-06-30"`)
	server := releaseServer(t)
	defer server.Close()

	result, err := checkRepository(root, server.URL, time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC), server.Client())
	if err != nil {
		t.Fatal(err)
	}
	if !containsFinding(result.Findings, "expired on 2026-06-30") {
		t.Fatalf("expected expired grace finding, got: %v", result.Findings)
	}
}

func TestSupportedLinesIgnorePrereleasesAndPatchVersions(t *testing.T) {
	lines, err := supportedLines([]goRelease{
		{Version: "go1.27rc1"},
		{Version: "go1.26.5"},
		{Version: "go1.26.4"},
		{Version: "go1.25.12"},
		{Version: "go1.24.18"},
	})
	if err != nil {
		t.Fatal(err)
	}
	got := fmt.Sprint(lines)
	if got != "[1.26 1.25 1.24]" {
		t.Fatalf("supportedLines() = %s", got)
	}
}

func releaseServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		fmt.Fprint(response, `[
			{"version":"go1.26.5"},
			{"version":"go1.25.12"},
			{"version":"go1.24.18"},
			{"version":"go1.27rc1"}
		]`)
	}))
}

func makeRepository(t *testing.T, minimum, graceUntil string) string {
	t.Helper()
	root := t.TempDir()

	files := map[string]string{
		".github/go-support-policy.json": fmt.Sprintf(
			`{"minimum":%q,"grace_until":%s,"reason":"test policy","approved_by":"SDK team"}`,
			minimum,
			graceUntil,
		),
		"go.mod":                            "module example.com/root\n\ngo " + minimum + ".0\n",
		"examples/go.mod":                   "module example.com/examples\n\ngo " + minimum + ".0\n",
		"internal/testdata/consumer/go.mod": "module example.com/consumer\n\ngo " + minimum + ".0\n",
		"README.md":                         "Requires Go " + minimum + ".\n",
		"CONTRIBUTING.md":                   "Install Go " + minimum + ".\n",
		"GO_VERSION_POLICY.md":              "The minimum is Go " + minimum + ".\n",
		".github/workflows/ci.yml":          "1.25.x\n1.26.x\nGOTOOLCHAIN: local\n./scripts/check-go-mod\ngovulncheck\n",
		"scripts/check-go-mod":              "go mod tidy -diff\n",
	}

	for relative, content := range files {
		path := filepath.Join(root, relative)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func containsFinding(findings []string, needle string) bool {
	for _, finding := range findings {
		if strings.Contains(finding, needle) {
			return true
		}
	}
	return false
}
