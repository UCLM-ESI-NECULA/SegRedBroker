package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"seg-red-broker/internal/app/dao"
)

const (
	authServiceBaseURL = "http://auth-service" // Replace with actual Auth Service URL
)

// AuthClient struct holds the HTTP client and the base URL for the Auth service
type AuthClient struct {
	HttpClient *http.Client
	BaseURL    string
}

// NewAuthClient creates a new instance of AuthClient
func NewAuthClient() *AuthClient {
	return &AuthClient{
		HttpClient: &http.Client{},
		BaseURL:    authServiceBaseURL,
	}
}

// makeRequest is a helper function to make an HTTP request with the base URL
func (client *AuthClient) makeRequest(method string, endpoint string, body *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest(method, client.BaseURL+endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return client.HttpClient.Do(req)
}

// Signup sends a signup request to the Auth service
func (client *AuthClient) Signup(username, password string) (string, error) {
	user := dao.User{Username: username, Password: password}
	userData, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	resp, err := client.makeRequest(http.MethodPost, "/signup", bytes.NewBuffer(userData))
	token, err := getToken(resp, err)
	return token.Token, err

}

// Login sends a login request to the Auth service
func (client *AuthClient) Login(username, password string) (string, error) {
	user := dao.User{Username: username, Password: password}
	userData, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	resp, err := client.makeRequest(http.MethodPost, "/login", bytes.NewBuffer(userData))
	token, err := getToken(resp, err)
	return token.Token, err
}

// ValidateToken sends a validate token request to the Auth service
func (client *AuthClient) ValidateToken(tokenString string) (bool, error) {
	token := dao.Token{Token: tokenString}
	tokenData, err := json.Marshal(token)
	if err != nil {
		return false, err
	}

	resp, err := client.makeRequest(http.MethodPost, "/checkToken", bytes.NewBuffer(tokenData))
	return err == nil && resp.StatusCode == http.StatusOK, nil
}

// getBody reads the response body and unmarshal JSON into the provided target interface
func getBody(resp *http.Response, target interface{}) error {
	if resp == nil {
		return nil // or an appropriate error
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Unmarshal JSON to the target type
	err = json.Unmarshal(body, &target)
	if err != nil {
		return err
	}

	return nil
}

func getToken(resp *http.Response, err error) (*dao.Token, error) {
	if err != nil {
		return nil, err
	}

	var token dao.Token
	err = getBody(resp, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
