package controllers

import (
	"task4/post/core/entities"
	models "task4/post/httpapi/models/comment"
	appresult "task4/shared/kernel/result"
	"task4/user/infrastructure/data"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommentController struct {
	data data.Database

	logger *zap.Logger
}

// 创建评论
func (ctrl *CommentController) Create(c *gin.Context) {
	var dto models.CommentCreateDto
	c.ShouldBind(&dto)
	// 获取当前用户ID
	userIdValue, exist := c.Get("user_id")
	if !exist {
		appresult.ErrorResponse(c, ctrl.logger, appresult.BadRequestError.WriteDetail("user_id not found in token"))
		return
	}
	userId := userIdValue.(uint64)

	comment := &entities.CommentEntity{
		Content: dto.Content,
		PostId:  dto.PostId,
		UserId:  userId,

		CreateAt: time.Now(),
		CreateBy: userId,
	}

	db := ctrl.data.GetDb()
	db.Create(comment)
}

func (ctrl *CommentController) GetList(c *gin.Context) {
	db := ctrl.data.GetDb()
	var comments []entities.CommentEntity
	db.Where("post_id = ? and is_deleted = 0").Find(&comments)

	// 这里没写完 需要从获取到请求参数 这个接口写完后 进行功能测试就好了。。
}
