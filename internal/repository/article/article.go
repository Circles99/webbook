package article

import (
	"context"
	"webbook/internal/domain"
	"webbook/internal/repository/dao"
)

type ArticleRepository interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}

type ArticleRepositoryImpl struct {
	dao dao.ArticleDAO
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &ArticleRepositoryImpl{
		dao: dao,
	}
}

func (a *ArticleRepositoryImpl) Save(ctx context.Context, art domain.Article) (int64, error) {
	return a.dao.Insert(ctx, dao.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}

func (a *ArticleRepositoryImpl) Update(ctx context.Context, art domain.Article) error {
	return a.dao.Update(ctx, dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}
