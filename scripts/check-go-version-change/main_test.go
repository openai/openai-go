package main

import (
	"strings"
	"testing"

	"github.com/openai/openai-go/v3/scripts/internal/gosupport"
)

func TestValidateAcceptsCompleteGoVersionChange(t *testing.T) {
	input := completeInput()
	if problems := validate(input); len(problems) != 0 {
		t.Fatalf("validate() returned problems: %v", problems)
	}
}

func TestValidateIgnoresUnchangedDirectives(t *testing.T) {
	input := validationInput{
		Base: directives{Root: "1.25.0", Examples: "1.25.0", Consumer: "1.25.0", Tools: "1.25.0"},
		Head: directives{Root: "1.25.0", Examples: "1.25.0", Consumer: "1.25.0", Tools: "1.25.0"},
	}
	if problems := validate(input); len(problems) != 0 {
		t.Fatalf("validate() returned problems: %v", problems)
	}
}

func TestValidateRequiresEveryPolicyArtifact(t *testing.T) {
	for _, path := range []string{
		"README.md",
		"CONTRIBUTING.md",
		"GO_VERSION_POLICY.md",
		".github/go-support-policy.json",
		".github/workflows/ci.yml",
		"internal/testdata/consumer/go.mod",
		"tools/go.mod",
	} {
		t.Run(path, func(t *testing.T) {
			input := completeInput()
			delete(input.Changed, path)
			if problems := validate(input); !containsProblem(problems, path) {
				t.Fatalf("expected %s problem, got: %v", path, problems)
			}
		})
	}
}

func TestValidateRequiresApprovedTextInsideReleaseNoteSection(t *testing.T) {
	input := completeInput()
	input.PRBody = "The summary mentions Go 1.25 and v3.45.0.\n\n## Release note\n\nNothing useful."
	problems := validate(input)
	if !containsProblem(problems, "release_note.text") {
		t.Fatalf("expected release-note problem, got: %v", problems)
	}
}

func TestValidateRequiresMatchingModuleFloors(t *testing.T) {
	for _, mutate := range []func(*validationInput){
		func(input *validationInput) { input.Head.Examples = "1.24.0" },
		func(input *validationInput) { input.Head.Consumer = "1.24.0" },
		func(input *validationInput) { input.Head.Tools = "1.24.0" },
	} {
		input := completeInput()
		mutate(&input)
		problems := validate(input)
		if !containsProblem(problems, "must use the same Go version") {
			t.Fatalf("expected module-floor problem, got: %v", problems)
		}
	}
}

func TestValidateRequiresPolicyMinimumToMatchGoMod(t *testing.T) {
	input := completeInput()
	input.Policy.Minimum = "1.24"
	problems := validate(input)
	if !containsProblem(problems, "policy.json declares Go 1.24") {
		t.Fatalf("expected policy-minimum problem, got: %v", problems)
	}
}

func TestGoDirectiveAllowsTrailingComment(t *testing.T) {
	got, err := goDirective([]byte("module example.com/test\n\ngo 1.25.0 // minimum\n"))
	if err != nil {
		t.Fatal(err)
	}
	if got != "1.25.0" {
		t.Fatalf("goDirective() = %q", got)
	}
}

func completeInput() validationInput {
	note := "openai-go v3.45.0 requires Go 1.25 or newer. " +
		"v3.44.0 is the final compatible release."
	return validationInput{
		Base: directives{Root: "1.22", Examples: "1.22.4"},
		Head: directives{Root: "1.25.0", Examples: "1.25.0", Consumer: "1.25.0", Tools: "1.25.0"},
		Changed: map[string]bool{
			"README.md":                         true,
			"CONTRIBUTING.md":                   true,
			"GO_VERSION_POLICY.md":              true,
			".github/go-support-policy.json":    true,
			".github/workflows/ci.yml":          true,
			"internal/testdata/consumer/go.mod": true,
			"tools/go.mod":                      true,
		},
		PRBody: "## Release note\n\n" + note,
		Codeowners: `
/go.mod @openai/sdks-team
/examples/go.mod @openai/sdks-team
/internal/testdata/consumer/go.mod @openai/sdks-team
/tools/go.mod @openai/sdks-team
/.github/go-support-policy.json @openai/sdks-team
/.github/workflows/ci.yml @openai/sdks-team
`,
		Policy: gosupport.Policy{
			Minimum: "1.25",
			Release: gosupport.ReleaseNote{
				SDKVersion:               "v3.45.0",
				LastCompatibleSDKVersion: "v3.44.0",
				Text:                     note,
			},
		},
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
