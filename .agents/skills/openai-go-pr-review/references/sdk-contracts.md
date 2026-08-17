# OpenAI Go SDK architecture and compatibility contracts

Follow only the sections touched by the proposed change. Confirm current
behavior against the exact PR-head implementation and repository documents.

## Ownership and generation

Most SDK models, resource methods, union accessors, generated tests, and some
shared scaffolding originate in Castiron. `lib/` and `examples/` are documented
as safe from generation, but the absence of a generated header does not prove
a file is handwritten: `packages/pagination/pagination.go` has previously
been identified as generated scaffolding.

For a proposed SDK-only fix, determine whether regeneration would overwrite
it or recreate a lint failure. Require an appropriate generator/template fix,
focused generator coverage, and regenerated output when generation owns the
behavior. Do not expand the public SDK surface merely to fix a quality warning.

## Request parameter encoding

Inspect `packages/param`, `internal/apijson`, `internal/apiquery`, and
`internal/apiform` when their boundaries are affected.

- Preserve four distinct states: omitted, explicit zero, explicit JSON null,
  and a supplied nonzero value.
- `param.Opt[T]`, `param.Null[T]()`, `param.NullStruct`, `param.NullMap`,
  `param.NullSlice`, `param.IsOmitted`, and `param.IsNull` deliberately have
  different representations and semantics.
- The non-nil empty map/slice returned by the null helpers is a sentinel. Do
  not mutate it or collapse it into an ordinary empty collection.
- Keep `omitzero`, required fields, embedded metadata, union variants,
  custom marshalers, explicit `param.Override`, `param.SetJSON`, and
  `SetExtraFields` compatible with existing wire output.
- A union accepts one present variant. Preserve clear errors for multiple
  variants and existing behavior for explicit null, overrides, and no variant.
- Extra-field keys must be escaped before use as `sjson` paths. Extra fields
  override same-named existing keys and can intentionally force omission; do
  not accept untrusted data into this escape hatch.
- Preserve query/form array and object encodings, map-key/value encoder
  errors, special characters, required multipart boundaries, file names,
  content types, and final writer-close errors.

## Response decoding and API compatibility

- `respjson.Field.Valid()` distinguishes valid values from missing, null, and
  invalid values. `Raw()` further distinguishes omitted (`""`), JSON null
  (`"null"`), and an invalid nonempty raw value.
- Preserve `JSON.ExtraFields`, `RawJSON`, unknown enum/union behavior,
  discriminators, structured error details, and generated accessor contracts.
- Generated union accessors may lack an error return. Do not change their
  exported signatures solely to surface an internal decoding error.
- URL path parameters must remain escaped exactly once. Review `%2F`, query
  characters, spaces, Unicode, Azure deployment IDs, `URL.Path`, and
  `URL.RawPath` together.

## Transport, options, and authentication

- Preserve client-level versus method-level option ordering, header set/add/
  delete behavior, API-key versus admin-key preference, custom HTTP clients,
  middleware ordering, query options, custom request bodies, and response
  capture ownership.
- Retry only replayable bodies. Preserve `GetBody`, server `x-should-retry`
  directives, retry-after units and dates, cancellation, per-attempt timeouts,
  retry limits, deterministic no-retry errors, and response-body cleanup.
- A middleware that consumes the request body must restore readability for
  downstream middleware while preserving the original `GetBody` and
  replayability contract. Do not synthesize retries for a one-shot body.
- Debug output must redact `Authorization`, `Api-Key`, `X-Api-Key`, AWS
  session tokens, cookies, and equivalent sensitive headers without mutating
  the headers of the live request or response.
- Workload identity must preserve request-owned caller contexts, cache
  expiry, timeout-bounded independently owned proactive refresh, concurrent
  refresh deduplication, provider errors, 401 refresh behavior, nonreplayable
  bodies, and close ownership.

## Server-sent events and streams

- Built-in SSE decoding ignores comments, keep-alives, retry-only blocks,
  and other blocks that contain no `data:` field.
