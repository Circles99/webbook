package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
	"webbook/internal/repository/dao"
	"webbook/pkg/logger"
)

func InitDB(l logger.Logger) *gorm.DB {

	type Config struct {
		Dsn string `json:"dsn"`
	}

	var cfg = Config{
		Dsn: "root:root@tcp(localhost:13316)/webook",
	}

	// 看起来 remote 不支持key的切割
	err := viper.UnmarshalKey("db", &cfg)

	dsn := viper.GetString(cfg.Dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询日志，只有执行时间超过这个阈值才会使用
			// 50Ms 100Ms 左右设置为慢查询日志
			SlowThreshold:             time.Millisecond * 10,
			IgnoreRecordNotFoundError: true,
			// 这里true就是日志中带有占位符
			ParameterizedQueries: true,
			LogLevel:             glogger.Info,
		})})
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, field ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{
		Key:   "args",
		Value: args,
	})
}
