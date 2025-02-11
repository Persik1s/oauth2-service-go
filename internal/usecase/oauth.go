package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Persik1s/oauth2-service-go/internal/adapter/postgres/database"
	"github.com/Persik1s/oauth2-service-go/internal/domain"
	"github.com/Persik1s/oauth2-service-go/internal/dto"
	"github.com/Persik1s/oauth2-service-go/internal/errorz"
	"github.com/Persik1s/oauth2-service-go/pkg/config"
	"github.com/Persik1s/oauth2-service-go/pkg/util/token"
	"github.com/google/go-querystring/query"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type OAuthUsecaseIn interface {
	OAuthorization(ctx context.Context, code string) (dto.TokenPair, error)
}

type OAuthUsecase struct {
	logger *slog.Logger

	userAdapter    UserAdapterIn
	sessionAdapter SessionAdapterIn
	roleAdapter    RoleAdapterIn
	googleAuthCfg  *config.OAuthGoogleConfig
	tokenManager   *token.TokenManager
	tokenConfig    *config.TokenConfig
}

func NewOAuthUsecase(userAdapter UserAdapterIn,
	roleAdapter RoleAdapterIn,
	sessionAdapter SessionAdapterIn,
	tokenConfig *config.TokenConfig,
	tokenManager *token.TokenManager,
	googleAuthCfg *config.OAuthGoogleConfig,
	logger *slog.Logger) *OAuthUsecase {

	return &OAuthUsecase{
		logger:         logger,
		userAdapter:    userAdapter,
		googleAuthCfg:  googleAuthCfg,
		tokenManager:   tokenManager,
		tokenConfig:    tokenConfig,
		sessionAdapter: sessionAdapter,
		roleAdapter:    roleAdapter,
	}
}

func (o *OAuthUsecase) OAuthorization(ctx context.Context, code string) (dto.TokenPair, error) {
	if code == "" {
		o.logger.Debug("code is null")
		return dto.TokenPair{}, errors.New("code is null")
	}

	requestData, err := query.Values(dto.OGoogleTokenRequestDto{
		ClientId:     o.googleAuthCfg.ClientId,
		ClientSecret: o.googleAuthCfg.ClientSecret,
		RedirectUri:  o.googleAuthCfg.RedirectUri,
		GrantType:    "authorization_code",
		Code:         code,
	})

	if err != nil {
		o.logger.Debug("query.Values", "err", err)
		return dto.TokenPair{}, err
	}

	responeToken, err := http.PostForm("https://oauth2.googleapis.com/token", requestData)
	if err != nil {
		o.logger.Debug("http.PostForm", "err", err)
		return dto.TokenPair{}, err
	}

	bodyToken, err := io.ReadAll(responeToken.Body)
	if err != nil {
		o.logger.Debug("io.ReadAll", "err", err)
		return dto.TokenPair{}, err
	}
	otoken := dto.OAuthTokenResponeDto{}

	if err := json.Unmarshal(bodyToken, &otoken); err != nil {
		o.logger.Debug("json.Unmarshal", "err", err)
		return dto.TokenPair{}, err
	}

	requstUserInfo, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		o.logger.Debug("http.NewRequest", "err", err)
		return dto.TokenPair{}, err
	}
	requstUserInfo.Header.Set("Authorization", fmt.Sprintf("Bearer %s", otoken.AccessToken))

	client := http.Client{}
	respone, err := client.Do(requstUserInfo)
	if err != nil {
		o.logger.Debug("client.Do", "err", err)
		return dto.TokenPair{}, err
	}

	if respone.StatusCode != http.StatusOK {
		o.logger.Debug("status google user-info != 200")
		return dto.TokenPair{}, errorz.ErrAuthorization
	}

	bodyUserInfo, err := io.ReadAll(respone.Body)
	if err != nil {
		o.logger.Debug("io.ReadAll", "err", err)
		return dto.TokenPair{}, err
	}

	ouser := dto.OGoogleUserDto{}
	if err := json.Unmarshal(bodyUserInfo, &ouser); err != nil {
		o.logger.Debug("json.Unmarshal", "err", err)
		return dto.TokenPair{}, err
	}

	user, err := o.userAdapter.GetUser(ctx, ouser.Username)
	if err != nil {
		if errors.Is(err, errorz.ErrUserNotFound) {
			userId, err := o.userAdapter.CreateUser(ctx, database.User{
				Username: ouser.Username,
				Email:    ouser.Email,
				Password: []byte("None"),
				Createat: time.Now(),
			})
			if err != nil {
				o.logger.Debug("o.userAdapter.CreateUser", "err", err)
				return dto.TokenPair{}, err
			}

			role, err := o.roleAdapter.GetRole(ctx, domain.UserRole)
			if err != nil {
				o.logger.Debug("o.roleAdapter.GetRole", "err", err)
				return dto.TokenPair{}, err
			}

			if err := o.userAdapter.CreateUserRole(ctx, database.UserRole{
				UserID: userId,
				RoleID: role.ID,
			}); err != nil {
				o.logger.Debug("o.userAdapter.CreateUserRole", "err", err)
				return dto.TokenPair{}, err
			}

			sessionToken, err := o.tokenManager.GenerateSessionToken(userId)
			if err != nil {
				o.logger.Debug("o.tokenManager.GenerateSessionToken", "err", err)
				return dto.TokenPair{}, err
			}
			if err := o.sessionAdapter.CreateSession(ctx, domain.Session{
				UserId:       userId,
				SessionToken: sessionToken,
				TTL:          o.tokenConfig.AccessExp,
			}); err != nil {
				o.logger.Debug("o.sessionAdapter.CreateSession", "err", err)
				return dto.TokenPair{}, err
			}

			fmt.Println("create user")

			refreshToken, err := o.tokenManager.GenerateRefreshToken(sessionToken)
			if err != nil {
				o.logger.Debug("o.tokenManager.GenerateRefreshToken", "err", err)
				return dto.TokenPair{}, err
			}

			accessToken, err := o.tokenManager.GenerateAccessToken(userId, role.ID)
			if err != nil {
				o.logger.Debug("o.tokenManager.GenerateAccessToken", "err", err)
				return dto.TokenPair{}, err
			}
			return dto.TokenPair{
				RefreshToken: refreshToken,
				AccessToken:  accessToken,
			}, nil
		}
		return dto.TokenPair{}, err
	}

	roleId, err := o.userAdapter.GetUserRole(ctx, user.ID)
	if err != nil {
		o.logger.Debug("o.userAdapter.GetUserRole", "err", err)
		return dto.TokenPair{}, err
	}

	sessionToken, err := o.tokenManager.GenerateSessionToken(user.ID)
	if err != nil {
		o.logger.Debug("o.tokenManager.GenerateSessionToken")
		return dto.TokenPair{}, err
	}

	if err := o.sessionAdapter.CreateSession(ctx, domain.Session{
		UserId:       user.ID,
		SessionToken: sessionToken,
		TTL:          o.tokenConfig.SessionExp,
	}); err != nil {
		o.logger.Debug("o.sessionAdapter.CreateSession", "err", err)
		return dto.TokenPair{}, err
	}

	refreshToken, err := o.tokenManager.GenerateRefreshToken(sessionToken)
	if err != nil {
		o.logger.Debug("o.tokenManager.GenerateSessionToken", "err", err)
		return dto.TokenPair{}, err
	}

	accessToken, err := o.tokenManager.GenerateAccessToken(user.ID, roleId.RoleID)
	if err != nil {
		o.logger.Debug("o.tokenManager.GenerateAccessToken", "err", err)
		return dto.TokenPair{}, err
	}

	return dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
