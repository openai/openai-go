# Lessons from actual OpenAI Go SDK changes

These are concrete examples, not blanket policies. Apply a lesson only when
the changed code has the same relevant preconditions. This guide was checked
against 500 repository PR-related comments across 96 PRs on 2026-08-17: all
283 available inline review comments plus 217 substantive PR-conversation
comments. Original findings matter more than acknowledgments, status bots,
or suggestions contradicted by the current implementation.

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

## Admin keys must follow operation security, not incidental ordering

[PR #652](https://github.com/openai/openai-go/pull/652) surfaced both sides of
credential selection: generated admin-only checkpoint methods could select
ordinary bearer auth, while generic `Client.Execute` could inherit an admin
key for an ordinary request. Later request options could also fail to replace
an already populated authorization header.

Review the actual endpoint security declaration, generic execution, both keys
configured together, and client-versus-request override order. Do not infer a
single globally correct preference from an admin-only or ordinary-only test.

## A cloned TLS transport inherits hidden security-sensitive behavior

[PR #741](https://github.com/openai/openai-go/pull/741) and
[PR #792](https://github.com/openai/openai-go/pull/792) showed that cloning
`http.DefaultTransport` retains HTTPS proxy routing, `InsecureSkipVerify`,
`DialTLS`/`DialTLSContext`, `ServerName`, `GetClientCertificate`, TLS session
caches, and timeout defaults. A proxy can request the client certificate before
the intended origin, while a custom TLS dialer can bypass the installed config.

Inspect the complete certificate path, inherited callbacks and cached identity,
plaintext API endpoints, explicit-base-URL precedence, and provider mixing.
Treat README recipes as real executable security surfaces, not cosmetic prose.

## Shared tokens must be scoped to real transport identity

[PR #792](https://github.com/openai/openai-go/pull/792) found that reusing an
X.509 option across two certificate-backed transports can reuse certificate A's
token for certificate B. Wrapping a doer can instead destroy stable identity;
non-comparable values can disable caching; and a comparable type containing an
interface with a non-comparable dynamic value can still panic on comparison.

Exercise direct auth constructors, fresh method-level options, comparable and
non-comparable custom doers, certificate rotation, concurrent waiters, bounded
cache ownership, and caller-context values. A canceled refresh leader should
not incorrectly fail other live requests.

## Authentication replay must preserve middleware and body ownership

[PR #792](https://github.com/openai/openai-go/pull/792) exposed several 401
replay traps: replay below user middleware omits body compression or signing,
replay from already-mutated state duplicates non-idempotent changes, and a
replay body may leak before transport or be closed twice after `http.Client`
takes ownership. Returning `http.ErrUseLastResponse` can also turn an unfollowed
3xx into apparent success for a no-response operation.

Trace pristine request state, complete middleware order, per-attempt body
ownership, custom doers, retry budgets, nonreplayable bodies, and final status
classification together.

## EOF is not an SSE event delimiter

[PR #697](https://github.com/openai/openai-go/pull/697) proposed dispatching
buffered SSE data at EOF, but the standard requires an event-ending blank line;
an incomplete final block is discarded. Dispatching an `event:`-only block also
recreates the empty-data JSON error documented in the earlier SSE lesson.

Separately, [PR #794](https://github.com/openai/openai-go/pull/794) showed that
dropping a schema-declared shell-output event from a generated stream union
can make an object-valued delta decode into the wrong scalar field and abort
the entire stream. Review framing, typed union coverage, and downstream
accumulation as one pipeline.

## Streaming examples must survive a complete tool-call exchange

[PR #706](https://github.com/openai/openai-go/pull/706) contained several
interacting example bugs: `JustFinishedToolCall` is unsafe with parallel tool
calls, unchecked argument assertions panic, a non-nil empty tools slice emits
the rejected `"tools":[]` wire shape, and a previous response's accumulator
rejects chunks from a new response ID.

Require strict function schemas with `additionalProperties: false`, typed
argument validation, checked accumulator results, a fresh accumulator per
response, and either a complete multi-round tool loop or correctly omitted
tools on the final request. Verify executable examples, not just compilation.

## Generated API hierarchy and conditional response shapes are contracts

[PR #704](https://github.com/openai/openai-go/pull/704) exposed a service
method at a different hierarchy from the documented `ServiceAccounts.APIKeys`
surface, removed exported enum/union names inside a stable major, and modeled
`api_key` as required even when a request explicitly skips key creation.

Check generated subservice placement, old exported types and aliases, role
enum assignments, schema-dependent response presence, and shared path escaping
against a real external consumer. A superficially equivalent JSON response
does not preserve Go source compatibility.

## Release workflows depend on live repository policy

[PR #747](https://github.com/openai/openai-go/pull/747) depended on an action
outside the repository's selected-action allowlist and a `GITHUB_TOKEN` that
repository policy did not permit to create release PRs. Token-created release
branches also would not automatically trigger the required CI workflows.
[PR #754](https://github.com/openai/openai-go/pull/754) additionally found an
environment requiring manual approval and missing release-app configuration.

When release automation changes, verify action allowlists, app/token authority,
environment reviewer rules and variables, event suppression, explicit check
dispatch, and real secret boundaries; valid YAML and passing unrelated CI do
not prove that the release path can run.

## Artifact layout and module grouping are operational contracts

[PR #705](https://github.com/openai/openai-go/pull/705) found hidden artifact
files omitted by default, archive paths rooted at the least common ancestor,
multi-module Dependabot groups that did not actually span directories, and Go
patch matrices that reused stale tool-cache versions.

Trace producer-to-consumer artifact names and paths, cross-directory grouping,
external-consumer `replace` behavior, latest-patch resolution, and whether
network-denied validation actually has warmed modules, scanner data, writable
temporary directories, and the loopback access its tests require.

## Generator metadata must match the producer's actual schema

[PR #731](https://github.com/openai/openai-go/pull/731) attempted to validate
the Steady mock against a transformed-spec hash that valid `.castiron.stats.yml`
files cannot contain. The producer accepts only its documented metadata keys,
so the new guard would reject every real generated candidate.

Validate generation hashes, metadata keys, fixture paths, and mock-server
startup against the actual producer/consumer contract. Do not propose adding
fields to a strict externally owned schema without generator support.

## Similar decoder and cache paths need sibling coverage

[PR #769](https://github.com/openai/openai-go/pull/769) found that root-shape
validation disappeared when decoding through `*Struct`, allowing arrays,
scalars, and other invalid object roots. Review on
[PR #610](https://github.com/openai/openai-go/pull/610) also identified the
same embedded-`reflect.Type` cache-key pattern in adjacent JSON, query, and
form encoders when only one decoder path was fixed.

For structural decoder or reflection changes, exercise pointer and value roots,
all JSON shape families, neighboring encoder/decoder caches, and actual
downstream build or linker behavior before declaring one local fix complete.
