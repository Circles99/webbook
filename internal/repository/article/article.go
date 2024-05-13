package article

import (
	"context"
	"webbook/internal/domain"
	"webbook/internal/repository/dao/article"
)

type ArticleRepository interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	// Sync
	Sync(ctx context.Context, art domain.Article) (int64, error)
}

type ArticleRepositoryImpl struct {
	dao article.ArticleDAO

	// v1 操作两个dao
	authorDao article.AuthorDao
	readerDao article.ReaderDao
}

func NewArticleRepository(dao article.ArticleDAO) ArticleRepository {
	return &ArticleRepositoryImpl{
		dao: dao,
	}
}

// repository 层面上解决事务问题
func (a *ArticleRepositoryImpl) Sync(ctx context.Context, art domain.Article) (int64, error) {
	// 先保存到制作库，在保存到线上库

	var (
		id  = art.Id
		err error
	)

	if art.Id > 0 {
		err = a.authorDao.Update(ctx, a.toEntity(art))
	} else {
		id, err = a.authorDao.Insert(ctx, a.toEntity(art))
	}
	if err != nil {
		return id, err
	}
	// 操作线上库，保存数据
	// 考虑到，线上可能有 可能没有，要有一个upset的写法
	id, err = a.readerDao.Upsert(ctx, a.toEntity(art))

	return id, err
}

func (a *ArticleRepositoryImpl) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	// 先保存到制作库，在保存到线上库

	var (
		id  = art.Id
		err error
	)

	if art.Id > 0 {
		err = a.authorDao.Update(ctx, a.toEntity(art))
	} else {
		id, err = a.authorDao.Insert(ctx, a.toEntity(art))
	}
	if err != nil {
		return id, err
	}
	// 操作线上库，保存数据
	// 考虑到，线上可能有 可能没有，要有一个upset的写法
	id, err = a.readerDao.Upsert(ctx, a.toEntity(art))

	return id, err
}

func (a *ArticleRepositoryImpl) Save(ctx context.Context, art domain.Article) (int64, error) {
	return a.dao.Insert(ctx, article.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}

func (a *ArticleRepositoryImpl) Update(ctx context.Context, art domain.Article) error {
	return a.dao.Update(ctx, article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}

func (a *ArticleRepositoryImpl) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
