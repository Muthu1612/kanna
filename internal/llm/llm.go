package llm

import "context"

type Message struct {
	Role    string
	Content string
}

type ToolDefinition struct {
	Name        string
	Description string
	Parameters  map[string]any
}

type Request struct {
	Messages []Message
	Tools    []ToolDefinition
}

type ToolCall struct {
	Name  string
	Input map[string]any
}

type Response struct {
	Content  string
	ToolCall *ToolCall
}

type Client interface {
	Generate(ctx context.Context, request Request) (Response, error)
}
