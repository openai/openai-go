# Publication lifecycle

Use this protocol only for a publication-capable run. Read it completely before
the run claims work. It is the canonical source for lifecycle, revision, and
external-write rules; do not duplicate them in `SKILL.md`.

This repository does not currently provide a generic implementation of this
protocol. The Go-version workflow is narrowly allowlisted and must not be reused.
Unless an approved host proves every boundary below, use audit/report-only mode.

## State machine and durable ownership

A trusted coordinator, not the model worker, owns this state machine:

`claim -> authorize -> implement -> validate -> publish -> dispatch -> steward -> finalize -> release`

After release, emit a read-only run report from the finalized durable record.
Reporting is not part of the mutation lifecycle.

Before selecting mutable work, atomically acquire one repository-scoped lease
that excludes all other runs of this skill. Record a run ID and issue a
monotonically increasing fencing epoch. Hold and renew the lease through
finalization. Every external writer must verify the epoch immediately before
credential issuance and again immediately before mutation. Lease loss must
invalidate in-flight broker authority, and a broker must reject every later
mutation from the stale epoch, including after a newer run reacquires the lease.

Persist every transition in trusted durable storage outside model-writable
paths. Before an external write, record an idempotency key, repository, exact
target, expected state, and epoch; afterward record the returned object and
observed state. On restart, reconcile remote branches, pull requests, workflow
runs, review threads, and Slack messages before continuing. Never blindly retry
an operation with an unknown result.

If the lease is unavailable, ambiguous, expired, or lost, stop without writing.
Finalize a terminal record, release the lease, and then report the result.

## Canonical revision record

Use these names everywhere:

- `target_base_sha`: target-branch comparison base for the complete pull request.
- `source_head_sha`: trusted starting head for this run's candidate artifact.
- `expected_remote_sha`: absent for a new branch; exact current remote head for
  a service run.
- `candidate_sha`: locally committed result produced by implementation.
- `published_sha`: head returned by the publisher after its compare-and-swap.

Require `source_head_sha` to be an ancestor of `candidate_sha`, its merge base to
equal `source_head_sha`, and the range to contain no merge commit unless the
authenticated work order explicitly allows one.

Bind validation to these exact ranges:

| Evidence | Exact range or revision |
| --- | --- |
| Artifact manifest, intended-path reconciliation, whitespace, focused regression | `source_head_sha...candidate_sha` |
| Full PR compatibility, security diff scan, `$openai-go-pr-review`, breaking-change input | `target_base_sha...candidate_sha`, then repeat or verify against `target_base_sha...published_sha` |
| Hosted CI, CodeQL, review replies and resolutions, review-ready state | `published_sha` only |

Pass `target_base_sha` as the breaking-change dispatch's base input. Never call
the service-run source head the PR base.

## Split authority

Keep these components independently constrained:

- **Coordinator:** owns state and read-only hosted inspection, but has no GitHub,
  Actions, or Slack write credential. Assign hosted security scans to this
  trusted read-only stage.
- **Model worker:** receives no ambient repository, pull-request, Actions,
  Slack, cloud, SSH-agent, credential-helper, or broker credential. It may read
  and edit only its disposable workspace after authorization. Deny command
  network and every external mutation tool or write-network route. Require the
  broker to reject the model identity even if a brokered app is visible. Warm
  reviewed dependencies in a separate trusted step before model execution.
- **Authenticated work order:** is signed by the coordinator or kept in trusted
  immutable storage outside model-writable paths. It fixes the run ID, mode,
  revision record, epoch, allowed paths and change kinds, and validation
  contract before implementation.
- **Artifact:** is a content-addressed credential-free bundle containing the
  patch, manifest, work-order identity and digest, revision record, epoch,
  allowlist, validation contract, and receipts. Model output cannot widen the
  work order.
- **Validator/publisher:** receives no model API credential and initially no
  publishing credential. It independently retrieves the authenticated work
  order and validates the artifact in a fresh, secret-free checkout without
  persistent caches, command network, or elevated privilege. Only its final
  step may receive a short-lived token scoped to one branch and pull request.
  It has no Actions-dispatch permission.
