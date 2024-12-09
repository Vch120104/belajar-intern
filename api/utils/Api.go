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
			// Implement exponential backoff to avoid flooding the API
			time.Sleep(time.Duration(math.Pow(2, float64(retry))) * time.Second) // Exponential backoff
			continue
		}

		// If there are server-side errors, retry with backoff, but avoid retries for 404
		if isServerError(err) {
			log.Printf("Retry %d/%d due to server error: %v", retry+1, maxRetries, err)
			time.Sleep(time.Duration(math.Pow(2, float64(retry))) * time.Second) // Exponential backoff
			continue
		}

		break
	}

	return fmt.Errorf("request failed after %d retries: %w", maxRetries, err)
}

// Helper function to detect server-side errors
func isServerError(err error) bool {
	return err != nil && (err.Error() == "500 Internal Server Error" || err.Error() == "503 Service Unavailable")
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

// handleResponse processes the HTTP response with specific checks
func handleResponse(resp *http.Response, result interface{}) error {
	// Check for specific HTTP status codes
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		// OK and Created responses are successful
		break
	case http.StatusBadRequest:
		// 400 Bad Request
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Bad Request (400) - Invalid syntax: %s", string(body))
		return fmt.Errorf("bad request (400): %s", string(body))
	case http.StatusUnauthorized:
		// 401 Unauthorized
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Unauthorized (401) - Authentication required: %s", string(body))
		return fmt.Errorf("unauthorized (401): %s", string(body))
	case http.StatusForbidden:
		// 403 Forbidden
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Forbidden (403) - Access denied: %s", string(body))
		return fmt.Errorf("forbidden (403): %s", string(body))
	case http.StatusNotFound:
		// 404 Not Found
		log.Printf("Not Found (404) - Resource not found for URL: %s", resp.Request.URL)
		return fmt.Errorf("not found (404): %s", resp.Request.URL)
	case http.StatusMethodNotAllowed:
		// 405 Method Not Allowed
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Method Not Allowed (405) - Invalid method for resource: %s", string(body))
		return fmt.Errorf("method not allowed (405): %s", string(body))
	case http.StatusInternalServerError:
		// 500 Internal Server Error
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Internal Server Error (500) - Server error at URL: %s, response: %s", resp.Request.URL, string(body))
		return fmt.Errorf("internal server error (500) at URL: %s: %s", resp.Request.URL, string(body))
	case http.StatusNotImplemented:
		// 501 Not Implemented
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Not Implemented (501) - Method not recognized: %s", string(body))
		return fmt.Errorf("not implemented (501): %s", string(body))
	case http.StatusBadGateway:
		// 502 Bad Gateway
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Bad Gateway (502) - Invalid response from upstream: %s", string(body))
		return fmt.Errorf("bad gateway (502): %s", string(body))
	case http.StatusServiceUnavailable:
		// 503 Service Unavailable
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Service Unavailable (503) - Service is temporarily unavailable: %s", string(body))
		return fmt.Errorf("service unavailable (503): %s", string(body))
	case http.StatusGatewayTimeout:
		// 504 Gateway Timeout
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Gateway Timeout (504) - Timeout from upstream: %s", string(body))
		return fmt.Errorf("gateway timeout (504): %s", string(body))
	default:
		// Handle any other unexpected status codes
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Unexpected status code %d for URL: %s, response: %s", resp.StatusCode, resp.Request.URL, string(body))
		return fmt.Errorf("unexpected status code %d at URL: %s: %s", resp.StatusCode, resp.Request.URL, string(body))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var generalResponse GeneralResponse
	if err := json.Unmarshal(bodyBytes, &generalResponse); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if err := json.Unmarshal(generalResponse.Data, result); err != nil {
		return fmt.Errorf("failed to unmarshal nested data: %w", err)
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
