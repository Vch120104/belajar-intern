package jsonchecker

import (
	exceptions "after-sales/api/exceptions"
	"encoding/json"
	"net/http"
)

func ReadFromRequestBody(request *http.Request, result interface{}) *exceptions.BaseErrorResponse {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        err,
		}
	}
	return nil
}

func WriteToResponseBody(writer http.ResponseWriter, response interface{}) *exceptions.BaseErrorResponse {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        err,
		}
	}
	return nil
}
