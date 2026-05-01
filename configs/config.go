package configs

import (
	"errors"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

const (
	envBasicAuthUsername string = "BASIC_AUTH_USERNAME"
	envBasicAuthPassword string = "BASIC_AUTH_PASSWORD"
	envServerAddr        string = "SERVER_ADDR"
)

var (
	ErrEmptyBasicAuthUsername = errors.New("BASIC_AUTH_USERNAME is empty")
	ErrEmptyBasicAuthPassword = errors.New("BASIC_AUTH_PASSWORD is empty")
	ErrEmptyServerAddr        = errors.New("SERVER_ADDR is empty")
)

type Config struct {
	ServerAddr string
	HTTPBasic  HTTPBasicConfig
}

type HTTPBasicConfig struct {
	Username string
	Password string
}

func Load() (*Config, error) {
	basicAuthUsername := os.Getenv(envBasicAuthUsername)
	if basicAuthUsername == "" {
		return nil, ErrEmptyBasicAuthUsername
	}

	basicAuthPassword := os.Getenv(envBasicAuthPassword)
	if basicAuthPassword == "" {
		return nil, ErrEmptyBasicAuthPassword
	}

	serverAddress := os.Getenv(envServerAddr)
	if serverAddress == "" {
		return nil, ErrEmptyServerAddr
	}

	cfg := Config{
		ServerAddr: serverAddress,
		HTTPBasic: HTTPBasicConfig{
			Username: basicAuthUsername,
			Password: basicAuthPassword,
		},
	}

	return &cfg, nil
}
