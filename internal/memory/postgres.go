package memory

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PostgresStore struct {
	conn *pgx.Conn
}

func NewPostgresStore(
	ctx context.Context,
	databaseURL string,
) (*PostgresStore, error) {
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		_ = conn.Close(ctx)

		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &PostgresStore{
		conn: conn,
	}, nil
}

func (s *PostgresStore) SaveConversation(
	ctx context.Context,
	conversation Conversation,
) error {
	_, err := s.conn.Exec(
		ctx,
		`
		INSERT INTO conversations (question, answer)
		VALUES ($1, $2)
		`,
		conversation.Question,
		conversation.Answer,
	)
	if err != nil {
		return fmt.Errorf("save conversation: %w", err)
	}

	return nil
}

func (s *PostgresStore) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}
