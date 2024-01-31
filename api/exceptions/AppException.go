package exceptions

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type OldError struct {
	Success bool `json:"Success"`
	Message string `json:"Message"`
	Data interface{} `json:"Data"`
}

func AppException(c *gin.Context, message string) {
	res := OldError{
		Success: false,
		Message: message,
		Data: nil,
	}

	c.JSON(http.StatusInternalServerError, res)
}

type AppExceptionError struct {
	Error string
}

func NewAppExceptionError(error string) AppExceptionError {
	return AppExceptionError{Error: error}
}