- An explicit empty `data:` field is not the same as a block with no `data:`
  field; preserve the existing invalid-JSON behavior where applicable.
- Do not silently change custom registered decoders or synthesized custom
  events while fixing the built-in SSE parser.
- Preserve multiline data, event-type reset, `[DONE]` draining, structured
  stream errors, decoder/scanner errors, malformed JSON, typed union decoding,
  `Current`/`Err` semantics, scanner limits, and stream close ownership.

## Pagination

- Preserve the contracts of nonpaginated, cursor, conversation-cursor,
  `NextCursorPage`, and any other existing page shapes; inspect the concrete
  generated page type.
- Respect the difference between an omitted `has_more` field and an explicit
  `false`: presence is tracked through response JSON metadata.
- Derive the next cursor from the correct final item, response `last_id`, or
  nullable response `next` field. Stop cleanly for empty pages, absent
  cursors, exhausted data, errors, and caller cancellation.
- Preserve request option propagation, original query parameters, page
  response configuration, iterator order, index semantics, and `Err()`.
- Verify whether changed pagination scaffolding is generator-owned before
  proposing local-only cleanup or new lint enforcement.

## Azure

- Azure API-key authentication uses `Api-Key`, not OpenAI's bearer header.
  Explicit deletion of `Authorization` suppresses automatic injection of
  ambient `OPENAI_API_KEY` and `OPENAI_ADMIN_KEY`; preserve this protection.
- Preserve API-version query handling, custom scopes, token credential
  refresh, endpoint trailing slashes, deployment routing, escaped model/path
  values, multipart model extraction, and restored request bodies.
- Do not turn explicitly supported custom HTTP transports or proxies into an
  unqualified plaintext canonical-provider allowance.

## Amazon Bedrock

- `EndpointMantle` remains the default and signs for `bedrock-mantle`;
  `EndpointRuntime` signs for `bedrock`. Their hosts and configured paths are
  compatibility-sensitive.
- Preserve AWS partition suffixes, canonical/FIPS/dual-stack hostname
  inference, region matching, endpoint-family matching, canonical HTTPS,
  explicitly configured custom hosts, and custom endpoint paths.
- Explicit Bedrock authentication modes are mutually exclusive: reject an API
  key combined with a token provider, bearer credentials combined with AWS
  credentials, multiple explicit AWS modes, and provider credentials combined
  with `SkipAuth`. A valid explicit mode overrides ambient bearer credentials;
  otherwise `AWS_BEARER_TOKEN_BEDROCK` precedes the implicit AWS credential
  chain.
- Prevent ambient OpenAI settings from leaking into any Bedrock request. In
  Bedrock-authenticated modes, reject cross-origin requests before resolving
  credentials, preserve redirect restrictions, reject conflicting custom
  `Authorization` headers, sign the final request after middleware, and
  re-sign every retry without replaying an unreplayable body.
- Preserve the documented `Config.SkipAuth` gateway mode: explicit gateway
  API/admin credentials and custom HTTP doers are supported, Bedrock signing
  is disabled, and ambient OpenAI credentials must still remain suppressed.
- Refresh bearer tokens and SigV4 credentials as appropriate for the selected
  Bedrock-authenticated mode.
- Distinguish local transport simulations from live AWS verification. Never
  claim a model, profile, API path, streaming mode, or authentication mode is
  supported by an AWS deployment without actual deployment evidence.

## Webhooks

- Verify the exact raw request body before parsing it. Preserve required
  webhook ID, timestamp, and signature headers; accepted signature variants;
  HMAC-SHA256; constant-time comparison; and `whsec_` secret decoding.
- Reject stale and future timestamps outside the configured tolerance.
- Preserve client/environment webhook-secret configuration, malformed-header
  errors, and the distinction between signature-verification and
  payload-decoding errors. Check method-level secret options explicitly:
  `VerifySignatureWithToleranceAndTime` currently accepts those options but
  does not apply them, so do not assume method-level precedence already works.
- Never log webhook secrets, signatures, sensitive payload data, or provider
  credentials while adding diagnostics.
