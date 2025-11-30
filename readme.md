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
go get gorm.io/driver/mysql
go get github.com/google/uuid@v1.6.0
go get github.com/aws/aws-lambda-go@v1.48.0
go get github.com/awslabs/aws-lambda-go-api-proxy@v0.16.2
go get github.com/renatofagalde/module-bitly@lates
go get github.com/google/uuid
go get github.com/testcontainers/testcontainers-go@v0.32.0
go get github.com/testcontainers/testcontainers-go/wait@v0.32.0
```

## 3# Estrutura padrão
```jshelllanguage
mkdir -p cmd/api
mkdir -p cmd/lambda

mkdir -p internal/config
mkdir -p internal/db
mkdir -p internal/http
mkdir -p internal/notes
mkdir -p internal/shared/errors

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
``internal/db/database.go``

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
``internal/db/postgres.go``
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
``internal/db/mysql.go``
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

## 8# Domínio

``internal/notes/model.go``

````go
package notes

import (
	"time"

	"gorm.io/datatypes"
)

type Note struct {
	ID        string         `gorm:"type:text;primaryKey"       json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Content   datatypes.JSON `gorm:"type:jsonb;not null"        json:"content"`
	CreatedAt time.Time      `gorm:"autoCreateTime"             json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"             json:"updated_at"`
	DeletedAt *time.Time     `gorm:"index"                      json:"deleted_at,omitempty"`
}

type CreateNoteRequest struct {
	Name    string         `json:"name"`
	Content datatypes.JSON `json:"content"`
}

type NoteResponse struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Content   datatypes.JSON `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time     `json:"deleted_at,omitempty"`
}

func (n *Note) ToResponse() *NoteResponse {
	return &NoteResponse{
		ID:        n.ID,
		Name:      n.Name,
		Content:   n.Content,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		DeletedAt: n.DeletedAt,
	}
}

````

## 9# Repository

```markdown
internal/notes/
  ├── repository.go        ← interface + struct base
  ├── repository/repository_create.go ← Create()
  ├── repository/repository_get.go    ← GetByID()
  ├── repository/repository_list.go   ← GetAll()
  └── model.go              ← struct Note
```

````go
package notes

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoteNotFound = errors.New("note not found")
)

type Repository interface {
	Create(ctx context.Context, n *Note) error
	GetByID(ctx context.Context, id string) (*Note, error)
	GetAll(ctx context.Context) ([]*Note, error)
}

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

````

``create``
```go
package notes

import "context"

func (r *gormRepository) Create(ctx context.Context, n *Note) error {
	return r.db.WithContext(ctx).Create(n).Error
}
```

``get``
```go
package notes

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

func (r *gormRepository) GetByID(ctx context.Context, id string) (*Note, error) {
	var note Note

	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&note).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNoteNotFound
	}
	if err != nil {
		return nil, err
	}

	return &note, nil
}
```

``get_all``
```go
package notes

import "context"

func (r *gormRepository) GetAll(ctx context.Context) ([]*Note, error) {
	var notes []Note

	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&notes).Error

	if err != nil {
		return nil, err
	}

	res := make([]*Note, 0, len(notes))
	for i := range notes {
		res = append(res, &notes[i])
	}

	return res, nil
}
```

## 10# Usecase e Service

```markdown
internal/
 └── notes/
      ├── usecase.go            ← interface, somente contratos
      ├── service/              ← implementação concreta
      ├── service/service.go    ← interface do repo
```

``internal/notes/repository/create.go``

````go
package usecase

import (
	"bootstrap/internal/notes/models"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	bitly "github.com/renatofagalde/module-bitly"
)

func (usecase *notesUsecase) CreateNote(ctx context.Context, note *models.CreateNoteRequest) (*models.NoteResponse, error) {

	var name string = strings.TrimSpace(note.Name)
	if len(name) < 1 || len(note.Content) < 1 {
		return nil, errInvalidInput
	}

	var n *models.Note = &models.Note{
		ID:        bitly.EncodeBytes([]byte(uuid.NewString())),
		Name:      name,
		Content:   note.Content,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: nil,
	}

	if err := usecase.repository.Create(ctx, n); err != nil {
		return nil, err
	}

	return n.ToResponse(), nil
}

````


``internal/notes/usecase/get.go``
```go
package usecase

