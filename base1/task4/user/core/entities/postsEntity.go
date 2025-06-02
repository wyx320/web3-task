package entities

import "time"

// 帖子实体
type PostEntity struct {
	Id      uint64 `gorm:"primary_key;autoIncrement"`
	Title   string
	Content string
	UserId  uint64

	CreateAt uint64
	UpdateAt time.Time
}

func (PostEntity) TableName() string {
	return "posts"
}
