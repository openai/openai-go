module example.com/openai-go-consumer

go 1.25.0

require github.com/openai/openai-go/v3 v3.0.0

require (
	github.com/tidwall/gjson v1.19.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
)

replace github.com/openai/openai-go/v3 => ../../..
