package jobs

import "time"

const PostCreated = "post.created"

type PostCreatedMessage struct {
	Type      string    `json:"type"`
	PostID    string    `json:"post_id"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}