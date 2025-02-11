package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"os"
)

const (
	_prodMode = "prod"
)

type LoggerConfig struct {
	modeLogger string `json:"mode"`
}

func LoadLoggerConfig() (*LoggerConfig, error) {
	var cfg struct {
		Logger LoggerConfig `json:"logger"`
	}

	if err := cleanenv.ReadConfig("config.json", &cfg); err != nil {
		return nil, err
	}
	return &cfg.Logger, nil
}

func NewLogger(cfg *LoggerConfig) *slog.Logger {
	switch cfg.modeLogger {
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case _prodMode:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
}
