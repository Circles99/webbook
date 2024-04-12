package jwt

import (
	"github.com/gin-gonic/gin"
)

type Handler interface {
	ExtractToken(ctx *gin.Context) string
	SetJwtToken(ctx *gin.Context, uid int64, ssid string) error
	SetLoginToken(ctx *gin.Context, uid int64) error
	ClearToken(ctx *gin.Context) error
	CheckSession(ctx *gin.Context, ssid string) error

	SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error
}
