package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string     `yaml:"log_level" env:"LOG_LEVEL"` // trace, debug, info, warn, error, fatal, panic, disabled
	Auth     AuthConfig `yaml:"auth"`
	HTTP     HTTPConfig `yaml:"http"`
	Log      LogConfig  `yaml:"log"`
}

type AuthConfig struct {
	Secret   string `yaml:"secret" env:"AUTH_SECRET"`
	Password string `yaml:"password" env:"PASSWORD"`
}

type HTTPConfig struct {
	Port    int    `yaml:"port" env:"HTTP_PORT"`
	GinMode string `yaml:"gin_mode" env:"GIN_MODE"`
}

type LogConfig struct {
	Dir       string `yaml:"dir" env:"LOG_DIR"`
	Retention int64  `yaml:"retention" env:"LOG_RETENTION"`
	ChunkSize int64  `yaml:"chunk_size" env:"LOG_CHUNK_SIZE"`
}

func LoadConfig(configFilePath string) (Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(configFilePath, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("error load config: %w", err)
	}
	return cfg, nil
}
