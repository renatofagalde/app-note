package http

import (
	"net/http"

	domainerror "bootstrap/internal/shared/errors"

	"github.com/gin-gonic/gin"
)

func (h *notesHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	res, err := h.service.GetNote(c.Request.Context(), id)
	if err != nil {

		if dErr, ok := err.(*domainerror.DomainError); ok {
			switch dErr.Code {
			case "NOT_FOUND":
				c.JSON(http.StatusNotFound, gin.H{"error": dErr})
				return
			case "INVALID_INPUT":
				c.JSON(http.StatusBadRequest, gin.H{"error": dErr})
				return
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, res)
}
