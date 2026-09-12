# Live

Params Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#BuiltInVoice">BuiltInVoice</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#ClientConfigParam">ClientConfigParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#ClientDelegationParam">ClientDelegationParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#CustomVoiceParam">CustomVoiceParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#DataChannelConfigParam">DataChannelConfigParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#FunctionToolParam">FunctionToolParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#InitialItemUnionParam">InitialItemUnionParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#MediaSessionConfigParam">MediaSessionConfigParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#MediaSessionForkConfigParam">MediaSessionForkConfigParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#ResponsesDelegationConfigParam">ResponsesDelegationConfigParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#ResponsesDelegationUpdateConfigParam">ResponsesDelegationUpdateConfigParam</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#ServerEventSelectorParam">ServerEventSelectorParam</a>

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#BuiltInVoice">BuiltInVoice</a>
- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#LiveNewResponse">LiveNewResponse</a>

Methods:

- <code title="post /live/sessions">client.Live.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#LiveService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#LiveNewParams">LiveNewParams</a>) (\*<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#LiveNewResponse">LiveNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Sideband

## Forks

## Sessions

Response Types:

- <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#SessionForkResponse">SessionForkResponse</a>

Methods:

- <code title="post /live/sessions/{session_id}/accept">client.Live.Sessions.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#SessionService.Accept">Accept</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, sessionID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#SessionAcceptParams">SessionAcceptParams</a>) <a href="https://pkg.go.dev/builtin#error">error</a></code>
- <code title="get /live/sessions/{session_id}/content">client.Live.Sessions.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#SessionService.DownloadRecording">DownloadRecording</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, sessionID <a href="https://pkg.go.dev/builtin#string">string</a>) (\*http.Response, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /live/sessions/{session_id}/fork">client.Live.Sessions.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#SessionService.Fork">Fork</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, sessionID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#SessionForkParams">SessionForkParams</a>) (\*<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#SessionForkResponse">SessionForkResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /live/sessions/{session_id}/hangup">client.Live.Sessions.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#SessionService.Hangup">Hangup</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, sessionID <a href="https://pkg.go.dev/builtin#string">string</a>) <a href="https://pkg.go.dev/builtin#error">error</a></code>
- <code title="post /live/sessions/{session_id}/refer">client.Live.Sessions.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#SessionService.Refer">Refer</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, sessionID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#SessionReferParams">SessionReferParams</a>) <a href="https://pkg.go.dev/builtin#error">error</a></code>
- <code title="post /live/sessions/{session_id}/reject">client.Live.Sessions.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#SessionService.Reject">Reject</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, sessionID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live">live</a>.<a href="https://pkg.go.dev/github.com/openai/openai-go/v3/live#SessionRejectParams">SessionRejectParams</a>) <a href="https://pkg.go.dev/builtin#error">error</a></code>
