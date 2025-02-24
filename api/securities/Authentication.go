package securities

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Client   string `json:"client"`
}

type LoginResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       struct {
		Token   string `json:"token"`
		UserID  int    `json:"user_id"`
		Company int    `json:"company"`
	} `json:"data"`
}

func LoginUser(username, password, client string) (string, error) {
	authServiceURL := fmt.Sprintf("%s/auth/login", config.EnvConfigs.UserServiceUrl)

	requestBody, _ := json.Marshal(LoginRequest{
		Username: username,
		Password: password,
		Client:   client,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authServiceURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "500 Internal Server Error - Unable to create authentication request",
			Err:        err,
		}
	}
	req.Header.Set("Content-Type", "application/json")

	clientHttp := &http.Client{}
	resp, err := clientHttp.Do(req)
	if err != nil {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "500 Internal Server Error - Unable to send authentication request",
			Err:        err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: resp.StatusCode,
			Message:    "401 Unauthorized - Authentication failed",
		}
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "500 Internal Server Error - Failed to process authentication response",
			Err:        err,
		}
	}

	if loginResp.StatusCode != 200 {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "401 Unauthorized - Invalid credentials",
		}
	}

	return loginResp.Data.Token, nil
}

// token JWT dari request
func GetAuthentication(r *http.Request) error {
	tokenString := ExtractToken(r)
	if tokenString == "" {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Missing authentication token",
		}
	}

	token, err := VerifyToken(tokenString)
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid authentication token",
			Err:        err,
		}
	}

	if !token.Valid {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Authentication token is not valid",
		}
	}

	return nil
}

// verifikasi token JWT dengan HS256 atau HS384
func VerifyToken(tokenString string) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Authentication token is required",
		}
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Unsupported token signing method",
			}
		}

		return []byte(config.EnvConfigs.JWTKey), nil
	})

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Authentication token verification failed",
			Err:        err,
		}
	}

	return token, nil
}

func ExtractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "No Authorization header found"
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) == 2 && strings.ToLower(bearerToken[0]) == "bearer" {
		return bearerToken[1]
	}

	return authHeader
}

// mengambil claims dari token JWT
func GetTokenClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid authentication token claims",
		}
	}

	return claims, nil
}
