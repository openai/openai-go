# Shared Params Types

- <a href="https://pkg.go.dev/github.com/openai/openai-go/shared">shared</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/shared#ChatModel">ChatModel</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/shared">shared</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/shared#FunctionDefinitionParam">FunctionDefinitionParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/shared">shared</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/shared#FunctionParameters">FunctionParameters</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/shared">shared</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/shared#MetadataParam">MetadataParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/shared">shared</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/shared#ResponseFormatJSONObjectParam">ResponseFormatJSONObjectParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/shared">shared</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/shared#ResponseFormatJSONSchemaParam">ResponseFormatJSONSchemaParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/shared">shared</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/shared#ResponseFormatTextParam">ResponseFormatTextParam</a>

# Shared Response Types

- <a href="https://pkg.go.dev/github.com/openai/openai-go/shared">shared</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/shared#ErrorObject">ErrorObject</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/shared">shared</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/shared#FunctionDefinition">FunctionDefinition</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/shared">shared</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/shared#FunctionParameters">FunctionParameters</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/shared">shared</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/shared#Metadata">Metadata</a>

# Completions

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Completion">Completion</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#CompletionChoice">CompletionChoice</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#CompletionUsage">CompletionUsage</a>

Methods:

- <code title="post /completions">client.Completions.<a href="https://pkg.go.dev/github.com/openai/openai-go#CompletionService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#CompletionNewParams">CompletionNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Completion">Completion</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Chat

## Completions

Params Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionAssistantMessageParam">ChatCompletionAssistantMessageParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionAudioParam">ChatCompletionAudioParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionContentPartUnionParam">ChatCompletionContentPartUnionParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionContentPartImageParam">ChatCompletionContentPartImageParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionContentPartInputAudioParam">ChatCompletionContentPartInputAudioParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionContentPartRefusalParam">ChatCompletionContentPartRefusalParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionContentPartTextParam">ChatCompletionContentPartTextParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionDeveloperMessageParam">ChatCompletionDeveloperMessageParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionFunctionCallOptionParam">ChatCompletionFunctionCallOptionParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionFunctionMessageParam">ChatCompletionFunctionMessageParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionMessageParamUnion">ChatCompletionMessageParamUnion</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionMessageToolCallParam">ChatCompletionMessageToolCallParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionModality">ChatCompletionModality</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionNamedToolChoiceParam">ChatCompletionNamedToolChoiceParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionPredictionContentParam">ChatCompletionPredictionContentParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionReasoningEffort">ChatCompletionReasoningEffort</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionStreamOptionsParam">ChatCompletionStreamOptionsParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionSystemMessageParam">ChatCompletionSystemMessageParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionToolParam">ChatCompletionToolParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionToolChoiceOptionUnionParam">ChatCompletionToolChoiceOptionUnionParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionToolMessageParam">ChatCompletionToolMessageParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionUserMessageParam">ChatCompletionUserMessageParam</a>

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletion">ChatCompletion</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionAudio">ChatCompletionAudio</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionChunk">ChatCompletionChunk</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionDeleted">ChatCompletionDeleted</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionMessage">ChatCompletionMessage</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionMessageToolCall">ChatCompletionMessageToolCall</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionStoreMessage">ChatCompletionStoreMessage</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionTokenLogprob">ChatCompletionTokenLogprob</a>

Methods:

- <code title="post /chat/completions">client.Chat.Completions.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionNewParams">ChatCompletionNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletion">ChatCompletion</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /chat/completions/{completion_id}">client.Chat.Completions.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, completionID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletion">ChatCompletion</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /chat/completions/{completion_id}">client.Chat.Completions.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionService.Update">Update</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, completionID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionUpdateParams">ChatCompletionUpdateParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletion">ChatCompletion</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /chat/completions">client.Chat.Completions.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionListParams">ChatCompletionListParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#CursorPage">CursorPage</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletion">ChatCompletion</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="delete /chat/completions/{completion_id}">client.Chat.Completions.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionService.Delete">Delete</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, completionID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionDeleted">ChatCompletionDeleted</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

