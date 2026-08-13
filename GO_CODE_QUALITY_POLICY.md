# Go Code Quality Policy

The OpenAI Go SDK improves code quality through small, independently
reviewable changes. Each change adds one formatter, analyzer, rule family, or
threshold reduction and leaves the repository buildable and testable.

This policy applies to repository-owned Go code in every module, including
generated source, handwritten source, examples, tests, internal consumer
fixtures, and development tools where the check is applicable.

## Principles

1. **Use the standard Go toolchain by default.** `gofmt` is the formatting
   authority. The compiler enforces import correctness, and generators must
   emit the exact imports their output needs. `goimports`, `gofumpt`, and other
   non-standard formatters are not required gates.
2. **Prefer correctness checks before style checks.** The rollout prioritizes
   `go vet`, `staticcheck` correctness checks, `ineffassign`, `errcheck`, and
   unused-code analysis before optional style or complexity rules.
3. **Treat generated and handwritten code consistently.** Generated code is
   not blanket-excluded from formatting or static analysis. A rule may be
   scoped when it is inapplicable to a code construct, but not merely because
   a generator emitted the file.
4. **Fix the source of generated findings.** When regeneration can recreate or
   overwrite a fix, update the generator or its shared scaffolding, add focused
   generator coverage, and then regenerate the SDK. An SDK-only patch is not a
   durable fix in that case.
5. **Ratchet in one direction.** Enabled checks remain enabled, numeric
   thresholds may only decrease, and exception lists may only shrink.
6. **Preserve compatibility.** Quality-only changes must not alter exported
   APIs, wire values, or documented behavior. A necessary compatibility change
   belongs in a separate, explicitly scoped pull request.

## Formatting

All Go source must conform to `gofmt`. Contributors should use the standard
`go fmt` command in each module that contains Go source:

```sh
go fmt ./...
go -C examples fmt ./...
go -C internal/testdata/consumer fmt ./...
```

Generated output must be `gofmt`-clean without a post-generation patch. A
developer may use `goimports` as an editor convenience, but CI and generation
do not depend on it; generated imports are fixed at their source instead.

## Static analysis

Static analysis is introduced incrementally with `golangci-lint` v2 as the
pinned runner. Its configuration uses an explicit empty default set so a tool
upgrade cannot silently enable new checks. Each analyzer or coherent rule
family is evaluated, fixed to zero untracked findings, and enabled in a
separate pull request.

The runner version used locally and in CI must match. It must also be built
with a Go release at least as new as the newest Go release analyzed by CI. Tool
updates are explicit policy changes and receive normal SDK CODEOWNER review.

The initial rollout order is:

1. `gofmt`;
2. suppression hygiene and `ineffassign`;
3. full `go vet` and `staticcheck` correctness checks;
4. `errcheck` and resource-lifecycle checks;
5. unused generated and handwritten code;
6. selected idiomatic style, security, and complexity checks.

Existing build, test, module-tidiness, supported-Go, vulnerability, and public
API checks remain authoritative throughout the rollout.

## Error handling

Every discarded error must be an intentional, documented ownership decision.

- Return actionable encoding, decoding, request, and finalization errors to the
  caller. On a successful multipart write, return any error from `Close`.
- If encoding already failed, close the multipart writer as best-effort cleanup
  and preserve the original failure. Use `_ = writer.Close()` rather than
  overwriting or obscuring the primary error.
- Close HTTP response bodies on every ownership path. If a response is already
  fully consumed or an earlier request error is being returned, explicitly
  discard a non-actionable cleanup error with `_ = body.Close()`.
- Generated union accessors without an error return cannot expose decode errors
  without changing their public API. Preserve their existing return behavior
  and make the intentional discard explicit in the Castiron template.
- In tests, fail setup, fixture writes, finalization, and required cleanup with
  `t.Fatal`, `t.Error`, or a checked `t.Cleanup` callback, as appropriate. Do
  not add an inline lint suppression when ordinary error handling is clearer.
- Fix generated models, resource methods, and generated tests in Castiron;
  correct repository-owned runtime, examples, and handwritten tests in the SDK.

## Exceptions

Exceptions are a last resort. They must:

- name the specific analyzer rather than use a bare `//nolint` directive;
- use the narrowest possible source or configuration scope;
- explain why the reported construct is correct or unavoidable; and
- identify an owner and removal condition when the exception is temporary.

Broad directory exclusions, generated-code exclusions, wildcard suppressions,
and unowned baseline files are not acceptable. Prefer clear, lint-compliant
code over a suppression whenever both express the same behavior.

## Pull request requirements

Each quality-ratchet pull request must:

- introduce one formatter, analyzer, rule family, or threshold reduction;
- record baseline and final finding counts;
- classify affected code as generated, generator-owned scaffolding, or
  SDK-owned handwritten code;
- document whether regeneration can recreate the finding;
- link the generator change when one is required;
- add focused tests for behavior-changing fixes;
- avoid broadening another exclusion or increasing another threshold; and
- pass formatting, linting, builds, tests, module-tidiness, vulnerability,
  supported-Go, and public API checks applicable to the change.

The SDK CODEOWNERS own analyzer selection, exceptions, and tool upgrades.
Generator CODEOWNERS must also review changes to templates, generated
scaffolding, or the canonical generated-code policy.
