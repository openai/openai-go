package gosupport

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// DefaultReleaseFeed includes retired Go release lines as well as the two
// supported lines. The retired line is needed to validate an approved grace
// period.
const DefaultReleaseFeed = "https://go.dev/dl/?mode=json&include=all"

var (
	stableVersionPattern = regexp.MustCompile(`^(?:go)?([0-9]+)\.([0-9]+)(?:\.([0-9]+))?$`)
	sdkVersionPattern    = regexp.MustCompile(`^v([0-9]+)\.([0-9]+)\.([0-9]+)$`)
)

type ReleaseNote struct {
	SDKVersion               string `json:"sdk_version"`
	LastCompatibleSDKVersion string `json:"last_compatible_sdk_version"`
	Text                     string `json:"text"`
}

type Policy struct {
	Minimum    string      `json:"minimum"`
	GraceUntil *string     `json:"grace_until"`
	Reason     string      `json:"reason"`
	ApprovedBy string      `json:"approved_by"`
	Release    ReleaseNote `json:"release_note"`
}

type GoRelease struct {
	Version string    `json:"version"`
	Time    time.Time `json:"time"`
}

type Version struct {
	Major int
	Minor int
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

type Review struct {
	GeneratedAt time.Time
	Current     Version
	Previous    Version
	Retired     Version
	RetiredAt   *time.Time
	Minimum     Version
	Policy      Policy
	Findings    []string
}

func (r Review) Markdown() string {
	var b strings.Builder
	fmt.Fprintln(&b, "# Go version support review")
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "Generated: %s\n\n", r.GeneratedAt.UTC().Format(time.RFC3339))

	if r.Current != (Version{}) {
		fmt.Fprintf(&b, "- Official current stable release: Go %s\n", r.Current)
		fmt.Fprintf(&b, "- Official previous stable release: Go %s\n", r.Previous)
		fmt.Fprintf(&b, "- Most recently retired release: Go %s\n", r.Retired)
	}
	if r.Minimum != (Version{}) {
		fmt.Fprintf(&b, "- Repository minimum: Go %s\n", r.Minimum)
	}
	if r.Policy.GraceUntil != nil {
		fmt.Fprintf(&b, "- Grace period ends: %s\n", *r.Policy.GraceUntil)
	}
	if r.Policy.Reason != "" {
		fmt.Fprintf(&b, "- Recorded reason: %s\n", r.Policy.Reason)
	}
	if r.Policy.ApprovedBy != "" {
		fmt.Fprintf(&b, "- Recorded approver: %s\n", r.Policy.ApprovedBy)
	}
	if r.Policy.Release.SDKVersion != "" {
		fmt.Fprintf(&b, "- Go-floor release: %s\n", r.Policy.Release.SDKVersion)
	}

	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "## Result")
	fmt.Fprintln(&b)
	if len(r.Findings) == 0 {
		fmt.Fprintln(&b, "No policy drift detected.")
		return b.String()
	}

	fmt.Fprintln(&b, "Review is required:")
	fmt.Fprintln(&b)
	for _, finding := range r.Findings {
		fmt.Fprintf(&b, "- %s\n", finding)
	}
	return b.String()
}

func ParseVersion(raw string) (Version, error) {
	match := stableVersionPattern.FindStringSubmatch(strings.TrimSpace(raw))
	if match == nil {
		return Version{}, fmt.Errorf("invalid stable Go version %q", raw)
	}
	major, err := strconv.Atoi(match[1])
	if err != nil {
		return Version{}, err
	}
	minor, err := strconv.Atoi(match[2])
	if err != nil {
		return Version{}, err
	}
	return Version{Major: major, Minor: minor}, nil
}

func SupportedLines(releases []GoRelease) ([]Version, error) {
	unique := make(map[Version]struct{})
	for _, release := range releases {
		parsed, err := ParseVersion(release.Version)
		if err != nil {
			continue // Ignore betas, release candidates, and malformed entries.
		}
		unique[parsed] = struct{}{}
	}

	versions := make([]Version, 0, len(unique))
	for release := range unique {
		versions = append(versions, release)
	}
	sort.Slice(versions, func(i, j int) bool {
		if versions[i].Major != versions[j].Major {
			return versions[i].Major > versions[j].Major
		}
		return versions[i].Minor > versions[j].Minor
	})
	if len(versions) < 3 {
		return nil, errors.New("official feed did not contain three stable Go release lines")
	}
	return versions[:3], nil
}

