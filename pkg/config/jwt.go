package config

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type TokenConfig struct {
	AccessExp  time.Duration `json:"access"`
	RefreshExp time.Duration `json:"refresh"`
	SessionExp time.Duration `json:"session"`
}

func LoadTokenConfig() (*TokenConfig, error) {
	var cfg struct {
		Token TokenConfig `json:"token"`
	}

	if err := cleanenv.ReadConfig("config.json", &cfg); err != nil {
		return nil, err
	}
	return &cfg.Token, nil
}

func LoadPublicKey() (*rsa.PublicKey, error) {
	body, err := os.ReadFile("public_key.pem")
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(body)
}

func LoadPrivateKey() (*rsa.PrivateKey, error) {
	body, err := os.ReadFile("private_key.pem")
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(body)
}
