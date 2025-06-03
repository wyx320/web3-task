package entities

import "time"

// 帖子实体
type PostEntity struct {
	Id      uint64 `gorm:"primary_key;autoIncrement"`
	Title   string
	Content string
	UserId  uint64

	CreateAt time.Time
	UpdateAt time.Time
	DeleteAt time.Time
	CreateBy uint64
	UpdateBy uint64
	DeleteBy uint64
}

func (PostEntity) TableName() string {
	return "posts"
}
