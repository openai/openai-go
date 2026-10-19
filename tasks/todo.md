# Mixed-file transition task ledger

See [plan.md](plan.md) for verified baseline, complete per-file counts,
compatibility requirements, compiler seams, dependencies, and acceptance criteria.
Frozen baseline: **28 files / 715 patch lines**, SDK `267806b5`.
Latest measured source: **27 files / 641 patch lines**, SDK `f85ad491` (2026-10-19).
Current main `85be4f59` adds only the merged plan; the source inventory is unchanged.
Statuses below describe work, not merge completion. Update links and measured
counts on Conductor state changes; do not mark a task done until merged and verified.

| Task | Scope | Status | Depends on |
|---|---|---|---|
| Audit | Hash-verified inventory and transition plan | [Merged PR #1106](https://github.com/openai/openai-go/pull/1106) at `85be4f59921fc04b44beb2c767ccfde5fad89bde`; post-merge push workflows separately confirmed green/complete | None |
| A1 | Vector-store file helper moves | [Merged PR #1104](https://github.com/openai/openai-go/pull/1104) at `4520c2dac8cbb2b0a6f94d559d34b10bb83019fa`; post-merge workflows green; local/independent verification retained | None |
| A2 | Chat/Responses/video helper moves | [Merged PR #1105](https://github.com/openai/openai-go/pull/1105) at `d1a720d84da12a6855a42e9e9ae9cf5f034a541f`; post-merge workflows green; local/independent verification retained | None |
| A3 | Vector-store batch helper moves | [Merged PR #1107](https://github.com/openai/openai-go/pull/1107) at `b5e9da2a994a8e2007863d906bf3d2642a4f4a37`; post-merge workflows green; local/independent verification retained | None |
| B | Multipart union metadata generation | Ready compiler change; stale CLI assertions repaired; 5 CLI multipart and 851 renderer tests plus negative controls pass; independent review and repeated exact-head six-file security scan complete; awaiting required reviews/CI/merge | None |
| C1 | Constructor deprecation generation | Ready compiler change; nine focused tests and independent review pass; awaiting merge | None |
| C2 | Transcription method docs | Queued | None |
| D1 | API-key expiry representation/sample | Queued | None |
| D2 | Model alias and arbitrary schema representation | Queued | None |
| D3 | Released naming/accessors | Queued | None |
| D4 | Struct and metadata field order | Queued | None |
| D5 | Responses direct content slice | Queued | Coordinate E |
| E1 | Stable/beta message dispatch | Queued | None |
| E2 | Shell delta variant membership | Queued | None |
| F1 | Batch concurrency field/default | Ready compiler change with authoritative opt-in; 853 renderer tests, config validation, emitted Go consumer/helper tests, negative controls and independent/security reviews pass | A3 merged; compiler must land before SDK integration |
| F2 | Agent helper field and delegate seams | Queued | Define typed helper contract |
| G1a | Public-entrypoint admin-auth coverage for New/Get/List/Delete | [Merged PR #1109](https://github.com/openai/openai-go/pull/1109) at `f85ad491c1bc0c3d08d5f867262a06a53548f2ab`; post-merge workflows green; local/independent verification retained | None |
| G1 | Generic Execute auth preference | Ready independent compiler change with authoritative opt-in; Full renderer suite, source-auth probes, config/Go/provider/negative checks and independent/security reviews pass | Compiler merge before SDK integration |
| G1b | Checkpoint-permission Go-only auth binding | Implementing typed method fact; shared auth input approach rejected because it breaks current Python/Ruby credential options | G1a merged; other language/parity output unchanged; security review |
| G2 | Webhook delegate and convenience helper/API docs | Queued | Typed verified helper contract; D4; security review |
| G3 | Image multipart null normalization | Queued | Focused wire tests; security review |
| H1 | Websocket discriminator generation | Referenced voice-generation change is merged; decoder remains absent in inspected master; investigate later history before new change | Exact upstream status |
| H2 | Routing and raw-method surface facts | Queued | Typed binding contract; security review |
| H3 | Responses websocket helper boundary | Queued | H2 |
| I1 | Budget script comment reconciliation | [Ready PR #1110](https://github.com/openai/openai-go/pull/1110); exact snapshot restoration, AST/compiled-code identity, Ruff, 95 offline tests (1 existing skip), independent review; candidate 26 files / 639 lines | Required CI/reviews/merge; no budget change |
| I2 | Generated test composite literal reconciliation | Snapshot restoration rejected by current formatter; investigating current compiler sample path | Preserve formatter and immutable fixtures |
| Integration | Regenerate each eligible compiler slice | Wait for corresponding compiler merge | Merged compiler and all compatibility gates |
| Finish | Zero patch, repeat generation, policy-only ratchet | Queued | All tracks merged and verified |

- [x] Freeze exact public baseline and validate generated hash/ownership.
- [x] Identify independent work and compatibility-sensitive compiler changes.
- [ ] Merge initial independent SDK and compiler PRs with required reviews/checks.
- [ ] Measure each regenerated slice against its newly verified generated snapshot.
- [ ] Demonstrate no developer-facing breaking changes across all slices.
- [ ] Reach zero mixed files/patch lines and reproduce with a second generation.
- [ ] Propose any lower budget separately, retaining required protections.

Initial A1/A2/A3 changes are independent PRs on the same baseline. Together they
removed 74 patch lines: verified **27 mixed files / 641 patch lines** on main
`f85ad491` against its own hash-verified snapshot. Neither compiler
change has been promoted into the SDK; wait for its merge before regeneration.

Review-request policy: every PR and compiler Delta requires a deduplicated
root-channel `#sdk-reviews` request explicitly including its URL. Thread-only
requests do not qualify. The current six changes (A1, A2, A3, Audit, B, C1) share
one consolidated root request, posted 2026-10-19. Keep its internal Slack
permalink in the owning conversation's execution record; do not publish
internal communication URLs in this public repository. Add later changes to a
root request that explicitly includes their URLs, or create a new root request
if none exists; use threads for status follow-ups. Retain required CODEOWNER
approval gates.

G1a (#1109) now has a separate root-channel review request, posted 2026-10-19,
after deduplication confirmed its earlier mentions were thread replies only.
The plan authorization and admin-coverage review findings were addressed;
required approvals and merge completion remain separate states.

On 2026-10-19 at 11:53 PDT, Conductor verified fresh exact-head/feedback/gate
checks and successful protected merge queue enrollment for #1104, #1105, #1107,
and #1109, in that order (positions 1-4 at observation). All four were approved,
CI-green, and had completed clean automated reviews at their recorded heads.
This records enrollment only: do not mark these tasks merged, reduce the measured
main baseline, or advance merge-dependent work until actual merge confirmation.
Wait for Conductor updates rather than polling or duplicating queue enrollment.

On 2026-10-19 at 12:02-12:03 PDT, Conductor confirmed all four actual merges
listed above and successful post-merge workflows. The earlier queue observation
is historical, not the current task status. A fresh trusted-report measurement
confirms 424 additions + 217 deletions = 641 residual patch lines across 27 files.

Plan PR #1106's OkTest SDK test job passed, but its publisher correctly rejected
an old synthetic merge when main advanced during the run. The tested merge was
`536c154c3d84229dde6ff7a056da3edb494227f0`, based on `132dcbe3`; the
publisher re-resolved a newer main. Refreshing the documentation branch onto
current main and rerunning a fresh candidate preserves that provenance gate.
The refreshed PR still changes only the two task Markdown files relative to
main; no SDK production changes are introduced by this repair. Approval was
recorded on head `4981429c`; fresh-head validation and normal approval gates
remain authoritative.

On 2026-10-19 at 12:23 PDT, Conductor confirmed #1106 actually merged at
`85be4f59921fc04b44beb2c767ccfde5fad89bde`. Its refreshed head had green
CI, renewed exact-head approval, clean reviews and resolved threads before
protected queue completion. At 12:26 PDT, Conductor separately confirmed all
post-merge push workflows succeeded; that watch is green/complete.

Second-wave publication records (internal compiler links stay in the owning
conversation): G1 generic Execute and F1 batch state are separate ready compiler
changes; PR #1110 is the separate one-comment restoration. Each has a root
`#sdk-reviews` request explicitly including its URL. These are not merged source
reductions. Compiler output must not be promoted before the relevant change lands.

G1b's shared-security proposal was never published. Independent review proved
that current Python/Ruby `api_key` callers would break under admin-only shared
metadata. Keep canonical security and derived monitoring filters unchanged;
implement the existing Go selection through a closed, validated, Go-only method
binding and verify other SDK/profile output remains unchanged.
