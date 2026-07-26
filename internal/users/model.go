package users

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserWithPassword struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"passwordHash"`
	CreatedAt    time.Time `json:"createdAt"`
}

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type createUserParams struct {
	Username     string
	PasswordHash string
}
