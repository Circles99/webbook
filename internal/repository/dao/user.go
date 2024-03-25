package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("邮箱冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDao interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	Insert(ctx context.Context, u User) error
	Edit(ctx context.Context, u User) error
}

type GORMUserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GORMUserDao{db: db}
}

func (dao *GORMUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

func (dao *GORMUserDao) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

func (dao *GORMUserDao) Insert(ctx context.Context, u User) error {
	// 毫秒数
	now := time.Now().UnixMilli()
	u.Updated = now
	u.Created = now

	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflicts uint16 = 1062
		if mysqlErr.Number == uniqueConflicts {
			// 邮箱冲突
			return ErrUserDuplicate
		}
	}
	return err
}

func (dao *GORMUserDao) Edit(ctx context.Context, u User) error {
	// 毫秒数
	now := time.Now().UnixMilli()
	u.Updated = now

	err := dao.db.WithContext(ctx).Updates(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflicts uint16 = 1062
		if mysqlErr.Number == uniqueConflicts {
			// 邮箱冲突
			return ErrUserDuplicate
		}
	}
	return err
}

// User 对应数据结构表， 相当于PO, 有些叫model，有些叫数据库层面的entity
type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement;id"`
	Email    sql.NullString `gorm:"unique"`
	Phone    sql.NullString `gorm:"unique"`
	Password string
	NickName string
	Birthday string
	Desc     string
	Created  int64
	Updated  int64
}
