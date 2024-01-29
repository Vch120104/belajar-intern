package exceptions

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func EntityException(c *gin.Context, message string) {
	res := Error{
		Success: false,
		Message: message,
		Data: nil,
	}

	c.JSON(http.StatusUnprocessableEntity, res)
}

type EntityError struct {
	Error string
}

func NewEntityError(error string) EntityError {
	return EntityError{Error: error}
}
