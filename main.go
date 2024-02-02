package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))

	// 登录校验
	server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/login").IgnorePaths("/users/signupt").Build())

	return server
}
