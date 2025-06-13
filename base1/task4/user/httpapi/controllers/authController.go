package controllers

import (
	"fmt"
	"net/http"
	"time"

	"task4/shared/config"
	appresult "task4/shared/kernel/result"
	"task4/user/core/entities"
	"task4/user/httpapi/models"
	"task4/user/infrastructure/data"
	"task4/user/infrastructure/tools"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type UserController struct {
	// 依赖接口
	dataDb data.Database

	logger *zap.Logger
}

// 显示初始化
func NewUserController(db data.Database, logger *zap.Logger) *UserController {
	return &UserController{
		dataDb: db,
		logger: logger,
	}
}

func (u *UserController) Register(c *gin.Context) {
	var dto models.AuthForRegisterDto
	if err := c.ShouldBind(&dto); err != nil {
		appresult.ErrorResponse(c, u.logger, appresult.BadRequestError)
		return
	}

	var user entities.UserEntity
	copier.Copy(&user, &dto)

	var err error
	user.Salt, err = tools.Md5EncodingOnly(user.Username)
	if err != nil {
		appresult.ErrorResponse(c, u.logger, appresult.InternalServerError)
		return
	}

	user.Password, err = tools.Md5EncodingWithSalt(user.Password, user.Salt)
	if err != nil {
		appresult.ErrorResponse(c, u.logger, appresult.InternalServerError)
		return
	}

	// 截止到这步 该去依赖注入生成一个Db 并在main的init做CodeFirst 再回来这里做下面的CRUD操作

	if u.dataDb == nil || u.dataDb.GetDb() == nil {
		appresult.ErrorResponse(c, u.logger, appresult.InternalServerError.WriteDetail("数据库连接失败"))
		return
	}

	err = u.dataDb.GetDb().Create(&user).Error
	if err != nil {
		appresult.ErrorResponse(c, u.logger, appresult.InternalServerError)
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
	c.BindJSON(&dto)

	salt, err := tools.Md5EncodingOnly(dto.Username)
	if err != nil {
		appresult.ErrorResponse(c, u.logger, appresult.InternalServerError)
		return
	}

	if u.dataDb == nil || u.dataDb.GetDb() == nil {
		appresult.ErrorResponse(c, u.logger, appresult.InternalServerError.WriteDetail("数据库连接失败"))
		return
	}

	var user entities.UserEntity
	if err = u.dataDb.GetDb().Where("username=?", dto.Username).First(&user).Error; err != nil {
		appresult.ErrorResponse(c, u.logger, appresult.InternalServerError)
		return
	}

	if user.Id == 0 {
		appresult.ErrorResponse(c, u.logger, appresult.BadRequestError.WriteDetail("用户名或密码错误"))
		return
	}

	password := dto.Password
	if password, err = tools.Md5EncodingWithSalt(password, salt); err != nil {
		appresult.ErrorResponse(c, u.logger, appresult.InternalServerError)
		return
	}

	if password != user.Password {
		appresult.ErrorResponse(c, u.logger, appresult.BadRequestError.WriteDetail("用户名或密码错误"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       user.Id,
		"usewrname": user.Username,
		"exp":       time.Now().Add(time.Hour * 8).Unix(),
	})
	var tokenString string
	if tokenString, err = token.SignedString([]byte(config.SecretKey)); err != nil {
		appresult.ErrorResponse(c, u.logger, appresult.InternalServerError.WriteDetail("JWT签名失败"))
		return
	}

	appresult.SuccessResponse(c, u.logger, map[string]string{"accessToken": tokenString})
}
