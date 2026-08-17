# Go correctness and review heuristics

Use these checks where the diff makes them relevant. Report defects and
material maintainability regressions, not a transcript of every checklist item.

## Errors and control flow

- Check every returned error at the point where its owner can handle it. Watch
  for named error results overwritten by later iterations, ignored encoder and
  decoder failures, successful returns after failed checks, and branches that
  check one error but return a different nil error.
- Preserve the primary failure when cleanup also fails. On a successful
  multipart encoding, propagate `writer.Close()` failures; after an earlier
  encoding failure, close best-effort without replacing the original error.
- Prefer `errors.Is` and `errors.As` for wrappers, but do not apply them
  mechanically. `io.Reader` EOF handling and `http.ErrUseLastResponse` use
  documented exact-identity contracts in specific standard-library paths.
- Adding `%w` to an existing exported error is observable: it can change
  `errors.Is` and `errors.As` results for downstream callers. Treat changed
  wrapping as an API decision, not a cosmetic linter fix.
- Check empty inputs, nil interfaces containing typed nil pointers, nil maps,
  nil slices, pointer/value receivers, zero values, boundary indexes, integer
  conversions, time units, and overflow when they affect changed behavior.
- A slice created with nonzero length already contains elements; appending to
  it adds after those zero values. Distinguish length from capacity.

## Context, concurrency, and resources

- Propagate the caller's `context.Context` through request-owned HTTP work,
  foreground token acquisition, retry waits, and paging. An independently
  owned proactive refresh may intentionally use its own timeout-bounded
  context so one caller cannot cancel a shared credential-cache refresh.
- Separate the overall context deadline from per-attempt request timeouts;
  verify cancellation and timer cleanup on success, error, retry, and stream
  termination.
- Establish who closes every request body, response body, multipart writer,
  stream decoder, file, timer, channel, and spawned goroutine. Check
  unexpected-response and early-error branches, not only the success path.
- Preserve the request body's existing replayability and `GetBody` contract
  when adding retries, authentication refreshes, middleware reads, or SigV4
  signing. Do not manufacture replayability for a caller-supplied one-shot body.
- Review shared maps, slices, mutable client options, global decoder
  registration, token caches, refresh deduplication, and concurrent requests
  for data races and ownership leaks. Use `go test -race` for an affected
  package when the change warrants it.
- Check channel close ownership, blocked sends, loop-variable capture, lock
  ordering, callbacks under locks, duplicate refresh work, and goroutine exit
  paths when the changed code introduces concurrency.

## Go API compatibility

- Compare exported names, method signatures, pointer versus value receiver
  method sets, interfaces, generic constraints, aliases, embedding, struct
  fields, constants, enum values, and package import paths with the base.
- Preserve positional struct-literal compatibility and schema field order;
  do not reorder public generated or handwritten structs just to satisfy a
  layout suggestion.
- Distinguish compile-time API compatibility from wire compatibility. A
  change can compile while changing JSON keys, omission behavior, URL paths,
  retry decisions, errors, or provider authentication.
- For behavior changes, inspect an external consumer or real call site rather
  than relying only on tests inside the package.
- Do not introduce general-purpose interfaces, wrappers, reflection,
  configuration booleans, or type parameters without a demonstrated owner and
  use. Prefer the simplest existing abstraction that preserves the contract.

## Tests that prove behavior

- A regression test should fail against the pre-fix implementation and pass
  against the proposed implementation. When practical, confirm both states;
  otherwise trace why the new assertion distinguishes them.
- Assert the behavior that matters: first-error propagation, call counts,
  request method/path/body/headers, close ownership, retry attempts, context
  cancellation, selected credentials, signer scope, response presence, or
  returned error identity.
- Do not rely on randomized Go map iteration, scheduler timing, wall-clock
  races, external provider availability, inherited environment variables, or
  a preexisting local server for determinism.
- Use `httptest.Server`, real `net/http` requests, and existing SDK helpers
  when they faithfully exercise the transport. Use the smallest concrete
  protocol implementation only when a framework mock cannot satisfy the real
  boundary.
- Scope environment changes with `t.Setenv`, cleanup with `t.Cleanup`, and
  temporary files with `t.TempDir`. Avoid `t.Parallel()` where process-wide
  environment, global registries, ports, or singleton state would conflict.
- Check setup, fixture writes, JSON encoding, multipart finalization,
  `http.Response.Body.Close`, and cleanup errors instead of adding avoidable
  inline linter suppressions.
- Exercise success, failure, malformed input, empty values, retries,
  cancellation, provider-specific edge cases, and backward compatibility only
  when those branches are meaningfully affected.

## Practical checks

Select the smallest checks that answer the actual review question. Execute
head-controlled code or scripts only when the trust and isolation requirements
in `SKILL.md` permit it; otherwise use static review:

```sh
git diff --check
go test ./path/to/changed/package
go test -race ./path/to/concurrent/package
./scripts/lint
./scripts/check-go-mod
go -C examples test ./...
go -C internal/testdata/consumer test -mod=readonly ./...
```

The full root suite uses the Steady OpenAPI mock server. `./scripts/test`
starts the mock when necessary; setting `SKIP_MOCK_TESTS=true` can skip
mock-dependent tests and is not equivalent to full SDK verification.
