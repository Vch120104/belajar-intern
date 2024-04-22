package exceptions

import (
	"encoding/json"
	"net/http"
)

// DuplicateError adalah struktur untuk kesalahan data duplikat
type DuplicateError struct {
	Error string
}

// NewDuplicateError membuat instance baru dari DuplicateError
func NewDuplicateError(error string) DuplicateError {
	return DuplicateError{Error: error}
}

// DuplicateException menangani kasus exception ketika terjadi duplikasi data
func DuplicateException(w http.ResponseWriter, message string) {
	errResponse := CustomError{
		StatusCode: http.StatusConflict,
		Message:    "Duplicate Entry",
		Error:      message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(errResponse)
}
