package exceptions

import (
	jsonresponse "after-sales/api/helper/json/json-response"
	"after-sales/api/utils"
	"errors"
	"net/http"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/sirupsen/logrus"
)

// BaseErrorResponse defines the general error response structure
type BaseErrorResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Err        error       `json:"-"` // The underlying error is not included in the response
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

	// If no message is set, try to assign a default message
	if err.Message == "" {
		if err.Err != nil && err.Err.Error() != "" {
			err.Message = translateErrorMessage(err.Err)
		} else {
			err.Message = defaultMessage
		}
	}

	// Log the error if there is an underlying error
	if err.Err != nil {
		logrus.Error(err)
	}

	// Build the response object
	res := &BaseErrorResponse{
		StatusCode: statusCode,
		Message:    err.Message,
	}

	// Write the response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	jsonresponse.WriteToResponseBody(writer, res)
}

// translateErrorMessage provides a human-readable error message for specific SQL errors
func translateErrorMessage(err error) string {
	var sqlErr mssql.Error
	if errors.As(err, &sqlErr) {
		switch sqlErr.Number {
		case 547:
			return "This record is associated with other data and cannot be deleted or modified. Please ensure there are no dependent records before proceeding."
		case 2601, 2627:
			return "Data already exists (Unique Constraint Violation)."
		case 1205:
			return "Deadlock detected, transaction stopped."
		case 208:
			return "Invalid object name. Table, view, or another object is not found."
		case 4060:
			return "Failed to access the database."
		case 53:
			return "Cannot connect to the server. Hostname or port is wrong."
		case 233:
			return "Database server refused the connection because resources are full."
		case 515, 8153:
			return "Cannot insert or update data to empty value."
		case 8152:
			return "Data exceeds column length or data type."
		case 102:
			return "SQL Syntax is not valid."
		case 207:
			return "Column name not found."
		case 209:
			return "Ambiguous column name."
		case 4104:
			return "Multi-part identifier could not be bound in SQL syntax."
		case 701:
			return "Server is running out of memory."
		case 8645:
			return "Resource request failed due to out of memory."
		case 1105:
			return "There is not enough disk space."
		default:
			return err.Error()
		}
	} else {
		return err.Error()
	}
}
