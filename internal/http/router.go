package http

import (
	"bootstrap/internal/notes/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	NotesService usecase.UseCase
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
