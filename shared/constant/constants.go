// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package constant

import (
	"encoding/json"
)

type Constant[T any] interface {
	Default() T
}

// ValueOf gives the default value of a constant from its type. It's helpful when
// constructing constants as variants in a one-of. Note that empty structs are
// marshalled by default. Usage: constant.ValueOf[constant.Foo]()
func ValueOf[T Constant[T]]() T {
	var t T
	return t.Default()
}

type Assistant string               // Always "assistant"
type AssistantDeleted string        // Always "assistant.deleted"
type Auto string                    // Always "auto"
type Batch string                   // Always "batch"
type ChatCompletion string          // Always "chat.completion"
type ChatCompletionChunk string     // Always "chat.completion.chunk"
type ChatCompletionDeleted string   // Always "chat.completion.deleted"
type CodeInterpreter string         // Always "code_interpreter"
type Content string                 // Always "content"
type Default2024_08_21 string       // Always "default_2024_08_21"
type Developer string               // Always "developer"
type Embedding string               // Always "embedding"
type Error string                   // Always "error"
type File string                    // Always "file"
type FileCitation string            // Always "file_citation"
type FilePath string                // Always "file_path"
type FileSearch string              // Always "file_search"
type FineTuningJob string           // Always "fine_tuning.job"
type FineTuningJobCheckpoint string // Always "fine_tuning.job.checkpoint"
type FineTuningJobEvent string      // Always "fine_tuning.job.event"
type Function string                // Always "function"
type Image string                   // Always "image"
type ImageFile string               // Always "image_file"
type ImageURL string                // Always "image_url"
type InputAudio string              // Always "input_audio"
type JSONObject string              // Always "json_object"
type JSONSchema string              // Always "json_schema"
type LastActiveAt string            // Always "last_active_at"
type List string                    // Always "list"
type Logs string                    // Always "logs"
type MessageCreation string         // Always "message_creation"
type Model string                   // Always "model"
type Other string                   // Always "other"
type Refusal string                 // Always "refusal"
type Static string                  // Always "static"
type SubmitToolOutputs string       // Always "submit_tool_outputs"
type System string                  // Always "system"
type Text string                    // Always "text"
type TextCompletion string          // Always "text_completion"
type Thread string                  // Always "thread"
type ThreadCreated string           // Always "thread.created"
type ThreadDeleted string           // Always "thread.deleted"
type ThreadMessage string           // Always "thread.message"
type ThreadMessageCompleted string  // Always "thread.message.completed"
type ThreadMessageCreated string    // Always "thread.message.created"
type ThreadMessageDeleted string    // Always "thread.message.deleted"
type ThreadMessageDelta string      // Always "thread.message.delta"
type ThreadMessageInProgress string // Always "thread.message.in_progress"
type ThreadMessageIncomplete string // Always "thread.message.incomplete"
type ThreadRun string               // Always "thread.run"
type ThreadRunCancelled string      // Always "thread.run.cancelled"
type ThreadRunCancelling string     // Always "thread.run.cancelling"
type ThreadRunCompleted string      // Always "thread.run.completed"
type ThreadRunCreated string        // Always "thread.run.created"
type ThreadRunExpired string        // Always "thread.run.expired"
type ThreadRunFailed string         // Always "thread.run.failed"
type ThreadRunInProgress string     // Always "thread.run.in_progress"
type ThreadRunIncomplete string     // Always "thread.run.incomplete"
type ThreadRunQueued string         // Always "thread.run.queued"
type ThreadRunRequiresAction string // Always "thread.run.requires_action"
type ThreadRunStep string           // Always "thread.run.step"
type ThreadRunStepCancelled string  // Always "thread.run.step.cancelled"
type ThreadRunStepCompleted string  // Always "thread.run.step.completed"
type ThreadRunStepCreated string    // Always "thread.run.step.created"
type ThreadRunStepDelta string      // Always "thread.run.step.delta"
type ThreadRunStepExpired string    // Always "thread.run.step.expired"
type ThreadRunStepFailed string     // Always "thread.run.step.failed"
type ThreadRunStepInProgress string // Always "thread.run.step.in_progress"
type Tool string                    // Always "tool"
type ToolCalls string               // Always "tool_calls"
type Upload string                  // Always "upload"
type UploadPart string              // Always "upload.part"
type User string                    // Always "user"
type VectorStore string             // Always "vector_store"
type VectorStoreDeleted string      // Always "vector_store.deleted"
type VectorStoreFile string         // Always "vector_store.file"
type VectorStoreFileDeleted string  // Always "vector_store.file.deleted"
type VectorStoreFilesBatch string   // Always "vector_store.files_batch"
type Wandb string                   // Always "wandb"

