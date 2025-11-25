## 1# Criar o projeto e inicializar o módulo

```shell
mkdir app-notes
cd app-notes
go mod init bootstrap

```


## 2# Adicionar dependências

```shell
go get github.com/gin-gonic/gin@v1.10.0
go get gorm.io/gorm@v1.25.9
go get gorm.io/driver/postgres@v1.5.7
go get github.com/google/uuid@v1.6.0
go get github.com/aws/aws-lambda-go@v1.48.0
go get github.com/awslabs/aws-lambda-go-api-proxy@v0.16.2
```

## 3# Estrutura padrão
```jshelllanguage
mkdir -p cmd/api
mkdir -p cmd/lambda

mkdir -p internal/config
mkdir -p internal/db
mkdir -p internal/http
mkdir -p internal/notes

mkdir -p sql
mkdir -p tests/testdata
```

desenho da estrutura:
```jshelllanguage
app-notes/
├── cmd
│   ├── api
│   │   └── main.go
│   └── lambda
│       └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── db
│   │   └── postgres.go
│   ├── http
│   │   ├── router.go
│   │   └── notes_handler.go
│   └── notes
│       ├── model.go
│       ├── repository.go
│       └── service.go
├── sql
│   └── notes.sql
├── tests
│   ├── docker-compose.test.yml
│   ├── e2e_setup.go
│   ├── notes_e2e_test.go
│   └── testdata
│       └── notes.sql
└── docker-compose.yml
```

## 4# SQL de criação de tabela + registros


criar os arquivos com o mesmo conteúdo:
```jshelllanguage
cp sql/notes.sql tests/testdata/notes.sql
```

```jql
-- DDL: tabela notes
CREATE TABLE IF NOT EXISTS notes (
    id          TEXT         PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    content     JSONB        NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ  NULL
);

CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes (created_at);
CREATE INDEX IF NOT EXISTS idx_notes_name ON notes (LOWER(name));

-- INSERT de exemplo 1
INSERT INTO notes (
    id, name, content, created_at, updated_at, deleted_at
) VALUES (
    '11111111-1111-1111-1111-111111111111',
    'Home - Meu Site',
    $${
      "url": "https://meusite.com/",
      "contentType": "text/html",
      "html": "<!DOCTYPE html><html><head><title>Meu Site</title></head><body><h1>Bem-vindo ao Meu Site</h1><p>Conteúdo qualquer...</p></body></html>"
    }$$::jsonb,
    NOW(),
    NOW(),
    NULL
);

-- INSERT de exemplo 2
INSERT INTO notes (
    id, name, content, created_at, updated_at, deleted_at
) VALUES (
    '22222222-2222-2222-2222-222222222222',
    'Artigo - Observabilidade em Go',
    $${
      "url": "https://blog.meusite.com/observabilidade-em-go",
      "contentType": "text/html",
      "html": "<html><head><title>Observabilidade em Go</title></head><body><h1>Observabilidade em Go</h1><p>Logs, métricas e traces são pilares importantes...</p><pre><code>func main() { /* ... */ }</code></pre></body></html>",
      "metadata": {
        "author": "Renato",
        "tags": ["go", "observabilidade", "logs"],
        "language": "pt-BR"
      }
    }$$::jsonb,
    NOW(),
    NOW(),
    NULL
);

-- INSERT de exemplo 3
INSERT INTO notes (
    id, name, content, created_at, updated_at, deleted_at
) VALUES (
    '33333333-3333-3333-3333-333333333333',
    'Dump HTML - Página de Login',
    $${
      "url": "https://app.meusite.com/login",
      "contentType": "text/html",
      "blob": "<html><head><title>Login</title></head><body><form><input type='email' name='email'/><input type='password' name='password'/></form></body></html>"
    }$$::jsonb,
    NOW(),
    NOW(),
    NULL
);
```

## 5# Docker compose
5# Docker Compose para desenvolvimento (docker-compose.yaml)
```dockerfile
version: "3.9"

services:
  notes-db:
    image: postgres:latest
    container_name: notes-db
    restart: unless-stopped
    environment:
      POSTGRES_DB: notes-db
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    ports:
      - "5433:5432"   # host:container
    # scripts DDL + seeds
    volumes:
      - ./sql:/docker-entrypoint-initdb.d
    logging:
      options:
        max-size: "10m"
        max-file: "3"

```

rodando o banco de dados
```jshelllanguage
docker compose up -d
```

## 6# Config/LoadEnv

internal/config/config.go
```go
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
```

## 7# GORM 

arquivo de interface

```go
package db

import "gorm.io/gorm"

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Database interface {
	Gorm() *gorm.DB
}

```

postgres:
```go
package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type gormDatabase struct {
	db *gorm.DB
}

func (d *gormDatabase) Gorm() *gorm.DB {
	return d.db
}

func NewPostgres(cfg Config) (Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return &gormDatabase{db: db}, nil
}
```

mysql
````go
package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Reaproveita o mesmo gormDatabase

func NewMySQL(cfg Config) (Database, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return &gormDatabase{db: db}, nil
}

````