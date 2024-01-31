package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	// 毫秒数
	now := time.Now().UnixMilli()
	u.Updated = now
	u.Created = now

	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflicts uint16 = 1062
		if mysqlErr.Number == uniqueConflicts {
			// 邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}

// User 对应数据结构表， 相当于PO, 有些叫model，有些叫数据库层面的entity
type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement;id"`
	Email    string `gorm:"unique"`
	Password string
	Created  int64
	Updated  int64
}
