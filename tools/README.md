# Development tools

This nested Go module pins repository-only tools without adding their
dependencies to the OpenAI SDK module graph.

`govulncheck` is installed from this module in CI. Dependabot monitors the
module weekly. To update it manually, run:

```sh
cd tools
go get -tool golang.org/x/vuln/cmd/govulncheck@latest
go mod tidy
```

Do not import packages from this module into the SDK.
