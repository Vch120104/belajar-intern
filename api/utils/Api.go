package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const serverUrl = ""
const SalesURL = "http://172.16.5.101/sales-service/v1"
const GeneralURL = "http://172.16.5.101/general-service/v1"

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