func (c Assistant) Default() Assistant                         { return "assistant" }
func (c AssistantDeleted) Default() AssistantDeleted           { return "assistant.deleted" }
func (c Auto) Default() Auto                                   { return "auto" }
func (c Batch) Default() Batch                                 { return "batch" }
func (c ChatCompletion) Default() ChatCompletion               { return "chat.completion" }
func (c ChatCompletionChunk) Default() ChatCompletionChunk     { return "chat.completion.chunk" }
func (c ChatCompletionDeleted) Default() ChatCompletionDeleted { return "chat.completion.deleted" }
func (c CodeInterpreter) Default() CodeInterpreter             { return "code_interpreter" }
func (c Content) Default() Content                             { return "content" }
func (c Default2024_08_21) Default() Default2024_08_21         { return "default_2024_08_21" }
func (c Developer) Default() Developer                         { return "developer" }
func (c Embedding) Default() Embedding                         { return "embedding" }
func (c Error) Default() Error                                 { return "error" }
func (c File) Default() File                                   { return "file" }
func (c FileCitation) Default() FileCitation                   { return "file_citation" }
func (c FilePath) Default() FilePath                           { return "file_path" }
func (c FileSearch) Default() FileSearch                       { return "file_search" }
func (c FineTuningJob) Default() FineTuningJob                 { return "fine_tuning.job" }
func (c FineTuningJobCheckpoint) Default() FineTuningJobCheckpoint {
	return "fine_tuning.job.checkpoint"
}
func (c FineTuningJobEvent) Default() FineTuningJobEvent         { return "fine_tuning.job.event" }
func (c Function) Default() Function                             { return "function" }
func (c Image) Default() Image                                   { return "image" }
func (c ImageFile) Default() ImageFile                           { return "image_file" }
func (c ImageURL) Default() ImageURL                             { return "image_url" }
func (c InputAudio) Default() InputAudio                         { return "input_audio" }
func (c JSONObject) Default() JSONObject                         { return "json_object" }
func (c JSONSchema) Default() JSONSchema                         { return "json_schema" }
func (c LastActiveAt) Default() LastActiveAt                     { return "last_active_at" }
func (c List) Default() List                                     { return "list" }
func (c Logs) Default() Logs                                     { return "logs" }
func (c MessageCreation) Default() MessageCreation               { return "message_creation" }
func (c Model) Default() Model                                   { return "model" }
func (c Other) Default() Other                                   { return "other" }
func (c Refusal) Default() Refusal                               { return "refusal" }
func (c Static) Default() Static                                 { return "static" }
func (c SubmitToolOutputs) Default() SubmitToolOutputs           { return "submit_tool_outputs" }
func (c System) Default() System                                 { return "system" }
func (c Text) Default() Text                                     { return "text" }
func (c TextCompletion) Default() TextCompletion                 { return "text_completion" }
func (c Thread) Default() Thread                                 { return "thread" }
func (c ThreadCreated) Default() ThreadCreated                   { return "thread.created" }
func (c ThreadDeleted) Default() ThreadDeleted                   { return "thread.deleted" }
func (c ThreadMessage) Default() ThreadMessage                   { return "thread.message" }
func (c ThreadMessageCompleted) Default() ThreadMessageCompleted { return "thread.message.completed" }
func (c ThreadMessageCreated) Default() ThreadMessageCreated     { return "thread.message.created" }
func (c ThreadMessageDeleted) Default() ThreadMessageDeleted     { return "thread.message.deleted" }
func (c ThreadMessageDelta) Default() ThreadMessageDelta         { return "thread.message.delta" }
func (c ThreadMessageInProgress) Default() ThreadMessageInProgress {
	return "thread.message.in_progress"
}
func (c ThreadMessageIncomplete) Default() ThreadMessageIncomplete {
	return "thread.message.incomplete"
}
func (c ThreadRun) Default() ThreadRun                     { return "thread.run" }
func (c ThreadRunCancelled) Default() ThreadRunCancelled   { return "thread.run.cancelled" }
func (c ThreadRunCancelling) Default() ThreadRunCancelling { return "thread.run.cancelling" }
func (c ThreadRunCompleted) Default() ThreadRunCompleted   { return "thread.run.completed" }
func (c ThreadRunCreated) Default() ThreadRunCreated       { return "thread.run.created" }
func (c ThreadRunExpired) Default() ThreadRunExpired       { return "thread.run.expired" }
func (c ThreadRunFailed) Default() ThreadRunFailed         { return "thread.run.failed" }
func (c ThreadRunInProgress) Default() ThreadRunInProgress { return "thread.run.in_progress" }
func (c ThreadRunIncomplete) Default() ThreadRunIncomplete { return "thread.run.incomplete" }
func (c ThreadRunQueued) Default() ThreadRunQueued         { return "thread.run.queued" }
func (c ThreadRunRequiresAction) Default() ThreadRunRequiresAction {
	return "thread.run.requires_action"
}
func (c ThreadRunStep) Default() ThreadRunStep                   { return "thread.run.step" }
func (c ThreadRunStepCancelled) Default() ThreadRunStepCancelled { return "thread.run.step.cancelled" }
func (c ThreadRunStepCompleted) Default() ThreadRunStepCompleted { return "thread.run.step.completed" }
func (c ThreadRunStepCreated) Default() ThreadRunStepCreated     { return "thread.run.step.created" }
func (c ThreadRunStepDelta) Default() ThreadRunStepDelta         { return "thread.run.step.delta" }
func (c ThreadRunStepExpired) Default() ThreadRunStepExpired     { return "thread.run.step.expired" }
func (c ThreadRunStepFailed) Default() ThreadRunStepFailed       { return "thread.run.step.failed" }
func (c ThreadRunStepInProgress) Default() ThreadRunStepInProgress {
	return "thread.run.step.in_progress"
}
func (c Tool) Default() Tool                                     { return "tool" }
func (c ToolCalls) Default() ToolCalls                           { return "tool_calls" }
func (c Upload) Default() Upload                                 { return "upload" }
func (c UploadPart) Default() UploadPart                         { return "upload.part" }
func (c User) Default() User                                     { return "user" }
func (c VectorStore) Default() VectorStore                       { return "vector_store" }
func (c VectorStoreDeleted) Default() VectorStoreDeleted         { return "vector_store.deleted" }
func (c VectorStoreFile) Default() VectorStoreFile               { return "vector_store.file" }
func (c VectorStoreFileDeleted) Default() VectorStoreFileDeleted { return "vector_store.file.deleted" }
func (c VectorStoreFilesBatch) Default() VectorStoreFilesBatch   { return "vector_store.files_batch" }
func (c Wandb) Default() Wandb                                   { return "wandb" }

