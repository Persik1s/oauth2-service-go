package token

import (
	"crypto/rsa"
	"github.com/Persik1s/oauth2-service-go/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type TokenManager struct {
	logger   *slog.Logger
	cfgToken *config.TokenConfig

	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewTokenManager(
	publicKey *rsa.PublicKey,
	privateKey *rsa.PrivateKey,
	cfgToken *config.TokenConfig,
	logger *slog.Logger) *TokenManager {
	return &TokenManager{
		logger:     logger,
		cfgToken:   cfgToken,
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

func (t *TokenManager) GenerateAccessToken(userId uuid.UUID, roleId uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"role_id": roleId.String(),
		"user_id": userId.String(),
		"exp":     time.Now().Add(t.cfgToken.AccessExp * time.Second).Unix(),
	})

	return token.SignedString(t.privateKey)
}

func (t *TokenManager) GenerateRefreshToken(sessionToken string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"session_id": sessionToken,
		"exp":        time.Now().Add(t.cfgToken.AccessExp * time.Second).Unix(),
	})

	return token.SignedString(t.privateKey)
}

func (t *TokenManager) GenerateSessionToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_id": userID.String(),
	})
	return token.SignedString(t.privateKey)
}
