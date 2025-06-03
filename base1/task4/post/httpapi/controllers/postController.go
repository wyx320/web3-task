package controllers

import (
	"net/http"
	"task4/post/core/entities"
	"task4/post/httpapi/models"
	"task4/user/infrastructure/data"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type PostController struct {
	data data.Database
}

// 创建文章
func (p *PostController) Create(c *gin.Context) {
	var dto models.PostForCreateDto
	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(400, gin.H{"message": "参数错误"})
		return
	}

	var post entities.PostEntity
	copier.Copy(&post, &dto)
	post.CreateAt = time.Now()

	userIdValue, exist := c.Get("user_id")
	if !exist {
		c.JSON(400, gin.H{"message": "cannot get user_id from token"})
		return
	}
	userId := userIdValue.(uint64)

	post.UserId = userId
	post.CreateBy = userId

	db := p.data.GetDb()
	db.Create(&post)

	c.JSON(http.StatusOK, &post)
}

// 更新文章
func (p *PostController) Update(c *gin.Context) {
	var dto models.PostForUpdateDto
	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(400, gin.H{"message": "参数错误"})
		return
	}

	db := p.data.GetDb()

	var post entities.PostEntity
	db.First(&post, dto.Id)
	if post.Id == 0 {
		c.JSON(400, gin.H{"message": "文章不存在"})
		return
	}

	copier.Copy(&post, &dto)

	userIdValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"message": "cannot get user_id from token"})
		return
	}
	userId := userIdValue.(uint64)

	post.UpdateBy = userId
	post.UpdateAt = time.Now()

	db.Model(&post).Updates(&post)

	c.JSON(http.StatusOK, &post)

	// 准备参照net项目做异常统一处理 然后修改为只有文章的作者才能操作 编写删除接口也是只有文章作者才能操作

}

// 获取文章明细
func (p *PostController) Get(c *gin.Context) {
	id := c.Param("id")
	db := p.data.GetDb()
	var post entities.PostEntity
	db.First(&post, id)
	if post.Id == 0 {
		c.JSON(400, gin.H{"message": "文章不存在"})
		return
	}
	var dto models.PostDto
	copier.Copy(&dto, &post)
	c.JSON(http.StatusOK, &dto)
}

// 获取文章列表
func (p *PostController) GetList(c *gin.Context) {
	db := p.data.GetDb()
	var posts []entities.PostEntity
	db.Find(&posts)
	var dtos []models.PostDto
	copier.Copy(&dtos, &posts)
	c.JSON(http.StatusOK, &dtos)
}
