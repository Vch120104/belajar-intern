package payloads

import (
	"after-sales/api/helper"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type ResponsePagination struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Page       int         `json:"page"`
	Limit      int         `json:"page_limit"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

// ErrorResponse represents the structure of an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

func NewHandleError(w http.ResponseWriter, errorMessage string, statusCode int) {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	// Set the status code
	w.WriteHeader(statusCode)
	// Create the error response payload
	errorResponse := ErrorResponse{Error: errorMessage}
	// Convert the error response to JSON
	jsonResponse, err := json.Marshal(errorResponse)
	if err != nil {
		// If there's an error in marshalling the JSON response, log it
		http.Error(w, "Failed to marshal error response", http.StatusInternalServerError)
		return
	}
	// Write the JSON response to the response writer
	_, err = w.Write(jsonResponse)
	if err != nil {
		// If there's an error in writing the response, log it
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// Deprecated: please change to the latest one without *gin.Context
func HandleSuccess(c *gin.Context, data interface{}, message string, status int) {
	res := Response{
		StatusCode: status,
		Message:    message,
		Data:       data,
	}

	c.JSON(status, res)
}

// Deprecated: please change to the latest one without *gin.Context
func HandleSuccessPagination(c *gin.Context, data interface{}, message string, status int, limit int, page int, totalRows int64, totalPages int) {
	res := ResponsePagination{
		StatusCode: status,
		Message:    message,
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Data:       data,
	}
	c.JSON(status, res)
}

func NewHandleSuccess(writer http.ResponseWriter, data interface{}, message string, status int) {
	res := Response{
		StatusCode: status,
		Message:    message,
		Data:       data,
	}

	helper.WriteToResponseBody(writer, res)
}

func NewHandleSuccessPagination(writer http.ResponseWriter, data interface{}, message string, status int, limit int, page int, totalRows int64, totalPages int) {
	res := ResponsePagination{
		StatusCode: status,
		Message:    message,
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Data:       data,
	}

	helper.WriteToResponseBody(writer, res)
}
