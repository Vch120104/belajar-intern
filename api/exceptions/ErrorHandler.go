package exceptions

import (
	jsonresponse "after-sales/api/helper_testt/json/json-response"
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

// NewAppException creates a new AppException with a customizable HTTP status code
func NewAppException(writer http.ResponseWriter, request *http.Request, err *BaseErrorResponse) {
	statusCode := http.StatusInternalServerError
	if err.Message == "" {
		err.Message = utils.SomethingWrong
	}
	if err.Err != nil {
		logrus.Info(err)
		res := &BaseErrorResponse{
			StatusCode: statusCode,
			Message:    err.Message,
			//Data:       err,
		}

		writer.WriteHeader(err.StatusCode)
		jsonresponse.WriteToResponseBody(writer, res)
		return
	}
}

func NewAuthorizationException(writer http.ResponseWriter, request *http.Request, err *BaseErrorResponse) {
	statusCode := http.StatusUnauthorized
	if err.Message == "" {
		err.Message = utils.SessionError
	}
	if err.Err != nil {
		logrus.Info(err)
		res := &BaseErrorResponse{
			StatusCode: statusCode,
			Message:    err.Message,
			//Data:       err,
		}

		writer.WriteHeader(statusCode)
		jsonresponse.WriteToResponseBody(writer, res)
		return
	}
}
func NewBadRequestException(writer http.ResponseWriter, request *http.Request, err *BaseErrorResponse) {
	statusCode := http.StatusBadRequest

	if err.Message == "" {
		err.Message = utils.BadRequestError
	}
	if err.Err != nil {
		logrus.Info(err)
		res := &BaseErrorResponse{
			StatusCode: statusCode,
			Message:    err.Err.Error(),
			//Data:       err,
		}

		writer.WriteHeader(statusCode)
		jsonresponse.WriteToResponseBody(writer, res)
		return
	}
}

func NewConflictException(writer http.ResponseWriter, request *http.Request, err *BaseErrorResponse) {
	statusCode := http.StatusConflict
	if err.Message == "" {
		err.Message = utils.DataExists
	}
	if err.Err != nil {
		logrus.Info(err)
		res := &BaseErrorResponse{
			StatusCode: statusCode,
			Message:    err.Message,
			//Data:       err,
		}

		writer.WriteHeader(statusCode)
		jsonresponse.WriteToResponseBody(writer, res)
		return
	}
}
func NewEntityException(writer http.ResponseWriter, request *http.Request, err *BaseErrorResponse) {
	statusCode := http.StatusUnprocessableEntity
	if err.Message == "" {
		err.Message = utils.JsonError
	}
	if err.Err != nil {
		logrus.Info(err)
		res := &BaseErrorResponse{
			StatusCode: statusCode,
			Message:    err.Err.Error(),
			//Data:       err,
		}

		writer.WriteHeader(statusCode)
		jsonresponse.WriteToResponseBody(writer, res)
		return
	}
}

func NewNotFoundException(writer http.ResponseWriter, request *http.Request, err *BaseErrorResponse) {
	statusCode := http.StatusNotFound
	if err.Message == "" {
		err.Message = utils.GetDataNotFound
	}
	if err.Err != nil {
		logrus.Info(err)
		res := &BaseErrorResponse{
			StatusCode: statusCode,
			Message:    err.Message,
			//Data:       err,
		}

		writer.WriteHeader(statusCode)
		jsonresponse.WriteToResponseBody(writer, res)
		return
	}
}
func NewRoleException(writer http.ResponseWriter, request *http.Request, err *BaseErrorResponse) {
	statusCode := http.StatusForbidden

	if err.Message == "" {
		err.Message = utils.PermissionError
	}
	if err.Err != nil {
		logrus.Info(err)
		res := &BaseErrorResponse{
			StatusCode: statusCode,
			Message:    err.Err.Error(),
			//Data:       err,
		}

		writer.WriteHeader(statusCode)
		jsonresponse.WriteToResponseBody(writer, res)
		return
	}
}
