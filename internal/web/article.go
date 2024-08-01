package web

import (
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
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
	svc     service.ArticleService
	l       logger.Logger
	intrSvc service.InteractiveService
	biz     string
}

func NewArticleHandler(svc service.ArticleService, l logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
		biz: "article",
	}
}

func (a *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", a.Edit)
	g.POST("/publish", a.Publish)
	g.POST("/withdraw", a.Withdraw)
	g.GET("/list", ginx.WrapBodyAndToken[ListReq, ijwt.UserClaims](a.List))
	g.GET("/detail/:id", ginx.WrapToken[ijwt.UserClaims](a.Detail))
	g.GET("/pubDetail/:id", ginx.WrapToken[ijwt.UserClaims](a.PubDetail))
	g.POST("/like", ginx.WrapBodyAndToken[LikeReq, ijwt.UserClaims](a.Like)) // 点赞和取消点赞都是这个接口
}

func (a *ArticleHandler) Like(ctx *gin.Context, req LikeReq, uc ijwt.UserClaims) (ginx.Result, error) {
	var err error
	if req.Like {
		err = a.intrSvc.IncrLike(ctx, a.biz, req.Id, uc.UserId)
	} else {

		err = a.intrSvc.CancelLike(ctx, a.biz, req.Id, uc.UserId)
	}

	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}

	return ginx.Result{Msg: "ok"}, nil
}

func (a *ArticleHandler) PubDetail(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	id := cast.ToInt64(ctx.Param("id"))
	res, err := a.svc.PubDetail(ctx, id)
	if err != nil {
		return ginx.Result{Code: 5, Msg: "数据错误"}, nil
	}
	if res.Author.Id != uc.UserId {
		return ginx.Result{Code: 5, Msg: "作者错误"}, nil
	}

	// 增加阅读计数
	go func() {
		er := a.intrSvc.IncrReadCnt(ctx, a.biz, res.Id)
		if er != nil {
			a.l.Error("增加阅读计数失败", logger.Int64("aid", res.Id), logger.Error(er))
		}
	}()
	return ginx.Result{Data: ArticleVo{
		Id:         res.Id,
		Title:      res.Title,
		Abstract:   res.Abstract(),
		Content:    res.Content,
		AuthorId:   res.Author.Id,
		AuthorName: res.Author.Name,
		Status:     res.Status.ToUint8(),
		Created:    res.Created.Format(time.DateTime),
		Updated:    res.Updated.Format(time.DateTime),
	}}, nil
}

func (a *ArticleHandler) Detail(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	id := cast.ToInt64(ctx.Param("id"))
	res, err := a.svc.Detail(ctx, id)
	if err != nil {
		return ginx.Result{Code: 5, Msg: "数据错误"}, nil
	}
	if res.Author.Id != uc.UserId {
		return ginx.Result{Code: 5, Msg: "作者错误"}, nil
	}

	return ginx.Result{Data: ArticleVo{
		Id:         res.Id,
		Title:      res.Title,
		Abstract:   res.Abstract(),
		Content:    res.Content,
		AuthorId:   res.Author.Id,
		AuthorName: res.Author.Name,
		Status:     res.Status.ToUint8(),
		Created:    res.Created.Format(time.DateTime),
		Updated:    res.Updated.Format(time.DateTime),
	}}, nil
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
