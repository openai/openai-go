package openai_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal/testutil"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

// Mock function to simulate weather data retrieval
func getWeather(_ string) string {
	// In a real implementation, this function would call a weather API
	return "Sunny, 25°C"
}

// Since the streamed response is hardcoded, we can hardcode the expected tool call
var expectedToolCall = openai.ChatCompletionMessageFunctionToolCallFunction{
	Arguments: `{"location":"Santorini, Greece"}`,
	Name:      "get_weather",
}

var expectedContents string = `Let's take a journey to the beautiful island of Santorini in Greece.

Santorini is a gem in the Aegean Sea, known for its stunning sunsets, white-washed buildings, and crystal-clear waters. The island's rugged landscape is the result of a volcanic eruption that took place around 3,600 years ago, which gave Santorini its distinct caldera. When you visit the island, you can't miss the picturesque town of Oia, perched high on the cliffs with its narrow streets, blue-domed churches, and panoramic views of the sea. Fira, the island's capital, is buzzing with life, offering a variety of shops, restaurants, and cafes. Don't forget to explore the archaeological site of Akrotiri, an ancient Minoan city preserved in volcanic ash, often referred to as the "Pompeii of the Aegean."

Now, let's check the weather in Santorini.`

func TestStreamingAccumulatorWithToolCalls(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}

	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			res, _ := next(req)
			res.Body = io.NopCloser(strings.NewReader(mockResponseBody))
			return res, nil
		}),
	)
	stream := client.Chat.Completions.NewStreaming(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			{OfSystem: &openai.ChatCompletionSystemMessageParam{
				Content: openai.ChatCompletionSystemMessageParamContentUnion{
					OfString: openai.String("Tell me a story about a place in Greece, then tell me the weather there."),
				},
				Name: openai.String("initialization"),
			}},
		},
		Model:    openai.ChatModelGPT4o,
		Logprobs: openai.Bool(true),
		StreamOptions: openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: openai.Bool(true),
		},
		Tools: []openai.ChatCompletionToolUnionParam{{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: shared.FunctionDefinitionParam{
					Name:        "get_weather",
					Description: openai.String("gets weather data"),
					Parameters: openai.FunctionParameters{
						"type": "object",
						"properties": map[string]interface{}{
							"location": map[string]string{
								"type": "string",
							},
						},
						"required":             []string{"location"},
						"additionalProperties": false,
					},
					Strict: openai.Bool(true),
				},
			},
		}},
		User: openai.String("user-1234"),
	})

	acc := openai.ChatCompletionAccumulator{}

	var err error

	anythingFinished := false

	for stream.Next() {
		chunk := stream.Current()
		if !acc.AddChunk(chunk) {
			err = errors.New("Chunk was not incorporated correctly")
			break
		}

		if _, ok := acc.JustFinishedContent(); ok {
			anythingFinished = true
		}
		if _, ok := acc.JustFinishedToolCall(); ok {
			anythingFinished = true
		}
		if _, ok := acc.JustFinishedRefusal(); ok {
			anythingFinished = true
		}
	}

	if err == nil {
		err = stream.Err()
	}

	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}

	if acc.Choices == nil || len(acc.Choices) == 0 {
		t.Fatal("No choices in accumulation")
	}

	if acc.Choices[0].Message.Content != expectedContents {
		t.Logf("%v", []byte(acc.Choices[0].Message.Content))
		t.Logf("%v", []byte(expectedContents))
		t.Fatalf("Found unexpected content")
	}

	if expectedToolCall.Arguments != acc.Choices[0].Message.ToolCalls[0].Function.Arguments || expectedToolCall.Name != acc.Choices[0].Message.ToolCalls[0].Function.Name {
		t.Fatalf("Found unexpected tool call %v %v", acc.Choices[0].Message.ToolCalls[0].Function.Arguments, acc.Choices[0].Message.ToolCalls[0].Function.Name)
	}

	if !anythingFinished {
		t.Fatalf("No finish events sent in accumulation")
	}
}

