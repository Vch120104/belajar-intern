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

// Error implements the error interface for BaseErrorResponse
func (e *BaseErrorResponse) Error() string {
	// If there is an underlying error, include it in the string
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	// Otherwise, just return the message
	return e.Message
}

// NewAppException creates a new AppException with a customizable HTTP status code
func NewAppException(writer http.ResponseWriter, request *http.Request, err *BaseErrorResponse) {
	handleError(writer, err, http.StatusInternalServerError, utils.SomethingWrong)
}

// NewAuthorizationException handles authorization errors
func NewAuthorizationException(writer http.ResponseWriter, request *http.Request, err *BaseErrorResponse) {
	handleError(writer, err, http.StatusUnauthorized, utils.SessionError)
}

// NewBadRequestException handles bad request errors
func NewBadRequestException(writer http.ResponseWriter, request *http.Request, err *BaseErrorResponse) {
	handleError(writer, err, http.StatusBadRequest, utils.BadRequestError)
}

// NewConflictException handles conflict errors
func NewConflictException(writer http.ResponseWriter, request *http.Request, err *BaseErrorResponse) {
	handleError(writer, err, http.StatusConflict, utils.DataExists)
}

// NewEntityException handles unprocessable entity errors
func NewEntityException(writer http.ResponseWriter, request *http.Request, err *BaseErrorResponse) {
	handleError(writer, err, http.StatusUnprocessableEntity, utils.JsonError)
}

// NewNotFoundException handles not found errors
func NewNotFoundException(writer http.ResponseWriter, request *http.Request, err *BaseErrorResponse) {
	handleError(writer, err, http.StatusNotFound, utils.GetDataNotFound)
}

// NewRoleException handles forbidden errors
func NewRoleException(writer http.ResponseWriter, request *http.Request, err *BaseErrorResponse) {
	handleError(writer, err, http.StatusForbidden, utils.PermissionError)
}

// handleError centralizes the error handling logic
func handleError(writer http.ResponseWriter, err *BaseErrorResponse, defaultStatusCode int, defaultMessage string) {
	statusCode := err.StatusCode
	if statusCode == 0 {
		statusCode = defaultStatusCode
	}

	if err.Message == "" {
		err.Message = defaultMessage
	}
	if err.Err != nil {
		logrus.Info(err)
	}

	res := &BaseErrorResponse{
		StatusCode: statusCode,
		Message:    err.Message,
	}

	writer.WriteHeader(statusCode)
	jsonresponse.WriteToResponseBody(writer, res)
}
