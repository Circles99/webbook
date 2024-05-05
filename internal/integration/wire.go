//go:build wireinject

package integration

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
		ioc.InitDB, ioc.InitRedis, ioc.InitSms, ioc.InitLogger,
		// dao, cache
		dao.NewUserDao, cache.NewUserCache, cache.NewCodeCache, dao.NewArticleDao,
		// repository
		repository.NewUserRepository, repository.NewCodeRepository, repository.NewArticleRepository,
		// service.go
		service.NewUserService, service.NewCodeService, service.NewArticleService,
		// web
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,

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

var thirdProvider = wire.NewSet(ioc.InitDB, ioc.InitRedis, ioc.InitLogger)

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(thirdProvider, service.NewArticleService, web.NewArticleHandler, repository.NewArticleRepository, dao.NewArticleDao)
	return &web.ArticleHandler{}
}
