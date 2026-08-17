# Repository security, workflow, dependency, and quality policy

Use this reference for affected workflows, credentials, module graphs, Go
versions, dependency updates, generated-code ownership, and quality gates.
The root `AGENTS.md`, `GO_VERSION_POLICY.md`, `GO_CODE_QUALITY_POLICY.md`, and
the current workflow files are authoritative.

## Codex Security escalation

Use `$codex-security:security-diff-scan` for every security-sensitive PR. Keep
the scan pinned to the review's captured base/head commits and let that skill
coordinate its `$threat-model`, `$finding-discovery`, `$validation`, and
`$attack-path-analysis` phases. Use Codex Security's native plugin tools when
available; never substitute a routine general-purpose code review for the
dedicated security scan.

Review these repository-specific attack surfaces even when the changed file is
not obviously named for security:

- OpenAI, Azure, and Bedrock credentials; workload identity; token refresh;
  AWS signing; credential precedence; environment fallbacks; header redaction;
  request origins; redirects; admin-only endpoints; HTTPS proxy/client-cert
  exposure; transport-bound token caches; and cross-provider isolation.
- Webhook signature verification, timestamp tolerance, replay resistance,
  cryptographic comparisons, request-size bounds, and malformed inputs.
- URLs, multipart/file paths, JSON and SSE parsing, pagination cursors,
  concurrency/resource exhaustion, retry safety, and caller-controlled input.
- GitHub Actions triggers and permissions, OIDC, checkout credentials,
  publisher boundaries, runner privileges, dependency caches, artifacts,
  action pinning, shell injection, generated patches, and secret exposure.
- Module/dependency changes, checksums, replacements, generator ownership,
  externally supplied code, review instructions, skills, and sandbox policy.

For broad cross-boundary changes or an explicitly exhaustive security audit,
escalate to `$codex-security:deep-security-scan` with appropriately bounded
target scope. Preserve untrusted-code isolation during dynamic validation.
Report attacker-controlled sources, vulnerable sinks, affected boundaries,
realistic exploit paths, defenses or counterevidence, confidence, validation
method, incomplete coverage, and any missing plugin capability. Do not publish
findings, modify code, or execute untrusted proofs without separate authority.

## GitHub Actions trust boundaries

- Keep `GITHUB_TOKEN` permissions least-privileged. Default to `contents:
  read`, grant write scopes only to the job that needs them, and do not give
  a publisher `actions: write` or a dispatcher repository-write permissions.
- Preserve `persist-credentials: false` whenever later code should not inherit
  checkout credentials. Do not move a secret into job-wide environment or a
  shell that executes untrusted repository inputs.
- Keep third-party actions pinned to full reviewed commit SHAs. Review action
  updates, runtime migrations, OIDC permissions, artifact handoffs, cache
  keys, runner labels, shell interpolation, and event trust independently.
- Distinguish trusted default-branch events from pull-request, fork, branch,
  merge-queue, schedule, and manually dispatched events.
- Keep PR descriptions, issue titles, branch names, artifact content, and
  model output out of directly interpolated shell code.
- Do not treat a renamed, missing, skipped, or separate-event status check as
  equivalent to the repository's stable required checks.
- Check live repository and organization policy when proposing a new action,
  publishing credential, or release environment: selected-action allowlists,
  full-SHA requirements, `GITHUB_TOKEN` PR-creation restrictions, environment
  reviewer gates, required variables, and GitHub App installation/permissions
  can reject an otherwise valid-looking workflow before its first step runs.
- `GITHUB_TOKEN`-authored pushes and pull requests do not automatically fire
  normal downstream workflow events. Confirm an authorized, least-privileged
  dispatch path actually runs every required generated-branch check.
- Artifact upload excludes hidden files by default and preserves paths under
  the inputs' least common ancestor. Verify download names, directory layout,
  executable/symlink modes, and actual publisher-consumed paths.

## Monthly Go-version automation

`.github/workflows/go-version-review.yml` intentionally separates proposal,
validation/publication, and Actions dispatch:

1. Proposal runs only for the trusted repository/default branch and gets
   `OPENAI_API_KEY` from the branch-restricted `ci` environment. It has
   read-only GitHub permissions and no persisted repository credential.
2. Before exposing the model credential, create workspace-local Go cache and
   temp directories, disable `GOENV`, set `GOTOOLCHAIN=local`, disable
   setup-go caching, and warm all four module graphs. `setup-go` calls
   `go env`, so those directories must already exist.