func ReadPolicy(path string) (Policy, error) {
	var policy Policy
	data, err := os.ReadFile(path)
	if err != nil {
		return policy, err
	}
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&policy); err != nil {
		return policy, fmt.Errorf("parse %s: %w", path, err)
	}
	return policy, nil
}

func ReadGoDirective(path string) (Version, error) {
	file, err := os.Open(path)
	if err != nil {
		return Version{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && fields[0] == "go" {
			return ParseVersion(fields[1])
		}
	}
	if err := scanner.Err(); err != nil {
		return Version{}, err
	}
	return Version{}, fmt.Errorf("%s has no go directive", path)
}

func FetchReleases(client *http.Client, feedURL string) ([]GoRelease, error) {
	response, err := client.Get(feedURL)
	if err != nil {
		return nil, fmt.Errorf("fetch official Go release feed: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, response.Body)
		return nil, fmt.Errorf("fetch official Go release feed: %s", response.Status)
	}

	var releases []GoRelease
	if err := json.NewDecoder(response.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("decode official Go release feed: %w", err)
	}
	return releases, nil
}

func CheckRepository(root, feedURL string, now time.Time, client *http.Client) (Review, error) {
	result := Review{GeneratedAt: now}

	policy, err := ReadPolicy(filepath.Join(root, ".github", "go-support-policy.json"))
	if err != nil {
		return result, err
	}
	result.Policy = policy

	minimum, err := ParseVersion(policy.Minimum)
	if err != nil {
		return result, fmt.Errorf("parse policy minimum: %w", err)
	}
	result.Minimum = minimum

	releases, err := FetchReleases(client, feedURL)
	if err != nil {
		return result, err
	}
	lines, err := SupportedLines(releases)
	if err != nil {
		return result, err
	}
	result.Current, result.Previous, result.Retired = lines[0], lines[1], lines[2]
	if retiredAt, found := retirementDate(releases, result.Current); found {
		result.RetiredAt = &retiredAt
	}

	result.Findings = append(result.Findings, validatePolicy(policy, minimum)...)

	switch minimum {
	case result.Previous:
		if policy.GraceUntil != nil {
			result.Findings = append(result.Findings, "grace_until must be null when the minimum is the previous stable release")
		}
	case result.Retired:
		result.Findings = append(result.Findings, validateGracePeriod(policy.GraceUntil, result.Retired, result.RetiredAt, now)...)
	default:
		result.Findings = append(result.Findings, fmt.Sprintf(
			"minimum Go %s is neither the previous stable release (Go %s) nor the grace-eligible retired release (Go %s)",
			minimum, result.Previous, result.Retired,
		))
	}

	moduleFiles := []string{
		"go.mod",
		filepath.Join("examples", "go.mod"),
		filepath.Join("internal", "testdata", "consumer", "go.mod"),
		filepath.Join("tools", "go.mod"),
	}
	for _, relative := range moduleFiles {
		directive, readErr := ReadGoDirective(filepath.Join(root, relative))
		if readErr != nil {
			result.Findings = append(result.Findings, readErr.Error())
			continue
		}
		if directive != minimum {
			result.Findings = append(result.Findings, fmt.Sprintf(
				"%s declares Go %s; policy declares Go %s", relative, directive, minimum,
			))
		}
	}

	documentNeedles := map[string]string{
		"README.md":            "requires Go " + minimum.String(),
		"CONTRIBUTING.md":      "install Go " + minimum.String(),
		"GO_VERSION_POLICY.md": "Go " + minimum.String(),
	}
	for relative, needle := range documentNeedles {
		data, readErr := os.ReadFile(filepath.Join(root, relative))
		if readErr != nil {
			result.Findings = append(result.Findings, readErr.Error())
			continue
		}
		if !strings.Contains(strings.ToLower(string(data)), strings.ToLower(needle)) {
			result.Findings = append(result.Findings, fmt.Sprintf("%s does not contain %q", relative, needle))
		}
	}

	expectedLines := []Version{result.Previous, result.Current}
	if minimum == result.Retired {
		// A grace period is a support promise. Test the retired minimum as well
		// as the two release lines still supported by the Go project.
		expectedLines = append([]Version{result.Retired}, expectedLines...)
	}
	result.Findings = append(result.Findings, checkCI(filepath.Join(root, ".github", "workflows", "ci.yml"), expectedLines)...)
	result.Findings = append(result.Findings, checkTidyScript(filepath.Join(root, "scripts", "check-go-mod"))...)

	return result, nil
}

func validatePolicy(policy Policy, minimum Version) []string {
	var findings []string
	if strings.TrimSpace(policy.Reason) == "" {
		findings = append(findings, "policy file has no reason")
	}
	if strings.TrimSpace(policy.ApprovedBy) == "" {
		findings = append(findings, "policy file has no approver")
	}

	if !sdkVersionPattern.MatchString(policy.Release.SDKVersion) {
		findings = append(findings, "release_note.sdk_version must be a v-prefixed semantic SDK version")
	}
	if !sdkVersionPattern.MatchString(policy.Release.LastCompatibleSDKVersion) {
		findings = append(findings, "release_note.last_compatible_sdk_version must be a v-prefixed semantic SDK version")
	}
	if compareSDKVersions(policy.Release.LastCompatibleSDKVersion, policy.Release.SDKVersion) >= 0 {
		findings = append(findings, "release_note.last_compatible_sdk_version must precede release_note.sdk_version")
	}

	note := normalizeWhitespace(policy.Release.Text)
	for _, required := range []string{
		"Go " + minimum.String(),
		policy.Release.SDKVersion,
		policy.Release.LastCompatibleSDKVersion,
	} {
		if required != "" && !strings.Contains(note, required) {
			findings = append(findings, fmt.Sprintf("release_note.text does not mention %q", required))
		}
	}
	return findings
}

func validateGracePeriod(raw *string, retired Version, retiredAt *time.Time, now time.Time) []string {
	if raw == nil {
		return []string{fmt.Sprintf("Go %s is retired and requires an approved grace_until date", retired)}
	}
	graceDate, err := time.Parse("2006-01-02", *raw)
	if err != nil {
		return []string{"grace_until must use YYYY-MM-DD"}
	}

	today := dateOnly(now)
	var findings []string
	if today.After(graceDate) {
		findings = append(findings, fmt.Sprintf("grace period for Go %s expired on %s", retired, *raw))
	}
	if retiredAt == nil {
		findings = append(findings, "official feed did not identify when the retired Go line left support")
		return findings
	}
	maximum := dateOnly(*retiredAt).AddDate(0, 6, 0)
	if graceDate.After(maximum) {
		findings = append(findings, fmt.Sprintf(
			"grace period for Go %s ends after the six-month maximum of %s", retired, maximum.Format("2006-01-02"),
		))
	}
	return findings
}

func retirementDate(releases []GoRelease, current Version) (time.Time, bool) {
	var earliest time.Time
	for _, release := range releases {
		line, err := ParseVersion(release.Version)
		if err != nil || line != current || release.Time.IsZero() {
			continue
		}
		if earliest.IsZero() || release.Time.Before(earliest) {
			earliest = release.Time
		}
	}
	if earliest.IsZero() {
		return time.Time{}, false
	}
	return dateOnly(earliest), true
}

func dateOnly(value time.Time) time.Time {
	utc := value.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}

func compareSDKVersions(left, right string) int {
	parse := func(raw string) [3]int {
		match := sdkVersionPattern.FindStringSubmatch(raw)
		if match == nil {
			return [3]int{}
		}
		var parsed [3]int
		for index := range parsed {
			parsed[index], _ = strconv.Atoi(match[index+1])
		}
		return parsed
	}
	a, b := parse(left), parse(right)
	for index := range a {
		if a[index] < b[index] {
			return -1
		}
		if a[index] > b[index] {
			return 1
		}
	}
	return 0
}

type workflowJob struct {
	Environment map[string]string
	GoVersions  []string
	Steps       []workflowStep
}

type workflowStep struct {
	WorkingDirectory string
	Run              string
}

func checkCI(path string, expected []Version) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{err.Error()}
	}
	lines := parseYAMLLines(data)

	var findings []string
	testJob, ok := parseWorkflowJob(lines, "test")
	if !ok {
		return []string{".github/workflows/ci.yml has no test job"}
	}
	if testJob.Environment["GOTOOLCHAIN"] != "local" {
		findings = append(findings, "CI test job does not set GOTOOLCHAIN to local")
	}

	actualLines := make(map[Version]bool)
	for _, raw := range testJob.GoVersions {
		raw = strings.TrimSuffix(strings.TrimSpace(raw), ".x")
		parsed, parseErr := ParseVersion(raw)
		if parseErr != nil {
			findings = append(findings, fmt.Sprintf("CI test matrix has invalid Go version %q", raw))
			continue
		}
		actualLines[parsed] = true
	}
	expectedLines := make(map[Version]bool)
	for _, line := range expected {
		expectedLines[line] = true
		if !actualLines[line] {
			findings = append(findings, fmt.Sprintf("CI test matrix does not test Go %s.x", line))
		}
	}
	for line := range actualLines {
		if !expectedLines[line] {
			findings = append(findings, fmt.Sprintf("CI test matrix unexpectedly tests Go %s.x", line))
		}
	}

	lintJob, ok := parseWorkflowJob(lines, "lint")
	if !ok {
		findings = append(findings, ".github/workflows/ci.yml has no lint job")
	} else {
		if lintJob.Environment["GOTOOLCHAIN"] != "local" {
			findings = append(findings, "CI lint job does not set GOTOOLCHAIN to local")
		}
		if !hasRunStep(lintJob, "", "./scripts/check-go-mod") {
			findings = append(findings, "CI lint job does not run ./scripts/check-go-mod")
		}
	}

	if !hasRunStep(testJob, "internal/testdata/consumer", "go test -mod=readonly ./...") {
		findings = append(findings, "CI test job does not test the consumer module with -mod=readonly")
	}

	vulnerabilityJob, ok := parseWorkflowJob(lines, "vulnerability")
	if !ok {
		findings = append(findings, ".github/workflows/ci.yml has no vulnerability job")
	} else {
		if vulnerabilityJob.Environment["GOTOOLCHAIN"] != "local" {
			findings = append(findings, "CI vulnerability job does not set GOTOOLCHAIN to local")
		}
		if !hasRunStep(vulnerabilityJob, "", "govulncheck ./...") {
			findings = append(findings, "CI vulnerability job does not scan the root module")
		}
		if !hasRunStep(vulnerabilityJob, "examples", "govulncheck ./...") {
			findings = append(findings, "CI vulnerability job does not scan the examples module")
		}
	}
	return findings
}

