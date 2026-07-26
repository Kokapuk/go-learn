package users

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type createUserDTO struct {
	Name string `json:"name"`
}
