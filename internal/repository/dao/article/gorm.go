package article

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GORMArticleDao struct {
	db *gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDAO {
	return &GORMArticleDao{
		db: db,
	}
}

func (dao *GORMArticleDao) SyncStatus(ctx context.Context, id int64, authorId int64, status uint8) error {
	now := time.Now().UnixMilli()

	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).Where("id = ? AND author_id = ?", id, authorId).Updates(map[string]any{
			"status":  status,
			"updated": now,
		})
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return errors.New("更新失败")
		}

		res = tx.Model(&PublishArticle{}).Where("id = ? AND author_id = ?", id, authorId).Updates(map[string]any{
			"status":  status,
			"updated": now,
		})

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return errors.New("更新失败")
		}

		return nil
	})

}

func (dao *GORMArticleDao) Transaction(ctx context.Context, bizFunc func(txDao ArticleDAO) error) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDao := NewArticleDao(tx)
		return bizFunc(txDao)
	})
}

func (dao *GORMArticleDao) Upsert(ctx context.Context, art PublishArticle) error {
	now := time.Now().UnixMilli()
	art.Created = now
	art.Updated = now
	// 相当于 inset into xxx Values xxx ON DUPLICATE KEY UPDATE
	err := dao.db.Clauses(clause.OnConflict{
		// mysql 只关心DoUpdates
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
			"updated": art.Updated,
			"status":  art.Status,
		}),
	}).Create(&art).Error
	return err
}

func (dao *GORMArticleDao) Sync(ctx context.Context, art Article) (int64, error) {

	var (
		id = art.Id
	)
	// 先操作制作表，后操作线上表
	// gorm帮助管理了事务的声明周期
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		var err error
		txDao := NewArticleDao(tx)
		if id > 0 {
			err = txDao.Update(ctx, art)
		} else {
			id, err = txDao.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		// 操作线上库，保存数据

		art.Id = id
		return txDao.Upsert(ctx, PublishArticle{Article: art})
	})
	return id, err
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
	res := dao.db.WithContext(ctx).Model(&art).Where("id=? AND author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"status":  art.Status,
			"updated": art.Updated,
		})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("更新失败")
	}

	return nil
}
