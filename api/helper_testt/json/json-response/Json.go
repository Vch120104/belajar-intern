package jsonresponse

import (
	"after-sales/api/utils"
	"encoding/json"
	"errors"
	"net/http"
)

func WriteToResponseBody(writer http.ResponseWriter, response interface{}) error {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	if err != nil {
		return errors.New(utils.JsonError)
	}
	return nil
}
