package entities

import "time"

// 评论实体
type CommentEntity struct {
	Id      uint64 `gorm:"primary_key;autoIncrement"`
	Content string
	UserId  uint64
	PostId  uint64

	CreateAt time.Time
	UpdateAt time.Time
	DeleteAt time.Time
	CreateBy uint64
	UpdateBy uint64
	DeleteBy uint64
}

func (CommentEntity) TableName() string {
	return "comments"
}
