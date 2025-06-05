package appresult

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AppError struct {
	Code    int
	Message string
	Detail  string

	zapLogger *zap.Logger
}

// New 创建新的错误
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// WriteDetail 添加错误详情
func (ar *AppError) WriteDetail(detail string) *AppError {
	return &AppError{
		Code:    ar.Code,
		Message: ar.Message,
		Detail:  detail,
	}
}

// 预定义错误
var (
	InternalServerError = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "服务器内部错误",
	}

	BadRequestError = &AppError{
		Code:    http.StatusBadRequest,
		Message: "请求参数错误",
	}

	UnauthorizedError = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "未授权访问",
	}

	ForbiddenError = &AppError{
		Code:    http.StatusForbidden,
		Message: "禁止访问",
	}

	NotFoundError = &AppError{
		Code:    http.StatusNotFound,
		Message: "资源不存在",
	}

	DatabaseConnectionError = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "数据库连接错误",
	}

	UserNotFoundError = &AppError{
		Code:    http.StatusNotFound,
		Message: "用户不存在",
	}

	ArticleNotFoundError = &AppError{
		Code:    http.StatusNotFound,
		Message: "文章不存在",
	}

	CommentNotFoundError = &AppError{
		Code:    http.StatusNotFound,
		Message: "评论不存在",
	}

	AuthenticationFailedError = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "认证失败",
	}

	InvalidTokenError = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "无效的令牌",
	}
)

// 错误响应函数
func ErrorResponse(c *gin.Context, logger *zap.Logger, appError *AppError) {
	c.JSON(appError.Code, gin.H{
		"code":    appError.Code,
		"message": appError.Message,
	})
	logger.Error("请求失败", zap.Int("code", appError.Code), zap.String("message", appError.Message))
}

// 成功响应函数
func SuccessResponse(c *gin.Context, logger *zap.Logger, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    data,
	})
	logger.Info("请求成功", zap.Any("data", data))
}
