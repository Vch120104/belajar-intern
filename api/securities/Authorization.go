package securities

import (
	"after-sales/api/payloads"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

// mengekstrak informasi pengguna dari token JWT
func ExtractAuthToken(w http.ResponseWriter, r *http.Request) (*payloads.UserDetail, error) {

	tokenString := ExtractToken(r)
	if tokenString == "" {
		log.Warn("No token found in request")
		return nil, errors.New("token not found")
	}

	token, err := VerifyToken(tokenString)
	if err != nil {
		log.Error("Token verification failed: ", err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Warn("Invalid token claims")
		return nil, errors.New("invalid token")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		log.Warn("Token missing expiration time")
		return nil, errors.New("token expiration time missing")
	}
	expTime := time.Unix(int64(expFloat), 0)
	if time.Now().After(expTime) {
		log.Warn("Token has expired")
		return nil, errors.New("token has expired")
	}

	userID, _ := strconv.Atoi(fmt.Sprintf("%v", claims["user_id"]))
	companyID, _ := strconv.Atoi(fmt.Sprintf("%v", claims["company_id"]))
	role, _ := strconv.Atoi(fmt.Sprintf("%v", claims["role"]))

	authorized, ok := claims["authorized"].(bool)
	if !ok {
		log.Warn("Invalid authorized field in token")
		return nil, errors.New("invalid authorized field in token")
	}

	authDetail := payloads.UserDetail{
		UserId:    int32(userID),
		Username:  fmt.Sprintf("%s", claims["username"]),
		Authorize: fmt.Sprintf("%t", authorized),
		Role:      uint16(role),
		CompanyId: fmt.Sprintf("%d", companyID),
		Client:    fmt.Sprintf("%s", claims["client"]),
		IpAddress: fmt.Sprintf("%s", claims["ip_address"]),
	}

	return &authDetail, nil
}
