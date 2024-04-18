package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webbook/internal/repository/dao"
)

func InitDB() *gorm.DB {
	dsn := viper.GetString("db.mysql.dsn")
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