type yamlLine struct {
	Indent int
	Text   string
}

func parseYAMLLines(data []byte) []yamlLine {
	var lines []yamlLine
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		raw := strings.TrimRight(scanner.Text(), " \t\r")
		raw = stripYAMLComment(raw)
		if strings.TrimSpace(raw) == "" {
			continue
		}
		indent := len(raw) - len(strings.TrimLeft(raw, " "))
		lines = append(lines, yamlLine{Indent: indent, Text: strings.TrimSpace(raw)})
	}
	return lines
}

func stripYAMLComment(value string) string {
	var singleQuoted, doubleQuoted bool
	for index, character := range value {
		switch character {
		case '\'':
			if !doubleQuoted {
				singleQuoted = !singleQuoted
			}
		case '"':
			if !singleQuoted {
				doubleQuoted = !doubleQuoted
			}
		case '#':
			if !singleQuoted && !doubleQuoted {
				return strings.TrimRight(value[:index], " \t")
			}
		}
	}
	return value
}

func parseWorkflowJob(lines []yamlLine, name string) (workflowJob, bool) {
	var jobLines []yamlLine
	foundJobs := false
	for index, line := range lines {
		if line.Indent == 0 && line.Text == "jobs:" {
			foundJobs = true
			continue
		}
		if !foundJobs || line.Indent != 2 || line.Text != name+":" {
			continue
		}
		for _, candidate := range lines[index+1:] {
			if candidate.Indent <= 2 {
				break
			}
			jobLines = append(jobLines, candidate)
		}
		break
	}
	if len(jobLines) == 0 {
		return workflowJob{}, false
	}

	job := workflowJob{Environment: make(map[string]string)}
	for index, line := range jobLines {
		switch {
		case line.Indent == 4 && line.Text == "env:":
			for _, child := range jobLines[index+1:] {
				if child.Indent <= line.Indent {
					break
				}
				if child.Indent == 6 {
					key, value, ok := yamlScalar(child.Text)
					if ok {
						job.Environment[key] = value
					}
				}
			}
		case line.Indent == 8 && strings.HasPrefix(line.Text, "go-version:"):
			_, inline, _ := yamlScalar(line.Text)
			if strings.HasPrefix(inline, "[") && strings.HasSuffix(inline, "]") {
				for _, value := range strings.Split(strings.Trim(inline, "[]"), ",") {
					if value = unquoteYAML(value); value != "" {
						job.GoVersions = append(job.GoVersions, value)
					}
				}
				continue
			}
			for _, child := range jobLines[index+1:] {
				if child.Indent <= line.Indent {
					break
				}
				if child.Indent == 10 && strings.HasPrefix(child.Text, "- ") {
					job.GoVersions = append(job.GoVersions, unquoteYAML(strings.TrimPrefix(child.Text, "- ")))
				}
			}
		case line.Indent == 4 && line.Text == "steps:":
			job.Steps = parseWorkflowSteps(jobLines[index+1:])
		}
	}
	return job, true
}

