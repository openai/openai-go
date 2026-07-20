package main

import (
	"strings"
	"testing"
)

func TestValidateAcceptsCompleteGoVersionChange(t *testing.T) {
	input := completeInput()
	if problems := validate(input); len(problems) != 0 {
		t.Fatalf("validate() returned problems: %v", problems)
	}
}

func TestValidateIgnoresUnchangedDirectives(t *testing.T) {
	input := validationInput{
		Base: directives{Root: "1.25.0", Examples: "1.25.0"},
		Head: directives{Root: "1.25.0", Examples: "1.25.0"},
	}
	if problems := validate(input); len(problems) != 0 {
		t.Fatalf("validate() returned problems: %v", problems)
	}
}

func TestValidateRequiresReleaseNote(t *testing.T) {
	input := completeInput()
	input.PRBody = ""
	problems := validate(input)
	if !containsProblem(problems, "## Release note") {
		t.Fatalf("expected release-note problem, got: %v", problems)
	}
}

func TestValidateRequiresMatchingModuleFloors(t *testing.T) {
	input := completeInput()
	input.Head.Examples = "1.24.0"
	problems := validate(input)
	if !containsProblem(problems, "same Go version") {
		t.Fatalf("expected module-floor problem, got: %v", problems)
	}
}

func TestGoDirective(t *testing.T) {
	got, err := goDirective([]byte("module example.com/test\n\ngo 1.25.0\n"))
	if err != nil {
		t.Fatal(err)
	}
	if got != "1.25.0" {
		t.Fatalf("goDirective() = %q", got)
	}
}

func completeInput() validationInput {
	return validationInput{
		Base: directives{Root: "1.22", Examples: "1.22.4"},
		Head: directives{Root: "1.25.0", Examples: "1.25.0"},
		Changed: map[string]bool{
			"README.md":            true,
			"CONTRIBUTING.md":      true,
			"GO_VERSION_POLICY.md": true,
		},
		PRBody: "## Release note\n\nThis release requires Go 1.25 or newer.",
		Codeowners: `
/go.mod @openai/sdks-team
/examples/go.mod @openai/sdks-team
`,
	}
}

func containsProblem(problems []string, needle string) bool {
	for _, problem := range problems {
		if strings.Contains(problem, needle) {
			return true
		}
	}
	return false
}
