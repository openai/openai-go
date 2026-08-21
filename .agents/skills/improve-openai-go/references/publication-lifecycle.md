# Publication lifecycle

Use this protocol only for a publication-capable run. Read it completely before
the run claims work. It is the canonical source for lifecycle, revision, and
external-write rules; do not duplicate them in `SKILL.md`.

This repository does not currently provide a generic implementation of this
protocol. The Go-version workflow is narrowly allowlisted and must not be reused.
Unless an approved host proves every boundary below, use audit/report-only mode.

## State machine and durable ownership

A trusted coordinator, not the model worker, owns this state machine:

`claim -> authorize -> implement -> validate -> publish-branch -> establish-pr -> dispatch -> attest -> report-status -> steward -> finalize -> release`

After a successful owned release, or after proving that a stale run cannot
release anything, emit a read-only run report from the durable terminal record.
Reporting is not part of the mutation lifecycle.

Before selecting mutable work, atomically acquire one repository-scoped lease
that excludes all other runs of this skill. On successful acquisition, return
an unforgeable acquisition receipt and durably record `lease_acquired = true`,
`run_id`, `lease_owner_id`, a unique acquisition ID, and a monotonically
increasing `fencing_epoch`. A run that did not receive that matching receipt
does not own the lease. Hold and renew the lease through finalization. Every
external writer must verify the complete receipt immediately before credential
issuance and again immediately before mutation. Lease loss must invalidate
in-flight broker authority, and a broker must reject every later mutation from
the stale epoch, including after a newer run reacquires the lease.
Each external writer must retain its credential inside the brokered operation
and keep that operation registered until the remote outcome is terminal.
Do not issue a newer epoch while an earlier remote operation is in flight or has
an unknown outcome. Enter a draining state, prevent new writes, and reconcile
each operation to a terminal result before lease reassignment.

Persist every transition in trusted durable storage outside model-writable
paths. Before an external write, record an idempotency key, repository, exact
target, expected state, and epoch; afterward record the returned object and
observed state. On restart, reconcile remote branches, pull requests, workflow
runs, check runs and commit statuses, review threads, and Slack messages before
continuing. Store check-run IDs and external IDs, or an equivalent stable
status-context lookup key, so recovery can identify the existing remote write.
Never blindly retry an operation with an unknown result. Use the remote
system's resource-specific compare-and-swap or idempotency control where
available. Otherwise serialize the brokered request, retain the lease until its
outcome is observed, and block epoch reassignment while that outcome is unknown.

If acquisition is unavailable or ambiguous, stop without external writes and
emit only a read-only acquisition-failure report; do not run a lifecycle
finalizer or release a lease this run did not acquire. If an acquired lease
expires or is lost, stop new writes and
drain or reconcile any in-flight operation. Finalize and release only when this
run still presents its matching acquisition receipt. Release must be a
compare-and-swap over `lease_acquired`, `run_id`, `lease_owner_id`, the unique
acquisition ID, and `fencing_epoch`; it clears only that exact ownership record.
A rejected release must never clear or modify a newer owner's lease. Emit the
read-only run report after the owned release succeeds, or after recording that
the stale receipt could not release anything.

## Canonical revision record

Use these names everywhere:

- `target_base_sha`: target-branch comparison base for the complete pull request.
- `target_head_sha`: authenticated target-branch tip used to construct the
  candidate integration tree.
- `source_head_sha`: trusted starting head for this run's candidate artifact.
- `expected_remote_sha`: absent for a new branch; exact current remote head for
  a service run.
- `candidate_sha`: locally committed result produced by implementation.
- `published_sha`: head returned by the publisher after its compare-and-swap.
- `integration_tree_digest`: independently measured tree produced by combining
  `target_head_sha` and the candidate or published head without executing either.

Require `source_head_sha` to be an ancestor of `candidate_sha`, its merge base to
equal `source_head_sha`, and the range to contain no merge commit unless the
authenticated work order explicitly allows one.

