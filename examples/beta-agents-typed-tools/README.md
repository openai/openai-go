# Typed application tools (beta)

Keep your application's argument structs and service methods. Supply a JSON schema
and adapt the handler to `AgentToolHandler`; the SDK's existing stream dispatches it.
This example uses the schema library already used by the structured-output example.
You can supply the schema directly or use another library.

```go
action := lookupAction{catalog: demoCatalog{id: "demo-catalog"}}
tool, err := action.tool()
if err != nil {
    return err
}
// Supply the hosted definition in agent.tools when creating a session:
agent := openai.BetaAgentSessionNewParamsAgent{
    Model: openai.String("YOUR_MODEL"),
    Tools: []openai.AgentToolParamUnion{tool},
}
// Reuse the local binding for each follow-up on that idle session:
handlers := map[string]openai.AgentToolHandler{
    tool.OfParamFunction.Name: action.handle,
}
```

The bound catalog is a local, read-only fixture. Replace it with your application's
service.
The handler decodes and validates arguments before calling the service and returns
the typed record as JSON text. Catalog credentials remain local.

From `examples/`, print the tool definition:

```sh
go run ./beta-agents-typed-tools
```

Configure a session with that definition, then run against it once it is idle:

```sh
export OPENAI_API_KEY=...
export AGENT_SESSION_ID=...
go run ./beta-agents-typed-tools
```

This recipe uses the existing-session stream helper. Creating or following an
active session with a raw event stream still requires application-managed dispatch.
