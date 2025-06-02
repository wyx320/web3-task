package main

import (
	"task4/user/core/entities"
	"task4/user/httpapi/controllers"
	"task4/user/infrastructure/data"

	"github.com/gin-gonic/gin"
)

func main() {

	// 初始化数据库
	dbContext, err := data.InitDb()
	if err != nil {
		panic(err)
	}

	db := dbContext.GetDb()
	db.AutoMigrate(&entities.UserEntity{}, &entities.PostEntity{}, &entities.CommentEntity{})
	db.Exec("ALTER TABLE users COMMENT '用户表'")
	db.Exec("ALTER TABLE posts COMMENT '文章表'")
	db.Exec("ALTER TABLE comments COMMENT '评论表'")

	r := gin.Default()

	// 初始化控制器
	userctrl := controllers.NewUserController(dbContext)
	// 绑定路由
	r.POST("/user", userctrl.Register)
	// r.GET("/user", userctrl.Test)

	r.Run(":8080")
}
