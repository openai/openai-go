# Files for beta Agents

These local-path helpers expect application-controlled paths and stable source
directories. They are not a sandbox for arbitrary user-supplied paths or hostile
filesystem writers. File contents may still come from users.

Prepare selected local files before creating a hosted session:

```go
prepared, err := client.Beta.Agents.Environments.Files.Prepare(ctx,
    map[string]string{"/workspace/source.pdf": "source.pdf"})
if err != nil { return err }
params.Environment.OfParamOpenAIHosted.Files = prepared.Files
```

For a directory, use `PrepareDirectory(ctx, "./data", "/workspace/data",
[]string{"*.csv"})`. Patterns select relative paths once; files are not synchronized.
To stage a file into a running environment:

```go
staged, err := client.Beta.Agents.Environments.Files.Upload(ctx,
    environmentID, "source.pdf", "/workspace/source.pdf")
```

`prepared.Uploads` and `staged.Upload` identify the Files API objects you own.
Partial failures expose successful uploads through
`BetaAgentFilePreparationError.Prepared`. Delete uploads explicitly with the
ordinary Files API when no longer needed; preparation never deletes them.

Download an artifact from the completed result's exact turn and published path
into memory:

```go
var content bytes.Buffer
artifact, err := client.Beta.Agents.Sessions.Artifacts.ForResult(result).
    Download(ctx, "/workspace/outputs/report.md", &content)
```

Or write to an application-owned, safe destination path:

```go
destination, err := os.Create("downloaded-report.md")
if err != nil { return err }
defer destination.Close()
artifact, err := client.Beta.Agents.Sessions.Artifacts.ForResult(result).
    Download(ctx, "/workspace/outputs/report.md", destination)
```

Lookup follows all pages and reports missing or ambiguous matches. The caller
chooses the local destination; remote paths never become local file paths.
