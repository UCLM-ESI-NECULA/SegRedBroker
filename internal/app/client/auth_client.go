package client

import (
	"github.com/go-resty/resty/v2"
	"os"
	"seg-red-broker/internal/app/common"
	"seg-red-broker/internal/app/dao"
)

// AuthClient struct holds the HTTP client and the base URL for the Auth service
type AuthClient struct {
	Client *resty.Client
}

// NewAuthClient creates a new instance of AuthClient
func NewAuthClient() *AuthClient {
	cl := resty.New()
	cl.
		SetBaseURL(os.Getenv("AUTH_SERVICE_BASE_URL")).
		SetHeader("Accept", "application/json").
		SetError(&common.APIError{})
	return &AuthClient{
		Client: cl,
	}
}

// Signup sends a signup request to the Auth service
func (client *AuthClient) Signup(username, password string) (*dao.Token, error) {
	resp, err := client.Client.R().
		SetBody(dao.User{Username: username, Password: password}).
		SetResult(&dao.Token{}).
		Post("/signup")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, resp.Error().(*common.APIError)
	}
	return resp.Result().(*dao.Token), nil
}

// Login sends a login request to the Auth service
func (client *AuthClient) Login(username, password string) (*dao.Token, error) {
	resp, err := client.Client.R().
		SetBody(dao.User{Username: username, Password: password}).
		SetResult(&dao.Token{}).
		Post("/login")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, resp.Error().(*common.APIError)
	}
	return resp.Result().(*dao.Token), nil
}

// ValidateToken sends a validate token request to the Auth service
func (client *AuthClient) ValidateToken(tokenString string) (*dao.User, error) {
	resp, err := client.Client.R().
		SetHeader("Authorization", tokenString).
		SetResult(&dao.User{}).
		Post("/checkToken")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, resp.Error().(*common.APIError)
	}
	return resp.Result().(*dao.User), nil
}
