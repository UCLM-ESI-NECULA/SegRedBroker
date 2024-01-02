package dao

type User struct {
	ID       int64
	Username string
	Password string
}

type Token struct {
	Token string `json:"access_token"`
}
