package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

// contextKey adalah tipe khusus untuk kunci konteks
type contextKey string

const (
	// ContextKeyDBTrx adalah kunci konteks untuk transaksi database
	ContextKeyDBTrx contextKey = "db_trx"
)

// StatusInList checks if the given status is in the list
func StatusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}

// DBTransactionMiddleware sets up the database transaction middleware
func DBTransactionMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			txHandle := db.Begin()
			log.Print("beginning database transaction")

			defer func() {
				if r := recover(); r != nil {
					txHandle.Rollback()
				}
			}()

			// Pass the transaction handle to the request context
			ctx := context.WithValue(r.Context(), ContextKeyDBTrx, txHandle)

			// Capture the response status using ResponseRecorder
			recorder := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Call the next handler
			next.ServeHTTP(recorder, r.WithContext(ctx))

			if StatusInList(recorder.Status(), []int{http.StatusOK, http.StatusCreated}) {
				log.Print("Committing transactions")
				if err := txHandle.Commit().Error; err != nil {
					log.Print("trx commit error: ", err)
				}
			} else {
				log.Print("Rolling back transaction due to status code: ", recorder.Status())
				txHandle.Rollback()
			}
		})
	}
}

// Router is a wrapper for chi Router
type Router struct {
	*chi.Mux
}

// NewRouter creates a new Router
func NewRouter() *Router {
	r := chi.NewRouter()
	return &Router{r}
}

// Use applies middleware to the router
func (r *Router) Use(middlewares ...func(http.Handler) http.Handler) {
	r.Mux.Use(middlewares...)
}
