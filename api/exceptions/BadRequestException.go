package exceptions

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BadRequestException(c *gin.Context, message string) {
	res := OldError{
		Success: false,
		Message: message,
		Data: nil,
	}

	c.JSON(http.StatusBadRequest, res)
}

type BadRequestError struct {
	Error string
}

func NewBadRequestError(error string) BadRequestError {
	return BadRequestError{Error: error}
}
