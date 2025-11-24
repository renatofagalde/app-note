package config

import (
	"os"
	"strconv"
)

type AppConfig struct {
	AppPort          string
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSSL      string
}

type Loader interface {
	Load() (*AppConfig, error)
}

type EnvLoader struct{}

func NewEnvLoader() Loader {
	return &EnvLoader{}
}

func (EnvLoader) Load() (*AppConfig, error) {
	return &AppConfig{
		AppPort:          getenv("APP_PORT", "8080"),
		PostgresHost:     getenv("POSTGRES_HOST", "localhost"),
		PostgresPort:     mustInt(getenv("POSTGRES_PORT", "5433")),
		PostgresUser:     getenv("POSTGRES_USER", "notes_user"),
		PostgresPassword: getenv("POSTGRES_PASSWORD", "notes_password"),
		PostgresDB:       getenv("POSTGRES_DB", "notes-db"),
		PostgresSSL:      getenv("POSTGRES_SSLMODE", "disable"),
	}, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func mustInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
