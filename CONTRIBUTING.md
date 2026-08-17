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

You can also [install Go 1.25 or later manually](https://go.dev/doc/install).
CI tests every supported Go release line with `GOTOOLCHAIN=local`, so
contributors should not rely on automatic toolchain downloads to satisfy the
repository's minimum version.

## Modifying/Adding code

Most of the SDK is generated code. Modifications to code will be persisted between generations, but may
result in merge conflicts between manual patches and changes from the generator. The generator will never
modify the contents of the `lib/` and `examples/` directories.

## Security requirements

- Never commit API keys, access tokens, Azure or AWS credentials, signing
  keys, `.env` files, or customer data. Read credentials from environment
  variables such as `OPENAI_API_KEY`; use synthetic values and payloads in
  examples, tests, and fixtures. Prefer `httptest` or the existing local mock
  server over requests to live services.
- Always redact all credential-bearing request and response headers and other
  metadata, including `Authorization`, `Cookie`, `Set-Cookie`, `Api-Key`,
  `X-Api-Key`, and `X-Amz-Security-Token`; credentials in URLs or query
  parameters; and webhook secrets. Redact real or sensitive request or response
  bodies from default or uncontrolled logs, errors, fixtures, recordings, and CI
  output. Preserve safe synthetic payloads and intentional public API error or
  opt-in debug diagnostics; sanitize sensitive bodies before sharing them.
- Review direct and transitive dependency changes in `go.mod` and `go.sum`
  across the root, `examples`, `api_reference`, `internal/testdata/consumer`,
  and `tools` modules. Review new `replace` directives, dependency origins,
  installation scripts, code generators, and npm-based mock tooling before
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

Most tests require you to [set up a mock server](https://github.com/dgellow/steady) against the OpenAPI spec to run the tests.

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
