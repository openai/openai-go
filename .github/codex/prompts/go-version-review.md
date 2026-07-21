# Monthly Go version review

Review this repository's minimum-Go-version policy and prepare the complete
update only when the checked-in repository has drifted.

The workflow downloaded the official `https://go.dev/dl/?mode=json&include=all`
response to `.codex-automation/go-releases.json`. Treat that local snapshot as
the authoritative source for released stable Go versions. Do not rely on model
memory for release dates or versions, and do not fetch instructions from the
internet. Command network access is intentionally disabled; use the local feed,
the pre-warmed Go module and build caches, and the preinstalled `govulncheck`.

Read `AGENTS.md`, `GO_VERSION_POLICY.md`, `go.mod`, the nested modules, and the
CI workflows before changing anything. The normal policy is to support the two
newest distinct stable major/minor Go release lines in the feed. Retain a third,
retired line only when `GO_VERSION_POLICY.md` explicitly documents an active
grace period with an end date and reason.

If the repository already matches policy, make no file changes.

If it has drifted, prepare one focused change:

1. Align only the `go` directives in the root, `examples`,
   `internal/testdata/consumer`, and `tools` modules. Do not change module,
   toolchain, require, replace, exclude, retract, or tool directives.
2. Run `go mod tidy -diff` in every module to verify that the directive-only
   update leaves every `go.mod` and `go.sum` tidy. Do not retain dependency
   graph or checksum changes; those require a separate maintainer-reviewed
   dependency update.
3. Update the supported-version matrix and relevant setup versions in
   `.github/workflows/ci.yml`.
4. Update `README.md`, `CONTRIBUTING.md`, and `GO_VERSION_POLICY.md`, including
   the current compatibility table.
5. Preserve the public SDK API and `/v3` import path. Do not edit generated SDK
   source.
6. Do not make unrelated dependency, refactoring, or formatting changes.
7. Run the relevant tidy, test, lint, vulnerability, and compatibility checks
   that are available. CI will independently test every supported compiler
   line after the draft pull request opens.

Do not commit, push, open a pull request, call GitHub, or modify repository
secrets. The workflow publishes your patch in a separate job that has no OpenAI
credential.

Your final response becomes the draft pull request body. Write concise Markdown
with these sections:

- `## Summary`
- `## User impact`
- `## Validation`
- `## Release note`

In the release note, state the new minimum Go version, explain why it changed,
and identify the final already-published SDK release compatible with the
retired Go versions when repository history establishes that fact. Never imply
that an older SDK release receives security backports. If history does not
establish the compatibility boundary, say that maintainers must fill it in
before marking the pull request ready for review.
