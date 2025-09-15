package models

import "time"

type UserInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

type UserResponse struct {
	ID        string `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}