Bind validation to these exact ranges:

| Evidence | Exact range or revision |
| --- | --- |
| Artifact manifest, intended-path reconciliation, whitespace, focused regression | `source_head_sha...candidate_sha` |
| Full PR compatibility, security diff scan, `$openai-go-pr-review`, breaking-change input | `target_base_sha...candidate_sha`, then repeat or verify against `target_base_sha...published_sha` |
| Integration build, tests, compatibility, and security | Measured merge tree from `target_head_sha` plus `candidate_sha`, then `published_sha` |
| Hosted CI, CodeQL, external check results, review replies and resolutions, review-ready state | `published_sha` plus its authenticated `integration_tree_digest` |

Pass `target_base_sha` as the breaking-change dispatch's base input. Never call
the service-run source head the PR base.

## Split authority

Keep these components independently constrained:

- **Coordinator:** owns state and read-only hosted inspection, but has no GitHub,
  Actions, or Slack write credential. Assign hosted security scans to this
  trusted read-only stage.
- **Model worker:** receives no ambient repository, pull-request, Actions,
  Slack, cloud, SSH-agent, credential-helper, or broker credential. It may read
  and edit only its disposable workspace after authorization. Use a standalone
  `--no-local`/`--no-hardlinks` disposable clone, or keep a linked worktree and
  its entire Git common directory inside the same per-run sandbox. Reject Git
  alternates, promisor or shared object directories, cross-boundary hardlinked
  objects, external config/includes, credential helpers, and external hooks;
  verify object inodes and resolved configuration before model access. Never
  expose shared refs, config, hooks, or objects. Deny command network and every
  external mutation tool or write-network route. Require the broker to reject
  the model identity even if a brokered app is visible. Warm reviewed
  dependencies in a separate trusted step before model execution.
- **Authenticated work order:** is signed by the coordinator or kept in trusted
  immutable storage outside model-writable paths. It fixes the run ID, mode,
  revision record, epoch, allowed paths and change kinds, and validation
  contract before implementation, including finite per-command resource
  budgets. It also limits the external operation types and targets the run may
  later request.
- **Authenticated mutation envelope:** is a separately coordinator-signed,
  append-only record for each external write. It binds the allowed method,
  normalized payload digest, repository and resource target, expected head,
  idempotency key, and epoch. Generate it only after secret and sensitive-
  disclosure checks pass; the broker rejects a different payload or operation.
- **Artifact:** is a content-addressed credential-free bundle containing the
  patch, manifest, work-order identity and digest, revision record, epoch,
  allowlist, validation contract, and receipts. Model output cannot widen the
  work order. Treat the bundle as untrusted structured input, not an archive to
  extract directly into a checkout.
- **Validator/publisher:** receives no model API credential and initially no
  publishing credential. It independently retrieves the authenticated work
  order and validates the artifact in a fresh, secret-free checkout without
  persistent caches, command network, or elevated privilege. Only its final
  brokered stages may obtain separate short-lived tokens: first one scoped to
  the single branch compare-and-swap, then one scoped to establishing the
  single pull request. Each broker retains its token and synchronously performs
  only its envelope's allowlisted request; it never returns reusable credentials
  to the publisher. Both credentials' event semantics must suppress every
  implicit push- or pull-request-triggered workflow, and neither has Actions-
  dispatch permission.
- **Dispatcher:** has Actions-dispatch permission but no repository-content or
  pull-request write permission. Its broker retains a short-lived token scoped
  to one allowlisted wrapper workflow, immutable target revision, exact inputs,
  and epoch; it records the returned run ID, and the dispatcher never receives
  a reusable credential.
- **Status reporter:** has checks/status write permission but no repository,
  pull-request, workflow-dispatch, or Slack write permission. Its broker accepts
  only a validated signed receipt, `published_sha`, `target_head_sha`,
  integration-tree digest, check name, conclusion, details URL, idempotency key,
  and current epoch, then attaches that result to `published_sha` with a
  one-operation token. It stores the check-run ID and a stable external ID or
  status-context lookup key for idempotent recovery while the lease is held.
