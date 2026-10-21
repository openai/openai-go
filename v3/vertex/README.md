# Vertex AI

Use the Responses API through Google Vertex AI with the OpenAI Go client. You need a Google Cloud project, a location, a model available to your project, and Google Application Default Credentials.

Install the Vertex module:

```sh
go get github.com/openai/openai-go/v3/vertex
```

The Vertex module brings in Google authentication dependencies. Applications that only use `github.com/openai/openai-go/v3` do not need to install it.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/vertex"
)

func main() {
	ctx := context.Background()
	client, err := vertex.NewClient(ctx, vertex.Config{
		ProjectID: "your-project-id",
		Location:  "your-location",
	})
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Model: "your-model-id",
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("Hello")},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.ID)
}
```

`ProjectID` and `Location` default to `GOOGLE_CLOUD_PROJECT` and `GOOGLE_CLOUD_LOCATION`. To use credentials managed by your application, set `Credentials` or `TokenProvider` instead of using Application Default Credentials. The two fields are mutually exclusive.

For a custom endpoint, set `BaseURL` to the complete provider base before `/responses`. HTTP endpoints can be used for testing. The client sends credentials to the configured URL; use HTTPS for production.
