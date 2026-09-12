package tools

import (
	"context"
	"time"
)

type TimeTool struct{}

func NewTimeTool() *TimeTool {
	return &TimeTool{}
}

func (t *TimeTool) Definition() Definition {
	return Definition{
		Name:        "get_current_time",
		Description: "Returns the current system time.",
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
	}
}

func (t *TimeTool) Execute(
	ctx context.Context,
	input Input,
) (Output, error) {
	return Output{
		Content: time.Now().Format(time.RFC1123),
	}, nil
}
