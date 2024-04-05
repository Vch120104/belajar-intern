package exceptions

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func AuthorizeException(w http.ResponseWriter, r *http.Request, message string) {
	res := ErrorResponse{
		Success: false,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(res)
}

type AuthorizationError struct {
	Error string
}

func NewAuthorizationError(error string) AuthorizationError {
	return AuthorizationError{Error: error}
}
