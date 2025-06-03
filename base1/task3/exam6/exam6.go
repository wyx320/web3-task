package exam6

import exam5 "task3/exam-5"

/*
题目六：
关联查询。
基于上述博客系统的模型定义。
要求 ：
编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
编写Go代码，使用Gorm查询评论数量最多的文章信息。
*/

var Db = exam5.Db

func GetPostWithCommentByUser(userId uint) []exam5.PostEntity {
	var posts []exam5.PostEntity
	Db.Preload("Comments").Where("user_id = ?", userId).Find(&posts)
	return posts
}

func GetPostWithMaxComment() (*exam5.PostEntity, error) {
	var post exam5.PostEntity
	sql := "SELECT * FROM posts WHERE id = (SELECT post_id FROM comments GROUP BY post_id ORDER BY COUNT(*) DESC LIMIT 1)"
	err := Db.Raw(sql).Scan(&post).Error
	return &post, err
}
