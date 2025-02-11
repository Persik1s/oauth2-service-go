package config

import (
	"context"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5"
	"go.uber.org/fx"
	"log/slog"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func (p *PostgresConfig) Addr() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", p.User, p.Password, p.Host, p.Port, p.Database)
}

func LoadPostgresConfig() (*PostgresConfig, error) {
	var cfg struct {
		Postgres PostgresConfig `json:"postgres"`
	}
	if err := cleanenv.ReadConfig("config.json", &cfg); err != nil {
		return nil, err
	}
	return &cfg.Postgres, nil
}

func NewPostgres(cfg *PostgresConfig) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), cfg.Addr())
}

func StartPostgres(lc fx.Lifecycle, connection *pgx.Conn, logger *slog.Logger) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := connection.Ping(ctx); err != nil {
					logger.Error("connection.Ping", "err", err)
					return err
				}
				return nil
			},

			OnStop: func(ctx context.Context) error {
				if err := connection.Close(ctx); err != nil {
					logger.Error("connection.Close", "err", err)
					return err
				}
				return nil
			},
		},
	)
}
