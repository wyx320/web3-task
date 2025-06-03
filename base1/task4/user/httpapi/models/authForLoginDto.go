package models

type AuthForLoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
