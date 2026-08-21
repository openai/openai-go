package openai

import (
	"strings"

	"github.com/openai/openai-go/v3/packages/respjson"
)

func cloneAccumulatorString[T ~string](src T) T {
	return T(strings.Clone(string(src)))
}

func assignAccumulatorString[T ~string](dst *T, src T) {
	if *dst != src {
		*dst = cloneAccumulatorString(src)
	}
}

func detachChatCompletionLogprobs(src []ChatCompletionTokenLogprob) []ChatCompletionTokenLogprob {
	detached := make([]ChatCompletionTokenLogprob, len(src), cap(src))
	copy(detached, src)
	return detached
}

func cloneChatCompletionTokenLogprob(src ChatCompletionTokenLogprob) ChatCompletionTokenLogprob {
	dst := src
	dst.Token = strings.Clone(src.Token)
	dst.Bytes = cloneAccumulatorSlice(src.Bytes)
	dst.TopLogprobs = make([]ChatCompletionTokenLogprobTopLogprob, len(src.TopLogprobs))
	for i := range src.TopLogprobs {
		dst.TopLogprobs[i] = cloneChatCompletionTokenLogprobTopLogprob(src.TopLogprobs[i])
	}
	dst.JSON.Token = cloneAccumulatorField(src.JSON.Token)
	dst.JSON.Bytes = cloneAccumulatorField(src.JSON.Bytes)
	dst.JSON.Logprob = cloneAccumulatorField(src.JSON.Logprob)
	dst.JSON.TopLogprobs = cloneAccumulatorField(src.JSON.TopLogprobs)
	dst.JSON.ExtraFields = cloneAccumulatorFields(src.JSON.ExtraFields)
	dst.JSON.raw = strings.Clone(src.JSON.raw)
	return dst
}

func cloneChatCompletionTokenLogprobTopLogprob(src ChatCompletionTokenLogprobTopLogprob) ChatCompletionTokenLogprobTopLogprob {
	dst := src
	dst.Token = strings.Clone(src.Token)
	dst.Bytes = cloneAccumulatorSlice(src.Bytes)
	dst.JSON.Token = cloneAccumulatorField(src.JSON.Token)
	dst.JSON.Bytes = cloneAccumulatorField(src.JSON.Bytes)
	dst.JSON.Logprob = cloneAccumulatorField(src.JSON.Logprob)
	dst.JSON.ExtraFields = cloneAccumulatorFields(src.JSON.ExtraFields)
	dst.JSON.raw = strings.Clone(src.JSON.raw)
	return dst
}

func cloneAccumulatorSlice[T any](src []T) []T {
	if src == nil {
		return nil
	}
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

func cloneAccumulatorField(src respjson.Field) respjson.Field {
	raw := strings.Clone(src.Raw())
	if src.Valid() || raw == respjson.Null {
		return respjson.NewField(raw)
	}
	if raw != respjson.Omitted {
		return respjson.NewInvalidField(raw)
	}
	return respjson.Field{}
}

func cloneAccumulatorFields(src map[string]respjson.Field) map[string]respjson.Field {
	if src == nil {
		return nil
	}
	dst := make(map[string]respjson.Field, len(src))
	for name, field := range src {
		dst[strings.Clone(name)] = cloneAccumulatorField(field)
	}
	return dst
}