- **Slack writer:** accepts only the allowlisted channel, pull-request URL,
  target message or thread, payload, idempotency key, and current epoch. It gets
  a short-lived credential inside one brokered, tracked mutation.

## Fail-first host assertions

The host implementation must keep executable tests beside the implementation.
For this Markdown-only repository policy, preserve a base-text trace showing the
missing rule and a final-text trace satisfying it; do not invent an executable
test for infrastructure the repository does not own.

| Assertion | Reject or stop when | Required passing proof |
| --- | --- | --- |
| Fenced lifecycle | Lease is held elsewhere, the run lacks its acquisition receipt, epoch or owner is stale, ownership is lost, or an earlier write is in flight or outcome-unknown | Adversarial tests reject stale epochs at credential issuance and mutation, drain unknown operations, prohibit reassignment until terminal reconciliation, and prove that only `lease_acquired`, `lease_owner_id`, and the matching `fencing_epoch` can release the lease |
| Uncredentialed model | Model or subprocess can reach a secret, helper, agent, external mutation tool, command network, or write broker | Environment/tool inspection is clean and the broker rejects model identity |
| Authenticated allowlist | Artifact differs from the trusted work order in revisions, paths, kinds, modes, contract, or digest | Independent validation records both digests, epoch, and all accepted allowlist items before token issuance |
| Authenticated ancestry | Candidate is not the approved linear descendant | Merge-base and no-unapproved-merge assertions pass |
| Current target integration | Target tip changed or its measured merge tree was not validated | Work order restarts on target drift; candidate and published integration-tree receipts bind the current `target_head_sha` and measured digest |
| Appropriate proof | Executable behavior lacks a base-failing regression, or a non-executable artifact lacks suitable validation | Proof distinguishes base from head and matches the artifact |
| Complete intended diff | Manifest differs from staged, worktree, untracked, or committed paths, or whitespace coverage is incomplete | Exact path reconciliation, staged check, and final range check pass |
| Honest resolution | Concern is disputed, blocked, informational, unvalidated, or not fixed on the published head | Exact-head evidence exists before resolving only that thread |
| Separate dispatch | Publisher can dispatch, a mutable candidate ref selects executable workflow code, or PR-only checks did not analyze the published tree | Dispatcher records run IDs, authenticates the trusted wrapper revision, and accepts its signed receipt only when the platform-measured input tree matches `published_sha` |
| No implicit execution | Publishing a branch or pull request can trigger candidate execution before brokered dispatch | Credential event semantics and trigger filters prove publication creates no workflow run; every check starts explicitly through the trusted wrapper |
| Persisted pull request | A published new-improvement branch has no uniquely reconciled pull request, PR creation is outcome-unknown, or a retry could create a duplicate | `publish-branch` is persisted before a separately fenced `create-pr` mutation whose idempotency key and canonical result bind `published_sha`; restart tests recover branch-created/PR-missing and PR-created/result-missing states without an orphan or duplicate |
| Trusted review instructions | Candidate content can supply or modify a mandatory review skill | Review receipts identify a trusted pre-candidate skill snapshot and immutable candidate blobs |
| Isolated Git metadata | Model-writable Git state shares refs, config, hooks, objects, alternates, promisor storage, or hardlinked inodes with a trusted or reusable checkout | Filesystem and resolved-config proof confines an independent no-local/no-hardlinks clone and all Git inputs to the per-run sandbox |
| Bounded candidate execution | Any candidate-controlled command lacks an explicit finite wall-clock, CPU, memory, PID, scratch-disk, or output budget, or descendants survive its terminal state | The trusted harness enforces every signed work-order budget, kills the complete process tree on exit or breach, and rejects over-limit results |
| Target-bound status | A successful candidate check names a stale `target_head_sha` or cannot be found idempotently during recovery | Reporter rechecks the target immediately before success, binds the target in its mutation, stores the remote check identity, and relies on strict up-to-date protection or a merge queue after release |
| Published-head stewardship | The canonical remote branch or pull-request head differs from `published_sha` before a post-publication mutation or finalization | Broker-side rechecks at credential issuance and mutation bind both heads to `published_sha`; any drift stops stewardship and restarts from a new work order without resolving feedback |
| Review-ready | Required checks are absent, stale, pending, unexpectedly skipped, for another revision, or not attached to the candidate | Strict up-to-date protection or a merge queue gates integration; the reporter attaches every validated result to `published_sha`, the current target integration tree is green, and no actionable feedback remains |

