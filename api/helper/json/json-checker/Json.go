package jsonchecker

import (
	"after-sales/api/exceptions"
	"encoding/json"
	"fmt"
	"net/http"
)

func ReadFromRequestBody(request *http.Request, result interface{}) *exceptions.BaseErrorResponse {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	fmt.Println(result)
	if err != nil {
		// errorMsg := fmt.Sprintf("Failed to decode request body: %s", err.Error())
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
		// errorMsg := fmt.Sprintf("Failed to encode response body: %s", err.Error())
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        err,
		}
	}
	return nil
}
