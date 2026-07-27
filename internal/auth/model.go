package auth

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserWithPassword struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"passwordHash"`
	CreatedAt    time.Time `json:"createdAt"`
}

type signUpRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30,username"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type signInRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30,username"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type createUserParams struct {
	Username     string
	PasswordHash string
}
