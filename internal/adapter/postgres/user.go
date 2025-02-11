package postgres

import (
	"context"
	"errors"
	"github.com/Persik1s/oauth2-service-go/internal/adapter/postgres/database"
	"github.com/Persik1s/oauth2-service-go/internal/errorz"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
)

type UserAdapter struct {
	logger *slog.Logger
	query  *database.Queries
}

func NewUserAdapter(query *database.Queries, logger *slog.Logger) *UserAdapter {
	return &UserAdapter{
		logger: logger,
		query:  query,
	}
}

func (u *UserAdapter) CreateUser(ctx context.Context, user database.User) (uuid.UUID, error) {
	id, err := u.query.CreateUser(ctx, database.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Createat: user.Createat,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return uuid.Nil, errorz.ErrUserAlreadyExists
			}
		}
		return uuid.Nil, err
	}
	return id, nil
}

func (u *UserAdapter) GetUser(ctx context.Context, username string) (database.User, error) {
	user, err := u.query.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.User{}, errorz.ErrUserNotFound
		}
		return database.User{}, err
	}

	return user, nil
}

func (u *UserAdapter) CreateUserRole(ctx context.Context, role database.UserRole) error {
	return u.query.CreateUserRole(ctx, database.CreateUserRoleParams{
		RoleID: role.RoleID,
		UserID: role.UserID,
	})
}

func (u *UserAdapter) GetUserRole(ctx context.Context, userId uuid.UUID) (database.UserRole, error) {
	return u.query.GetUserRole(ctx, userId)
}