## Validate and publish

Before commit, reconcile the complete intended-path manifest with staged,
unstaged, and untracked paths. Stage intended paths explicitly and run
`git diff --cached --check`. After commit, require the committed path set to
equal the manifest and run the final range checks defined above.

Before materializing model output, parse it with explicit byte, file-count, and
expansion limits. Require unique canonical repository-relative paths and reject
absolute paths, parent traversal, case or Unicode collisions, links, devices,
special files, and nested or decompression-bomb archives. Apply validated file
content to the fresh checkout without following filesystem links.

The validator/publisher must reject unapproved creates, deletes, renames,
symlinks, executable bits, workflows, dependencies, generated source, file
types or modes, and any validation-contract mismatch. It must rerun every
artifact-appropriate check required by the work order before obtaining a token.

Immediately before publication, recount skill-owned pull requests and repeat
duplicate and overlap searches. Abort if the budget or selection rules no
longer permit publication. Fetch the live target-branch tip and require it to
equal `target_head_sha`; if it moved, discard the candidate authorization and
restart with a new work order. In a trusted checkout, construct and measure the
integration tree without executing candidate code, reject merge conflicts, and
run the contract's integration checks against that exact tree. For a new branch
require the remote ref to remain absent; for service mode require it to equal
`expected_remote_sha`. Publish the branch with a non-force compare-and-swap and
abort on mismatch. Independently prove that the chosen credential's event
semantics and repository trigger filters cannot start any workflow from this
push. If publication could implicitly execute candidate content, stop instead.
Discard the branch token and durably record `published_sha`, the remote ref,
artifact digest, and compare-and-swap result before attempting a pull-request
mutation.

The `establish-pr` transition then binds one canonical pull request to that
persisted branch result. In new-improvement mode, first reconcile open and
closed pull requests by repository, base, source repository and ref,
`published_sha`, and a stable skill-owned marker. If none exists and the prior
outcome is definitively absent, perform a separately fenced `create-pr`
mutation. Its authenticated envelope, short-lived token, and deterministic
idempotency key must bind the normalized title and body, target base, source
ref, stable marker, and `published_sha`. Persist the returned node ID, number,
URL, base, source ref, head SHA, and idempotency key. In service mode, reconcile
and bind the authenticated existing pull request to the same fields instead of
creating one. On restart after branch publication, repeat reconciliation before
any retry: resume `create-pr` when creation is definitively absent, adopt the
one exact matching result when creation succeeded but its response was lost,
and stop on ambiguity. Never leave the lifecycle ready to dispatch with an
orphan branch, and never create a second pull request to recover an unknown
outcome. Prove separately that the pull-request mutation also cannot trigger
implicit candidate execution.

## Dispatch and steward the exact published head

Do not dispatch workflow code selected by the mutable candidate branch: the ref
can move between verification and dispatch. Use a reviewed wrapper workflow from
a protected trusted revision. Pass `published_sha` as an immutable input and
make the wrapper explicitly fetch and check out that SHA for analysis without
exposing secrets or write credentials to candidate code. An equivalently
protected immutable ref is acceptable. The current repository CI and CodeQL
workflows accept only mutable refs and therefore do not satisfy this protocol
without such a wrapper.

