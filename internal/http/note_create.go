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
