package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"webbook/internal/domain"
	"webbook/internal/repository/dao/article"
	ijwt "webbook/internal/web/jwt"
	"webbook/ioc"
)

// 测试套件
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *ArticleTestSuite) SetupSuite() {
	// 所有测试执行之前，初始化一些内容
	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("claims", &ijwt.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{},
			UserId:           123,
		})
	})
	s.db = ioc.InitDB(ioc.InitLogger())
	aHdl := InitArticleHandler(article.NewArticleDao(s.db))
	aHdl.RegisterRoutes(s.server)
}

// TearDownTest 每一个都会执行
func (s *ArticleTestSuite) TearDownTest() {
	// 清空所有数据，并且自增主键恢复到1
	s.db.Exec("TRUNCATE TABLE articles")
}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name string

		// 集成测试准备数据
		before func(t *testing.T)
		// 集成测试验证数据
		after func(t *testing.T)

		art Article
		// http 响应码
		wantCode int
		// 返回带上帖子的ID
		wantRes Result[int64]
	}{
		{
			name: "保存成功，新建帖子",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art article.Article
				err := s.db.Where("id=?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Created > 0)
				assert.True(t, art.Updated > 0)
				art.Created = 0
				art.Updated = 0
				assert.Equal(t, article.Article{
					Id:       1,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
				}, art)
			},
			art: Article{
				Title:   "我的标题",
				Content: "我的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "OK",
				Data: 1,
			},
		},
		{
			name: "修改已有帖子，并保存",
			before: func(t *testing.T) {
				err := s.db.Create(article.Article{Id: 2, Title: "我的标题", Content: "我的内容", AuthorId: 123, Created: 123, Updated: 234}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := s.db.Where("id=?").First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Updated > 234)

				art.Updated = 0
				assert.Equal(t, article.Article{
					Created:  123,
					Id:       2,
					Title:    "新的标题",
					Content:  "新的内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
				}, art)
			},
			art: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "OK",
				Data: 2,
			},
		},
		{
			name: "修改别人的帖子",
			before: func(t *testing.T) {
				err := s.db.Create(article.Article{Id: 2, Title: "我的标题", Content: "我的内容", AuthorId: 789, Created: 123, Updated: 234}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := s.db.Where("id=?", 2).First(&art).Error
				assert.NoError(t, err)

				assert.Equal(t, article.Article{
					Created:  123,
					Id:       3,
					Title:    "新的标题",
					Content:  "新的内容",
					Updated:  234,
					AuthorId: 789,
				}, art)
			},
			art: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "OK",
				Data: 3,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)

			reqBody, err := json.Marshal(tc.art)

			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("content-type", "application/json")

			resp := httptest.NewRecorder()
			t.Log(resp)

			s.server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)

			if resp.Code != 200 {
				return
			}

			var webResult Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webResult)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webResult)

			tc.after(t)
		})
	}
}

func (s *ArticleTestSuite) TestAbc() {
	s.T().Log("hello， 这是一个套件")
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
