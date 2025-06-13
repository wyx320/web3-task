package controllers

import (
	"math"
	"task4/post/core/entities"
	models "task4/post/httpapi/models/comment"
	"task4/shared/kernel/page"
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

type GetListQuery struct {
	PostId uint64 `form:"post_id"`
}

func (ctrl *CommentController) GetList(c *gin.Context) {

	// 获取分页参数
	var query GetListQuery
	c.BindQuery(&query)

	// 获取评论列表
	db := ctrl.data.GetDb()
	var comments []entities.CommentEntity
	db.Where("post_id = ? and is_deleted = 0", query.PostId).Find(&comments)

	// 分页处理
	pageObject := page.NewPageObject(math.MaxInt64, 1)
	pageList := page.NewPageList(comments, *pageObject)

	// 返回结果
	appresult.SuccessResponse(c, ctrl.logger, pageList)
}
