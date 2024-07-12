package payloads

import (
	"encoding/json"
	"net/http"
)

type ResponseAuth struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

func WriteResponseToken(w http.ResponseWriter, message string, token string, status int) {
	res := ResponseAuth{
		Status:  status,
		Message: message,
		Token:   token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}
