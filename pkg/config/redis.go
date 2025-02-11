package config

import (
	"context"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"log/slog"
)

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func LoadRedisConfig() (*RedisConfig, error) {
	var cfg struct {
		Redis RedisConfig `json:"redis"`
	}

	if err := cleanenv.ReadConfig("config.json", &cfg); err != nil {
		return nil, err
	}

	return &cfg.Redis, nil
}

func NewRedis(cfg *RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       0,
	})
}

func StartRedis(lc fx.Lifecycle, redis *redis.Client, logger *slog.Logger) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := redis.Ping(ctx).Err(); err != nil {
					logger.Error("redis.Ping(ctx)", "err", err)
					return err
				}
				return nil
			},
			//OnStop: func(ctx context.Context) error {
			//	if err := redis.Shutdown(ctx).Err(); err != nil {
			//		logger.Error("redis.Shutdown(ctx)", "err", err)
			//		return err
			//	}
			//	return nil
			//},
		},
	)
}
