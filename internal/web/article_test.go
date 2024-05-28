package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"webbook/internal/domain"
	"webbook/internal/service"
	svcmocks "webbook/internal/service/mocks"
	ijwt "webbook/internal/web/jwt"
	"webbook/pkg/ginx"
	"webbook/pkg/logger"
)

func TestArticleHandler_Publish(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (service.ArticleService, logger.Logger)
		reqBody  string
		wantCode int
		wantRes  ginx.Result
	}{
		{
			name: "新建并发表",
			mock: func(ctrl *gomock.Controller) (service.ArticleService, logger.Logger) {

				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      0,
					Title:   "我的标题1号",
					Content: "我的内容1号",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return nil, &logger.ZapLogger{}
			},
			reqBody: `
{
	"title":"我的标题1号"
	"content":"我的内容1号"
}
`,
			wantCode: 200,
			wantRes: ginx.Result{
				Data: float64(1),
				Msg:  "OK",
			},
		},
		{
			name: "publish失败",
			mock: func(ctrl *gomock.Controller) (service.ArticleService, logger.Logger) {

				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      0,
					Title:   "我的标题1号",
					Content: "我的内容1号",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("publish失败"))
				return nil, &logger.ZapLogger{}
			},
			reqBody: `
{
	"title":"我的标题1号"
	"content":"我的内容1号"
}
`,
			wantCode: 200,
			wantRes: ginx.Result{
				Data: float64(5),
				Msg:  "系统错误",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("claims", &ijwt.UserClaims{
					RegisteredClaims: jwt.RegisteredClaims{},
					UserId:           123,
				})
			})

			h := NewArticleHandler(tt.mock(ctrl))
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/article/publish", bytes.NewBuffer([]byte(tt.reqBody)))

			require.NoError(t, err)
			req.Header.Set("content-type", "application/json")

			resp := httptest.NewRecorder()
			t.Log(resp)

			server.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)

			var webResult ginx.Result
			err = json.NewDecoder(resp.Body).Decode(&webResult)
			require.NoError(t, err)

			assert.Equal(t, tt.wantRes, webResult)

		})
	}
}
