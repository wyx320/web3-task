package models

type UserRegisterDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
