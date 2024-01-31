package exceptions

import (
	"after-sales/api/helper"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	StatusCode uint16      `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {

	if notFoundError(writer, request, err) {
		return
	}

	if notAuthorizedError(writer, request, err) {
		return
	}

	if conflictError(writer, request, err) {
		return
	}

	if entityError(writer, request, err) {
		return
	}

	if noContentError(writer, request, err) {
		return
	}

	if badRequestError(writer, request, err) {
		return
	}

	if validationErrors(writer, request, err) {
		return
	}

	if roleError(writer, request, err) {
		return
	}

	internalServerError(writer, request, err)
}
func validationErrors(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := Error{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad Request",
			Data:       exception.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func badRequestError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(BadRequestError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := Error{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad Request",
			Data:       exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}
func noContentError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(NoContentError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNoContent)

		webResponse := Error{
			StatusCode: http.StatusNoContent,
			Message:    "Data Not Found",
			Data:       exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func notAuthorizedError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(AuthorizationError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)

		webResponse := Error{
			StatusCode: http.StatusUnauthorized,
			Message:    "You don't have permission",
			Data:       exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func conflictError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(ConflictError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusConflict)

		webResponse := Error{
			StatusCode: http.StatusConflict,
			Message:    "Data Already Exists",
			Data:       exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}
func entityError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(EntityError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnprocessableEntity)

		webResponse := Error{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Data Error, please check your input",
			Data:       exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}
func notFoundError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		webResponse := Error{
			StatusCode: http.StatusNotFound,
			Message:    "Data Not Found",
			Data:       exception.Error,
		}
		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func internalServerError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(AppExceptionError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusInternalServerError)

		webResponse := Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			Data:       exception.Error,
		}
		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func roleError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(RoleError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusForbidden)

		webResponse := Error{
			StatusCode: http.StatusForbidden,
			Message:    "You don't have permission",
			Data:       exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

// func ReadFromRequestBody(request *http.Request, result interface{}) {
// 	decoder := json.NewDecoder(request.Body)
// 	err := decoder.Decode(result)
// 	if err != nil {
// 		panic(NewEntityError("Invalid Input"))
// 	}
// }

// func helper.WriteToResponseBody(writer http.ResponseWriter, response interface{}) {
// 	writer.Header().Add("Content-Type", "application/json")
// 	encoder := json.NewEncoder(writer)
// 	err := encoder.Encode(response)
// 	if err != nil {
// 		panic(NewEntityError("Invalid Output"))
// 	}
// }
