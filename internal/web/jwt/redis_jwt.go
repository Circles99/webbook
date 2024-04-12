package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

var (
	// access_token key
	AtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igyf0")

	// refresh_token key
	RtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igyf0")
)

type RedisJwt struct {
	cmd redis.Cmdable
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	UserId int64
	Ssid   string
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId    int64
	Ssid      string
	UserAgent string
}

func NewRedisJwtHandler(cmd redis.Cmdable) Handler {
	return &RedisJwt{
		cmd: cmd,
	}
}

func (h RedisJwt) ExtractToken(ctx *gin.Context) string {
	// 用jwt校验
	tokenHeader := ctx.GetHeader("Authorization")

	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

func (h RedisJwt) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJwtToken(ctx, uid, ssid)
	if err != nil {
		return err
	}

	err = h.SetRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	return nil
}

func (h RedisJwt) ClearToken(ctx *gin.Context) error {
	c, _ := ctx.Get("claims")

	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")

	}

	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")

	return h.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid), "", time.Hour*24*7).Err()
}

func (h RedisJwt) CheckSession(ctx *gin.Context, ssid string) error {
	_, err := h.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	return err
}

func (h RedisJwt) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	//TODO implement me
	panic("implement me")
}

func (h RedisJwt) SetJwtToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		UserId: uid,
		Ssid:   ssid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenStr, err := token.SignedString([]byte(AtKey))
	if err != nil {
		return err
	}

	// 返回token
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}
