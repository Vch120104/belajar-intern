package securities

import (
	"after-sales/api/config"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
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

// mengautentikasi user dan mendapatkan token
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
		log.Error("Failed to create request: ", err)
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	clientHttp := &http.Client{}
	resp, err := clientHttp.Do(req)
	if err != nil {
		log.Error("Failed to send request: ", err)
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn("Authentication failed with status: ", resp.StatusCode)
		return "", errors.New("failed to authenticate user")
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		log.Error("Failed to decode response: ", err)
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if loginResp.StatusCode != 200 {
		log.Warn("Login failed: ", loginResp.Message)
		return "", errors.New(loginResp.Message)
	}

	log.Info("User authenticated successfully")
	return loginResp.Data.Token, nil
}

// memvalidasi token JWT dari request
func GetAuthentication(r *http.Request) error {
	tokenString := ExtractToken(r)
	if tokenString == "" {
		log.Warn("Token not found")
		return errors.New("token not found")
	}

	log.Info("Extracted Token: ", tokenString)

	token, err := VerifyToken(tokenString)
	if err != nil {
		log.Error("Token verification failed: ", err)
		return err
	}

	if !token.Valid {
		log.Warn("Invalid token detected")
		return errors.New("invalid token")
	}

	log.Info("Token is valid")
	return nil
}

// memverifikasi token JWT dengan HS256
func VerifyToken(tokenString string) (*jwt.Token, error) {
	if tokenString == "" {
		log.Warn("Empty token string")
		return nil, errors.New("token not found")
	}

	log.Info("Parsing token: ", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			log.Warn("Unexpected signing method: ", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		log.Info("Using JWT secret key")
		return []byte(config.EnvConfigs.JWTKey), nil
	})

	if err != nil {
		log.Error("JWT parsing failed: ", err)
		return nil, err
	}

	return token, nil
}

func ExtractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	log.Info("Authorization Header: ", authHeader)

	if authHeader == "" {
		log.Warn("No Authorization header found")
		return ""
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) == 2 && strings.ToLower(bearerToken[0]) == "bearer" {
		log.Info("Extracted Bearer Token: ", bearerToken[1])
		return bearerToken[1]
	}

	log.Warn("Invalid Bearer token format")
	return ""
}
