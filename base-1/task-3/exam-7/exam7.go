package exam7

import (
	exam5 "task3/exam-5"

	"gorm.io/gorm"
)

/*
题目七：
钩子函数,
继续使用博客系统的模型。
要求 ：
为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。,
为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
*/

func (post *PostEntity) BeforeCreate(tx *gorm.DB) (err error) {
	if err = Db.Model(&exam5.UserEntity{}).Where("id = ?", post.UserID).Update("post_count", gorm.Expr("post_count + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

var Db = exam5.Db

type PostEntity struct {
	ID      uint   `gorm:"primaryKey;autoIncrement;column:id"`
	Title   string `gorm:"column:title"`
	Content string `gorm:"column:content"`
	UserID  uint   `gorm:"column:user_id"`

	Comments []exam5.CommentEntity `gorm:"foreignKey:PostID"`
}

type UserEntity struct {
	ID        uint   `gorm:"primaryKey;autoIncrement;column:id"`
	Username  string `gorm:"column:username"`
	Password  string `gorm:"column:password"`
	PostCount uint   `gorm:"column:post_count"` // 用户的文章数量统计字段
	// Posts     []PostEntity `gorm:"foreignKey:UserID"` // 与 PostEntity 的关联关系
}
