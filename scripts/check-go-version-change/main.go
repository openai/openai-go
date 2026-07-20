package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/openai/openai-go/v3/scripts/internal/gosupport"
)

var releaseNoteHeading = regexp.MustCompile(`(?mi)^##[ \t]+Release note[ \t]*\r?$`)

type directives struct {
	Root     string
	Examples string
	Consumer string
	Tools    string
}

type validationInput struct {
	Base       directives
	Head       directives
	Changed    map[string]bool
	PRBody     string
	Codeowners string
	Policy     gosupport.Policy
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./scripts/check-go-version-change BASE_SHA")
		os.Exit(2)
	}

	root, err := gitOutput("", "rev-parse", "--show-toplevel")
	if err != nil {
		fatal(err)
	}

	base, err := directivesAt(root, os.Args[1], true)
	if err != nil {
		fatal(err)
	}
	head, err := directivesAt(root, "", false)
	if err != nil {
		fatal(err)
	}

	if base == head {
		fmt.Println("Go directives are unchanged.")
		return
	}

	fmt.Println("Go directive change detected:")
	fmt.Printf("  root:     %s -> %s\n", displayDirective(base.Root), head.Root)
	fmt.Printf("  examples: %s -> %s\n", displayDirective(base.Examples), head.Examples)
	fmt.Printf("  consumer: %s -> %s\n", displayDirective(base.Consumer), head.Consumer)
	fmt.Printf("  tools:    %s -> %s\n", displayDirective(base.Tools), head.Tools)

	changed, err := changedFiles(root, os.Args[1])
	if err != nil {
		fatal(err)
	}
	codeowners, err := os.ReadFile(filepath.Join(root, ".github", "CODEOWNERS"))
	if err != nil {
		fatal(err)
	}
	policy, err := gosupport.ReadPolicy(filepath.Join(root, ".github", "go-support-policy.json"))
	if err != nil {
		fatal(err)
	}

	problems := validate(validationInput{
		Base:       base,
		Head:       head,
		Changed:    changed,
		PRBody:     os.Getenv("PR_BODY"),
		Codeowners: string(codeowners),
		Policy:     policy,
	})
	if len(problems) > 0 {
		fmt.Fprintln(os.Stderr, "A Go directive change is missing required policy artifacts:")
		for _, problem := range problems {
			fmt.Fprintf(os.Stderr, "  - %s\n", problem)
		}
		os.Exit(1)
	}

	// A Go-floor PR is rare and high impact. Validate its complete repository
	// state against the live official release feed before allowing it to merge.
	review, err := gosupport.CheckRepository(
		root,
		gosupport.DefaultReleaseFeed,
		time.Now(),
		&http.Client{Timeout: 30 * time.Second},
	)
	if err != nil {
		fatal(err)
	}
	if len(review.Findings) > 0 {
		fmt.Fprint(os.Stderr, review.Markdown())
		os.Exit(1)
	}

	fmt.Println("Go directive policy checks passed.")
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func directivesAt(root, revision string, allowMissingNested bool) (directives, error) {
	read := func(relative string, allowMissing bool) (string, error) {
		var data []byte
		var err error
		if revision == "" {
			data, err = os.ReadFile(filepath.Join(root, relative))
		} else {
			var output string
			output, err = gitOutput(root, "show", revision+":"+relative)
			data = []byte(output)
		}
		if err != nil {
			if allowMissing {
				return "", nil
			}
			return "", err
		}
		return goDirective(data)
	}

	rootDirective, err := read("go.mod", false)
	if err != nil {
		return directives{}, err
	}
	examplesDirective, err := read(filepath.ToSlash(filepath.Join("examples", "go.mod")), false)
	if err != nil {
		return directives{}, err
	}
	consumerDirective, err := read(
		filepath.ToSlash(filepath.Join("internal", "testdata", "consumer", "go.mod")),
		allowMissingNested,
	)
	if err != nil {
		return directives{}, err
	}
	toolsDirective, err := read(filepath.ToSlash(filepath.Join("tools", "go.mod")), allowMissingNested)
	if err != nil {
		return directives{}, err
	}
	return directives{
		Root:     rootDirective,
		Examples: examplesDirective,
		Consumer: consumerDirective,
		Tools:    toolsDirective,
	}, nil
}

