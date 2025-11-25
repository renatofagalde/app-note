package config

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
