package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

const (
	errParamNotDefined = "parameter is not defined"
)

type Config struct {
	Port         string
	ExtensionURL string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")

	if port == "" {
		return nil, errors.New("PORT " + errParamNotDefined)
	}

	extensionURL := os.Getenv("EXTENSION_URL")

	if extensionURL == "" {
		return nil, errors.New("EXTENSION_URL " + errParamNotDefined)
	}

	return &Config{
		Port:         port,
		ExtensionURL: extensionURL,
	}, nil
}
