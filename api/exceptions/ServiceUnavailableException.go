package exceptions

import (
	"encoding/json"
	"net/http"
)

// ServiceUnavailableException menangani kasus exception ketika layanan tidak tersedia
func ServiceUnavailableException(w http.ResponseWriter, message string) {
	errResponse := CustomError{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service Unavailable",
		Error:      message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	json.NewEncoder(w).Encode(errResponse)
}

// ServiceUnavailableError adalah struktur untuk kesalahan layanan tidak tersedia
type ServiceUnavailableError struct {
	Error string
}

// NewServiceUnavailableError membuat instance baru dari ServiceUnavailableError
func NewServiceUnavailableError(error string) ServiceUnavailableError {
	return ServiceUnavailableError{Error: error}
}
