package memory

import "context"

type Conversation struct {
	Question string
	Answer   string
}

type Store interface {
	SaveConversation(ctx context.Context, conversation Conversation) error
}
