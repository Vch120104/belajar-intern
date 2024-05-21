package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func Get(url string, responseData interface{}, client *http.Client) error {
	if client == nil {
		client = &http.Client{}
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch data from the external API, status code: %d", response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(responseData); err != nil {
		return err
	}

	return nil
}
