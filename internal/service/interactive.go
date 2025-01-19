package service

import (
	"context"
	"golang.org/x/sync/errgroup"
	"webbook/internal/domain"
	"webbook/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	// Like 点赞
	IncrLike(ctx context.Context, biz string, bizId int64, uid int64) error
	// CancelLike 取消点赞
	CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error
	// Collect 收藏
	Collection(ctx context.Context, biz string, bizId, cid, uid int64) error
	Get(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository
}

func (i interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizId)
}

func (i interactiveService) IncrLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	return i.repo.IncrLike(ctx, biz, bizId, uid)
}

func (i interactiveService) CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	return i.repo.CancelLike(ctx, biz, bizId, uid)
}

func (i interactiveService) Collection(ctx context.Context, biz string, bizId, cid, uid int64) error {
	// service层面上还叫收藏
	// repository 层面上就应该是增加一个项
	return i.repo.AddCollectionItem(ctx, biz, bizId, cid, uid)
}

func (i interactiveService) Get(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error) {
	var (
		eg        errgroup.Group
		intr      domain.Interactive
		liked     bool
		collected bool
	)
	eg.Go(func() error {
		var err error
		intr, err = i.repo.Get(ctx, biz, bizId, uid)
		return err
	})

	eg.Go(func() error {
		var err error
		liked, err = i.repo.Liked(ctx, biz, bizId, uid)
		return err
	})

	eg.Go(func() error {
		var err error
		liked, err = i.repo.Collected(ctx, biz, bizId, uid)
		return err
	})

	err := eg.Wait()
	if err != nil {
		return domain.Interactive{}, err
	}
	intr.Liked = liked
	intr.Collected = collected
	return intr, nil
}