Run candidate-controlled build, test, analysis, and scripts in a fresh
disposable job with no command network, secrets, write token, persistent cache
restore or write capability, or receipt-signing authority. A trusted stage must
measure the input tree and mount source, toolchain, analyzers, and warmed
dependencies read-only; confine candidate writes to fresh scratch storage. Use
a fixed trusted harness to invoke reviewed analyzers against that immutable
snapshot. The authenticated work order must assign every candidate-controlled
command explicit finite numeric ceilings for wall-clock duration, CPU time or
quota, memory bytes, PID/process count, scratch-disk bytes and inodes, and
aggregate stdout, stderr, and result bytes. Do not inherit a platform default or
allow a candidate to raise a limit. Enforce the ceilings in the sandbox or job
runtime, terminate the complete process group and all descendants when the
command exits, times out, or exceeds any ceiling, and reclaim its scratch
storage. Treat a limit breach or incomplete descendant teardown as validation
failure: discard candidate-provided results, do not publish a pre-publication
candidate, and do not issue a successful receipt for a published one.

After the candidate job exits, have a separate trusted job or service obtain
job identity, exit status, measured input-tree digest, and bounded canonical
non-executable results from a platform-attested channel, never candidate output.
Before signing, verify that digest against the tree for `published_sha`. Bind the
trusted wrapper revision, tool and analyzer identities, measured tree digest,
current `target_head_sha`, integration-tree digest, and result into the receipt.

Verify the epoch and authenticated mutation envelope before the broker performs
each dispatch. Record the workflow-run ID, authenticate the trusted wrapper's
own `head_sha`, and require its signed receipt to bind the measured tree for
`published_sha`. A mismatch is unusable and requires a fresh lifecycle.

Immediately before every post-publication external write--including dispatch,
status reporting, review replies or resolutions, review requests, readiness
announcements, and Slack messages--and again before finalization, read the
canonical remote branch and pull-request head through the coordinator. Require
both to equal `published_sha` and the pull request to match the persisted
`establish-pr` identity. Bind those observations in the mutation envelope. The
broker must recheck them immediately before credential issuance and mutation;
if either head or identity drifts, perform no stewardship mutation, resolve no
thread, and restart from a new work order.

The separately authorized status reporter must validate the signed receipt and
re-read the target branch immediately before attaching a successful check run or
commit status to `published_sha`, not the wrapper revision. Bind
`target_head_sha`, the stored check-run ID and external ID or stable status
context, and the payload digest in the authenticated mutation. If the target
differs, do not write success. Re-read it again before lifecycle finalization and
restart with a new work order if it moved while the lease was held. Require
strict up-to-date branch protection or a merge queue for any required reporter
check; that merge gate, not a post-release automation writer, invalidates
eligibility after later target drift. Otherwise use an informational check name
that cannot satisfy the merge gate. Treat wrapper checks that exist only on the
trusted wrapper revision as execution receipts, never candidate status.

Validate each review concern against `published_sha`. After publishing a real
fix, reply with individualized exact-head evidence. Resolve only that addressed
thread. Bind each reply or resolution to repository, pull-request number,
GraphQL thread ID, operation type, normalized payload digest, `published_sha`,
idempotency key, and current epoch in its authenticated mutation envelope.
Recheck the live pull-request head as required above separately before the reply
and before the resolution; a successful reply does not authorize a resolution
after drift.
Leave disputed, blocked, informational, unvalidated, and unfixed concerns open.

Do not claim green or review-ready until every required push and dispatched
PR-only check succeeds on `published_sha` and no actionable feedback remains.
Never merge automatically.

For `#sdk-reviews`, reuse the retained thread reference or successfully search
the channel for the pull-request URL. Create exactly one root message only after
a successful search finds none, then store its thread identifier. If search is
unavailable, failed, or ambiguous, stop rather than risk a duplicate. Send all
follow-ups in that thread through the separate Slack writer.
