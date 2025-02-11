package config

import "go.uber.org/fx"

var ModuleConfig = fx.Module(
	"config",

	fx.Provide(
		LoadLoggerConfig, NewLogger,
	),

	fx.Provide(LoadPrivateKey, LoadPublicKey, LoadTokenConfig),

	fx.Provide(LoadServerConfig, NewServer),

	fx.Provide(LoadOAuthGoogleConfig),

	fx.Provide(LoadPostgresConfig),

	fx.Provide(LoadApplicationConfig),

	fx.Provide(LoadRedisConfig, NewRedis),

	fx.Provide(NewValidator),
)
