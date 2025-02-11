package config

import (
	"context"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"log/slog"
)

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (s *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func LoadServerConfig() (*ServerConfig, error) {
	var cfg struct {
		Server ServerConfig `json:"server"`
	}

	if err := cleanenv.ReadConfig("config.json", &cfg); err != nil {
		return nil, err
	}

	return &cfg.Server, nil
}

func NewServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		HandleError:      true,
		LogLatency:       true,
		LogProtocol:      true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogURIPath:       true,
		LogRoutePath:     true,
		LogRequestID:     true,
		LogReferer:       true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			return nil
		},
	}))

	return e
}

func StartServer(lc fx.Lifecycle, ec *echo.Echo, cfg *ServerConfig, logger *slog.Logger) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := ec.Start(cfg.Addr()); err != nil {
						logger.Error("start server", "err", err)
						return
					}
				}()

				return nil
			},

			OnStop: func(ctx context.Context) error {
				go func() {
					if err := ec.Shutdown(ctx); err != nil {
						logger.Error("ec shutdown", "err", err)
						return
					}
				}()
				return nil
			},
		},
	)
}