import (
	"bootstrap/internal/notes/models"
	"context"
)

func (usecase *notesUsecase) GetNote(ctx context.Context, id string) (*models.NoteResponse, error) {

	if len(id) < 1 {
		return nil, errInvalidInput
	}

	n, err := usecase.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return n.ToResponse(), nil
}

```

``internal/notes/usecase/get_all.go``
```go
package usecase

import (
	"bootstrap/internal/notes/models"
	"context"
)

func (usecase *notesUsecase) ListNotes(ctx context.Context) ([]*models.NoteResponse, error) {
	notes, err := usecase.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*models.NoteResponse, 0, len(notes))
	for _, n := range notes {
		res = append(res, n.ToResponse())
	}
	return res, nil
}
```


## 11 HTTP com Gin

``internal\http\router.go``

```go
package http

import (
	"bootstrap/internal/notes/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotesService struct {
	usecase.UseCase
}

type RouterConfig struct {
	NotesService NotesService
}

func NewRouter(cfg RouterConfig) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	notesHandler := NewNotesHandler(cfg.NotesService)

	notesGroup := r.Group("/notes")
	{
		notesGroup.POST("", notesHandler.Create)
		notesGroup.GET("", notesHandler.GetAll)
		notesGroup.GET("/:id", notesHandler.GetByID)
	}

	return r
}

```

``internal/http/notes_handler.go``

```go
package http

import (
	"bootstrap/internal/notes/usecase"

	"github.com/gin-gonic/gin"
)

type NotesHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	GetAll(c *gin.Context)
}

type notesHandler struct {
	service usecase.UseCase
}

func NewNotesHandler(service usecase.UseCase) NotesHandler {
	return &notesHandler{service: service}

}

```


``internal/http/note_create.go``

````go
package http

import (
	"bootstrap/internal/notes/models"
	"bootstrap/internal/notes/usecase"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (n notesHandler) Create(c *gin.Context) {

	var request models.CreateNoteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	response, err := n.service.CreateNote(c.Request.Context(), &request)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, response)
}

````

``internal/http/note_by_id.go``
````go
package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *notesHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	res, err := h.service.GetNote(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, res)
}

````

``internal/http/note_getall.go``

````go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *notesHandler) GetAll(c *gin.Context) {
	res, err := h.service.ListNotes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, res)
}

````

## 12 main.go - normal app


``cmd/api/main.go``

```go
package main

import (
	"bootstrap/internal/notes/models"
	"bootstrap/internal/notes/repository"
	"bootstrap/internal/notes/usecase"
	"log"

	"bootstrap/internal/config"
	appdb "bootstrap/internal/db"
	apphttp "bootstrap/internal/http"
)

