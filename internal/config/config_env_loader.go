package config

import "os"
import "strconv"

type EnvLoader struct {
}

func NewEnvLoader() Loader {
	return &EnvLoader{}
}

func (e *EnvLoader) Load() (*AppConfig, error) {
	return &AppConfig{
		AppPort:          getenv("APP_PORT", "8080"),
		PostgresHost:     getenv("POSTGRES_HOST", "localhost"),
		PostgresPort:     mustInt(getenv("POSTGRES_PORT", "5433")),
		PostgresUser:     getenv("POSTGRES_USER", "user"),
		PostgresPassword: getenv("POSTGRES_PASSWORD", "password"),
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
