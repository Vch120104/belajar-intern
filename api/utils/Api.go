package utils

import (
	"after-sales/api/exceptions"
	"bytes"
	"encoding/json"
	"net/http"
)

// const serverUrl = "http://10.1.32.26:8000/general-service"
const serverUrl = ""

// const serverUrl = "http://127.0.0.1:8000"

type APIResponse struct {
	Data interface{} `json:"data"`
}

type APIPaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	TotalPages int         `json:"total_pages"`
	Limit      int         `json:"limit"`
	TotalRows  int64       `json:"total_rows"`
}

// get data from url
func Get(url string, data interface{}, body interface{}) error {
	client := &http.Client{}
	var buf bytes.Buffer

	// Jika ada parameter Body/body request untuk getnya
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		panic(exceptions.NewBadRequestError(err.Error()))
	}

	var responseBody APIResponse

	newRequest, err := http.NewRequest("GET", serverUrl+url, &buf)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	newResponse, err := client.Do(newRequest)

	if err != nil {
		return nil
	}

	defer newResponse.Body.Close()
	defer client.CloseIdleConnections()

	responseBody = APIResponse{
		Data: data,
	}

	//jika status != ok, maka return nothing
	if newResponse.StatusCode != http.StatusOK {
		return nil
		// c.JSON(newResponse.StatusCode, gin.H{"error": "Failed to fetch data from the external API"})
		// return err
	}

	//decode body response
	err = json.NewDecoder(newResponse.Body).Decode(&responseBody)

	if err != nil {
		panic(exceptions.NewBadRequestError(err.Error()))
	}

	return nil
}

// get data from url with pagination, the returned data is in form of APIPaginationResponse
func GetWithPagination(url string, pagination APIPaginationResponse, body interface{}) (APIPaginationResponse, error) {
	client := &http.Client{}
	var buf bytes.Buffer

	// Jika ada parameter Body
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		panic(exceptions.NewBadRequestError(err.Error()))

	}

	var responseBody APIPaginationResponse

	newRequest, err := http.NewRequest("GET", serverUrl+url, &buf)

	if err != nil {
		panic(exceptions.NewBadRequestError(err.Error()))

	}

	newResponse, err := client.Do(newRequest)

	if err != nil {
		panic(exceptions.NewBadRequestError(err.Error()))

	}

	defer newResponse.Body.Close()
	defer client.CloseIdleConnections()

	responseBody = APIPaginationResponse{
		Data:       pagination.Data,
		Page:       pagination.Page,
		TotalPages: pagination.TotalPages,
		Limit:      pagination.Limit,
		TotalRows:  pagination.TotalRows,
	}

	if newResponse.StatusCode != http.StatusOK {
		return pagination, err
		// c.JSON(newResponse.StatusCode, gin.H{"error": "Failed to fetch data from the external API"})
		// return err
	}

	err = json.NewDecoder(newResponse.Body).Decode(&responseBody)

	if err != nil {
		panic(exceptions.NewBadRequestError(err.Error()))
	}

	return responseBody, err
}
