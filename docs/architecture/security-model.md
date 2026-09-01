# Security Model

This is the repository's canonical detailed threat model for Codex Security
scans and security review. Read it from the scanned revision.
[SECURITY.md](../../SECURITY.md) remains the disclosure and reportability entry
point and does not independently redefine trust boundaries.

### Authority and inventory

At adoption, the repository contains one root SECURITY.md, no nested
SECURITY.md files, and no other checked-in threat-model document. This document
alone governs the repository's detailed threat model, trust boundaries, and
severity calibration. SECURITY.md governs private disclosure and reportability;
AGENTS.md and repository skills govern review procedure; workflow files are
enforcement evidence rather than competing threat-model authority. The Codex
Security cloud scan field is external configuration, not repository guidance,
and should contain only a short pointer to this file once the scanned revision
includes it.

## 1. Overview

openai-go is the official Go client SDK for OpenAI APIs. It is a library
embedded into caller-owned services, workers, CLIs, and applications—not an API
server, sandbox, or local command-execution engine. Generated API services
delegate request construction, authentication, retries, transport, response
decoding, and streaming to shared runtime packages. Optional Azure, Amazon
Bedrock, workload-identity, mTLS, and webhook helpers add security-relevant
boundaries. (README.md:1-13, client.go:19-62, client.go:103-178)

The principal assets are OpenAI, organization-admin, Azure, AWS, and exchanged
bearer credentials; workload subject tokens and TLS private keys; webhook
secrets; customer prompts, files, responses, and raw diagnostics; tenant and
project routing; release credentials; official SDK artifacts; availability; and
API spend.

| Component | Role and security relevance | Evidence |
| --- | --- | --- |
| client.go, option/, internal/requestconfig/ | Builds credentialed HTTP requests, selects endpoints and credentials, applies middleware and retries, validates request origins, and decodes responses. | client.go:65-178; option/requestoption.go:26-61; internal/requestconfig/requestconfig.go:48-61; internal/requestconfig/origin.go:46-140 |
| Generated services and helpers | Expose typed OpenAI API operations, including administrative, uploads, files, Responses, Realtime, and tool descriptions; they do not execute model-returned tools locally. | client.go:24-62; client.go:175-210; README.md:163-219 |
| azure/, bedrock/ | Select provider-specific endpoints and authentication while preventing ambient OpenAI credentials from being sent to another provider. | azure/azure.go:58-76, azure/azure.go:164-184; bedrock/bedrock.go:19-110; bedrock/auth.go:97-114 |
| auth/ | Obtains and exchanges workload identity, manages short-lived bearers, and enforces X.509 transport constraints. | auth/subjecttokenprovider.go:1-113; auth/x509transport.go:1-120; option/x509workloadidentity.go:18-55 |
| webhooks/ | Authenticates webhook bytes before callers decode or act on events. | webhooks/webhook.go:1-44; webhooks/webhook_helpers.go:1-143; webhooks/webhook_secret.go:1-34 |
| internal/apijson/, internal/apiform/, packages/ssestream/, streamaccumulator* | Parses attacker-influenced API data, serializes caller input, and accumulates streaming output. | internal/apijson/decoder.go:1-120; internal/apiform/encoder.go:1-180; packages/ssestream/ssestream.go:1-170; streamaccumulator.go:1-120 |
| .github/workflows/, scripts/, tools/ | Builds, tests, scans, and publishes the repository; privileged publication must remain isolated from pull-request-controlled execution. | .github/workflows/ci.yml:1-183; .github/workflows/codeql.yml:1-72; .github/workflows/castiron-custom-code-comment.yml:1-237; .github/workflows/create-releases.yml:1-55; .github/workflows/go-version-review.yml:200-411; scripts/castiron/CUSTOM_CODE.md:24-78 |

