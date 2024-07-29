package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const serverUrl = ""
const SalesURL = "http://10.1.32.26:8000/sales-service/v1"
const GeneralURL = "http://10.1.32.26:8000/general-service/v1"

type ResponseBody struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

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

// Get function for JSON object response
func Get(url string, data interface{}, body interface{}) error {

	client := &http.Client{}
	var buf bytes.Buffer

	// Jika ada parameter Body/body request untuk getnya
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return err
	}

	var responseBody APIResponse

	newRequest, err := http.NewRequest("GET", serverUrl+url, &buf)

	if err != nil {
		return err
	}

	newResponse, err := client.Do(newRequest)

	if err != nil {
		return err
	}

	defer newResponse.Body.Close()
	defer client.CloseIdleConnections()

	responseBody = APIResponse{
		Data: data,
	}

	//jika status != ok, maka return nothing
	if newResponse.StatusCode != http.StatusOK {
		return nil
	}

	if err := json.NewDecoder(newResponse.Body).Decode(&responseBody); err != nil {
		return err
	}

	return nil

}

// GetWithPagination function for JSON object response
// get data from url with pagination, the returned data is in form of APIPaginationResponse
func GetWithPagination(url string, pagination APIPaginationResponse, body interface{}) (APIPaginationResponse, error) {

	client := &http.Client{}
	var buf bytes.Buffer

	// Jika ada parameter Body/body request untuk getnya
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return APIPaginationResponse{}, err
	}

	var responseBody APIPaginationResponse

	newRequest, err := http.NewRequest("GET", serverUrl+url, &buf)

	if err != nil {
		return APIPaginationResponse{}, err
	}

	newResponse, err := client.Do(newRequest)

	if err != nil {
		return APIPaginationResponse{}, err
	}

	defer newResponse.Body.Close()
	defer client.CloseIdleConnections()

	responseBody = APIPaginationResponse{
		Data: pagination.Data,
	}

	//jika status != ok, maka return nothing
	if newResponse.StatusCode != http.StatusOK {
		return APIPaginationResponse{}, nil
	}

	if err := json.NewDecoder(newResponse.Body).Decode(&responseBody); err != nil {
		return APIPaginationResponse{}, err
	}

	return responseBody, nil

}

// GetArray function for JSON array response
func GetArray(url string, data interface{}, body interface{}) error {
	client := &http.Client{}
	var buf bytes.Buffer

	// Jika ada parameter Body/body request untuk getnya
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return err
		}
	}

	newRequest, err := http.NewRequest("GET", url, &buf)
	if err != nil {
		return err
	}

	newResponse, err := client.Do(newRequest)
	if err != nil {
		return err
	}
	defer newResponse.Body.Close()
	defer client.CloseIdleConnections()

	//jika status != ok, maka return nothing
	if newResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", newResponse.StatusCode)
	}

	// Decode the response body into a temporary structure to extract the array
	var responseBody struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(newResponse.Body).Decode(&responseBody); err != nil {
		return err
	}

	// Decode the array into the provided data interface
	if err := json.Unmarshal(responseBody.Data, &data); err != nil {
		return err
	}

	return nil
}

// GetWithPaginationArray function for JSON array response
func GetWithPaginationArray(url string, pagination APIPaginationResponse, body interface{}) (APIPaginationResponse, error) {

	client := &http.Client{}
	var buf bytes.Buffer

	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return APIPaginationResponse{}, err
		}
	}

	newRequest, err := http.NewRequest("GET", serverUrl+url, &buf)
	if err != nil {
		return APIPaginationResponse{}, err
	}

	newResponse, err := client.Do(newRequest)
	if err != nil {
		return APIPaginationResponse{}, err
	}
	defer newResponse.Body.Close()
	defer client.CloseIdleConnections()

	responseBody := APIPaginationResponse{
		Data: pagination.Data,
	}

	if newResponse.StatusCode != http.StatusOK {
		return APIPaginationResponse{}, nil
	}

	if err := json.NewDecoder(newResponse.Body).Decode(&responseBody.Data); err != nil {
		return APIPaginationResponse{}, err
	}

	return responseBody, nil
}
