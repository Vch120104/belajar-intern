package exceptions

import (
	"encoding/json"
	"net/http"
)

// CustomError represents a custom error response structure
type CustomError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}

// NotFoundError represents the not found error
type NotFoundError struct {
	Message string
}

// Implementasi metode Error untuk NotFoundError
func (e NotFoundError) Error() string {
	return e.Message
}

// NewNotFoundError creates a new NotFoundError instance
func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{Message: message}
}

// NotFoundException handles the not found exception
func NotFoundException(w http.ResponseWriter, message string) {
	errResponse := CustomError{
		StatusCode: http.StatusNotFound,
		Message:    "Not Found",
		Error:      message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(errResponse)
}
