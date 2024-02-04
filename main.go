package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"webbook/internal/respository"
	"webbook/internal/respository/dao"
	"webbook/internal/service"
	"webbook/internal/web"
	"webbook/internal/web/middleware"
)

func main() {

	db := initDB()
	server := initWebServer()
	u := initUser(db)
	u.RegisterRoutes(server)
	server.Run(":8080")

}

func initDB() *gorm.DB {

	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initUser(db *gorm.DB) *web.UserHandler {
	repo := respository.NewUserRepository(dao.NewUserDao(db))
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "你的域名.com")
		},
	}))

	//store := cookie.NewStore([]byte("secret"))

	store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("dsadsadsaeeeeeeqq-1"), []byte("dsadsadsaeeeeeeqq-2"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("mysession", store))

	// 登录校验
	//server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/login").IgnorePaths("/users/signupt").Build())
	server.Use(middleware.NewLoginJwtMiddlewareBuilder().IgnorePaths("/users/login").IgnorePaths("/users/signupt").Build())

	return server
}
