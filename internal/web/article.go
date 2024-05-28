package web

import (
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"webbook/internal/domain"
	"webbook/internal/service"
	ijwt "webbook/internal/web/jwt"
	"webbook/pkg/ginx"
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
	g.POST("/withdraw", a.Withdraw)
	g.POST("/list", ginx.WrapBodyAndToken[ListReq, ijwt.UserClaims](a.List))
}

func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
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
	err := a.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claims.UserId,
		},
	})

	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("保存失败", logger.Error(err))
	}

	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "OK",
	})

}

func (a *ArticleHandler) Edit(ctx *gin.Context) {

	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
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
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("保存失败", logger.Error(err))
	}
	// 返回结果

	ctx.JSON(http.StatusOK, ginx.Result{
		Msg:  "OK",
		Data: id,
	})

}

func (a *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
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
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("发表帖子失败", logger.Error(err))
	}
	// 返回结果

	ctx.JSON(http.StatusOK, ginx.Result{
		Msg:  "OK",
		Data: id,
	})

}

func (a *ArticleHandler) List(ctx *gin.Context, req ListReq, uc ijwt.UserClaims) (ginx.Result, error) {
	res, err := a.svc.List(ctx, uc.UserId, req.Offset, req.Limit)
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, nil
	}
	return ginx.Result{
		Code: 2,
		Msg:  "",
		Data: slice.Map[domain.Article, ArticleVo](res, func(idx int, src domain.Article) ArticleVo {
			return ArticleVo{
				Id:       src.Id,
				Title:    src.Title,
				Abstract: src.Abstract(),
				//AuthorId:   src.Author.Id,
				//AuthorName: src.Author.Name,
				Status:  src.Status.ToUint8(),
				Created: src.Created.Format(time.DateTime),
				Updated: src.Updated.Format(time.DateTime),
			}
		}),
	}, nil

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
