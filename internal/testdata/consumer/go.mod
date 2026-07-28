module example.com/openai-go-consumer

go 1.25.0

require github.com/openai/openai-go/v3 v3.0.0

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.22.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.12.0 // indirect
	github.com/tidwall/gjson v1.19.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/text v0.40.0 // indirect
)

replace github.com/openai/openai-go/v3 => ../../..
