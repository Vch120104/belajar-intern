package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"reflect"
	"time"
)

const (
	FinanceURL     = "https://testing-backendims.indomobil.co.id/finance-service/v1/"
	SalesURL       = "https://testing-backendims.indomobil.co.id/sales-service/v1/"
	GeneralURL     = "https://testing-backendims.indomobil.co.id/general-service/v1/"
	AftersalesURL  = "https://testing-backendims.indomobil.co.id/aftersales-service/v1/"
	requestTimeout = 10 * time.Second
	maxRetries     = 3 // Number of retries for failed requests
)

type ResponseBody struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type GeneralResponse struct {
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
	Data       json.RawMessage `json:"data"`
}

type APIPaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	TotalPages int         `json:"total_pages"`
	Limit      int         `json:"limit"`
	TotalRows  int64       `json:"total_rows"`
}

// Reusable HTTP client with timeout and transport settings
var httpClient = &http.Client{
	Timeout: requestTimeout,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	},
}

// handleResponse handles the HTTP response, checking status codes and decoding response bodies
func handleResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	// Log the status code for better debugging
	//log.Printf("Received HTTP status: %d", resp.StatusCode)

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResponse ResponseBody
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("error decoding error response: %w", err)
		}
		//log.Printf("Error response: %s, status code: %d", errorResponse.Message, resp.StatusCode)
		return fmt.Errorf("error response: %s, status code: %d", errorResponse.Message, resp.StatusCode)
	}

	// Decode response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// Log the raw response body for debugging
	//log.Printf("Raw response body: %s", bodyBytes)

	// Unmarshal into GeneralResponse to check the data field
	var generalResponse GeneralResponse
	if err := json.Unmarshal(bodyBytes, &generalResponse); err != nil {
		return fmt.Errorf("error unmarshalling general response: %w, body: %s", err, bodyBytes)
	}

	// Determine if result is a slice or a single object
	if reflect.TypeOf(result).Kind() == reflect.Slice {
		// Expecting an array
		if err := json.Unmarshal(generalResponse.Data, result); err != nil {
			return fmt.Errorf("error unmarshalling nested data into slice: %w", err)
		}
	} else {
		// If the result is not a slice, check if data is an array
		var tempData json.RawMessage
		if err := json.Unmarshal(generalResponse.Data, &tempData); err != nil {
			return fmt.Errorf("error unmarshalling nested data into temp data: %w", err)
		}

		// Try unmarshalling into the expected struct
		if err := json.Unmarshal(tempData, result); err != nil {
			return fmt.Errorf("error unmarshalling nested data: %w", err)
		}
	}

	return nil
}

// CallAPI is a generic function for making API calls with retry logic
func CallAPI(method, url string, body interface{}, result interface{}) error {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	for retry := 0; retry < maxRetries; retry++ {
		err = makeRequest(method, url, reqBody, result)
		if err == nil {
			return nil
		}

		log.Printf("Retry attempt %d for %s request to %s failed: %v", retry+1, method, url, err)

		// Use exponential backoff
		time.Sleep(time.Duration(math.Pow(2, float64(retry))) * time.Second)
	}

	return fmt.Errorf("request failed after %d retries: %w", maxRetries, err)
}

// Helper function for making the actual HTTP request
func makeRequest(method, url string, reqBody []byte, result interface{}) error {

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}

	return handleResponse(resp, result)
}

// Helper functions for CRUD operations
// GET request
func Get(url string, result interface{}, params interface{}) error {
	return CallAPI("GET", url, params, result)
}

// POST request
func Post(url string, body interface{}, result interface{}) error {
	return CallAPI("POST", url, body, result)
}

// PUT request
func Put(url string, body interface{}, result interface{}) error {
	return CallAPI("PUT", url, body, result)
}

// DELETE request
func Delete(url string, body interface{}, result interface{}) error {
	return CallAPI("DELETE", url, body, result)
}

// GetArray handles array responses
func GetArray(url string, body interface{}, response interface{}) error {
	// Set up HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil) // Adjust method and body as necessary
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	// Unmarshal the response body
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// BatchRequest supports sending multiple requests in one call
func BatchRequest(url string, requests []interface{}, results []interface{}) error {
	if len(requests) != len(results) {
		return fmt.Errorf("requests and results length must match")
	}

	for i, req := range requests {
		if err := Post(url, req, &results[i]); err != nil {
			return fmt.Errorf("error processing request %d: %w", i, err)
		}
	}

	return nil
}

// GetWithPagination makes a paginated GET request
func GetWithPagination(url string, pagination *APIPaginationResponse, params interface{}) error {
	return CallAPI("GET", url, params, pagination)
}

// GetArrayWithPagination retrieves an array with pagination
func GetArrayWithPagination(url string, pagination *APIPaginationResponse, params interface{}) error {
	return CallAPI("GET", url, params, pagination)
}
