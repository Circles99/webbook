package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
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
	return dao.db.WithContext(ctx).Create(&u).Error
}

// User 对应数据结构表， 相当于PO, 有些叫model，有些叫数据库层面的entity
type User struct {
	Id       int `gorm:"id"`
	Email    string
	Password string
	Created  int64
	Updated  int64
}
