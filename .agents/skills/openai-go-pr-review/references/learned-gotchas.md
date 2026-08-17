# Lessons from actual OpenAI Go SDK changes

These are concrete examples, not blanket policies. Apply a lesson only when
the changed code has the same relevant preconditions.

## Credentials can leak through a different provider's defaults

[PR #698](https://github.com/openai/openai-go/pull/698) fixed Azure requests
that correctly included `Api-Key` but also inherited OpenAI bearer/admin
credentials from `OPENAI_API_KEY` or `OPENAI_ADMIN_KEY`.

The durable fix reused `option.WithHeaderDel("Authorization")`, which also
suppresses automatic credential injection. A new auth flag or a test that
checks only the Azure key misses the actual invariant. For any provider or
custom base URL, inspect ambient credentials, option precedence, redirects,
cross-origin requests, debug logs, and the complete outgoing header set.

## A regression test must fail before the fix

[PR #774](https://github.com/openai/openai-go/pull/774) fixed map-element
encoding errors. An initial regression test passed against both the old and
new implementations because it asserted only the final error. Review feedback
required asserting the marshaler call count so the test also proved that
encoding stopped immediately after the first failure.

Ask which assertion distinguishes base from head. For randomized Go map
iteration, make the distinguishing behavior independent of key order.

## Shared scaffolding can be generated without looking generated

A review on [PR #785](https://github.com/openai/openai-go/pull/785) identified
`packages/pagination/pagination.go` as Castiron-owned scaffolding even though
it does not carry an obvious generated-file header. An SDK-only cleanup would
have been overwritten on the next generation while the newly enabled linter
continued rejecting regenerated output.

Classify ownership using generator knowledge, repository history, and
generation behavior. Do not infer durable hand ownership solely from file
location, missing headers, or a clean current lint run.

## Policy order is part of correctness

The same review on [PR #785](https://github.com/openai/openai-go/pull/785)
rejected an `S1002` style rule added before the documented unused-code stage.
The change was replaced with the correctness-focused `nilnesserr` analyzer,
and downstream changes remained correctness checks.

Validate a proposed linter against both `.golangci.yml` and
`GO_CODE_QUALITY_POLICY.md`. A clean baseline does not justify skipping an
explicit rollout stage.

## Merged stacked PRs do not guarantee merged cumulative changes

[PR #782](https://github.com/openai/openai-go/pull/782) recovered approved
`ineffassign`, `govet`, and Staticcheck changes that had merged into
intermediate stacked branches but never reached `main` when their parent
merged independently.

For stacks, compare the final target branch and merge base, not just GitHub's
merged badges. Confirm that quality gates, policy updates, and prerequisite
fixes survive in the actual resulting default-branch tree.

## Broad-looking exceptions can protect real Go compatibility

[PR #782](https://github.com/openai/openai-go/pull/782) also challenged the
global `fieldalignment` exclusion. Investigation showed that upstream does not
include it in ordinary `go vet`, its advice can introduce false sharing, and
reordering exported structs can break positional literals and schema order.

Similarly, [PR #784](https://github.com/openai/openai-go/pull/784) preserved
exact-identity handling for `io.EOF` and `http.ErrUseLastResponse` and avoided
adding `%w` where public `errors.Is`/`errors.As` behavior would change.

Read the real contract and documented exception before demanding a narrower
rule. An aesthetically stricter configuration can be technically wrong.

## SSE comments are not empty data events

[PR #621](https://github.com/openai/openai-go/pull/621) fixed providers that
send SSE comments or retry-only blocks. The parser had emitted an event with
no data, causing `unexpected end of JSON input` in typed streams.

The fix belongs in the built-in decoder: blocks with no `data:` field are
ignored. An explicit empty `data:` field retains its existing invalid-JSON
behavior, and custom decoders/synthesized events keep their existing contracts.
Verify parser-layer behavior, downstream consumers, and provider-specific
transport together.

## Cleanup errors have different owners on success and failure

[PR #778](https://github.com/openai/openai-go/pull/778) removed ignored errors
across generated union accessors, multipart cleanup, generated tests, and
handwritten providers without changing exported signatures.

A successful multipart writer must propagate `Close()` failures. If encoding
already failed, best-effort `_ = writer.Close()` preserves the original error.
Generated accessors without an error result cannot expose decode failures
without changing their public API. Inspect the operation's actual error
contract instead of applying one blanket cleanup rule.

## Response bodies leak on non-obvious paths

[PR #780](https://github.com/openai/openai-go/pull/780) enabled `bodyclose` and
fixed a mock-server health check plus Bedrock middleware test responses.

Check successful health checks, unexpected middleware responses, retry
discard paths, token-refresh failures, failed response parsing, and tests—not
just the most obvious production `http.Get` call.

## A provider's endpoint family determines authentication and signing

[PR #793](https://github.com/openai/openai-go/pull/793) added Amazon Bedrock
Runtime while preserving the existing Mantle default. Runtime uses the
`bedrock` SigV4 service; Mantle uses `bedrock-mantle`.

Review endpoint inference, AWS partition suffixes, FIPS/dual-stack hosts,
HTTPS, canonical region/family matching, explicit custom proxies, credential
precedence, ambient OpenAI isolation, redirect restrictions, final middleware
ordering, and per-attempt signing as one boundary. Passing loopback tests are
not evidence that a real AWS account supports a particular model or route.

## Workflow setup can execute before the step that looks relevant

[PR #725](https://github.com/openai/openai-go/pull/725) fixed a production
failure caused by `actions/setup-go` invoking `go env` before the configured
workspace-local `GOTMPDIR` had been created.

Trace action pre/post execution, inherited job-wide environment, caches,
artifact handling, and secret exposure in actual order. A directory that is
created in the next obvious shell step can still be too late.

## PR metadata can drift away from the actual change

A review on [PR #788](https://github.com/openai/openai-go/pull/788) initially
interpreted a deliberately restacked `durationcheck` change as an intended
spelling-linter change. The implementation was correct after the title and
description were updated to reflect the new scope.

Compare the latest head, title, description, base branch, prior review
comments, and user goal. Do not report stale review context as a live code
defect after a stack has been repurposed.
