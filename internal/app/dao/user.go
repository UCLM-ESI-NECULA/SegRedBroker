package dao

type User struct {
	ID       int64  `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type Token struct {
	Token string `json:"access_token"`
}
