package ioc

import (
	"go.uber.org/zap"
	"webbook/pkg/logger"
)

func InitLogger() logger.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return logger.NewZapLogger(l)
}
