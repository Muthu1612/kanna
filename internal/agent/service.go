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

	tools := make([]llm.ToolDefinition, 0, len(toolDefinitions))

	for _, tool := range toolDefinitions {
		tools = append(tools, llm.ToolDefinition{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  tool.Parameters,
		})
	}

	messages := []llm.Message{
		{
			Role:    "user",
			Content: request.Message,
		},
	}

	result, err := s.llm.Generate(
		ctx,
		llm.Request{
			Messages: messages,
			Tools:    tools,
		},
	)
	if err != nil {
		return Response{}, fmt.Errorf("initial LLM generation failed: %w", err)
	}

	// No tool requested.
	if result.ToolCall == nil {
		return Response{
			Content: result.Content,
		}, nil
	}

	// Execute requested tool.
	toolResult, err := s.toolRegistry.Execute(
		ctx,
		result.ToolCall.Name,
		result.ToolCall.Input,
	)
	if err != nil {
		return Response{}, fmt.Errorf("execute tool: %w", err)
	}

	// Give the tool result back to the LLM.
	messages = append(messages,
		llm.Message{
			Role:    "assistant",
			Content: "",
		},
		llm.Message{
			Role:    "tool",
			Content: toolResult.Content,
		},
	)

	finalResult, err := s.llm.Generate(
		ctx,
		llm.Request{
			Messages: messages,
		},
	)
	if err != nil {
		return Response{}, fmt.Errorf(
			"final LLM generation failed: %w",
			err,
		)
	}

	return Response{
		Content: finalResult.Content,
	}, nil
}
