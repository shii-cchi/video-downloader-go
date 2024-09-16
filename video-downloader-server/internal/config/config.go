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
	DbName       string
	DbUser       string
	DbPassword   string
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

	dbName := os.Getenv("DB_NAME")

	if dbName == "" {
		return nil, errors.New("DB_NAME " + errParamNotDefined)
	}

	dbUser := os.Getenv("DB_USER")

	if dbUser == "" {
		return nil, errors.New("DB_USER " + errParamNotDefined)
	}

	dbPassword := os.Getenv("DB_PASSWORD")

	if dbPassword == "" {
		return nil, errors.New("DB_PASSWORD " + errParamNotDefined)
	}

	return &Config{
		Port:         port,
		ExtensionURL: extensionURL,
		DbName:       dbName,
		DbUser:       dbUser,
		DbPassword:   dbPassword,
	}, nil
}