func parseWorkflowSteps(lines []yamlLine) []workflowStep {
	var steps []workflowStep
	var current *workflowStep
	for _, line := range lines {
		if line.Indent <= 4 {
			break
		}
		if line.Indent == 6 && strings.HasPrefix(line.Text, "- ") {
			steps = append(steps, workflowStep{})
			current = &steps[len(steps)-1]
			applyStepScalar(current, strings.TrimPrefix(line.Text, "- "))
			continue
		}
		if current != nil && line.Indent == 8 {
			applyStepScalar(current, line.Text)
		}
	}
	return steps
}

func applyStepScalar(step *workflowStep, text string) {
	key, value, ok := yamlScalar(text)
	if !ok {
		return
	}
	switch key {
	case "working-directory":
		step.WorkingDirectory = value
	case "run":
		step.Run = value
	}
}

func yamlScalar(text string) (string, string, bool) {
	key, value, found := strings.Cut(text, ":")
	if !found {
		return "", "", false
	}
	return strings.TrimSpace(key), unquoteYAML(value), true
}

func unquoteYAML(value string) string {
	return strings.Trim(strings.TrimSpace(value), `"'`)
}

func hasRunStep(job workflowJob, directory, command string) bool {
	for _, step := range job.Steps {
		if step.WorkingDirectory == directory && strings.TrimSpace(step.Run) == command {
			return true
		}
	}
	return false
}

