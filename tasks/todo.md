# Mixed-file transition task ledger

See [plan.md](plan.md) for verified baseline, complete per-file counts,
compatibility requirements, compiler seams, dependencies, and acceptance criteria.
Frozen baseline: **28 files / 715 patch lines**, SDK `267806b5`.
Statuses below describe work, not merge completion. Update links and measured
counts on Conductor state changes; do not mark a task done until merged and verified.

| Task | Scope | Status | Depends on |
|---|---|---|---|
| Audit | Hash-verified inventory and transition plan | [Ready PR #1106](https://github.com/openai/openai-go/pull/1106); independent evidence review passed | None |
| A1 | Vector-store file helper moves | [Ready PR #1104](https://github.com/openai/openai-go/pull/1104); local checks/review pass; CI green at `2ba8f08d`; approval pending | None |
| A2 | Chat/Responses/video helper moves | [Ready PR #1105](https://github.com/openai/openai-go/pull/1105); source identity/tests/review pass; CI green at `1d1a4184`; approval pending | None |
| A3 | Vector-store batch helper moves | [Ready PR #1107](https://github.com/openai/openai-go/pull/1107); source identity, focused/race/consumer/lint checks and independent review pass; CI green at `0a1e8398`; approval pending | None |
| B | Multipart union metadata generation | Ready compiler change; 851 renderer tests, generated Go race tests, and independent review pass; formal security scan unavailable (tool transport); awaiting required reviews/CI/merge | None |
| C1 | Constructor deprecation generation | Ready compiler change; nine focused tests and independent review pass; awaiting merge | None |
| C2 | Transcription method docs | Queued | None |
| D1 | API-key expiry representation/sample | Queued | None |
| D2 | Model alias and arbitrary schema representation | Queued | None |
| D3 | Released naming/accessors | Queued | None |
| D4 | Struct and metadata field order | Queued | None |
| D5 | Responses direct content slice | Queued | Coordinate E |
| E1 | Stable/beta message dispatch | Queued | None |
| E2 | Shell delta variant membership | Queued | None |
| F1 | Batch concurrency field/default | Queued | A3 before final SDK integration |
| F2 | Agent helper field and delegate seams | Queued | Define typed helper contract |
| G1a | Public-entrypoint admin-auth coverage for New/Get/List/Delete | [Ready test-only PR #1109](https://github.com/openai/openai-go/pull/1109); focused/race checks and four independent negative controls pass; independent review clear; security scan running; CI/approval pending | None |
| G1 | Generic Execute and admin-operation auth facts | Queued | G1a; security review |
| G2 | Webhook delegate and convenience helper/API docs | Queued | Typed verified helper contract; D4; security review |
| G3 | Image multipart null normalization | Queued | Focused wire tests; security review |
| H1 | Websocket discriminator generation | Referenced voice-generation change is merged; decoder remains absent in inspected master; investigate later history before new change | Exact upstream status |
| H2 | Routing and raw-method surface facts | Queued | Typed binding contract; security review |
| H3 | Responses websocket helper boundary | Queued | H2 |
| I1 | Budget script comment reconciliation | Queued | Shared template/lint validation |
| I2 | Generated test composite literal reconciliation | Queued | Emitter/formatter validation |
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
remove 74 patch lines: projected **27 mixed files / 641 patch lines** after all
three merge, subject to remeasurement against current main. Neither compiler
change has been promoted into the SDK; wait for its merge before regeneration.

Review-request policy: every PR and compiler Delta requires a deduplicated
`#sdk-reviews` request. The current six changes (A1, A2, A3, Audit, B, C1) share
one consolidated root request, posted 2026-10-19. Keep its internal Slack
permalink in the owning conversation's execution record; do not publish
internal communication URLs in this public repository. Add later PR requests
in that same thread and retain required CODEOWNER approval gates.

G1a (#1109) was requested for review in the existing consolidated internal
review thread. The plan authorization and admin-coverage review findings were
addressed; required PR approvals remain separate and outstanding.
