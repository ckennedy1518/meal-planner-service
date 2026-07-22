package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type AuthService struct {
	jwtSecret string
}

func NewAuthService(secret string) (*AuthService, error) {
	authService := &AuthService{jwtSecret: secret}

	return authService, nil
}

func (a *AuthService) ValidateAuthHeader(authHeader string) (map[string]any, error) {
	godotenv.Load()

	req, err := http.NewRequest("GET", "https://djxvcplftiwhdluniobh.supabase.co/auth/v1/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token: %s", resp.Status)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}
