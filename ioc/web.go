package ioc

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"webbook/internal/web"
	ijwt "webbook/internal/web/jwt"
	"webbook/internal/web/middleware"
	logger2 "webbook/pkg/logger"
	"webbook/pkg/middlewares/logger"
	mdwratelimit "webbook/pkg/middlewares/ratelimit"
	"webbook/pkg/ratelimit"
)

func InitWebServer(mdls []gin.HandlerFunc, hdl *web.UserHandler, oauth2WechatHdl *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server)
	oauth2WechatHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable, jwtHdl ijwt.Handler, l logger2.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		middleware.NewLoginJwtMiddlewareBuilder(jwtHdl).
			IgnorePaths("/users/login").
			IgnorePaths("/users/signupt").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("/oauth2/wechat/callback").
			Build(),
		mdwratelimit.NewRedisSlidingWindowLimiter(ratelimit.NewRedisSlidingWindowLimiter(redisClient, time.Second, 100)).Build(),
		logger.NewMiddlewareBuilder(func(ctx context.Context, al *logger.AccessLog) {
			l.Debug("Http请求", logger2.Field{
				Key:   "al",
				Value: al,
			})
		}).AllowReqBody().AllowRespBody().Builder(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		// 是否允许携带cookie
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "你的域名.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
