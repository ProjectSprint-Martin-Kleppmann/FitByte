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
	Name       string  `json:"name" validate:"omitempty,min=2,max=60"`
	ImageURI   string  `json:"imageUri" validate:"omitempty,uri"`
	Preference string  `json:"preference" validate:"omitempty,oneof=CARDIO WEIGHT"`
	WeightUnit string  `json:"weightUnit" validate:"omitempty,oneof=KG LBS"`
	HeightUnit string  `json:"heightUnit" validate:"omitempty,oneof=CM INCH"`
	Weight     float64 `json:"weight" validate:"omitempty,min=10,max=1000"`
	Height     float64 `json:"height" validate:"omitempty,min=3,max=250"`
}

type PatchProfileRequest struct {
	Preference string  `json:"preference" validate:"required,oneof=CARDIO WEIGHT"`
	WeightUnit string  `json:"weightUnit" validate:"required,oneof=KG LBS"`
	HeightUnit string  `json:"heightUnit" validate:"required,oneof=CM INCH"`
	Weight     float64 `json:"weight" validate:"required,min=10,max=1000"`
	Height     float64 `json:"height" validate:"required,min=3,max=250"`
	Name       string  `json:"name" validate:"required,min=2,max=60"`
	ImageURI   string  `json:"imageUri" validate:"required,uri"`
}

type ProfileResponse struct {
	Preference string  `json:"preference"`
	WeightUnit string  `json:"weightUnit"`
	HeightUnit string  `json:"heightUnit"`
	Weight     float64 `json:"weight"`
	Height     float64 `json:"height"`
	Name       string  `json:"name"`
	Email      string  `json:"email,omitempty"`
	ImageURI   string  `json:"imageUri"`
}
