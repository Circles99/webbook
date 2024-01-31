package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webbook/internal/respository"
	"webbook/internal/respository/dao"
	"webbook/internal/service"
	"webbook/internal/web"
)

func main() {

	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	server := gin.Default()

	repo := respository.NewUserRepository(dao.NewUserDao(db))
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	u.RegisterRoutes(server)
	server.Run(":8080")

}
