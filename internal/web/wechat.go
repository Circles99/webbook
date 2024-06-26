package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"time"
	"webbook/internal/service"
	"webbook/internal/service/oauth2/wechat"
	ijwt "webbook/internal/web/jwt"
	"webbook/pkg/ginx"
)

type OAuth2WechatHandler struct {
	svc      wechat.Service
	userSvc  service.UserService
	stateKey []byte
	cfg      WechatConfig
	ijwt.Handler
}

type WechatConfig struct {
	Secure bool
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService, cfg WechatConfig, jwtHdl ijwt.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:      svc,
		userSvc:  userSvc,
		stateKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igyf0"),
		cfg:      cfg,
		Handler:  jwtHdl,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()

	url, err := h.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "构造登录URL失败",
		})
		return
	}

	if err = h.setStateCookie(ctx, state); err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}

	ctx.JSON(http.StatusOK, ginx.Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 3)),
		},
	})

	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		return err
	}

	ctx.SetCookie("jwt-state", tokenStr, 600, "oauth2/wechat/callback", "", h.cfg.Secure, true)
	return nil
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "登录失败",
		})
		return
	}

	info, err := h.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 从userService 里面拿userId
	err = h.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "OK",
	})
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")

	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("拿不到state的cookie %w", err)
	}

	var sc StateClaims

	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {

		return fmt.Errorf("token已经过期 %w", err)
	}

	if sc.State != state {
		return errors.New("state 不相等")
	}
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}
