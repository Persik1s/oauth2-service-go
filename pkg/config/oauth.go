package config

import "github.com/ilyakaznacheev/cleanenv"

type OAuthGoogleConfig struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
}

func LoadOAuthGoogleConfig() (*OAuthGoogleConfig, error) {
	var cfg struct {
		Google OAuthGoogleConfig `json:"oauth_google"`
	}

	if err := cleanenv.ReadConfig("config.json", &cfg); err != nil {
		return nil, err
	}

	return &cfg.Google, nil
}