func main() {
	cfgLoader := config.NewEnvLoader()
	appCfg, err := cfgLoader.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dbCfg := appdb.Config{
		Host:     appCfg.PostgresHost,
		Port:     appCfg.PostgresPort,
		User:     appCfg.PostgresUser,
		Password: appCfg.PostgresPassword,
		DBName:   appCfg.PostgresDB,
		SSLMode:  appCfg.PostgresSSL,
	}

	dbInstance, err := appdb.NewPostgres(dbCfg)
	if err != nil {
		log.Fatalf("failed to init postgres: %v", err)
	}

	gormDB := dbInstance.Gorm()

	if err := gormDB.AutoMigrate(&models.Note{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	repo := repository.NewNoteRepository(gormDB)
	service := usecase.NewService(repo)

	router := apphttp.NewRouter(apphttp.RouterConfig{NotesService: service})

	log.Printf("API listening on :%s", appCfg.AppPort)
	if err := router.Run(":" + appCfg.AppPort); err != nil {
		log.Fatal(err)
	}
}

```


## 13 main.go - lambda

```go
package main

import (
	"bootstrap/internal/notes/repository"
	"bootstrap/internal/notes/usecase"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"bootstrap/internal/config"
	appdb "bootstrap/internal/db"
	apphttp "bootstrap/internal/http"
)

var ginLambda *ginadapter.GinLambda

func init() {
	cfgLoader := config.NewEnvLoader()
	appCfg, err := cfgLoader.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dbCfg := appdb.Config{
		Host:     appCfg.PostgresHost,
		Port:     appCfg.PostgresPort,
		User:     appCfg.PostgresUser,
		Password: appCfg.PostgresPassword,
		DBName:   appCfg.PostgresDB,
		SSLMode:  appCfg.PostgresSSL,
	}

	dbInstance, err := appdb.NewPostgres(dbCfg)
	if err != nil {
		log.Fatalf("failed to init postgres: %v", err)
	}

	gormDB := dbInstance.Gorm()

	repo := repository.NewNoteRepository(gormDB)
	service := usecase.NewService(repo)

	router := apphttp.NewRouter(apphttp.RouterConfig{
		NotesService: service,
	})

	ginLambda = ginadapter.New(router)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

```

## 14 - Test container
``tests/e2e_setup_test.go``

```go
package tests

import (
	"bootstrap/internal/notes/repository"
	"bootstrap/internal/notes/usecase"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	appdb "bootstrap/internal/db"
	apphttp "bootstrap/internal/http"

	"github.com/gin-gonic/gin"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

var (
	testDB      *gorm.DB
	testRouter  *gin.Engine
	pgContainer tc.Container
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. Sobe container Postgres via Testcontainers
	req := tc.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "notes_test",
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").
			WithStartupTimeout(60 * time.Second),
	}

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatalf("falha ao subir container Postgres de teste: %v", err)
	}
	pgContainer = container
	defer func() {
		_ = pgContainer.Terminate(ctx)
	}()

	// 2. Descobre host/port mapeados
	host, err := pgContainer.Host(ctx)
	if err != nil {
		log.Fatalf("erro ao obter host do container: %v", err)
	}

	mappedPort, err := pgContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatalf("erro ao obter porta mapeada do container: %v", err)
	}

	port, err := strconv.Atoi(mappedPort.Port())
	if err != nil {
		log.Fatalf("porta inválida: %v", err)
	}

	fmt.Printf("🧪 Postgres de teste em %s:%d\n", host, port)

	// 3. Conecta usando sua camada de DB (internal/db)
	dbCfg := appdb.Config{
		Host:     host,
		Port:     port,
		User:     "test",
		Password: "test",
		DBName:   "notes_test",
		SSLMode:  "disable",
	}

	dbInstance, err := appdb.NewPostgres(dbCfg)
	if err != nil {
		log.Fatalf("falha ao conectar no Postgres de teste: %v", err)
	}

	testDB = dbInstance.Gorm()

	// 4. Executa o SQL de criação + inserts usando sql/notes.sql
	if err := applySQL(testDB); err != nil {
		log.Fatalf("falha ao aplicar SQL de teste: %v", err)
	}

	// 5. Monta Repo + Service + Router (igual app real)
	repo := repository.NewNoteRepository(testDB)
	svc := usecase.NewService(repo)

	testRouter = apphttp.NewRouter(apphttp.RouterConfig{
		NotesService: svc,
	})

	// 6. Executa os testes
	code := m.Run()
	os.Exit(code)
}

// aplica o conteúdo de sql/notes.sql no banco de teste
func applySQL(db *gorm.DB) error {
	// caminho relativo a partir da pasta tests/
	sqlPath := filepath.Join("..", "sql", "notes.sql")

	fmt.Printf("📄 Aplicando SQL de teste: %s\n", sqlPath)

	bytes, err := os.ReadFile(sqlPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo SQL: %w", err)
	}

	if err := db.Exec(string(bytes)).Error; err != nil {
		return fmt.Errorf("erro ao executar SQL: %w", err)
	}

	return nil
}

```

``tests/notes_create_e2e_test.go``
````go
package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_CreateNote_Success(t *testing.T) {
	body := `{
		"name": "Teste E2E - Sucesso",
		"content": { "html": "<h1>Hello Test</h1>", "lang": "pt-BR" }
	}`

	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("esperado status 201, recebeu %d. Body: %s", w.Code, w.Body.String())
	}

	var res noteResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("falha ao parsear resposta JSON: %v", err)
	}

	if res.ID == "" {
		t.Fatalf("esperava ID preenchido")
	}
	if res.Name != "Teste E2E - Sucesso" {
		t.Fatalf("nome diferente do esperado. got=%q", res.Name)
	}
	if len(res.Content) == 0 {
		t.Fatalf("content não pode ser vazio")
	}
}

func Test_CreateNote_InvalidJSON(t *testing.T) {
	body := `{
		"name": "Teste JSON inválido",
		"content": { "html": "<h1>Oops</h1>" }` // JSON quebrado

	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperado status 400 para JSON inválido, recebeu %d. Body: %s", w.Code, w.Body.String())
	}
}

func Test_CreateNote_MissingFields(t *testing.T) {
	// name em branco e content vazio → deve bater no ErrInvalidInput
	body := `{
		"name": "   ",
		"content": {}
	}`

	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperado status 400 para input inválido, recebeu %d. Body: %s", w.Code, w.Body.String())
	}
}

````

``tests/notes_get_by_id_e2e_test.go``
```go
package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetNoteByID_Existing(t *testing.T) {
	// ID semeado em sql/notes.sql
	const seededID = "11111111-1111-1111-1111-111111111111"

	req := httptest.NewRequest(http.MethodGet, "/notes/"+seededID, nil)
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperado status 200, recebeu %d. Body: %s", w.Code, w.Body.String())
	}

	var res noteResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("falha ao parsear resposta: %v", err)
	}
	if res.ID != seededID {
		t.Fatalf("ID diferente do esperado. got=%q want=%q", res.ID, seededID)
	}
	if res.Name == "" {
		t.Fatalf("nome não deveria estar vazio")
	}
}

