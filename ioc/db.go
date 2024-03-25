package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webbook/internal/repository/dao"
)

func InitDB() *gorm.DB {
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
