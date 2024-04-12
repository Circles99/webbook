//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webbook/internal/repository"
	"webbook/internal/repository/cache"
	"webbook/internal/repository/dao"
	"webbook/internal/service"
	"webbook/internal/web"
	"webbook/internal/web/jwt"
	"webbook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方服务
		ioc.InitDB, ioc.InitRedis, ioc.InitSms,
		// dao, cache
		dao.NewUserDao, cache.NewUserCache, cache.NewCodeCache,
		// repository
		repository.NewUserRepository, repository.NewCodeRepository,
		// service.go
		service.NewUserService, service.NewCodeService,
		// web
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,

		// jwt
		jwt.NewRedisJwtHandler,

		// init
		ioc.InitWebServer,
		ioc.InitMiddlewares,
		ioc.InitOAuth2WechatService,
		ioc.NewWechatHandlerConfig,
	)
	return new(gin.Engine)
}
