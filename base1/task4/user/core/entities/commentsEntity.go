package entities

// 评论实体
type CommentEntity struct {
	Id      uint64 `gorm:"primary_key;autoIncrement"`
	Content string
	UserId  uint64
	PostId  uint64

	CreateAt uint64
}

func (CommentEntity) TableName() string {
	return "comments"
}
