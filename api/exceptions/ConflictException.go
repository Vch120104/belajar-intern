package exceptions

import (
	"encoding/json"
	"net/http"
)

// ConflictError adalah struktur untuk kesalahan konflik
type ConflictError struct {
	Error string
}

// NewConflictError membuat instance baru dari ConflictError
func NewConflictError(error string) ConflictError {
	return ConflictError{Error: error}
}

// ConflictException menangani kasus exception ketika terjadi konflik
func ConflictException(w http.ResponseWriter, message string) {
	errResponse := CustomError{
		StatusCode: http.StatusConflict,
		Message:    "Conflict",
		Error:      message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(errResponse)
}
