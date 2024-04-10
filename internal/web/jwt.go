package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type JwtHandler struct {
	// access_token key
	atKey []byte
	// refresh_token key
	rtKey []byte
}

func NewJwtHandler() JwtHandler {
	return JwtHandler{
		atKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igyf0"),
		rtKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igyf0"),
	}
}

func (j JwtHandler) setJwtToken(ctx *gin.Context, uid int64) error {
	claims := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		UserId: uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenStr, err := token.SignedString([]byte(j.atKey))
	if err != nil {
		return err
	}

	// 返回token
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (j JwtHandler) setRefreshToken(ctx *gin.Context, uid int64) error {
	claims := &RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		UserId: uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenStr, err := token.SignedString([]byte(j.rtKey))
	if err != nil {
		return err
	}

	// 返回token
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func ExtractToken(ctx *gin.Context) string {
	// 用jwt校验
	tokenHeader := ctx.GetHeader("Authorization")

	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	UserId int64
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId int64
}
