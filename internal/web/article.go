package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webbook/internal/domain"
	"webbook/internal/service"
	ijwt "webbook/internal/web/jwt"
	"webbook/pkg/logger"
)

var _ Handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.Logger
}

func NewArticleHandler(svc service.ArticleService, l logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (a *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", a.Edit)
	g.POST("/publish", a.Publish)
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {

	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "参数错误",
		})
		return
	}

	// 检测输入
	c := ctx.MustGet("claims")

	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, "系统错误")
		return
	}
	// 调用service
	id, err := a.svc.Save(ctx, req.toDomain(claims.UserId))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("保存失败", logger.Error(err))
	}
	// 返回结果

	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})

}

func (a *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "参数错误",
		})
		return
	}

	// 检测输入
	c := ctx.MustGet("claims")

	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, "系统错误")
		return
	}
	// 调用service
	id, err := a.svc.Publish(ctx, req.toDomain(claims.UserId))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("发表帖子失败", logger.Error(err))
	}
	// 返回结果

	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})

}

func (req ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}
