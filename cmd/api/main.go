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
