package main

import (
	postEntities "task4/post/core/entities"
	postControllers "task4/post/httpapi/controllers"
	"task4/shared/httpapi/auth"
	applogger "task4/shared/kernel/logger"
	userEntities "task4/user/core/entities"
	userControllers "task4/user/httpapi/controllers"
	"task4/user/infrastructure/data"

	"github.com/gin-gonic/gin"
)

func main() {

	//初始化日志
	logger, apperr := applogger.NewLogger()
	if apperr != nil {
		panic(apperr)
	}
	defer logger.Sync()

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
	userCtl := userControllers.NewUserController(dbContext, logger)
	postCtl := postControllers.NewPostController(dbContext, logger)
	// v1 绑定路由
	v1.GET("/auth", userCtl.Test)

	v1.GET("/post", postCtl.GetList)
	v1.GET("/post/:id", postCtl.Get)
	v1.POST("/post", postCtl.Create)
	v1.PUT("/post/:id", postCtl.Update)
	v1.DELETE("/post/:id", postCtl.Delete)

	v1.GET("/comment", postCtl.GetList)
	v1.POST("/comment", postCtl.Create)

	// 默认路由
	r.POST("/auth/register", userCtl.Register)
	r.POST("/auth/login", userCtl.Login)

	r.Run(":8080")
}
