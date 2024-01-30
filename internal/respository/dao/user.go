package dao

import "context"

type UserDao struct {
}

func (u *UserDao) Insert(ctx context.Context) error {
	return nil
}

type User struct {
	Id       int `gorm:"id"`
	Email    string
	Password string
	Created  int64
	Updated  int64
}