func TestAccumulateTokenDetails(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}

	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			res, err := next(req)
			if err != nil {
				return nil, err
			}
			res.Body = io.NopCloser(strings.NewReader(mockResponseBody))
			return res, nil
		}),
	)

	stream := client.Chat.Completions.NewStreaming(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			{OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("Tell me about Greece"),
				},
			}},
		},
		Model: openai.ChatModelGPT4o,
		StreamOptions: openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: openai.Bool(true),
		},
	})

	acc := openai.ChatCompletionAccumulator{}
	for stream.Next() {
		chunk := stream.Current()
		if !acc.AddChunk(chunk) {
			t.Fatal("Failed to accumulate chunk")
		}
	}

	if err := stream.Err(); err != nil {
		t.Fatalf("Stream error: %v", err)
	}

	// First chunk: prompt=10, completion=100, total=110
	// + Second chunk: prompt=0, completion=50, total=50
	// = Accumulated: prompt=10, completion=150, total=160
	if acc.Usage.PromptTokens != 10 {
		t.Errorf("PromptTokens: expected 10, got %d", acc.Usage.PromptTokens)
	}
	if acc.Usage.CompletionTokens != 150 {
		t.Errorf("CompletionTokens: expected 150, got %d", acc.Usage.CompletionTokens)
	}
	if acc.Usage.TotalTokens != 160 {
		t.Errorf("TotalTokens: expected 160, got %d", acc.Usage.TotalTokens)
	}

	// First chunk: reasoning=5, audio=2, accepted=8, rejected=1
	// + Second chunk: reasoning=3, audio=1, accepted=4, rejected=2
	// = Accumulated: reasoning=8, audio=3, accepted=12, rejected=3
	if acc.Usage.CompletionTokensDetails.ReasoningTokens != 8 {
		t.Errorf("CompletionTokensDetails.ReasoningTokens: expected 8, got %d", acc.Usage.CompletionTokensDetails.ReasoningTokens)
	}
	if acc.Usage.CompletionTokensDetails.AudioTokens != 3 {
		t.Errorf("CompletionTokensDetails.AudioTokens: expected 3, got %d", acc.Usage.CompletionTokensDetails.AudioTokens)
	}
	if acc.Usage.CompletionTokensDetails.AcceptedPredictionTokens != 12 {
		t.Errorf("CompletionTokensDetails.AcceptedPredictionTokens: expected 12, got %d", acc.Usage.CompletionTokensDetails.AcceptedPredictionTokens)
	}
	if acc.Usage.CompletionTokensDetails.RejectedPredictionTokens != 3 {
		t.Errorf("CompletionTokensDetails.RejectedPredictionTokens: expected 3, got %d", acc.Usage.CompletionTokensDetails.RejectedPredictionTokens)
	}

	// First chunk: audio=3, cached=20
	// + Second chunk: audio=0, cached=10
	// = Accumulated: audio=3, cached=30
	if acc.Usage.PromptTokensDetails.AudioTokens != 3 {
		t.Errorf("PromptTokensDetails.AudioTokens: expected 3, got %d", acc.Usage.PromptTokensDetails.AudioTokens)
	}
	if acc.Usage.PromptTokensDetails.CachedTokens != 30 {
		t.Errorf("PromptTokensDetails.CachedTokens: expected 30, got %d", acc.Usage.PromptTokensDetails.CachedTokens)
	}
}