| Deployment or workflow | Resource or capability | Configuration and precedence | Safe effective value or location | Readers, writers, or recipients | Enforcing control | Evidence or unknowns |
| --- | --- | --- | --- | --- | --- | --- |
| Ordinary SDK client | OpenAI API key, admin key, organization and project IDs, webhook secret | Explicit request or client options layer over environment-derived defaults | Caller process memory and outbound headers to the configured OpenAI origin | Caller process; configured OpenAI endpoint | Option layering, auth selection, and origin validation | client.go:65-100; option/requestoption.go:299-337; internal/requestconfig/option_layers.go:229-306; internal/requestconfig/origin.go:46-140 |
| Custom endpoint or HTTP client | Credentialed request destination and redirect behavior | WithBaseURL and WithHTTPClient are operator or developer escape hatches | Caller-selected endpoint; bespoke Do implementations own redirects | Configured endpoint or custom transport | Native http.Client path retains SDK redirect checks; custom implementations must enforce their own | option/requestoption.go:26-61; internal/requestconfig/requestconfig.go:632-660 |
| Azure or Bedrock client | Provider credentials and endpoint | Provider options override or suppress ambient OpenAI defaults | Configured Azure or Bedrock origin; authenticated redirects constrained | Azure or AWS endpoint selected by operator | Provider finalizers, HTTPS and origin checks, per-attempt authentication or signing | azure/azure.go:164-184, azure/azure.go:286-307; bedrock/auth.go:97-114, bedrock/auth.go:462-577 |
| Workload identity and X.509 | Subject tokens, exchanged bearer tokens, client certificates and private keys | Operator-selected token source, issuer, signer, and transport | Fixed metadata or exchange endpoints and attested OpenAI mTLS origins | Cloud metadata service, issuer, OpenAI endpoint | Issuer and origin checks, TLS constraints, bounded responses, refresh bookkeeping | auth/subjecttokenprovider.go:1-113; auth/x509transport.go:1-120; option/x509workloadidentity.go:18-55 |
| Webhook verification | Webhook secret and raw request bytes | Caller supplies secret and exact body; default tolerance may be overridden | In-process verification before event decoding | Caller process only until verification succeeds | HMAC verification, timestamp window, canonical secret parsing | webhooks/webhook_helpers.go:1-143; webhooks/webhook_secret.go:1-34 |
| Pull request CI and release | Read-only test execution versus protected publication credentials | Workflow event, job permissions, and protected environments | PR jobs without release credentials; release job only on main in the release environment | GitHub Actions jobs, release app, package consumers | Pinned actions, minimal permissions, trusted recomputation and publisher separation | .github/workflows/ci.yml:1-183; .github/workflows/codeql.yml:1-72; .github/workflows/castiron-custom-code-comment.yml:1-237; .github/workflows/create-releases.yml:1-55; .github/workflows/go-version-review.yml:200-411; scripts/castiron/CUSTOM_CODE.md:24-78 |

## 2. Threat Model, Trust Boundaries, and Assumptions

### Actors, authority, and protected objectives

- Application maintainer or operator: chooses SDK options, credentials,
  endpoints, proxies, trust roots, middleware, custom transports, workload
  providers, and downstream handling. These are trusted configuration or
  in-process code, not ordinary remote attacker input.
- Application user: may control prompts, uploaded bytes and filenames, resource
  IDs, model or tool arguments, and selected typed fields only to the extent the
  embedding application exposes them. The application—not this SDK—must
  authorize users and tenant or resource access.
- API, provider, proxy, or network peer: can return HTTP status, headers, JSON,
  SSE frames, cursors, retry hints, and errors. Even across authenticated TLS,
  protocol data remains untrusted runtime input to parsers and application
  consumers.
- Webhook sender: controls headers and body bytes until signature and timestamp
  verification succeed.
- Repository contributor: controls candidate pull-request content. Reviewed,
  checked-in repository source—including examples, tests, fixtures, build
  scripts, generators, and other executable checkout files—runs with
  repository-code authority when the repository's own workflows intentionally
  execute it. A contributor who can change tracked executable code does not
  acquire a new privilege merely because candidate-safe CI intentionally runs
  that code. This is vulnerability and severity classification, not permission
  to execute an untrusted PR checkout during review: developer commands apply
  only to trusted revisions or checkouts. Treat it as a security boundary only
  when independently mutable lower-trust data crosses a parser or evaluator
  boundary, when untrusted runtime, API, or network data reaches a sensitive
  sink, or when candidate PR code or artifacts can reach protected CI or release
  credentials.
- GitHub Actions or release operator: owns workflow permissions, protected
  environments, release credentials, and publication of official artifacts.

The objectives are to keep credentials on intended origins and out of
diagnostics, preserve authentication and provider separation, verify webhooks
before trusted use, parse hostile protocol data without unsafe effects, preserve
tenant and project routing semantics, and prevent candidate PR content from
influencing trusted publication or protected credentials.

### Material trust boundaries and invariants

1. Application input to request construction. Path parameters must remain data
   rather than alternate authorities or URLs; typed serialization must preserve
   intended wire shapes; caller-selected raw or extra-field escape hatches are
   trusted-input APIs, not sanitizers. (internal/requestconfig/requestconfig.go:48-61,
   option/requestoption.go:197-244, packages/param/param.go:25-55,
   packages/param/param.go:141-159)