3. A pinned Codex action/runtime removes runner sudo access and executes with
   the workspace permission profile and no command network access.
4. The output artifact is untrusted. A fresh publisher independently rejects
   unexpected paths, creates/deletes, renames, symlinks, mode changes,
   executable edits, dependency-graph edits, and workflow changes beyond
   quoted Go-version literals.
5. Complete tidy, lint, build, tests, and vulnerability scans run after
   allowlist validation but before any publishing token enters their shell.
6. Only the final publishing step receives a short-lived GitHub credential;
   it must avoid overwriting an existing generated draft.
7. A distinct dispatcher can trigger required CI, compatibility, and CodeQL
   workflows but cannot write repository contents or pull requests.

Treat cache restoration across model/publisher boundaries, token persistence,
secret exposure to untrusted code, privilege recombination, widened patch
allowlists, or executable untrusted workflow edits as consequential defects.

## Go versions and module graphs

- Keep `go.mod`, `examples/go.mod`, `internal/testdata/consumer/go.mod`, and
  `tools/go.mod` aligned when the minimum supported Go version changes.
- Verify `README.md`, `CONTRIBUTING.md`, `GO_VERSION_POLICY.md`, and the CI
  matrix describe the same minimum/current supported release lines.
- A Go-floor increase is an SDK minor-release change and requires SDK
  CODEOWNER approval plus a PR release note naming the new minimum and last
  compatible SDK release. Do not promise unsupported security backports.
- Run affected supported lines with `GOTOOLCHAIN=local`; automatic toolchain
  selection can hide an accidental minimum-version increase.
- Keep root, examples, and external-consumer dependency updates coordinated.
  Keep development tools isolated from SDK customers' module graphs.
- A separate Go package does not isolate dependency graphs; a separately
  versioned module is required when optional provider dependencies must not
  enter every SDK customer's module graph.
- For multi-directory Dependabot updates, verify grouping really spans the
  coupled modules rather than independently updating identically named groups.
- Review `replace`, `toolchain`, `tool`, direct/indirect dependency, checksum,
  and Azure/AWS dependency changes rather than focusing only on `go` lines.
- Use `./scripts/check-go-mod`; verify the external consumer with
  `go -C internal/testdata/consumer test -mod=readonly ./...`.
- Dependabot's cooldown, grouping, isolated tools updates, and independent
  GitHub Actions review exist to preserve reviewability and supply-chain
  discipline; do not silently weaken them.

## Quality ratchet

- `GO_CODE_QUALITY_POLICY.md` governs rollout order: formatting, correctness
  analyzers, error/resource handling, unused-code analysis, then selected
  style/security/complexity checks. Do not advance to a style stage before its
  documented prerequisites.
- Enabled analyzers stay enabled, thresholds only decrease, and exception
  lists only shrink unless the SDK owners explicitly approve a policy change.
- Generated and handwritten Go remain under the same checks. Broad generated
  exclusions, blanket directory exclusions, bare `nolint`, wildcard
  suppressions, and unowned baseline files are not acceptable.
- Explicitly inspect linter defaults as well as configured values: generated
  files can be excluded by default even when no exclusion is written.
- Verify that enabled checks from previously approved stacked branches
  actually survive in the proposed merge into `main`; branch ancestry and a
  green isolated stack do not prove that its cumulative changes will land.
- A rule change should cover one analyzer or coherent family, identify
  baseline/final counts and affected ownership, preserve API/wire behavior,
  and include a generator-side fix when regeneration recreates a finding.
- Full `govet` does not imply every optional upstream analyzer is appropriate.
  The documented global `fieldalignment` exclusion avoids false sharing and
  compatibility-breaking public-struct reordering; do not demand its removal
  merely because it is global.
- `errorlint` deliberately allows exact-identity EOF and redirect sentinels,
  and does not require `%w` where adding wrapping would alter public error
  contracts. Review the documented exception before proposing a change.

## Validation and compatibility

Select applicable checks based on what changed. Execute head-controlled code
or scripts only when the trust and isolation requirements in `SKILL.md`
permit it; otherwise use static review:

```sh
./scripts/lint
./scripts/check-go-mod
./scripts/test
go -C examples test ./...
go -C internal/testdata/consumer test -mod=readonly ./...
govulncheck ./...
```

Run `govulncheck` for security-relevant dependency or call-graph changes when
the pinned scanner is available. The full test script requires its Steady mock
server. `scripts/detect-breaking-changes BASE` modifies checked-out test
files, so use it only in an explicitly disposable worktree.
