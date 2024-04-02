package middlewares

import (
	"after-sales/api/exceptions"
	"after-sales/api/securities"
	"net/http"
)

func SetupAuthenticationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := securities.GetAuthentication(r)

			if err != nil {
				exceptions.AuthorizeException(w, r, err.Error())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

func (middleware *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err := securities.GetAuthentication(r)
	if err != nil {
		exceptions.AuthorizeException(w, r, err.Error())
		return
	}

	middleware.Handler.ServeHTTP(w, r)
}
