package config

import "github.com/ilyakaznacheev/cleanenv"

type ApplicationConfig struct {
	ClientId int `json:"client_id"`
}

func LoadApplicationConfig() (*ApplicationConfig, error) {
	var cfg struct {
		Application ApplicationConfig `json:"application"`
	}

	if err := cleanenv.ReadConfig("config.json", &cfg); err != nil {
		return nil, err
	}
	return &cfg.Application, nil
}
