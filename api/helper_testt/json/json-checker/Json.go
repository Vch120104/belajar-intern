package jsonchecker

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"encoding/json"
	"net/http"
)

func ReadFromRequestBody(request *http.Request, result interface{}) *exceptionsss_test.BaseErrorResponse {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        err,
		}
	}
	return nil
}

func WriteToResponseBody(writer http.ResponseWriter, response interface{}) *exceptionsss_test.BaseErrorResponse {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        err,
		}
	}
	return nil
}
