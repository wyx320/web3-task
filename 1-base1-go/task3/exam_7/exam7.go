package exam7

import (
	exam5 "task3/exam_5"

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

// 创建文章
func CreatePost(post *PostEntity) error {
	return Db.Create(post).Error
}

// 添加评论
func AddComment(comment *CommentEntity) error {
	return Db.Create(comment).Error
}

// 删除评论
func DeleteComment(comment *CommentEntity) error {
	return Db.Delete(comment).Error
}

// 钩子函数：文章创建时自动更新用户的文章数量统计字段
func (post *PostEntity) BeforeCreate(tx *gorm.DB) (err error) {
	if err = tx.Model(&exam5.UserEntity{}).Where("id = ?", post.UserId).Update("post_count", gorm.Expr("post_count + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// 钩子函数：评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
func (comment *CommentEntity) AfterDelete(tx *gorm.DB) (err error) {
	var commentCount int64
	err = tx.Model(&CommentEntity{}).Where("post_id = ?", comment.PostId).Count(&commentCount).Error

	if err != nil {
		return err
	}

	if commentCount == 0 {
		err = tx.Model(&PostEntity{}).Where("id = ?", comment.PostId).Update("comment_status", "无评论").Error
		if err != nil {
			return err
		}
	}

	return nil
}

var Db = exam5.Db

type PostEntity struct {
	Id            uint   `gorm:"primaryKey;autoIncrement;column:id"`
	Title         string `gorm:"column:title"`
	Content       string `gorm:"column:content"`
	UserId        uint   `gorm:"column:user_id"`
	CommentStatus string

	Comments []exam5.CommentEntity `gorm:"foreignKey:PostId"`
}

type UserEntity struct {
	Id        uint   `gorm:"primaryKey;autoIncrement;column:id"`
	Username  string `gorm:"column:username"`
	Password  string `gorm:"column:password"`
	PostCount uint   `gorm:"column:post_count"` // 用户的文章数量统计字段
}

type CommentEntity struct {
	Id      uint   `gorm:"primaryKey;autoIncrement;column:id"`
	Content string `gorm:"column:content"`
	PostId  uint   `gorm:"column:post_id"`
	UserId  uint   `gorm:"column:user_id"`
}

func (UserEntity) TableName() string {
	return "users"
}

func (PostEntity) TableName() string {
	return "posts"
}

func (CommentEntity) TableName() string {
	return "comments"
}
