package repository

import (
	"context"
	"webbook/internal/domain"
	"webbook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
}

type ArticleRepositoryImpl struct {
	dao dao.ArticleDAO
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &ArticleRepositoryImpl{
		dao: dao,
	}
}

func (a *ArticleRepositoryImpl) Create(ctx context.Context, art domain.Article) (int64, error) {
	return a.dao.Insert(ctx, dao.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}
