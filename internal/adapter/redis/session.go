package redis

import (
	"context"
	"github.com/Persik1s/oauth2-service-go/internal/domain"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type SessionAdapter struct {
	logger     *slog.Logger
	connection *redis.Client
}

func NewSessionAdapter(connection *redis.Client, logger *slog.Logger) *SessionAdapter {
	return &SessionAdapter{
		logger:     logger,
		connection: connection,
	}
}

func (s *SessionAdapter) CreateSession(ctx context.Context, session domain.Session) error {
	return s.connection.Set(ctx, session.UserId.String(), session.SessionToken, session.TTL).Err()
}
