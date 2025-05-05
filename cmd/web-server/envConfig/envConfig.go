package appconfig

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	scheme string
	host   string
}

func NewEnvConfig() *EnvConfig {
	godotenv.Load()

	scheme := os.Getenv("R_SCHEME")
	if scheme == "" {
		scheme = "http"
	}
	host := os.Getenv("R_HOST")
	if host == "" {
		host = "localhost"
	}
	return &EnvConfig{scheme, host}
}

func (c *EnvConfig) GetScheme() string {
	return c.scheme
}

func (c *EnvConfig) GetHost() string {
	return c.host
}
