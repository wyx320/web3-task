package exam5

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
题目五：
模型定义,
假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
要求 ：
使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
编写Go代码，使用Gorm创建这些模型对应的数据库表。
*/

func Test() {
	fmt.Println("CodeFirst执行成功！")
}

type UserEntity struct {
	Id       uint   `gorm:"primaryKey;autoIncrement;column:id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

type PostEntity struct {
	Id      uint   `gorm:"primaryKey;autoIncrement;column:id"`
	Title   string `gorm:"column:title"`
	Content string `gorm:"column:content"`
	UserId  uint   `gorm:"column:user_id"`

	Comments []CommentEntity `gorm:"foreignKey:PostId"`
}
type CommentEntity struct {
	Id      uint   `gorm:"primaryKey;autoIncrement;column:id"`
	Content string `gorm:"column:content"`
	PostId  uint   `gorm:"column:post_id"`
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

var Db *gorm.DB

func init() {
	host := "localhost"
	port := 3306
	user := "root"
	password := "1"
	database := "web3_task3"
	timeout := "10s"

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", user, password, host, port, database, timeout)

	var err error
	Db, err = gorm.Open(mysql.Open(connStr))
	if err != nil {
		fmt.Println("数据库连接失败：", err)
		return
	}

	// 连接成功
	fmt.Println(Db)

	// // 自动迁移表结构
	// Db.AutoMigrate(&UserEntity{}, &PostEntity{}, &CommentEntity{})
	// Db.Exec("ALTER TABLE users COMMENT '用户表'")
	// Db.Exec("ALTER TABLE posts COMMENT '文章表'")
	// Db.Exec("ALTER TABLE comments COMMENT '评论表'")
}
