package service

import (
	"context"
	"webbook/internal/domain"
	"webbook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
}

type ArticleServiceImpl struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &ArticleServiceImpl{
		repo: repo,
	}
}

func (a *ArticleServiceImpl) Save(ctx context.Context, art domain.Article) (int64, error) {

	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		if err != nil {
			return 0, err
		}
		return art.Id, nil
	}

	return a.repo.Create(ctx, art)
}
