package account

import db "go-rest-api-boilerplate/internal/db/sqlc"

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(email Email, password Password) *User {
	return &User{Email: email.String(), Password: password.String()}
}

func (u *User) ToParams() *db.AddParams {
	return &db.AddParams{
		Email:    u.Email,
		Password: u.Password,
	}
}