### Messages

Methods:

- <code title="get /chat/completions/{completion_id}/messages">client.Chat.Completions.Messages.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionMessageService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, completionID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionMessageListParams">ChatCompletionMessageListParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#CursorPage">CursorPage</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ChatCompletionStoreMessage">ChatCompletionStoreMessage</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Embeddings

Params Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#EmbeddingModel">EmbeddingModel</a>

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#CreateEmbeddingResponse">CreateEmbeddingResponse</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Embedding">Embedding</a>

Methods:

- <code title="post /embeddings">client.Embeddings.<a href="https://pkg.go.dev/github.com/openai/openai-go#EmbeddingService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#EmbeddingNewParams">EmbeddingNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#CreateEmbeddingResponse">CreateEmbeddingResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Files

Params Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FilePurpose">FilePurpose</a>

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileDeleted">FileDeleted</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileObject">FileObject</a>

Methods:

- <code title="post /files">client.Files.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileNewParams">FileNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileObject">FileObject</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /files/{file_id}">client.Files.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, fileID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileObject">FileObject</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /files">client.Files.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileListParams">FileListParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#CursorPage">CursorPage</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileObject">FileObject</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="delete /files/{file_id}">client.Files.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileService.Delete">Delete</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, fileID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileDeleted">FileDeleted</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /files/{file_id}/content">client.Files.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileService.Content">Content</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, fileID <a href="https://pkg.go.dev/builtin#string">string</a>) (http.Response, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Images

Params Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageModel">ImageModel</a>

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Image">Image</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImagesResponse">ImagesResponse</a>

Methods:

- <code title="post /images/variations">client.Images.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageService.NewVariation">NewVariation</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageNewVariationParams">ImageNewVariationParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImagesResponse">ImagesResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /images/edits">client.Images.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageService.Edit">Edit</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageEditParams">ImageEditParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImagesResponse">ImagesResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /images/generations">client.Images.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageService.Generate">Generate</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageGenerateParams">ImageGenerateParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImagesResponse">ImagesResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Audio

Params Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AudioModel">AudioModel</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AudioResponseFormat">AudioResponseFormat</a>

## Transcriptions

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Transcription">Transcription</a>

Methods:

