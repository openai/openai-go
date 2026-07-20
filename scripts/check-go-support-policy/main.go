package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
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

// include=all is needed to identify not only the current and previous stable
// lines, but also the just-retired line that may receive a grace period.
const defaultReleaseFeed = "https://go.dev/dl/?mode=json&include=all"

var stableVersionPattern = regexp.MustCompile(`^(?:go)?([0-9]+)\.([0-9]+)(?:\.[0-9]+)?$`)

type supportPolicy struct {
	Minimum    string  `json:"minimum"`
	GraceUntil *string `json:"grace_until"`
	Reason     string  `json:"reason"`
	ApprovedBy string  `json:"approved_by"`
}

type goRelease struct {
	Version string `json:"version"`
}

type version struct {
	Major int
	Minor int
}

func (v version) String() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

type review struct {
	GeneratedAt time.Time
	Current     version
	Previous    version
	Retired     version
	Minimum     version
	GraceUntil  *string
	Reason      string
	ApprovedBy  string
	Findings    []string
}

func (r review) Markdown() string {
	var b strings.Builder
	fmt.Fprintln(&b, "# Go version support review")
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "Generated: %s\n\n", r.GeneratedAt.UTC().Format(time.RFC3339))

	if r.Current != (version{}) {
		fmt.Fprintf(&b, "- Official current stable release: Go %s\n", r.Current)
		fmt.Fprintf(&b, "- Official previous stable release: Go %s\n", r.Previous)
		fmt.Fprintf(&b, "- Most recently retired release: Go %s\n", r.Retired)
	}
	if r.Minimum != (version{}) {
		fmt.Fprintf(&b, "- Repository minimum: Go %s\n", r.Minimum)
	}
	if r.GraceUntil != nil {
		fmt.Fprintf(&b, "- Grace period ends: %s\n", *r.GraceUntil)
	}
	if r.Reason != "" {
		fmt.Fprintf(&b, "- Recorded reason: %s\n", r.Reason)
	}
	if r.ApprovedBy != "" {
		fmt.Fprintf(&b, "- Recorded approver: %s\n", r.ApprovedBy)
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

func parseVersion(raw string) (version, error) {
	match := stableVersionPattern.FindStringSubmatch(strings.TrimSpace(raw))
	if match == nil {
		return version{}, fmt.Errorf("invalid stable Go version %q", raw)
	}
	major, err := strconv.Atoi(match[1])
	if err != nil {
		return version{}, err
	}
	minor, err := strconv.Atoi(match[2])
	if err != nil {
		return version{}, err
	}
	return version{Major: major, Minor: minor}, nil
}

func supportedLines(releases []goRelease) ([]version, error) {
	unique := make(map[version]struct{})
	for _, release := range releases {
		parsed, err := parseVersion(release.Version)
		if err != nil {
			continue // Ignore betas, release candidates, and malformed feed entries.
		}
		unique[parsed] = struct{}{}
	}

	versions := make([]version, 0, len(unique))
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

func readJSON(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

func readGoDirective(path string) (version, error) {
	file, err := os.Open(path)
	if err != nil {
		return version{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 && fields[0] == "go" {
			return parseVersion(fields[1])
		}
	}
	if err := scanner.Err(); err != nil {
		return version{}, err
	}
	return version{}, fmt.Errorf("%s has no go directive", path)
}

func fetchReleases(client *http.Client, feedURL string) ([]goRelease, error) {
	response, err := client.Get(feedURL)
	if err != nil {
		return nil, fmt.Errorf("fetch official Go release feed: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, response.Body)
		return nil, fmt.Errorf("fetch official Go release feed: %s", response.Status)
	}

	var releases []goRelease
	if err := json.NewDecoder(response.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("decode official Go release feed: %w", err)
	}
	return releases, nil
}

func checkRepository(root, feedURL string, now time.Time, client *http.Client) (review, error) {
	result := review{GeneratedAt: now}

	var policy supportPolicy
	if err := readJSON(filepath.Join(root, ".github", "go-support-policy.json"), &policy); err != nil {
		return result, err
	}
	result.GraceUntil = policy.GraceUntil
	result.Reason = policy.Reason
	result.ApprovedBy = policy.ApprovedBy

	minimum, err := parseVersion(policy.Minimum)
	if err != nil {
		return result, fmt.Errorf("parse policy minimum: %w", err)
	}
	result.Minimum = minimum

	releases, err := fetchReleases(client, feedURL)
	if err != nil {
		return result, err
	}
	lines, err := supportedLines(releases)
	if err != nil {
		return result, err
	}
	result.Current, result.Previous, result.Retired = lines[0], lines[1], lines[2]

	if strings.TrimSpace(policy.Reason) == "" {
		result.Findings = append(result.Findings, "policy file has no reason")
	}
	if strings.TrimSpace(policy.ApprovedBy) == "" {
		result.Findings = append(result.Findings, "policy file has no approver")
	}

	switch minimum {
	case result.Previous:
		if policy.GraceUntil != nil {
			result.Findings = append(result.Findings, "grace_until must be null when the minimum is the previous stable release")
		}
	case result.Retired:
		if policy.GraceUntil == nil {
			result.Findings = append(result.Findings, fmt.Sprintf("Go %s is retired and requires an approved grace_until date", minimum))
			break
		}
		graceDate, parseErr := time.Parse("2006-01-02", *policy.GraceUntil)
		if parseErr != nil {
			result.Findings = append(result.Findings, "grace_until must use YYYY-MM-DD")
			break
		}
		today := time.Date(now.UTC().Year(), now.UTC().Month(), now.UTC().Day(), 0, 0, 0, 0, time.UTC)
		if today.After(graceDate) {
			result.Findings = append(result.Findings, fmt.Sprintf("grace period for Go %s expired on %s", minimum, *policy.GraceUntil))
		}
	default:
		result.Findings = append(result.Findings, fmt.Sprintf(
			"minimum Go %s is neither the previous stable release (Go %s) nor the grace-eligible retired release (Go %s)",
			minimum, result.Previous, result.Retired,
		))
	}

	moduleFiles := []string{
		"go.mod",
		filepath.Join("examples", "go.mod"),
		// The consumer fixture is a separate module and must not quietly acquire
		// a different compiler floor.
		filepath.Join("internal", "testdata", "consumer", "go.mod"),
	}
	for _, relative := range moduleFiles {
		directive, readErr := readGoDirective(filepath.Join(root, relative))
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

	documentation := []string{"README.md", "CONTRIBUTING.md", "GO_VERSION_POLICY.md"}
	documentNeedle := "Go " + minimum.String()
	for _, relative := range documentation {
		data, readErr := os.ReadFile(filepath.Join(root, relative))
		if readErr != nil {
			result.Findings = append(result.Findings, readErr.Error())
			continue
		}
		if !strings.Contains(string(data), documentNeedle) {
			result.Findings = append(result.Findings, fmt.Sprintf("%s does not mention %q", relative, documentNeedle))
		}
	}

	ciPath := filepath.Join(root, ".github", "workflows", "ci.yml")
	ci, err := os.ReadFile(ciPath)
	if err != nil {
		result.Findings = append(result.Findings, err.Error())
		return result, nil
	}
	ciText := string(ci)
	// These strings intentionally make CI part of the policy contract. A new Go
	// release therefore creates a review issue until maintainers update the
	// tested versions rather than silently changing the support claim.
	for _, required := range []string{
		result.Previous.String() + ".x",
		result.Current.String() + ".x",
		"GOTOOLCHAIN: local",
		"./scripts/check-go-mod",
		"govulncheck",
	} {
		if !strings.Contains(ciText, required) {
			result.Findings = append(result.Findings, fmt.Sprintf(".github/workflows/ci.yml does not contain %q", required))
		}
	}

	tidyPath := filepath.Join(root, "scripts", "check-go-mod")
	tidy, err := os.ReadFile(tidyPath)
	if err != nil {
		result.Findings = append(result.Findings, err.Error())
	} else if !strings.Contains(string(tidy), "go mod tidy -diff") {
		result.Findings = append(result.Findings, "scripts/check-go-mod does not use go mod tidy -diff")
	}

	return result, nil
}

func main() {
	root := flag.String("root", ".", "repository root")
	feedURL := flag.String("feed-url", defaultReleaseFeed, "official Go JSON release feed")
	reportPath := flag.String("report", "", "optional path for the Markdown report")
	flag.Parse()

	result, err := checkRepository(*root, *feedURL, time.Now(), &http.Client{Timeout: 30 * time.Second})
	if err != nil {
		result.Findings = append(result.Findings, err.Error())
	}
	markdown := result.Markdown()
	fmt.Print(markdown)

	if *reportPath != "" {
		if writeErr := os.WriteFile(*reportPath, []byte(markdown), 0o644); writeErr != nil {
			fmt.Fprintln(os.Stderr, writeErr)
			os.Exit(1)
		}
	}
	if len(result.Findings) > 0 {
		os.Exit(1)
	}
}
