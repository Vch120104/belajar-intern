package exceptions

import (
	"encoding/json"
	"net/http"
)

type OldError struct {
	Success bool        `json:"Success"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

func AppException(w http.ResponseWriter, message string) {
	res := OldError{
		Success: false,
		Message: message,
		Data:    nil,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(res)
}

type AppExceptionError struct {
	Error string
}

func NewAppExceptionError(error string) AppExceptionError {
	return AppExceptionError{Error: error}
}
