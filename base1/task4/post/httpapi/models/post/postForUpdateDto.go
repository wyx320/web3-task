package models

type PostForUpdateDto struct {
	Id      uint64 `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
