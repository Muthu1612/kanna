package llm

import "context"

type Request struct {
	Prompt string
}

type Response struct {
	Content string
}

type Client interface {
	Generate(ctx context.Context, request Request) (Response, error)
}
