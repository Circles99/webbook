package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
	"webbook/internal/repository"
	"webbook/internal/repository/cache"
	"webbook/internal/repository/dao"
	"webbook/internal/service"
	"webbook/internal/service/sms/tencent"
	"webbook/internal/web"
	"webbook/internal/web/middleware"
)

func main() {

	//db := initDB()
	//server := initWebServer()
	//u := initUser(db)
	//u.RegisterRoutes(server)

	server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好")
	})

	server.Run(":8080")

}

func initDB() *gorm.DB {

	db, err := gorm.Open(mysql.Open("root:root@tcp(webook-mysql:3309)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	repo := repository.NewUserRepository(dao.NewUserDao(db), cache.NewUserCache(rdb))
	svc := service.NewUserService(repo)

	codeRepo := repository.NewCodeRepository(cache.NewCodeCache(rdb))
	codeSvc := service.NewCodeService(codeRepo, tencent.NewService(nil, "", ""))
	u := web.NewUserHandler(svc, codeSvc)
	return u
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	_ = redis.NewClient(&redis.Options{
		Addr: "webook-redis:16379",
	})

	//server.Use(ratelimit.NewRedisSlidingWindowLimiter(r, time.Second, 100).Limit())
	server.Use(cors.New(cors.Config{
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"x-jwt-token"},
		// 是否允许携带cookie
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "你的域名.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	//store := cookie.NewStore([]byte("secret"))

	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("dsadsadsaeeeeeeqq-1"), []byte("dsadsadsaeeeeeeqq-2"))
	//if err != nil {
	//	panic(err)
	//}
	//server.Use(sessions.Sessions("mysession", store))

	// 登录校验
	//server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/login").IgnorePaths("/users/signupt").Build())
	server.Use(middleware.NewLoginJwtMiddlewareBuilder().IgnorePaths("/users/login").IgnorePaths("/users/signupt").Build())

	return server
}
