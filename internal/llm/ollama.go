package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type OllamaClient struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewOllamaClient(baseURL, model string) *OllamaClient {
	return &OllamaClient{
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type ollamaMessage struct {
	Role      string           `json:"role"`
	Content   string           `json:"content,omitempty"`
	ToolCalls []ollamaToolCall `json:"tool_calls,omitempty"`
}

type ollamaToolCall struct {
	Function ollamaToolCallFunction `json:"function"`
}

type ollamaToolCallFunction struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type ollamaTool struct {
	Type     string         `json:"type"`
	Function ollamaFunction `json:"function"`
}

type ollamaFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Tools    []ollamaTool    `json:"tools,omitempty"`
	Stream   bool            `json:"stream"`
}

type ollamaResponse struct {
	Message ollamaMessage `json:"message"`
}

func (c *OllamaClient) Generate(
	ctx context.Context,
	request Request,
) (Response, error) {
	messages := make([]ollamaMessage, 0, len(request.Messages))

	for _, message := range request.Messages {
		messages = append(messages, ollamaMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	ollamaTools := make([]ollamaTool, 0, len(request.Tools))

	for _, tool := range request.Tools {
		ollamaTools = append(ollamaTools, ollamaTool{
			Type: "function",
			Function: ollamaFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.Parameters,
			},
		})
	}

	payload := ollamaRequest{
		Model:    c.model,
		Messages: messages,
		Tools:    ollamaTools,
		Stream:   false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf(
			"marshal Ollama request: %w",
			err,
		)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/api/chat",
		bytes.NewReader(body),
	)
	if err != nil {
		return Response{}, fmt.Errorf(
			"create Ollama request: %w",
			err,
		)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf(
			"call Ollama: %w",
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf(
			"Ollama returned status %d",
			resp.StatusCode,
		)
	}

	var result ollamaResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Response{}, fmt.Errorf(
			"decode Ollama response: %w",
			err,
		)
	}

	response := Response{
		Content: result.Message.Content,
	}

	// Native Ollama tool call.
	if len(result.Message.ToolCalls) > 0 {
		call := result.Message.ToolCalls[0]

		response.ToolCall = &ToolCall{
			Name:  call.Function.Name,
			Input: call.Function.Arguments,
		}

		return response, nil
	}

	// Some models return the tool call as JSON text.
	content := strings.TrimSpace(result.Message.Content)

	var toolCall struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}

	if err := json.Unmarshal([]byte(content), &toolCall); err == nil {
		if toolCall.Name != "" {
			response.ToolCall = &ToolCall{
				Name:  toolCall.Name,
				Input: toolCall.Arguments,
			}

			response.Content = ""

			return response, nil
		}
	}

	return response, nil
}
