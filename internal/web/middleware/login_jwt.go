package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	ijwt "webbook/internal/web/jwt"
)

type LoginJwtMiddlewareBuilder struct {
	paths []string
	ijwt.Handler
}

func NewLoginJwtMiddlewareBuilder(jwthdl ijwt.Handler) *LoginJwtMiddlewareBuilder {
	return &LoginJwtMiddlewareBuilder{}
}

func (l *LoginJwtMiddlewareBuilder) IgnorePaths(path string) *LoginJwtMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJwtMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 用jwt校验
		tokenStr := l.ExtractToken(ctx)

		claims := &ijwt.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return ijwt.AtKey, nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.UserId == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserAgent != ctx.Request.UserAgent() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = l.CheckSession(ctx, claims.Ssid)
		if err != nil {
			// 要么redis有问题，要么已经退出登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
		//
		//tokenStr, err = token.SignedString([]byte(web.SaltKey))
		//if err != nil {
		//	log.Println("JWT 续约失败", err)
		//}
		//ctx.Header("x-jwt-token", tokenStr)
		ctx.Set("claims", claims)

	}
}
