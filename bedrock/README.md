# Amazon Bedrock

Use the `bedrock` package to configure the normal OpenAI client for Amazon Bedrock's OpenAI-compatible API. The default Mantle endpoint and its existing authentication behavior remain unchanged.

| Endpoint | Default API root | SigV4 service |
| --- | --- | --- |
| `bedrock.EndpointMantle` (default) | `https://bedrock-mantle.<region>.api.aws/openai/v1` | `bedrock-mantle` |
| `bedrock.EndpointRuntime` | `https://bedrock-runtime.<region>.amazonaws.com/openai/v1` | `bedrock` |

Runtime hostnames use the correct suffix for the selected AWS partition. Canonical Runtime, FIPS, and dual-stack `BaseURL` overrides automatically select the endpoint family when `Endpoint` is omitted. Canonical AWS hosts must use HTTPS and match the configured region and endpoint. Custom or proxy hosts default to the Mantle signer; set `EndpointRuntime` explicitly when a custom host requires Runtime signing.

## Runtime Chat Completions

```go
import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/bedrock"
)

client, err := bedrock.NewClient(context.Background(), bedrock.Config{
	Endpoint:  bedrock.EndpointRuntime,
	AWSRegion: "us-east-1",
})
if err != nil {
	panic(err)
}

completion, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
	Model:    openai.ChatModel("us.openai.gpt-5.6-sol"),
	Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage("Say hello!")},
})
if err != nil {
	panic(err)
}
fmt.Println(completion.Choices[0].Message.Content)
```

The target US cross-Region inference profile IDs are `us.openai.gpt-5.6-sol`, `us.openai.gpt-5.6-terra`, and `us.openai.gpt-5.6-luna`. A global profile such as `global.openai.gpt-5.6-sol` requires account, role, and AWS Organizations policy access. The deployment requires the complete inference profile ID, not a bare model ID.

For the target US profiles, non-streaming Chat Completions at `/openai/v1/chat/completions` with the SigV4 `bedrock` signing service is the live-verified integration. The SDK also supports Responses and streaming, but support for a specific model, profile, authentication method, route, and stream is controlled by AWS and must be verified for that deployment.

AWS documentation describes both `/openai/v1` and `/v1` for different Bedrock OpenAI-compatible routes. If a deployment requires the alternative path, configure it explicitly:

```go
bedrock.Config{
	Endpoint:  bedrock.EndpointRuntime,
	AWSRegion: "us-east-1",
	BaseURL:   "https://bedrock-runtime.us-east-1.amazonaws.com/v1",
}
```

## Authentication

Credential precedence is unchanged:

1. Explicit `APIKey` or `BedrockTokenProvider`.
2. Explicit static AWS credentials, `AWSProfile`, or `AWSCredentialsProvider`.
3. `AWS_BEARER_TOKEN_BEDROCK`.
4. The standard AWS credential chain.

An expired `AWS_BEARER_TOKEN_BEDROCK` takes precedence over the implicit AWS credential chain. Unset it or configure an explicit AWS profile or credentials provider when SigV4 is required. Bearer providers are evaluated before each request attempt; AWS credentials and SigV4 signatures are refreshed as needed for retries. AWS determines whether a given deployment accepts bearer authentication.

The runnable Runtime Chat example defaults to SigV4 and clears any ambient Bedrock bearer token so it cannot shadow the default AWS credentials provider:

```sh
AWS_REGION=us-east-1 go -C examples run ./bedrock-runtime
AWS_REGION=us-east-1 BEDROCK_MODEL=us.openai.gpt-5.6-terra BEDROCK_STREAM=1 go -C examples run ./bedrock-runtime
AWS_REGION=us-east-1 BEDROCK_AUTH=bearer BEDROCK_MODEL=us.openai.gpt-5.6-luna go -C examples run ./bedrock-runtime
```

Run these commands from the repository root. Bearer mode requires `AWS_BEARER_TOKEN_BEDROCK`. Streaming is opt-in and requires support from the selected AWS deployment.

## Optional live verification

The Runtime live test never contacts AWS unless `BEDROCK_LIVE_TEST=1` is explicitly set. Choose one authentication mode from `bearer`, `environment-bearer`, `token-provider`, `default-chain`, `profile`, `static`, or `custom-provider`:

```sh
BEDROCK_LIVE_TEST=1 BEDROCK_LIVE_AUTH=default-chain AWS_REGION=us-east-1 go test -run TestRuntimeLive -v ./bedrock
```

By default, it sends non-streaming Chat requests for the three US inference profiles. `BEDROCK_MODEL` selects one profile; `BEDROCK_LIVE_MODELS` selects a comma-separated list. Set `BEDROCK_LIVE_RESPONSES=1` and `BEDROCK_LIVE_STREAM=1` only when that deployment supports the corresponding paths.
