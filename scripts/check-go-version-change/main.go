package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var releaseNoteHeading = regexp.MustCompile(`(?mi)^##[ \t]+Release note[ \t]*\r?$`)

type directives struct {
	Root     string
	Examples string
}

type validationInput struct {
	Base       directives
	Head       directives
	Changed    map[string]bool
	PRBody     string
	Codeowners string
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

	base, err := directivesAt(root, os.Args[1])
	if err != nil {
		fatal(err)
	}
	head, err := directivesAt(root, "")
	if err != nil {
		fatal(err)
	}

	if base == head {
		fmt.Println("Go directives are unchanged.")
		return
	}

	fmt.Println("Go directive change detected:")
	fmt.Printf("  root:     %s -> %s\n", base.Root, head.Root)
	fmt.Printf("  examples: %s -> %s\n", base.Examples, head.Examples)

	changed, err := changedFiles(root, os.Args[1])
	if err != nil {
		fatal(err)
	}
	codeowners, err := os.ReadFile(filepath.Join(root, ".github", "CODEOWNERS"))
	if err != nil {
		fatal(err)
	}

	problems := validate(validationInput{
		Base:       base,
		Head:       head,
		Changed:    changed,
		PRBody:     os.Getenv("PR_BODY"),
		Codeowners: string(codeowners),
	})
	if len(problems) > 0 {
		fmt.Fprintln(os.Stderr, "A Go directive change is missing required policy artifacts:")
		for _, problem := range problems {
			fmt.Fprintf(os.Stderr, "  - %s\n", problem)
		}
		os.Exit(1)
	}

	fmt.Println("Go directive policy checks passed.")
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func directivesAt(root, revision string) (directives, error) {
	read := func(relative string) (string, error) {
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
			return "", err
		}
		return goDirective(data)
	}

	rootDirective, err := read("go.mod")
	if err != nil {
		return directives{}, err
	}
	examplesDirective, err := read(filepath.ToSlash(filepath.Join("examples", "go.mod")))
	if err != nil {
		return directives{}, err
	}
	return directives{Root: rootDirective, Examples: examplesDirective}, nil
}

func goDirective(data []byte) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 && fields[0] == "go" {
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
	commands := [][]string{
		{"diff", "--name-only", base + "...HEAD"},
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
	if input.Head.Root != input.Head.Examples {
		problems = append(problems, "the root and examples modules must use the same Go version")
	}

	for _, required := range []string{"README.md", "CONTRIBUTING.md", "GO_VERSION_POLICY.md"} {
		if !input.Changed[required] {
			problems = append(problems, required)
		}
	}

	minimum, err := majorMinor(input.Head.Root)
	if err != nil {
		problems = append(problems, err.Error())
	} else if !releaseNoteHeading.MatchString(input.PRBody) ||
		!strings.Contains(input.PRBody, "Go "+minimum) {
		problems = append(problems, fmt.Sprintf(
			"a PR description with a '## Release note' section mentioning Go %s", minimum,
		))
	}

	for _, path := range []string{"/go.mod", "/examples/go.mod"} {
		if !hasCodeowner(input.Codeowners, path, "@openai/sdks-team") {
			problems = append(problems, fmt.Sprintf(
				"an explicit %s CODEOWNERS entry for @openai/sdks-team", path,
			))
		}
	}
	return problems
}

func majorMinor(raw string) (string, error) {
	parts := strings.Split(raw, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid Go version %q", raw)
	}
	return strings.Join(parts[:2], "."), nil
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
