package exceptions

import (
	"github.com/gin-gonic/gin"
	"net/http"
)
// Deprecated: please change to the latest one without *gin.Context
//
func AuthorizeException(c *gin.Context, message string) {
	res := OldError{
		Success: false,
		Message: message,
		Data: nil,
	}

	c.JSON(http.StatusUnauthorized, res)
}

type AuthorizationError struct {
	Error string
}

func NewAuthorizationError(error string) AuthorizationError {
	return AuthorizationError{Error: error}
}