func (c Assistant) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c AssistantDeleted) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c Auto) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c Batch) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c ChatCompletion) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c ChatCompletionChunk) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c ChatCompletionDeleted) MarshalJSON() ([]byte, error)   { return marshalString(c) }
func (c CodeInterpreter) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c Content) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c Default2024_08_21) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c Developer) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c Embedding) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c Error) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c File) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c FileCitation) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c FilePath) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c FileSearch) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c FineTuningJob) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c FineTuningJobCheckpoint) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c FineTuningJobEvent) MarshalJSON() ([]byte, error)      { return marshalString(c) }
func (c Function) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c Image) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c ImageFile) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c ImageURL) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c InputAudio) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c JSONObject) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c JSONSchema) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c LastActiveAt) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c List) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c Logs) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c MessageCreation) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c Model) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c Other) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c Refusal) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c Static) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c SubmitToolOutputs) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c System) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c Text) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c TextCompletion) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c Thread) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c ThreadCreated) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c ThreadDeleted) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c ThreadMessage) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c ThreadMessageCompleted) MarshalJSON() ([]byte, error)  { return marshalString(c) }
func (c ThreadMessageCreated) MarshalJSON() ([]byte, error)    { return marshalString(c) }
func (c ThreadMessageDeleted) MarshalJSON() ([]byte, error)    { return marshalString(c) }
func (c ThreadMessageDelta) MarshalJSON() ([]byte, error)      { return marshalString(c) }
func (c ThreadMessageInProgress) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c ThreadMessageIncomplete) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c ThreadRun) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c ThreadRunCancelled) MarshalJSON() ([]byte, error)      { return marshalString(c) }
func (c ThreadRunCancelling) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c ThreadRunCompleted) MarshalJSON() ([]byte, error)      { return marshalString(c) }
func (c ThreadRunCreated) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c ThreadRunExpired) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c ThreadRunFailed) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c ThreadRunInProgress) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c ThreadRunIncomplete) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c ThreadRunQueued) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c ThreadRunRequiresAction) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c ThreadRunStep) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c ThreadRunStepCancelled) MarshalJSON() ([]byte, error)  { return marshalString(c) }
func (c ThreadRunStepCompleted) MarshalJSON() ([]byte, error)  { return marshalString(c) }
func (c ThreadRunStepCreated) MarshalJSON() ([]byte, error)    { return marshalString(c) }
func (c ThreadRunStepDelta) MarshalJSON() ([]byte, error)      { return marshalString(c) }
func (c ThreadRunStepExpired) MarshalJSON() ([]byte, error)    { return marshalString(c) }
func (c ThreadRunStepFailed) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c ThreadRunStepInProgress) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c Tool) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c ToolCalls) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c Upload) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c UploadPart) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c User) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c VectorStore) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c VectorStoreDeleted) MarshalJSON() ([]byte, error)      { return marshalString(c) }
func (c VectorStoreFile) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c VectorStoreFileDeleted) MarshalJSON() ([]byte, error)  { return marshalString(c) }
func (c VectorStoreFilesBatch) MarshalJSON() ([]byte, error)   { return marshalString(c) }
func (c Wandb) MarshalJSON() ([]byte, error)                   { return marshalString(c) }

type constant[T any] interface {
	Constant[T]
	*T
}

func marshalString[T ~string, PT constant[T]](v T) ([]byte, error) {
	var zero T
	if v == zero {
		v = PT(&v).Default()
	}
	return json.Marshal(string(v))
}
