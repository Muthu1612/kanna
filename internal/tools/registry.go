package tools

import (
	"context"
	"fmt"
)

type Registry interface {
	Register(tool Tool) error
	Get(name string) (Tool, error)
	List() []Definition
	Execute(ctx context.Context, name string, input Input) (Output, error)
}

type registry struct {
	tools map[string]Tool
}

func NewRegistry() Registry {
	return &registry{
		tools: make(map[string]Tool),
	}
}

func (r *registry) Register(tool Tool) error {
	name := tool.Definition().Name

	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool already registered: %s", name)
	}

	r.tools[name] = tool

	return nil
}

func (r *registry) Get(name string) (Tool, error) {
	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	return tool, nil
}

func (r *registry) List() []Definition {
	definitions := make([]Definition, 0, len(r.tools))

	for _, tool := range r.tools {
		definitions = append(definitions, tool.Definition())
	}

	return definitions
}

func (r *registry) Execute(
	ctx context.Context,
	name string,
	input Input,
) (Output, error) {
	tool, err := r.Get(name)
	if err != nil {
		return Output{}, err
	}

	return tool.Execute(ctx, input)
}
