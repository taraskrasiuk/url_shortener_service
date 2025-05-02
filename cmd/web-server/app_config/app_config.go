package appconfig

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	scheme string
	host   string
}

func NewConfig() *AppConfig {
	godotenv.Load()

	scheme := os.Getenv("scheme")
	if scheme == "" {
		scheme = "http"
	}
	host := os.Getenv("host")
	if host == "" {
		host = "localhost"
	}
	return &AppConfig{scheme, host}
}

func (c *AppConfig) GetScheme() string {
	return c.scheme
}

func (c *AppConfig) GetHost() string {
	return c.host
}
