package gosupport_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/scripts/internal/gosupport"
)

func TestCheckRepositoryAcceptsCurrentSupportWindow(t *testing.T) {
	root := makeRepository(t, "1.25", "null", []string{"1.25.x", "1.26.x"})
	server := releaseServer(t)
	defer server.Close()

	result, err := gosupport.CheckRepository(
		root,
		server.URL,
		time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC),
		server.Client(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 0 {
		t.Fatalf("unexpected findings: %v", result.Findings)
	}
}

func TestCheckRepositoryRejectsExpiredGrace(t *testing.T) {
	root := makeRepository(t, "1.24", `"2026-06-30"`, []string{"1.24.x", "1.25.x", "1.26.x"})
	server := releaseServer(t)
	defer server.Close()

	result, err := gosupport.CheckRepository(
		root,
		server.URL,
		time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC),
		server.Client(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !containsFinding(result.Findings, "expired on 2026-06-30") {
		t.Fatalf("expected expired grace finding, got: %v", result.Findings)
	}
}

func TestCheckRepositoryRejectsGraceLongerThanSixMonths(t *testing.T) {
	root := makeRepository(t, "1.24", `"2027-01-01"`, []string{"1.24.x", "1.25.x", "1.26.x"})
	server := releaseServer(t)
	defer server.Close()

	result, err := gosupport.CheckRepository(
		root,
		server.URL,
		time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC),
		server.Client(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !containsFinding(result.Findings, "six-month maximum of 2026-08-10") {
		t.Fatalf("expected maximum-grace finding, got: %v", result.Findings)
	}
}

func TestCheckRepositoryRequiresRetiredLineInGraceMatrix(t *testing.T) {
	root := makeRepository(t, "1.24", `"2026-08-01"`, []string{"1.25.x", "1.26.x"})
	server := releaseServer(t)
	defer server.Close()

	result, err := gosupport.CheckRepository(
		root,
		server.URL,
		time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC),
		server.Client(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !containsFinding(result.Findings, "does not test Go 1.24.x") {
		t.Fatalf("expected grace-matrix finding, got: %v", result.Findings)
	}
}

func TestCheckRepositoryDoesNotCountCommentedMatrixVersion(t *testing.T) {
	root := makeRepository(t, "1.25", "null", []string{"1.25.x"})
	workflow := filepath.Join(root, ".github", "workflows", "ci.yml")
	file, err := os.OpenFile(workflow, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fmt.Fprintln(file, "# The old checker incorrectly counted 1.26.x in comments."); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	server := releaseServer(t)
	defer server.Close()

	result, err := gosupport.CheckRepository(
		root,
		server.URL,
		time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC),
		server.Client(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !containsFinding(result.Findings, "does not test Go 1.26.x") {
		t.Fatalf("expected semantic-matrix finding, got: %v", result.Findings)
	}
}

func TestSupportedLinesIgnorePrereleasesAndPatchVersions(t *testing.T) {
	lines, err := gosupport.SupportedLines([]gosupport.GoRelease{
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
		t.Fatalf("SupportedLines() = %s", got)
	}
}

func releaseServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		fmt.Fprint(response, `[
			{"version":"go1.26.5","time":"2026-07-07T00:00:00Z"},
			{"version":"go1.26.0","time":"2026-02-10T00:00:00Z"},
			{"version":"go1.25.12","time":"2026-07-07T00:00:00Z"},
			{"version":"go1.25.0","time":"2025-08-12T00:00:00Z"},
			{"version":"go1.24.18","time":"2026-02-03T00:00:00Z"},
			{"version":"go1.27rc1","time":"2026-07-15T00:00:00Z"}
		]`)
	}))
}

func makeRepository(t *testing.T, minimum, graceUntil string, matrix []string) string {
	t.Helper()
	root := t.TempDir()

	matrixLines := make([]string, 0, len(matrix))
	for _, line := range matrix {
		matrixLines = append(matrixLines, "          - '"+line+"'")
	}
	workflow := fmt.Sprintf(`jobs:
  lint:
    env:
      GOTOOLCHAIN: local
    steps:
      - run: ./scripts/check-go-mod
  test:
    env:
      GOTOOLCHAIN: local
    strategy:
      matrix:
        go-version:
%s
    steps:
      - working-directory: internal/testdata/consumer
        run: go test -mod=readonly ./...
  vulnerability:
    env:
      GOTOOLCHAIN: local
    steps:
      - run: govulncheck ./...
      - working-directory: examples
        run: govulncheck ./...
`, strings.Join(matrixLines, "\n"))

	releaseText := "openai-go v3.45.0 requires Go " + minimum +
		" or newer. v3.44.0 is the final compatible release."
	files := map[string]string{
		".github/go-support-policy.json": fmt.Sprintf(
			`{
				"minimum":%q,
				"grace_until":%s,
				"reason":"test policy",
				"approved_by":"SDK team",
				"release_note":{
					"sdk_version":"v3.45.0",
					"last_compatible_sdk_version":"v3.44.0",
					"text":%q
				}
			}`,
			minimum,
			graceUntil,
			releaseText,
		),
		"go.mod":                            "module example.com/root\n\ngo " + minimum + ".0\n",
		"examples/go.mod":                   "module example.com/examples\n\ngo " + minimum + ".0\n",
		"internal/testdata/consumer/go.mod": "module example.com/consumer\n\ngo " + minimum + ".0\n",
		"tools/go.mod":                      "module example.com/tools\n\ngo " + minimum + ".0\n",
		"README.md":                         "This library requires Go " + minimum + ".\n",
		"CONTRIBUTING.md":                   "Install Go " + minimum + ".\n",
		"GO_VERSION_POLICY.md":              "The minimum is Go " + minimum + ".\n",
		".github/workflows/ci.yml":          workflow,
		"scripts/check-go-mod": strings.Join([]string{
			"go mod tidy -diff",
			"(cd examples && go mod tidy -diff)",
			"(cd internal/testdata/consumer && go mod tidy -diff)",
			"(cd tools && go mod tidy -diff)",
			"",
		}, "\n"),
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
