package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      App      `yaml:"app"`
		HTTP     HTTP     `yaml:"http"`
		Log      Log      `yaml:"logger"`
		Download Download `yaml:"download"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"SERVER_PORT"`
	}
	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}
	Download struct {
		Path         string `env-required:"true" yaml:"downloadPath" env:"DOWNLOAD_PATH"`
		WorkersCount int    `env-required:"true" yaml:"workersCount" env:"WORKERS_COUNT"`
	}
)

func New(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config - NewConfig - cleanenv.ReadConfig: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("config - NewConfig - cleanenv.UpdateEnv: %w", err)
	}

	return cfg, nil
}
