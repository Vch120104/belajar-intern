package middlewares

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/securities"

	// "encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// Inisialisasi konfigurasi logger
func init() {
	logger.Formatter = &logrus.JSONFormatter{} // Ubah formatter sesuai kebutuhan
	logger.Level = logrus.InfoLevel            // Ubah level log sesuai kebutuhan
}

func SetupAuthenticationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := securities.GetAuthentication(r)

			if err != nil {
				exceptions.NewAuthorizationException(w, r, &exceptions.BaseErrorResponse{
					Err: err,
				})
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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Authorization-Key, accept, origin, Cache-Control, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err := securities.GetAuthentication(r)
	if err != nil {
		exceptions.NewAuthorizationException(w, r, &exceptions.BaseErrorResponse{
			Err: err,
		})
		return
	}

	middleware.Handler.ServeHTTP(w, r)
}

// func NotFoundHandler(next http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		defer func() {
// 			if r := recover(); r != nil {
// 				notFoundErr, ok := r.(exceptions.NotFoundError)
// 				if !ok {
// 					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 					return
// 				}

// 				w.Header().Set("Content-Type", "application/json")
// 				w.WriteHeader(http.StatusNotFound)
// 				errResponse := exceptions.CustomError{
// 					StatusCode: http.StatusNotFound,
// 					Message:    "Not Found",
// 					Error:      notFoundErr.Error(), // Panggil metode Error() untuk mendapatkan pesan kesalahan
// 				}
// 				json.NewEncoder(w).Encode(errResponse)
// 			}
// 		}()

// 		next.ServeHTTP(w, r)
// 	}

// 	return http.HandlerFunc(fn)
// }

// Logger adalah middleware untuk logging setiap request yang masuk
func Logger(next http.Handler) http.Handler {
	// Create a new logger middleware with the default log formatter and logger
	handler := middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger})
	// Then, call the middleware with the provided handler and return the result
	return handler(next)
}

// Recoverer adalah middleware untuk pemulihan dari panic dan pengiriman tanggapan 500
func Recoverer(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}
