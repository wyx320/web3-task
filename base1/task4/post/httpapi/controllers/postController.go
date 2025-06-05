package controllers

import (
	"task4/post/core/entities"
	models "task4/post/httpapi/models/post"
	appresult "task4/shared/kernel/result"
	"task4/user/infrastructure/data"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type PostController struct {
	data data.Database

	logger *zap.Logger
}

// 创建文章
func (p *PostController) Create(c *gin.Context) {
	var dto models.PostForCreateDto
	if err := c.ShouldBind(&dto); err != nil {
		appresult.ErrorResponse(c, p.logger, appresult.BadRequestError)
		return
	}

	var post entities.PostEntity
	copier.Copy(&post, &dto)
	post.CreateAt = time.Now()

	userIdValue, exist := c.Get("user_id")
	if !exist {
		appresult.ErrorResponse(c, p.logger, appresult.BadRequestError.WriteDetail("user_id not found in token"))
		return
	}
	userId := userIdValue.(uint64)

	post.UserId = userId
	post.CreateBy = userId

	db := p.data.GetDb()
	db.Create(&post)

	appresult.SuccessResponse(c, p.logger, &post)
}

// 更新文章
func (p *PostController) Update(c *gin.Context) {
	var dto models.PostForUpdateDto
	if err := c.ShouldBind(&dto); err != nil {
		appresult.ErrorResponse(c, p.logger, appresult.BadRequestError)
		return
	}

	db := p.data.GetDb()

	var post entities.PostEntity
	db.First(&post, dto.Id)
	if post.Id == 0 {
		appresult.ErrorResponse(c, p.logger, appresult.ArticleNotFoundError)
		return
	}

	userIdValue, exists := c.Get("user_id")
	if !exists {
		appresult.ErrorResponse(c, p.logger, appresult.BadRequestError.WriteDetail("user_id not found in token"))
		return
	}
	userId := userIdValue.(uint64)

	if post.UserId != userId {
		appresult.ErrorResponse(c, p.logger, appresult.UnauthorizedError.WriteDetail("只有文章的作者才能更新自己的文章"))
		return
	}

	copier.Copy(&post, &dto)

	post.UpdateBy = userId
	post.UpdateAt = time.Now()

	db.Model(&post).Updates(&post)

	appresult.SuccessResponse(c, p.logger, &post)
}

// 获取文章明细
func (p *PostController) Get(c *gin.Context) {
	id := c.Param("id")
	db := p.data.GetDb()
	var post entities.PostEntity
	db.First(&post, id)
	if post.Id == 0 {
		appresult.ErrorResponse(c, p.logger, appresult.ArticleNotFoundError)
		return
	}
	var dto models.PostDto
	copier.Copy(&dto, &post)

	appresult.SuccessResponse(c, p.logger, &dto)
}

// 获取文章列表
func (p *PostController) GetList(c *gin.Context) {
	db := p.data.GetDb()
	var posts []entities.PostEntity
	db.Find(&posts)
	var dtos []models.PostDto
	copier.Copy(&dtos, &posts)

	appresult.SuccessResponse(c, p.logger, &dtos)
}

// 删除文章
func (p *PostController) Delete(c *gin.Context) {
	id := c.Param("id")
	db := p.data.GetDb()
	var post entities.PostEntity
	db.First(&post, id)
	if post.Id == 0 {
		appresult.ErrorResponse(c, p.logger, appresult.ArticleNotFoundError)
		return
	}

	userIdValue, exists := c.Get("user_id")
	if !exists {
		appresult.ErrorResponse(c, p.logger, appresult.BadRequestError.WriteDetail("user_id not found in token"))
		return
	}
	userId := userIdValue.(uint64)

	if post.UserId != userId {
		appresult.ErrorResponse(c, p.logger, appresult.UnauthorizedError.WriteDetail("只有文章的作者才能删除自己的文章"))
		return
	}

	post.IsDeleted = true
	db.Model(&post).Updates(&post)
}