- **Dispatcher:** has Actions-dispatch permission but no repository-content or
  pull-request write permission.
- **Slack writer:** accepts only the allowlisted channel, pull-request URL,
  target message or thread, payload, idempotency key, and current epoch. It gets
  a short-lived credential for one verified mutation.

## Fail-first host assertions

The host implementation must keep executable tests beside the implementation.
For this Markdown-only repository policy, preserve a base-text trace showing the
missing rule and a final-text trace satisfying it; do not invent an executable
test for infrastructure the repository does not own.

| Assertion | Reject or stop when | Required passing proof |
| --- | --- | --- |
| Fenced lifecycle | Lease is held elsewhere, epoch is stale, or ownership is lost | Adversarial tests reject stale epochs at credential issuance and mutation, including after reacquisition |
| Uncredentialed model | Model or subprocess can reach a secret, helper, agent, external mutation tool, command network, or write broker | Environment/tool inspection is clean and the broker rejects model identity |
| Authenticated allowlist | Artifact differs from the trusted work order in revisions, paths, kinds, modes, contract, or digest | Independent validation records both digests, epoch, and all accepted allowlist items before token issuance |
| Authenticated ancestry | Candidate is not the approved linear descendant | Merge-base and no-unapproved-merge assertions pass |
| Appropriate proof | Executable behavior lacks a base-failing regression, or a non-executable artifact lacks suitable validation | Proof distinguishes base from head and matches the artifact |
| Complete intended diff | Manifest differs from staged, worktree, untracked, or committed paths, or whitespace coverage is incomplete | Exact path reconciliation, staged check, and final range check pass |
| Honest resolution | Concern is disputed, blocked, informational, unvalidated, or not fixed on the published head | Exact-head evidence exists before resolving only that thread |
| Separate dispatch | Publisher can dispatch, ref moved, or PR-only checks did not start for the published revision | Dispatcher records run IDs and accepts only runs whose `head_sha` equals `published_sha` |
| Review-ready | Required checks are absent, stale, pending, unexpectedly skipped, or for another revision | All required push and PR-only checks are green on `published_sha` with no unresolved actionable feedback |

## Validate and publish

Before commit, reconcile the complete intended-path manifest with staged,
unstaged, and untracked paths. Stage intended paths explicitly and run
`git diff --cached --check`. After commit, require the committed path set to
equal the manifest and run the final range checks defined above.

The validator/publisher must reject unapproved creates, deletes, renames,
symlinks, executable bits, workflows, dependencies, generated source, file
types or modes, and any validation-contract mismatch. It must rerun every
artifact-appropriate check required by the work order before obtaining a token.

Immediately before publication, recount skill-owned pull requests and repeat
duplicate and overlap searches. Abort if the budget or selection rules no
longer permit publication. For a new branch require the remote ref to remain
absent; for service mode require it to equal `expected_remote_sha`. Publish with
a non-force compare-and-swap and abort on mismatch. Discard the token and record
`published_sha`, artifact digest, and compare-and-swap result.

## Dispatch and steward the exact published head

GitHub workflow dispatch uses a mutable ref. Immediately before each dispatch,
verify the epoch and atomically confirm the ref still equals `published_sha`.
Dispatch every required PR-only workflow with the separate dispatcher and
record each workflow-run ID. Accept a run only if its recorded `head_sha`
equals `published_sha`; a moved ref or mismatch requires a fresh lifecycle.

Validate each review concern against `published_sha`. After publishing a real
fix, reply with individualized exact-head evidence. Resolve only that addressed
thread. Bind each reply or resolution to repository, pull-request number,
GraphQL thread ID, `published_sha`, and current epoch. Leave disputed, blocked,
informational, unvalidated, and unfixed concerns open.

Do not claim green or review-ready until every required push and dispatched
PR-only check succeeds on `published_sha` and no actionable feedback remains.
Never merge automatically.

For `#sdk-reviews`, reuse the retained thread reference or successfully search
the channel for the pull-request URL. Create exactly one root message only after
a successful search finds none, then store its thread identifier. If search is
unavailable, failed, or ambiguous, stop rather than risk a duplicate. Send all
follow-ups in that thread through the separate Slack writer.
