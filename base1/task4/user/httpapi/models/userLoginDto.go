package models

type UserLoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
