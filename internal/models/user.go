package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
}
