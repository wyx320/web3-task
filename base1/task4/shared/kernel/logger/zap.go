package applogger

import (
	appresult "task4/shared/kernel/result"

	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, *appresult.AppError) {
	config := zap.NewProductionConfig()

	// 设置日志级别
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	// 设置输出格式
	config.Encoding = "json"

	// 设置输出路径
	config.OutputPaths = []string{"stdout", "logs/app.log"}
	config.ErrorOutputPaths = []string{"stderr", "logs/app.log"}

	logger, err := config.Build()
	if err != nil {
		return nil, appresult.InternalServerError.WriteDetail("failed to create logger")
	}
	return logger, nil
}