2. Credentialed SDK to configured endpoint. Initial requests and native
   redirects must retain the configured origin; credentials for OpenAI, Azure,
   and Bedrock must not bleed into another provider; admin-only operations must
   not silently use incompatible workload credentials.
   (internal/requestconfig/origin.go:46-140, azure/azure.go:286-307,
   bedrock/auth.go:97-114, bedrock/auth.go:577-580,
   option/requestoption.go:355-364)
3. Webhook network bytes to trusted event. Verify the exact body, signature, and
   timestamp before decoding or acting. Callers still own request-size limits,
   deduplication, tenant mapping, and business authorization.
   (webhooks/webhook_helpers.go:1-143, README.md:430-520)
4. Provider or API response to parser and caller. JSON, SSE, pagination cursors,
   headers, and error bodies are untrusted data. Parsing must not create code
   execution, credential disclosure, unsafe routing, or disproportionate
   resource exhaustion; applications must treat model output and tool arguments
   as data until they validate a downstream sink. (packages/ssestream/ssestream.go:1-170,
   internal/apijson/decoder.go:1-120, README.md:163-219)
5. Operator-selected identity material to credential exchange. Workload token
   paths, metadata services, issuers, signers, trust roots, and X.509 transports
   are privileged configuration. The SDK must enforce its documented endpoint,
   redirect, TLS, bearer-lifetime, and refresh constraints.
   (auth/subjecttokenprovider.go:1-113, auth/x509transport.go:1-120,
   option/x509workloadidentity.go:18-55)
6. PR candidate to protected CI or release authority. Candidate-controlled
   source, filenames, artifacts, reports, and patches must not be evaluated by
   a job holding write or release credentials unless independently constrained
   and validated by trusted code. (.github/workflows/codeql.yml:1-72,
   .github/workflows/castiron-custom-code-comment.yml:1-237,
   .github/workflows/create-releases.yml:1-55,
   .github/workflows/go-version-review.yml:200-411,
   scripts/castiron/CUSTOM_CODE.md:24-78)

### Assumptions, exclusions, and caller obligations

- The SDK runs in the caller's process with the caller's privileges. Malicious
  in-process middleware, arbitrary custom Do implementations, attacker-selected
  base URLs or token paths, disabled TLS verification, and concurrent mutation
  contrary to ownership contracts already grant authority beyond ordinary SDK
  inputs.
- OpenAI, Azure, and AWS enforce server-side authentication, authorization,
  ownership, retention, and revocation. A shared SDK credential is not an
  end-user authorization boundary.
- The repository's tracked executable files are reviewed repository code, not a
  separate lower-trust data channel merely because a PR author can edit them.
  This does not exempt dependency provenance, generated artifacts, PR metadata,
  network-fetched content, or protected publisher boundaries from review.
  Candidate PR checkouts remain operationally untrusted until reviewed and
  merged; this classification does not authorize local execution contrary to
  the repository's trusted-checkout review procedure.
- The SDK does not sanitize model output for HTML, SQL, or shell sinks, execute
  returned tool calls, provide application quotas, deduplicate webhooks, or
  authorize application users.
- Large JSON bodies and SSE events are normal supported API output. Existing
  compatibility must be preserved; new arbitrary limits are not a security fix
  without an explicit API contract.

## 3. Attack Surface, Mitigations, and Attacker Stories

The scenarios below are review hypotheses and severity calibration aids, not
confirmed findings.