- <code title="post /audio/transcriptions">client.Audio.Transcriptions.<a href="https://pkg.go.dev/github.com/openai/openai-go#AudioTranscriptionService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AudioTranscriptionNewParams">AudioTranscriptionNewParams</a>) (Transcription, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Translations

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Translation">Translation</a>

Methods:

- <code title="post /audio/translations">client.Audio.Translations.<a href="https://pkg.go.dev/github.com/openai/openai-go#AudioTranslationService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AudioTranslationNewParams">AudioTranslationNewParams</a>) (Translation, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Speech

Params Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#SpeechModel">SpeechModel</a>

Methods:

- <code title="post /audio/speech">client.Audio.Speech.<a href="https://pkg.go.dev/github.com/openai/openai-go#AudioSpeechService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AudioSpeechNewParams">AudioSpeechNewParams</a>) (http.Response, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Moderations

Params Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ModerationImageURLInputParam">ModerationImageURLInputParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ModerationModel">ModerationModel</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ModerationMultiModalInputUnionParam">ModerationMultiModalInputUnionParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ModerationTextInputParam">ModerationTextInputParam</a>

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Moderation">Moderation</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ModerationNewResponse">ModerationNewResponse</a>

Methods:

- <code title="post /moderations">client.Moderations.<a href="https://pkg.go.dev/github.com/openai/openai-go#ModerationService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ModerationNewParams">ModerationNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ModerationNewResponse">ModerationNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Models

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Model">Model</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ModelDeleted">ModelDeleted</a>

Methods:

- <code title="get /models/{model}">client.Models.<a href="https://pkg.go.dev/github.com/openai/openai-go#ModelService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, model <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Model">Model</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /models">client.Models.<a href="https://pkg.go.dev/github.com/openai/openai-go#ModelService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#Page">Page</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Model">Model</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="delete /models/{model}">client.Models.<a href="https://pkg.go.dev/github.com/openai/openai-go#ModelService.Delete">Delete</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, model <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ModelDeleted">ModelDeleted</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# FineTuning

## Jobs

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJob">FineTuningJob</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobEvent">FineTuningJobEvent</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobWandbIntegration">FineTuningJobWandbIntegration</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobWandbIntegrationObject">FineTuningJobWandbIntegrationObject</a>

Methods:

- <code title="post /fine_tuning/jobs">client.FineTuning.Jobs.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobNewParams">FineTuningJobNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJob">FineTuningJob</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /fine_tuning/jobs/{fine_tuning_job_id}">client.FineTuning.Jobs.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, fineTuningJobID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJob">FineTuningJob</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /fine_tuning/jobs">client.FineTuning.Jobs.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobListParams">FineTuningJobListParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#CursorPage">CursorPage</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJob">FineTuningJob</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /fine_tuning/jobs/{fine_tuning_job_id}/cancel">client.FineTuning.Jobs.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobService.Cancel">Cancel</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, fineTuningJobID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJob">FineTuningJob</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /fine_tuning/jobs/{fine_tuning_job_id}/events">client.FineTuning.Jobs.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobService.ListEvents">ListEvents</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, fineTuningJobID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobListEventsParams">FineTuningJobListEventsParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#CursorPage">CursorPage</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobEvent">FineTuningJobEvent</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

### Checkpoints

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobCheckpoint">FineTuningJobCheckpoint</a>

Methods:

- <code title="get /fine_tuning/jobs/{fine_tuning_job_id}/checkpoints">client.FineTuning.Jobs.Checkpoints.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobCheckpointService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, fineTuningJobID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobCheckpointListParams">FineTuningJobCheckpointListParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#CursorPage">CursorPage</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FineTuningJobCheckpoint">FineTuningJobCheckpoint</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Beta

## VectorStores

Params Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AutoFileChunkingStrategyParam">AutoFileChunkingStrategyParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileChunkingStrategyParamUnion">FileChunkingStrategyParamUnion</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#StaticFileChunkingStrategyParam">StaticFileChunkingStrategyParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#StaticFileChunkingStrategyObjectParam">StaticFileChunkingStrategyObjectParam</a>

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileChunkingStrategy">FileChunkingStrategy</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#OtherFileChunkingStrategyObject">OtherFileChunkingStrategyObject</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#StaticFileChunkingStrategy">StaticFileChunkingStrategy</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#StaticFileChunkingStrategyObject">StaticFileChunkingStrategyObject</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStore">VectorStore</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStoreDeleted">VectorStoreDeleted</a>

Methods:

- <code title="post /vector_stores">client.Beta.VectorStores.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreNewParams">BetaVectorStoreNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStore">VectorStore</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /vector_stores/{vector_store_id}">client.Beta.VectorStores.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, vectorStoreID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStore">VectorStore</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /vector_stores/{vector_store_id}">client.Beta.VectorStores.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreService.Update">Update</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, vectorStoreID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreUpdateParams">BetaVectorStoreUpdateParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStore">VectorStore</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /vector_stores">client.Beta.VectorStores.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreListParams">BetaVectorStoreListParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#CursorPage">CursorPage</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStore">VectorStore</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="delete /vector_stores/{vector_store_id}">client.Beta.VectorStores.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreService.Delete">Delete</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, vectorStoreID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStoreDeleted">VectorStoreDeleted</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

### Files

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStoreFile">VectorStoreFile</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStoreFileDeleted">VectorStoreFileDeleted</a>

Methods:

- <code title="post /vector_stores/{vector_store_id}/files">client.Beta.VectorStores.Files.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreFileService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, vectorStoreID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreFileNewParams">BetaVectorStoreFileNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStoreFile">VectorStoreFile</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /vector_stores/{vector_store_id}/files/{file_id}">client.Beta.VectorStores.Files.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreFileService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, vectorStoreID <a href="https://pkg.go.dev/builtin#string">string</a>, fileID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStoreFile">VectorStoreFile</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /vector_stores/{vector_store_id}/files">client.Beta.VectorStores.Files.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreFileService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, vectorStoreID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreFileListParams">BetaVectorStoreFileListParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#CursorPage">CursorPage</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStoreFile">VectorStoreFile</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="delete /vector_stores/{vector_store_id}/files/{file_id}">client.Beta.VectorStores.Files.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreFileService.Delete">Delete</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, vectorStoreID <a href="https://pkg.go.dev/builtin#string">string</a>, fileID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStoreFileDeleted">VectorStoreFileDeleted</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

### FileBatches

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStoreFileBatch">VectorStoreFileBatch</a>

Methods:

- <code title="post /vector_stores/{vector_store_id}/file_batches">client.Beta.VectorStores.FileBatches.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreFileBatchService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, vectorStoreID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreFileBatchNewParams">BetaVectorStoreFileBatchNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStoreFileBatch">VectorStoreFileBatch</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /vector_stores/{vector_store_id}/file_batches/{batch_id}">client.Beta.VectorStores.FileBatches.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreFileBatchService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, vectorStoreID <a href="https://pkg.go.dev/builtin#string">string</a>, batchID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStoreFileBatch">VectorStoreFileBatch</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /vector_stores/{vector_store_id}/file_batches/{batch_id}/cancel">client.Beta.VectorStores.FileBatches.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreFileBatchService.Cancel">Cancel</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, vectorStoreID <a href="https://pkg.go.dev/builtin#string">string</a>, batchID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStoreFileBatch">VectorStoreFileBatch</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /vector_stores/{vector_store_id}/file_batches/{batch_id}/files">client.Beta.VectorStores.FileBatches.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreFileBatchService.ListFiles">ListFiles</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, vectorStoreID <a href="https://pkg.go.dev/builtin#string">string</a>, batchID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaVectorStoreFileBatchListFilesParams">BetaVectorStoreFileBatchListFilesParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#CursorPage">CursorPage</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#VectorStoreFile">VectorStoreFile</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Assistants

Params Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AssistantToolUnionParam">AssistantToolUnionParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#CodeInterpreterToolParam">CodeInterpreterToolParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileSearchToolParam">FileSearchToolParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FunctionToolParam">FunctionToolParam</a>

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Assistant">Assistant</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AssistantDeleted">AssistantDeleted</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AssistantStreamEvent">AssistantStreamEvent</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AssistantTool">AssistantTool</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#CodeInterpreterTool">CodeInterpreterTool</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileSearchTool">FileSearchTool</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FunctionTool">FunctionTool</a>

Methods:

- <code title="post /assistants">client.Beta.Assistants.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaAssistantService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaAssistantNewParams">BetaAssistantNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Assistant">Assistant</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /assistants/{assistant_id}">client.Beta.Assistants.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaAssistantService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, assistantID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Assistant">Assistant</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /assistants/{assistant_id}">client.Beta.Assistants.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaAssistantService.Update">Update</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, assistantID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaAssistantUpdateParams">BetaAssistantUpdateParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Assistant">Assistant</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /assistants">client.Beta.Assistants.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaAssistantService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaAssistantListParams">BetaAssistantListParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#CursorPage">CursorPage</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Assistant">Assistant</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="delete /assistants/{assistant_id}">client.Beta.Assistants.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaAssistantService.Delete">Delete</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, assistantID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AssistantDeleted">AssistantDeleted</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Threads

Params Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AssistantToolChoiceParam">AssistantToolChoiceParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AssistantToolChoiceFunctionParam">AssistantToolChoiceFunctionParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AssistantToolChoiceOptionUnionParam">AssistantToolChoiceOptionUnionParam</a>

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AssistantToolChoice">AssistantToolChoice</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AssistantToolChoiceFunction">AssistantToolChoiceFunction</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AssistantToolChoiceOptionUnion">AssistantToolChoiceOptionUnion</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Thread">Thread</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ThreadDeleted">ThreadDeleted</a>

Methods:

- <code title="post /threads">client.Beta.Threads.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadNewParams">BetaThreadNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Thread">Thread</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /threads/{thread_id}">client.Beta.Threads.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Thread">Thread</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /threads/{thread_id}">client.Beta.Threads.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadService.Update">Update</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadUpdateParams">BetaThreadUpdateParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Thread">Thread</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="delete /threads/{thread_id}">client.Beta.Threads.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadService.Delete">Delete</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ThreadDeleted">ThreadDeleted</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /threads/runs">client.Beta.Threads.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadService.NewAndRun">NewAndRun</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadNewAndRunParams">BetaThreadNewAndRunParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Run">Run</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

### Runs

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#RequiredActionFunctionToolCall">RequiredActionFunctionToolCall</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Run">Run</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#RunStatus">RunStatus</a>

Methods:

- <code title="post /threads/{thread_id}/runs">client.Beta.Threads.Runs.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadRunService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadRunNewParams">BetaThreadRunNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Run">Run</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /threads/{thread_id}/runs/{run_id}">client.Beta.Threads.Runs.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadRunService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>, runID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Run">Run</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /threads/{thread_id}/runs/{run_id}">client.Beta.Threads.Runs.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadRunService.Update">Update</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>, runID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadRunUpdateParams">BetaThreadRunUpdateParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Run">Run</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /threads/{thread_id}/runs">client.Beta.Threads.Runs.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadRunService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadRunListParams">BetaThreadRunListParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#CursorPage">CursorPage</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Run">Run</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /threads/{thread_id}/runs/{run_id}/cancel">client.Beta.Threads.Runs.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadRunService.Cancel">Cancel</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>, runID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Run">Run</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /threads/{thread_id}/runs/{run_id}/submit_tool_outputs">client.Beta.Threads.Runs.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadRunService.SubmitToolOutputs">SubmitToolOutputs</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>, runID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadRunSubmitToolOutputsParams">BetaThreadRunSubmitToolOutputsParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Run">Run</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

#### Steps

Params Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#RunStepInclude">RunStepInclude</a>

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#CodeInterpreterLogs">CodeInterpreterLogs</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#CodeInterpreterOutputImage">CodeInterpreterOutputImage</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#CodeInterpreterToolCall">CodeInterpreterToolCall</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#CodeInterpreterToolCallDelta">CodeInterpreterToolCallDelta</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileSearchToolCall">FileSearchToolCall</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileSearchToolCallDelta">FileSearchToolCallDelta</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FunctionToolCall">FunctionToolCall</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FunctionToolCallDelta">FunctionToolCallDelta</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#MessageCreationStepDetails">MessageCreationStepDetails</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#RunStep">RunStep</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#RunStepDelta">RunStepDelta</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#RunStepDeltaEvent">RunStepDeltaEvent</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#RunStepDeltaMessageDelta">RunStepDeltaMessageDelta</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ToolCall">ToolCall</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ToolCallDelta">ToolCallDelta</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ToolCallDeltaObject">ToolCallDeltaObject</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ToolCallsStepDetails">ToolCallsStepDetails</a>

Methods:

- <code title="get /threads/{thread_id}/runs/{run_id}/steps/{step_id}">client.Beta.Threads.Runs.Steps.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadRunStepService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>, runID <a href="https://pkg.go.dev/builtin#string">string</a>, stepID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadRunStepGetParams">BetaThreadRunStepGetParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#RunStep">RunStep</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /threads/{thread_id}/runs/{run_id}/steps">client.Beta.Threads.Runs.Steps.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadRunStepService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>, runID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadRunStepListParams">BetaThreadRunStepListParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#CursorPage">CursorPage</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#RunStep">RunStep</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

### Messages

Params Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageFileParam">ImageFileParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageFileContentBlockParam">ImageFileContentBlockParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageURLParam">ImageURLParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageURLContentBlockParam">ImageURLContentBlockParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#MessageContentPartParamUnion">MessageContentPartParamUnion</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#TextContentBlockParam">TextContentBlockParam</a>

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Annotation">Annotation</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#AnnotationDelta">AnnotationDelta</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileCitationAnnotation">FileCitationAnnotation</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FileCitationDeltaAnnotation">FileCitationDeltaAnnotation</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FilePathAnnotation">FilePathAnnotation</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#FilePathDeltaAnnotation">FilePathDeltaAnnotation</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageFile">ImageFile</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageFileContentBlock">ImageFileContentBlock</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageFileDelta">ImageFileDelta</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageFileDeltaBlock">ImageFileDeltaBlock</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageURL">ImageURL</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageURLContentBlock">ImageURLContentBlock</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageURLDelta">ImageURLDelta</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#ImageURLDeltaBlock">ImageURLDeltaBlock</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Message">Message</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#MessageContent">MessageContent</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#MessageContentDelta">MessageContentDelta</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#MessageDeleted">MessageDeleted</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#MessageDelta">MessageDelta</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#MessageDeltaEvent">MessageDeltaEvent</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#RefusalContentBlock">RefusalContentBlock</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#RefusalDeltaBlock">RefusalDeltaBlock</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Text">Text</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#TextContentBlock">TextContentBlock</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#TextDelta">TextDelta</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#TextDeltaBlock">TextDeltaBlock</a>

Methods:

- <code title="post /threads/{thread_id}/messages">client.Beta.Threads.Messages.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadMessageService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadMessageNewParams">BetaThreadMessageNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Message">Message</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /threads/{thread_id}/messages/{message_id}">client.Beta.Threads.Messages.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadMessageService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>, messageID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Message">Message</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /threads/{thread_id}/messages/{message_id}">client.Beta.Threads.Messages.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadMessageService.Update">Update</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>, messageID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadMessageUpdateParams">BetaThreadMessageUpdateParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Message">Message</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /threads/{thread_id}/messages">client.Beta.Threads.Messages.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadMessageService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadMessageListParams">BetaThreadMessageListParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#CursorPage">CursorPage</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Message">Message</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="delete /threads/{thread_id}/messages/{message_id}">client.Beta.Threads.Messages.<a href="https://pkg.go.dev/github.com/openai/openai-go#BetaThreadMessageService.Delete">Delete</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, threadID <a href="https://pkg.go.dev/builtin#string">string</a>, messageID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#MessageDeleted">MessageDeleted</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Batches

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Batch">Batch</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BatchError">BatchError</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BatchRequestCounts">BatchRequestCounts</a>

Methods:

- <code title="post /batches">client.Batches.<a href="https://pkg.go.dev/github.com/openai/openai-go#BatchService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BatchNewParams">BatchNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Batch">Batch</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /batches/{batch_id}">client.Batches.<a href="https://pkg.go.dev/github.com/openai/openai-go#BatchService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, batchID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Batch">Batch</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /batches">client.Batches.<a href="https://pkg.go.dev/github.com/openai/openai-go#BatchService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, query <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#BatchListParams">BatchListParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination">pagination</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/packages/pagination#CursorPage">CursorPage</a>[<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Batch">Batch</a>], <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /batches/{batch_id}/cancel">client.Batches.<a href="https://pkg.go.dev/github.com/openai/openai-go#BatchService.Cancel">Cancel</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, batchID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Batch">Batch</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Uploads

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Upload">Upload</a>

Methods:

- <code title="post /uploads">client.Uploads.<a href="https://pkg.go.dev/github.com/openai/openai-go#UploadService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#UploadNewParams">UploadNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Upload">Upload</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /uploads/{upload_id}/cancel">client.Uploads.<a href="https://pkg.go.dev/github.com/openai/openai-go#UploadService.Cancel">Cancel</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, uploadID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Upload">Upload</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /uploads/{upload_id}/complete">client.Uploads.<a href="https://pkg.go.dev/github.com/openai/openai-go#UploadService.Complete">Complete</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, uploadID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#UploadCompleteParams">UploadCompleteParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#Upload">Upload</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Parts

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#UploadPart">UploadPart</a>

Methods:

- <code title="post /uploads/{upload_id}/parts">client.Uploads.Parts.<a href="https://pkg.go.dev/github.com/openai/openai-go#UploadPartService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, uploadID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#UploadPartNewParams">UploadPartNewParams</a>) (<a href="https://pkg.go.dev/github.com/openai/openai-go">openai</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go#UploadPart">UploadPart</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
