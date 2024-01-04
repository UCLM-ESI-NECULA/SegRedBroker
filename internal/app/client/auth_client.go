package client

import (
	"bytes"
	"net/http"
	"os"
	"seg-red-broker/internal/app/dao"
)

// AuthClient struct holds the HTTP client and the base URL for the Auth service
type AuthClient struct {
	Client
}

// NewAuthClient creates a new instance of AuthClient
func NewAuthClient() *AuthClient {
	return &AuthClient{
		Client{
			HttpClient: &http.Client{},
			BaseURL:    os.Getenv("AUTH_SERVICE_BASE_URL"),
		},
	}
}

// Signup sends a signup request to the Auth service
func (client *AuthClient) Signup(username, password string) (string, error) {
	resp, requestErr := client.makeRequest(http.MethodPost, "/signup", bytes.NewBuffer(dao.UserToJson(username, password)))
	if requestErr != nil {
		return "", requestErr
	}

	// Unmarshal JSON to the target type
	var token dao.Token
	marshalError := getBody(resp, &token)
	if marshalError != nil {
		return "", marshalError
	}
	return token.Token, nil
}

// Login sends a login request to the Auth service
func (client *AuthClient) Login(username, password string) (string, error) {
	resp, err := client.makeRequest(http.MethodPost, "/login", bytes.NewBuffer(dao.UserToJson(username, password)))
	token, err := getToken(resp, err)
	return token.Token, err
}

// ValidateToken sends a validate token request to the Auth service
func (client *AuthClient) ValidateToken(tokenString string) (string, error) {
	resp, err := client.makeRequest(http.MethodPost, "/checkToken", bytes.NewBuffer(dao.TokenToJson(tokenString)))
	var user dao.User
	err = getBody(resp, &user)
	return user.Username, err
}
