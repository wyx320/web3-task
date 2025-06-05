package models

type CommentCreateDto struct {
	Content string `json:"content"`
	PostId  uint64 `json:"post_id"`
}
