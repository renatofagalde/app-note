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
		Request: req,
		Started: true,
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
