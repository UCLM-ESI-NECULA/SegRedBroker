package service

import (
	"seg-red-broker/internal/app/client"
	"seg-red-broker/internal/app/dao"
)

type AuthServiceImpl struct {
	ac client.AuthClient
}

func NewAuthService(ac client.AuthClient) *AuthServiceImpl {
	return &AuthServiceImpl{ac}
}

type AuthService interface {
	Signup(username, password string) (*dao.Token, error)
	Login(username, password string) (*dao.Token, error)
	ValidateToken(tokenString string) (*dao.User, error)
}

// Signup handles the /signup endpoint
func (svc *AuthServiceImpl) Signup(username, password string) (*dao.Token, error) {
	return svc.ac.Signup(username, password)
}

// Login handles the /login endpoint
func (svc *AuthServiceImpl) Login(username, password string) (*dao.Token, error) {
	return svc.ac.Login(username, password)
}

// ValidateToken checks if the provided token string is valid and returns the corresponding user.
func (svc *AuthServiceImpl) ValidateToken(tokenString string) (*dao.User, error) {
	return svc.ac.ValidateToken(tokenString)
}
