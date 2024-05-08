package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"time"
)

type Config struct {
	HTTP `yaml:"http"`
	JWT  `yaml:"jwt"`
	PG
	Hasher
	Admin
}

type HTTP struct {
	Port string `yaml:"port"`
}

type JWT struct {
	SignKey  string        `env:"JWT_SIGN_KEY"`
	TokenTTL time.Duration `yaml:"token_ttl"`
}

type PG struct {
	Host     string `env:"POSTGRES_HOST"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `env:"POSTGRES_NAME"`
	Port     string `env:"POSTGRES_PORT"`
	SSLMode  string `env:"POSTGRES_SSL_MODE"`
}

type Hasher struct {
	Salt string `env:"HASH_SALT"`
}

type Admin struct {
	Username string `env:"ADMIN_USERNAME"`
	Password string `env:"ADMIN_PASSWORD"`
}

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}
	err = godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("error loading env variables: %w", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading env variables: %w", err)
	}

	return cfg, nil
}
