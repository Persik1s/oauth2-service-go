package postgres

import (
	"context"
	"errors"
	"github.com/Persik1s/oauth2-service-go/internal/adapter/postgres/database"
	"github.com/Persik1s/oauth2-service-go/internal/errorz"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
)

type RoleAdapter struct {
	logger *slog.Logger
	query  *database.Queries
}

func NewRoleAdapter(query *database.Queries, logger *slog.Logger) *RoleAdapter {
	return &RoleAdapter{
		logger: logger,
		query:  query,
	}
}

func (r *RoleAdapter) GetRole(ctx context.Context, name string) (database.Role, error) {
	role, err := r.query.GetRole(ctx, name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UndefinedColumn {
				return database.Role{}, errorz.ErrRoleNotFound
			}
		}
		return database.Role{}, err
	}
	return role, nil
}
