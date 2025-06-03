package controllers

import (
	"fmt"
	"net/http"
	"time"

	"task4/user/core/entities"
	"task4/user/httpapi/models"
	"task4/user/infrastructure/data"
	"task4/user/infrastructure/tools"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserController struct {
	// 依赖接口
	dataDb data.Database
}

// 显示初始化
func NewUserController(db data.Database) *UserController {
	return &UserController{dataDb: db}
}

func (u *UserController) Register(c *gin.Context) {
	var dto models.AuthForRegisterDto
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
	user.Salt, err = tools.Md5EncodingOnly(user.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误",
		})
		return
	}

	user.Password, err = tools.Md5EncodingWithSalt(user.Password, user.Salt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误",
		})
		return
	}

	// 截止到这步 该去依赖注入生成一个Db 并在main的init做CodeFirst 再回来这里做下面的CRUD操作

	if u.dataDb == nil || u.dataDb.GetDb() == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "数据库连接失败",
		})
		return
	}

	err = u.dataDb.GetDb().Create(&user).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误",
		})
		return
	}
}

func (u *UserController) Test(c *gin.Context) {
	var users []entities.UserEntity
	db := u.dataDb.GetDb()
	fmt.Println(db)
	db.Model(&entities.UserEntity{}).Find(&users)
	fmt.Printf("%v", users)
	c.String(http.StatusOK, "Hello World")
}

func (u *UserController) Login(c *gin.Context) {
	var dto models.AuthForLoginDto
	c.ShouldBind(&dto)

	salt, err := tools.Md5EncodingOnly(dto.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误",
		})
		return
	}

	if u.dataDb == nil || u.dataDb.GetDb() == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "数据库连接失败",
		})
		return
	}

	var user entities.UserEntity
	if err = u.dataDb.GetDb().Where("username=?", dto.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误",
		})
		return
	}

	if user.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "用户名或密码错误",
		})
		return
	}

	password := dto.Password
	if password, err = tools.Md5EncodingWithSalt(password, salt); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "系统错误",
		})
	}

	if password != user.Password {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "用户名或密码错误",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       user.Id,
		"usewrname": user.Username,
		"exp":       time.Now().Add(time.Hour * 8).Unix(),
	})
	var tokenString string
	if tokenString, err = token.SignedString([]byte("jx9KdLm3sQwRtP5BvG8XzY2NcF6EaH0")); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "JWT签名失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        200,
		"msg":         "登录成功",
		"accessToken": tokenString,
	})
}
