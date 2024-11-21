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
	"os"
	"reflect"
	"time"
)

const (
	FinanceURL     = "https://testing-backendims.indomobil.co.id/finance-service/v1/"
	SalesURL       = "https://testing-backendims.indomobil.co.id/sales-service/v1/"
	GeneralURL     = "https://testing-backendims.indomobil.co.id/general-service/v1/"
	AftersalesURL  = "https://testing-backendims.indomobil.co.id/aftersales-service/v1/"
	requestTimeout = 10 * time.Second
	maxRetries     = 2 // Retry limit for failed requests
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

// HTTP client with transport settings
var httpClient = &http.Client{
	Timeout: requestTimeout,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	},
}

// CallAPI makes a request with retries and handles errors
func CallAPI(method, url string, body interface{}, result interface{}) error {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	for retry := 0; retry <= maxRetries; retry++ {
		err = makeRequest(method, url, reqBody, result)
		if err == nil {
			return nil
		}

		// Log and retry only for specific errors
		if os.IsTimeout(err) || err.Error() == "context deadline exceeded" {
			log.Printf("Retry %d/%d: %v", retry+1, maxRetries, err)
			time.Sleep(time.Duration(math.Pow(2, float64(retry))) * time.Second) // Exponential backoff
			continue
		}

		break // Stop retrying for other errors
	}

	return fmt.Errorf("request failed after %d retries: %w", maxRetries, err)
}

// makeRequest executes a single HTTP request
func makeRequest(method, url string, reqBody []byte, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	return handleResponse(resp, result)
}

// handleResponse processes the HTTP response
func handleResponse(resp *http.Response, result interface{}) error {
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusNotFound {
			log.Printf("Data not found (404) for URL: %s", resp.Request.URL)
			return fmt.Errorf("data not found")
		}

		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var generalResponse GeneralResponse
	if err := json.Unmarshal(bodyBytes, &generalResponse); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Unmarshal data into result
	if reflect.TypeOf(result).Kind() == reflect.Slice {
		if err := json.Unmarshal(generalResponse.Data, result); err != nil {
			return fmt.Errorf("failed to unmarshal nested data (slice): %w", err)
		}
	} else {
		if err := json.Unmarshal(generalResponse.Data, result); err != nil {
			return fmt.Errorf("failed to unmarshal nested data: %w", err)
		}
	}

	return nil
}

// Helper functions for CRUD operations
func Get(url string, result interface{}, params interface{}) error {
	return CallAPI("GET", url, params, result)
}

func Post(url string, body interface{}, result interface{}) error {
	return CallAPI("POST", url, body, result)
}

func Put(url string, body interface{}, result interface{}) error {
	return CallAPI("PUT", url, body, result)
}

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
