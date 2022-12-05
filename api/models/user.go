package models

import (
	"time"
)

type User struct {
	ID              int64     `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	DateOfBirth     string    `json:"dob"`
	Email           string    `json:"email"`
	PhoneNumber     *string   `json:"phone_number"`
	Gender          string    `json:"gender"`
	ProfileImageUrl *string   `json:"profile_image_url"`
	Address         *string   `json:"address"`
	Type            string    `json:"type"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	FirstName       string  `json:"first_name" binding:"required,min=2,max=30"`
	LastName        string  `json:"last_name" binding:"required,min=2,max=30"`
	DateOfBirth     string  `json:"dob" binding:"required"`
	Email           string  `json:"email" binding:"required,email"`
	PhoneNumber     *string `json:"phone_number"`
	Gender          string  `json:"gender" binding:"required,oneof=male female"`
	Password        string  `json:"password" binding:"required,min=6,max=16"`
	ProfileImageUrl *string `json:"profile_image_url"`
	Type            string  `json:"type" binding:"required,oneof=superadmin guest owner"`
	Address         *string `json:"address"`
}

type GetUsersResponse struct {
	Users []*User `json:"users"`
	Count int32   `json:"count"`
}
