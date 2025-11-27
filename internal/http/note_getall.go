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
