package service

import (
	"context"
	"webbook/internal/domain"
	"webbook/internal/repository/article"
	"webbook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
}

type ArticleServiceImpl struct {
	repo article.ArticleRepository

	// v1 操作两个repository
	authRepo   article.ArticleAuthorRepository
	readerRepo article.ArticleReaderRepository
	l          logger.Logger
}

func NewArticleService(authRepo article.ArticleAuthorRepository, readerRepo article.ArticleReaderRepository, l logger.Logger) ArticleService {
	return &ArticleServiceImpl{
		authRepo:   authRepo,
		readerRepo: readerRepo,
		l:          l,
	}
}

func (a *ArticleServiceImpl) Save(ctx context.Context, art domain.Article) (int64, error) {

	art.Status = domain.ArticleStatusPublished
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Save(ctx, art)
}

func (a *ArticleServiceImpl) SaveV1(ctx context.Context, art domain.Article) (int64, error) {

	var (
		id  = art.Id
		err error
	)

	if art.Id > 0 {
		err = a.authRepo.Update(ctx, art)

	} else {
		id, err = a.authRepo.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id

	for i := 0; i < 3; i++ {
		id, err = a.readerRepo.Save(ctx, art)
		if err == nil {
			break
		}
		a.l.Error("部分失败,保存到线上库失败", logger.Int64("art_id", art.Id), logger.Error(err))
	}

	if err != nil {
		a.l.Error("部分失败,重试彻底失败 ", logger.Int64("art_id", art.Id), logger.Error(err))
	}

	return id, err
}

func (a *ArticleServiceImpl) Publish(ctx context.Context, art domain.Article) (int64, error) {

	art.Status = domain.ArticleStatusPublished

	return a.repo.Sync(ctx, art)
}

func (a *ArticleServiceImpl) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	id, err := a.authRepo.Create(ctx, art)
	if err != nil {
		return 0, err
	}
	art.Id = id
	return a.readerRepo.Save(ctx, art)
}
