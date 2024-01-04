package dao

import "encoding/json"

type User struct {
	ID       int64
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserToJson(username, password string) []byte {
	userData, _ := json.Marshal(User{Username: username, Password: password})
	return userData
}

type Token struct {
	Token string `json:"access_token"`
}

func TokenToJson(token string) []byte {
	tokenData, _ := json.Marshal(Token{Token: token})
	return tokenData
}
