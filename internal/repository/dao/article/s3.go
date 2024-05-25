package article

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ecodeclub/ekit"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"
	"webbook/internal/domain"
)

type S3DAO struct {
	oss *s3.S3
	GORMArticleDao
	bucket *string
}

func NewOssDAO(oss *s3.S3, db *gorm.DB) ArticleDAO {
	return &S3DAO{
		oss:            nil,
		GORMArticleDao: GORMArticleDao{},
		bucket:         nil,
	}
}

func (d *S3DAO) Sync(ctx context.Context, art Article) (int64, error) {
	var (
		id = art.Id
	)
	// 先操作制作表，后操作线上表
	// gorm帮助管理了事务的声明周期
	err := d.db.Transaction(func(tx *gorm.DB) error {
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
		now := time.Now().UnixMilli()
		publishArt := PublishArticle{art}
		publishArt.Created = now
		publishArt.Updated = now
		publishArt.Content = ""
		return txDao.Upsert(ctx, PublishArticle{Article: art})
	})

	if err != nil {
		return 0, err
	}

	// 保存到oss 需要有容错机制，重试等
	_, err = d.oss.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      ekit.ToPtr[string]("webook-xxx"),
		Key:         ekit.ToPtr[string](cast.ToString(art.Id)),
		Body:        bytes.NewReader([]byte(art.Content)),
		ContentType: ekit.ToPtr[string]("text/plain;charset=utf-8"),
	})
	if err != nil {
		return 0, err
	}

	return id, err
}

func (d *S3DAO) SyncStatus(ctx context.Context, id int64, authorId int64, status uint8) error {
	now := time.Now().UnixMilli()

	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
	if err != nil {
		return err
	}

	if status == domain.ArticleStatusPrivate.ToUint8() {
		// 删除资源, 需要有容错机制，重试等
		_, err = d.oss.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: ekit.ToPtr[string]("webook-xxx"),
			Key:    ekit.ToPtr[string](cast.ToString(id)),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
