package openai

import (
	"math/bits"
	"strings"
	"unsafe"

	"github.com/openai/openai-go/v3/packages/respjson"
)

type chatCompletionAccumulatorLogprobState struct {
	choices            []chatCompletionChoiceLogprobState
	activeChoices      []int
	bytes              int
	reconciliationWork int
}

type chatCompletionLogprobReconcilePlan struct {
	state        chatCompletionAccumulatorLogprobState
	detaches     []uint8
	measured     [maxChatCompletionAccumulatorStructuralSlots / 64]uint64
	contentBytes [maxChatCompletionAccumulatorStructuralSlots]int
	refusalBytes [maxChatCompletionAccumulatorStructuralSlots]int
}

const (
	detachContentLogprobs uint8 = 1 << iota
	detachRefusalLogprobs
)

type chatCompletionChoiceLogprobState struct {
	content chatCompletionLogprobSliceState
	refusal chatCompletionLogprobSliceState
}

type chatCompletionLogprobSliceState struct {
	// A uintptr fingerprints published storage without keeping cleared backing alive.
	data     uintptr
	length   int
	capacity int
	bytes    int
}

type chatCompletionLogprobProjection struct {
	contentCount int
	contentBytes int
	refusalCount int
	refusalBytes int
}

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
	for i := range src {
		detached[i] = cloneChatCompletionTokenLogprob(src[i])
	}
	return detached
}

