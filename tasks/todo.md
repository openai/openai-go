# Mixed-file transition task ledger

See [plan.md](plan.md) for verified baseline, complete per-file counts,
compatibility requirements, compiler seams, dependencies, and acceptance criteria.
Frozen baseline: **28 files / 715 patch lines**, SDK `267806b5`.
Statuses below describe work, not merge completion. Update links and measured
counts on Conductor state changes; do not mark a task done until merged and verified.

| Task | Scope | Status | Depends on |
|---|---|---|---|
| Audit | Hash-verified inventory and transition plan | [Ready PR #1106](https://github.com/openai/openai-go/pull/1106); independent evidence review passed | None |
| A1 | Vector-store file helper moves | [Ready PR #1104](https://github.com/openai/openai-go/pull/1104); local checks/review pass; CI/approval pending | None |
| A2 | Chat/Responses/video helper moves | [Ready PR #1105](https://github.com/openai/openai-go/pull/1105); source identity/tests/review pass; CI/approval pending | None |
| A3 | Vector-store batch helper moves | Implemented; existing focused/race/consumer checks running | None |
| B | Multipart union metadata generation | Implementing compiler change | None |
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
| G1 | Generic Execute and admin-operation auth facts | Queued | Security review |
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
