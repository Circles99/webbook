package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, art Article) error
}

type GORMArticleDao struct {
	db *gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDAO {
	return &GORMArticleDao{
		db: db,
	}
}

func (dao *GORMArticleDao) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Created = now
	art.Updated = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (dao *GORMArticleDao) Update(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Updated = now
	err := dao.db.WithContext(ctx).Model(&art).Where("id=?", art.Id).Updates(map[string]any{
		"title":   art.Title,
		"content": art.Content,
		"updated": art.Updated,
	}).Error
	return err
}

// Article 制作库的
type Article struct {
	Id       int64  `gorm:"primaryKey;autoIncrement"`
	Title    string `gorm:"type=varchar(1024)"`
	Content  string `gorm:"type=BLOB"`
	AuthorId int64  `gorm:"index"`

	//AuthorId int64  `gorm:"index=aid_ctime"`
	//Created  int64  `gorm:"index=aid_ctime"`
	Created int64
	Updated int64
}
