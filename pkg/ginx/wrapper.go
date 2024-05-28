package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"webbook/pkg/logger"
)

// 受制于泛型，我们这里只能使用包变量，我深恶痛绝的包变量
var log logger.Logger = logger.NewNoOpLogger()

func SetLogger(l logger.Logger) {
	log = l
}

// WrapToken
func WrapToken[c jwt.Claims](fn func(ctx *gin.Context, uc c) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 可以用包变量来配置，还是那句话，因为泛型的限制，这里只能用包变量
		rawVal, ok := ctx.Get("user")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			log.Error("无法获得 claims",
				logger.String("path", ctx.Request.URL.Path))
			return
		}
		// 注意，这里要求放进去 ctx 的不能是*UserClaims，这是常见的一个错误
		claims, ok := rawVal.(c)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			log.Error("无法获得 claims",
				logger.String("path", ctx.Request.URL.Path))
			return
		}
		res, err := fn(ctx, claims)
		if err != nil {
			log.Error("执行业务逻辑失败",
				logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

// WrapBodyAndToken 如果做成中间件来源出去，那么直接耦合 UserClaims 也是不好的。
func WrapBodyAndToken[T any, c jwt.Claims](fn func(ctx *gin.Context, req T, uc c) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			log.Error("解析请求失败", logger.Error(err))
			return
		}
		// 可以用包变量来配置，还是那句话，因为泛型的限制，这里只能用包变量
		rawVal, ok := ctx.Get("user")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			log.Error("无法获得 claims",
				logger.String("path", ctx.Request.URL.Path))
			return
		}
		// 注意，这里要求放进去 ctx 的不能是*UserClaims，这是常见的一个错误
		claims, ok := rawVal.(c)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			log.Error("无法获得 claims",
				logger.String("path", ctx.Request.URL.Path))
			return
		}
		res, err := fn(ctx, req, claims)
		if err != nil {
			log.Error("执行业务逻辑失败",
				logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

// WrapBody 如果做成中间件来源出去，那么直接耦合 UserClaims 也是不好的。
func WrapBody[T any](fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			log.Error("解析请求失败", logger.Error(err))
			return
		}

		res, err := fn(ctx, req)
		if err != nil {
			log.Error("执行业务逻辑失败",
				logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}