func checkTidyScript(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{err.Error()}
	}
	commands := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			commands[line] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return []string{err.Error()}
	}

	var findings []string
	for _, required := range []string{
		"go mod tidy -diff",
		"(cd examples && go mod tidy -diff)",
		"(cd internal/testdata/consumer && go mod tidy -diff)",
		"(cd tools && go mod tidy -diff)",
	} {
		if !commands[required] {
			findings = append(findings, fmt.Sprintf("scripts/check-go-mod does not run %q", required))
		}
	}
	return findings
}

func CheckReleaseNote(root string) []string {
	policy, err := ReadPolicy(filepath.Join(root, ".github", "go-support-policy.json"))
	if err != nil {
		return []string{err.Error()}
	}
	minimum, err := ParseVersion(policy.Minimum)
	if err != nil {
		return []string{err.Error()}
	}
	findings := validatePolicy(policy, minimum)

	changelog, err := os.ReadFile(filepath.Join(root, "CHANGELOG.md"))
	if err != nil {
		return append(findings, err.Error())
	}
	version := strings.TrimPrefix(policy.Release.SDKVersion, "v")
	heading := regexp.MustCompile(`(?m)^##[ \t]+` + regexp.QuoteMeta(version) + `(?:[ \t(]|$)`)
	location := heading.FindIndex(changelog)
	if location == nil {
		return append(findings, fmt.Sprintf("CHANGELOG.md has no %s release section", policy.Release.SDKVersion))
	}
	sectionEnd := len(changelog)
	if next := regexp.MustCompile(`(?m)^##[ \t]+`).FindIndex(changelog[location[1]:]); next != nil {
		sectionEnd = location[1] + next[0]
	}
	section := changelog[location[0]:sectionEnd]
	if !strings.Contains(normalizeWhitespace(string(section)), normalizeWhitespace(policy.Release.Text)) {
		findings = append(findings, fmt.Sprintf(
			"CHANGELOG.md %s section does not contain the approved Go-version release note", policy.Release.SDKVersion,
		))
	}
	return findings
}

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
