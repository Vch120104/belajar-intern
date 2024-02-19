package exceptions

import (
	"github.com/gin-gonic/gin"
	"net/http"
)
// Deprecated: please change to the latest one without *gin.Context
//
func NotFoundException(c *gin.Context, message string) {
	res := OldError{
		Success: false,
		Message: message,
		Data: nil,
	}

	c.JSON(http.StatusNotFound, res)
}

type NotFoundError struct {
	Error string
}

func NewNotFoundError(error string) NotFoundError {
	return NotFoundError{Error: error}
}
