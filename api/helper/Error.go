package helper

import (
	"after-sales/api/exceptions"
	"errors"
	"net/http"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func ReturnError(writer http.ResponseWriter, request *http.Request, err *exceptions.BaseErrorResponse) {
	switch err.StatusCode {
	case http.StatusUnauthorized:
		exceptions.NewAuthorizationException(writer, request,err)
	case http.StatusBadRequest:
		exceptions.NewBadRequestException(writer, request, err)
	case http.StatusUnprocessableEntity:
		exceptions.NewEntityException(writer, request, err)
	case http.StatusNotFound:
		exceptions.NewNotFoundException(writer, request, err)
	case http.StatusForbidden:
		exceptions.NewRoleException(writer, request, err)
	case http.StatusConflict:
		exceptions.NewConflictException(writer, request, err)
	default:
		exceptions.NewAppException(writer, request, err)
	}
}

// ConvertBaseErrorResponseToError converts a BaseErrorResponse to a standard error
func ConvertBaseErrorResponseToError(baseErr *exceptions.BaseErrorResponse) error {
	if baseErr == nil {
		return nil
	}
	return errors.New(baseErr.Message)
}
