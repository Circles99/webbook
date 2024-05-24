package article

import (
	"context"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	Upsert(ctx context.Context, art PublishArticle) error
	SyncStatus(ctx context.Context, id int64, authorId int64, status uint8) error
}

// Article 制作库的
type Article struct {
	Id       int64  `gorm:"primaryKey;autoIncrement" bson:"id"`
	Title    string `gorm:"type=varchar(1024)" bson:"title"`
	Content  string `gorm:"type=BLOB" bson:"content"`
	AuthorId int64  `gorm:"index" bson:"authorId"`

	Status uint8

	//AuthorId int64  `gorm:"index=aid_ctime"`
	//Created  int64  `gorm:"index=aid_ctime"`
	Created int64
	Updated int64
}
