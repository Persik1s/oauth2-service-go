package usecase

import (
	"context"
	"github.com/Persik1s/oauth2-service-go/internal/adapter/postgres/database"
	"github.com/Persik1s/oauth2-service-go/internal/domain"
	"github.com/Persik1s/oauth2-service-go/internal/dto"
	"github.com/Persik1s/oauth2-service-go/internal/errorz"
	"github.com/Persik1s/oauth2-service-go/pkg/config"
	"github.com/Persik1s/oauth2-service-go/pkg/util/token"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type UserAdapterIn interface {
	CreateUser(ctx context.Context, user database.User) (uuid.UUID, error)
	GetUser(ctx context.Context, username string) (database.User, error)

	CreateUserRole(ctx context.Context, role database.UserRole) error
	GetUserRole(ctx context.Context, userId uuid.UUID) (database.UserRole, error)
}

type RoleAdapterIn interface {
	GetRole(ctx context.Context, name string) (database.Role, error)
}

type AuthUsecaseIn interface {
	Registration(ctx context.Context, data dto.SignUpRequestDto) (dto.SignUpResponeDto, error)
	Authorization(ctx context.Context, requestDto dto.SignInRequestDto) (dto.TokenPair, error)
}

type SessionAdapterIn interface {
	CreateSession(ctx context.Context, session domain.Session) error
}

type AuthUsecase struct {
	logger *slog.Logger

	userAdapter    UserAdapterIn
	roleAdapter    RoleAdapterIn
	sessionAdapter SessionAdapterIn

	applicationConfig *config.ApplicationConfig

	tokenManager *token.TokenManager
	cfgToken     *config.TokenConfig
}

func NewAuthUsecase(sessionAdapter SessionAdapterIn,
	userAdapter UserAdapterIn,
	roleAdapter RoleAdapterIn,
	cfgToken *config.TokenConfig,
	tokenManager *token.TokenManager,
	applicationConfig *config.ApplicationConfig,
	logger *slog.Logger) *AuthUsecase {
	return &AuthUsecase{
		logger:            logger,
		userAdapter:       userAdapter,
		applicationConfig: applicationConfig,
		roleAdapter:       roleAdapter,
		tokenManager:      tokenManager,
		sessionAdapter:    sessionAdapter,
		cfgToken:          cfgToken,
	}
}

func (a *AuthUsecase) Registration(ctx context.Context, data dto.SignUpRequestDto) (dto.SignUpResponeDto, error) {
	// create role

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Debug("bcrypt.GenerateFromPassword", "err", err)
		return dto.SignUpResponeDto{}, err
	}

	userId, err := a.userAdapter.CreateUser(ctx, database.User{
		Username: data.Username,
		Email:    data.Email,
		Password: passwordHash,
		Createat: time.Now(),
	})

	roleUser, err := a.roleAdapter.GetRole(ctx, domain.UserRole)
	if err != nil {
		a.logger.Debug("a.roleAdapter.GetRole", "err", err)
		return dto.SignUpResponeDto{}, err
	}

	if err := a.userAdapter.CreateUserRole(ctx, database.UserRole{
		UserID: userId,
		RoleID: roleUser.ID,
	}); err != nil {
		a.logger.Debug("a.userAdapter.CreateUserRole", "err")
		return dto.SignUpResponeDto{}, err
	}

	return dto.SignUpResponeDto{
		a.applicationConfig.ClientId,
		userId,
	}, nil
}

func (a *AuthUsecase) Authorization(ctx context.Context, requestDto dto.SignInRequestDto) (dto.TokenPair, error) {
	user, err := a.userAdapter.GetUser(ctx, requestDto.Username)
	if err != nil {
		a.logger.Debug("a.userAdapter.GetUser", "err", err)
		return dto.TokenPair{}, err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(requestDto.Password)); err != nil {
		a.logger.Debug("bcrypt.CompareHashAndPassword", "err", err)
		return dto.TokenPair{}, errorz.ErrPasswordNotValid
	}

	userRole, err := a.userAdapter.GetUserRole(ctx, user.ID)
	if err != nil {
		a.logger.Debug("a.userAdapter.GetUserRole", "err", err)
		return dto.TokenPair{}, err
	}

	sessionToken, err := a.tokenManager.GenerateSessionToken(userRole.UserID)
	if err != nil {
		a.logger.Debug("a.tokenManager.GenerateSessionToken", "err", err)
		return dto.TokenPair{}, err
	}

	if err := a.sessionAdapter.CreateSession(ctx, domain.Session{
		UserId:       user.ID,
		SessionToken: sessionToken,
		TTL:          a.cfgToken.SessionExp * time.Second,
	}); err != nil {
		a.logger.Debug("a.sessionAdapter.CreateSession", "err", err)
		return dto.TokenPair{}, err
	}

	accessToken, err := a.tokenManager.GenerateAccessToken(user.ID, userRole.RoleID)
	if err != nil {
		a.logger.Debug("a.tokenManager.GenerateAccessToken", "err", err)
		return dto.TokenPair{}, err
	}

	refreshToken, err := a.tokenManager.GenerateRefreshToken(sessionToken)
	if err != nil {
		a.logger.Debug("a.tokenManager.GenerateRefreshToken", "err", err)
		return dto.TokenPair{}, err
	}

	return dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
