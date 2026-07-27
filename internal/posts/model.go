package posts

import "time"

type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
}

type createPostRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=64"`
	Content string `json:"content" binding:"required,min=3,max=2048"`
}

type createPostParams struct {
	Title    string
	Content  string
	AuthorID string
}

type listParams struct {
	Limit  int `form:"limit" binding:"numeric"`
	Offset int `form:"offset" binding:"numeric"`
}
