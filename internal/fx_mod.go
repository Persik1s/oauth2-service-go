package internal

import (
	"github.com/Persik1s/oauth2-service-go/internal/adapter/postgres"
	"github.com/Persik1s/oauth2-service-go/internal/adapter/postgres/database"
	"github.com/Persik1s/oauth2-service-go/internal/adapter/redis"
	v1 "github.com/Persik1s/oauth2-service-go/internal/delivery/rest/v1"
	"github.com/Persik1s/oauth2-service-go/internal/usecase"
	"github.com/Persik1s/oauth2-service-go/pkg/config"
	"github.com/Persik1s/oauth2-service-go/pkg/util/token"
	"go.uber.org/fx"
)

var ModuleApp = fx.Module(
	"app",
	config.ModuleConfig,

	fx.Provide(token.NewTokenManager),

	fx.Provide(
		config.NewPostgres,
		fx.Annotate(
			config.NewPostgres,
			fx.As(new(database.DBTX)),
		),
	),
	fx.Provide(database.New),

	fx.Provide(
		postgres.NewUserAdapter,
		fx.Annotate(
			postgres.NewUserAdapter,
			fx.As(new(usecase.UserAdapterIn)),
		),
	),

	fx.Provide(
		postgres.NewRoleAdapter,
		fx.Annotate(
			postgres.NewRoleAdapter,
			fx.As(new(usecase.RoleAdapterIn)),
		),
	),

	fx.Provide(
		redis.NewSessionAdapter,
		fx.Annotate(
			redis.NewSessionAdapter,
			fx.As(new(usecase.SessionAdapterIn)),
		),
	),

	fx.Provide(
		usecase.NewAuthUsecase,
		fx.Annotate(
			usecase.NewAuthUsecase,
			fx.As(new(usecase.AuthUsecaseIn)),
		),
	),

	fx.Provide(
		usecase.NewOAuthUsecase,
		fx.Annotate(
			usecase.NewOAuthUsecase,
			fx.As(new(usecase.OAuthUsecaseIn)),
		),
	),

	fx.Provide(v1.NewDelivery),

	fx.Invoke(v1.NewAuthDelivery, v1.NewOAuthDelivery, config.StartPostgres, config.StartRedis, config.StartServer),
)
