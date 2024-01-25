package middlewares

import (
	"after-sales/api/exceptions"
	"after-sales/api/securities"

	"github.com/gin-gonic/gin"
)

func SetupAuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := securities.GetAuthentication(c)

		if err != nil {
			exceptions.AuthorizeException(c, err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}

// type AuthMiddleware struct {
// 	Handler httprouter.Handle
// }

// func NewAuthMiddleware(handler httprouter.Handle) *AuthMiddleware {
// 	return &AuthMiddleware{Handler: handler}
// }
// func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
// 	var services services.AuthService
// 	writer.Header().Set("Access-Control-Allow-Origin", "*")
// 	writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 	writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH")

// 	if request.Method == "OPTIONS" {
// 		writer.WriteHeader(http.StatusNoContent)
// 		return
// 	}

// 	err := securities.GetAuthentication(request, services)
// 	if err != nil {
// 		panic(exceptions.NewAuthorizationError("You don't have access"))
// 	}

// 	// middleware.Handler.ServeHTTP(writer, request)
// 	middleware.Handler(writer, request, params)

// }