func cloneChatCompletionTokenLogprob(src ChatCompletionTokenLogprob) ChatCompletionTokenLogprob {
	dst := src
	dst.Token = strings.Clone(src.Token)
	dst.Bytes = cloneAccumulatorSlice(src.Bytes)
	dst.TopLogprobs = cloneAccumulatorSlice(src.TopLogprobs)
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

func (acc *ChatCompletionAccumulator) planLogprobReconciliation(plan *chatCompletionLogprobReconcilePlan, includePopulatedChoices bool) (bool, bool) {
	state := &acc.logprobState
	indices := state.activeChoices
	if includePopulatedChoices {
		indices = acc.stringState.activeChoices
	}
	headerChanged := false
	for _, i := range indices {
		if i >= len(acc.Choices) {
			headerChanged = true
			continue
		}
		var current chatCompletionChoiceLogprobState
		if i < len(state.choices) {
			current = state.choices[i]
		}
		logprobs := &acc.Choices[i].Logprobs
		headerChanged = headerChanged || !current.content.matches(logprobs.Content) || !current.refusal.matches(logprobs.Refusal)
	}
	if !headerChanged {
		return false, true
	}

	plan.state = *state
	if headerChanged {
		plan.state.choices = make([]chatCompletionChoiceLogprobState, len(acc.Choices))
		copy(plan.state.choices, state.choices)
		plan.state.activeChoices = make([]int, 0, len(indices))
		plan.detaches = make([]uint8, len(acc.Choices))
	}

	plan.state.bytes = 0
	for _, i := range indices {
		if !addAccumulatorLogprobWork(&plan.state.reconciliationWork, 1) {
			return false, false
		}
		if i >= len(acc.Choices) {
			continue
		}
		logprobs := &acc.Choices[i].Logprobs
		var current chatCompletionChoiceLogprobState
		if i < len(state.choices) {
			current = state.choices[i]
		}
		content, contentDetach, ok := plan.state.projectSlice(current.content, logprobs.Content)
		if !ok {
			return false, false
		}
		refusal, refusalDetach, ok := plan.state.projectSlice(current.refusal, logprobs.Refusal)
		if !ok {
			return false, false
		}
		if !addAccumulatorLogprobBytes(&plan.state.bytes, content.bytes) ||
			!addAccumulatorLogprobBytes(&plan.state.bytes, refusal.bytes) {
			return false, false
		}
		plan.measured[i/64] |= uint64(1) << (i % 64)
		plan.contentBytes[i] = content.bytes
		plan.refusalBytes[i] = refusal.bytes
		if !headerChanged {
			continue
		}
		choiceState := &plan.state.choices[i]
		choiceState.content = content
		choiceState.refusal = refusal
		if content.data != 0 || refusal.data != 0 {
			plan.state.activeChoices = append(plan.state.activeChoices, i)
		}
		if contentDetach {
			plan.detaches[i] |= detachContentLogprobs
		}
		if refusalDetach {
			plan.detaches[i] |= detachRefusalLogprobs
		}
	}
	return true, true
}

func (state chatCompletionLogprobSliceState) matches(logprobs []ChatCompletionTokenLogprob) bool {
	header := chatCompletionLogprobHeader(logprobs)
	return state.data == header.data && state.length == header.length && state.capacity == header.capacity
}

func (state *chatCompletionAccumulatorLogprobState) projectSlice(current chatCompletionLogprobSliceState, logprobs []ChatCompletionTokenLogprob) (chatCompletionLogprobSliceState, bool, bool) {
	header := chatCompletionLogprobHeader(logprobs)
	detach := chatCompletionLogprobSliceRetainsBacking(current, header)

	bytes := 0
	retained := true
	if detach {
		retained = false
	}
	if !addAccumulatorLogprobStorage(&bytes, cap(logprobs), int(unsafe.Sizeof(ChatCompletionTokenLogprob{}))) ||
		!addMeasuredChatCompletionLogprobData(&bytes, logprobs, retained, &state.reconciliationWork) {
		return chatCompletionLogprobSliceState{}, false, false
	}
	if current.matches(logprobs) && bytes < current.bytes {
		// An unchanged outer slice may still retain nested backing that a public
		// reslice has hidden. Only an outer replacement or deep detachment proves
		// that storage has been released.
		bytes = current.bytes
	}
	header.bytes = bytes
	return header, detach, true
}

func chatCompletionChunkHasLogprobs(chunk *ChatCompletionChunk) bool {
	for i := range chunk.Choices {
		logprobs := &chunk.Choices[i].Logprobs
		if len(logprobs.Content) > 0 || len(logprobs.Refusal) > 0 {
			return true
		}
	}
	return false
}

func (acc *ChatCompletionAccumulator) applyLogprobReconciliation(plan *chatCompletionLogprobReconcilePlan) {
	if plan == nil {
		return
	}
	for i, detaches := range plan.detaches {
		logprobs := &acc.Choices[i].Logprobs
		choiceState := &plan.state.choices[i]
		if detaches&detachContentLogprobs != 0 {
			logprobs.Content = detachChatCompletionLogprobs(logprobs.Content)
			setChatCompletionLogprobHeader(&choiceState.content, logprobs.Content)
		}
		if detaches&detachRefusalLogprobs != 0 {
			logprobs.Refusal = detachChatCompletionLogprobs(logprobs.Refusal)
			setChatCompletionLogprobHeader(&choiceState.refusal, logprobs.Refusal)
		}
	}
	for word, measured := range plan.measured {
		for measured != 0 {
			bit := bits.TrailingZeros64(measured)
			i := word*64 + bit
			if i < len(plan.state.choices) {
				plan.state.choices[i].content.bytes = plan.contentBytes[i]
				plan.state.choices[i].refusal.bytes = plan.refusalBytes[i]
			}
			measured &^= uint64(1) << bit
		}
	}
	acc.logprobState = plan.state
}

func chatCompletionLogprobSliceRetainsBacking(current chatCompletionLogprobSliceState, public chatCompletionLogprobSliceState) bool {
	if current.data == 0 || current.capacity == 0 || public.data == 0 {
		return false
	}
	allocationEnd := current.data + uintptr(current.capacity)*unsafe.Sizeof(ChatCompletionTokenLogprob{})
	withinCurrentAllocation := public.data >= current.data && public.data < allocationEnd
	return withinCurrentAllocation && (public.data != current.data || public.capacity < current.capacity)
}

func (state *chatCompletionAccumulatorLogprobState) chunkWithinLimit(chunk *ChatCompletionChunk) bool {
	if chatCompletionChunkHasLogprobs(chunk) {
		return state.nonEmptyChunkWithinLimit(chunk)
	}
	return true
}

func (state *chatCompletionAccumulatorLogprobState) nonEmptyChunkWithinLimit(chunk *ChatCompletionChunk) bool {
	var projections [maxChatCompletionAccumulatorStructuralSlots]chatCompletionLogprobProjection
	var touched [maxChatCompletionAccumulatorStructuralSlots]int
	touchedCount := 0
	for i := range chunk.Choices {
		choice := &chunk.Choices[i]
		projection := &projections[choice.Index]
		wasEmpty := projection.contentCount == 0 && projection.refusalCount == 0
		if !addChatCompletionLogprobProjection(&projection.contentCount, &projection.contentBytes, choice.Logprobs.Content) ||
			!addChatCompletionLogprobProjection(&projection.refusalCount, &projection.refusalBytes, choice.Logprobs.Refusal) {
			return false
		}
		if wasEmpty && (projection.contentCount > 0 || projection.refusalCount > 0) {
			touched[touchedCount] = int(choice.Index)
			touchedCount++
		}
	}

	total := state.bytes
	logprobSize := int(unsafe.Sizeof(ChatCompletionTokenLogprob{}))
	for _, i := range touched[:touchedCount] {
		projection := &projections[i]
		var current chatCompletionChoiceLogprobState
		if i < len(state.choices) {
			current = state.choices[i]
		}
		if !addProjectedLogprobSliceBytes(&total, current.content, projection.contentCount, projection.contentBytes, logprobSize) ||
			!addProjectedLogprobSliceBytes(&total, current.refusal, projection.refusalCount, projection.refusalBytes, logprobSize) {
			return false
		}
	}
	return true
}

func addChatCompletionLogprobProjection(count *int, bytes *int, logprobs []ChatCompletionTokenLogprob) bool {
	maxLogprobs := maxChatCompletionAccumulatorLogprobBytes / int(unsafe.Sizeof(ChatCompletionTokenLogprob{}))
	if len(logprobs) > maxLogprobs-*count {
		return false
	}
	*count += len(logprobs)
	return addClonedChatCompletionLogprobData(bytes, logprobs)
}

func addProjectedLogprobSliceBytes(total *int, current chatCompletionLogprobSliceState, count int, dataBytes int, logprobSize int) bool {
	maxLogprobs := maxChatCompletionAccumulatorLogprobBytes / logprobSize
	if count > maxLogprobs-current.length {
		return false
	}
	capacity := projectedLogprobCapacity(current.capacity, current.length+count, maxLogprobs)
	if !addAccumulatorLogprobStorage(total, capacity-current.capacity, logprobSize) {
		return false
	}
	return addAccumulatorLogprobBytes(total, dataBytes)
}

func (state *chatCompletionAccumulatorLogprobState) acceptChunk(completion *ChatCompletion, chunk *ChatCompletionChunk) {
	if !chatCompletionChunkHasLogprobs(chunk) {
		return
	}
	if len(completion.Choices) > len(state.choices) {
		state.choices = expandToFit(state.choices, len(completion.Choices)-1)
	}
	logprobSize := int(unsafe.Sizeof(ChatCompletionTokenLogprob{}))
	maxLogprobs := maxChatCompletionAccumulatorLogprobBytes / logprobSize
	for i := range chunk.Choices {
		choice := &chunk.Choices[i]
		choiceState := &state.choices[choice.Index]
		wasInactive := choiceState.content.length == 0 && choiceState.content.capacity == 0 &&
			choiceState.refusal.length == 0 && choiceState.refusal.capacity == 0
		state.acceptSlice(&choiceState.content, choice.Logprobs.Content, logprobSize, maxLogprobs)
		state.acceptSlice(&choiceState.refusal, choice.Logprobs.Refusal, logprobSize, maxLogprobs)
		if wasInactive && (choiceState.content.length > 0 || choiceState.refusal.length > 0) {
			state.activeChoices = append(state.activeChoices, int(choice.Index))
		}
	}
	for i := range chunk.Choices {
		choiceIndex := int(chunk.Choices[i].Index)
		logprobs := &completion.Choices[choiceIndex].Logprobs
		setChatCompletionLogprobHeader(&state.choices[choiceIndex].content, logprobs.Content)
		setChatCompletionLogprobHeader(&state.choices[choiceIndex].refusal, logprobs.Refusal)
	}
}

func (state *chatCompletionAccumulatorLogprobState) acceptSlice(current *chatCompletionLogprobSliceState, logprobs []ChatCompletionTokenLogprob, logprobSize int, maxLogprobs int) {
	if len(logprobs) == 0 {
		return
	}
	dataBytes := 0
	// preflightChunk ran this same bounded calculation before accumulation.
	_ = addClonedChatCompletionLogprobData(&dataBytes, logprobs)
	capacity := projectedLogprobCapacity(current.capacity, current.length+len(logprobs), maxLogprobs)
	current.bytes += (capacity-current.capacity)*logprobSize + dataBytes
	state.bytes += (capacity-current.capacity)*logprobSize + dataBytes
	current.length += len(logprobs)
	current.capacity = capacity
}

func appendChatCompletionLogprobs(dst []ChatCompletionTokenLogprob, src []ChatCompletionTokenLogprob) []ChatCompletionTokenLogprob {
	if len(src) == 0 {
		return dst
	}
	required := len(dst) + len(src)
	if required > cap(dst) {
		maxLogprobs := maxChatCompletionAccumulatorLogprobBytes / int(unsafe.Sizeof(ChatCompletionTokenLogprob{}))
		capacity := projectedLogprobCapacity(cap(dst), required, maxLogprobs)
		grown := make([]ChatCompletionTokenLogprob, len(dst), capacity)
		copy(grown, dst)
		dst = grown
	}
	for i := range src {
		dst = append(dst, cloneChatCompletionTokenLogprob(src[i]))
	}
	return dst
}

func projectedLogprobCapacity(current int, required int, maximum int) int {
	if required <= current {
		return current
	}
	capacity := max(1, current)
	// Deterministic doubling keeps append amortized-linear while making the
	// retained logical capacity exactly predictable during preflight.
	for capacity < required {
		if capacity > maximum/2 {
			return required
		}
		capacity *= 2
	}
	return capacity
}

func chatCompletionLogprobHeader(logprobs []ChatCompletionTokenLogprob) chatCompletionLogprobSliceState {
	return chatCompletionLogprobSliceState{
		data:     uintptr(unsafe.Pointer(unsafe.SliceData(logprobs))),
		length:   len(logprobs),
		capacity: cap(logprobs),
	}
}

func setChatCompletionLogprobHeader(state *chatCompletionLogprobSliceState, logprobs []ChatCompletionTokenLogprob) {
	state.data = uintptr(unsafe.Pointer(unsafe.SliceData(logprobs)))
	state.length = len(logprobs)
	state.capacity = cap(logprobs)
}

func addClonedChatCompletionLogprobData(total *int, logprobs []ChatCompletionTokenLogprob) bool {
	return addMeasuredChatCompletionLogprobData(total, logprobs, false, nil)
}

func addMeasuredChatCompletionLogprobData(total *int, logprobs []ChatCompletionTokenLogprob, retainCapacity bool, work *int) bool {
	if retainCapacity {
		logprobs = logprobs[:cap(logprobs)]
	}
	for i := range logprobs {
		if work != nil && !addAccumulatorLogprobWork(work, 1) {
			return false
		}
		logprob := &logprobs[i]
		byteCount := len(logprob.Bytes)
		topLogprobCount := len(logprob.TopLogprobs)
		if retainCapacity {
			byteCount = cap(logprob.Bytes)
			topLogprobCount = cap(logprob.TopLogprobs)
		}
		if !addAccumulatorLogprobBytes(total, len(logprob.Token)) ||
			!addAccumulatorLogprobStorage(total, byteCount, int(unsafe.Sizeof(int64(0)))) ||
			!addAccumulatorLogprobStorage(total, topLogprobCount, int(unsafe.Sizeof(ChatCompletionTokenLogprobTopLogprob{}))) ||
			!addMeasuredChatCompletionLogprobMetadata(total, logprob.RawJSON(), logprob.JSON.ExtraFields, work,
				logprob.JSON.Token, logprob.JSON.Bytes, logprob.JSON.Logprob, logprob.JSON.TopLogprobs) {
			return false
		}
		topLogprobs := logprob.TopLogprobs
		if retainCapacity {
			topLogprobs = topLogprobs[:cap(topLogprobs)]
		}
		for j := range topLogprobs {
			if work != nil && !addAccumulatorLogprobWork(work, 1) {
				return false
			}
			topLogprob := &topLogprobs[j]
			topLogprobByteCount := len(topLogprob.Bytes)
			if retainCapacity {
				topLogprobByteCount = cap(topLogprob.Bytes)
			}
			if !addAccumulatorLogprobBytes(total, len(topLogprob.Token)) ||
				!addAccumulatorLogprobStorage(total, topLogprobByteCount, int(unsafe.Sizeof(int64(0)))) ||
				!addMeasuredChatCompletionLogprobMetadata(total, topLogprob.RawJSON(), topLogprob.JSON.ExtraFields, work,
					topLogprob.JSON.Token, topLogprob.JSON.Bytes, topLogprob.JSON.Logprob) {
				return false
			}
		}
	}
	return true
}

func addMeasuredChatCompletionLogprobMetadata(total *int, raw string, extraFields map[string]respjson.Field, work *int, fields ...respjson.Field) bool {
	if work != nil && !addAccumulatorLogprobWork(work, len(extraFields)) {
		return false
	}
	return addChatCompletionLogprobMetadata(total, raw, extraFields, fields...)
}

func addChatCompletionLogprobMetadata(total *int, raw string, extraFields map[string]respjson.Field, fields ...respjson.Field) bool {
	if len(extraFields) > 0 && !addAccumulatorLogprobBytes(total, chatCompletionAccumulatorMapOverheadBytes) {
		return false
	}
	entrySize := int(unsafe.Sizeof(string(""))) + int(unsafe.Sizeof(respjson.Field{}))
	if !addAccumulatorLogprobStorage(total, len(extraFields), entrySize) {
		return false
	}
	if !addAccumulatorLogprobStorage(total, len(extraFields), entrySize) {
		return false
	}
	if !addAccumulatorLogprobBytes(total, len(raw)) {
		return false
	}
	for i := range fields {
		if !addAccumulatorLogprobBytes(total, len(fields[i].Raw())) {
			return false
		}
	}
	for name, field := range extraFields {
		if !addAccumulatorLogprobBytes(total, len(name)) ||
			!addAccumulatorLogprobBytes(total, len(field.Raw())) {
			return false
		}
	}
	return true
}

func addAccumulatorLogprobStorage(total *int, count int, size int) bool {
	if count > (maxChatCompletionAccumulatorLogprobBytes-*total)/size {
		return false
	}
	*total += count * size
	return true
}

func addAccumulatorLogprobBytes(total *int, count int) bool {
	if count > maxChatCompletionAccumulatorLogprobBytes-*total {
		return false
	}
	*total += count
	return true
}

func addAccumulatorLogprobWork(total *int, count int) bool {
	if count > maxChatCompletionAccumulatorLogprobWork-*total {
		return false
	}
	*total += count
	return true
}

// Updates the internal response state and returns the previous state if
