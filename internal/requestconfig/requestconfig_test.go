package requestconfig_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
)

type requestBodySizeTestCase struct {
	name          string
	size          int
	jsonPayload   []byte
	responseInput string
	requestBody   responses.ResponseNewParams
}

func makeJSONPayload(size int) []byte {
	const prefix = `{"data":"`
	const suffix = `"}`
	fillLen := size - len(prefix) - len(suffix)
	if fillLen < 0 {
		panic("payload size too small")
	}

	payload := make([]byte, size)
	copy(payload, prefix)
	for i := 0; i < fillLen; i++ {
		payload[len(prefix)+i] = 'a'
	}
	copy(payload[len(prefix)+fillLen:], suffix)
	return payload
}

func makeTextPayload(size int) string {
	return strings.Repeat("a", size)
}

func BenchmarkNewRequestConfig(b *testing.B) {
	requestBodySizeTestCases := []requestBodySizeTestCase{
		{name: "small_2KiB", size: 2 * 1024},
		{name: "medium_1MiB", size: 1024 * 1024},
		{name: "large_2MiB", size: 2 * 1024 * 1024},
	}

	for i := range requestBodySizeTestCases {
		input := makeTextPayload(requestBodySizeTestCases[i].size)
		requestBodySizeTestCases[i].jsonPayload = makeJSONPayload(requestBodySizeTestCases[i].size)
		requestBodySizeTestCases[i].responseInput = input
		requestBodySizeTestCases[i].requestBody = responses.ResponseNewParams{
			Model: shared.ResponsesModel("gpt-4o-mini"),
			Input: responses.ResponseNewParamsInputUnion{
				OfString: param.NewOpt(input),
			},
		}
	}

	withRequestBodyOptionCases := []struct {
		name string
		with bool
	}{
		{name: "without_with_request_body_option", with: false},
		{name: "with_with_request_body_option", with: true},
	}

	for _, requestBodySizeTestCase := range requestBodySizeTestCases {
		b.Run(requestBodySizeTestCase.name, func(b *testing.B) {
			for _, opt := range withRequestBodyOptionCases {
				b.Run(opt.name, func(b *testing.B) {
					body := requestBodySizeTestCase.requestBody
					var opts []requestconfig.RequestOption
					if opt.with {
						opts = []requestconfig.RequestOption{
							option.WithRequestBody("application/json", requestBodySizeTestCase.jsonPayload),
						}
					}

					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						_, err := requestconfig.NewRequestConfig(
							context.Background(),
							http.MethodPost,
							"https://example.com",
							body,
							nil,
							opts...,
						)
						if err != nil {
							b.Fatal(err)
						}
					}
				})
			}
		})
	}
}