func goDirective(data []byte) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && fields[0] == "go" {
			if _, err := gosupport.ParseVersion(fields[1]); err != nil {
				return "", err
			}
			return fields[1], nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", errors.New("go.mod has no go directive")
}

func changedFiles(root, base string) (map[string]bool, error) {
	changed := make(map[string]bool)
	// The first command covers committed PR changes. The remaining commands
	// make the checker useful before commit by including staged, unstaged, and
	// untracked files.
	commands := [][]string{
		{"diff", "--name-only", base + "...HEAD"},
		{"diff", "--cached", "--name-only"},
		{"diff", "--name-only"},
		{"ls-files", "--others", "--exclude-standard"},
	}
	for _, args := range commands {
		output, err := gitOutput(root, args...)
		if err != nil {
			return nil, err
		}
		for _, path := range strings.Split(output, "\n") {
			if path = strings.TrimSpace(path); path != "" {
				changed[path] = true
			}
		}
	}
	return changed, nil
}

func gitOutput(root string, args ...string) (string, error) {
	command := exec.Command("git", args...)
	command.Dir = root
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, output)
	}
	return strings.TrimSpace(string(output)), nil
}

func validate(input validationInput) []string {
	if input.Base == input.Head {
		return nil
	}

	var problems []string
	if input.Head.Root != input.Head.Examples ||
		input.Head.Root != input.Head.Consumer ||
		input.Head.Root != input.Head.Tools {
		problems = append(problems, "the root, examples, consumer, and tools modules must use the same Go version")
	}

	for _, required := range []string{
		"README.md",
		"CONTRIBUTING.md",
		"GO_VERSION_POLICY.md",
		".github/go-support-policy.json",
		".github/workflows/ci.yml",
		"internal/testdata/consumer/go.mod",
		"tools/go.mod",
	} {
		if !input.Changed[required] {
			problems = append(problems, required)
		}
	}

	minimum, err := majorMinor(input.Head.Root)
	if err != nil {
		problems = append(problems, err.Error())
	} else if input.Policy.Minimum != minimum {
		problems = append(problems, fmt.Sprintf(
			".github/go-support-policy.json declares Go %s; go.mod declares Go %s",
			input.Policy.Minimum,
			minimum,
		))
	}

	releaseSection := releaseNoteSection(input.PRBody)
	if releaseSection == "" || !containsNormalized(releaseSection, input.Policy.Release.Text) {
		problems = append(problems, "the PR's '## Release note' section must contain release_note.text from .github/go-support-policy.json")
	}

	for _, path := range []string{
		"/go.mod",
		"/examples/go.mod",
		"/internal/testdata/consumer/go.mod",
		"/tools/go.mod",
		"/.github/go-support-policy.json",
		"/.github/workflows/ci.yml",
	} {
		if !hasCodeowner(input.Codeowners, path, "@openai/sdks-team") {
			problems = append(problems, fmt.Sprintf(
				"an explicit %s CODEOWNERS entry for @openai/sdks-team", path,
			))
		}
	}
	return problems
}

func majorMinor(raw string) (string, error) {
	parsed, err := gosupport.ParseVersion(raw)
	if err != nil {
		return "", err
	}
	return parsed.String(), nil
}

func releaseNoteSection(body string) string {
	location := releaseNoteHeading.FindStringIndex(body)
	if location == nil {
		return ""
	}
	section := body[location[1]:]
	nextHeading := regexp.MustCompile(`(?m)^##[ \t]+`).FindStringIndex(section)
	if nextHeading != nil {
		section = section[:nextHeading[0]]
	}
	return section
}

func containsNormalized(haystack, needle string) bool {
	normalize := func(value string) string {
		return strings.Join(strings.Fields(value), " ")
	}
	return strings.Contains(normalize(haystack), normalize(needle))
}

func hasCodeowner(contents, path, owner string) bool {
	scanner := bufio.NewScanner(strings.NewReader(contents))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 || fields[0] != path {
			continue
		}
		for _, candidate := range fields[1:] {
			if candidate == owner {
				return true
			}
		}
	}
	return false
}

func displayDirective(value string) string {
	if value == "" {
		return "(missing)"
	}
	return value
}
