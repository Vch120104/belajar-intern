package exceptions

import (
	"encoding/json"
	"net/http"
)

func EntityException(w http.ResponseWriter, message string) {
	res := OldError{
		Success: false,
		Message: message,
		Data:    nil,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(res)
}

type EntityError struct {
	Error string
}

func NewEntityError(error string) EntityError {
	return EntityError{Error: error}
}
