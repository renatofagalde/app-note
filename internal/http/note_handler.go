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
