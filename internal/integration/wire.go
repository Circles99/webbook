//go:build wireinject

package integration

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webbook/internal/repository"
	"webbook/internal/repository/article"
	"webbook/internal/repository/cache"
	"webbook/internal/repository/dao"
	articleDao "webbook/internal/repository/dao/article"
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
		dao.NewUserDao, cache.NewUserCache, cache.NewCodeCache, articleDao.NewArticleDao,
		// repository
		repository.NewUserRepository, repository.NewCodeRepository, article.NewArticleRepository,
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

var thirdProvider = wire.NewSet(ioc.InitDB, ioc.InitRedis, ioc.InitLogger, InitMongoDB)

func InitArticleHandler(dao article2.ArticleDAO) *web.ArticleHandler {
	wire.Build(thirdProvider, service.NewArticleService, web.NewArticleHandler, article.NewArticleRepository)
	return &web.ArticleHandler{}
}

//func InitArticleHandler(dao article.ArticleDAO) *web.ArticleHandler {
//	wire.Build(thirdProvider,
//		userSvcProvider,
//		interactiveSvcProvider,
//		cache.NewRedisArticleCache,
//		//wire.InterfaceValue(new(article.ArticleDAO), dao),
//		repository.NewArticleRepository,
//		service.NewArticleService,
//		web.NewArticleHandler)
//	return new(web.ArticleHandler)
//}
