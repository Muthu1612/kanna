package tools

import "context"

type Definition struct {
	Name        string
	Description string
	Parameters  map[string]any
}

type Input map[string]any

type Output struct {
	Content string
}

type Tool interface {
	Definition() Definition
	Execute(ctx context.Context, input Input) (Output, error)
}
