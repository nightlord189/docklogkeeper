package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogDirectory string     `yaml:"log_directory" env:"LOG_DIRECTORY"`
	LogLevel     string     `yaml:"log_level" env:"LOG_LEVEL"` // trace, debug, info, warn, error, fatal, panic, disabled
	HTTP         HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
	Port    int    `yaml:"port" env:"HTTP_PORT"`
	GinMode string `yaml:"gin_mode" env:"GIN_MODE"`
}

func LoadConfig(configFilePath string) (Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(configFilePath, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("error load config: %w", err)
	}
	return cfg, nil
}