| Priority | Scenario and capability gain | Prerequisites | Impact | Existing controls | Mitigation or review focus | Evidence |
| --- | --- | --- | --- | --- | --- | --- |
| High | A path, URL, redirect, or option-layering bug sends a credentialed request to an unintended origin or with the wrong credential. | Attacker controls a normally exposed request field or a hostile endpoint response; not merely trusted WithBaseURL configuration. | Usable credential disclosure or cross-tenant or admin action. | Segment escaping, initial-origin and redirect checks, provider-specific guards, credential-layer tracking. | Test encoded separators, authority mutations, redirects, pagination, retries, and provider mixing. | internal/requestconfig/requestconfig.go:48-61; internal/requestconfig/origin.go:46-140; azure/azure.go:286-307; bedrock/auth.go:97-114 |
| High | Forged or replayed webhook input is accepted as a trusted event. | Attacker reaches a caller webhook endpoint; signature bypass or caller omission or replay handling is demonstrated. | Unauthorized downstream action in the embedding application. | HMAC verification, constant-time comparison, timestamp window, canonical secret parsing. | Preserve verify-before-decode; distinguish SDK signature failures from caller-owned deduplication and authorization. | webhooks/webhook_helpers.go:1-143; webhooks/webhook_secret.go:1-34; README.md:430-520 |
| High | Workload or X.509 exchange leaks or mints a bearer through an unsafe destination, redirect, TLS policy, or refresh race. | Supported identity configuration plus a concrete bypass; not simply an operator-selected malicious signer or path. | Credential theft or unauthorized API access. | Fixed or validated destinations, redirect restrictions, TLS constraints, bounded issuer parsing, bearer history. | Review lifecycle, cancellation, rotation, replay, and concurrent refresh paths. | auth/subjecttokenprovider.go:1-113; auth/x509transport.go:1-120; option/x509workloadidentity.go:18-55 |
| Medium | Hostile JSON, SSE, error, cursor, or retry data causes disproportionate resource use, parser disagreement, or sensitive leakage. | Malicious compatible endpoint, compromised proxy, or attacker-influenced legitimate response reaches parser or diagnostic sink. | Process availability loss or limited disclosure. | SSE token limit, context cancellation, retry caps, typed decoding, safe debug logger. | Preserve large-payload compatibility; examine unbounded aggregation and explicit raw dump sinks separately. | packages/ssestream/ssestream.go:1-170; default_http_client.go:8-31; internal/apierror/apierror.go:1-56; option/middleware.go:1-114 |
| Medium | Caller passes end-user maps or model output into raw JSON or extra-field APIs or then dispatches a tool, renderer, or shell sink. | Embedding application exposes a trusted escape hatch or unsafe downstream sink. | Application-specific policy bypass or injection. | Typed APIs and explicit escape-hatch naming; SDK itself does not execute returned tools. | Require a demonstrated SDK-owned sink before reporting as an SDK flaw; otherwise document caller obligation. | README.md:163-219; option/requestoption.go:197-244; packages/param/param.go:25-55; packages/param/param.go:141-159 |
| Critical | Candidate PR source, artifact, report, or patch reaches a write token or release credential and alters official artifacts. | A concrete workflow path crosses from candidate-controlled execution or data into a protected job. | Supply-chain compromise of official SDK releases. | SHA-pinned actions, minimal permissions, protected release environment, trusted recomputation and publisher validation. | Trace event, artifact provenance, permissions, checkout, and shell evaluation end to end. | .github/workflows/codeql.yml:1-72; .github/workflows/castiron-custom-code-comment.yml:1-237; .github/workflows/create-releases.yml:1-55; .github/workflows/go-version-review.yml:200-411; scripts/castiron/CUSTOM_CODE.md:24-78 |
| Not a boundary by itself | A contributor edits a checked-in test, fixture, example, build script, or generator and candidate-safe CI deliberately executes that tracked file. | Contributor already has authority to propose that repository-code change. | No new capability solely from execution as repository code. | Code review and ordinary CI treat the candidate tree as the proposed program. | Reclassify only if lower-trust mutable data is parsed or evaluated, runtime data reaches a sensitive sink, or PR code reaches protected credentials; do not execute untrusted checkouts locally. | CONTRIBUTING.md:1-34; .github/workflows/ci.yml:1-183; .github/workflows/go-version-review.yml:200-411 |

## 4. Severity Calibration (Critical, High, Medium, Low)

- Critical: A remotely reachable or PR-controlled path compromises
  organization-admin credentials, protected release credentials, or official
  distributed artifacts without already possessing equivalent authority.
  Example: candidate-controlled data is evaluated in a release-credentialed
  publisher and changes published SDK content. Counterexample: editing a tracked
  test that a candidate-safe PR job intentionally runs is not a new privilege by
  itself; local review execution must still follow trusted-checkout rules.
- High: A realistic supported path bypasses webhook authenticity,
  credential-origin or provider isolation, or tenant or admin credential
  selection and yields substantial unauthorized action or usable credential
  theft. Counterexample: an operator deliberately configures an
  attacker-controlled base URL or malicious custom transport.
- Medium: Reachable hostile runtime data causes disproportionate process
  exhaustion, persistent leaks, or limited sensitive disclosure without broader
  compromise. Severity drops when exploitation requires unusual trusted
  configuration or caller misuse; it rises with default reachability and durable
  impact.
- Low: Nonsensitive metadata disclosure, bounded development-tooling failure, or
  a minor diagnostic issue with limited operational impact. Pure style,
  unsupported speculation, and behavior requiring already-equivalent
  repository-code authority are not security findings.

Rate each story by demonstrated reachability, attacker starting capability,
authority gained, blast radius, and effective controls. Missing evidence is an
investigation gap, not proof of safety; hypotheses are not findings until
validated.
