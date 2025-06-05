package main

import (
	postEntities "task4/post/core/entities"
	"task4/shared/httpapi/auth"
	applogger "task4/shared/kernel/logger"
	userEntities "task4/user/core/entities"
	"task4/user/httpapi/controllers"
	"task4/user/infrastructure/data"

	"github.com/gin-gonic/gin"
)

func main() {

	//初始化日志
	log, apperr := applogger.NewLogger()
	if apperr != nil {
		panic(apperr)
	}
	defer log.Sync()

	// 初始化数据库
	dbContext, err := data.InitDb()
	if err != nil {
		panic(err)
	}

	db := dbContext.GetDb()
	db.AutoMigrate(&userEntities.UserEntity{}, &postEntities.PostEntity{}, &postEntities.CommentEntity{})
	db.Exec("ALTER TABLE users COMMENT '用户表'")
	db.Exec("ALTER TABLE posts COMMENT '文章表'")
	db.Exec("ALTER TABLE comments COMMENT '评论表'")

	r := gin.Default()

	// 记录请求日志
	r.Use(gin.LoggerWithWriter(gin.DefaultWriter))

	// v1 路由分组
	v1 := r.Group("/v1")
	v1.Use(auth.JwtAuthMiddleware())
	// 初始化控制器
	userctrl := controllers.NewUserController(dbContext)
	// 绑定路由
	v1.GET("/auth", userctrl.Test)

	// 默认路由分组
	r.POST("/auth/register", userctrl.Register)
	r.POST("/auth/login", userctrl.Login)

	r.Run(":8080")
}
