# Development tools

This nested Go module pins repository-only tools without adding their
dependencies to the OpenAI SDK module graph.

`golangci-lint` and `govulncheck` are run from this module in CI. Dependabot
monitors the module weekly. To update them manually, run:

```sh
cd tools
go get -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go get -tool golang.org/x/vuln/cmd/govulncheck@latest
go mod tidy
```

Do not import packages from this module into the SDK.