func Test_GetNoteByID_NotFound(t *testing.T) {
	const nonExistingID = "99999999-9999-9999-9999-999999999999"

	req := httptest.NewRequest(http.MethodGet, "/notes/"+nonExistingID, nil)
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("esperado status 404, recebeu %d. Body: %s", w.Code, w.Body.String())
	}
}

```

``tests/notes_get_all_e2e_test.go``
````go
package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_GetAllNotes_ReturnsSeededData(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperado status 200, recebeu %d. Body: %s", w.Code, w.Body.String())
	}

	var res []noteResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("falha ao parsear lista: %v", err)
	}

	if len(res) == 0 {
		t.Fatalf("esperava pelo menos 1 note dos seeds")
	}
}

func Test_GetAllNotes_AfterCreate_IncludesNewNote(t *testing.T) {
	// 1) Cria uma note nova
	body := `{
		"name": "Note criada no teste GetAll",
		"content": { "section": "test", "value": 123 }
	}`

	reqCreate := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()

	testRouter.ServeHTTP(wCreate, reqCreate)

	if wCreate.Code != http.StatusCreated {
		t.Fatalf("falha ao criar note. status=%d body=%s", wCreate.Code, wCreate.Body.String())
	}

	var created noteResponse
	if err := json.Unmarshal(wCreate.Body.Bytes(), &created); err != nil {
		t.Fatalf("erro ao parsear resposta de criação: %v", err)
	}

	// 2) Chama GET /notes e verifica se novo ID está na lista
	reqList := httptest.NewRequest(http.MethodGet, "/notes", nil)
	wList := httptest.NewRecorder()

	testRouter.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Fatalf("esperado status 200, recebeu %d. Body: %s", wList.Code, wList.Body.String())
	}

	var list []noteResponse
	if err := json.Unmarshal(wList.Body.Bytes(), &list); err != nil {
		t.Fatalf("erro ao parsear lista: %v", err)
	}

	if len(list) == 0 {
		t.Fatalf("lista não deveria estar vazia")
	}

	found := false
	for _, n := range list {
		if n.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("note criada (id=%s) não encontrada na lista", created.ID)
	}
}

````