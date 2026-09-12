package agent

import (
	"context"
	"fmt"

	"github.com/Muthu1612/kanna/internal/llm"
	"github.com/Muthu1612/kanna/internal/tools"
)

type service struct {
	llm          llm.Client
	toolRegistry tools.Registry
}

func NewService(
	llmClient llm.Client,
	toolRegistry tools.Registry,
) Agent {
	return &service{
		llm:          llmClient,
		toolRegistry: toolRegistry,
	}
}

func (s *service) Run(
	ctx context.Context,
	request Request,
) (Response, error) {

	toolDefinitions := s.toolRegistry.List()

	llmTools := make([]llm.ToolDefinition, 0, len(toolDefinitions))

	for _, tool := range toolDefinitions {
		llmTools = append(llmTools, llm.ToolDefinition{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  tool.Parameters,
		})
	}

	result, err := s.llm.Generate(
		ctx,
		llm.Request{
			Messages: []llm.Message{
				{
					Role:    "user",
					Content: request.Message,
				},
			},
			Tools: llmTools,
		},
	)
	if err != nil {
		return Response{}, fmt.Errorf(
			"agent generation failed: %w",
			err,
		)
	}

	if result.ToolCall != nil {
		output, err := s.toolRegistry.Execute(
			ctx,
			result.ToolCall.Name,
			tools.Input(result.ToolCall.Input),
		)
		if err != nil {
			return Response{}, fmt.Errorf(
				"tool execution failed: %w",
				err,
			)
		}

		return Response{
			Content: output.Content,
		}, nil
	}

	return Response{
		Content: result.Content,
	}, nil
}