// manually created on 11/3/2024
var mockResponseBody = `data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"role":"assistant","content":"","refusal":null},"logprobs":{"content":[],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"Let's"},"logprobs":{"content":[{"token":"Let's","logprob":-2.3433902,"bytes":[76,101,116,39,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" take"},"logprobs":{"content":[{"token":" take","logprob":-2.0225642,"bytes":[32,116,97,107,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" a"},"logprobs":{"content":[{"token":" a","logprob":-0.009433285,"bytes":[32,97],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" journey"},"logprobs":{"content":[{"token":" journey","logprob":-0.10697952,"bytes":[32,106,111,117,114,110,101,121],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" to"},"logprobs":{"content":[{"token":" to","logprob":-0.0023115498,"bytes":[32,116,111],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" the"},"logprobs":{"content":[{"token":" the","logprob":-0.5783461,"bytes":[32,116,104,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" beautiful"},"logprobs":{"content":[{"token":" beautiful","logprob":-1.2464501,"bytes":[32,98,101,97,117,116,105,102,117,108],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" island"},"logprobs":{"content":[{"token":" island","logprob":-0.11082827,"bytes":[32,105,115,108,97,110,100],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" of"},"logprobs":{"content":[{"token":" of","logprob":-5.3193703e-6,"bytes":[32,111,102],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" Sant"},"logprobs":{"content":[{"token":" Sant","logprob":-0.058394875,"bytes":[32,83,97,110,116],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"orini"},"logprobs":{"content":[{"token":"orini","logprob":0.0,"bytes":[111,114,105,110,105],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" in"},"logprobs":{"content":[{"token":" in","logprob":-0.76883507,"bytes":[32,105,110],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" Greece"},"logprobs":{"content":[{"token":" Greece","logprob":-0.0003801489,"bytes":[32,71,114,101,101,99,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":".\n\n"},"logprobs":{"content":[{"token":".\n\n","logprob":-0.17490852,"bytes":[46,10,10],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"Sant"},"logprobs":{"content":[{"token":"Sant","logprob":-1.6702999,"bytes":[83,97,110,116],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"orini"},"logprobs":{"content":[{"token":"orini","logprob":-5.5122365e-7,"bytes":[111,114,105,110,105],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" is"},"logprobs":{"content":[{"token":" is","logprob":-1.0545638,"bytes":[32,105,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" a"},"logprobs":{"content":[{"token":" a","logprob":-1.7354689,"bytes":[32,97],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" gem"},"logprobs":{"content":[{"token":" gem","logprob":-1.9755802,"bytes":[32,103,101,109],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" in"},"logprobs":{"content":[{"token":" in","logprob":-0.5036636,"bytes":[32,105,110],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" the"},"logprobs":{"content":[{"token":" the","logprob":-0.000013067608,"bytes":[32,116,104,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" A"},"logprobs":{"content":[{"token":" A","logprob":-0.0115876645,"bytes":[32,65],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"ege"},"logprobs":{"content":[{"token":"ege","logprob":0.0,"bytes":[101,103,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"an"},"logprobs":{"content":[{"token":"an","logprob":0.0,"bytes":[97,110],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" Sea"},"logprobs":{"content":[{"token":" Sea","logprob":-0.00016647171,"bytes":[32,83,101,97],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-0.06300592,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" known"},"logprobs":{"content":[{"token":" known","logprob":-0.5851157,"bytes":[32,107,110,111,119,110],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" for"},"logprobs":{"content":[{"token":" for","logprob":-0.00074875605,"bytes":[32,102,111,114],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" its"},"logprobs":{"content":[{"token":" its","logprob":-9.253091e-6,"bytes":[32,105,116,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" stunning"},"logprobs":{"content":[{"token":" stunning","logprob":-0.11636195,"bytes":[32,115,116,117,110,110,105,110,103],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" sunsets"},"logprobs":{"content":[{"token":" sunsets","logprob":-0.08226424,"bytes":[32,115,117,110,115,101,116,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-0.00254772,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" white"},"logprobs":{"content":[{"token":" white","logprob":-0.052607585,"bytes":[32,119,104,105,116,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"-w"},"logprobs":{"content":[{"token":"-w","logprob":-0.14288215,"bytes":[45,119],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"ashed"},"logprobs":{"content":[{"token":"ashed","logprob":-0.000030471343,"bytes":[97,115,104,101,100],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" buildings"},"logprobs":{"content":[{"token":" buildings","logprob":-0.11438937,"bytes":[32,98,117,105,108,100,105,110,103,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-1.1407294,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" and"},"logprobs":{"content":[{"token":" and","logprob":-0.051912103,"bytes":[32,97,110,100],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" crystal"},"logprobs":{"content":[{"token":" crystal","logprob":-0.9625359,"bytes":[32,99,114,121,115,116,97,108],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"-clear"},"logprobs":{"content":[{"token":"-clear","logprob":-0.05790539,"bytes":[45,99,108,101,97,114],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" waters"},"logprobs":{"content":[{"token":" waters","logprob":-0.20442869,"bytes":[32,119,97,116,101,114,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"."},"logprobs":{"content":[{"token":".","logprob":-0.00007755679,"bytes":[46],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" The"},"logprobs":{"content":[{"token":" The","logprob":-0.2418657,"bytes":[32,84,104,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" island"},"logprobs":{"content":[{"token":" island","logprob":-0.0027353284,"bytes":[32,105,115,108,97,110,100],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"'s"},"logprobs":{"content":[{"token":"'s","logprob":-1.6695559,"bytes":[39,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" rugged"},"logprobs":{"content":[{"token":" rugged","logprob":-5.071614,"bytes":[32,114,117,103,103,101,100],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" landscape"},"logprobs":{"content":[{"token":" landscape","logprob":-0.02823662,"bytes":[32,108,97,110,100,115,99,97,112,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" is"},"logprobs":{"content":[{"token":" is","logprob":-1.1156895,"bytes":[32,105,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" the"},"logprobs":{"content":[{"token":" the","logprob":-0.4554348,"bytes":[32,116,104,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" result"},"logprobs":{"content":[{"token":" result","logprob":-0.008874194,"bytes":[32,114,101,115,117,108,116],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" of"},"logprobs":{"content":[{"token":" of","logprob":-1.147242e-6,"bytes":[32,111,102],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" a"},"logprobs":{"content":[{"token":" a","logprob":-0.013579912,"bytes":[32,97],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" volcanic"},"logprobs":{"content":[{"token":" volcanic","logprob":-0.7000349,"bytes":[32,118,111,108,99,97,110,105,99],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" eruption"},"logprobs":{"content":[{"token":" eruption","logprob":-0.002232571,"bytes":[32,101,114,117,112,116,105,111,110],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" that"},"logprobs":{"content":[{"token":" that","logprob":-0.18755092,"bytes":[32,116,104,97,116],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" took"},"logprobs":{"content":[{"token":" took","logprob":-2.0938475,"bytes":[32,116,111,111,107],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" place"},"logprobs":{"content":[{"token":" place","logprob":-2.3392786e-6,"bytes":[32,112,108,97,99,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" around"},"logprobs":{"content":[{"token":" around","logprob":-3.3542585,"bytes":[32,97,114,111,117,110,100],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" "},"logprobs":{"content":[{"token":" ","logprob":-0.0046334034,"bytes":[32],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"3"},"logprobs":{"content":[{"token":"3","logprob":-0.19865946,"bytes":[51],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-6.704273e-7,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"600"},"logprobs":{"content":[{"token":"600","logprob":-0.006090777,"bytes":[54,48,48],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" years"},"logprobs":{"content":[{"token":" years","logprob":-2.220075e-6,"bytes":[32,121,101,97,114,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" ago"},"logprobs":{"content":[{"token":" ago","logprob":-6.704273e-7,"bytes":[32,97,103,111],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-0.35183582,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" which"},"logprobs":{"content":[{"token":" which","logprob":-1.2507571,"bytes":[32,119,104,105,99,104],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" gave"},"logprobs":{"content":[{"token":" gave","logprob":-3.6327841,"bytes":[32,103,97,118,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" Sant"},"logprobs":{"content":[{"token":" Sant","logprob":-0.89531523,"bytes":[32,83,97,110,116],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"orini"},"logprobs":{"content":[{"token":"orini","logprob":0.0,"bytes":[111,114,105,110,105],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" its"},"logprobs":{"content":[{"token":" its","logprob":-0.00020890454,"bytes":[32,105,116,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" distinct"},"logprobs":{"content":[{"token":" distinct","logprob":-3.0858865,"bytes":[32,100,105,115,116,105,110,99,116],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" cal"},"logprobs":{"content":[{"token":" cal","logprob":-1.9204621,"bytes":[32,99,97,108],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"dera"},"logprobs":{"content":[{"token":"dera","logprob":-0.00006742448,"bytes":[100,101,114,97],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"."},"logprobs":{"content":[{"token":".","logprob":-0.683659,"bytes":[46],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" When"},"logprobs":{"content":[{"token":" When","logprob":-4.303914,"bytes":[32,87,104,101,110],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" you"},"logprobs":{"content":[{"token":" you","logprob":-0.073804155,"bytes":[32,121,111,117],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" visit"},"logprobs":{"content":[{"token":" visit","logprob":-1.4989231,"bytes":[32,118,105,115,105,116],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" the"},"logprobs":{"content":[{"token":" the","logprob":-2.5110056,"bytes":[32,116,104,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" island"},"logprobs":{"content":[{"token":" island","logprob":-0.31087697,"bytes":[32,105,115,108,97,110,100],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-0.00087918504,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" you"},"logprobs":{"content":[{"token":" you","logprob":-0.43109176,"bytes":[32,121,111,117],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" can't"},"logprobs":{"content":[{"token":" can't","logprob":-2.211998,"bytes":[32,99,97,110,39,116],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" miss"},"logprobs":{"content":[{"token":" miss","logprob":-0.2015489,"bytes":[32,109,105,115,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" the"},"logprobs":{"content":[{"token":" the","logprob":-0.26818505,"bytes":[32,116,104,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" picturesque"},"logprobs":{"content":[{"token":" picturesque","logprob":-1.7780524,"bytes":[32,112,105,99,116,117,114,101,115,113,117,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" town"},"logprobs":{"content":[{"token":" town","logprob":-0.7454153,"bytes":[32,116,111,119,110],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" of"},"logprobs":{"content":[{"token":" of","logprob":-1.504853e-6,"bytes":[32,111,102],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" O"},"logprobs":{"content":[{"token":" O","logprob":-0.0036102824,"bytes":[32,79],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"ia"},"logprobs":{"content":[{"token":"ia","logprob":-0.00015860428,"bytes":[105,97],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-0.1405737,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" perched"},"logprobs":{"content":[{"token":" perched","logprob":-0.38644546,"bytes":[32,112,101,114,99,104,101,100],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" high"},"logprobs":{"content":[{"token":" high","logprob":-1.3419728,"bytes":[32,104,105,103,104],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" on"},"logprobs":{"content":[{"token":" on","logprob":-0.044032417,"bytes":[32,111,110],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" the"},"logprobs":{"content":[{"token":" the","logprob":-0.0126658585,"bytes":[32,116,104,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" cliffs"},"logprobs":{"content":[{"token":" cliffs","logprob":-0.08698677,"bytes":[32,99,108,105,102,102,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" with"},"logprobs":{"content":[{"token":" with","logprob":-1.8817015,"bytes":[32,119,105,116,104],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" its"},"logprobs":{"content":[{"token":" its","logprob":-0.4731294,"bytes":[32,105,116,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" narrow"},"logprobs":{"content":[{"token":" narrow","logprob":-0.97249043,"bytes":[32,110,97,114,114,111,119],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" streets"},"logprobs":{"content":[{"token":" streets","logprob":-0.46202877,"bytes":[32,115,116,114,101,101,116,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-0.44733757,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" blue"},"logprobs":{"content":[{"token":" blue","logprob":-0.5288072,"bytes":[32,98,108,117,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"-dom"},"logprobs":{"content":[{"token":"-dom","logprob":-0.018796053,"bytes":[45,100,111,109],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"ed"},"logprobs":{"content":[{"token":"ed","logprob":-3.650519e-6,"bytes":[101,100],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" churches"},"logprobs":{"content":[{"token":" churches","logprob":-0.0002193908,"bytes":[32,99,104,117,114,99,104,101,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-0.00009674858,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" and"},"logprobs":{"content":[{"token":" and","logprob":-0.00003166338,"bytes":[32,97,110,100],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" panoramic"},"logprobs":{"content":[{"token":" panoramic","logprob":-3.464097,"bytes":[32,112,97,110,111,114,97,109,105,99],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" views"},"logprobs":{"content":[{"token":" views","logprob":-0.05716486,"bytes":[32,118,105,101,119,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" of"},"logprobs":{"content":[{"token":" of","logprob":-0.19405605,"bytes":[32,111,102],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" the"},"logprobs":{"content":[{"token":" the","logprob":-0.000013306016,"bytes":[32,116,104,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" sea"},"logprobs":{"content":[{"token":" sea","logprob":-0.41413975,"bytes":[32,115,101,97],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"."},"logprobs":{"content":[{"token":".","logprob":-0.5598801,"bytes":[46],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" F"},"logprobs":{"content":[{"token":" F","logprob":-2.284997,"bytes":[32,70],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"ira"},"logprobs":{"content":[{"token":"ira","logprob":-0.0026779182,"bytes":[105,114,97],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-0.0013514261,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" the"},"logprobs":{"content":[{"token":" the","logprob":-0.005124177,"bytes":[32,116,104,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" island"},"logprobs":{"content":[{"token":" island","logprob":-0.27476153,"bytes":[32,105,115,108,97,110,100],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"'s"},"logprobs":{"content":[{"token":"'s","logprob":-0.11290983,"bytes":[39,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" capital"},"logprobs":{"content":[{"token":" capital","logprob":-0.3653767,"bytes":[32,99,97,112,105,116,97,108],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-8.418666e-6,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" is"},"logprobs":{"content":[{"token":" is","logprob":-0.6775553,"bytes":[32,105,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" buzzing"},"logprobs":{"content":[{"token":" buzzing","logprob":-5.6459823,"bytes":[32,98,117,122,122,105,110,103],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" with"},"logprobs":{"content":[{"token":" with","logprob":-3.2929079e-6,"bytes":[32,119,105,116,104],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" life"},"logprobs":{"content":[{"token":" life","logprob":-0.8478267,"bytes":[32,108,105,102,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-0.19484033,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" offering"},"logprobs":{"content":[{"token":" offering","logprob":-0.24219286,"bytes":[32,111,102,102,101,114,105,110,103],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" a"},"logprobs":{"content":[{"token":" a","logprob":-0.34762233,"bytes":[32,97],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" variety"},"logprobs":{"content":[{"token":" variety","logprob":-3.6589396,"bytes":[32,118,97,114,105,101,116,121],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" of"},"logprobs":{"content":[{"token":" of","logprob":-2.3392786e-6,"bytes":[32,111,102],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" shops"},"logprobs":{"content":[{"token":" shops","logprob":-0.2816828,"bytes":[32,115,104,111,112,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-0.0000118755715,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" restaurants"},"logprobs":{"content":[{"token":" restaurants","logprob":-0.35918993,"bytes":[32,114,101,115,116,97,117,114,97,110,116,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-0.0000258224,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" and"},"logprobs":{"content":[{"token":" and","logprob":-0.00096065435,"bytes":[32,97,110,100],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" cafes"},"logprobs":{"content":[{"token":" cafes","logprob":-2.4229627,"bytes":[32,99,97,102,101,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"."},"logprobs":{"content":[{"token":".","logprob":-0.47348917,"bytes":[46],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" Don't"},"logprobs":{"content":[{"token":" Don't","logprob":-2.2735655,"bytes":[32,68,111,110,39,116],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" forget"},"logprobs":{"content":[{"token":" forget","logprob":-0.010000905,"bytes":[32,102,111,114,103,101,116],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" to"},"logprobs":{"content":[{"token":" to","logprob":-0.0008115323,"bytes":[32,116,111],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" explore"},"logprobs":{"content":[{"token":" explore","logprob":-0.4913598,"bytes":[32,101,120,112,108,111,114,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" the"},"logprobs":{"content":[{"token":" the","logprob":-0.008387392,"bytes":[32,116,104,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" archaeological"},"logprobs":{"content":[{"token":" archaeological","logprob":-1.1182892,"bytes":[32,97,114,99,104,97,101,111,108,111,103,105,99,97,108],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" site"},"logprobs":{"content":[{"token":" site","logprob":-0.16406576,"bytes":[32,115,105,116,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" of"},"logprobs":{"content":[{"token":" of","logprob":-0.00072994747,"bytes":[32,111,102],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" Ak"},"logprobs":{"content":[{"token":" Ak","logprob":-0.00030984072,"bytes":[32,65,107],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"rot"},"logprobs":{"content":[{"token":"rot","logprob":0.0,"bytes":[114,111,116],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"iri"},"logprobs":{"content":[{"token":"iri","logprob":-2.4584822e-6,"bytes":[105,114,105],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-0.014347774,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" an"},"logprobs":{"content":[{"token":" an","logprob":-0.38988778,"bytes":[32,97,110],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" ancient"},"logprobs":{"content":[{"token":" ancient","logprob":-0.0001006823,"bytes":[32,97,110,99,105,101,110,116],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" M"},"logprobs":{"content":[{"token":" M","logprob":-0.08886115,"bytes":[32,77],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"ino"},"logprobs":{"content":[{"token":"ino","logprob":-4.3202e-7,"bytes":[105,110,111],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"an"},"logprobs":{"content":[{"token":"an","logprob":-5.5122365e-7,"bytes":[97,110],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" city"},"logprobs":{"content":[{"token":" city","logprob":-0.040044695,"bytes":[32,99,105,116,121],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" preserved"},"logprobs":{"content":[{"token":" preserved","logprob":-0.0767503,"bytes":[32,112,114,101,115,101,114,118,101,100],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" in"},"logprobs":{"content":[{"token":" in","logprob":-0.20755702,"bytes":[32,105,110],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" volcanic"},"logprobs":{"content":[{"token":" volcanic","logprob":-0.007286554,"bytes":[32,118,111,108,99,97,110,105,99],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" ash"},"logprobs":{"content":[{"token":" ash","logprob":-0.00014180024,"bytes":[32,97,115,104],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-0.2912434,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" often"},"logprobs":{"content":[{"token":" often","logprob":-0.580934,"bytes":[32,111,102,116,101,110],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" referred"},"logprobs":{"content":[{"token":" referred","logprob":-0.08569331,"bytes":[32,114,101,102,101,114,114,101,100],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" to"},"logprobs":{"content":[{"token":" to","logprob":-9.729906e-6,"bytes":[32,116,111],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" as"},"logprobs":{"content":[{"token":" as","logprob":-5.5122365e-7,"bytes":[32,97,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" the"},"logprobs":{"content":[{"token":" the","logprob":-0.004702584,"bytes":[32,116,104,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" \""},"logprobs":{"content":[{"token":" \"","logprob":-0.025759168,"bytes":[32,34],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"Pom"},"logprobs":{"content":[{"token":"Pom","logprob":-0.46578047,"bytes":[80,111,109],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"pe"},"logprobs":{"content":[{"token":"pe","logprob":-0.000203898,"bytes":[112,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"ii"},"logprobs":{"content":[{"token":"ii","logprob":-1.9361265e-7,"bytes":[105,105],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" of"},"logprobs":{"content":[{"token":" of","logprob":-4.8425554e-6,"bytes":[32,111,102],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" the"},"logprobs":{"content":[{"token":" the","logprob":-0.0011754631,"bytes":[32,116,104,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" A"},"logprobs":{"content":[{"token":" A","logprob":-0.000026895234,"bytes":[32,65],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"ege"},"logprobs":{"content":[{"token":"ege","logprob":-4.3202e-7,"bytes":[101,103,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"an"},"logprobs":{"content":[{"token":"an","logprob":-7.89631e-7,"bytes":[97,110],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":".\"\n\n"},"logprobs":{"content":[{"token":".\"\n\n","logprob":-0.13864219,"bytes":[46,34,10,10],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"Now"},"logprobs":{"content":[{"token":"Now","logprob":-0.07598482,"bytes":[78,111,119],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":","},"logprobs":{"content":[{"token":",","logprob":-0.023346568,"bytes":[44],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" let's"},"logprobs":{"content":[{"token":" let's","logprob":-0.10193493,"bytes":[32,108,101,116,39,115],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" check"},"logprobs":{"content":[{"token":" check","logprob":-0.6507268,"bytes":[32,99,104,101,99,107],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" the"},"logprobs":{"content":[{"token":" the","logprob":-0.03336788,"bytes":[32,116,104,101],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" weather"},"logprobs":{"content":[{"token":" weather","logprob":-1.1369112,"bytes":[32,119,101,97,116,104,101,114],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" in"},"logprobs":{"content":[{"token":" in","logprob":-0.03188107,"bytes":[32,105,110],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":" Sant"},"logprobs":{"content":[{"token":" Sant","logprob":-0.00010473523,"bytes":[32,83,97,110,116],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"orini"},"logprobs":{"content":[{"token":"orini","logprob":-1.2664457e-6,"bytes":[111,114,105,110,105],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"content":"."},"logprobs":{"content":[{"token":".","logprob":-0.20214307,"bytes":[46],"top_logprobs":[]}],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_FXoAjBUMcVv1k40fficJ9cSs","type":"function","function":{"name":"get_weather","arguments":""}}]},"logprobs":{"content":[],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\""}}]},"logprobs":{"content":[],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"location"}}]},"logprobs":{"content":[],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"\":\""}}]},"logprobs":{"content":[],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"Sant"}}]},"logprobs":{"content":[],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"orini"}}]},"logprobs":{"content":[],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":","}}]},"logprobs":{"content":[],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":" Greece"}}]},"logprobs":{"content":[],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"\"}"}}]},"logprobs":{"content":[],"refusal":null},"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[{"index":0,"delta":{},"logprobs":null,"finish_reason":"tool_calls"}],"usage":null}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[],"usage":{"prompt_tokens":10,"completion_tokens":100,"total_tokens":110,"completion_tokens_details":{"reasoning_tokens":5,"audio_tokens":2,"accepted_prediction_tokens":8,"rejected_prediction_tokens":1},"prompt_tokens_details":{"audio_tokens":3,"cached_tokens":20}}}

data: {"id":"chatcmpl-A3Tguz3LSXTHBTY2NAPBCSyfBltxF","object":"chat.completion.chunk","created":1725392480,"model":"gpt-4o-2024-05-13","system_fingerprint":"fp_157b3831f5","choices":[],"usage":{"prompt_tokens":0,"completion_tokens":50,"total_tokens":50,"completion_tokens_details":{"reasoning_tokens":3,"audio_tokens":1,"accepted_prediction_tokens":4,"rejected_prediction_tokens":2},"prompt_tokens_details":{"audio_tokens":0,"cached_tokens":10}}}

data: [DONE]

`
