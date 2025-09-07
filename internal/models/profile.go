package models

import "gorm.io/gorm"

type RegisterResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type Profile struct {
	gorm.Model
	Email      string  `json:"email" validate:"required,email"`
	Password   string  `json:"password" validate:"required,min=8,max=32"`
	Name       string  `json:"name" validate:"omitempty,min=1,max=255"`
	ImageURI   string  `json:"imageUri" validate:"omitempty,url"`
	Preference string  `json:"preference" validate:"omitempty,max=255"`
	WeightUnit string  `json:"weightUnit" validate:"omitempty,oneof=kg lbs"`
	HeightUnit string  `json:"heightUnit" validate:"omitempty,oneof=cm ft in"`
	Weight     float64 `json:"weight" validate:"omitempty,min=0,max=999.99"`
	Height     float64 `json:"height" validate:"omitempty,min=0,max=999.99"`
}

type ProfileResponse struct {
	Email      string   `json:"email"`
	Name       string   `json:"name"`
	ImageURI   string   `json:"imageUri"`
	Preference string   `json:"preference"`
	WeightUnit string   `json:"weightUnit"`
	HeightUnit string   `json:"heightUnit"`
	Weight     *float64 `json:"weight"`
	Height     *float64 `json:"height"`
}
