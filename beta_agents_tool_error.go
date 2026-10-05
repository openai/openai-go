package openai

// BetaAgentToolError describes a local tool failure observed by OnToolError.
// Err is the original error and may contain sensitive application information.
// It is never included in the generic failed tool result sent to the API.
// This diagnostic helper is beta.
type BetaAgentToolError struct {
	Err       error
	ToolName  string
	SessionID string
	TurnID    string
	CallID    string
	Stage     BetaAgentToolErrorStage
}

// BetaAgentToolErrorStage identifies the part of a tool invocation that failed.
type BetaAgentToolErrorStage string

const (
	BetaAgentToolErrorStageArguments BetaAgentToolErrorStage = "arguments"
	BetaAgentToolErrorStageExecution BetaAgentToolErrorStage = "execution"
	BetaAgentToolErrorStageOutput    BetaAgentToolErrorStage = "output"
)
