package models

import "time"

type PostDto struct {
	Id      uint64 `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserId  uint64 `json:"user_id"`

	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
	CreateBy uint64    `json:"create_by"`
	UpdateBy uint64    `json:"update_by"`
}
