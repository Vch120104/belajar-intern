package exceptions

import (
	jsonresponse "after-sales/api/helper/json/json-response"
	"after-sales/api/utils"
	"net/http"

	"github.com/sirupsen/logrus"
)

// BaseErrorResponse defines the general error response structure
type BaseErrorResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Err        error       `json:"-"`
}

// CustomError defines a custom error structure
type CustomError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}

// NotFoundError defines a custom error for resource not found
type NotFoundError struct {
	Resource string // Resource yang tidak ditemukan
}

// Error returns the error message
func (e NotFoundError) Error() string {
	return e.Resource + " not found"
}

// NewBaseErrorResponse creates a new BaseErrorResponse from an error
func NewBaseErrorResponse(statusCode int, defaultMessage string, err error) *BaseErrorResponse {
	message := defaultMessage
	if err != nil {
		message = err.Error()
	}

	return &BaseErrorResponse{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

// WriteErrorResponse writes the error response to the HTTP response writer
func WriteErrorResponse(writer http.ResponseWriter, err interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	var statusCode int
	var message string
	var errorString string

	switch e := err.(type) {
	case *BaseErrorResponse:
		statusCode = e.StatusCode
		message = e.Message
		errorString = e.Err.Error()
	case CustomError:
		statusCode = e.StatusCode
		message = e.Message
		errorString = e.Error
	default:
		statusCode = http.StatusInternalServerError
		message = "Internal Server Error"
		errorString = ""
	}

	writer.WriteHeader(statusCode)
	response := CustomError{
		StatusCode: statusCode,
		Message:    message,
		Error:      errorString,
	}
	jsonresponse.WriteToResponseBody(writer, response)
	logrus.Info(err)
}

// NewAppException creates a new AppException with a customizable HTTP status code
func NewAppException(writer http.ResponseWriter, request *http.Request, err error) {
	statusCode := http.StatusInternalServerError
	message := utils.SomethingWrong

	if err != nil {
		logrus.Error(err)
		message = err.Error()
	}

	res := NewBaseErrorResponse(statusCode, message, err)
	WriteErrorResponse(writer, res)
}

// NewAuthorizationException creates a new AuthorizationException with a customizable HTTP status code
func NewAuthorizationException(writer http.ResponseWriter, request *http.Request, err error) {
	statusCode := http.StatusUnauthorized
	message := utils.SessionError

	if err != nil {
		logrus.Error(err)
		message = err.Error()
	}

	res := NewBaseErrorResponse(statusCode, message, err)
	WriteErrorResponse(writer, res)
}

// NewBadRequestException creates a new BadRequestException with a customizable HTTP status code
func NewBadRequestException(writer http.ResponseWriter, request *http.Request, err error) {
	statusCode := http.StatusBadRequest
	message := utils.BadRequestError

	if err != nil {
		logrus.Error(err)
		message = err.Error()
	}

	res := NewBaseErrorResponse(statusCode, message, err)
	WriteErrorResponse(writer, res)
}

// NewConflictException creates a new ConflictException with a customizable HTTP status code
func NewConflictException(writer http.ResponseWriter, request *http.Request, err error) {
	statusCode := http.StatusConflict
	message := utils.DataExists

	if err != nil {
		logrus.Error(err)
		message = err.Error()
	}

	res := NewBaseErrorResponse(statusCode, message, err)
	WriteErrorResponse(writer, res)
}

// NewEntityException creates a new EntityException with a customizable HTTP status code
func NewEntityException(writer http.ResponseWriter, request *http.Request, err error) {
	statusCode := http.StatusUnprocessableEntity
	message := utils.JsonError

	if err != nil {
		logrus.Error(err)
		message = err.Error()
	}

	res := NewBaseErrorResponse(statusCode, message, err)
	WriteErrorResponse(writer, res)
}

// NewNotFoundException creates a new NotFoundException with a customizable HTTP status code
func NewNotFoundException(writer http.ResponseWriter, request *http.Request, err error) {
	statusCode := http.StatusNotFound
	message := utils.GetDataNotFound

	if err != nil {
		logrus.Error(err)
		message = err.Error()
	}

	res := NewBaseErrorResponse(statusCode, message, err)
	WriteErrorResponse(writer, res)
}

// NewRoleException creates a new RoleException with a customizable HTTP status code
func NewRoleException(writer http.ResponseWriter, request *http.Request, err error) {
	statusCode := http.StatusForbidden
	message := utils.PermissionError

	if err != nil {
		logrus.Error(err)
		message = err.Error()
	}

	res := NewBaseErrorResponse(statusCode, message, err)
	WriteErrorResponse(writer, res)
}
