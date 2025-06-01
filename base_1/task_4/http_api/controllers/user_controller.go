package user_controller

import (
	"net/http"
	md5_helper "task4/Infrastructure/Tools"
	"task4/core/entities"
	user_register_dto "task4/http_api/models/user"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserController struct{}

func (u *UserController) Register(c *gin.Context) {
	var dto user_register_dto.UserRegisterDto
	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var user entities.UserEntity
	copier.Copy(&user, &dto)

	var err error
	user.Salt, err = md5_helper.Md5EncodingOnly(user.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误",
		})
		return
	}

	user.Password, err = md5_helper.Md5EncodingWithSalt(user.Password, user.Salt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误",
		})
		return
	}

	// 该去生成一个Db 并在main的init做CodeFirst 再回来这里做CRUD操作
}
