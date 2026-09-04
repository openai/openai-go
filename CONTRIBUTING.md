## Setting up the environment

Before running any setup commands, review dependency origins, installation
tooling, and the [security requirements](#security-requirements). Bootstrap and
lint scripts can download Go, Homebrew, and direct or transitive dependencies.

To set up the repository, run:

```sh
$ ./scripts/bootstrap
$ ./scripts/lint
```

This will install all the required dependencies and build the SDK.

Contributors need [Go 1.25 or later](https://go.dev/doc/install) and
[Node.js 14 or later](https://nodejs.org/). Homebrew can install both from
the repository's `Brewfile`; they can also be installed manually.
CI tests every supported Go release line with `GOTOOLCHAIN=local`, so
contributors should not rely on automatic toolchain downloads to satisfy the
repository's minimum version.

## Modifying/Adding code

Most of the SDK is generated code. Modifications to code will be persisted between generations, but may
result in merge conflicts between manual patches and changes from the generator. The generator will never
modify the contents of the `lib/` and `examples/` directories.

## Custom-code budget

The custom-code budget counts additions plus deletions in the remaining patch
against verified generated output. `.castiron-ratchet.json` defines this repository's
ceiling. CI uses the checker and budget on main, not the PR's proposed versions.

Budget changes must be in a separate PR modifying **only `.castiron-ratchet.json`**.
Justify the current usage, proposed ceiling, and why fixing generation is not
appropriate in the PR description. Increases require a **human approving review**
and must merge before an SDK change relies on them. Agents may draft proposals,
but must not approve increases or bypass the gate. Keep default CODEOWNERS.
Lower the ceiling after cleanup while retaining headroom; decreases must still
fit the measured usage.

See [custom-code technical details](scripts/castiron/CUSTOM_CODE.md) for accounting,
local checks, trusted CI, and activation instructions.

## Security requirements

- Never commit API keys, access tokens, Azure or AWS credentials, signing
  keys, `.env` files, or customer data. Read credentials from environment
  variables such as `OPENAI_API_KEY`; use synthetic values and payloads in
  examples, tests, and fixtures. Prefer `httptest` or the existing local mock
  server over requests to live services.
- Redact all credential-bearing request and response headers and other
  metadata, including `Authorization`, `Cookie`, `Set-Cookie`, `Api-Key`,
  `X-Api-Key`, and `X-Amz-Security-Token`; credentials in URLs or query
  parameters; webhook secrets; and sensitive request or response bodies before
  logging, recording, forwarding, or disclosure. `Error.DumpRequest`,
  `Error.DumpResponse`, and `Error.Error` are intentionally raw, opt-in
  diagnostics and may include credentials, URLs, headers, or bodies; sanitize
  them before sharing or sending them to an untrusted sink. Preserve safe
  synthetic payloads and intentional public API error or debug diagnostics.
- Review direct and transitive dependency changes in `go.mod` and `go.sum`
  across the root, `examples`, `api_reference`, `internal/testdata/consumer`,
  and `tools` modules. Review new `replace` directives, dependency origins,
  installation scripts, code generators, and source-based mock tooling before
  running them; preserve Go checksum verification and run the existing
  vulnerability checks.
- Keep third-party GitHub Actions pinned to reviewed full commit SHAs, grant
  CI and publishing jobs only their required permissions, and isolate release
  credentials from untrusted pull requests or other untrusted code.
- Obtain SDK CODEOWNER review for changes involving authentication, webhook
  verification, base URLs, redirects, proxies, TLS, file paths or uploads, JSON
  or event-stream decoding, dependency installation, code generation, CI, or
  release tooling. Add focused public-entrypoint regression or security tests
  for executable behavior changes; use artifact-appropriate validation for
  docs, dependency, generated, CI, release, or policy-only changes.
- Report suspected vulnerabilities privately according to [SECURITY.md](SECURITY.md).
  Do not open public issues, pull requests, or discussions about them.

## Adding and running examples

All files in the `examples/` directory are not modified by the generator and can be freely edited or added to.

```go
# add an example to examples/<your-example>/main.go

package main

func main() {
  // ...
}
```

```sh
$ go run ./examples/<your-example>
```

## Using the repository from source

To use a local version of this library from source in another project, edit the `go.mod` with a replace
directive. This can be done through the CLI with the following:

```sh
$ go mod edit -replace github.com/openai/openai-go=/path/to/openai-go
```

## Running tests

The mock server uses [the OpenAI Steady fork](https://github.com/openai-oss-forks/steady).
`scripts/steady/manifest.json` is the single source of dependency pins: the
Steady Git commit and source digest, plus the Deno version and runtime checksums. `./scripts/steady/install` fetches that source, verifies the runtime,
and caches dependencies using the fork's frozen Deno lockfile. It requires
Git, Node.js, curl, unzip, and sha256sum or shasum. The installation supports
macOS and Linux on x64/ARM64, and Windows x64 through Git Bash.

`./scripts/run-steady` verifies the local source and runtime, then runs without
downloading dependencies. Pass a local OpenAPI specification path. To update
Steady, review the fork commit and run
`node scripts/steady/update.cjs <full-commit-sha>`. This updates the manifest
with the commit and its source digest; no launcher or test edits are needed.
Then run `./scripts/steady/install`. Review the release checksums when changing Deno.
Run `node scripts/steady/test.cjs` to check the
installation, integrity checks, and mock-server lifecycle.


Most tests require you to [set up a mock server](https://github.com/openai-oss-forks/steady) against the OpenAPI spec to run the tests.

`./scripts/bootstrap` and `./scripts/mock` install the pinned Steady source
and runtime with `./scripts/steady/install`.

```sh
$ ./scripts/mock
```

```sh
$ ./scripts/test
```

## Formatting

This library uses the standard `go fmt` command. Run it in each module that
contains Go source:

```sh
$ go fmt ./...
$ go -C examples fmt ./...
$ go -C internal/testdata/consumer fmt ./...
```

The [Go code quality policy](GO_CODE_QUALITY_POLICY.md) describes the staged
formatting and static-analysis rules for generated and handwritten code.
