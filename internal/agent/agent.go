package agent

import "context"

type Request struct {
	Message string
}

type Response struct {
	Content string
}

type Agent interface {
	Run(ctx context.Context, request Request) (Response, error)
}
