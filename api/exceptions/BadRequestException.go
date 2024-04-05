package exceptions

import (
	"encoding/json"
	"net/http"
)

// BadRequestException menangani kasus exception ketika permintaan tidak valid
func BadRequestException(w http.ResponseWriter, message string) {
	errResponse := CustomError{
		StatusCode: http.StatusBadRequest,
		Message:    "Bad Request",
		Error:      message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errResponse)
}

// BadRequestError adalah struktur untuk kesalahan permintaan yang tidak valid
type BadRequestError struct {
	Error string
}

// NewBadRequestError membuat instance baru dari BadRequestError
func NewBadRequestError(error string) BadRequestError {
	return BadRequestError{Error: error}
}
