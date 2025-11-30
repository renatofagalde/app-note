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
