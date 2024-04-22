package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webbook/internal/repository/dao"
)

func InitDB() *gorm.DB {

	type Config struct {
		Dsn string `json:"dsn"`
	}

	var cfg = Config{
		Dsn: "root:root@tcp(localhost:13316)/webook",
	}

	// 看起来 remote 不支持key的切割
	err := viper.UnmarshalKey("db", &cfg)

	dsn := viper.GetString(cfg.Dsn)
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
