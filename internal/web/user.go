package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"webbook/internal/domain"
	"webbook/internal/service"
)

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func (u UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)

}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = "^(?=.*[A-Za-z])(?=.*\\d)(?=.*[$@$!%*#?&])[A-Za-z\\d$@$!%*#?&]{8,}$"
	)

	emailExp := regexp.MustCompile(emailRegexPattern)
	passwordExp := regexp.MustCompile(passwordRegexPattern)

	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq

	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidEmailOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
	return

}

func (u *UserHandler) Signup(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq

	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok := u.emailExp.MatchString(req.Email)
	if !ok {
		ctx.String(http.StatusOK, "邮箱错误")
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入密码不一致")
		return
	}

	ok = u.passwordExp.MatchString(req.Password)
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8为，包含数字，字母，特殊符号")
		return
	}
	err := u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}
