package exceptions

import (
	"after-sales/api/helper"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type BaseErrorResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Err        string      `json:"error"`
}

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	baseError := &BaseErrorResponse{}
	switch e := err.(type) {
	case validator.ValidationErrors:
		baseError.StatusCode = http.StatusBadRequest
		baseError.Message = "Bad Request"
		baseError.Data = e.Error()
		baseError.Err = "Validation error"
	case BadRequestError:
		baseError.StatusCode = http.StatusBadRequest
		baseError.Message = "Bad Request"
		baseError.Data = e.Error
		baseError.Err = "Bad request error"
	case NoContentError:
		baseError.StatusCode = http.StatusNoContent
		baseError.Message = "Data Not Found"
		baseError.Data = e.Error
		baseError.Err = "No content error"
	case AuthorizationError:
		baseError.StatusCode = http.StatusUnauthorized
		baseError.Message = "You don't have permission"
		baseError.Data = e.Error
		baseError.Err = "Authorization error"
	case ConflictError:
		baseError.StatusCode = http.StatusConflict
		baseError.Message = "Data Already Exists"
		baseError.Data = e.Error
		baseError.Err = "Conflict error"
	case EntityError:
		baseError.StatusCode = http.StatusUnprocessableEntity
		baseError.Message = "Data Error, please check your input"
		baseError.Data = e.Error
		baseError.Err = "Entity error"
	case NotFoundError:
		baseError.StatusCode = http.StatusNotFound
		baseError.Message = "Data Not Found"
		baseError.Data = e.Error
		baseError.Err = "Not found error"
	case AppExceptionError:
		baseError.StatusCode = http.StatusInternalServerError
		baseError.Message = "Internal Server Error"
		baseError.Data = e.Error
		baseError.Err = "App exception error"
	case RoleError:
		baseError.StatusCode = http.StatusForbidden
		baseError.Message = "You don't have permission"
		baseError.Data = e.Error
		baseError.Err = "Role error"
	default:
		baseError.StatusCode = http.StatusInternalServerError
		baseError.Message = "Internal Server Error"
		baseError.Data = "An unexpected error occurred"
		baseError.Err = "Unexpected error"
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(baseError.StatusCode)
	helper.WriteToResponseBody(writer, baseError)
}
